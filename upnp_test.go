package main

import (
	"testing"
)

func TestDiscover(t *testing.T) {
	nat,err:=Discover()
	if err!=nil{
		t.Fatal(err)
	}
	t.Log(nat)

	listenPort,err:=nat.AddPortMapping("tcp", 6666, 8888,
		"hcd listen port", 20*60)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("listenPort:",listenPort)
}
