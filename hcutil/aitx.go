package hcutil

import (
	"bytes"
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/wire"
	bs "github.com/HcashOrg/hcd/crypto/bliss"
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

func (aiTxVote *AiTxVote) GetPubKey() []byte{
	return aiTxVote.msgAiTxVote.PubKey
}

type AiTx struct {
	Tx
}

// MsgAiTx returns the underlying wire.MsgAiTx for the transaction.
func (t *AiTx) MsgAiTx() *wire.MsgAiTx {
	// Return transaction.
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



func VerifyMessage(msg string, addr Address, sig []byte, pubKey []byte) (bool, error) {
	// Validate the signature - this just shows that it was valid for any pubkey
	// at all. Whether the pubkey matches is checked below.
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Hc Signed Message:\n")
	wire.WriteVarString(&buf, 0, msg)
	expectedMessageHash := chainhash.HashB(buf.Bytes())
	pk, wasCompressed, err := chainec.Secp256k1.RecoverCompact(sig,
		expectedMessageHash)
	if err != nil {
		//maby bliss address
		pSig, err := bs.BlissDSA.ParseDERSignature(sig)
		if err != nil {
			return false, err
		}

		restoredPK, err := bs.Bliss.ParsePubKey(pubKey)
		if err != nil{
			return false, err
		}
		return bs.Bliss.Verify(restoredPK, expectedMessageHash, pSig), nil
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
	//ok := bs.Bliss.Verify(pubKey, hash, signature)
}
