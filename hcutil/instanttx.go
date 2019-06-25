package hcutil

import (
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/wire"
	"bytes"
)


//
//type InstantTx struct {
//	msgTx *wire.MsgInstantTx // Underlying MsgInstantTx
//}

type InstantTxVote struct {
	msgInstantTxVote *wire.MsgInstantTxVote
}

func NewInstantTxVote(vote *wire.MsgInstantTxVote) *InstantTxVote {
	return &InstantTxVote{
		msgInstantTxVote:vote,
	}
}

func (instantTxVote *InstantTxVote)Hash() *chainhash.Hash {
	return instantTxVote.msgInstantTxVote.Hash()
}

func (instantTxVote *InstantTxVote)MsgInstantTxVote()*wire.MsgInstantTxVote  {
	return instantTxVote.msgInstantTxVote
}



type InstantTx struct {
	Tx
}


func NewInstantTx(msgInstantTx *wire.MsgInstantTx) *InstantTx {
	return &InstantTx{
		Tx:Tx{
			hash:    msgInstantTx.TxHash(),
			msgTx:   &msgInstantTx.MsgTx,
			txTree:  wire.TxTreeUnknown,
			txIndex: TxIndexUnknown,
		},
	}
}

func NewInstantTxFromTx(tx *Tx)*InstantTx  {
	return &InstantTx{
		Tx:*tx,
	}
}

//func (instantTx *InstantTx) Hash() *chainhash.Hash {
//	ret:=instantTx.msgTx.TxHash()
//	return &ret
//}
//func (instantTx *InstantTx) MsgTx() *wire.MsgInstantTx {
//	// Return the cached transaction.
//	return instantTx.msgTx
//}


func VerifyMessage(msg string, addr Address, sig []byte) (bool, error) {
	// Validate the signature - this just shows that it was valid for any pubkey
	// at all. Whether the pubkey matches is checked below.
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Hc Signed Message:\n")
	wire.WriteVarString(&buf, 0, msg)
	expectedMessageHash := chainhash.HashB(buf.Bytes())
	pk, wasCompressed, err := chainec.Secp256k1.RecoverCompact(sig,
		expectedMessageHash)
	if err != nil {
		return false, err
	}

	// Reconstruct the address from the recovered pubkey.
	var serializedPK []byte
	if wasCompressed {
		serializedPK = pk.SerializeCompressed()
	} else {
		serializedPK = pk.SerializeUncompressed()
	}
	recoveredAddr, err := NewAddressSecpPubKey(serializedPK, addr.Net())
	if err != nil {
		return false, err
	}

	// Return whether addresses match.
	return recoveredAddr.EncodeAddress() == addr.EncodeAddress(), nil
}