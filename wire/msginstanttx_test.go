package wire

import (
	"bytes"
	"testing"
)

func TestInstantTx(t *testing.T) {
	pver := ProtocolVersion

	// Block 100000 hash.
	//hashStr := "3ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"
	//hash, err := chainhash.NewHashFromStr(hashStr)
	//if err != nil {
	//	t.Errorf("NewHashFromStr: %v", err)
	//}

	// Ensure the command is expected value.
	wantCmd := "instanttx"
	msg := NewMsgInstantTx()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgAddr: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	// Num addresses (varInt) + max allowed addresses.
	wantPayload := uint32(1310720)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}

	// Ensure max payload is expected value for protocol version 3.
	wantPayload = uint32(1000000)
	maxPayload = msg.MaxPayloadLength(3)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", 3,
			maxPayload, wantPayload)
	}
}

func TestInstantTxDecode(t *testing.T) {
	msg := NewMsgInstantTx()
	var buf bytes.Buffer
	err := msg.BtcEncode(&buf, ProtocolVersion)
	if err != nil {
		t.Errorf("BtcEncode error %v", err)
	}

	var msg2 MsgInstantTx
	rbuf := bytes.NewReader(buf.Bytes())
	err = msg2.BtcDecode(rbuf, ProtocolVersion)
	if err != nil {
		t.Errorf("BtcDecode error %v",err)
	}
}

func TestInstantTxVoteDecode(t *testing.T) {
	msg:=NewMsgInstantTxVote()
	var buf bytes.Buffer
	err := msg.BtcEncode(&buf, ProtocolVersion)
	if err != nil {
		t.Errorf("BtcEncode error %v", err)
	}

	var msg2 MsgInstantTxVote
	rbuf := bytes.NewReader(buf.Bytes())
	err = msg2.BtcDecode(rbuf, ProtocolVersion)
	if err != nil {
		t.Errorf("BtcDecode error %v",err)
	}
}