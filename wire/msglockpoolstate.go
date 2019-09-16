// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"

	"github.com/HcashOrg/hcd/chaincfg/chainhash"
)

// MsgLockPoolState is the maximum number of block hashes allowed
// per message.
const MaxAiTx = 20
const MaxAiTxVote = 100

// MsgLockPoolState implements the Message interface and represents a mining state
// message.  It is used to request a list of blocks located at the chain tip
// along with all votes for those blocks.  The list is returned is limited by
// the maximum number of blocks per message and the maximum number of votes per
// message.
type MsgLockPoolState struct {
	AiTxHashes     []*chainhash.Hash
	AiTxVoteHashes []*chainhash.Hash
}

// AddAiTxHash adds a new AiTx hash to the message.
func (msg *MsgLockPoolState) AddAiTxHash(hash *chainhash.Hash) error {
	if len(msg.AiTxHashes)+1 > MaxAiTx {
		str := fmt.Sprintf("too many aitx hashes for message [max %v]",
			MaxAiTx)
		return messageError("MsgLockPoolState.AddBlockHash", str)
	}

	msg.AiTxHashes = append(msg.AiTxHashes, hash)
	return nil
}

// AddVoteHash adds a new vote hash to the message.
func (msg *MsgLockPoolState) AddAiTxVoteHash(hash *chainhash.Hash) error {
	if len(msg.AiTxVoteHashes)+1 > MaxAiTxVote {
		str := fmt.Sprintf("too many vote hashes for message [max %v]",
			MaxAiTxVote)
		return messageError("MsgLockPoolState.AddVoteHash", str)
	}

	msg.AiTxVoteHashes = append(msg.AiTxVoteHashes, hash)
	return nil
}

// BtcDecode decodes r using the protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgLockPoolState) BtcDecode(r io.Reader, pver uint32) error {
	// Read num block hashes and limit to max.
	count, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}
	if count > MaxAiTx {
		str := fmt.Sprintf("too many aiTx hashes for message "+
			"[count %v, max %v]", count, MaxAiTx)
		return messageError("MsgLockPoolState.BtcDecode", str)
	}

	msg.AiTxHashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := chainhash.Hash{}
		err := readElement(r, &hash)
		if err != nil {
			return err
		}
		msg.AddAiTxHash(&hash)
	}

	// Read num vote hashes and limit to max.
	count, err = ReadVarInt(r, pver)
	if err != nil {
		return err
	}
	if count > MaxAiTxVote {
		str := fmt.Sprintf("too many vote hashes for message "+
			"[count %v, max %v]", count, MaxAiTxVote)
		return messageError("MsgLockPoolState.BtcDecode", str)
	}

	msg.AiTxVoteHashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := chainhash.Hash{}
		err := readElement(r, &hash)
		if err != nil {
			return err
		}
		err = msg.AddAiTxVoteHash(&hash)
		if err != nil {
			return err
		}
	}

	return nil
}

// BtcEncode encodes the receiver to w using the protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgLockPoolState) BtcEncode(w io.Writer, pver uint32) error {
	// Write block hashes.
	count := len(msg.AiTxHashes)
	if count > MaxAiTx {
		str := fmt.Sprintf("too many aiTx hashes for message "+
			"[count %v, max %v]", count, MaxAiTx)
		return messageError("MsgLockPoolState.BtcEncode", str)
	}

	err := WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, hash := range msg.AiTxHashes {
		err = writeElement(w, hash)
		if err != nil {
			return err
		}
	}

	// Write vote hashes.
	count = len(msg.AiTxVoteHashes)
	if count > MaxAiTxVote {
		str := fmt.Sprintf("too many vote hashes for message "+
			"[count %v, max %v]", count, MaxAiTxVote)
		return messageError("MsgLockPoolState.BtcEncode", str)
	}

	err = WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, hash := range msg.AiTxVoteHashes {
		err = writeElement(w, hash)
		if err != nil {
			return err
		}
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgLockPoolState) Command() string {
	return CmdLockPoolState
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgLockPoolState) MaxPayloadLength(pver uint32) uint32 {
	//  + num block hashes (varInt) +
	// block hashes + num vote hashes (varInt) + vote hashes
	return MaxVarIntPayload + (MaxAiTx *
		chainhash.HashSize) + MaxVarIntPayload + (MaxAiTxVote *
		chainhash.HashSize)
}

// NewMsgLockPoolState returns a new hcd MsgLockPoolState message that conforms to
// the Message interface using the defaults for the fields.
func NewMsgLockPoolState() *MsgLockPoolState {
	return &MsgLockPoolState{
		AiTxHashes:     make([]*chainhash.Hash, 0, MaxAiTx),
		AiTxVoteHashes: make([]*chainhash.Hash, 0, MaxAiTxVote),
	}
}
