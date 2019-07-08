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

// MaxMSBlocksAtHeadPerMsg is the maximum number of block hashes allowed
// per message.
const MaxInstantTx = 20
const MaxInstantTxVote = 100

// MsgMiningState implements the Message interface and represents a mining state
// message.  It is used to request a list of blocks located at the chain tip
// along with all votes for those blocks.  The list is returned is limited by
// the maximum number of blocks per message and the maximum number of votes per
// message.
type MsgLockPoolState struct {
	InstantTxHashes     []*chainhash.Hash
	InstantTxVoteHashes []*chainhash.Hash
}

// AddBlockHash adds a new block hash to the message.
func (msg *MsgLockPoolState) AddInstantTxHash(hash *chainhash.Hash) error {
	if len(msg.InstantTxHashes)+1 > MaxInstantTx {
		str := fmt.Sprintf("too many instanttx hashes for message [max %v]",
			MaxInstantTx)
		return messageError("MsgLockPoolState.AddBlockHash", str)
	}

	msg.InstantTxHashes = append(msg.InstantTxHashes, hash)
	return nil
}

// AddVoteHash adds a new vote hash to the message.
func (msg *MsgLockPoolState) AddInstantTxVoteHash(hash *chainhash.Hash) error {
	if len(msg.InstantTxVoteHashes)+1 > MaxInstantTxVote {
		str := fmt.Sprintf("too many vote hashes for message [max %v]",
			MaxInstantTxVote)
		return messageError("MsgLockPoolState.AddVoteHash", str)
	}

	msg.InstantTxVoteHashes = append(msg.InstantTxVoteHashes, hash)
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
	if count > MaxInstantTx {
		str := fmt.Sprintf("too many instantTx hashes for message "+
			"[count %v, max %v]", count, MaxInstantTx)
		return messageError("MsgLockPoolState.BtcDecode", str)
	}

	msg.InstantTxHashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := chainhash.Hash{}
		err := readElement(r, &hash)
		if err != nil {
			return err
		}
		msg.AddInstantTxHash(&hash)
	}

	// Read num vote hashes and limit to max.
	count, err = ReadVarInt(r, pver)
	if err != nil {
		return err
	}
	if count > MaxInstantTxVote {
		str := fmt.Sprintf("too many vote hashes for message "+
			"[count %v, max %v]", count, MaxInstantTxVote)
		return messageError("MsgLockPoolState.BtcDecode", str)
	}

	msg.InstantTxVoteHashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := chainhash.Hash{}
		err := readElement(r, &hash)
		if err != nil {
			return err
		}
		err = msg.AddInstantTxVoteHash(&hash)
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
	count := len(msg.InstantTxHashes)
	if count > MaxInstantTx {
		str := fmt.Sprintf("too many instantTx hashes for message "+
			"[count %v, max %v]", count, MaxInstantTx)
		return messageError("MsgLockPoolState.BtcEncode", str)
	}

	err := WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, hash := range msg.InstantTxHashes {
		err = writeElement(w, hash)
		if err != nil {
			return err
		}
	}

	// Write vote hashes.
	count = len(msg.InstantTxVoteHashes)
	if count > MaxInstantTxVote {
		str := fmt.Sprintf("too many vote hashes for message "+
			"[count %v, max %v]", count, MaxInstantTxVote)
		return messageError("MsgLockPoolState.BtcEncode", str)
	}

	err = WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, hash := range msg.InstantTxVoteHashes {
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
	return MaxVarIntPayload + (MaxInstantTx *
		chainhash.HashSize) + MaxVarIntPayload + (MaxInstantTxVote *
		chainhash.HashSize)
}

// NewMsgMiningState returns a new hcd MsgLockPoolState message that conforms to
// the Message interface using the defaults for the fields.
func NewMsgLockPoolState() *MsgLockPoolState {
	return &MsgLockPoolState{
		InstantTxHashes:     make([]*chainhash.Hash, 0, MaxInstantTx),
		InstantTxVoteHashes: make([]*chainhash.Hash, 0, MaxInstantTxVote),
	}
}
