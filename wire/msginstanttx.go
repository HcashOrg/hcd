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
		MsgTx:*NewMsgTx(),
	}
}

func NewMsgInstantTxFromMsgTx(msgTx *MsgTx) *MsgInstantTx{
	return &MsgInstantTx{
		MsgTx:*msgTx,
	}
}
//func (msg *MsgInstantTx)BtcDecode(r io.Reader, pver uint32) error{
//	err:=readElement(r,&msg.LotteryHash)
//	if err!=nil{
//		return err
//	}
//	return msg.MsgTx.BtcDecode(r,pver)
//}
//
//func (msg *MsgInstantTx)BtcEncode(w io.Writer,pver uint32) error{
//	err:=writeElement(w,&msg.LotteryHash)
//	if err!=nil {
//		return err
//	}
//	return msg.MsgTx.BtcEncode(w,pver)
//}

func (msg *MsgInstantTx)Command() string{
	return CmdInstantTx
}

//func (msg *MsgInstantTx)MaxPayloadLength(pver uint32) uint32{
//	return msg.MsgTx.MaxPayloadLength(pver)
//}


// serialize returns the serialization of the transaction for the provided
// serialization type without modifying the original transaction.
//func (msg *MsgInstantTx) serialize(serType TxSerializeType) ([]byte, error) {
//	// Shallow copy so the serialization type can be changed without
//	// modifying the original transaction.
//	mtxCopy := *msg
//	mtxCopy.SerType = serType
//	buf := bytes.NewBuffer(make([]byte, 0, mtxCopy.SerializeSize()))
//	err := mtxCopy.Serialize(buf)
//	if err != nil {
//		return nil, err
//	}
//	return buf.Bytes(), nil
//}
//
//
//func (msg *MsgInstantTx) SerializeSize() int {
//	return chainhash.HashSize+msg.MsgTx.SerializeSize()
//}
//
//func (msg *MsgInstantTx) Serialize(w io.Writer) error {
//	// At the current time, there is no difference between the wire encoding
//	// at protocol version 0 and the stable long-term storage format.  As
//	// a result, make use of BtcEncode.
//	return msg.BtcEncode(w, 0)
//}
//
//
//// mustSerialize returns the serialization of the transaction for the provided
//// serialization type without modifying the original transaction.  It will panic
//// if any errors occur.
//func (msg *MsgInstantTx) mustSerialize(serType TxSerializeType) []byte {
//	serialized, err := msg.serialize(serType)
//	if err != nil {
//		panic(fmt.Sprintf("msgTx failed serializing for type %v",
//			serType))
//	}
//	return serialized
//}
//
//// TxHash generates the hash for the transaction prefix.  Since it does not
//// contain any witness data, it is not malleable and therefore is stable for
//// use in unconfirmed transaction chains.
//func (msg *MsgInstantTx) TxHash() chainhash.Hash {
//	// TxHash should always calculate a non-witnessed hash.
//	return chainhash.HashH(msg.mustSerialize(TxSerializeNoWitness))
//}
//
//
//func (msg *MsgInstantTx) Deserialize(r io.Reader) error {
//	// At the current time, there is no difference between the wire encoding
//	// at protocol version 0 and the stable long-term storage format.  As
//	// a result, make use of BtcDecode.
//	return msg.BtcDecode(r, 0)
//}
