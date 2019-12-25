package peer_test

import (
	"errors"
	"github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/peer"
	"github.com/HcashOrg/hcd/wire"
	"io"
	"sync"
	"testing"
	"time"
)

func init() {
	peer.TstAllowSelfWitnessConns()
}

// testWitnessPeer tests the given witness peer's flags and stats
func testWitnessPeer(t *testing.T, p *peer.WitnessPeer, s peerStats) {
	if p.UserAgent() != s.wantUserAgent {
		t.Errorf("testPeer: wrong UserAgent - got %v, want %v", p.UserAgent(), s.wantUserAgent)
		return
	}

	if p.Services() != s.wantServices {
		t.Errorf("testPeer: wrong Services - got %v, want %v", p.Services(), s.wantServices)
		return
	}

	if !p.LastPingTime().Equal(s.wantLastPingTime) {
		t.Errorf("testPeer: wrong LastPingTime - got %v, want %v", p.LastPingTime(), s.wantLastPingTime)
		return
	}

	if p.LastPingNonce() != s.wantLastPingNonce {
		t.Errorf("testPeer: wrong LastPingNonce - got %v, want %v", p.LastPingNonce(), s.wantLastPingNonce)
		return
	}

	if p.LastPingMicros() != s.wantLastPingMicros {
		t.Errorf("testPeer: wrong LastPingMicros - got %v, want %v", p.LastPingMicros(), s.wantLastPingMicros)
		return
	}

	if p.VerAckReceived() != s.wantVerAckReceived {
		t.Errorf("testPeer: wrong VerAckReceived - got %v, want %v", p.VerAckReceived(), s.wantVerAckReceived)
		return
	}

	if p.VersionKnown() != s.wantVersionKnown {
		t.Errorf("testPeer: wrong VersionKnown - got %v, want %v", p.VersionKnown(), s.wantVersionKnown)
		return
	}

	if p.ProtocolVersion() != s.wantProtocolVersion {
		t.Errorf("testPeer: wrong ProtocolVersion - got %v, want %v", p.ProtocolVersion(), s.wantProtocolVersion)
		return
	}

	// Allow for a deviation of 1s, as the second may tick when the message is
	// in transit and the protocol doesn't support any further precision.
	if p.TimeOffset() != s.wantTimeOffset && p.TimeOffset() != s.wantTimeOffset-1 {
		t.Errorf("testPeer: wrong TimeOffset - got %v, want %v or %v", p.TimeOffset(),
			s.wantTimeOffset, s.wantTimeOffset-1)
		return
	}

	if p.BytesSent() != s.wantBytesSent {
		t.Errorf("testPeer: wrong BytesSent - got %v, want %v", p.BytesSent(), s.wantBytesSent)
		return
	}

	if p.BytesReceived() != s.wantBytesReceived {
		t.Errorf("testPeer: wrong BytesReceived - got %v, want %v", p.BytesReceived(), s.wantBytesReceived)
		return
	}

	if p.StartingHeight() != s.wantStartingHeight {
		t.Errorf("testPeer: wrong StartingHeight - got %v, want %v", p.StartingHeight(), s.wantStartingHeight)
		return
	}

	if p.Connected() != s.wantConnected {
		t.Errorf("testPeer: wrong Connected - got %v, want %v", p.Connected(), s.wantConnected)
		return
	}

	stats := p.WitnessStatsSnapshot()

	if p.ID() != stats.ID {
		t.Errorf("testPeer: wrong ID - got %v, want %v", p.ID(), stats.ID)
		return
	}

	if p.Addr() != stats.Addr {
		t.Errorf("testPeer: wrong Addr - got %v, want %v", p.Addr(), stats.Addr)
		return
	}

	if p.LastSend() != stats.LastSend {
		t.Errorf("testPeer: wrong LastSend - got %v, want %v", p.LastSend(), stats.LastSend)
		return
	}

	if p.LastRecv() != stats.LastRecv {
		t.Errorf("testPeer: wrong LastRecv - got %v, want %v", p.LastRecv(), stats.LastRecv)
		return
	}
}

// TestWitnessPeerConnection tests connection between inbound and outbound witness peers.
func TestWitnessPeerConnection(t *testing.T) {
	var pause sync.Mutex
	verack := make(chan struct{})
	peerCfg := &peer.WitnessConfig{
		Listeners: peer.WitnessMessageListeners{
			OnVerAck: func(p *peer.WitnessPeer, msg *wire.MsgVerAck) {
				verack <- struct{}{}
			},
			OnWrite: func(p *peer.WitnessPeer, bytesWritten int, msg wire.Message,
				err error) {
				if _, ok := msg.(*wire.MsgVerAck); ok {
					verack <- struct{}{}
				}
				pause.Lock()
				pause.Unlock()
			},
		},
		UserAgentName:    "witnesspeer",
		UserAgentVersion: "1.0",
		ChainParams:      &chaincfg.MainNetParams,
		Services:         0,
	}
	wantStats := peerStats{
		wantUserAgent:       wire.DefaultUserAgent + "witnesspeer:1.0/",
		wantServices:        0,
		wantProtocolVersion: peer.MaxWitnessProtocolVersion,
		wantConnected:       true,
		wantVersionKnown:    true,
		wantVerAckReceived:  true,
		wantLastPingTime:    time.Time{},
		wantLastPingNonce:   uint64(0),
		wantLastPingMicros:  int64(0),
		wantTimeOffset:      int64(0),
		wantBytesSent:       164, // 140 version + 24 verack
		wantBytesReceived:   164,
	}
	tests := []struct {
		name  string
		setup func() (*peer.WitnessPeer, *peer.WitnessPeer, error)
	}{
		{
			"basic handshake",
			func() (*peer.WitnessPeer, *peer.WitnessPeer, error) {
				inConn, outConn := pipe(
					&conn{raddr: "10.0.0.1:8333"},
					&conn{raddr: "10.0.0.2:8333"},
				)
				inPeer := peer.NewInboundWitnessPeer(peerCfg)
				inPeer.AssociateConnection(inConn)

				outPeer, err := peer.NewOutboundWitnessPeer(peerCfg, "10.0.0.2:8333")
				if err != nil {
					return nil, nil, err
				}
				outPeer.AssociateConnection(outConn)

				for i := 0; i < 4; i++ {
					select {
					case <-verack:
					case <-time.After(time.Second):
						return nil, nil, errors.New("verack timeout")
					}
				}
				return inPeer, outPeer, nil
			},
		},
		{
			"socks proxy",
			func() (*peer.WitnessPeer, *peer.WitnessPeer, error) {
				inConn, outConn := pipe(
					&conn{raddr: "10.0.0.1:8333", proxy: true},
					&conn{raddr: "10.0.0.2:8333"},
				)
				inPeer := peer.NewInboundWitnessPeer(peerCfg)
				inPeer.AssociateConnection(inConn)

				outPeer, err := peer.NewOutboundWitnessPeer(peerCfg, "10.0.0.2:8333")
				if err != nil {
					return nil, nil, err
				}
				outPeer.AssociateConnection(outConn)

				for i := 0; i < 4; i++ {
					select {
					case <-verack:
					case <-time.After(time.Second):
						return nil, nil, errors.New("verack timeout")
					}
				}
				return inPeer, outPeer, nil
			},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		inPeer, outPeer, err := test.setup()
		if err != nil {
			t.Errorf("TestPeerConnection setup #%d: unexpected err %v", i, err)
			return
		}

		pause.Lock()
		testWitnessPeer(t, inPeer, wantStats)
		testWitnessPeer(t, outPeer, wantStats)
		pause.Unlock()

		inPeer.Disconnect()
		outPeer.Disconnect()
		inPeer.WaitForDisconnect()
		outPeer.WaitForDisconnect()
	}
}

// TestWitnessPeerListeners tests that the witness peer listeners are called as expected.
func TestWitnessPeerListeners(t *testing.T) {
	verack := make(chan struct{}, 1)
	version := make(chan wire.Message, 1)
	ok := make(chan wire.Message, 15)

	peerCfg := &peer.WitnessConfig{
		Listeners: peer.WitnessMessageListeners{
			OnGetAddr: func(p *peer.WitnessPeer, msg *wire.MsgGetAddr) {
				ok <- msg
			},
			OnAddr: func(p *peer.WitnessPeer, msg *wire.MsgAddr) {
				ok <- msg
			},
			OnPing: func(p *peer.WitnessPeer, msg *wire.MsgPing) {
				ok <- msg
			},
			OnPong: func(p *peer.WitnessPeer, msg *wire.MsgPong) {
				ok <- msg
			},
			OnAlert: func(p *peer.WitnessPeer, msg *wire.MsgAlert) {
				ok <- msg
			},
			OnTx: func(p *peer.WitnessPeer, msg *wire.MsgTx) {
				ok <- msg
			},

			OnInv: func(p *peer.WitnessPeer, msg *wire.MsgInv) {
				ok <- msg
			},

			OnNotFound: func(p *peer.WitnessPeer, msg *wire.MsgNotFound) {
				ok <- msg
			},
			OnGetData: func(p *peer.WitnessPeer, msg *wire.MsgGetData) {
				ok <- msg
			},

			OnFeeFilter: func(p *peer.WitnessPeer, msg *wire.MsgFeeFilter) {
				ok <- msg
			},
			OnFilterAdd: func(p *peer.WitnessPeer, msg *wire.MsgFilterAdd) {
				ok <- msg
			},
			OnFilterClear: func(p *peer.WitnessPeer, msg *wire.MsgFilterClear) {
				ok <- msg
			},
			OnFilterLoad: func(p *peer.WitnessPeer, msg *wire.MsgFilterLoad) {
				ok <- msg
			},
			OnVersion: func(p *peer.WitnessPeer, msg *wire.MsgVersion) {
				version <- msg
			},
			OnVerAck: func(p *peer.WitnessPeer, msg *wire.MsgVerAck) {
				verack <- struct{}{}
			},
			OnReject: func(p *peer.WitnessPeer, msg *wire.MsgReject) {
				ok <- msg
			},
		},
		UserAgentName:    "witnesspeer",
		UserAgentVersion: "1.0",
		ChainParams:      &chaincfg.MainNetParams,
		Services:         wire.SFNodeBloom,
	}
	inConn, outConn := pipe(
		&conn{raddr: "10.0.0.1:8333"},
		&conn{raddr: "10.0.0.2:8333"},
	)
	inPeer := peer.NewInboundWitnessPeer(peerCfg)
	inPeer.AssociateConnection(inConn)

	peerCfg.Listeners = peer.WitnessMessageListeners{
		OnVerAck: func(p *peer.WitnessPeer, msg *wire.MsgVerAck) {
			verack <- struct{}{}
		},
	}
	outPeer, err := peer.NewOutboundWitnessPeer(peerCfg, "10.0.0.1:8333")
	if err != nil {
		t.Errorf("NewOutboundPeer: unexpected err %v\n", err)
		return
	}
	outPeer.AssociateConnection(outConn)

	select {
	case <-version:
	case <-time.After(time.Second * 1):
		t.Errorf("TestPeerListeners: verack timeout\n")
		return
	}

	for i := 0; i < 2; i++ {
		select {
		case <-verack:
		case <-time.After(time.Second * 1):
			t.Errorf("TestPeerListeners: verack timeout\n")
			return
		}
	}

	tests := []struct {
		listener string
		msg      wire.Message
	}{
		{
			"OnGetAddr",
			wire.NewMsgGetAddr(),
		},
		{
			"OnAddr",
			wire.NewMsgAddr(),
		},
		{
			"OnPing",
			wire.NewMsgPing(42),
		},
		{
			"OnPong",
			wire.NewMsgPong(42),
		},
		{
			"OnAlert",
			wire.NewMsgAlert([]byte("payload"), []byte("signature")),
		},
		{
			"OnTx",
			wire.NewMsgTx(),
		},
		{
			"OnInv",
			wire.NewMsgInv(),
		},
		{
			"OnNotFound",
			wire.NewMsgNotFound(),
		},
		{
			"OnGetData",
			wire.NewMsgGetData(),
		},
		{
			"OnFeeFilter",
			wire.NewMsgFeeFilter(15000),
		},
		{
			"OnFilterAdd",
			wire.NewMsgFilterAdd([]byte{0x01}),
		},
		{
			"OnFilterClear",
			wire.NewMsgFilterClear(),
		},
		{
			"OnFilterLoad",
			wire.NewMsgFilterLoad([]byte{0x01}, 10, 0, wire.BloomUpdateNone),
		},
		// only one version message is allowed
		// only one verack message is allowed
		{
			"OnReject",
			wire.NewMsgReject("block", wire.RejectDuplicate, "dupe block"),
		},
	}
	t.Logf("Running %d tests", len(tests))
	for _, test := range tests {
		// Queue the test message
		outPeer.QueueMessage(test.msg, nil)
		//wait until ok or fail , then next test
		select {
		case msg := <-ok:
			t.Log("ok", msg.Command(), test.listener)
		case <-time.After(time.Second * 1):
			t.Errorf("TestPeerListeners: %s timeout", test.listener)
			return
		}
	}
	inPeer.Disconnect()
	outPeer.Disconnect()
}

// TestOutboundWitnessPeer tests that the outbound peer works as expected.
func TestOutboundWitnessPeer(t *testing.T) {
	peerCfg := &peer.WitnessConfig{
		NewestBlock: func() (*chainhash.Hash, int64, error) {
			return nil, 0, errors.New("newest block not found")
		},
		UserAgentName:    "witnesspeer",
		UserAgentVersion: "1.0",
		ChainParams:      &chaincfg.MainNetParams,
		Services:         0,
	}

	//test connect
	r, w := io.Pipe()
	c := &conn{raddr: "10.0.0.1:8333", Writer: w, Reader: r}

	p, err := peer.NewOutboundWitnessPeer(peerCfg, "10.0.0.1:8333")
	if err != nil {
		t.Errorf("NewOutboundPeer: unexpected err - %v\n", err)
		return
	}

	// Test trying to connect twice.
	p.AssociateConnection(c)
	p.AssociateConnection(c)


	disconnected := make(chan struct{})
	go func() {
		p.WaitForDisconnect()
		disconnected <- struct{}{}
	}()

	select {
	case <-disconnected:
		close(disconnected)
	case <-time.After(time.Second):
		t.Fatal("Peer did not automatically disconnect.")
	}

	if p.Connected() {
		t.Fatalf("Should not be connected as c.close error.")
	}

	// Test Queue Inv
	fakeTxHash := &chainhash.Hash{0: 0x00, 1: 0x01}

	fakeInv := wire.NewInvVect(wire.InvTypeTx, fakeTxHash)

	// Should be noops as the peer could not connect.
	p.QueueInventory(fakeInv)
	p.AddKnownWitnessInventory(fakeInv)
	p.QueueInventory(fakeInv)

	fakeMsg := wire.NewMsgVerAck()
	p.QueueMessage(fakeMsg, nil)
	done := make(chan struct{})
	p.QueueMessage(fakeMsg, done)
	<-done
	p.Disconnect()

	r1, w1 := io.Pipe()
	c1 := &conn{raddr: "10.0.0.1:8333", Writer: w1, Reader: r1}
	p1, err := peer.NewOutboundWitnessPeer(peerCfg, "10.0.0.1:8333")
	if err != nil {
		t.Errorf("NewOutboundPeer: unexpected err - %v\n", err)
		return
	}
	p1.AssociateConnection(c1)

	// Test Queue Inv after connection
	p1.QueueInventory(fakeInv)
	p1.Disconnect()

	// Test testnet
	peerCfg.ChainParams = &chaincfg.TestNet2Params
	peerCfg.Services = wire.SFNodeBloom
	r2, w2 := io.Pipe()
	c2 := &conn{raddr: "10.0.0.1:8333", Writer: w2, Reader: r2}
	p2, err := peer.NewOutboundWitnessPeer(peerCfg, "10.0.0.1:8333")
	if err != nil {
		t.Errorf("NewOutboundPeer: unexpected err - %v\n", err)
		return
	}
	p2.AssociateConnection(c2)

	// Test PushXXX
	var addrs []*wire.NetAddress
	for i := 0; i < 5; i++ {
		na := wire.NetAddress{}
		addrs = append(addrs, &na)
	}
	if _, err := p2.PushAddrMsg(addrs); err != nil {
		t.Errorf("PushAddrMsg: unexpected err %v\n", err)
		return
	}

	p2.PushRejectMsg("block", wire.RejectMalformed, "malformed", nil, false)
	p2.PushRejectMsg("block", wire.RejectInvalid, "invalid", nil, false)

	// Test Queue Messages
	p2.QueueMessage(wire.NewMsgGetAddr(), nil)
	p2.QueueMessage(wire.NewMsgPing(1), nil)
	p2.QueueMessage(wire.NewMsgMemPool(), nil)
	p2.QueueMessage(wire.NewMsgGetData(), nil)
	p2.QueueMessage(wire.NewMsgGetHeaders(), nil)
	p2.QueueMessage(wire.NewMsgFeeFilter(20000), nil)

	p2.Disconnect()
}
