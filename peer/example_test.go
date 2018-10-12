// Copyright (c) 2015-2016 The btcsuite developers
// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package peer_test

import (
	"fmt"
	"net"

	"github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/peer"
)

// mockRemotePeer creates a basic inbound peer listening on the simnet port for
// use with Example_peerConnection.  It does not return until the listner is
// active.
func mockRemotePeer() error {
	// Configure peer to act as a simnet node that offers no services.
	peerCfg := &peer.Config{
		UserAgentName:    "peer",  // User agent name to advertise.
		UserAgentVersion: "1.0.0", // User agent version to advertise.
		ChainParams:      &chaincfg.SimNetParams,
	}

	// Accept connections on the simnet port.
	listener, err := net.Listen("tcp", "127.0.0.1:18555")
	if err != nil {
		return err
	}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept: error %v\n", err)
			return
		}

		// Create and start the inbound peer.
		p := peer.NewInboundPeer(peerCfg)
		p.AssociateConnection(conn)
	}()

	return nil
}
