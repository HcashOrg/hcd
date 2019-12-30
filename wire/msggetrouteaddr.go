// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
"io"
)

// MsgGetAddr implements the Message interface and represents a hcd
// getaddr message.  It is used to request a list of known active peers on the
// network from a peer to help identify potential nodes.  The list is returned
// via one or more addr messages (MsgAddr).
//
// This message has no payload.
type MsgGetRouteAddr struct{}

// BtcDecode decodes r using the hcd protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetRouteAddr) BtcDecode(r io.Reader, pver uint32) error {
	return nil
}

// BtcEncode encodes the receiver to w using the hcd protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgGetRouteAddr) BtcEncode(w io.Writer, pver uint32) error {
	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgGetRouteAddr) Command() string {
	return CmdGetRouteAddr
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetRouteAddr) MaxPayloadLength(pver uint32) uint32 {
	return 0
}

// NewMsgGetAddr returns a new hcd getaddr message that conforms to the
// Message interface.  See MsgGetAddr for details.
func NewMsgGetRouteAddr() *MsgGetRouteAddr {
	return &MsgGetRouteAddr{}
}

