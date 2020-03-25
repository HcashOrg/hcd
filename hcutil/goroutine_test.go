package hcutil

import "testing"

func TestCurGoroutineID(t *testing.T) {
	id:=CurGoroutineID()
	t.Log(id)
}
