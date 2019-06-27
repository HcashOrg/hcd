// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

// BlockOneLedgerMainNet is the block one output ledger for the main
// network.
var BlockOneLedgerMainNet = []*TokenPayout{
	{"HsXr6yF1EBx4KPpNagXCGXavzZB4jUHvHF9",41326010 * 1e8},
	{"HsVcu6C3Xmb4iNin3uB5cGXkHZmzHDxUbXn", 2000000 * 1e8},
}

// BlockOneLedgerTestNet is the block one output ledger for the test
// network.
var BlockOneLedgerTestNet = []*TokenPayout{
	{"TspAtg3jUAieeHhR9wgWso5CET3tZxQmvLq", 40000000 * 1e8},
}
// BlockOneLedgerTestNet2 is the block one output ledger for the 2nd test
// network.
var BlockOneLedgerTestNet2 = []*TokenPayout{
	{"TspAtg3jUAieeHhR9wgWso5CET3tZxQmvLq", 40000000 * 1e8},
}

// BlockOneLedgerSimNet is the block one output ledger for the simulation
// network. See under "Hcd organization related parameters" in params.go
// for information on how to spend these outputs.
var BlockOneLedgerSimNet = []*TokenPayout{
	{"SsbQiVBFLP79dC7gZT5PWH7gEfMhUkExPTv", 100000 * 1e8},
	{"Ssf3CjCEU28p8BNm2k53DqC2FZDJ6PdEHe8", 100000 * 1e8},
	{"SsW6wQWN8wiVY4eqQY3sQi4Eu6kjTqmCH1q", 100000 * 1e8},
}

var AIEnableHeightSimNet = uint64(186)
