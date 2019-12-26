// Copyright (c) 2014-2015 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Copyright (c) 2018-2020 The Hcd developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/HcashOrg/hcd/hcutil"
	hcrpcclient "github.com/HcashOrg/hcd/hcrpcclient"
)

func main() {
	// Only override the handlers for notifications you care about.
	// Also note most of these handlers will only be called if you register
	// for notifications.  See the documentation of the hcrpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := hcrpcclient.NotificationHandlers{
		OnBlockConnected: func(blockHeader []byte, transactions [][]byte) {
			log.Printf("Block connected: %v %v", blockHeader, transactions)
		},
		OnBlockDisconnected: func(blockHeader []byte) {
			log.Printf("Block disconnected: %v", blockHeader)
		},
	}

	// Connect to local hcd RPC server using websockets.
	hcdHomeDir := hcutil.AppDataDir("hcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(hcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &hcrpcclient.ConnConfig{
		Host:         "localhost:14009", //"localhost:12009", testnet
		Endpoint:     "ws",
		User:         "admin", //"yourrpcuser",
		Pass:         "123",   //"yourrpcpass",
		Certificates: certs,
	}
	client, err := hcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	// Register for block connect and disconnect notifications.
	if err := client.NotifyBlocks(); err != nil {
		log.Fatal(err)
	}
	log.Println("NotifyBlocks: Registration Complete")

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)

	fDifficulty, err := client.GetDifficulty()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("fDifficulty : %f", fDifficulty)

	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.
	log.Println("Client shutdown in 10 seconds...")
	time.AfterFunc(time.Second*10, func() {
		log.Println("Client shutting down...")
		client.Shutdown()
		log.Println("Client shutdown complete.")
	})

	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	client.WaitForShutdown()
}
