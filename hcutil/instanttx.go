package hcutil

import (
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/wire"
)

type InstantTx struct {
	msgTx *wire.MsgInstantTx // Underlying MsgInstantTx
}

func NewInstantTx(msgTx *wire.MsgInstantTx) *InstantTx {
	return &InstantTx{
		msgTx: msgTx,
	}
}
func (instantTx *InstantTx) Hash() *chainhash.Hash {
	ret:=instantTx.msgTx.TxHash()
	return &ret
}
func (instantTx *InstantTx) MsgTx() *wire.MsgInstantTx {
	// Return the cached transaction.
	return instantTx.msgTx
}
