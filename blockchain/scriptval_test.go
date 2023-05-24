// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain_test

//	"fmt"
//	"runtime"
import (
	"testing"
)

//	"github.com/james-ray/hcd/blockchain"
//	"github.com/james-ray/hcd/txscript"

// TestCheckBlockScripts ensures that validating the all of the scripts in a
// known-good block doesn't return an error.
func TestCheckBlockScripts(t *testing.T) {
	/*
		// TODO In the future, add a block here with a lot of tx to validate.
		// The blockchain tests already validate a ton of scripts with signatures,
		// so we don't really need to make a new test for this immediately.
		runtime.GOMAXPROCS(runtime.NumCPU())

		testBlockNum := 277647
		blockDataFile := fmt.Sprintf("%d.dat.bz2", testBlockNum)
		blocks, err := loadBlocks(blockDataFile)
		if err != nil {
			t.Errorf("Error loading file: %v\n", err)
			return
		}
		if len(blocks) > 1 {
			t.Errorf("The test block file must only have one block in it")
			return
		}
		if len(blocks) == 0 {
			t.Errorf("The test block file may not be empty")
			return
		}

		storeDataFile := fmt.Sprintf("%d.utxostore.bz2", testBlockNum)
		view, err := loadUtxoView(storeDataFile)
		if err != nil {
			t.Errorf("Error loading txstore: %v\n", err)
			return
		}

		scriptFlags := txscript.ScriptBip16
		err = blockchain.TstCheckBlockScripts(blocks[0], view, scriptFlags,
			nil)
		if err != nil {
			t.Errorf("Transaction script validation failed: %v\n", err)
			return
		}
	*/
}
