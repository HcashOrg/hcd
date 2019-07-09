// Copyright (c) 2013, 2014 The btcsuite developers
// Copyright (c) 2015 The Decred developers
// Copyright (c) 2018-2020 The Hcd developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hcutil_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	. "github.com/HcashOrg/hcd/hcutil"
)

//TestEncodeDecodeWIF tests encode and decode wallet import format
func TestEncodeDecodeWIF(t *testing.T) {

	priv := "ab9715425cb447f24dcd64df7099fea0dd41986fa917c569a255858dcbf26147"
	buf,_ := hex.DecodeString(priv)
	privKey, _ := chainec.Secp256k1.PrivKeyFromBytes(buf)
	algType := chainec.ECTypeSecp256k1

//	netIDTest := [2]byte{25,171}
	wifTest, _ := NewWIF(privKey, &chaincfg.MainNetParams, algType)

	fmt.Println(wifTest.String())
	return


	suites := []int{
		chainec.ECTypeSecp256k1,
		chainec.ECTypeEdwards,
		chainec.ECTypeSecSchnorr,
	}
	for _, suite := range suites {
		var priv1, priv2 chainec.PrivateKey
		switch suite {
		case chainec.ECTypeSecp256k1:
			priv1, _ = chainec.Secp256k1.PrivKeyFromBytes([]byte{
				0x0c, 0x28, 0xfc, 0xa3, 0x86, 0xc7, 0xa2, 0x27,
				0x60, 0x0b, 0x2f, 0xe5, 0x0b, 0x7c, 0xae, 0x11,
				0xec, 0x86, 0xd3, 0xbf, 0x1f, 0xbe, 0x47, 0x1b,
				0xe8, 0x98, 0x27, 0xe1, 0x9d, 0x72, 0xaa, 0x1d})

			priv2, _ = chainec.Secp256k1.PrivKeyFromBytes([]byte{
				0xdd, 0xa3, 0x5a, 0x14, 0x88, 0xfb, 0x97, 0xb6,
				0xeb, 0x3f, 0xe6, 0xe9, 0xef, 0x2a, 0x25, 0x81,
				0x4e, 0x39, 0x6f, 0xb5, 0xdc, 0x29, 0x5f, 0xe9,
				0x94, 0xb9, 0x67, 0x89, 0xb2, 0x1a, 0x03, 0x98})
		case chainec.ECTypeEdwards:
			priv1, _ = chainec.Edwards.PrivKeyFromScalar([]byte{
				0x0c, 0x28, 0xfc, 0xa3, 0x86, 0xc7, 0xa2, 0x27,
				0x60, 0x0b, 0x2f, 0xe5, 0x0b, 0x7c, 0xae, 0x11,
				0xec, 0x86, 0xd3, 0xbf, 0x1f, 0xbe, 0x47, 0x1b,
				0xe8, 0x98, 0x27, 0xe1, 0x9d, 0x72, 0xaa, 0x1d})

			priv2, _ = chainec.Edwards.PrivKeyFromScalar([]byte{
				0x0c, 0xa3, 0x5a, 0x14, 0x88, 0xfb, 0x97, 0xb6,
				0xeb, 0x3f, 0xe6, 0xe9, 0xef, 0x2a, 0x25, 0x81,
				0x4e, 0x39, 0x6f, 0xb5, 0xdc, 0x29, 0x5f, 0xe9,
				0x94, 0xb9, 0x67, 0x89, 0xb2, 0x1a, 0x03, 0x98})
		case chainec.ECTypeSecSchnorr:
			priv1, _ = chainec.SecSchnorr.PrivKeyFromBytes([]byte{
				0x0c, 0x28, 0xfc, 0xa3, 0x86, 0xc7, 0xa2, 0x27,
				0x60, 0x0b, 0x2f, 0xe5, 0x0b, 0x7c, 0xae, 0x11,
				0xec, 0x86, 0xd3, 0xbf, 0x1f, 0xbe, 0x47, 0x1b,
				0xe8, 0x98, 0x27, 0xe1, 0x9d, 0x72, 0xaa, 0x1d})

			priv2, _ = chainec.SecSchnorr.PrivKeyFromBytes([]byte{
				0xdd, 0xa3, 0x5a, 0x14, 0x88, 0xfb, 0x97, 0xb6,
				0xeb, 0x3f, 0xe6, 0xe9, 0xef, 0x2a, 0x25, 0x81,
				0x4e, 0x39, 0x6f, 0xb5, 0xdc, 0x29, 0x5f, 0xe9,
				0x94, 0xb9, 0x67, 0x89, 0xb2, 0x1a, 0x03, 0x98})
		}

		wif1, err := NewWIF(priv1, &chaincfg.MainNetParams, suite)
		if err != nil {
			t.Fatal(err)
		}
		wif2, err := NewWIF(priv2, &chaincfg.TestNet2Params, suite)
		if err != nil {
			t.Fatal(err)
		}
		wif3, err := NewWIF(priv2, &chaincfg.SimNetParams, suite)
		if err != nil {
			t.Fatal(err)
		}

		var tests []struct {
			wif     *WIF
			encoded string
		}

		switch suite {
		case chainec.ECTypeSecp256k1:
			tests = []struct {
				wif     *WIF
				encoded string
			}{
				{
					wif1,
					"Hm37RK5N6Pv1d1GZ8QFKuNjpR6vtLpPGTaoa1cf9yNWruUXwm3ySf",
				},
				{
					wif2,
					"PtWVDUidYaiiNT5e2Sfb1Ah4evbaSopZJkkpFBuzkJYcYteugvdFg",
				},
				{
					wif3,
					"PsURoUb7FMeJQdTYea8pkbUQFBZAsxtfDcfTLGja5sCLZvLZWRtjK",
				},
			}
		case chainec.ECTypeEdwards:
			tests = []struct {
				wif     *WIF
				encoded string
			}{
				{
					wif1,
					"Hm39N4FbYRmSJEQNi79uXNh91sqzxY2yJTbD5iuJfGWWus2dSyYDz",
				},
				{
					wif2,
					"PtWVaBGeCfbFQfgqFew8YvdrSH5TH439K7rvpo3aWnSfDvyK8ijbK",
				},
				{
					wif3,
					"PsUSAB97uSWqSr4jsnQNJMRC2Y33iD7FDymZuss9rM6PExexSPyTQ",
				},
			}
		case chainec.ECTypeSecSchnorr:
			tests = []struct {
				wif     *WIF
				encoded string
			}{
				{
					wif1,
					"Hm3BJoRpzTcryTYCHp4V9NeTcem7aFgg9LNr9q9TMAWAvFXF8pu2e",
				},
				{
					wif2,
					"PtWZ6y56SeRZiuMHBrUkFAbhrURogF7xzWL6PQQJ86XvZfeE3jf1a",
				},
				{
					wif3,
					"PsUVgxwa9RM9m5jBoywyzbP3SjPQ7QC4uNEjUVDsTfBeahKkmETvQ",
				},
			}
		}

		for _, test := range tests {
			// Test that encoding the WIF structure matches the expected string.
			s := test.wif.String()
			if s != test.encoded {
				t.Errorf("TestEncodeDecodePrivateKey failed: want '%s', got '%s'",
					test.encoded, s)
				continue
			}

			// Test that decoding the expected string results in the original WIF
			// structure.
			w, err := DecodeWIF(test.encoded)
			if err != nil {
				t.Error(err)
				continue
			}
			if got := w.String(); got != test.encoded {
				t.Errorf("NewWIF failed: want '%v', got '%v'", test.wif, got)
			}

			w.SerializePubKey()
		}
	}
}
