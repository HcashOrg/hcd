// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2016-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package peer

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/wire"
	"github.com/davecgh/go-spew/spew"
	"io"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// witness nodeCount is the total number of WitnessPeer connections made since startup
	// and is used to assign an id to a WitnessPeer.
	witnessNodeCount int32

	// zeroHash is the zero value hash (all zeros).  It is defined as a
	// convenience.
	witnessZeroHash chainhash.Hash

	//witness sentNonces houses the unique nonces that are generated when pushing
	// version messages that are used to detect self connections.
	witnessSentNonces = newMruNonceMap(50)

	// witness allowSelfConns is only used to allow the tests to bypass the self
	// connection detecting and disconnect logic since they intentionally
	// do so for testing purposes.
	allowSelfWitnessConns bool
)

// witnessMessageListeners defines callback function pointers to invoke with message
// listeners for a WitnessPeer. Any listener which is not set to a concrete callback
// during WitnessPeer initialization is ignored. Execution of multiple message
// listeners occurs serially, so one callback blocks the execution of the next.
//
// NOTE: Unless otherwise documented, these listeners must NOT directly call any
// blocking calls (such as WaitForShutdown) on the WitnessPeer instance since the input
// handler goroutine blocks until the callback has completed.  Doing so will
// result in a deadlock.
type WitnessMessageListeners struct {
	// OnGetAddr is invoked when a WitnessPeer receives a getaddr wire message.
	OnGetAddr func(p *WitnessPeer, msg *wire.MsgGetAddr)

	// OnAddr is invoked when a WitnessPeer receives an addr wire message.
	OnAddr func(p *WitnessPeer, msg *wire.MsgAddr)

	// OnPing is invoked when a WitnessPeer receives a ping wire message.
	OnPing func(p *WitnessPeer, msg *wire.MsgPing)

	// OnPong is invoked when a WitnessPeer receives a pong wire message.
	OnPong func(p *WitnessPeer, msg *wire.MsgPong)

	// OnAlert is invoked when a WitnessPeer receives an alert wire message.
	OnAlert func(p *WitnessPeer, msg *wire.MsgAlert)


	// OnTx is invoked when a WitnessPeer receives a tx wire message.
	OnTx func(p *WitnessPeer, msg *wire.MsgTx)

	// OnInv is invoked when a WitnessPeer receives an inv wire message.
	OnInv func(p *WitnessPeer, msg *wire.MsgInv)


	// OnNotFound is invoked when a WitnessPeer receives a notfound wire message.
	OnNotFound func(p *WitnessPeer, msg *wire.MsgNotFound)

	// OnGetData is invoked when a WitnessPeer receives a getdata wire message.
	OnGetData func(p *WitnessPeer, msg *wire.MsgGetData)

	// OnFeeFilter is invoked when a WitnessPeer receives a feefilter wire message.
	OnFeeFilter func(p *WitnessPeer, msg *wire.MsgFeeFilter)

	// OnFilterAdd is invoked when a WitnessPeer receives a filteradd wire message.
	OnFilterAdd func(p *WitnessPeer, msg *wire.MsgFilterAdd)

	// OnFilterClear is invoked when a WitnessPeer receives a filterclear wire
	// message.
	OnFilterClear func(p *WitnessPeer, msg *wire.MsgFilterClear)

	// OnFilterLoad is invoked when a WitnessPeer receives a filterload wire
	// message.
	OnFilterLoad func(p *WitnessPeer, msg *wire.MsgFilterLoad)

	// OnVersion is invoked when a WitnessPeer receives a version wire message.
	OnVersion func(p *WitnessPeer, msg *wire.MsgVersion)

	// OnVerAck is invoked when a WitnessPeer receives a verack wire message.
	OnVerAck func(p *WitnessPeer, msg *wire.MsgVerAck)

	// OnReject is invoked when a WitnessPeer receives a reject wire message.
	OnReject func(p *WitnessPeer, msg *wire.MsgReject)

	// OnRead is invoked when a WitnessPeer receives a wire message.  It consists
	// of the number of bytes read, the message, and whether or not an error
	// in the read occurred.  Typically, callers will opt to use the
	// callbacks for the specific message types, however this can be useful
	// for circumstances such as keeping track of server-wide byte counts or
	// working with custom message types for which the WitnessPeer does not
	// directly provide a callback.
	OnRead func(p *WitnessPeer, bytesRead int, msg wire.Message, err error)

	// OnWrite is invoked when we write a wire message to a WitnessPeer.  It
	// consists of the number of bytes written, the message, and whether or
	// not an error in the write occurred.  This can be useful for
	// circumstances such as keeping track of server-wide byte counts.
	OnWrite func(p *WitnessPeer, bytesWritten int, msg wire.Message, err error)
}

// witnessConfig is the struct to hold configuration options useful to WitnessPeer.
type WitnessConfig struct {
	// NewestBlock specifies a callback which provides the newest block
	// details to the WitnessPeer as needed.  This can be nil in which case the
	// WitnessPeer will report a block height of 0, however it is good practice for
	// WitnessPeers to specify this so their currently best known is accurately
	// reported.
	NewestBlock HashFunc

	// HostToNetAddress returns the netaddress for the given host. This can be
	// nil in  which case the host will be parsed as an IP address.
	HostToNetAddress HostToNetAddrFunc

	// Proxy indicates a proxy is being used for connections.  The only
	// effect this has is to prevent leaking the tor proxy address, so it
	// only needs to specified if using a tor proxy.
	Proxy string

	// UserAgentName specifies the user agent name to advertise.  It is
	// highly recommended to specify this value.
	UserAgentName string

	// UserAgentVersion specifies the user agent version to advertise.  It
	// is highly recommended to specify this value and that it follows the
	// form "major.minor.revision" e.g. "2.6.41".
	UserAgentVersion string

	// ChainParams identifies which chain parameters the WitnessPeer is associated
	// with.  It is highly recommended to specify this field, however it can
	// be omitted in which case the test network will be used.
	ChainParams *chaincfg.Params

	// Services specifies which services to advertise as supported by the
	// local WitnessPeer.  This field can be omitted in which case it will be 0
	// and therefore advertise no supported services.
	Services wire.ServiceFlag

	// ProtocolVersion specifies the maximum protocol version to use and
	// advertise.  This field can be omitted in which case
	// WitnessPeer.MaxProtocolVersion will be used.
	ProtocolVersion uint32

	// DisableRelayTx specifies if the remote WitnessPeer should be informed to
	// not send inv messages for transactions.
	DisableRelayTx bool

	// Listeners houses callback functions to be invoked on receiving WitnessPeer
	// messages.
	Listeners WitnessMessageListeners
}

// WitnessStatsSnap is a snapshot of WitnessPeer stats at a point in time.
type WitnessStatsSnap struct {
	ID             int32
	Addr           string
	Services       wire.ServiceFlag
	LastSend       time.Time
	LastRecv       time.Time
	BytesSent      uint64
	BytesRecv      uint64
	ConnTime       time.Time
	TimeOffset     int64
	Version        uint32
	UserAgent      string
	Inbound        bool
	StartingHeight int64

	LastPingNonce  uint64
	LastPingTime   time.Time
	LastPingMicros int64
}

// NOTE: The overall data flow of a WitnessPeer is split into 3 goroutines.  Inbound
// messages are read via the inHandler goroutine and generally dispatched to
// their own handler.  For inbound data-related messages such as blocks,
// transactions, and inventory, the data is handled by the corresponding
// message handlers.  The data flow for outbound messages is split into 2
// goroutines, queueHandler and outHandler.  The first, queueHandler, is used
// as a way for external entities to queue messages, by way of the QueueMessage
// function, quickly regardless of whether the WitnessPeer is currently sending or not.
// It acts as the traffic cop between the external world and the actual
// goroutine which writes to the network socket.

// WitnessPeer provides a basic concurrent safe hcd WitnessPeer for handling hcd
// communications via the WitnessPeer-to-WitnessPeer protocol.  It provides full duplex
// reading and writing, automatic handling of the initial handshake process,
// querying of usage statistics and other information about the remote WitnessPeer such
// as its address, user agent, and protocol version, output message queuing,
// inventory trickling, and the ability to dynamically register and unregister
// callbacks for handling hcd protocol messages.
//
// Outbound messages are typically queued via QueueMessage or QueueInventory.
// QueueMessage is intended for all messages, including responses to data such
// as blocks and transactions.  QueueInventory, on the other hand, is only
// intended for relaying inventory as it employs a trickling mechanism to batch
// the inventory together.  However, some helper functions for pushing messages
// of specific types that typically require common special handling are
// provided as a convenience.
type WitnessPeer struct {
	// The following variables must only be used atomically.
	bytesReceived uint64
	bytesSent     uint64
	lastRecv      int64
	lastSend      int64
	connected     int32
	disconnect    int32

	conn net.Conn

	// These fields are set at creation time and never modified, so they are
	// safe to read from concurrently without a mutex.
	addr    string
	cfg     WitnessConfig
	inbound bool

	flagsMtx             sync.Mutex // protects the peer flags below
	na                   *wire.NetAddress
	id                   int32
	userAgent            string
	services             wire.ServiceFlag
	versionKnown         bool
	advertisedProtoVer   uint32 // protocol version advertised by remote
	protocolVersion      uint32 // negotiated protocol version
	sendHeadersPreferred bool   // WitnessPeer sent a sendheaders message
	versionSent          bool
	verAckReceived       bool

	knownInventory *mruInventoryMap

	// These fields keep track of statistics for the WitnessPeer and are protected
	// by the statsMtx mutex.
	statsMtx       sync.RWMutex
	timeOffset     int64
	timeConnected  time.Time
	startingHeight int64

	lastPingNonce  uint64    // Set to nonce if we have a pending ping.
	lastPingTime   time.Time // Time we sent last ping.
	lastPingMicros int64     // Time for last ping to return.

	stallControl  chan stallControlMsg
	outputQueue   chan outMsg
	sendQueue     chan outMsg
	sendDoneQueue chan struct{}
	outputInvChan chan *wire.InvVect
	inQuit        chan struct{}
	queueQuit     chan struct{}
	outQuit       chan struct{}
	quit          chan struct{}
}

// String returns the WitnessPeer's address and directionality as a human-readable
// string.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) String() string {
	return fmt.Sprintf("%s (%s)", p.addr, directionString(p.inbound))
}

// AddKnownWitnessInventory adds the passed inventory to the cache of known inventory
// for the WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) AddKnownWitnessInventory(invVect *wire.InvVect) {
	p.knownInventory.Add(invVect)
}

// WitnessStatsSnapshot returns a snapshot of the current WitnessPeer flags and statistics.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) WitnessStatsSnapshot() *WitnessStatsSnap {
	p.statsMtx.RLock()

	p.flagsMtx.Lock()
	id := p.id
	addr := p.addr
	userAgent := p.userAgent
	services := p.services
	protocolVersion := p.advertisedProtoVer
	p.flagsMtx.Unlock()

	// Get a copy of all relevant flags and stats.
	witnessStatsSnap := &WitnessStatsSnap{
		ID:             id,
		Addr:           addr,
		UserAgent:      userAgent,
		Services:       services,
		LastSend:       p.LastSend(),
		LastRecv:       p.LastRecv(),
		BytesSent:      p.BytesSent(),
		BytesRecv:      p.BytesReceived(),
		ConnTime:       p.timeConnected,
		TimeOffset:     p.timeOffset,
		Version:        protocolVersion,
		Inbound:        p.inbound,
		StartingHeight: p.startingHeight,
		LastPingNonce:  p.lastPingNonce,
		LastPingMicros: p.lastPingMicros,
		LastPingTime:   p.lastPingTime,
	}

	p.statsMtx.RUnlock()
	return witnessStatsSnap
}

// ID returns the WitnessPeer id.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) ID() int32 {
	p.flagsMtx.Lock()
	id := p.id
	p.flagsMtx.Unlock()

	return id
}

// NA returns the WitnessPeer network address.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) NA() *wire.NetAddress {
	p.flagsMtx.Lock()
	na := p.na
	p.flagsMtx.Unlock()

	return na
}

// Addr returns the WitnessPeer address.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) Addr() string {
	// The address doesn't change after initialization, therefore it is not
	// protected by a mutex.
	return p.addr
}

// Inbound returns whether the WitnessPeer is inbound.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) Inbound() bool {
	return p.inbound
}

// Services returns the services flag of the remote WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) Services() wire.ServiceFlag {
	p.flagsMtx.Lock()
	services := p.services
	p.flagsMtx.Unlock()

	return services
}

// UserAgent returns the user agent of the remote WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) UserAgent() string {
	p.flagsMtx.Lock()
	userAgent := p.userAgent
	p.flagsMtx.Unlock()

	return userAgent
}

// LastPingNonce returns the last ping nonce of the remote WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) LastPingNonce() uint64 {
	p.statsMtx.RLock()
	lastPingNonce := p.lastPingNonce
	p.statsMtx.RUnlock()

	return lastPingNonce
}

// LastPingTime returns the last ping time of the remote WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) LastPingTime() time.Time {
	p.statsMtx.RLock()
	lastPingTime := p.lastPingTime
	p.statsMtx.RUnlock()

	return lastPingTime
}

// LastPingMicros returns the last ping micros of the remote WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) LastPingMicros() int64 {
	p.statsMtx.RLock()
	lastPingMicros := p.lastPingMicros
	p.statsMtx.RUnlock()

	return lastPingMicros
}

// VersionKnown returns the whether or not the version of a WitnessPeer is known
// locally.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) VersionKnown() bool {
	p.flagsMtx.Lock()
	versionKnown := p.versionKnown
	p.flagsMtx.Unlock()

	return versionKnown
}

// VerAckReceived returns whether or not a verack message was received by the
// WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) VerAckReceived() bool {
	p.flagsMtx.Lock()
	verAckReceived := p.verAckReceived
	p.flagsMtx.Unlock()

	return verAckReceived
}

// ProtocolVersion returns the negotiated WitnessPeer protocol version.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) ProtocolVersion() uint32 {
	p.flagsMtx.Lock()
	protocolVersion := p.protocolVersion
	p.flagsMtx.Unlock()

	return protocolVersion
}

// LastSend returns the last send time of the WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) LastSend() time.Time {
	return time.Unix(atomic.LoadInt64(&p.lastSend), 0)
}

// LastRecv returns the last recv time of the WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) LastRecv() time.Time {
	return time.Unix(atomic.LoadInt64(&p.lastRecv), 0)
}

// BytesSent returns the total number of bytes sent by the WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) BytesSent() uint64 {
	return atomic.LoadUint64(&p.bytesSent)
}

// BytesReceived returns the total number of bytes received by the WitnessPeer.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) BytesReceived() uint64 {
	return atomic.LoadUint64(&p.bytesReceived)
}

// TimeConnected returns the time at which the WitnessPeer connected.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) TimeConnected() time.Time {
	p.statsMtx.RLock()
	timeConnected := p.timeConnected
	p.statsMtx.RUnlock()

	return timeConnected
}

// TimeOffset returns the number of seconds the local time was offset from the
// time the WitnessPeer reported during the initial negotiation phase.  Negative values
// indicate the remote WitnessPeer's time is before the local time.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) TimeOffset() int64 {
	p.statsMtx.RLock()
	timeOffset := p.timeOffset
	p.statsMtx.RUnlock()

	return timeOffset
}

// StartingHeight returns the last known height the WitnessPeer reported during the
// initial negotiation phase.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) StartingHeight() int64 {
	p.statsMtx.RLock()
	startingHeight := p.startingHeight
	p.statsMtx.RUnlock()

	return startingHeight
}

// WantsHeaders returns if the WitnessPeer wants header messages instead of
// inventory vectors for blocks.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) WantsHeaders() bool {
	p.flagsMtx.Lock()
	sendHeadersPreferred := p.sendHeadersPreferred
	p.flagsMtx.Unlock()

	return sendHeadersPreferred
}

// localVersionMsg creates a version message that can be used to send to the
// remote WitnessPeer.
func (p *WitnessPeer) localVersionMsg() (*wire.MsgVersion, error) {
	var blockNum int64
	if p.cfg.NewestBlock != nil {
		var err error
		_, blockNum, err = p.cfg.NewestBlock()
		if err != nil {
			return nil, err
		}
	}

	theirNA := p.na

	// If we are behind a proxy and the connection comes from the proxy then
	// we return an unroutable address as their address. This is to prevent
	// leaking the tor proxy address.
	if p.cfg.Proxy != "" {
		proxyaddress, _, err := net.SplitHostPort(p.cfg.Proxy)
		// invalid proxy means poorly configured, be on the safe side.
		if err != nil || p.na.IP.String() == proxyaddress {
			theirNA = wire.NewNetAddressIPPort(net.IP([]byte{0, 0, 0, 0}), 0, theirNA.Services)
		}
	}

	// Create a wire.NetAddress with only the services set to use as the
	// "addrme" in the version message.
	//
	// Older nodes previously added the IP and port information to the
	// address manager which proved to be unreliable as an inbound
	// connection from a WitnessPeer didn't necessarily mean the WitnessPeer itself
	// accepted inbound connections.
	//
	// Also, the timestamp is unused in the version message.
	ourNA := &wire.NetAddress{
		Services: p.cfg.Services,
	}

	// Generate a unique nonce for this WitnessPeer so self connections can be
	// detected.  This is accomplished by adding it to a size-limited map of
	// recently seen nonces.
	nonce, err := wire.RandomUint64()
	if err != nil {
		return nil, err
	}
	witnessSentNonces.Add(nonce)

	// Version message.
	msg := wire.NewMsgVersion(ourNA, theirNA, nonce, int32(blockNum))
	msg.AddUserAgent(p.cfg.UserAgentName, p.cfg.UserAgentVersion)

	// XXX: bitcoind appears to always enable the full node services flag
	// of the remote WitnessPeer netaddress field in the version message regardless
	// of whether it knows it supports it or not.  Also, bitcoind sets
	// the services field of the local WitnessPeer to 0 regardless of support.
	//
	// Realistically, this should be set as follows:
	// - For outgoing connections:
	//    - Set the local netaddress services to what the local WitnessPeer
	//      actually supports
	//    - Set the remote netaddress services to 0 to indicate no services
	//      as they are still unknown
	// - For incoming connections:
	//    - Set the local netaddress services to what the local WitnessPeer
	//      actually supports
	//    - Set the remote netaddress services to the what was advertised by
	//      by the remote WitnessPeer in its version message
	//msg.AddrYou.Services = wire.SFNodeNetwork

	// Advertise local services.
	// Advertise the services flag
	msg.Services = p.cfg.Services

	// Advertise our max supported protocol version.
	msg.ProtocolVersion = int32(p.ProtocolVersion())

	// Advertise if inv messages for transactions are desired.
	msg.DisableRelayTx = p.cfg.DisableRelayTx

	return msg, nil
}

// PushAddrMsg sends an addr message to the connected WitnessPeer using the provided
// addresses.  This function is useful over manually sending the message via
// QueueMessage since it automatically limits the addresses to the maximum
// number allowed by the message and randomizes the chosen addresses when there
// are too many.  It returns the addresses that were actually sent and no
// message will be sent if there are no entries in the provided addresses slice.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) PushAddrMsg(addresses []*wire.NetAddress) ([]*wire.NetAddress, error) {

	// Nothing to send.
	if len(addresses) == 0 {
		return nil, nil
	}

	msg := wire.NewMsgAddr()
	msg.AddrList = make([]*wire.NetAddress, len(addresses))
	copy(msg.AddrList, addresses)

	// Randomize the addresses sent if there are more than the maximum allowed.
	if len(msg.AddrList) > wire.MaxAddrPerMsg {
		// Shuffle the address list.
		for i := range msg.AddrList {
			j := rand.Intn(i + 1)
			msg.AddrList[i], msg.AddrList[j] = msg.AddrList[j], msg.AddrList[i]
		}

		// Truncate it to the maximum size.
		msg.AddrList = msg.AddrList[:wire.MaxAddrPerMsg]
	}

	p.QueueMessage(msg, nil)
	return msg.AddrList, nil
}

// PushRejectMsg sends a reject message for the provided command, reject code,
// reject reason, and hash.  The hash will only be used when the command is a tx
// or block and should be nil in other cases.  The wait parameter will cause the
// function to block until the reject message has actually been sent.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) PushRejectMsg(command string, code wire.RejectCode, reason string, hash *chainhash.Hash, wait bool) {
	msg := wire.NewMsgReject(command, code, reason)
	if command == wire.CmdTx || command == wire.CmdBlock {
		if hash == nil {
			log.Warnf("Sending a reject message for command "+
				"type %v which should have specified a hash "+
				"but does not", command)
			hash = &witnessZeroHash
		}
		msg.Hash = *hash
	}

	// Send the message without waiting if the caller has not requested it.
	if !wait {
		p.QueueMessage(msg, nil)
		return
	}

	// Send the message and block until it has been sent before returning.
	doneChan := make(chan struct{}, 1)
	p.QueueMessage(msg, doneChan)
	<-doneChan
}

// handleRemoteVersionMsg is invoked when a version wire message is received
// from the remote WitnessPeer.  It will return an error if the remote WitnessPeer's version
// is not compatible with ours.
func (p *WitnessPeer) handleRemoteVersionMsg(msg *wire.MsgVersion) error {
	// Detect self connections.
	if !allowSelfWitnessConns && witnessSentNonces.Exists(msg.Nonce) {
		return errors.New("disconnecting WitnessPeer connected to self")
	}

	if msg.ProtocolVersion <int32(wire.WitnessVersion){
		reason := fmt.Sprintf("witness protocol version must be %d or greater",
			wire.WitnessVersion)
		rejectMsg := wire.NewMsgReject(msg.Command(), wire.RejectObsolete,
			reason)
		return p.writeMessage(rejectMsg)
	}

	// Notify and disconnect clients that have a protocol version that is
	// too old.
	if msg.ProtocolVersion < int32(wire.InitialProcotolVersion) {
		// Send a reject message indicating the protocol version is
		// obsolete and wait for the message to be sent before
		// disconnecting.
		reason := fmt.Sprintf("protocol version must be %d or greater",
			wire.InitialProcotolVersion)
		rejectMsg := wire.NewMsgReject(msg.Command(), wire.RejectObsolete,
			reason)
		return p.writeMessage(rejectMsg)
	}

	// Limit to one version message per WitnessPeer.
	// No read lock is necessary because versionKnown is not written to in any
	// other goroutine
	if p.versionKnown {
		// Send a reject message indicating the version message was
		// incorrectly sent twice and wait for the message to be sent
		// before disconnecting.
		p.PushRejectMsg(msg.Command(), wire.RejectDuplicate,
			"duplicate version message", nil, true)
		return errors.New("only one version message per WitnessPeer is allowed")
	}

	// Updating a bunch of stats.
	p.statsMtx.Lock()

	p.startingHeight = int64(msg.LastBlock)

	// Set the WitnessPeer's time offset.
	p.timeOffset = msg.Timestamp.Unix() - time.Now().Unix()
	p.statsMtx.Unlock()

	// Negotiate the protocol version.
	p.flagsMtx.Lock()
	p.advertisedProtoVer = uint32(msg.ProtocolVersion)
	p.protocolVersion = minUint32(p.protocolVersion, p.advertisedProtoVer)
	p.versionKnown = true
	log.Debugf("Negotiated protocol version %d for WitnessPeer %s",
		p.protocolVersion, p)
	// Set the WitnessPeer's ID.
	p.id = atomic.AddInt32(&witnessNodeCount, 1)
	// Set the supported services for the WitnessPeer to what the remote WitnessPeer
	// advertised.
	p.services = msg.Services

	p.na.Services = msg.Services
	// Set the remote WitnessPeer's user agent.
	p.userAgent = msg.UserAgent
	p.flagsMtx.Unlock()
	return nil
}

// handlePingMsg is invoked when a WitnessPeer receives a ping wire message.  For
// recent clients (protocol version > BIP0031Version), it replies with a pong
// message.  For older clients, it does nothing and anything other than failure
// is considered a successful ping.
func (p *WitnessPeer) handlePingMsg(msg *wire.MsgPing) {
	// Include nonce from ping so pong can be identified.
	p.QueueMessage(wire.NewMsgPong(msg.Nonce), nil)
}

// handlePongMsg is invoked when a WitnessPeer receives a pong wire message.  It
// updates the ping statistics as required for recent clients (protocol
// version > BIP0031Version).  There is no effect for older clients or when a
// ping was not previously sent.
func (p *WitnessPeer) handlePongMsg(msg *wire.MsgPong) {
	// Arguably we could use a buffered channel here sending data
	// in a fifo manner whenever we send a ping, or a list keeping track of
	// the times of each ping. For now we just make a best effort and
	// only record stats if it was for the last ping sent. Any preceding
	// and overlapping pings will be ignored. It is unlikely to occur
	// without large usage of the ping rpc call since we ping infrequently
	// enough that if they overlap we would have timed out the WitnessPeer.
	p.statsMtx.Lock()
	if p.lastPingNonce != 0 && msg.Nonce == p.lastPingNonce {
		p.lastPingMicros = time.Since(p.lastPingTime).Nanoseconds()
		p.lastPingMicros /= 1000 // convert to usec.
		p.lastPingNonce = 0
	}
	p.statsMtx.Unlock()
}

// readMessage reads the next wire message from the WitnessPeer with logging.
func (p *WitnessPeer) readMessage() (wire.Message, []byte, error) {
	n, msg, buf, err := wire.ReadMessageN(p.conn, p.ProtocolVersion(),
		p.cfg.ChainParams.Net)
	atomic.AddUint64(&p.bytesReceived, uint64(n))
	if p.cfg.Listeners.OnRead != nil {
		p.cfg.Listeners.OnRead(p, n, msg, err)
	}
	if err != nil {
		return nil, nil, err
	}

	// Use closures to log expensive operations so they are only run when
	// the logging level requires it.
	log.Debugf("%v", newLogClosure(func() string {
		// Debug summary of message.
		summary := messageSummary(msg)
		if len(summary) > 0 {
			summary = " (" + summary + ")"
		}
		return fmt.Sprintf("Received %v%s from %s",
			msg.Command(), summary, p)
	}))
	log.Tracef("%v", newLogClosure(func() string {
		return spew.Sdump(msg)
	}))
	log.Tracef("%v", newLogClosure(func() string {
		return spew.Sdump(buf)
	}))

	return msg, buf, nil
}

// writeMessage sends a wire message to the WitnessPeer with logging.
func (p *WitnessPeer) writeMessage(msg wire.Message) error {
	// Don't do anything if we're disconnecting.
	if atomic.LoadInt32(&p.disconnect) != 0 {
		return nil
	}

	// Use closures to log expensive operations so they are only run when
	// the logging level requires it.
	log.Debugf("%v", newLogClosure(func() string {
		// Debug summary of message.
		summary := messageSummary(msg)
		if len(summary) > 0 {
			summary = " (" + summary + ")"
		}
		return fmt.Sprintf("Sending %v%s to %s", msg.Command(),
			summary, p)
	}))
	log.Tracef("%v", newLogClosure(func() string {
		return spew.Sdump(msg)
	}))
	log.Tracef("%v", newLogClosure(func() string {
		var buf bytes.Buffer
		err := wire.WriteMessage(&buf, msg, p.ProtocolVersion(),
			p.cfg.ChainParams.Net)
		if err != nil {
			return err.Error()
		}
		return spew.Sdump(buf.Bytes())
	}))

	// Write the message to the WitnessPeer.
	n, err := wire.WriteMessageN(p.conn, msg, p.ProtocolVersion(),
		p.cfg.ChainParams.Net)
	atomic.AddUint64(&p.bytesSent, uint64(n))
	if p.cfg.Listeners.OnWrite != nil {
		p.cfg.Listeners.OnWrite(p, n, msg, err)
	}
	return err
}

// shouldHandleReadError returns whether or not the passed error, which is
// expected to have come from reading from the remote WitnessPeer in the inHandler,
// should be logged and responded to with a reject message.
func (p *WitnessPeer) shouldHandleReadError(err error) bool {
	// No logging or reject message when the WitnessPeer is being forcibly
	// disconnected.
	if atomic.LoadInt32(&p.disconnect) != 0 {
		return false
	}

	// No logging or reject message when the remote WitnessPeer has been
	// disconnected.
	if err == io.EOF {
		return false
	}
	if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
		return false
	}

	return true
}

// maybeAddDeadline potentially adds a deadline for the appropriate expected
// response for the passed wire protocol command to the pending responses map.
func (p *WitnessPeer) maybeAddDeadline(pendingResponses map[string]time.Time, msgCmd string) {
	// Setup a deadline for each message being sent that expects a response.
	//
	// NOTE: Pings are intentionally ignored here since they are typically
	// sent asynchronously and as a result of a long backlock of messages,
	// such as is typical in the case of initial block download, the
	// response won't be received in time.
	log.Debugf("Adding deadline for command %s for WitnessPeer %s", msgCmd, p.addr)

	deadline := time.Now().Add(stallResponseTimeout)
	switch msgCmd {
	case wire.CmdVersion:
		// Expects a verack message.
		pendingResponses[wire.CmdVerAck] = deadline

	case wire.CmdGetData:
		// Expects tx, or notfound message.
		pendingResponses[wire.CmdTx] = deadline
		pendingResponses[wire.CmdNotFound] = deadline

	}
}

// stallHandler handles stall detection for the WitnessPeer.  This entails keeping
// track of expected responses and assigning them deadlines while accounting for
// the time spent in callbacks.  It must be run as a goroutine.
func (p *WitnessPeer) stallHandler() {
	// These variables are used to adjust the deadline times forward by the
	// time it takes callbacks to execute.  This is done because new
	// messages aren't read until the previous one is finished processing
	// (which includes callbacks), so the deadline for receiving a response
	// for a given message must account for the processing time as well.
	var handlerActive bool
	var handlersStartTime time.Time
	var deadlineOffset time.Duration

	// pendingResponses tracks the expected response deadline times.
	pendingResponses := make(map[string]time.Time)

	// stallTicker is used to periodically check pending responses that have
	// exceeded the expected deadline and disconnect the WitnessPeer due to
	// stalling.
	stallTicker := time.NewTicker(stallTickInterval)
	defer stallTicker.Stop()

	// ioStopped is used to detect when both the input and output handler
	// goroutines are done.
	var ioStopped bool
out:
	for {
		select {
		case msg := <-p.stallControl:
			switch msg.command {
			case sccSendMessage:
				// Add a deadline for the expected response
				// message if needed.
				p.maybeAddDeadline(pendingResponses,
					msg.message.Command())

			case sccReceiveMessage:
				// Remove received messages from the expected
				// response map.  Since certain commands expect
				// one of a group of responses, remove
				// everything in the expected group accordingly.
				switch msgCmd := msg.message.Command(); msgCmd {

				case wire.CmdTx:
					fallthrough
				case wire.CmdNotFound:
					delete(pendingResponses, wire.CmdTx)
					delete(pendingResponses, wire.CmdNotFound)

				default:
					delete(pendingResponses, msgCmd)
				}

			case sccHandlerStart:
				// Warn on unbalanced callback signalling.
				if handlerActive {
					log.Warn("Received handler start " +
						"control command while a " +
						"handler is already active")
					continue
				}

				handlerActive = true
				handlersStartTime = time.Now()

			case sccHandlerDone:
				// Warn on unbalanced callback signalling.
				if !handlerActive {
					log.Warn("Received handler done " +
						"control command when a " +
						"handler is not already active")
					continue
				}

				// Extend active deadlines by the time it took
				// to execute the callback.
				duration := time.Since(handlersStartTime)
				deadlineOffset += duration
				handlerActive = false

			default:
				log.Warnf("Unsupported message command %v",
					msg.command)
			}

		case <-stallTicker.C:
			// Calculate the offset to apply to the deadline based
			// on how long the handlers have taken to execute since
			// the last tick.
			now := time.Now()
			offset := deadlineOffset
			if handlerActive {
				offset += now.Sub(handlersStartTime)
			}

			// Disconnect the WitnessPeer if any of the pending responses
			// don't arrive by their adjusted deadline.
			for command, deadline := range pendingResponses {
				if now.Before(deadline.Add(offset)) {
					log.Debugf("Stall ticker rolling over for WitnessPeer %s on "+
						"cmd %s (deadline for data: %s)", p, command,
						deadline.String())
					continue
				}

				log.Infof("WitnessPeer %s appears to be stalled or "+
					"misbehaving, %s timeout -- "+
					"disconnecting", p, command)
				p.Disconnect()

				break
			}

			// Reset the deadline offset for the next tick.
			deadlineOffset = 0

		case <-p.inQuit:
			// The stall handler can exit once both the input and
			// output handler goroutines are done.
			if ioStopped {
				break out
			}
			ioStopped = true

		case <-p.outQuit:
			// The stall handler can exit once both the input and
			// output handler goroutines are done.
			if ioStopped {
				break out
			}
			ioStopped = true
		}
	}

	// Drain any wait channels before going away so there is nothing left
	// waiting on this goroutine.
cleanup:
	for {
		select {
		case <-p.stallControl:
		default:
			break cleanup
		}
	}
	log.Tracef("WitnessPeer stall handler done for %s", p)
}

// inHandler handles all incoming messages for the WitnessPeer.  It must be run as a
// goroutine.
func (p *WitnessPeer) inHandler() {
	// WitnessPeers must complete the initial version negotiation within a shorter
	// timeframe than a general idle timeout.  The timer is then reset below
	// to idleTimeout for all future messages.
	idleTimer := time.AfterFunc(idleTimeout, func() {
		log.Warnf("WitnessPeer %s no answer for %s -- disconnecting", p, idleTimeout)
		p.Disconnect()
	})

out:
	for atomic.LoadInt32(&p.disconnect) == 0 {
		// Read a message and stop the idle timer as soon as the read
		// is done.  The timer is reset below for the next iteration if
		// needed.
		rmsg, _, err := p.readMessage()
		idleTimer.Stop()
		if err != nil {
			// Only log the error and send reject message if the
			// local WitnessPeer is not forcibly disconnecting and the
			// remote WitnessPeer has not disconnected.
			if p.shouldHandleReadError(err) {
				errMsg := fmt.Sprintf("Can't read message from %s: %v", p, err)
				log.Errorf(errMsg)

				// Push a reject message for the malformed message and wait for
				// the message to be sent before disconnecting.
				//
				// NOTE: Ideally this would include the command in the header if
				// at least that much of the message was valid, but that is not
				// currently exposed by wire, so just used malformed for the
				// command.
				p.PushRejectMsg("malformed", wire.RejectMalformed, errMsg, nil,
					true)
			}
			break out
		}
		atomic.StoreInt64(&p.lastRecv, time.Now().Unix())
		p.stallControl <- stallControlMsg{sccReceiveMessage, rmsg}

		// Handle each supported message type.
		p.stallControl <- stallControlMsg{sccHandlerStart, rmsg}
		switch msg := rmsg.(type) {
		case *wire.MsgVersion:
			p.PushRejectMsg(msg.Command(), wire.RejectDuplicate,
				"duplicate version message", nil, true)
			break out

		case *wire.MsgVerAck:

			// No read lock is necessary because verAckReceived is not written
			// to in any other goroutine.
			if p.verAckReceived {
				log.Infof("Already received 'verack' from WitnessPeer %v -- "+
					"disconnecting", p)
				break out
			}
			p.flagsMtx.Lock()
			p.verAckReceived = true
			p.flagsMtx.Unlock()
			if p.cfg.Listeners.OnVerAck != nil {
				p.cfg.Listeners.OnVerAck(p, msg)
			}

		case *wire.MsgGetAddr:
			if p.cfg.Listeners.OnGetAddr != nil {
				p.cfg.Listeners.OnGetAddr(p, msg)
			}

		case *wire.MsgAddr:
			if p.cfg.Listeners.OnAddr != nil {
				p.cfg.Listeners.OnAddr(p, msg)
			}

		case *wire.MsgPing:
			p.handlePingMsg(msg)
			if p.cfg.Listeners.OnPing != nil {
				p.cfg.Listeners.OnPing(p, msg)
			}

		case *wire.MsgPong:
			p.handlePongMsg(msg)
			if p.cfg.Listeners.OnPong != nil {
				p.cfg.Listeners.OnPong(p, msg)
			}

		case *wire.MsgAlert:
			if p.cfg.Listeners.OnAlert != nil {
				p.cfg.Listeners.OnAlert(p, msg)
			}

		case *wire.MsgTx:
			if p.cfg.Listeners.OnTx != nil {
				p.cfg.Listeners.OnTx(p, msg)
			}

		case *wire.MsgInv:
			if p.cfg.Listeners.OnInv != nil {
				p.cfg.Listeners.OnInv(p, msg)
			}


		case *wire.MsgNotFound:
			if p.cfg.Listeners.OnNotFound != nil {
				p.cfg.Listeners.OnNotFound(p, msg)
			}

		case *wire.MsgGetData:
			if p.cfg.Listeners.OnGetData != nil {
				p.cfg.Listeners.OnGetData(p, msg)
			}

		case *wire.MsgFeeFilter:
			if p.cfg.Listeners.OnFeeFilter != nil {
				p.cfg.Listeners.OnFeeFilter(p, msg)
			}

		case *wire.MsgFilterAdd:
			if p.cfg.Listeners.OnFilterAdd != nil {
				p.cfg.Listeners.OnFilterAdd(p, msg)
			}

		case *wire.MsgFilterClear:
			if p.cfg.Listeners.OnFilterClear != nil {
				p.cfg.Listeners.OnFilterClear(p, msg)
			}

		case *wire.MsgFilterLoad:
			if p.cfg.Listeners.OnFilterLoad != nil {
				p.cfg.Listeners.OnFilterLoad(p, msg)
			}

		case *wire.MsgReject:
			if p.cfg.Listeners.OnReject != nil {
				p.cfg.Listeners.OnReject(p, msg)
			}

		default:
			log.Debugf("Received unhandled message of type %v "+
				"from %v", rmsg.Command(), p)
		}
		p.stallControl <- stallControlMsg{sccHandlerDone, rmsg}

		// A message was received so reset the idle timer.
		idleTimer.Reset(idleTimeout)
	}

	// Ensure the idle timer is stopped to avoid leaking the resource.
	idleTimer.Stop()

	// Ensure connection is closed.
	p.Disconnect()

	close(p.inQuit)
	log.Tracef("WitnessPeer input handler done for %s", p)
}

//KnownInverntory list
func (p *WitnessPeer) KnownInventory() *list.List {
	return p.knownInventory.invList
}

// queueHandler handles the queuing of outgoing data for the WitnessPeer. This runs as
// a muxer for various sources of input so we can ensure that server and WitnessPeer
// handlers will not block on us sending a message.  That data is then passed on
// to outHandler to be actually written.
func (p *WitnessPeer) queueHandler() {
	//pendingMsgs := list.New()
	//invSendQueue := list.New()
	var pendingMsgs []outMsg
	var invSendQueue []*wire.InvVect
	trickleTicker := time.NewTicker(trickleTimeout)
	defer trickleTicker.Stop()

	// We keep the waiting flag so that we know if we have a message queued
	// to the outHandler or not.  We could use the presence of a head of
	// the list for this but then we have rather racy concerns about whether
	// it has gotten it at cleanup time - and thus who sends on the
	// message's done channel.  To avoid such confusion we keep a different
	// flag and pendingMsgs only contains messages that we have not yet
	// passed to outHandler.
	waiting := false

	// To avoid duplication below.
	queuePacket := func(msg outMsg, list *[]outMsg, waiting bool) bool {
		if !waiting {
			p.sendQueue <- msg
		} else {
			//list.PushBack(msg)
			*list = append(*list, msg)
		}
		// we are always waiting now.
		return true
	}
out:
	for {
		select {
		case msg := <-p.outputQueue:
			//waiting = queuePacket(msg, pendingMsgs, waiting)
			waiting = queuePacket(msg, &pendingMsgs, waiting)

			// This channel is notified when a message has been sent across
			// the network socket.
		case <-p.sendDoneQueue:
			// No longer waiting if there are no more messages
			// in the pending messages queue.
			//next := pendingMsgs.Front()
			//if next == nil {
			if len(pendingMsgs) == 0 {
				waiting = false
				continue
			}

			// Notify the outHandler about the next item to
			// asynchronously send.
			//val := pendingMsgs.Remove(next)
			//p.sendQueue <- val.(outMsg)

			next := pendingMsgs[0]
			pendingMsgs = pendingMsgs[1:]
			p.sendQueue <- next

		case iv := <-p.outputInvChan:
			// No handshake?  They'll find out soon enough.
			if p.VersionKnown() {
				//invSendQueue.PushBack(iv)
				invSendQueue = append(invSendQueue, iv)
			}

		case <-trickleTicker.C:
			// Don't send anything if we're disconnecting or there
			// is no queued inventory.
			// version is known if send queue has any entries.
			if atomic.LoadInt32(&p.disconnect) != 0 ||
			//invSendQueue.Len() == 0 {
				len(invSendQueue) == 0 {
				continue
			}

			// Create and send as many inv messages as needed to
			// drain the inventory send queue.
			//invMsg := wire.NewMsgInvSizeHint(uint(invSendQueue.Len()))
			invMsg := wire.NewMsgInvSizeHint(uint(len(invSendQueue)))
			//for e := invSendQueue.Front(); e != nil; e = invSendQueue.Front() {
			for _, iv := range invSendQueue {
				//iv := invSendQueue.Remove(e).(*wire.InvVect)

				// Don't send inventory that became known after
				// the initial check.
				if p.knownInventory.Exists(iv) {
					continue
				}

				invMsg.AddInvVect(iv)
				if len(invMsg.InvList) >= maxInvTrickleSize {
					waiting = queuePacket(
						outMsg{msg: invMsg},
						&pendingMsgs, waiting)
					//invMsg = wire.NewMsgInvSizeHint(uint(invSendQueue.Len()))
					invMsg = wire.NewMsgInvSizeHint(uint(len(invSendQueue)))
				}

				// Add the inventory that is being relayed to
				// the known inventory for the WitnessPeer.
				p.AddKnownWitnessInventory(iv)
			}
			if len(invMsg.InvList) > 0 {
				waiting = queuePacket(outMsg{msg: invMsg},
					&pendingMsgs, waiting)
			}
			invSendQueue = nil

		case <-p.quit:
			break out
		}
	}

	// Drain any wait channels before we go away so we don't leave something
	// waiting for us.
	//for e := pendingMsgs.Front(); e != nil; e = pendingMsgs.Front() {
	//	val := pendingMsgs.Remove(e)
	//msg := val.(outMsg)
	for _, msg := range pendingMsgs {
		if msg.doneChan != nil {
			msg.doneChan <- struct{}{}
		}
	}
cleanup:
	for {
		select {
		case msg := <-p.outputQueue:
			if msg.doneChan != nil {
				msg.doneChan <- struct{}{}
			}
		case <-p.outputInvChan:
			// Just drain channel
			// sendDoneQueue is buffered so doesn't need draining.
		default:
			break cleanup
		}
	}
	close(p.queueQuit)
	log.Tracef("WitnessPeer queue handler done for %s", p)
}

// shouldLogWriteError returns whether or not the passed error, which is
// expected to have come from writing to the remote WitnessPeer in the outHandler,
// should be logged.
func (p *WitnessPeer) shouldLogWriteError(err error) bool {
	// No logging when the WitnessPeer is being forcibly disconnected.
	if atomic.LoadInt32(&p.disconnect) != 0 {
		return false
	}

	// No logging when the remote WitnessPeer has been disconnected.
	if err == io.EOF {
		return false
	}
	if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
		return false
	}

	return true
}

// outHandler handles all outgoing messages for the WitnessPeer.  It must be run as a
// goroutine.  It uses a buffered channel to serialize output messages while
// allowing the sender to continue running asynchronously.
func (p *WitnessPeer) outHandler() {
	// pingTicker is used to periodically send pings to the remote WitnessPeer.
	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

out:
	for {
		select {
		case msg := <-p.sendQueue:
			switch m := msg.msg.(type) {
			case *wire.MsgPing:
				// Setup ping statistics.
				p.statsMtx.Lock()
				p.lastPingNonce = m.Nonce
				p.lastPingTime = time.Now()
				p.statsMtx.Unlock()
			}

			p.stallControl <- stallControlMsg{sccSendMessage, msg.msg}
			if err := p.writeMessage(msg.msg); err != nil {
				p.Disconnect()
				if p.shouldLogWriteError(err) {
					log.Errorf("Failed to send message to "+
						"%s: %v", p, err)
				}
				if msg.doneChan != nil {
					msg.doneChan <- struct{}{}
				}
				continue
			}

			// At this point, the message was successfully sent, so
			// update the last send time, signal the sender of the
			// message that it has been sent (if requested), and
			// signal the send queue to the deliver the next queued
			// message.
			atomic.StoreInt64(&p.lastSend, time.Now().Unix())
			if msg.doneChan != nil {
				msg.doneChan <- struct{}{}
			}
			p.sendDoneQueue <- struct{}{}

		case <-pingTicker.C:
			nonce, err := wire.RandomUint64()
			if err != nil {
				log.Errorf("Not sending ping to %s: %v", p, err)
				continue
			}
			p.QueueMessage(wire.NewMsgPing(nonce), nil)

		case <-p.quit:
			break out
		}
	}

	<-p.queueQuit

	// Drain any wait channels before we go away so we don't leave something
	// waiting for us. We have waited on queueQuit and thus we can be sure
	// that we will not miss anything sent on sendQueue.
cleanup:
	for {
		select {
		case msg := <-p.sendQueue:
			if msg.doneChan != nil {
				msg.doneChan <- struct{}{}
			}
			// no need to send on sendDoneQueue since queueHandler
			// has been waited on and already exited.
		default:
			break cleanup
		}
	}
	close(p.outQuit)
	log.Tracef("WitnessPeer output handler done for %s", p)
}

// QueueMessage adds the passed wire message to the WitnessPeer send queue.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) QueueMessage(msg wire.Message, doneChan chan<- struct{}) {
	// Avoid risk of deadlock if goroutine already exited.  The goroutine
	// we will be sending to hangs around until it knows for a fact that
	// it is marked as disconnected and *then* it drains the channels.
	if !p.Connected() {
		if doneChan != nil {
			go func() {
				doneChan <- struct{}{}
			}()
		}
		return
	}
	p.outputQueue <- outMsg{msg: msg, doneChan: doneChan}
}

// QueueInventory adds the passed inventory to the inventory send queue which
// might not be sent right away, rather it is trickled to the WitnessPeer in batches.
// Inventory that the WitnessPeer is already known to have is ignored.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) QueueInventory(invVect *wire.InvVect) {
	// Don't add the inventory to the send queue if the WitnessPeer is already
	// known to have it.
	if p.knownInventory.Exists(invVect) {
		return
	}

	// Avoid risk of deadlock if goroutine already exited.  The goroutine
	// we will be sending to hangs around until it knows for a fact that
	// it is marked as disconnected and *then* it drains the channels.
	if !p.Connected() {
		return
	}

	p.outputInvChan <- invVect
}

// AssociateConnection associates the given conn to the WitnessPeer.   Calling this
// function when the WitnessPeer is already connected will have no effect.
func (p *WitnessPeer) AssociateConnection(conn net.Conn) {
	// Already connected?
	if !atomic.CompareAndSwapInt32(&p.connected, 0, 1) {
		return
	}

	p.conn = conn
	p.timeConnected = time.Now()

	if p.inbound {
		p.addr = p.conn.RemoteAddr().String()

		// Set up a NetAddress for the WitnessPeer to be used with AddrManager.  We
		// only do this inbound because outbound set this up at connection time
		// and no point recomputing.
		na, err := newNetAddress(p.conn.RemoteAddr(), p.services)
		if err != nil {
			log.Errorf("Cannot create remote net address: %v", err)
			p.Disconnect()
			return
		}
		p.na = na
	}

	go func(WitnessPeer *WitnessPeer) {
		if err := WitnessPeer.start(); err != nil {
			log.Debugf("Cannot start WitnessPeer %v: %v", WitnessPeer, err)
			WitnessPeer.Disconnect()
		}
	}(p)
}

// Connected returns whether or not the WitnessPeer is currently connected.
//
// This function is safe for concurrent access.
func (p *WitnessPeer) Connected() bool {
	return atomic.LoadInt32(&p.connected) != 0 &&
		atomic.LoadInt32(&p.disconnect) == 0
}

// Disconnect disconnects the WitnessPeer by closing the connection.  Calling this
// function when the WitnessPeer is already disconnected or in the process of
// disconnecting will have no effect.
func (p *WitnessPeer) Disconnect() {
	if atomic.AddInt32(&p.disconnect, 1) != 1 {
		return
	}

	log.Tracef("Disconnecting %s", p)
	if atomic.LoadInt32(&p.connected) != 0 {
		p.conn.Close()
	}
	close(p.quit)
}

// start begins processing input and output messages.
func (p *WitnessPeer) start() error {
	log.Tracef("Starting WitnessPeer %s", p)

	negotiateErr := make(chan error)
	go func() {
		if p.inbound {
			negotiateErr <- p.negotiateInboundProtocol()
		} else {
			negotiateErr <- p.negotiateOutboundProtocol()
		}
	}()

	// Negotiate the protocol within the specified negotiateTimeout.
	select {
	case err := <-negotiateErr:
		if err != nil {
			return err
		}
	case <-time.After(negotiateTimeout):
		return errors.New("protocol negotiation timeout")
	}
	log.Debugf("Connected to %s", p.Addr())

	// The protocol has been negotiated successfully so start processing input
	// and output messages.
	go p.stallHandler()
	go p.inHandler()
	go p.queueHandler()
	go p.outHandler()

	// Send our verack message now that the IO processing machinery has started.
	p.QueueMessage(wire.NewMsgVerAck(), nil)
	return nil
}

// WaitForDisconnect waits until the WitnessPeer has completely disconnected and all
// resources are cleaned up.  This will happen if either the local or remote
// side has been disconnected or the WitnessPeer is forcibly disconnected via
// Disconnect.
func (p *WitnessPeer) WaitForDisconnect() {
	<-p.quit
}

// readRemoteVersionMsg waits for the next message to arrive from the remote
// WitnessPeer.  If the next message is not a version message or the version is not
// acceptable then return an error.
func (p *WitnessPeer) readRemoteVersionMsg() error {
	log.Tracef("readRemoteVersionMsg %s", p)
	// Read their version message.
	msg, _, err := p.readMessage()
	if err != nil {
		return err
	}

	remoteVerMsg, ok := msg.(*wire.MsgVersion)
	if !ok {
		errStr := "A version message must precede all others"
		log.Errorf(errStr)

		rejectMsg := wire.NewMsgReject(msg.Command(), wire.RejectMalformed,
			errStr)
		return p.writeMessage(rejectMsg)
	}

	if err := p.handleRemoteVersionMsg(remoteVerMsg); err != nil {
		return err
	}

	if p.cfg.Listeners.OnVersion != nil {
		p.cfg.Listeners.OnVersion(p, remoteVerMsg)
	}
	return nil
}

// writeLocalVersionMsg writes our version message to the remote WitnessPeer.
func (p *WitnessPeer) writeLocalVersionMsg() error {
	localVerMsg, err := p.localVersionMsg()
	if err != nil {
		return err
	}

	if err := p.writeMessage(localVerMsg); err != nil {
		return err
	}

	p.flagsMtx.Lock()
	p.versionSent = true
	p.flagsMtx.Unlock()
	return nil
}

// negotiateInboundProtocol waits to receive a version message from the WitnessPeer
// then sends our version message. If the events do not occur in that order then
// it returns an error.
func (p *WitnessPeer) negotiateInboundProtocol() error {
	if err := p.readRemoteVersionMsg(); err != nil {
		return err
	}

	return p.writeLocalVersionMsg()
}

// negotiateOutboundProtocol sends our version message then waits to receive a
// version message from the WitnessPeer.  If the events do not occur in that order then
// it returns an error.
func (p *WitnessPeer) negotiateOutboundProtocol() error {
	if err := p.writeLocalVersionMsg(); err != nil {
		return err
	}

	return p.readRemoteVersionMsg()
}

// newWitnessPeerBase returns a new base hcd WitnessPeer based on the inbound flag.  This
// is used by the NewInboundWitnessPeer and NewOutboundWitnessPeer functions to perform base
// setup needed by both types of WitnessPeers.
func newWitnessPeerBase(cfg *WitnessConfig, inbound bool) *WitnessPeer {
	// Default to the max supported protocol version.  Override to the
	// version specified by the caller if configured.
	protocolVersion := MaxWitnessProtocolVersion
	if cfg.ProtocolVersion != 0 {
		protocolVersion = cfg.ProtocolVersion
	}

	// Set the chain parameters to testnet if the caller did not specify any.
	if cfg.ChainParams == nil {
		cfg.ChainParams = &chaincfg.TestNet2Params
	}

	p := WitnessPeer{
		inbound:         inbound,
		knownInventory:  newMruInventoryMap(maxKnownInventory),
		stallControl:    make(chan stallControlMsg, 1), // nonblocking sync
		outputQueue:     make(chan outMsg, outputBufferSize),
		sendQueue:       make(chan outMsg, 1),   // nonblocking sync
		sendDoneQueue:   make(chan struct{}, 1), // nonblocking sync
		outputInvChan:   make(chan *wire.InvVect, outputBufferSize),
		inQuit:          make(chan struct{}),
		queueQuit:       make(chan struct{}),
		outQuit:         make(chan struct{}),
		quit:            make(chan struct{}),
		cfg:             *cfg, // Copy so caller can't mutate.
		services:        cfg.Services,
		protocolVersion: protocolVersion,
	}
	return &p
}

// NewInboundWitnessPeer returns a new inbound hcd WitnessPeer. Use Start to begin
// processing incoming and outgoing messages.
func NewInboundWitnessPeer(cfg *WitnessConfig) *WitnessPeer {
	return newWitnessPeerBase(cfg, true)
}

// NewOutboundWitnessPeer returns a new outbound hcd WitnessPeer.
func NewOutboundWitnessPeer(cfg *WitnessConfig, addr string) (*WitnessPeer, error) {
	p := newWitnessPeerBase(cfg, false)
	p.addr = addr

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}

	if cfg.HostToNetAddress != nil {
		na, err := cfg.HostToNetAddress(host, uint16(port), 0)
		if err != nil {
			return nil, err
		}
		p.na = na
	} else {
		p.na = wire.NewNetAddressIPPort(net.ParseIP(host), uint16(port),
			0)
	}

	return p, nil
}
