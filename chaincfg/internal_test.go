// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"testing"

	"github.com/HcashOrg/hcd/chaincfg/chainhash"
)

// TestInvalidHashStr tests against the NewHashFromStr function.
func TestInvalidHashStr(t *testing.T) {
	_, err := chainhash.NewHashFromStr("banana")
	if err == nil {
		t.Error("Invalid string should fail.")
	}
}

//TestIsPubKeyHashAddrId
func TestIsPubKeyHashAddrID(t *testing.T) {
	is:=IsPubKeyAddrID([2]byte{34,28})
	if !is {
		t.Log("not pubKeyHashAddr")
	}
}
