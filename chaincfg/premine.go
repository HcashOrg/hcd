// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

// BlockOneLedgerMainNet is the block one output ledger for the main
// network.
var BlockOneLedgerMainNet = []*TokenPayout{
	{"HsW2kpP9EfCpHyRnhRDNQxtzpf9ZJFKohyv",31326010 * 1e8},
	{"HsDJjXwWB8DDZAujwsQszHM5RxxoGpTBxrk", 2000000 * 1e8},
	{"HsHd1ZyBLenU677NgxweTrNqVPoqcFkXq5t", 2000000 * 1e8},
	{"HsKmFZG7tPwysiN2G4UoRjgjsNNoGfLK3Ef", 2000000 * 1e8},
	{"HsZuN7E8Pk9ZN7QAbZkrmXbjeNx73QSnq1e", 2000000 * 1e8},
	{"HsUFPNENKsevjoqg28Hkj6R8pUkrRNCNJpq", 2000000 * 1e8},
	{"HsGL9BprKZKR2X4HuuSHZYjf5PDiHJtGeZc", 2000000 * 1e8},
	{"HsBxxtGn3Ggyx5eKxCMQ6ggtFrHRwWMA4kh", 2000000 * 1e8},
	{"HsPf1VZy2JDSA8WMFyUTdSnXaZctAoeqsMq", 2000000 * 1e8},
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
var AIEnableHeightMainNet = uint64(206666)
var AIEnableHeightTestNet = uint64(366560)
var AIEnableHeightSimNet = uint64(186)
