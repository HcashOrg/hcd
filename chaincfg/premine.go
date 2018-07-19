// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hcd developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

// BlockOneLedgerMainNet is the block one output ledger for the main
// network.
var BlockOneLedgerMainNet = []*TokenPayout{
	{"HsJwqHsbVjb4SdgdYQBDxcWQCsyJxiJMqda", 4000000 * 1e8},
	{"HsH2U86zVmqT6p17fsfxhx7CP5w73XhHcs3", 4000000 * 1e8},
	{"HsLhsSN5ziYgK9XQosx8VjM59wEZjXAuzZC", 4000000 * 1e8},
	{"HsPo3AJrYgtvHwpqQKNCEkrTPrLcY1irMqe", 4000000 * 1e8},
	{"HsKrJGTiQKERFQmTu3UPN8c8toUgquMiY2i", 4000000 * 1e8},
	{"HsJvzUrW3fcDbxCskspHYCj6zbRS6aDcDw1", 4000000 * 1e8},
	{"HsFf124jsfwNaUwbrDU22BcuKJjim7nF6Y3", 4000000 * 1e8},
	{"HsMJiuGD5afsr37yKrHUtkXoZTCCoj5sPMb", 4000000 * 1e8},
	{"HsGXZUDR3L35TXH7UNXMhCTuRpNVHF6BsNr", 4000000 * 1e8},
	{"HsNRWXPC7bRrfCakNNXF3izJoyQXTTNg8ik", 4000000 * 1e8},
	{"HsLBVmVQ3RhdrqFJWJoFY8qEH31iKwtuapL", 4478231 * 1e8},
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
	{"Sshw6S86G2bV6W32cbc7EhtFy8f93rU6pae", 100000 * 1e8},
	{"SsjXRK6Xz6CFuBt6PugBvrkdAa4xGbcZ18w", 100000 * 1e8},
	{"SsfXiYkYkCoo31CuVQw428N6wWKus2ZEw5X", 100000 * 1e8},
}
