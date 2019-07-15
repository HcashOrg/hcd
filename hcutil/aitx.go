package hcutil

import (
	"bytes"
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/wire"
)



type AiTxVote struct {
	msgAiTxVote *wire.MsgAiTxVote
}

func NewAiTxVote(vote *wire.MsgAiTxVote) *AiTxVote {
	return &AiTxVote{
		msgAiTxVote: vote,
	}
}

func (aiTxVote *AiTxVote) Hash() *chainhash.Hash {
	return aiTxVote.msgAiTxVote.Hash()
}

func (aiTxVote *AiTxVote) MsgAiTxVote() *wire.MsgAiTxVote {
	return aiTxVote.msgAiTxVote
}

type AiTx struct {
	Tx
}

// MsgTx returns the underlying wire.MsgTx for the transaction.
func (t *AiTx) MsgAiTx() *wire.MsgAiTx {
	// Return the cached transaction.
	return wire.NewMsgAiTxFromMsgTx(t.msgTx)
}

func NewAiTx(msgAiTx *wire.MsgAiTx) *AiTx {
	return &AiTx{
		Tx: Tx{
			hash:    msgAiTx.TxHash(),
			msgTx:   &msgAiTx.MsgTx,
			txTree:  wire.TxTreeUnknown,
			txIndex: TxIndexUnknown,
		},
	}
}

func NewAiTxFromTx(tx *Tx) *AiTx {
	return &AiTx{
		Tx: *tx,
	}
}



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
