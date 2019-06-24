package hcutil

import (
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/wire"
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
//func (instantTx *InstantTx) Hash() *chainhash.Hash {
//	ret:=instantTx.msgTx.TxHash()
//	return &ret
//}
//func (instantTx *InstantTx) MsgTx() *wire.MsgInstantTx {
//	// Return the cached transaction.
//	return instantTx.msgTx
//}

