// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

type MsgAiTx struct {
	MsgTx
}

func NewMsgAiTx() *MsgAiTx {
	return &MsgAiTx{
		MsgTx: *NewMsgTx(),
	}
}

func NewMsgAiTxFromMsgTx(msgTx *MsgTx) *MsgAiTx {
	return &MsgAiTx{
		MsgTx: *msgTx,
	}
}


func (msg *MsgAiTx) Command() string {
	return CmdAiTx
}

