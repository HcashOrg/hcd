package wire

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestRouteAddr(t *testing.T) {

	// Ensure the command is expected value.
	wantCmd := "routeaddr"
	msg := NewMsgRouteAddr()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgRouteAddr: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	na:="HsTJckn6hjhP4QYHF7CE87ok3y5TDA2gd6D"
	err := msg.AddAddress(na)
	if err != nil {
		t.Errorf("AddAddress: %v", err)
	}

	if msg.AddrList[0] != na {
		t.Errorf("AddAddress: wrong address added - got %v, want %v",
			spew.Sprint(msg.AddrList[0]), spew.Sprint(na))
	}

	// Ensure the address list is cleared properly.
	msg.ClearAddresses()
	if len(msg.AddrList) != 0 {
		t.Errorf("ClearAddresses: address list is not empty - "+
			"got %v [%v], want %v", len(msg.AddrList),
			spew.Sprint(msg.AddrList[0]), 0)
	}

	// Ensure adding more than the max allowed addresses per message returns
	// error.
	for i := 0; i < MaxAddrPerMsg+1; i++ {
		err = msg.AddAddress(na)
	}
	if err == nil {
		t.Errorf("AddAddress: expected error on too many addresses " +
			"not received")
	}
	err = msg.AddAddresses(na)
	if err == nil {
		t.Errorf("AddAddresses: expected error on too many addresses " +
			"not received")
	}

	return
}

