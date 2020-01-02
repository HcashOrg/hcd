// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
"fmt"
"io"
)

// MaxRouteAddrPerMsg is the maximum number of addresses that can be in a single
// bitcoin addr message (MsgAddr).
const MaxRouteAddrPerMsg = 1000

// MsgRouteAddr implements the Message interface and represents a bitcoin
// addr message.  It is used to provide a list of known active peers on the
// network.  An active peer is considered one that has transmitted a message
// within the last 3 hours.  Nodes which have not transmitted in that time
// frame should be forgotten.  Each message is limited to a maximum number of
// addresses, which is currently 1000.  As a result, multiple messages must
// be used to relay the full list.
//
// Use the AddAddress function to build up the list of known addresses when
// sending an addr message to another peer.
type MsgRouteAddr struct {
	AddrList []string
}

// AddAddress adds a known active peer to the message.
func (msg *MsgRouteAddr) AddAddress(na string) error {
	if len(msg.AddrList)+1 > MaxRouteAddrPerMsg {
		str := fmt.Sprintf("too many route addresses in message [max %v]",
			MaxRouteAddrPerMsg)
		return messageError("MsgRouteAddr.AddAddress", str)
	}

	msg.AddrList = append(msg.AddrList, na)
	return nil
}

// AddAddresses adds multiple known active peers to the message.
func (msg *MsgRouteAddr) AddAddresses(netAddrs ...string) error {
	for _, na := range netAddrs {
		err := msg.AddAddress(na)
		if err != nil {
			return err
		}
	}
	return nil
}

// ClearAddresses removes all addresses from the message.
func (msg *MsgRouteAddr) ClearAddresses() {
	msg.AddrList = []string{}
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgRouteAddr) BtcDecode(r io.Reader, pver uint32) error {
	count, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	// Limit to max addresses per message.
	if count > MaxRouteAddrPerMsg {
		str := fmt.Sprintf("too many route addresses for message "+
			"[count %v, max %v]", count, MaxRouteAddrPerMsg)
		return messageError("MsgRouteAddr.BtcDecode", str)
	}


	msg.AddrList = make([]string, 0, count)
	for i := uint64(0); i < count; i++ {
		addr, err := ReadVarString(r, pver)
		if err != nil {
			return err
		}

		msg.AddAddress(addr)
	}
	return nil
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgRouteAddr) BtcEncode(w io.Writer, pver uint32) error {
	// Protocol versions before MultipleAddressVersion only allowed 1 address
	// per message.
	count := len(msg.AddrList)
	if count > MaxRouteAddrPerMsg {
		str := fmt.Sprintf("too many Route addresses for message "+
			"[count %v, max %v]", count, MaxRouteAddrPerMsg)
		return messageError("MsgRouteAddr.BtcEncode", str)
	}

	err := WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, na := range msg.AddrList {
		err = WriteVarString(w,pver,na)
		if err != nil {
			return err
		}
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgRouteAddr) Command() string {
	return CmdRouteAddr
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgRouteAddr) MaxPayloadLength(pver uint32) uint32 {

	return MaxMessagePayload
}

// NewMsgAddr returns a new bitcoin addr message that conforms to the
// Message interface.  See MsgAddr for details.
func NewMsgRouteAddr() *MsgRouteAddr {
	return &MsgRouteAddr{
		AddrList: make([]string, 0, MaxRouteAddrPerMsg),
	}
}

