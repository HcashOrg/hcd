package bliss

import (
	"github.com/HcashOrg/bliss"
	hccrypto "github.com/HcashOrg/hcd/crypto"
)

type PublicKey struct {
	hccrypto.PublicKeyAdapter
	bliss.PublicKey
}

func (p PublicKey) GetType() int {
	return pqcTypeBliss
}

func (p PublicKey) Serialize() []byte {
	return p.PublicKey.Serialize()
}

func (p PublicKey) SerializeCompressed() []byte {
	return p.Serialize()
}

func (p PublicKey) SerializeUnCompressed() []byte {
	return p.Serialize()
}
