package main

import (
	"errors"
	"github.com/HcashOrg/hcd/addrmgr"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/connmgr"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcd/hcutil/bloom"
	"github.com/HcashOrg/hcd/peer"
	"github.com/HcashOrg/hcd/wire"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// serverPeer extends the peer to maintain state shared by the server and
// the blockmanager.
type serverWitnessPeer struct {
	*peer.WitnessPeer

	connReq        *connmgr.ConnReq
	server         *server
	persistent     bool
	continueHash   *chainhash.Hash
	relayMtx       sync.Mutex
	disableRelayTx bool
	isWhitelisted  bool
	requestQueue   []*wire.InvVect
	requestedTxns  map[chainhash.Hash]struct{}
	filter         *bloom.Filter
	knownAddresses map[string]struct{}
	banScore       connmgr.DynamicBanScore
	quit           chan struct{}
	// It is used to prevent more than one response per connection.
	addrsSent bool
	// addrsSent  track whether or not the peer
	// has already sent the respective request.  It is used to prevent more
	// than one response per connection.

	// The following chans are used to sync blockmanager and server.
	txProcessed chan struct{}
}

// Only respond with addresses once per connection
//if sp.addrsSent {
//	peerLog.Tracef("Ignoring getaddr from %v - already sent", sp.Peer)
//	return
//}

//sp.addrsSent = true

// newserverWitnessPeer returns a new serverWitnessPeer instance. The peer needs to be set by
// the caller.
func newserverWitnessPeer(s *server, isPersistent bool) *serverWitnessPeer {
	return &serverWitnessPeer{
		server:         s,
		persistent:     isPersistent,
		requestedTxns:  make(map[chainhash.Hash]struct{}),
		filter:         bloom.LoadFilter(nil),
		knownAddresses: make(map[string]struct{}),
		quit:           make(chan struct{}),
		txProcessed:    make(chan struct{}, 1),
	}
}

type broadcastWitnessMsg struct {
	message      wire.Message
	excludePeers []*serverWitnessPeer
}
// addKnownAddresses adds the given addresses to the set of known addreses to
// the peer to prevent sending duplicate addresses.
func (sp *serverWitnessPeer) addKnownAddresses(addresses []*wire.NetAddress) {
	for _, na := range addresses {
		sp.knownAddresses[addrmgr.NetAddressKey(na)] = struct{}{}
	}
}

// addressKnown true if the given address is already known to the peer.
func (sp *serverWitnessPeer) addressKnown(na *wire.NetAddress) bool {
	_, exists := sp.knownAddresses[addrmgr.NetAddressKey(na)]
	return exists
}

// setDisableRelayTx toggles relaying of transactions for the given peer.
// It is safe for concurrent access.
func (sp *serverWitnessPeer) setDisableRelayTx(disable bool) {
	sp.relayMtx.Lock()
	sp.disableRelayTx = disable
	sp.relayMtx.Unlock()
}

// relayTxDisabled returns whether or not relaying of transactions for the given
// peer is disabled.
// It is safe for concurrent access.
func (sp *serverWitnessPeer) relayTxDisabled() bool {
	sp.relayMtx.Lock()
	isDisabled := sp.disableRelayTx
	sp.relayMtx.Unlock()

	return isDisabled
}

// pushAddrMsg sends an addr message to the connected peer using the provided
// addresses.
func (sp *serverWitnessPeer) pushAddrMsg(addresses []*wire.NetAddress) {
	// Filter addresses already known to the peer.
	addrs := make([]*wire.NetAddress, 0, len(addresses))
	for _, addr := range addresses {
		if !sp.addressKnown(addr) {
			addrs = append(addrs, addr)
		}
	}
	known, err := sp.PushAddrMsg(addrs)
	if err != nil {
		peerLog.Errorf("Can't push address message to %s: %v", sp.WitnessPeer, err)
		sp.Disconnect()
		return
	}
	sp.addKnownAddresses(known)
}

// addBanScore increases the persistent and decaying ban score fields by the
// values passed as parameters. If the resulting score exceeds half of the ban
// threshold, a warning is logged including the reason provided. Further, if
// the score is above the ban threshold, the peer will be banned and
// disconnected.
func (sp *serverWitnessPeer) addBanScore(persistent, transient uint32, reason string) {
	// No warning is logged and no score is calculated if banning is disabled.
	if cfg.DisableBanning {
		return
	}
	if sp.isWhitelisted {
		peerLog.Debugf("Misbehaving whitelisted peer %s: %s", sp, reason)
		return
	}

	warnThreshold := cfg.BanThreshold >> 1
	if transient == 0 && persistent == 0 {
		// The score is not being increased, but a warning message is still
		// logged if the score is above the warn threshold.
		score := sp.banScore.Int()
		if score > warnThreshold {
			peerLog.Warnf("Misbehaving peer %s: %s -- ban score is %d, "+
				"it was not increased this time", sp, reason, score)
		}
		return
	}
	score := sp.banScore.Increase(persistent, transient)
	if score > warnThreshold {
		peerLog.Warnf("Misbehaving peer %s: %s -- ban score increased to %d",
			sp, reason, score)
		if score > cfg.BanThreshold {
			peerLog.Warnf("Misbehaving peer %s -- banning and disconnecting",
				sp)
			sp.server.BanWitnessPeer(sp)
			sp.Disconnect()
		}
	}
}

// OnVersion is invoked when a peer receives a version wire message and is used
// to negotiate the protocol version details as well as kick start the
// communications.
func (sp *serverWitnessPeer) OnVersion(p *peer.Peer, msg *wire.MsgVersion) {
	// Update the address manager with the advertised services for outbound
	// connections in case they have changed.  This is not done for inbound
	// connections to help prevent malicious behavior and is skipped when
	// running on the simulation test network since it is only intended to
	// connect to specified peers and actively avoids advertising and
	// connecting to discovered peers.
	//
	// NOTE: This is done before rejecting peers that are too old to ensure
	// it is updated regardless in the case a new minimum protocol version is
	// enforced and the remote node has not upgraded yet.
	addrManager := sp.server.addrManager
	isInbound := sp.Inbound()
	remoteAddr := sp.NA()
	if !cfg.SimNet && !isInbound {
		addrManager.SetServices(remoteAddr, msg.Services)
	}
	// Ignore peers that have a protcol version that is too old.  The peer
	// negotiation logic will disconnect it after this callback returns.
	if msg.ProtocolVersion < int32(wire.InitialProcotolVersion) {
		return
	}
	// Add the remote peer time as a sample for creating an offset against
	// the local clock to keep the network time in sync.
	sp.server.timeSource.AddTimeSample(p.Addr(), msg.Timestamp)


	//format example /hcd:2.0.0/
	var valid = regexp.MustCompile("hcd:[0-9]*.[0-9]*.[0-9]*")
	val := valid.FindAllStringSubmatch(p.UserAgent(), 1)
	if !(len(val) != 0 && len(val[0]) != 0) {
		peerLog.Warnf("peer has no hcd agentVersion %s ", sp)
		sp.server.BanWitnessPeer(sp)
		sp.Disconnect()
		return
	}

	receiveVerisonStr := strings.TrimLeft(val[0][0], "hcd:")
	versionArray := strings.Split(receiveVerisonStr, ".")
	if len(versionArray) != 3 {
		peerLog.Warnf("parser remote app version %s fail", sp)
		sp.server.BanWitnessPeer(sp)
		sp.Disconnect()
		return
	}

	oldAppMajor, err := strconv.ParseInt(versionArray[0], 10, 32)
	if err != nil {
		peerLog.Warnf("parser remote app version %s fail", sp)
		sp.server.BanWitnessPeer(sp)
		sp.Disconnect()
		return
	}
	oldAppMinor, err := strconv.ParseInt(versionArray[1], 10, 32)
	if err != nil {
		peerLog.Warnf("parser remote app version %s fail", sp)
		sp.server.BanWitnessPeer(sp)
		sp.Disconnect()
		return
	}
	oldAppPatch, err := strconv.ParseInt(versionArray[2], 10, 32)
	if err != nil {
		peerLog.Warnf("parser remote app version %s fail", sp)
		sp.server.BanWitnessPeer(sp)
		sp.Disconnect()
		return
	}

	oldVersion := int32(1000000*oldAppMajor + 10000*oldAppMinor + 100*oldAppPatch)
	currVersion := int32(1000000*appMajor + 10000*appMinor + 100*appPatch)

	if oldVersion < currVersion {
		peerLog.Warnf("too old version peer %s ", sp)
		sp.server.BanWitnessPeer(sp)
		sp.Disconnect()
		return
	}

	// Choose whether or not to relay transactions before a filter command
	// is received.
	sp.setDisableRelayTx(msg.DisableRelayTx)

	// Update the address manager and request known addresses from the
	// remote peer for outbound connections.  This is skipped when running
	// on the simulation test network since it is only intended to connect
	// to specified peers and actively avoids advertising and connecting to
	// discovered peers.
	if !cfg.SimNet {
		//addrManager := sp.server.addrManager
		// Outbound connections.
		if !p.Inbound() {
			// TODO(davec): Only do this if not doing the initial block
			// download and the local address is routable.
			//if !cfg.DisableListen /* && isCurrent? */ {
			if !cfg.DisableListen && sp.server.blockManager.IsCurrent() {
				// Get address that best matches.
				lna := addrManager.GetBestLocalAddress(remoteAddr)
				if addrmgr.IsRoutable(lna) {
					// Filter addresses the peer already knows about.
					addresses := []*wire.NetAddress{lna}
					sp.pushAddrMsg(addresses)
				}
			}

			// Request known addresses if the server address manager
			// needs more.
			if addrManager.NeedMoreAddresses() {
				p.QueueMessage(wire.NewMsgGetAddr(), nil)
			}

			// Mark the address as a known good address.
			addrManager.Good(remoteAddr)
		}
	}

	// Add valid peer to the server.
	sp.server.AddWitnessPeer(sp)
}

// OnTx is invoked when a peer receives a tx wire message.  It blocks until the
// transaction has been fully processed.  Unlock the block handler this does not
// serialize all transactions through a single thread transactions don't rely on
// the previous one in a linear fashion like blocks.
func (sp *serverWitnessPeer) OnWitnessTx(p *peer.Peer, msg *wire.MsgTx) {

	// Add the transaction to the known inventory for the peer.
	// Convert the raw MsgTx to a hcutil.Tx which provides some convenience
	// methods and things such as hash caching.
	tx := hcutil.NewTx(msg)
	iv := wire.NewInvVect(wire.InvTypeTx, tx.Hash())
	p.AddKnownInventory(iv)

	// Queue the transaction up to be handled by the block manager and
	// intentionally block further receives until the transaction is fully
	// processed and known good or bad.  This helps prevent a malicious peer
	// from queuing up a bunch of bad transactions before disconnecting (or
	// being disconnected) and wasting memory.
	sp.server.blockManager.QueueWitnessTx(tx, sp)
	<-sp.txProcessed
}

// OnInv is invoked when a peer receives an inv wire message and is used to
// examine the inventory being advertised by the remote peer and react
// accordingly.  We pass the message down to blockmanager which will call
// QueueMessage with any appropriate responses.
func (sp *serverWitnessPeer) OnWitnessInv(p *peer.Peer, msg *wire.MsgInv) {

	if len(msg.InvList) > 0 {
		sp.server.blockManager.QueueWitnessInv(msg, sp)
	}
	return

}

// handleGetData is invoked when a peer receives a getdata wire message and is
// used to deliver block and transaction information.
func (sp *serverWitnessPeer) OnGetData(p *peer.Peer, msg *wire.MsgGetData) {
	// Ignore empty getdata messages.
	if len(msg.InvList) == 0 {
		return
	}

	numAdded := 0
	notFound := wire.NewMsgNotFound()

	length := len(msg.InvList)
	// A decaying ban score increase is applied to prevent exhausting resources
	// with unusually large inventory queries.
	// Requesting more than the maximum inventory vector length within a short
	// period of time yields a score above the default ban threshold. Sustained
	// bursts of small requests are not penalized as that would potentially ban
	// peers performing IBD.
	// This incremental score decays each minute to half of its value.
	sp.addBanScore(0, uint32(length)*99/wire.MaxInvPerMsg, "getdata")

	// We wait on this wait channel periodically to prevent queuing
	// far more data than we can send in a reasonable time, wasting memory.
	// The waiting occurs after the database fetch for the next one to
	// provide a little pipelining.
	var waitChan chan struct{}
	doneChan := make(chan struct{}, 1)

	for i, iv := range msg.InvList {
		var c chan struct{}
		// If this will be the last message we send.
		if i == length-1 && len(notFound.InvList) == 0 {
			c = doneChan
		} else if (i+1)%3 == 0 {
			// Buffered so as to not make the send goroutine block.
			c = make(chan struct{}, 1)
		}
		var err error
		switch iv.Type {
		case wire.InvTypeTx:
			err = sp.server.pushWitnessTxMsg(sp, &iv.Hash, c, waitChan)
		default:
			peerLog.Warnf("Unknown type in inventory request %d",
				iv.Type)
			continue
		}
		if err != nil {
			notFound.AddInvVect(iv)

			// When there is a failure fetching the final entry
			// and the done channel was sent in due to there
			// being no outstanding not found inventory, consume
			// it here because there is now not found inventory
			// that will use the channel momentarily.
			if i == len(msg.InvList)-1 && c != nil {
				<-c
			}
		}
		numAdded++
		waitChan = c
	}
	if len(notFound.InvList) != 0 {
		p.QueueMessage(notFound, doneChan)
	}

	// Wait for messages to be sent. We can send quite a lot of data at this
	// point and this will keep the peer busy for a decent amount of time.
	// We don't process anything else by them in this time so that we
	// have an idea of when we should hear back from them - else the idle
	// timeout could fire when we were only half done sending the blocks.
	if numAdded > 0 {
		<-doneChan
	}
}



// AddRebroadcastInventory adds 'iv' to the list of inventories to be
// rebroadcasted at random intervals until they show up in a block.
func (s *server) AddRebroadcastWitnessInventory(iv *wire.InvVect, data interface{}) {
	// Ignore if shutting down.
	if atomic.LoadInt32(&s.shutdown) != 0 {
		return
	}

	s.modifyRebroadcastWitnessInv <- broadcastInventoryAdd{invVect: iv, data: data}
}

// RemoveRebroadcastInventory removes 'iv' from the list of items to be
// rebroadcasted if present.
func (s *server) RemoveRebroadcastWitnessInventory(iv *wire.InvVect) {
	// Ignore if shutting down.
	if atomic.LoadInt32(&s.shutdown) != 0 {
		return
	}

	s.modifyRebroadcastWitnessInv <- broadcastInventoryDel(iv)
}





func (s *server) pushWitnessTxMsg(sp *serverWitnessPeer, hash *chainhash.Hash, doneChan chan<- struct{}, waitChan <-chan struct{}) error {
	// Attempt to fetch the requested transaction from the pool.  A
	// call could be made to check for existence first, but simply trying
	// to fetch a missing transaction results in the same behavior.
	// Do not allow peers to request transactions already in a block
	// but are unconfirmed, as they may be expensive. Restrict that
	// to the authenticated RPC only.
	tx, err := s.txMemPool.FetchTransaction(hash, false)

	if err != nil {
		peerLog.Tracef("Unable to fetch tx %v from transaction "+
			"pool: %v", hash, err)

		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}

	// Once we have fetched data wait for any previous operation to finish.
	if waitChan != nil {
		<-waitChan
	}

	sp.QueueMessage(tx.MsgTx(), doneChan)

	return nil
}




// handleAddPeerMsg deals with adding new peers.  It is invoked from the
// peerHandler goroutine.
func (s *server) handleAddWitnessPeerMsg(state *witnessPeerState, sp *serverWitnessPeer) bool {
	if sp == nil {
		return false
	}

	// Ignore new peers if we're shutting down.
	if atomic.LoadInt32(&s.shutdown) != 0 {
		srvrLog.Infof("New peer %s ignored - server is shutting down", sp)
		sp.Disconnect()
		return false
	}

	// Disconnect banned peers.
	host, _, err := net.SplitHostPort(sp.Addr())
	if err != nil {
		srvrLog.Debugf("can't split hostport %v", err)
		sp.Disconnect()
		return false
	}
	if banEnd, ok := state.banned[host]; ok {
		if time.Now().Before(banEnd) {
			srvrLog.Debugf("Peer %s is banned for another %v - disconnecting",
				host, banEnd.Sub(time.Now()))
			sp.Disconnect()
			return false
		}

		srvrLog.Infof("Peer %s is no longer banned", host)
		delete(state.banned, host)
	}

	// TODO: Check for max peers from a single IP.

	// Limit max number of total peers.
	// allow whitelisted inbound peers regardless.
	if state.Count() >= cfg.MaxPeers && !(sp.Inbound() && sp.isWhitelisted) {
		srvrLog.Infof("Max peers reached [%d] - disconnecting peer %s",
			cfg.MaxPeers, sp)
		sp.Disconnect()
		// TODO(oga) how to handle permanent peers here?
		// they should be rescheduled.
		return false
	}

	// Add the new peer and start it.
	srvrLog.Debugf("New peer %s", sp)
	if sp.Inbound() {
		state.inboundPeers[sp.ID()] = sp
	} else {
		state.outboundGroups[addrmgr.GroupKey(sp.NA())]++
		if sp.persistent {
			state.persistentPeers[sp.ID()] = sp
		} else {
			state.outboundPeers[sp.ID()] = sp
		}
	}

	return true
}

// handleDonePeerMsg deals with peers that have signalled they are done.  It is
// invoked from the peerHandler goroutine.
func (s *server) handleDoneWitnessPeerMsg(state *witnessPeerState, sp *serverWitnessPeer) {
	var list map[int32]*serverWitnessPeer
	if sp.persistent {
		list = state.persistentPeers
	} else if sp.Inbound() {
		list = state.inboundPeers
	} else {
		list = state.outboundPeers
	}
	if _, ok := list[sp.ID()]; ok {
		if !sp.Inbound() && sp.VersionKnown() {
			state.outboundGroups[addrmgr.GroupKey(sp.NA())]--
		}
		if !sp.Inbound() && sp.connReq != nil {
			s.witnessConnManager.Disconnect(sp.connReq.ID())
		}
		delete(list, sp.ID())
		srvrLog.Debugf("Removed peer %s", sp)
		return
	}

	if sp.connReq != nil {
		s.witnessConnManager.Disconnect(sp.connReq.ID())
	}

	// Update the address' last seen time if the peer has acknowledged
	// our version and has sent us its version as well.
	if sp.VerAckReceived() && sp.VersionKnown() && sp.NA() != nil {
		s.witnessAddrManager.Connected(sp.NA())
	}

	// If we get here it means that either we didn't know about the peer
	// or we purposefully deleted it.
}

// handleBanPeerMsg deals with banning peers.  It is invoked from the
// peerHandler goroutine.
func (s *server) handleBanWitnessPeerMsg(state *witnessPeerState, sp *serverWitnessPeer) {
	host, _, err := net.SplitHostPort(sp.Addr())
	if err != nil {
		srvrLog.Debugf("can't split ban peer %s %v", sp.Addr(), err)
		return
	}
	direction := directionString(sp.Inbound())
	srvrLog.Infof("Banned peer %s (%s) for %v", host, direction,
		cfg.BanDuration)
	state.banned[host] = time.Now().Add(cfg.BanDuration)
}

// handleRelayInvMsg deals with relaying inventory to peers that are not already
// known to have it.  It is invoked from the peerHandler goroutine.
func (s *server) handleRelayWitnessInvMsg(state *witnessPeerState, msg relayMsg) {
	state.forAllPeers(func(sp *serverWitnessPeer) {
		if !sp.Connected() {
			return
		}


		if msg.invVect.Type == wire.InvTypeTx {
			// Don't relay the transaction to the peer when it has
			// transaction relaying disabled.
			if sp.relayTxDisabled() {
				return
			}
			// Don't relay the transaction if there is a bloom
			// filter loaded and the transaction doesn't match it.
			if sp.filter.IsLoaded() {
				tx, ok := msg.data.(*hcutil.Tx)
				if !ok {
					peerLog.Warnf("Underlying data for tx" +
						" inv relay is not a transaction")
					return
				}

				if !sp.filter.MatchTxAndUpdate(tx) {
					return
				}
			}
		}

		// Queue the inventory to be relayed with the next batch.
		// It will be ignored if the peer is already known to
		// have the inventory.
		sp.QueueInventory(msg.invVect)
	})
}

// handlebroadcastWitnessMsg deals with broadcasting messages to peers.  It is invoked
// from the peerHandler goroutine.
func (s *server) handleBroadcastWitnessMsg(state *witnessPeerState, bmsg *broadcastWitnessMsg) {
	state.forAllPeers(func(sp *serverWitnessPeer) {
		if !sp.Connected() {
			return
		}

		for _, ep := range bmsg.excludePeers {
			if sp == ep {
				return
			}
		}

		sp.QueueMessage(bmsg.message, nil)
	})
}



type getWitnessPeersMsg struct {
	reply chan []*serverWitnessPeer
}



type getAddedWitnessNodesMsg struct {
	reply chan []*serverWitnessPeer
}

type disconnectWitnessNodeMsg struct {
	cmp   func(*serverWitnessPeer) bool
	reply chan error
}

type connectWitnessNodeMsg struct {
	addr      string
	permanent bool
	reply     chan error
}

type removeWitnessNodeMsg struct {
	cmp   func(*serverWitnessPeer) bool
	reply chan error
}

// handleQuery is the central handler for all queries and commands from other
// goroutines related to peer state.
func (s *server) handleQueryWitness(state *witnessPeerState, querymsg interface{}) {
	switch msg := querymsg.(type) {
	case getConnCountMsg:
		nconnected := int32(0)
		state.forAllPeers(func(sp *serverWitnessPeer) {
			if sp.Connected() {
				nconnected++
			}
		})
		msg.reply <- nconnected

	case getWitnessPeersMsg:
		peers := make([]*serverWitnessPeer, 0, state.Count())
		state.forAllPeers(func(sp *serverWitnessPeer) {
			if !sp.Connected() {
				return
			}
			peers = append(peers, sp)
		})
		msg.reply <- peers

	case connectWitnessNodeMsg:
		// XXX(oga) duplicate oneshots?
		// Limit max number of total peers.
		if state.Count() >= cfg.MaxPeers {
			msg.reply <- errors.New("max peers reached")
			return
		}
		for _, peer := range state.persistentPeers {
			if peer.Addr() == msg.addr {
				if msg.permanent {
					msg.reply <- errors.New("peer already connected")
				} else {
					msg.reply <- errors.New("peer exists as a permanent peer")
				}
				return
			}
		}

		netAddr, err := addrStringToNetAddr(msg.addr)
		if err != nil {
			msg.reply <- err
			return
		}

		// TODO(oga) if too many, nuke a non-perm peer.
		go s.witnessConnManager.Connect(&connmgr.ConnReq{
			Addr:      netAddr,
			Permanent: msg.permanent,
		})
		msg.reply <- nil
	case removeWitnessNodeMsg:
		found := disconnectWitnessPeer(state.persistentPeers, msg.cmp, func(sp *serverWitnessPeer) {
			// Keep group counts ok since we remove from
			// the list now.
			state.outboundGroups[addrmgr.GroupKey(sp.NA())]--
		})

		if found {
			msg.reply <- nil
		} else {
			msg.reply <- errors.New("witness peer not found")
		}
	case getOutboundGroup:
		count, ok := state.outboundGroups[msg.key]
		if ok {
			msg.reply <- count
		} else {
			msg.reply <- 0
		}
		// Request a list of the persistent (added) peers.
	case getAddedWitnessNodesMsg:
		// Respond with a slice of the relavent peers.
		peers := make([]*serverWitnessPeer, 0, len(state.persistentPeers))
		for _, sp := range state.persistentPeers {
			peers = append(peers, sp)
		}
		msg.reply <- peers
	case disconnectWitnessNodeMsg:
		// Check inbound peers. We pass a nil callback since we don't
		// require any additional actions on disconnect for inbound peers.
		found := disconnectWitnessPeer(state.inboundPeers, msg.cmp, nil)
		if found {
			msg.reply <- nil
			return
		}

		// Check outbound peers.
		found = disconnectWitnessPeer(state.outboundPeers, msg.cmp, func(sp *serverWitnessPeer) {
			// Keep group counts ok since we remove from
			// the list now.
			state.outboundGroups[addrmgr.GroupKey(sp.NA())]--
		})
		if found {
			// If there are multiple outbound connections to the same
			// ip:port, continue disconnecting them all until no such
			// peers are found.
			for found {
				found = disconnectWitnessPeer(state.outboundPeers, msg.cmp, func(sp *serverWitnessPeer) {
					state.outboundGroups[addrmgr.GroupKey(sp.NA())]--
				})
			}
			msg.reply <- nil
			return
		}

		msg.reply <- errors.New("peer not found")
	}
}

// disconnectPeer attempts to drop the connection of a tageted peer in the
// passed peer list. Targets are identified via usage of the passed
// `compareFunc`, which should return `true` if the passed peer is the target
// peer. This function returns true on success and false if the peer is unable
// to be located. If the peer is found, and the passed callback: `whenFound'
// isn't nil, we call it with the peer as the argument before it is removed
// from the peerList, and is disconnected from the server.
func disconnectWitnessPeer(peerList map[int32]*serverWitnessPeer, compareFunc func(*serverWitnessPeer) bool, whenFound func(*serverWitnessPeer)) bool {
	for addr, peer := range peerList {
		if compareFunc(peer) {
			if whenFound != nil {
				whenFound(peer)
			}

			// This is ok because we are not continuing
			// to iterate so won't corrupt the loop.
			delete(peerList, addr)
			peer.Disconnect()
			return true
		}
	}
	return false
}

// newPeerConfig returns the configuration for the given serverPeer.
func newWitnessPeerConfig(sp *serverPeer) *peer.Config {
	return &peer.Config{
		Listeners: peer.MessageListeners{
			OnVersion:        sp.OnVersion,
			OnMemPool:        sp.OnMemPool,
			OnGetMiningState: sp.OnGetMiningState,
			OnMiningState:    sp.OnMiningState,
			OnTx:             sp.OnTx,
			OnBlock:          sp.OnBlock,
			OnInv:            sp.OnInv,
			OnHeaders:        sp.OnHeaders,
			OnGetData:        sp.OnGetData,
			OnGetBlocks:      sp.OnGetBlocks,
			OnGetHeaders:     sp.OnGetHeaders,
			OnFilterAdd:      sp.OnFilterAdd,
			OnFilterClear:    sp.OnFilterClear,
			OnFilterLoad:     sp.OnFilterLoad,
			OnGetAddr:        sp.OnGetAddr,
			OnAddr:           sp.OnAddr,
			OnRead:           sp.OnRead,
			OnWrite:          sp.OnWrite,
		},
		NewestBlock:      sp.newestBlock,
		HostToNetAddress: sp.server.addrManager.HostToNetAddress,
		Proxy:            cfg.Proxy,
		UserAgentName:    userAgentName,
		UserAgentVersion: userAgentVersion,
		ChainParams:      sp.server.chainParams,
		Services:         sp.server.services,
		DisableRelayTx:   cfg.BlocksOnly,
		ProtocolVersion:  maxProtocolVersion,
	}
}

// inboundPeerConnected is invoked by the connection manager when a new inbound
// connection is established.  It initializes a new inbound server peer
// instance, associates it with the connection, and starts a goroutine to wait
// for disconnection.
func (s *server) inboundWitnessPeerConnected(conn net.Conn) {
	sp := newServerPeer(s, false)
	sp.isWhitelisted = isWhitelisted(conn.RemoteAddr())
	sp.Peer = peer.NewInboundPeer(newPeerConfig(sp))
	sp.AssociateConnection(conn)
	go s.peerDoneHandler(sp)
}

// outboundPeerConnected is invoked by the connection manager when a new
// outbound connection is established.  It initializes a new outbound server
// peer instance, associates it with the relevant state such as the connection
// request instance and the connection itself, and finally notifies the address
// manager of the attempt.
func (s *server) outboundWitnessPeerConnected(c *connmgr.ConnReq, conn net.Conn) {
	sp := newServerPeer(s, c.Permanent)
	p, err := peer.NewOutboundPeer(newPeerConfig(sp), c.Addr.String())
	if err != nil {
		srvrLog.Debugf("Cannot create outbound peer %s: %v", c.Addr, err)
		s.connManager.Disconnect(c.ID())
		return
	}
	sp.Peer = p
	sp.connReq = c
	sp.isWhitelisted = isWhitelisted(conn.RemoteAddr())
	sp.AssociateConnection(conn)
	go s.peerDoneHandler(sp)
	s.addrManager.Attempt(sp.NA())
}

// peerDoneHandler handles peer disconnects by notifiying the server that it's
// done.
func (s *server) witnessPeerDoneHandler(sp *serverPeer) {
	sp.WaitForDisconnect()
	s.donePeers <- sp

	// Only tell block manager we are gone if we ever told it we existed.
	if sp.VersionKnown() {
		s.blockManager.DonePeer(sp)
	}
	close(sp.quit)
}

// peerHandler is used to handle peer operations such as adding and removing
// peers to and from the server, banning peers, and broadcasting messages to
// peers.  It must be run in a goroutine.
func (s *server) witnessPeerHandler() {
	// Start the address manager and block manager, both of which are needed
	// by peers.  This is done here since their lifecycle is closely tied
	// to this handler and rather than adding more channels to sychronize
	// things, it's easier and slightly faster to simply start and stop them
	// in this handler.
	s.witnessAddrManager.Start()

	srvrLog.Tracef("Starting witness peer handler")

	state := &witnessPeerState{
		inboundPeers:    make(map[int32]*serverWitnessPeer),
		persistentPeers: make(map[int32]*serverWitnessPeer),
		outboundPeers:   make(map[int32]*serverWitnessPeer),
		banned:          make(map[string]time.Time),
		outboundGroups:  make(map[string]int),
	}

	if !cfg.DisableDNSSeed {
		// Add peers discovered through DNS to the address manager.
		connmgr.SeedFromWitnessDNS(activeNetParams.Params, hcdLookup, func(addrs []*wire.NetAddress) {
			// Bitcoind uses a lookup of the dns seeder here. This
			// is rather strange since the values looked up by the
			// DNS seed lookups will vary quite a lot.
			// to replicate this behaviour we put all addresses as
			// having come from the first one.
			s.witnessAddrManager.AddAddresses(addrs, addrs[0])
		})
	}
	go s.witnessConnManager.Start()

out:
	for {
		select {
		// New peers connected to the server.
		case p := <-s.newWitnessPeers:
			s.handleAddWitnessPeerMsg(state, p)

			// Disconnected peers.
		case p := <-s.doneWitnessPeers:
			s.handleDoneWitnessPeerMsg(state, p)
			// Peer to ban.
		case p := <-s.banWitnessPeers:
			s.handleBanWitnessPeerMsg(state, p)

			// New inventory to potentially be relayed to other peers.
		case invMsg := <-s.relayWitnessInv:
			s.handleRelayWitnessInvMsg(state, invMsg)

			// Message to broadcastWitness to all connected peers except those
			// which are excluded by the message.
		case bmsg := <-s.broadcastWitness:
			s.handleBroadcastWitnessMsg(state, &bmsg)

		case qmsg := <-s.queryWitness:
			s.handleQueryWitness(state, qmsg)

		case <-s.quit:
			// Disconnect all peers on server shutdown.
			state.forAllPeers(func(sp *serverWitnessPeer) {
				srvrLog.Tracef("Shutdown witness peer %s", sp)
				sp.Disconnect()
			})
			break out
		}
	}

	s.witnessConnManager.Stop()
	s.witnessAddrManager.Stop()

	// Drain channels before exiting so nothing is left waiting around
	// to send.
cleanup:
	for {
		select {
		case <-s.newWitnessPeers:
		case <-s.doneWitnessPeers:
		case <-s.relayWitnessInv:
		case <-s.broadcastWitness:
		case <-s.queryWitness:
		default:
			break cleanup
		}
	}
	s.wg.Done()
	srvrLog.Tracef("Witness Peer handler done")
}

// AddPeer adds a new peer that has already been connected to the server.
func (s *server) AddWitnessPeer(sp *serverWitnessPeer) {
	s.newWitnessPeers <- sp
}

// BanPeer bans a peer that has already been connected to the server by ip.
func (s *server) BanWitnessPeer(sp *serverWitnessPeer) {
	s.banWitnessPeers <- sp
}



// RelayInventory relays the passed inventory vector to all connected peers
// that are not already known to have it.
func (s *server) RelayWitnessInventory(invVect *wire.InvVect, data interface{}) {
	s.relayWitnessInv <- relayMsg{invVect: invVect, data: data}
}

// BroadcastMessage sends msg to all peers currently connected to the server
// except those in the passed peers to exclude.
func (s *server) BroadcastWitnessMessage(msg wire.Message, exclPeers ...*serverWitnessPeer) {
	// XXX: Need to determine if this is an alert that has already been
	// broadcast and refrain from broadcasting again.
	bmsg := broadcastWitnessMsg{message: msg, excludePeers: exclPeers}
	s.broadcastWitness <- bmsg
}

// ConnectedCount returns the number of currently connected peers.
func (s *server) WitnessConnectedCount() int32 {
	replyChan := make(chan int32)

	s.queryWitness <- getConnCountMsg{reply: replyChan}

	return <-replyChan
}

// OutboundGroupCount returns the number of peers connected to the given
// outbound group key.
func (s *server) WitnessOutboundGroupCount(key string) int {
	replyChan := make(chan int)
	s.queryWitness <- getOutboundGroup{key: key, reply: replyChan}
	return <-replyChan
}

// AddedNodeInfo returns an array of hcjson.GetAddedNodeInfoResult structures
// describing the persistent (added) nodes.
func (s *server) AddedWitnessNodeInfo() []*serverWitnessPeer {
	replyChan := make(chan []*serverWitnessPeer)
	s.queryWitness <- getAddedWitnessNodesMsg{reply: replyChan}
	return <-replyChan
}

// Peers returns an array of all connected peers.
func (s *server) WitnessPeers() []*serverWitnessPeer {
	replyChan := make(chan []*serverWitnessPeer)

	s.queryWitness <- getWitnessPeersMsg{reply: replyChan}

	return <-replyChan
}

// DisconnectNodeByAddr disconnects a peer by target address. Both outbound and
// inbound nodes will be searched for the target node. An error message will
// be returned if the peer was not found.
func (s *server) DisconnectWitnessNodeByAddr(addr string) error {
	replyChan := make(chan error)

	s.queryWitness <- disconnectWitnessNodeMsg{
		cmp:   func(sp *serverWitnessPeer) bool { return sp.Addr() == addr },
		reply: replyChan,
	}

	return <-replyChan
}

// DisconnectNodeByID disconnects a peer by target node id. Both outbound and
// inbound nodes will be searched for the target node. An error message will be
// returned if the peer was not found.
func (s *server) WitnessDisconnectNodeByID(id int32) error {
	replyChan := make(chan error)

	s.queryWitness <- disconnectWitnessNodeMsg{
		cmp:   func(sp *serverWitnessPeer) bool { return sp.ID() == id },
		reply: replyChan,
	}

	return <-replyChan
}

// RemoveNodeByAddr removes a peer from the list of persistent peers if
// present. An error will be returned if the peer was not found.
func (s *server) RemoveWitnessNodeByAddr(addr string) error {
	replyChan := make(chan error)

	s.queryWitness <- removeWitnessNodeMsg{
		cmp:   func(sp *serverWitnessPeer) bool { return sp.Addr() == addr },
		reply: replyChan,
	}

	return <-replyChan
}

// RemoveNodeByID removes a peer by node ID from the list of persistent peers
// if present. An error will be returned if the peer was not found.
func (s *server) RemoveWitnessNodeByID(id int32) error {
	replyChan := make(chan error)

	s.queryWitness <- removeWitnessNodeMsg{
		cmp:   func(sp *serverWitnessPeer) bool { return sp.ID() == id },
		reply: replyChan,
	}

	return <-replyChan
}

// ConnectNode adds `addr' as a new outbound peer. If permanent is true then the
// peer will be persistent and reconnect if the connection is lost.
// It is an error to call this with an already existing peer.
func (s *server) ConnectWitnessNode(addr string, permanent bool) error {
	replyChan := make(chan error)

	s.queryWitness <- connectNodeMsg{addr: addr, permanent: permanent, reply: replyChan}

	return <-replyChan
}

// AddBytesSent adds the passed number of bytes to the total bytes sent counter
// for the server.  It is safe for concurrent access.
func (s *server) AddBytesWitnessSent(bytesSent uint64) {
	atomic.AddUint64(&s.bytesWitnessSent, bytesSent)
}

// AddBytesReceived adds the passed number of bytes to the total bytes received
// counter for the server.  It is safe for concurrent access.
func (s *server) AddBytesWitnessReceived(bytesReceived uint64) {
	atomic.AddUint64(&s.bytesWitnessReceived, bytesReceived)
}

// NetTotals returns the sum of all bytes received and sent across the network
// for all peers.  It is safe for concurrent access.
func (s *server) WitnessNetTotals() (uint64, uint64) {
	return atomic.LoadUint64(&s.bytesWitnessReceived),
		atomic.LoadUint64(&s.bytesWitnessSent)
}



// rebroadcastHandler keeps track of user submitted inventories that we have
// sent out but have not yet made it into a block. We periodically rebroadcast
// them in case our peers restarted or otherwise lost track of them.
func (s *server) rebroadcastWitnessHandler() {
	// Wait 5 min before first tx rebroadcast.
	timer := time.NewTimer(5 * time.Minute)
	pendingInvs := make(map[wire.InvVect]interface{})

out:
	for {
		select {
		case riv := <-s.modifyRebroadcastInv:
			switch msg := riv.(type) {
			// Incoming InvVects are added to our map of RPC txs.
			case broadcastInventoryAdd:
				srvrLog.Debugf("Add inventory : %v", msg.invVect)
				pendingInvs[*msg.invVect] = msg.data

				// When an InvVect has been added to a block, we can
				// now remove it, if it was present.
			case broadcastInventoryDel:
				if _, ok := pendingInvs[*msg]; ok {
					srvrLog.Debugf("Remove inventory : %v", msg)
					delete(pendingInvs, *msg)
				}
			}

		case <-timer.C:
			// Any inventory we have has not made it into a block
			// yet. We periodically resubmit them until they have.
			srvrLog.Debugf("Start relay inventory")
			for iv, data := range pendingInvs {
				srvrLog.Debugf("Relay inventory : %v", iv)
				ivCopy := iv
				s.RelayWitnessInventory(&ivCopy, data)
			}
			srvrLog.Debugf("Start relay complete")
			// Process at a random time up to 30mins (in seconds)
			// in the future.
			//timer.Reset(time.Second * time.Duration(randomUint16Number(1800)))
			timer.Reset(time.Second * 300)
		case <-s.quit:
			break out
		}
	}

	timer.Stop()

	// Drain channels before exiting so nothing is left waiting around
	// to send.
cleanup:
	for {
		select {
		case <-s.modifyRebroadcastWitnessInv:
		default:
			break cleanup
		}
	}
	s.wg.Done()
}
type witnessPeerState struct {
	inboundPeers    map[int32]*serverWitnessPeer
	outboundPeers   map[int32]*serverWitnessPeer
	persistentPeers map[int32]*serverWitnessPeer
	banned          map[string]time.Time
	outboundGroups  map[string]int
}

// Count returns the count of all known peers.
func (ps *witnessPeerState) Count() int {
	return len(ps.inboundPeers) + len(ps.outboundPeers) +
		len(ps.persistentPeers)
}

// forAllOutboundPeers is a helper function that runs closure on all outbound
// peers known to peerState.
func (ps *witnessPeerState) forAllOutboundPeers(closure func(sp *serverWitnessPeer)) {
	for _, e := range ps.outboundPeers {
		closure(e)
	}
	for _, e := range ps.persistentPeers {
		closure(e)
	}
}

// forAllPeers is a helper function that runs closure on all peers known to
// peerState.
func (ps *witnessPeerState) forAllPeers(closure func(sp *serverWitnessPeer)) {
	for _, e := range ps.inboundPeers {
		closure(e)
	}
	ps.forAllOutboundPeers(closure)
}
