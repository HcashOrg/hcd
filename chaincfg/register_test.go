// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg_test

import (
	"bytes"
	"reflect"
	"testing"
)

// Define some of the required parameters for a user-registered
// network.  This is necessary to test the registration of and
// lookup of encoding magics from the network.
var mockNetParams = Params{
	Name:             "mocknet",
	Net:              1<<32 - 1,
	PubKeyHashAddrID: [2]byte{0x9f},
	ScriptHashAddrID: [2]byte{0xf9},
	HDPrivateKeyID:   [4]byte{0x01, 0x02, 0x03, 0x04},
	HDPublicKeyID:    [4]byte{0x05, 0x06, 0x07, 0x08},
}

// TestRegister test registered network
func TestRegister(t *testing.T) {
	type registerTest struct {
		name   string
		params *Params
		err    error
	}
	type magicTest struct {
		magic [2]byte
		valid bool
	}
	type hdTest struct {
		priv []byte
		want []byte
		err  error
	}

	tests := []struct {
		name        string
		register    []registerTest
		p2pkhMagics []magicTest
		p2shMagics  []magicTest
		hdMagics    []hdTest
	}{
		{
			name: "default networks",
			register: []registerTest{
				{
					name:   "duplicate mainnet",
					params: &MainNetParams,
					err:    ErrDuplicateNet,
				},
				{
					name:   "duplicate testnet",
					params: &TestNet2Params,
					err:    ErrDuplicateNet,
				},
				{
					name:   "duplicate simnet",
					params: &SimNetParams,
					err:    ErrDuplicateNet,
				},
			},
			p2pkhMagics: []magicTest{
				{
					magic: MainNetParams.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: TestNet2Params.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: SimNetParams.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: mockNetParams.PubKeyHashAddrID,
					valid: false,
				},
				{
					magic: [2]byte{0xFF},
					valid: false,
				},
			},
			p2shMagics: []magicTest{
				{
					magic: MainNetParams.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: TestNet2Params.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: SimNetParams.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: mockNetParams.ScriptHashAddrID,
					valid: false,
				},
				{
					magic: [2]byte{0xFF},
					valid: false,
				},
			},
			hdMagics: []hdTest{
				{
					priv: MainNetParams.HDPrivateKeyID[:],
					want: MainNetParams.HDPublicKeyID[:],
					err:  nil,
				},
				{
					priv: TestNet2Params.HDPrivateKeyID[:],
					want: TestNet2Params.HDPublicKeyID[:],
					err:  nil,
				},
				{
					priv: SimNetParams.HDPrivateKeyID[:],
					want: SimNetParams.HDPublicKeyID[:],
					err:  nil,
				},
				{
					priv: mockNetParams.HDPrivateKeyID[:],
					err:  ErrUnknownHDKeyID,
				},
				{
					priv: []byte{0xff, 0xff, 0xff, 0xff},
					err:  ErrUnknownHDKeyID,
				},
				{
					priv: []byte{0xff},
					err:  ErrUnknownHDKeyID,
				},
			},
		},
		{
			name: "register mocknet",
			register: []registerTest{
				{
					name:   "mocknet",
					params: &mockNetParams,
					err:    nil,
				},
			},
			p2pkhMagics: []magicTest{
				{
					magic: MainNetParams.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: TestNet2Params.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: SimNetParams.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: mockNetParams.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: [2]byte{0xFF},
					valid: false,
				},
			},
			p2shMagics: []magicTest{
				{
					magic: MainNetParams.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: TestNet2Params.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: SimNetParams.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: mockNetParams.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: [2]byte{0xFF},
					valid: false,
				},
			},
			hdMagics: []hdTest{
				{
					priv: mockNetParams.HDPrivateKeyID[:],
					want: mockNetParams.HDPublicKeyID[:],
					err:  nil,
				},
			},
		},
		{
			name: "more duplicates",
			register: []registerTest{
				{
					name:   "duplicate mainnet",
					params: &MainNetParams,
					err:    ErrDuplicateNet,
				},
				{
					name:   "duplicate testnet",
					params: &TestNet2Params,
					err:    ErrDuplicateNet,
				},
				{
					name:   "duplicate simnet",
					params: &SimNetParams,
					err:    ErrDuplicateNet,
				},
				{
					name:   "duplicate mocknet",
					params: &mockNetParams,
					err:    ErrDuplicateNet,
				},
			},
			p2pkhMagics: []magicTest{
				{
					magic: MainNetParams.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: TestNet2Params.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: SimNetParams.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: mockNetParams.PubKeyHashAddrID,
					valid: true,
				},
				{
					magic: [2]byte{0xFF},
					valid: false,
				},
			},
			p2shMagics: []magicTest{
				{
					magic: MainNetParams.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: TestNet2Params.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: SimNetParams.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: mockNetParams.ScriptHashAddrID,
					valid: true,
				},
				{
					magic: [2]byte{0xFF},
					valid: false,
				},
			},
			hdMagics: []hdTest{
				{
					priv: MainNetParams.HDPrivateKeyID[:],
					want: MainNetParams.HDPublicKeyID[:],
					err:  nil,
				},
				{
					priv: TestNet2Params.HDPrivateKeyID[:],
					want: TestNet2Params.HDPublicKeyID[:],
					err:  nil,
				},
				{
					priv: SimNetParams.HDPrivateKeyID[:],
					want: SimNetParams.HDPublicKeyID[:],
					err:  nil,
				},
				{
					priv: mockNetParams.HDPrivateKeyID[:],
					want: mockNetParams.HDPublicKeyID[:],
					err:  nil,
				},
				{
					priv: []byte{0xff, 0xff, 0xff, 0xff},
					err:  ErrUnknownHDKeyID,
				},
				{
					priv: []byte{0xff},
					err:  ErrUnknownHDKeyID,
				},
			},
		},
	}

	for _, test := range tests {
		for _, regTest := range test.register {
			err := Register(regTest.params)
			if err != regTest.err {
				t.Errorf("%s:%s: Registered network with unexpected error: got %v expected %v",
					test.name, regTest.name, err, regTest.err)
			}
		}
		for i, magTest := range test.p2pkhMagics {
			valid := IsPubKeyHashAddrID(magTest.magic)
			if valid != magTest.valid {
				t.Errorf("%s: P2PKH magic %d valid mismatch: got %v expected %v",
					test.name, i, valid, magTest.valid)
			}
		}
		for i, magTest := range test.p2shMagics {
			valid := IsScriptHashAddrID(magTest.magic)
			if valid != magTest.valid {
				t.Errorf("%s: P2SH magic %d valid mismatch: got %v expected %v",
					test.name, i, valid, magTest.valid)
			}
		}
		for i, magTest := range test.hdMagics {
			pubKey, err := HDPrivateKeyToPublicKeyID(magTest.priv[:])
			if !reflect.DeepEqual(err, magTest.err) {
				t.Errorf("%s: HD magic %d mismatched error: got %v expected %v ",
					test.name, i, err, magTest.err)
				continue
			}
			if magTest.err == nil && !bytes.Equal(pubKey, magTest.want[:]) {
				t.Errorf("%s: HD magic %d private and public mismatch: got %v expected %v ",
					test.name, i, pubKey, magTest.want[:])
			}
		}
	}
}
