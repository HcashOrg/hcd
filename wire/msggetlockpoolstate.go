// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"io"
)

// MsgGetMiningState implements the Message interface and represents a
// getminingstate message.  It is used to request the current mining state
// from a peer.
type MsgGetLockPoolState struct{}

// BtcDecode decodes r using the hcd protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetLockPoolState) BtcDecode(r io.Reader, pver uint32) error {
	return nil
}

// BtcEncode encodes the receiver to w using the hcd protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgGetLockPoolState) BtcEncode(w io.Writer, pver uint32) error {
	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgGetLockPoolState) Command() string {
	return CmdGetLockPoolState
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetLockPoolState) MaxPayloadLength(pver uint32) uint32 {
	return 0
}

// NewMsgGetMiningState returns a new hcd pong message that conforms to the Message
// interface.  See MsgPong for details.
func NewMsgGetLockPoolState() *MsgGetLockPoolState {
	return &MsgGetLockPoolState{}
}
