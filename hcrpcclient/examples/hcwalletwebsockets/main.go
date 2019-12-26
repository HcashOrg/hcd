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
	"github.com/davecgh/go-spew/spew"
)

func main() {
	// Only override the handlers for notifications you care about.
	// Also note most of the handlers will only be called if you register
	// for notifications.  See the documentation of the hcrpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := hcrpcclient.NotificationHandlers{
		OnAccountBalance: func(account string, balance hcutil.Amount, confirmed bool) {
			log.Printf("New balance for account %s: %v", account,
				balance)
		},
	}

	// Connect to local hcwallet RPC server using websockets.
	certHomeDir := hcutil.AppDataDir("hcwallet", false)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &hcrpcclient.ConnConfig{
		Host:         "localhost:12010", //testnet "localhost:12010",
		Endpoint:     "ws",
		User:         "admin", //"yourrpcuser",
		Pass:         "123",   //"yourrpcpass",
		Certificates: certs,
	}
	client, err := hcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	// Get the list of unspent transaction outputs (utxos) that the
	// connected wallet has at least one private key for.
	unspent, err := client.ListUnspent()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Num unspent outputs (utxos): %d", len(unspent))
	if len(unspent) > 0 {
		log.Printf("First utxo:\n%v", spew.Sdump(unspent[0]))
	}

	mp := make(map[hcutil.Address]hcutil.Amount)
	addr1,_ := hcutil.DecodeAddress("TsRj7wFGcFrpWhWztm4kdG5g1WZWonoFtXm")
	mp[addr1] = hcutil.Amount(200000000)//2 HCD


	mp[addr1] = hcutil.Amount(200000000)//2 HCD
	hash, err := client.SendManyV2("default",mp)
	log.Println("hash:",hash,err)

	hash, err = client.SendManyV2ChangeAddr("default",mp,"TbMpHSpnkw2CVi6cbTLcMYydeGn3LnoZG8q")
	log.Println("hash:",hash,err)

	
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
