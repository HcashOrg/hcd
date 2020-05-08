// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript_test

import (
	"crypto/rand"
	"errors"
	"fmt"
	mrand "math/rand"
	"testing"

	"github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/txscript"
	"github.com/HcashOrg/hcd/wire"
	"github.com/HcashOrg/hcd/hcutil"
)

// testingParams defines the chain params to use throughout these tests so it
// can more easily be changed if desired.
var testingParams = &chaincfg.TestNet2Params

const testValueIn = 12345

type addressToKey struct {
	key        *chainec.PrivateKey
	compressed bool
}

func mkGetKey(keys map[string]addressToKey) txscript.KeyDB {
	if keys == nil {
		return txscript.KeyClosure(func(addr hcutil.Address) (chainec.PrivateKey,
			bool, error) {
			return nil, false, errors.New("nope 1")
		})
	}
	return txscript.KeyClosure(func(addr hcutil.Address) (chainec.PrivateKey,
		bool, error) {
		a2k, ok := keys[addr.EncodeAddress()]
		if !ok {
			return nil, false, errors.New("nope 2")
		}
		return *a2k.key, a2k.compressed, nil
	})
}

func mkGetKeyPub(keys map[string]addressToKey) txscript.KeyDB {
	if keys == nil {
		return txscript.KeyClosure(func(addr hcutil.Address) (chainec.PrivateKey,
			bool, error) {
			return nil, false, errors.New("nope 1")
		})
	}
	return txscript.KeyClosure(func(addr hcutil.Address) (chainec.PrivateKey,
		bool, error) {
		a2k, ok := keys[addr.String()]
		if !ok {
			return nil, false, errors.New("nope 2")
		}
		return *a2k.key, a2k.compressed, nil
	})
}

func mkGetScript(scripts map[string][]byte) txscript.ScriptDB {
	if scripts == nil {
		return txscript.ScriptClosure(func(addr hcutil.Address) (
			[]byte, error) {
			return nil, errors.New("nope 3")
		})
	}
	return txscript.ScriptClosure(func(addr hcutil.Address) ([]byte,
		error) {
		script, ok := scripts[addr.EncodeAddress()]
		if !ok {
			return nil, errors.New("nope 4")
		}
		return script, nil
	})
}

func checkScripts(msg string, tx *wire.MsgTx, idx int, sigScript, pkScript []byte) error {
	tx.TxIn[idx].SignatureScript = sigScript
	vm, err := txscript.NewEngine(pkScript, tx, idx,
		txscript.ScriptBip16|txscript.ScriptVerifyDERSignatures, 0, nil)
	if err != nil {
		return fmt.Errorf("failed to make script engine for %s: %v",
			msg, err)
	}

	err = vm.Execute()
	if err != nil {
		return fmt.Errorf("invalid script signature for %s: %v", msg,
			err)
	}

	return nil
}

func signAndCheck(msg string, tx *wire.MsgTx, idx int, pkScript []byte,
	hashType txscript.SigHashType, kdb txscript.KeyDB, sdb txscript.ScriptDB,
	previousScript []byte, suite int) error {

	sigScript, err := txscript.SignTxOutput(testingParams, tx, idx,
		pkScript, hashType, kdb, sdb, nil, suite)
	if err != nil {
		return fmt.Errorf("failed to sign output %s: %v", msg, err)
	}

	return checkScripts(msg, tx, idx, sigScript, pkScript)
}

func signBadAndCheck(msg string, tx *wire.MsgTx, idx int, pkScript []byte,
	hashType txscript.SigHashType, kdb txscript.KeyDB, sdb txscript.ScriptDB,
	previousScript []byte, suite int) error {
	// Setup a PRNG.
	randScriptHash := chainhash.HashB(pkScript)
	tRand := mrand.New(mrand.NewSource(int64(randScriptHash[0])))

	// Test SigHashAllValue by corrupting the transaction's ValueIn so that
	// the signature becomes invalid.
	if hashType == txscript.SigHashAllValue {
		tx.TxIn[0].ValueIn = 1
	}

	sigScript, err := txscript.SignTxOutput(testingParams, tx,
		idx, pkScript, hashType, kdb, sdb, nil, suite)
	if err != nil {
		return fmt.Errorf("failed to sign output %s: %v", msg, err)
	}

	// Be sure to reset the value in when we're done creating the
	// corrupted signature for that flag.
	tx.TxIn[0].ValueIn = testValueIn

	// Corrupt a random bit in the signature.
	if hashType != txscript.SigHashAllValue {
		pos := tRand.Intn(len(sigScript) - 1)
		bitPos := tRand.Intn(7)
		sigScript[pos] ^= 1 << uint8(bitPos)
	}

	return checkScripts(msg, tx, idx, sigScript, pkScript)
}

func TestSignTxOutput(t *testing.T) {
	t.Parallel()

	// make key
	// make script based on key.
	// sign with magic pixie dust.
	hashTypes := []txscript.SigHashType{
		txscript.SigHashOld, // no longer used but should act like all
		txscript.SigHashAll,
		txscript.SigHashNone,
		txscript.SigHashSingle,
		txscript.SigHashAllValue,
		txscript.SigHashAll | txscript.SigHashAnyOneCanPay,
		txscript.SigHashNone | txscript.SigHashAnyOneCanPay,
		txscript.SigHashSingle | txscript.SigHashAnyOneCanPay,
		txscript.SigHashAllValue | txscript.SigHashAnyOneCanPay,
	}
	signatureSuites := []int{
		secp,
		edwards,
		secSchnorr,
	}
	tx := &wire.MsgTx{
		SerType: wire.TxSerializeFull,
		Version: 1,
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 0,
					Tree:  0,
				},
				Sequence:    4294967295,
				ValueIn:     testValueIn,
				BlockHeight: 78901,
				BlockIndex:  23456,
			},
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 1,
					Tree:  0,
				},
				Sequence:    4294967295,
				ValueIn:     testValueIn,
				BlockHeight: 78901,
				BlockIndex:  23456,
			},
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 2,
					Tree:  0,
				},
				Sequence:    4294967295,
				ValueIn:     testValueIn,
				BlockHeight: 78901,
				BlockIndex:  23456,
			},
		},
		TxOut: []*wire.TxOut{
			{
				Version: wire.DefaultPkScriptVersion,
				Value:   1,
			},
			{
				Version: wire.DefaultPkScriptVersion,
				Value:   2,
			},
			{
				Version: wire.DefaultPkScriptVersion,
				Value:   3,
			},
		},
		LockTime: 0,
		Expiry:   0,
	}

	// Pay to Pubkey Hash (uncompressed)
	secp256k1 := chainec.Secp256k1
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(
						rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
				}

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				address, err := hcutil.NewAddressPubKeyHash(
					hcutil.Hash160(pkBytes), testingParams,
					suite)

				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				if err := signAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), nil, suite); err != nil {
					t.Error(err)
					break
				}

				if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), nil, suite); err == nil {
					t.Errorf("corrupted signature validated %s: %v",
						msg, err)
					break
				}
			}
		}
	}

	// Pay to Pubkey Hash (uncompressed) (merging with correct)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
				}

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				address, err := hcutil.NewAddressPubKeyHash(
					hcutil.Hash160(pkBytes), testingParams,
					suite)
				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				sigScript, err := txscript.SignTxOutput(
					testingParams, tx, i, pkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s: %v", msg,
						err)
					break
				}

				// by the above loop, this should be valid, now sign
				// again and merge.
				sigScript, err = txscript.SignTxOutput(
					testingParams, tx, i, pkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), sigScript, suite)
				if err != nil {
					t.Errorf("failed to sign output %s a "+
						"second time: %v", msg, err)
					break
				}

				err = checkScripts(msg, tx, i, sigScript, pkScript)
				if err != nil {
					t.Errorf("twice signed script invalid for "+
						"%s: %v", msg, err)
					break
				}
			}
		}
	}

	// Pay to Pubkey Hash (compressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				}

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				address, err := hcutil.NewAddressPubKeyHash(
					hcutil.Hash160(pkBytes), testingParams,
					suite)
				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				if err := signAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(nil), nil, suite); err != nil {
					t.Error(err)
					break
				}

				if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(nil), nil, suite); err == nil {
					t.Errorf("corrupted signature validated %s: %v",
						msg, err)
					break
				}
			}
		}
	}

	// Pay to Pubkey Hash (compressed) with duplicate merge
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				}

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				address, err := hcutil.NewAddressPubKeyHash(
					hcutil.Hash160(pkBytes), testingParams,
					suite)
				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				sigScript, err := txscript.SignTxOutput(
					testingParams, tx, i, pkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(nil), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s: %v", msg,
						err)
					break
				}

				// by the above loop, this should be valid, now sign
				// again and merge.
				sigScript, err = txscript.SignTxOutput(
					testingParams, tx, i, pkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(nil), sigScript, suite)
				if err != nil {
					t.Errorf("failed to sign output %s a "+
						"second time: %v", msg, err)
					break
				}

				err = checkScripts(msg, tx, i, sigScript, pkScript)
				if err != nil {
					t.Errorf("twice signed script invalid for "+
						"%s: %v", msg, err)
					break
				}
			}
		}
	}

	// Pay to Pubkey Hash for a ticket(SStx) (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)

			keyDB, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key, pk := secp256k1.PrivKeyFromBytes(keyDB)
			pkBytes := pk.SerializeCompressed()

			address, err := hcutil.NewAddressPubKeyHash(
				hcutil.Hash160(pkBytes), testingParams, secp)
			if err != nil {
				t.Errorf("failed to make address for %s: %v",
					msg, err)
				break
			}

			pkScript, err := txscript.PayToSStx(address)
			if err != nil {
				t.Errorf("failed to make pkscript "+
					"for %s: %v", msg, err)
			}

			if err := signAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address.EncodeAddress(): {&key, true},
				}), mkGetScript(nil), nil, secp); err != nil {
				t.Error(err)
				break
			}

			if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address.EncodeAddress(): {&key, true},
				}), mkGetScript(nil), nil, secp); err == nil {
				t.Errorf("corrupted signature validated %s: %v",
					msg, err)
				break
			}
		}
	}

	// Pay to Pubkey Hash for a ticket change (SStx change) (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)

			keyDB, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key, pk := secp256k1.PrivKeyFromBytes(keyDB)
			pkBytes := pk.SerializeCompressed()

			address, err := hcutil.NewAddressPubKeyHash(
				hcutil.Hash160(pkBytes), testingParams, secp)
			if err != nil {
				t.Errorf("failed to make address for %s: %v",
					msg, err)
				break
			}

			pkScript, err := txscript.PayToSStxChange(address)
			if err != nil {
				t.Errorf("failed to make pkscript "+
					"for %s: %v", msg, err)
			}

			if err := signAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address.EncodeAddress(): {&key, true},
				}), mkGetScript(nil), nil, secp); err != nil {
				t.Error(err)
				break
			}

			if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address.EncodeAddress(): {&key, true},
				}), mkGetScript(nil), nil, secp); err == nil {
				t.Errorf("corrupted signature validated %s: %v",
					msg, err)
				break
			}
		}
	}

	// Pay to Pubkey Hash for a ticket spending (SSGen) (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)

			keyDB, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key, pk := secp256k1.PrivKeyFromBytes(keyDB)
			pkBytes := pk.SerializeCompressed()

			address, err := hcutil.NewAddressPubKeyHash(
				hcutil.Hash160(pkBytes), testingParams, secp)
			if err != nil {
				t.Errorf("failed to make address for %s: %v",
					msg, err)
				break
			}

			pkScript, err := txscript.PayToSSGen(address)
			if err != nil {
				t.Errorf("failed to make pkscript "+
					"for %s: %v", msg, err)
			}

			if err := signAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address.EncodeAddress(): {&key, true},
				}), mkGetScript(nil), nil, secp); err != nil {
				t.Error(err)
				break
			}

			if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address.EncodeAddress(): {&key, true},
				}), mkGetScript(nil), nil, secp); err == nil {
				t.Errorf("corrupted signature validated %s: %v",
					msg, err)
				break
			}
		}
	}

	// Pay to Pubkey Hash for a ticket revocation (SSRtx) (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)

			keyDB, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key, pk := secp256k1.PrivKeyFromBytes(keyDB)
			pkBytes := pk.SerializeCompressed()

			address, err := hcutil.NewAddressPubKeyHash(
				hcutil.Hash160(pkBytes), testingParams, secp)
			if err != nil {
				t.Errorf("failed to make address for %s: %v",
					msg, err)
				break
			}

			pkScript, err := txscript.PayToSSRtx(address)
			if err != nil {
				t.Errorf("failed to make pkscript "+
					"for %s: %v", msg, err)
			}

			if err := signAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address.EncodeAddress(): {&key, true},
				}), mkGetScript(nil), nil, secp); err != nil {
				t.Error(err)
				break
			}

			if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address.EncodeAddress(): {&key, true},
				}), mkGetScript(nil), nil, secp); err == nil {
				t.Errorf("corrupted signature validated %s: %v",
					msg, err)
				break
			}
		}
	}

	// Pay to PubKey (uncompressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				keyDB, _, _, err := secp256k1.GenerateKey(rand.Reader)
				key, pk := secp256k1.PrivKeyFromBytes(keyDB)
				// For address generation, consensus rules require using
				// a compressed public key. Look up ExtractPkScriptAddrs
				// for more details
				address, err := hcutil.NewAddressSecpPubKeyCompressed(pk,
					testingParams)
				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				if err := signAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), nil, suite); err != nil {
					t.Error(err)
					break
				}

				if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), nil, suite); err == nil {
					t.Errorf("corrupted signature validated %s: %v",
						msg, err)
					break
				}
			}
		}
	}

	// Pay to PubKey (uncompressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey
				var address hcutil.Address
				var err error

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					// For address generation, consensus rules require using
					// a compressed public key. Look up ExtractPkScriptAddrs
					// for more details
					address, err = hcutil.NewAddressSecpPubKeyCompressed(pk,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
					address, err = hcutil.NewAddressEdwardsPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
					address, err = hcutil.NewAddressSecSchnorrPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)

				}

				sigScript, err := txscript.SignTxOutput(
					testingParams, tx, i, pkScript,
					hashType, mkGetKeyPub(map[string]addressToKey{
						address.String(): {&key, false},
					}), mkGetScript(nil), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s: %v", msg,
						err)
					break
				}

				// by the above loop, this should be valid, now sign
				// again and merge.
				sigScript, err = txscript.SignTxOutput(
					testingParams, tx, i, pkScript,
					hashType, mkGetKeyPub(map[string]addressToKey{
						address.String(): {&key, false},
					}), mkGetScript(nil), sigScript, suite)
				if err != nil {
					t.Errorf("failed to sign output %s a "+
						"second time: %v", msg, err)
					break
				}

				err = checkScripts(msg, tx, i, sigScript, pkScript)
				if err != nil {
					t.Errorf("twice signed script invalid for "+
						"%s: %v", msg, err)
					break
				}
			}
		}
	}

	// Pay to PubKey (compressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey
				var address hcutil.Address
				var err error

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					// For address generation, consensus rules require using
					// a compressed public key. Look up ExtractPkScriptAddrs
					// for more details
					address, err = hcutil.NewAddressSecpPubKeyCompressed(pk,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
					address, err = hcutil.NewAddressEdwardsPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
					address, err = hcutil.NewAddressSecSchnorrPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				if err := signAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKeyPub(map[string]addressToKey{
						address.String(): {&key, true},
					}), mkGetScript(nil), nil, suite); err != nil {
					t.Error(err)
					break
				}

				if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKeyPub(map[string]addressToKey{
						address.String(): {&key, true},
					}), mkGetScript(nil), nil, suite); err == nil {
					t.Errorf("corrupted signature validated %s: %v",
						msg, err)
					break
				}
			}
		}
	}

	// Pay to PubKey (compressed) with duplicate merge
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey
				var address hcutil.Address
				var err error

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					address, err = hcutil.NewAddressSecpPubKeyCompressed(pk,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
					address, err = hcutil.NewAddressEdwardsPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
					address, err = hcutil.NewAddressSecSchnorrPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				sigScript, err := txscript.SignTxOutput(
					testingParams, tx, i, pkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(nil), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s: %v", msg,
						err)
					break
				}

				// by the above loop, this should be valid, now sign
				// again and merge.
				sigScript, err = txscript.SignTxOutput(
					testingParams, tx, i, pkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(nil), sigScript, suite)
				if err != nil {
					t.Errorf("failed to sign output %s a "+
						"second time: %v", msg, err)
					break
				}

				err = checkScripts(msg, tx, i, sigScript, pkScript)
				if err != nil {
					t.Errorf("twice signed script invalid for "+
						"%s: %v", msg, err)
					break
				}
			}
		}
	}

	// As before, but with p2sh now.
	// Pay to Pubkey Hash (uncompressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
				}

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				address, err := hcutil.NewAddressPubKeyHash(
					hcutil.Hash160(pkBytes), testingParams, suite)
				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
					break
				}

				scriptAddr, err := hcutil.NewAddressScriptHash(
					pkScript, testingParams)
				if err != nil {
					t.Errorf("failed to make p2sh addr for %s: %v",
						msg, err)
					break
				}

				scriptPkScript, err := txscript.PayToAddrScript(
					scriptAddr)
				if err != nil {
					t.Errorf("failed to make script pkscript for "+
						"%s: %v", msg, err)
					break
				}

				if err := signAndCheck(msg, tx, i, scriptPkScript,
					hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite); err != nil {
					t.Error(err)
					break
				}

				if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), nil, suite); err == nil {
					t.Errorf("corrupted signature validated %s: %v",
						msg, err)
					break
				}
			}
		}
	}

	// Pay to Pubkey Hash (uncompressed) with duplicate merge
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
				}

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				address, err := hcutil.NewAddressPubKeyHash(
					hcutil.Hash160(pkBytes), testingParams, suite)
				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
					break
				}

				scriptAddr, err := hcutil.NewAddressScriptHash(
					pkScript, testingParams)
				if err != nil {
					t.Errorf("failed to make p2sh addr for %s: %v",
						msg, err)
					break
				}

				scriptPkScript, err := txscript.PayToAddrScript(
					scriptAddr)
				if err != nil {
					t.Errorf("failed to make script pkscript for "+
						"%s: %v", msg, err)
					break
				}

				sigScript, err := txscript.SignTxOutput(
					testingParams, tx, i, scriptPkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s: %v", msg,
						err)
					break
				}

				// by the above loop, this should be valid, now sign
				// again and merge.
				sigScript, err = txscript.SignTxOutput(
					testingParams, tx, i, scriptPkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s a "+
						"second time: %v", msg, err)
					break
				}

				err = checkScripts(msg, tx, i, sigScript, scriptPkScript)
				if err != nil {
					t.Errorf("twice signed script invalid for "+
						"%s: %v", msg, err)
					break
				}
			}
		}
	}

	// Pay to Pubkey Hash (compressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				}

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				address, err := hcutil.NewAddressPubKeyHash(
					hcutil.Hash160(pkBytes), testingParams, suite)
				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				scriptAddr, err := hcutil.NewAddressScriptHash(
					pkScript, testingParams)
				if err != nil {
					t.Errorf("failed to make p2sh addr for %s: %v",
						msg, err)
					break
				}

				scriptPkScript, err := txscript.PayToAddrScript(
					scriptAddr)
				if err != nil {
					t.Errorf("failed to make script pkscript for "+
						"%s: %v", msg, err)
					break
				}

				if err := signAndCheck(msg, tx, i, scriptPkScript,
					hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite); err != nil {
					t.Error(err)
					break
				}

				if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(nil), nil, suite); err == nil {
					t.Errorf("corrupted signature validated %s: %v",
						msg, err)
					break
				}
			}
		}
	}

	// Pay to Pubkey Hash (compressed) with duplicate merge
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
				}

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				address, err := hcutil.NewAddressPubKeyHash(
					hcutil.Hash160(pkBytes), testingParams, suite)
				if err != nil {
					t.Errorf("failed to make address for %s: %v",
						msg, err)
					break
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				scriptAddr, err := hcutil.NewAddressScriptHash(
					pkScript, testingParams)
				if err != nil {
					t.Errorf("failed to make p2sh addr for %s: %v",
						msg, err)
					break
				}

				scriptPkScript, err := txscript.PayToAddrScript(
					scriptAddr)
				if err != nil {
					t.Errorf("failed to make script pkscript for "+
						"%s: %v", msg, err)
					break
				}

				sigScript, err := txscript.SignTxOutput(
					testingParams, tx, i, scriptPkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s: %v", msg,
						err)
					break
				}

				// by the above loop, this should be valid, now sign
				// again and merge.
				sigScript, err = txscript.SignTxOutput(
					testingParams, tx, i, scriptPkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s a "+
						"second time: %v", msg, err)
					break
				}

				err = checkScripts(msg, tx, i, sigScript, scriptPkScript)
				if err != nil {
					t.Errorf("twice signed script invalid for "+
						"%s: %v", msg, err)
					break
				}
			}
		}
	}

	// Pay to PubKey (uncompressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey
				var address hcutil.Address
				var err error

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					// For address generation, consensus rules require using
					// a compressed public key. Look up ExtractPkScriptAddrs
					// for more details
					address, err = hcutil.NewAddressSecpPubKeyCompressed(pk,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
					address, err = hcutil.NewAddressEdwardsPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
					address, err = hcutil.NewAddressSecSchnorrPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				scriptAddr, err := hcutil.NewAddressScriptHash(
					pkScript, testingParams)
				if err != nil {
					t.Errorf("failed to make p2sh addr for %s: %v",
						msg, err)
					break
				}

				scriptPkScript, err := txscript.PayToAddrScript(
					scriptAddr)
				if err != nil {
					t.Errorf("failed to make script pkscript for "+
						"%s: %v", msg, err)
					break
				}

				if err := signAndCheck(msg, tx, i, scriptPkScript,
					hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite); err != nil {
					t.Error(err)
					break
				}

				if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), nil, suite); err == nil {
					t.Errorf("corrupted signature validated %s: %v",
						msg, err)
					break
				}
			}
		}
	}

	// Pay to PubKey (uncompressed) with duplicate merge
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey
				var address hcutil.Address
				var err error

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					// For address generation, consensus rules require using
					// a compressed public key. Look up ExtractPkScriptAddrs
					// for more details
					address, err = hcutil.NewAddressSecpPubKeyCompressed(pk,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeUncompressed()
					address, err = hcutil.NewAddressEdwardsPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
					address, err = hcutil.NewAddressSecSchnorrPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				scriptAddr, err := hcutil.NewAddressScriptHash(
					pkScript, testingParams)
				if err != nil {
					t.Errorf("failed to make p2sh addr for %s: %v",
						msg, err)
					break
				}

				scriptPkScript, err := txscript.PayToAddrScript(
					scriptAddr)
				if err != nil {
					t.Errorf("failed to make script pkscript for "+
						"%s: %v", msg, err)
					break
				}

				sigScript, err := txscript.SignTxOutput(
					testingParams, tx, i, scriptPkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s: %v", msg,
						err)
					break
				}

				// by the above loop, this should be valid, now sign
				// again and merge.
				sigScript, err = txscript.SignTxOutput(
					testingParams, tx, i, scriptPkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s a "+
						"second time: %v", msg, err)
					break
				}

				err = checkScripts(msg, tx, i, sigScript, scriptPkScript)
				if err != nil {
					t.Errorf("twice signed script invalid for "+
						"%s: %v", msg, err)
					break
				}
			}
		}
	}

	// Pay to PubKey (compressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey
				var address hcutil.Address
				var err error

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					address, err = hcutil.NewAddressSecpPubKeyCompressed(pk,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
					address, err = hcutil.NewAddressEdwardsPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
					address, err = hcutil.NewAddressSecSchnorrPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				scriptAddr, err := hcutil.NewAddressScriptHash(
					pkScript, testingParams)
				if err != nil {
					t.Errorf("failed to make p2sh addr for %s: %v",
						msg, err)
					break
				}

				scriptPkScript, err := txscript.PayToAddrScript(
					scriptAddr)
				if err != nil {
					t.Errorf("failed to make script pkscript for "+
						"%s: %v", msg, err)
					break
				}

				if err := signAndCheck(msg, tx, i, scriptPkScript,
					hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite); err != nil {
					t.Error(err)
					break
				}

				if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
					mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, false},
					}), mkGetScript(nil), nil, suite); err == nil {
					t.Errorf("corrupted signature validated %s: %v",
						msg, err)
					break
				}
			}
		}
	}

	// Pay to PubKey (compressed)
	for _, hashType := range hashTypes {
		for _, suite := range signatureSuites {
			for i := range tx.TxIn {
				var keyDB, pkBytes []byte
				var key chainec.PrivateKey
				var pk chainec.PublicKey
				var address hcutil.Address
				var err error

				msg := fmt.Sprintf("%d:%d:%d", hashType, i, suite)

				switch suite {
				case secp:
					keyDB, _, _, _ = secp256k1.GenerateKey(rand.Reader)
					key, pk = secp256k1.PrivKeyFromBytes(keyDB)
					address, err = hcutil.NewAddressSecpPubKeyCompressed(pk,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case edwards:
					keyDB, _, _, _ = chainec.Edwards.GenerateKey(rand.Reader)
					key, pk = chainec.Edwards.PrivKeyFromBytes(keyDB)
					pkBytes = pk.SerializeCompressed()
					address, err = hcutil.NewAddressEdwardsPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}

				case secSchnorr:
					keyDB, _, _, _ = chainec.SecSchnorr.GenerateKey(rand.Reader)
					key, pk = chainec.SecSchnorr.PrivKeyFromBytes(keyDB)
					pkBytes = pk.Serialize()
					address, err = hcutil.NewAddressSecSchnorrPubKey(pkBytes,
						testingParams)
					if err != nil {
						t.Errorf("failed to make address for %s: %v",
							msg, err)
						break
					}
				}

				pkScript, err := txscript.PayToAddrScript(address)
				if err != nil {
					t.Errorf("failed to make pkscript "+
						"for %s: %v", msg, err)
				}

				scriptAddr, err := hcutil.NewAddressScriptHash(
					pkScript, testingParams)
				if err != nil {
					t.Errorf("failed to make p2sh addr for %s: %v",
						msg, err)
					break
				}

				scriptPkScript, err := txscript.PayToAddrScript(
					scriptAddr)
				if err != nil {
					t.Errorf("failed to make script pkscript for "+
						"%s: %v", msg, err)
					break
				}

				sigScript, err := txscript.SignTxOutput(
					testingParams, tx, i, scriptPkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s: %v", msg,
						err)
					break
				}

				// by the above loop, this should be valid, now sign
				// again and merge.
				sigScript, err = txscript.SignTxOutput(
					testingParams, tx, i, scriptPkScript,
					hashType, mkGetKey(map[string]addressToKey{
						address.EncodeAddress(): {&key, true},
					}), mkGetScript(map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					}), nil, suite)
				if err != nil {
					t.Errorf("failed to sign output %s a "+
						"second time: %v", msg, err)
					break
				}

				err = checkScripts(msg, tx, i, sigScript, scriptPkScript)
				if err != nil {
					t.Errorf("twice signed script invalid for "+
						"%s: %v", msg, err)
					break
				}
			}
		}
	}

	// Basic Multisig
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)

			keyDB1, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key1, pk1 := secp256k1.PrivKeyFromBytes(keyDB1)

			address1, err := hcutil.NewAddressSecpPubKeyCompressed(pk1,
				testingParams)
			if err != nil {
				t.Errorf("failed to make address for %s: %v",
					msg, err)
				break
			}

			keyDB2, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key2, pk2 := secp256k1.PrivKeyFromBytes(keyDB2)

			address2, err := hcutil.NewAddressSecpPubKeyCompressed(pk2,
				testingParams)
			if err != nil {
				t.Errorf("failed to make address 2 for %s: %v",
					msg, err)
				break
			}

			pkScript, err := txscript.MultiSigScript(
				[]*hcutil.AddressSecpPubKey{address1, address2},
				2)
			if err != nil {
				t.Errorf("failed to make pkscript "+
					"for %s: %v", msg, err)
			}

			scriptAddr, err := hcutil.NewAddressScriptHash(
				pkScript, testingParams)
			if err != nil {
				t.Errorf("failed to make p2sh addr for %s: %v",
					msg, err)
				break
			}

			scriptPkScript, err := txscript.PayToAddrScript(
				scriptAddr)
			if err != nil {
				t.Errorf("failed to make script pkscript for "+
					"%s: %v", msg, err)
				break
			}

			if err := signAndCheck(msg, tx, i, scriptPkScript,
				hashType,
				mkGetKey(map[string]addressToKey{
					address1.EncodeAddress(): {&key1, true},
					address2.EncodeAddress(): {&key2, true},
				}), mkGetScript(map[string][]byte{
					scriptAddr.EncodeAddress(): pkScript,
				}), nil, secp); err != nil {
				t.Error(err)
				break
			}

			if err := signBadAndCheck(msg, tx, i, pkScript, hashType,
				mkGetKey(map[string]addressToKey{
					address1.EncodeAddress(): {&key1, true},
					address2.EncodeAddress(): {&key2, true},
				}), mkGetScript(nil), nil, secp); err == nil {
				t.Errorf("corrupted signature validated %s: %v",
					msg, err)
				break
			}
		}
	}

	// Two part multisig, sign with one key then the other.
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)

			keyDB1, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key1, pk1 := secp256k1.PrivKeyFromBytes(keyDB1)

			address1, err := hcutil.NewAddressSecpPubKeyCompressed(pk1,
				testingParams)
			if err != nil {
				t.Errorf("failed to make address for %s: %v",
					msg, err)
				break
			}

			keyDB2, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key2, pk2 := secp256k1.PrivKeyFromBytes(keyDB2)

			address2, err := hcutil.NewAddressSecpPubKeyCompressed(pk2,
				testingParams)
			if err != nil {
				t.Errorf("failed to make address 2 for %s: %v",
					msg, err)
				break
			}

			pkScript, err := txscript.MultiSigScript(
				[]*hcutil.AddressSecpPubKey{address1, address2},
				2)
			if err != nil {
				t.Errorf("failed to make pkscript "+
					"for %s: %v", msg, err)
			}

			scriptAddr, err := hcutil.NewAddressScriptHash(
				pkScript, testingParams)
			if err != nil {
				t.Errorf("failed to make p2sh addr for %s: %v",
					msg, err)
				break
			}

			scriptPkScript, err := txscript.PayToAddrScript(
				scriptAddr)
			if err != nil {
				t.Errorf("failed to make script pkscript for "+
					"%s: %v", msg, err)
				break
			}

			sigScript, err := txscript.SignTxOutput(
				testingParams, tx, i, scriptPkScript,
				hashType, mkGetKey(map[string]addressToKey{
					address1.EncodeAddress(): {&key1, true},
				}), mkGetScript(map[string][]byte{
					scriptAddr.EncodeAddress(): pkScript,
				}), nil, secp)
			if err != nil {
				t.Errorf("failed to sign output %s: %v", msg,
					err)
				break
			}

			// Only 1 out of 2 signed, this *should* fail.
			if checkScripts(msg, tx, i, sigScript,
				scriptPkScript) == nil {
				t.Errorf("part signed script valid for %s", msg)
				break
			}

			// Sign with the other key and merge
			sigScript, err = txscript.SignTxOutput(
				testingParams, tx, i, scriptPkScript,
				hashType, mkGetKey(map[string]addressToKey{
					address2.EncodeAddress(): {&key2, true},
				}), mkGetScript(map[string][]byte{
					scriptAddr.EncodeAddress(): pkScript,
				}), sigScript, secp)
			if err != nil {
				t.Errorf("failed to sign output %s: %v", msg, err)
				break
			}

			err = checkScripts(msg, tx, i, sigScript,
				scriptPkScript)
			if err != nil {
				t.Errorf("fully signed script invalid for "+
					"%s: %v", msg, err)
				break
			}
		}
	}

	// Two part multisig, sign with one key then both, check key dedup
	// correctly.
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)

			keyDB1, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key1, pk1 := secp256k1.PrivKeyFromBytes(keyDB1)

			address1, err := hcutil.NewAddressSecpPubKeyCompressed(pk1,
				testingParams)
			if err != nil {
				t.Errorf("failed to make address for %s: %v",
					msg, err)
				break
			}

			keyDB2, _, _, err := secp256k1.GenerateKey(rand.Reader)
			key2, pk2 := secp256k1.PrivKeyFromBytes(keyDB2)
			address2, err := hcutil.NewAddressSecpPubKeyCompressed(pk2,
				testingParams)
			if err != nil {
				t.Errorf("failed to make address 2 for %s: %v",
					msg, err)
				break
			}

			pkScript, err := txscript.MultiSigScript(
				[]*hcutil.AddressSecpPubKey{address1, address2},
				2)
			if err != nil {
				t.Errorf("failed to make pkscript "+
					"for %s: %v", msg, err)
			}

			scriptAddr, err := hcutil.NewAddressScriptHash(
				pkScript, testingParams)
			if err != nil {
				t.Errorf("failed to make p2sh addr for %s: %v",
					msg, err)
				break
			}

			scriptPkScript, err := txscript.PayToAddrScript(
				scriptAddr)
			if err != nil {
				t.Errorf("failed to make script pkscript for "+
					"%s: %v", msg, err)
				break
			}

			sigScript, err := txscript.SignTxOutput(
				testingParams, tx, i, scriptPkScript,
				hashType, mkGetKey(map[string]addressToKey{
					address1.EncodeAddress(): {&key1, true},
				}), mkGetScript(map[string][]byte{
					scriptAddr.EncodeAddress(): pkScript,
				}), nil, secp)
			if err != nil {
				t.Errorf("failed to sign output %s: %v", msg,
					err)
				break
			}

			// Only 1 out of 2 signed, this *should* fail.
			if checkScripts(msg, tx, i, sigScript,
				scriptPkScript) == nil {
				t.Errorf("part signed script valid for %s", msg)
				break
			}

			// Sign with the other key and merge
			sigScript, err = txscript.SignTxOutput(
				testingParams, tx, i, scriptPkScript,
				hashType, mkGetKey(map[string]addressToKey{
					address1.EncodeAddress(): {&key1, true},
					address2.EncodeAddress(): {&key2, true},
				}), mkGetScript(map[string][]byte{
					scriptAddr.EncodeAddress(): pkScript,
				}), sigScript, secp)
			if err != nil {
				t.Errorf("failed to sign output %s: %v", msg, err)
				break
			}

			// Now we should pass.
			err = checkScripts(msg, tx, i, sigScript,
				scriptPkScript)
			if err != nil {
				t.Errorf("fully signed script invalid for "+
					"%s: %v", msg, err)
				break
			}
		}
	}
}

type tstInput struct {
	txout              *wire.TxOut
	sigscriptGenerates bool
	inputValidates     bool
	indexOutOfRange    bool
}

type tstSigScript struct {
	name               string
	inputs             []tstInput
	hashType           txscript.SigHashType
	compress           bool
	scriptAtWrongIndex bool
}

var coinbaseOutPoint = &wire.OutPoint{
	Index: (1 << 32) - 1,
}

// Pregenerated private key, with associated public key and pkScripts
// for the uncompressed and compressed hash160.
var (
	privKeyD = []byte{0x6b, 0x0f, 0xd8, 0xda, 0x54, 0x22, 0xd0, 0xb7,
		0xb4, 0xfc, 0x4e, 0x55, 0xd4, 0x88, 0x42, 0xb3, 0xa1, 0x65,
		0xac, 0x70, 0x7f, 0x3d, 0xa4, 0x39, 0x5e, 0xcb, 0x3b, 0xb0,
		0xd6, 0x0e, 0x06, 0x92}
	pubkeyX = []byte{0xb2, 0x52, 0xf0, 0x49, 0x85, 0x78, 0x03, 0x03, 0xc8,
		0x7d, 0xce, 0x51, 0x7f, 0xa8, 0x69, 0x0b, 0x91, 0x95, 0xf4,
		0xf3, 0x5c, 0x26, 0x73, 0x05, 0x05, 0xa2, 0xee, 0xbc, 0x09,
		0x38, 0x34, 0x3a}
	pubkeyY = []byte{0xb7, 0xc6, 0x7d, 0xb2, 0xe1, 0xff, 0xc8, 0x43, 0x1f,
		0x63, 0x32, 0x62, 0xaa, 0x60, 0xc6, 0x83, 0x30, 0xbd, 0x24,
		0x7e, 0xef, 0xdb, 0x6f, 0x2e, 0x8d, 0x56, 0xf0, 0x3c, 0x9f,
		0x6d, 0xb6, 0xf8}
	_, thisPubKey     = chainec.Secp256k1.PrivKeyFromBytes(privKeyD)
	thisAddressUnc, _ = hcutil.NewAddressPubKeyHash(
		hcutil.Hash160(thisPubKey.SerializeUncompressed()),
		testingParams, secp)
	uncompressedPkScript, _ = txscript.PayToAddrScript(thisAddressUnc)
	thisAddressCom, _       = hcutil.NewAddressPubKeyHash(
		hcutil.Hash160(thisPubKey.SerializeCompressed()),
		testingParams, secp)
	compressedPkScript, _ = txscript.PayToAddrScript(thisAddressCom)
	shortPkScript         = []byte{0x76, 0xa9, 0x14, 0xd1, 0x7c, 0xb5,
		0xeb, 0xa4, 0x02, 0xcb, 0x68, 0xe0, 0x69, 0x56, 0xbf, 0x32,
		0x53, 0x90, 0x0e, 0x0a, 0x88, 0xac}
	uncompressedAddrStr = "1L6fd93zGmtzkK6CsZFVVoCwzZV3MUtJ4F"
	compressedAddrStr   = "14apLppt9zTq6cNw8SDfiJhk9PhkZrQtYZ"
)

// Pretend output amounts.
const coinbaseVal = 2500000000
const fee = 5000000

var sigScriptTests = []tstSigScript{
	{
		name: "one input uncompressed",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAll,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "two inputs uncompressed",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
			{
				txout:              wire.NewTxOut(coinbaseVal+fee, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAll,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "one input compressed",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, compressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAll,
		compress:           true,
		scriptAtWrongIndex: false,
	},
	{
		name: "two inputs compressed",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, compressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
			{
				txout:              wire.NewTxOut(coinbaseVal+fee, compressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAll,
		compress:           true,
		scriptAtWrongIndex: false,
	},
	{
		name: "hashType SigHashNone",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashNone,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "hashType SigHashSingle",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashSingle,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "hashType SigHashAnyoneCanPay",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAnyOneCanPay,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "hashType non-standard",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           0x04,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "invalid compression",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     false,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAll,
		compress:           true,
		scriptAtWrongIndex: false,
	},
	{
		name: "short PkScript",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, shortPkScript),
				sigscriptGenerates: false,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAll,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "valid script at wrong index",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
			{
				txout:              wire.NewTxOut(coinbaseVal+fee, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAll,
		compress:           false,
		scriptAtWrongIndex: true,
	},
	{
		name: "index out of range",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
			{
				txout:              wire.NewTxOut(coinbaseVal+fee, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           txscript.SigHashAll,
		compress:           false,
		scriptAtWrongIndex: true,
	},
}

// Test the sigscript generation for valid and invalid inputs, all
// hashTypes, and with and without compression.  This test creates
// sigscripts to spend fake coinbase inputs, as sigscripts cannot be
// created for the MsgTxs in txTests, since they come from the blockchain
// and we don't have the private keys.
func TestSignatureScript(t *testing.T) {
	t.Parallel()
	privKey, _ := chainec.Secp256k1.PrivKeyFromBytes(privKeyD)

nexttest:
	for i := range sigScriptTests {
		tx := wire.NewMsgTx()

		output := wire.NewTxOut(500, []byte{txscript.OP_RETURN})
		tx.AddTxOut(output)

		for range sigScriptTests[i].inputs {
			txin := wire.NewTxIn(coinbaseOutPoint, nil)
			tx.AddTxIn(txin)
		}

		var script []byte
		var err error
		for j := range tx.TxIn {
			var idx int
			if sigScriptTests[i].inputs[j].indexOutOfRange {
				t.Errorf("at test %v", sigScriptTests[i].name)
				idx = len(sigScriptTests[i].inputs)
			} else {
				idx = j
			}
			script, err = txscript.SignatureScript(tx, idx,
				sigScriptTests[i].inputs[j].txout.PkScript,
				sigScriptTests[i].hashType, privKey,
				sigScriptTests[i].compress)

			if (err == nil) != sigScriptTests[i].inputs[j].sigscriptGenerates {
				if err == nil {
					t.Errorf("passed test '%v' incorrectly",
						sigScriptTests[i].name)
				} else {
					t.Errorf("failed test '%v': %v",
						sigScriptTests[i].name, err)
				}
				continue nexttest
			}
			if !sigScriptTests[i].inputs[j].sigscriptGenerates {
				// done with this test
				continue nexttest
			}

			tx.TxIn[j].SignatureScript = script
		}

		// If testing using a correct sigscript but for an incorrect
		// index, use last input script for first input.  Requires > 0
		// inputs for test.
		if sigScriptTests[i].scriptAtWrongIndex {
			tx.TxIn[0].SignatureScript = script
			sigScriptTests[i].inputs[0].inputValidates = false
		}

		// Validate tx input scripts
		scriptFlags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures
		for j := range tx.TxIn {
			vm, err := txscript.NewEngine(sigScriptTests[i].
				inputs[j].txout.PkScript, tx, j, scriptFlags, 0,
				nil)
			if err != nil {
				t.Errorf("cannot create script vm for test %v: %v",
					sigScriptTests[i].name, err)
				continue nexttest
			}
			err = vm.Execute()
			if (err == nil) != sigScriptTests[i].inputs[j].inputValidates {
				if err == nil {
					t.Errorf("passed test '%v' validation incorrectly: %v",
						sigScriptTests[i].name, err)
				} else {
					t.Errorf("failed test '%v' validation: %v",
						sigScriptTests[i].name, err)
				}
				continue nexttest
			}
		}
	}
}
