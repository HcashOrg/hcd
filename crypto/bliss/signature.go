package bliss

import (
	"crypto/rand"

	"github.com/HCashOrg/bliss"
	"github.com/HCashOrg/bliss/sampler"
	hxcrypto "github.com/james-ray/hcd/crypto"
)

type Signature struct {
	hxcrypto.SignatureAdapter
	bliss.Signature
}

func (s Signature) GetType() int {
	return pqcTypeBliss
}

func (s Signature) Serialize() []byte {
	return s.Signature.Serialize()
}

func SignCompact(key hxcrypto.PrivateKey, hash []byte) ([]byte, error) {

	seed := make([]byte, sampler.SHA_512_DIGEST_LENGTH)
	rand.Read(seed)
	entropy, err := sampler.NewEntropy(seed)
	if err != nil {
		return nil, err
	}
	var sig *bliss.Signature
	switch pv := key.(type) {
	case PrivateKey:
		sig, err = pv.Sign(hash, entropy)
	case *PrivateKey:
		sig, err = pv.Sign(hash, entropy)
	}

	if err != nil {
		return nil, err
	}

	result := sig.Serialize()
	return result, err
}

func VerifyCompact(key hxcrypto.PublicKey, messageHash, sign []byte) (bool, error) {

	sig, _ := bliss.DeserializeBlissSignature(sign)
	result, err := key.(*PublicKey).Verify(messageHash, sig)

	return result, err
}
