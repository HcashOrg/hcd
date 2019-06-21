package wire

import (
	"bytes"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"io"
	"fmt"
)

type MsgInstantTxVote struct {
	InstantTxHash chainhash.Hash
	TicketHash    chainhash.Hash
	Vote          bool
	Sig           []byte
	PubKey        []byte
}

func NewMsgInstantTxVote() *MsgInstantTxVote {
	return &MsgInstantTxVote{}
}

func (msg *MsgInstantTxVote) Hash() *chainhash.Hash {
	ret := chainhash.HashH(msg.MustSerialize())
	return &ret
}

func (msg *MsgInstantTxVote) BtcDecode(r io.Reader, pver uint32) error {
	err := readElements(r, &msg.InstantTxHash, &msg.TicketHash, &msg.Vote)
	if err != nil {
		return err
	}
	msg.Sig, err = ReadVarBytes(r, pver, 300, "read instantTxVote sig")
	if err != nil {
		return err
	}
	msg.PubKey, err = ReadVarBytes(r, pver, 300, "read instantTxVote sig")
	if err != nil {
		return err
	}
	return nil
}

func (msg *MsgInstantTxVote) BtcEncode(w io.Writer, pver uint32) error {
	err := writeElements(w, &msg.InstantTxHash, &msg.TicketHash, msg.Vote)
	if err != nil {
		return err
	}

	err = WriteVarBytes(w, pver, msg.Sig)
	if err != nil {
		return err
	}
	return WriteVarBytes(w, pver, msg.PubKey)
}

func (msg *MsgInstantTxVote) Command() string {
	return CmdInstantTxVote
}

func (msg *MsgInstantTxVote) MaxPayloadLength(pver uint32) uint32 {
	//return
	if pver <= 3 {
		return MaxBlockPayloadV3
	}

	return MaxBlockPayload
}

// serialize returns the serialization of the transaction for the provided
// serialization type without modifying the original transaction.
func (msg *MsgInstantTxVote) serialize() ([]byte, error) {
	// Shallow copy so the serialization type can be changed without
	// modifying the original transaction.
	mtxCopy := *msg
	buf := bytes.NewBuffer(make([]byte, 0, mtxCopy.SerializeSize()))
	err := mtxCopy.Serialize(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (msg *MsgInstantTxVote) SerializeSize() int {
	return 32 + 32 + 1 + VarIntSerializeSize(uint64(len(msg.Sig))) + len(msg.Sig) + VarIntSerializeSize(uint64(len(msg.PubKey))) + len(msg.PubKey)
}

func (msg *MsgInstantTxVote) Serialize(w io.Writer) error {
	// At the current time, there is no difference between the wire encoding
	// at protocol version 0 and the stable long-term storage format.  As
	// a result, make use of BtcEncode.
	return msg.BtcEncode(w, 0)
}

// mustSerialize returns the serialization of the transaction for the provided
// serialization type without modifying the original transaction.  It will panic
// if any errors occur.
func (msg *MsgInstantTxVote) MustSerialize() []byte {
	serialized, err := msg.serialize()
	if err != nil {
		panic(fmt.Sprintf("msgInstantTx failed serializing for type"))
	}
	return serialized
}
func (msg *MsgInstantTxVote) Deserialize(r io.Reader) error {
	// At the current time, there is no difference between the wire encoding
	// at protocol version 0 and the stable long-term storage format.  As
	// a result, make use of BtcDecode.
	return msg.BtcDecode(r, 0)
}
