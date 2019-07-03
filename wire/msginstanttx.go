// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

type MsgInstantTx struct {
	MsgTx
}

func NewMsgInstantTx() *MsgInstantTx {
	return &MsgInstantTx{
		MsgTx: *NewMsgTx(),
	}
}

func NewMsgInstantTxFromMsgTx(msgTx *MsgTx) *MsgInstantTx {
	return &MsgInstantTx{
		MsgTx: *msgTx,
	}
}


func (msg *MsgInstantTx) Command() string {
	return CmdInstantTx
}

