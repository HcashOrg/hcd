// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Copyright (c) 2018-2020 The Hcd developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hdkeychain

// References:
//   [BIP32]: BIP0032 - Hierarchical Deterministic Wallets
//   https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/HCashOrg/bliss"
	"github.com/HCashOrg/bliss/sampler"
	"github.com/james-ray/hcd/chaincfg"
	"github.com/james-ray/hcd/chaincfg/chainec"
	"github.com/james-ray/hcd/chaincfg/chainhash"
	hccrypto "github.com/james-ray/hcd/crypto/bliss"
	"github.com/james-ray/hcd/hcutil"
	"github.com/james-ray/hcd/hcutil/base58"
	"golang.org/x/crypto/sha3"
)

const (
	// RecommendedSeedLen is the recommended length in bytes for a seed
	// to a master node.
	RecommendedSeedLen = 32 // 256 bits

	// HardenedKeyStart is the index at which a hardended key starts.  Each
	// extended key has 2^31 normal child keys and 2^31 hardned child keys.
	// Thus the range for normal child keys is [0, 2^31 - 1] and the range
	// for hardened child keys is [2^31, 2^32 - 1].
	HardenedKeyStart = 0x80000000 // 2^31

	// MinSeedBytes is the minimum number of bytes allowed for a seed to
	// a master node.
	MinSeedBytes = 16 // 128 bits

	// MaxSeedBytes is the maximum number of bytes allowed for a seed to
	// a master node.
	MaxSeedBytes = 64 // 512 bits

	// serializedKeyLen is the length of a serialized public or private
	// extended key.  It consists of 4 bytes version, 1 byte depth, 4 bytes
	// fingerprint, 4 bytes child number, 32 bytes chain code, and 33 bytes
	// public/private key data.
	serializedKeyLen                = 4 + 1 + 1 + 4 + 4 + 32 + 33 // 79 bytes
	serializedKeyLenForTest         = 4 + 1 + 4 + 4 + 32 + 33     // 78 bytes
	blissserializedPubKeyLen        = 4 + 1 + 1 + 4 + 4 + 32 + 897
	blissserializedPrivKeyLen       = 4 + 1 + 1 + 4 + 4 + 32 + 386
	keyEc                     uint8 = 0
	keyBliss                  uint8 = 1
	BlissPubKeyLen                  = 897
)

var (
	// ErrDeriveHardFromPublic describes an error in which the caller
	// attempted to derive a hardened extended key from a public key.
	ErrDeriveHardFromPublic = errors.New("cannot derive a hardened key " +
		"from a public key")

	ErrUnknownAlg = errors.New("unkown algtype")

	ErrDerivePublicFromPublic = errors.New("cannot derive a public key " +
		"from a public key")

	// ErrNotPrivExtKey describes an error in which the caller attempted
	// to extract a private key from a public extended key.
	ErrNotPrivExtKey = errors.New("unable to create private keys from a " +
		"public extended key")

	// ErrInvalidChild describes an error in which the child at a specific
	// index is invalid due to the derived key falling outside of the valid
	// range for secp256k1 private keys.  This error indicates the caller
	// should simply ignore the invalid child extended key at this index and
	// increment to the next index.
	ErrInvalidChild = errors.New("the extended key at this index is invalid")

	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// secp256k1 private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = errors.New("unusable seed")

	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf("seed length must be between %d and %d "+
		"bits", MinSeedBytes*8, MaxSeedBytes*8)

	// ErrBadChecksum describes an error in which the checksum encoded with
	// a serialized extended key does not match the calculated value.
	ErrBadChecksum = errors.New("bad extended key checksum")

	// ErrInvalidKeyLen describes an error in which the provided serialized
	// key is not the expected length.
	ErrInvalidKeyLen = errors.New("the provided serialized extended key " +
		"length is invalid")
)

// masterKey is the master key used along with a random seed used to generate
// the master node in the hierarchical tree.
var masterKey = []byte("Bitcoin seed")

// ExtendedKey houses all the information needed to support a hierarchical
// deterministic extended key.  See the package overview documentation for
// more details on how to use extended keys.
type ExtendedKey struct {
	key       []byte // This will be the pubkey for extended pub keys
	pubKey    []byte // This will only be set for extended priv keys
	chainCode []byte
	depth     uint16
	parentFP  []byte
	childNum  uint32
	version   []byte
	isPrivate bool
	algtype   uint8
}

// newExtendedKey returns a new instance of an extended key with the given
// fields.  No error checking is performed here as it's only intended to be a
// convenience method used to create a populated struct.
func newExtendedKey(version, key, chainCode, parentFP []byte, depth uint16,
	childNum uint32, isPrivate bool, algtype uint8) *ExtendedKey {

	// NOTE: The pubKey field is intentionally left nil so it is only
	// computed and memoized as required.
	return &ExtendedKey{
		key:       key,
		chainCode: chainCode,
		depth:     depth,
		parentFP:  parentFP,
		childNum:  childNum,
		version:   version,
		isPrivate: isPrivate,
		algtype:   algtype,
	}
}

// pubKeyBytes returns bytes for the serialized compressed public key associated
// with this extended key in an efficient manner including memoization as
// necessary.
//
// When the extended key is already a public key, the key is simply returned as
// is since it's already in the correct form.  However, when the extended key is
// a private key, the public key will be calculated and memoized so future
// accesses can simply return the cached result.
func (k *ExtendedKey) pubKeyBytes() []byte {
	// Just return the key if it's already an extended public key.
	if !k.isPrivate {
		return k.key
	}

	// This is a private extended key, so calculate and memoize the public
	// key if needed.
	if len(k.pubKey) == 0 {
		if k.algtype == keyBliss {
			privkey, err := bliss.DeserializePrivateKey(k.key)
			if err != nil {
				return nil
			}
			k.pubKey = privkey.PublicKey().Serialize()

		} else {
			pkx, pky := chainec.Secp256k1.ScalarBaseMult(k.key)
			pubKey := chainec.Secp256k1.NewPublicKey(pkx, pky)
			k.pubKey = pubKey.SerializeCompressed()
		}
	}

	return k.pubKey
}

// IsPrivate returns whether or not the extended key is a private extended key.
//
// A private extended key can be used to derive both hardened and non-hardened
// child private and public extended keys.  A public extended key can only be
// used to derive non-hardened child public extended keys.
func (k *ExtendedKey) IsPrivate() bool {
	return k.isPrivate
}

// ParentFingerprint returns a fingerprint of the parent extended key from which
// this one was derived.
func (k *ExtendedKey) ParentFingerprint() uint32 {
	return binary.BigEndian.Uint32(k.parentFP)
}

// Child returns a derived child extended key at the given index.  When this
// extended key is a private extended key (as determined by the IsPrivate
// function), a private extended key will be derived.  Otherwise, the derived
// extended key will be also be a public extended key.
//
// When the index is greater to or equal than the HardenedKeyStart constant, the
// derived extended key will be a hardened extended key.  It is only possible to
// derive a hardended extended key from a private extended key.  Consequently,
// this function will return ErrDeriveHardFromPublic if a hardened child
// extended key is requested from a public extended key.
//
// A hardened extended key is useful since, as previously mentioned, it requires
// a parent private extended key to derive.  In other words, normal child
// extended public keys can be derived from a parent public extended key (no
// knowledge of the parent private key) whereas hardened extended keys may not
// be.
//
// NOTE: There is an extremely small chance (< 1 in 2^127) the specific child
// index does not derive to a usable child.  The ErrInvalidChild error will be
// returned if this should occur, and the caller is expected to ignore the
// invalid child and simply increment to the next index.
func (k *ExtendedKey) Child(i uint32) (*ExtendedKey, error) {
	var isPrivate bool
	var childKey []byte
	childChainCode := make([]byte, 32)
	switch k.algtype {
	case keyBliss:
		if !k.isPrivate {
			return nil, ErrDerivePublicFromPublic
		}
		isPrivate = true
		keyLen := BlissPubKeyLen
		data := make([]byte, keyLen+4)
		copy(data, k.pubKeyBytes())
		binary.BigEndian.PutUint32(data[keyLen:], i)
		hmac512 := hmac.New(sha512.New, k.chainCode)
		hmac512.Write(data)
		ilr := hmac512.Sum(nil)
		il := ilr[:len(ilr)/2]
		childChainCode = ilr[len(ilr)/2:]
		entropyrand := sha3.Sum512(il)
		entropy, err := sampler.NewEntropy(entropyrand[:])
		if err != nil {
			return nil, err
		}
		privKey, err := bliss.GeneratePrivateKey(1, entropy)
		if err != nil && strings.Contains(err.Error(), "invertible polynomial") {
			return nil, ErrInvalidChild
		}
		if err != nil {
			return nil, err
		}
		childKey = privKey.Serialize()
	default:
		// There are four scenarios that could happen here:
		// 1) Private extended key -> Hardened child private extended key
		// 2) Private extended key -> Non-hardened child private extended key
		// 3) Public extended key -> Non-hardened child public extended key
		// 4) Public extended key -> Hardened child public extended key (INVALID!)

		// Case #4 is invalid, so error out early.
		// A hardened child extended key may not be created from a public
		// extended key.
		isChildHardened := i >= HardenedKeyStart
		k.algtype = keyEc
		if !k.isPrivate && isChildHardened {
			return nil, ErrDeriveHardFromPublic
		}

		// The data used to derive the child key depends on whether or not the
		// child is hardened per [BIP32].
		//
		// For hardened children:
		//   0x00 || ser256(parentKey) || ser32(i)
		//
		// For normal children:
		//   serP(parentPubKey) || ser32(i)
		keyLen := 33
		data := make([]byte, keyLen+4)
		if isChildHardened {
			// Case #1.
			// When the child is a hardened child, the key is known to be a
			// private key due to the above early return.  Pad it with a
			// leading zero as required by [BIP32] for deriving the child.
			copy(data[1:], k.key)
		} else {
			// Case #2 or #3.
			// This is either a public or private extended key, but in
			// either case, the data which is used to derive the child key
			// starts with the secp256k1 compressed public key bytes.
			copy(data, k.pubKeyBytes())
		}
		binary.BigEndian.PutUint32(data[keyLen:], i)

		// Take the HMAC-SHA512 of the current key's chain code and the derived
		// data:
		//   I = HMAC-SHA512(Key = chainCode, Data = data)
		hmac512 := hmac.New(sha512.New, k.chainCode)
		hmac512.Write(data)
		ilr := hmac512.Sum(nil)
		// Split "I" into two 32-byte sequences Il and Ir where:
		//   Il = intermediate key used to derive the child
		//   Ir = child chain code
		il := ilr[:len(ilr)/2]
		copy(childChainCode, ilr[len(ilr)/2:])
		// Both derived public or private keys rely on treating the left 32-byte
		// sequence calculated above (Il) as a 256-bit integer that must be
		// within the valid range for a secp256k1 private key.  There is a small
		// chance (< 1 in 2^127) this condition will not hold, and in that case,
		// a child extended key can't be created for this index and the caller
		// should simply increment to the next index.
		ilNum := new(big.Int).SetBytes(il)
		if ilNum.Cmp(chainec.Secp256k1.GetN()) >= 0 || ilNum.Sign() == 0 {
			return nil, ErrInvalidChild
		}

		// The algorithm used to derive the child key depends on whether or not
		// a private or public child is being derived.
		//
		// For private children:
		//   childKey = parse256(Il) + parentKey
		//
		// For public children:
		//   childKey = serP(point(parse256(Il)) + parentKey)
		if k.isPrivate {
			// Case #1 or #2.
			// Add the parent private key to the intermediate private key to
			// derive the final child key.
			//
			// childKey = parse256(Il) + parenKey
			keyNum := new(big.Int).SetBytes(k.key)
			ilNum.Add(ilNum, keyNum)
			ilNum.Mod(ilNum, chainec.Secp256k1.GetN())
			childKey = ilNum.Bytes()
			isPrivate = true
		} else {
			// Case #3.
			// Calculate the corresponding intermediate public key for
			// intermediate private key.
			ilx, ily := chainec.Secp256k1.ScalarBaseMult(il)
			if ilx.Sign() == 0 || ily.Sign() == 0 {
				return nil, ErrInvalidChild
			}

			// Convert the serialized compressed parent public key into X
			// and Y coordinates so it can be added to the intermediate
			// public key.
			pubKey, err := chainec.Secp256k1.ParsePubKey(k.key)
			if err != nil {
				return nil, err
			}

			// Add the intermediate public key to the parent public key to
			// derive the final child key.
			//
			// childKey = serP(point(parse256(Il)) + parentKey)
			childX, childY := chainec.Secp256k1.Add(ilx, ily, pubKey.GetX(),
				pubKey.GetY())
			pk := chainec.Secp256k1.NewPublicKey(childX, childY)
			childKey = pk.SerializeCompressed()
		}
	}
	// The fingerprint of the parent for the derived child is the first 4
	// bytes of the RIPEMD160(SHA256(parentPubKey)).
	parentFP := hcutil.Hash160(k.pubKeyBytes())[:4]
	return newExtendedKey(k.version, childKey, childChainCode, parentFP,
		k.depth+1, i, isPrivate, k.algtype), nil
}

// Neuter returns a new extended public key from this extended private key.  The
// same extended key will be returned unaltered if it is already an extended
// public key.
//
// As the name implies, an extended public key does not have access to the
// private key, so it is not capable of signing transactions or deriving
// child extended private keys.  However, it is capable of deriving further
// child extended public keys.
func (k *ExtendedKey) Neuter() (*ExtendedKey, error) {
	// Already an extended public key.
	if !k.isPrivate {
		return k, nil
	}

	// Get the associated public extended key version bytes.
	version, err := chaincfg.HDPrivateKeyToPublicKeyID(k.version)
	if err != nil {
		return nil, err
	}

	// Convert it to an extended public key.  The key for the new extended
	// key will simply be the pubkey of the current extended private key.
	//
	// This is the function N((k,c)) -> (K, c) from [BIP32].
	return newExtendedKey(version, k.pubKeyBytes(), k.chainCode, k.parentFP,
		k.depth, k.childNum, false, k.algtype), nil
}

// ECPubKey converts the extended key to a hcec public key and returns it.
func (k *ExtendedKey) ECPubKey() (chainec.PublicKey, error) {
	if k.algtype == keyEc {
		return chainec.Secp256k1.ParsePubKey(k.pubKeyBytes())
	} else if k.algtype == keyBliss {
		return hccrypto.Bliss.ParsePubKey(k.pubKeyBytes())
	}
	return nil, ErrUnknownAlg
}

// ECPrivKey converts the extended key to a hcec private key and returns it.
// As you might imagine this is only possible if the extended key is a private
// extended key (as determined by the IsPrivate function).  The ErrNotPrivExtKey
// error will be returned if this function is called on a public extended key.
func (k *ExtendedKey) ECPrivKey() (chainec.PrivateKey, error) {
	if !k.isPrivate {
		return nil, ErrNotPrivExtKey
	}

	if k.algtype == keyEc {
		privKey, _ := chainec.Secp256k1.PrivKeyFromBytes(k.key)
		return privKey, nil
	} else if k.algtype == keyBliss {
		privKey, _ := hccrypto.Bliss.PrivKeyFromBytes(k.key)
		return privKey, nil
	}
	return nil, ErrUnknownAlg
}

// Address converts the extended key to a standard hcd pay-to-pubkey-hash
// address for the passed network.
func (k *ExtendedKey) Address(net *chaincfg.Params, addrtype uint8) (*hcutil.AddressPubKeyHash, error) {
	pkHash := hcutil.Hash160(k.pubKeyBytes())
	if addrtype == 1 {
		return hcutil.NewAddressPubKeyHash(pkHash, net, 4)
	}
	return hcutil.NewAddressPubKeyHash(pkHash, net, chainec.ECTypeSecp256k1)
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// String returns the extended key as a human-readable base58-encoded string.
func (k *ExtendedKey) String() (string, error) {
	if len(k.key) == 0 {
		return "", fmt.Errorf("zeroed extended key")
	}

	var childNumBytes [4]byte
	depthByte := byte(k.depth % 256)
	typeByte := byte(k.algtype)
	binary.BigEndian.PutUint32(childNumBytes[:], k.childNum)

	// The serialized format is:
	//   version (4) || depth (1) || parent fingerprint (4)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)
	serializedBytes := make([]byte, 0, serializedKeyLen+4)
	serializedBytes = append(serializedBytes, k.version...)
	serializedBytes = append(serializedBytes, depthByte)
	if k.algtype == keyBliss {
		serializedBytes = append(serializedBytes, typeByte)
	}
	serializedBytes = append(serializedBytes, k.parentFP...)
	serializedBytes = append(serializedBytes, childNumBytes[:]...)
	serializedBytes = append(serializedBytes, k.chainCode...)
	if k.isPrivate {
		if k.algtype == keyEc {
			serializedBytes = append(serializedBytes, 0x00)
			serializedBytes = paddedAppend(32, serializedBytes, k.key)
		} else if k.algtype == keyBliss {
			serializedBytes = append(serializedBytes, 0x00)
			serializedBytes = append(serializedBytes, k.key[:]...)
		} else {
			return "", ErrUnknownAlg
		}
	} else {
		serializedBytes = append(serializedBytes, k.pubKeyBytes()...)
	}

	checkSum := chainhash.HashB(chainhash.HashB(serializedBytes))[:4]
	serializedBytes = append(serializedBytes, checkSum...)
	return base58.Encode(serializedBytes), nil
}

// IsForNet returns whether or not the extended key is associated with the
// passed hcd network.
func (k *ExtendedKey) IsForNet(net *chaincfg.Params) bool {
	return bytes.Equal(k.version, net.HDPrivateKeyID[:]) ||
		bytes.Equal(k.version, net.HDPublicKeyID[:])
}

// SetNet associates the extended key, and any child keys yet to be derived from
// it, with the passed network.
func (k *ExtendedKey) SetNet(net *chaincfg.Params) {
	if k.isPrivate {
		k.version = net.HDPrivateKeyID[:]
	} else {
		k.version = net.HDPublicKeyID[:]
	}
}

// zero sets all bytes in the passed slice to zero.  This is used to
// explicitly clear private key material from memory.
func zero(b []byte) {
	lenb := len(b)
	for i := 0; i < lenb; i++ {
		b[i] = 0
	}
}

// Zero manually clears all fields and bytes in the extended key.  This can be
// used to explicitly clear key material from memory for enhanced security
// against memory scraping.  This function only clears this particular key and
// not any children that have already been derived.
func (k *ExtendedKey) Zero() {
	zero(k.key)
	zero(k.pubKey)
	zero(k.chainCode)
	zero(k.parentFP)
	k.version = nil
	k.key = nil
	k.depth = 0
	k.childNum = 0
	k.algtype = 0
	k.isPrivate = false
}

func (k *ExtendedKey) GetAlgType() uint8 {
	return k.algtype
}

func (k *ExtendedKey) SetAlgType(i uint8) {
	k.algtype = i
}

// NewMaster creates a new master node for use in creating a hierarchical
// deterministic key chain.  The seed must be between 128 and 512 bits and
// should be generated by a cryptographically secure random generation source.
//
// NOTE: There is an extremely small chance (< 1 in 2^127) the provided seed
// will derive to an unusable secret key.  The ErrUnusable error will be
// returned if this should occur, so the caller must check for it and generate a
// new seed accordingly.
func NewMaster(seed []byte, net *chaincfg.Params) (*ExtendedKey, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "Bitcoin seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	hmac512.Write(seed)
	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master secret key
	//   Ir = master chain code
	secretKey := lr[:len(lr)/2]
	chainCode := lr[len(lr)/2:]

	// Ensure the key in usable.
	secretKeyNum := new(big.Int).SetBytes(secretKey)
	if secretKeyNum.Cmp(chainec.Secp256k1.GetN()) >= 0 ||
		secretKeyNum.Sign() == 0 {
		return nil, ErrUnusableSeed
	}

	parentFP := []byte{0x00, 0x00, 0x00, 0x00}
	return newExtendedKey(net.HDPrivateKeyID[:], secretKey, chainCode,
		parentFP, 0, 0, true, 0), nil
}

// NewKeyFromString returns a new extended key instance from a base58-encoded
// extended key.
func NewKeyFromString(key string) (*ExtendedKey, error) {
	// The base58-decoded extended key must consist of a serialized payload
	// plus an additional 4 bytes for the checksum.
	decoded := base58.Decode(key)
	if len(decoded) != serializedKeyLen+4 && len(decoded) != blissserializedPubKeyLen+4 && len(decoded) != blissserializedPrivKeyLen+4 && len(decoded) != serializedKeyLenForTest+4 {
		return nil, ErrInvalidKeyLen
	}

	// The serialized format is:
	//   version (4) || depth (1) || parent fingerprint (4)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)

	// Split the payload and checksum up and ensure the checksum matches.
	payload := decoded[:len(decoded)-4]
	checkSum := decoded[len(decoded)-4:]
	expectedCheckSum := chainhash.HashB(chainhash.HashB(payload))[:4]
	if !bytes.Equal(checkSum, expectedCheckSum) {
		return nil, ErrBadChecksum
	}

	// Deserialize each of the payload fields.
	version := payload[:4]
	depth := uint16(payload[4:5][0])
	algtype := uint8(payload[5])
	parentFP := payload[6:10]
	childNum := binary.BigEndian.Uint32(payload[10:14])
	chainCode := payload[14:46]
	keyData := payload[46:]
	if len(decoded) == serializedKeyLenForTest+4 {
		version = payload[:4]
		depth = uint16(payload[4:5][0])
		algtype = uint8(keyEc)
		parentFP = payload[5:9]
		childNum = binary.BigEndian.Uint32(payload[9:13])
		chainCode = payload[13:45]
		keyData = payload[45:]
	}
	// The key data is a private key if it starts with 0x00.  Serialized
	// compressed pubkeys either start with 0x02 or 0x03.
	isPrivate := keyData[0] == 0x00
	if algtype == keyBliss {
		isPrivate = len(decoded) == blissserializedPrivKeyLen+4
	}
	if isPrivate {
		// Ensure the private key is valid.  It must be within the range
		// of the order of the secp256k1 curve and not be 0.
		keyData = keyData[1:]
		switch {
		case algtype == keyEc:
			keyNum := new(big.Int).SetBytes(keyData)
			if keyNum.Cmp(chainec.Secp256k1.GetN()) >= 0 || keyNum.Sign() == 0 {
				return nil, ErrUnusableSeed
			}
		case algtype == keyBliss:
			//TODO
		default:
			return nil, ErrUnknownAlg
		}
	} else {
		switch {
		case algtype == keyEc:
			// Ensure the public key parses correctly and is actually on the
			// secp256k1 curve.
			_, err := chainec.Secp256k1.ParsePubKey(keyData)
			if err != nil {
				return nil, err
			}
		case algtype == keyBliss:
			// TODO
			_, err := hccrypto.Bliss.ParsePubKey(keyData)
			if err != nil {
				return nil, err
			}
		default:
			return nil, ErrUnknownAlg
		}
	}

	return newExtendedKey(version, keyData, chainCode, parentFP, depth,
		childNum, isPrivate, algtype), nil
}

// GenerateSeed returns a cryptographically secure random seed that can be used
// as the input for the NewMaster function to generate a new master node.
//
// The length is in bytes and it must be between 16 and 64 (128 to 512 bits).
// The recommended length is 32 (256 bits) as defined by the RecommendedSeedLen
// constant.
func GenerateSeed(length uint8) ([]byte, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if length < MinSeedBytes || length > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (k *ExtendedKey) SwitchChild(i uint32, acctype uint8) (*ExtendedKey, error) {
	var childKey []byte
	var isPrivate = true
	childChainCode := make([]byte, 32)
	if k.algtype != keyEc && !k.isPrivate {
		return nil, ErrDerivePublicFromPublic
	}
	keyLen := 33
	data := make([]byte, keyLen+4)
	copy(data[1:], k.key)
	binary.BigEndian.PutUint32(data[keyLen:], i)
	hmac512 := hmac.New(sha512.New, k.chainCode)
	hmac512.Write(data)
	ilr := hmac512.Sum(nil)
	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = intermediate key used to derive the child
	//   Ir = child chain code
	il := ilr[:len(ilr)/2]
	copy(childChainCode, ilr[len(ilr)/2:])
	switch acctype {
	case keyBliss:
		entropyrand := sha3.Sum512(il)
		entropy, err := sampler.NewEntropy(entropyrand[:])
		if err != nil {
			return nil, err
		}
		privKey, err := bliss.GeneratePrivateKey(1, entropy)
		if err != nil {
			return nil, err
		}
		childKey = privKey.Serialize()

	default:
		// Both derived public or private keys rely on treating the left 32-byte
		// sequence calculated above (Il) as a 256-bit integer that must be
		// within the valid range for a secp256k1 private key.  There is a small
		// chance (< 1 in 2^127) this condition will not hold, and in that case,
		// a child extended key can't be created for this index and the caller
		// should simply increment to the next index.
		ilNum := new(big.Int).SetBytes(il)
		if ilNum.Cmp(chainec.Secp256k1.GetN()) >= 0 || ilNum.Sign() == 0 {
			return nil, ErrInvalidChild
		}
		keyNum := new(big.Int).SetBytes(k.key)
		ilNum.Add(ilNum, keyNum)
		ilNum.Mod(ilNum, chainec.Secp256k1.GetN())
		childKey = ilNum.Bytes()

	}
	// The fingerprint of the parent for the derived child is the first 4
	// bytes of the RIPEMD160(SHA256(parentPubKey)).
	parentFP := hcutil.Hash160(k.pubKeyBytes())[:4]
	return newExtendedKey(k.version, childKey, childChainCode, parentFP,
		k.depth+1, i, isPrivate, acctype), nil
}
