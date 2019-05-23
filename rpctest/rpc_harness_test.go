// Copyright (c) 2016 The btcsuite developers
// Copyright (c) 2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpctest

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcd/txscript"
	"github.com/HcashOrg/hcd/wire"
	"github.com/HcashOrg/hcd/hcutil"
)

const (
	numMatureOutputs = 25
)

func testSendOutputs(r *Harness, t *testing.T) {
	genSpend := func(amt hcutil.Amount) *chainhash.Hash {
		// Grab a fresh address from the wallet.
		addr, err := r.NewAddress()
		if err != nil {
			t.Fatalf("unable to get new address: %v", err)
		}

		// Next, send amt to this address, spending from one of our
		// mature coinbase outputs.
		addrScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			t.Fatalf("unable to generate pkscript to addr: %v", err)
		}
		output := wire.NewTxOut(int64(amt), addrScript)
		txid, err := r.SendOutputs([]*wire.TxOut{output}, 10)
		if err != nil {
			t.Fatalf("coinbase spend failed: %v", err)
		}
		return txid
	}

	assertTxMined := func(txid *chainhash.Hash, blockHash *chainhash.Hash) {
		block, err := r.Node.GetBlock(blockHash)
		if err != nil {
			t.Fatalf("unable to get block: %v", err)
		}

		numBlockTxns := len(block.Transactions)
		if numBlockTxns < 2 {
			t.Fatalf("crafted transaction wasn't mined, block should have "+
				"at least %v transactions instead has %v", 2, numBlockTxns)
		}

		minedTx := block.Transactions[1]
		txHash := minedTx.TxHash()
		if txHash != *txid {
			t.Fatalf("txid's don't match, %v vs %v", txHash, txid)
		}
	}

	// First, generate a small spend which will require only a single
	// input.
	txid := genSpend(hcutil.Amount(5 * hcutil.AtomsPerCoin))

	// Generate a single block, the transaction the wallet created should
	// be found in this block.
	blockHashes, err := r.Node.Generate(1)
	if err != nil {
		t.Fatalf("unable to generate single block: %v", err)
	}
	assertTxMined(txid, blockHashes[0])

	// Next, generate a spend much greater than the block reward. This
	// transaction should also have been mined properly.
	txid = genSpend(hcutil.Amount(5000 * hcutil.AtomsPerCoin))
	blockHashes, err = r.Node.Generate(1)
	if err != nil {
		t.Fatalf("unable to generate single block: %v", err)
	}
	assertTxMined(txid, blockHashes[0])

	// Generate another block to ensure the transaction is removed from the
	// mempool.
	if _, err := r.Node.Generate(1); err != nil {
		t.Fatalf("unable to generate block: %v", err)
	}
}

func assertConnectedTo(t *testing.T, nodeA *Harness, nodeB *Harness) {
	nodeAPeers, err := nodeA.Node.GetPeerInfo()
	if err != nil {
		t.Fatalf("unable to get nodeA's peer info")
	}

	nodeAddr := nodeB.node.config.listen
	addrFound := false
	for _, peerInfo := range nodeAPeers {
		if peerInfo.Addr == nodeAddr {
			addrFound = true
			break
		}
	}

	if !addrFound {
		t.Fatal("nodeA not connected to nodeB")
	}
}

func testConnectNode(r *Harness, t *testing.T) {
	// Create a fresh test harness.
	harness, err := New(&chaincfg.SimNetParams, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := harness.SetUp(false, 0); err != nil {
		t.Fatalf("unable to complete rpctest setup: %v", err)
	}
	defer harness.TearDown()

	// Establish a p2p connection from our new local harness to the main
	// harness.
	if err := ConnectNode(harness, r); err != nil {
		t.Fatalf("unable to connect local to main harness: %v", err)
	}

	// The main harness should show up in our local harness' peer's list,
	// and vice verse.
	assertConnectedTo(t, harness, r)
}

func testTearDownAll(t *testing.T) {
	// Grab a local copy of the currently active harnesses before
	// attempting to tear them all down.
	initialActiveHarnesses := ActiveHarnesses()

	// Tear down all currently active harnesses.
	if err := TearDownAll(); err != nil {
		t.Fatalf("unable to teardown all harnesses: %v", err)
	}

	// The global testInstances map should now be fully purged with no
	// active test harnesses remaining.
	if len(ActiveHarnesses()) != 0 {
		t.Fatalf("test harnesses still active after TearDownAll")
	}
	for _, harness := range initialActiveHarnesses {
		// Ensure all test directories have been deleted.
		if _, err := os.Stat(harness.testNodeDir); err == nil {
			t.Errorf("created test datadir was not deleted.")
		}
	}

}

func testActiveHarnesses(r *Harness, t *testing.T) {
	numInitialHarnesses := len(ActiveHarnesses())

	// Create a single test harness.
	harness1, err := New(&chaincfg.SimNetParams, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer harness1.TearDown()

	// With the harness created above, a single harness should be detected
	// as active.
	numActiveHarnesses := len(ActiveHarnesses())
	if !(numActiveHarnesses > numInitialHarnesses) {
		t.Fatalf("ActiveHarnesses not updated, should have an " +
			"additional test harness listed.")
	}
}

func testJoinMempools(r *Harness, t *testing.T) {
	// Assert main test harness has no transactions in its mempool.
	pooledHashes, err := r.Node.GetRawMempool(hcjson.GRMAll)
	if err != nil {
		t.Fatalf("unable to get mempool for main test harness: %v", err)
	}
	if len(pooledHashes) != 0 {
		t.Fatal("main test harness mempool not empty")
	}

	// Create a local test harness with only the genesis block.  The nodes
	// will be synced below so the same transaction can be sent to both
	// nodes without it being an orphan.
	harness, err := New(&chaincfg.SimNetParams, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := harness.SetUp(false, 0); err != nil {
		t.Fatalf("unable to complete rpctest setup: %v", err)
	}
	defer harness.TearDown()

	nodeSlice := []*Harness{r, harness}

	// Both mempools should be considered synced as they are empty.
	// Therefore, this should return instantly.
	if err := JoinNodes(nodeSlice, Mempools); err != nil {
		t.Fatalf("unable to join node on mempools: %v", err)
	}

	// Generate a coinbase spend to a new address within the main harness'
	// mempool.
	addr, err := r.NewAddress()
	addrScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		t.Fatalf("unable to generate pkscript to addr: %v", err)
	}
	output := wire.NewTxOut(5e8, addrScript)
	testTx, err := r.CreateTransaction([]*wire.TxOut{output}, 10)
	if err != nil {
		t.Fatalf("coinbase spend failed: %v", err)
	}
	if _, err := r.Node.SendRawTransaction(testTx, true); err != nil {
		t.Fatalf("send transaction failed: %v", err)
	}

	// Wait until the transaction shows up to ensure the two mempools are
	// not the same.
	harnessSynced := make(chan struct{})
	go func() {
		for {
			poolHashes, err := r.Node.GetRawMempool(hcjson.GRMAll)
			if err != nil {
				t.Fatalf("failed to retrieve harness mempool: %v", err)
			}
			if len(poolHashes) > 0 {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		harnessSynced <- struct{}{}
	}()
	select {
	case <-harnessSynced:
	case <-time.After(time.Minute):
		t.Fatalf("harness node never received transaction")
	}

	// This select case should fall through to the default as the goroutine
	// should be blocked on the JoinNodes call.
	poolsSynced := make(chan struct{})
	go func() {
		if err := JoinNodes(nodeSlice, Mempools); err != nil {
			t.Fatalf("unable to join node on mempools: %v", err)
		}
		poolsSynced <- struct{}{}
	}()
	select {
	case <-poolsSynced:
		t.Fatalf("mempools detected as synced yet harness has a new tx")
	default:
	}

	// Establish an outbound connection from the local harness to the main
	// harness and wait for the chains to be synced.
	if err := ConnectNode(harness, r); err != nil {
		t.Fatalf("unable to connect harnesses: %v", err)
	}
	if err := JoinNodes(nodeSlice, Blocks); err != nil {
		t.Fatalf("unable to join node on blocks: %v", err)
	}

	// Send the transaction to the local harness which will result in synced
	// mempools.
	if _, err := harness.Node.SendRawTransaction(testTx, true); err != nil {
		t.Fatalf("send transaction failed: %v", err)
	}

	// Select once again with a special timeout case after 1 minute. The
	// goroutine above should now be blocked on sending into the unbuffered
	// channel. The send should immediately succeed. In order to avoid the
	// test hanging indefinitely, a 1 minute timeout is in place.
	select {
	case <-poolsSynced:
		// fall through
	case <-time.After(time.Minute):
		t.Fatalf("mempools never detected as synced")
	}
}

func testMemWalletReorg(r *Harness, t *testing.T) {
	// Create a fresh harness, we'll be using the main harness to force a
	// re-org on this local harness.
	harness, err := New(&chaincfg.SimNetParams, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := harness.SetUp(true, 5); err != nil {
		t.Fatalf("unable to complete rpctest setup: %v", err)
	}
	defer harness.TearDown()

	// Ensure the internal wallet has the expected balance.
	expectedBalance := hcutil.Amount(5 * 300 * hcutil.AtomsPerCoin)
	walletBalance := harness.ConfirmedBalance()
	if expectedBalance != walletBalance {
		t.Fatalf("wallet balance incorrect: expected %v, got %v",
			expectedBalance, walletBalance)
	}

	// Now connect this local harness to the main harness then wait for
	// their chains to synchronize.
	if err := ConnectNode(harness, r); err != nil {
		t.Fatalf("unable to connect harnesses: %v", err)
	}
	nodeSlice := []*Harness{r, harness}
	if err := JoinNodes(nodeSlice, Blocks); err != nil {
		t.Fatalf("unable to join node on blocks: %v", err)
	}

	// The original wallet should now have a balance of 0 Coin as its entire
	// chain should have been decimated in favor of the main harness'
	// chain.
	expectedBalance = hcutil.Amount(0)
	walletBalance = harness.ConfirmedBalance()
	if expectedBalance != walletBalance {
		t.Fatalf("wallet balance incorrect: expected %v, got %v",
			expectedBalance, walletBalance)
	}
}

// testJoinBlocks tests new block join blockchain
func testJoinBlocks(r *Harness, t *testing.T) {
	// Create a second harness with only the genesis block so it is behind
	// the main harness.
	harness, err := New(&chaincfg.SimNetParams, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := harness.SetUp(false, 0); err != nil {
		t.Fatalf("unable to complete rpctest setup: %v", err)
	}
	defer harness.TearDown()

	nodeSlice := []*Harness{r, harness}
	blocksSynced := make(chan struct{})
	go func() {
		if err := JoinNodes(nodeSlice, Blocks); err != nil {
			t.Fatalf("unable to join node on blocks: %v", err)
		}
		blocksSynced <- struct{}{}
	}()

	// This select case should fall through to the default as the goroutine
	// should be blocked on the JoinNodes calls.
	select {
	case <-blocksSynced:
		t.Fatalf("blocks detected as synced yet local harness is behind")
	default:
	}

	// Connect the local harness to the main harness which will sync the
	// chains.
	if err := ConnectNode(harness, r); err != nil {
		t.Fatalf("unable to connect harnesses: %v", err)
	}

	// Select once again with a special timeout case after 1 minute. The
	// goroutine above should now be blocked on sending into the unbuffered
	// channel. The send should immediately succeed. In order to avoid the
	// test hanging indefinitely, a 1 minute timeout is in place.
	select {
	case <-blocksSynced:
		// fall through
	case <-time.After(time.Minute):
		t.Fatalf("blocks never detected as synced")
	}
}

func testMemWalletLockedOutputs(r *Harness, t *testing.T) {
	// Obtain the initial balance of the wallet at this point.
	startingBalance := r.ConfirmedBalance()
	
	fmt.Println("startingBalance:",startingBalance)

	// First, create a signed transaction spending some outputs.
	addr, err := r.NewAddress()
	if err != nil {
		t.Fatalf("unable to generate new address: %v", err)
	}
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		t.Fatalf("unable to create script: %v", err)
	}
	outputAmt := hcutil.Amount(50 * hcutil.AtomsPerCoin)
	output := wire.NewTxOut(int64(outputAmt), pkScript)
	tx, err := r.CreateTransaction([]*wire.TxOut{output}, 10)
	if err != nil {
		t.Fatalf("unable to create transaction: %v", err)
	}

	// The current wallet balance should now be at least 50 Coin less
	// (accounting for fees) than the period balance
	currentBalance := r.ConfirmedBalance()
	if !(currentBalance <= startingBalance-outputAmt) {
		t.Fatalf("spent outputs not locked: previous balance %v, "+
			"current balance %v", startingBalance, currentBalance)
	}

	// Now unlocked all the spent inputs within the unbroadcast signed
	// transaction. The current balance should now be exactly that of the
	// starting balance.
	r.UnlockOutputs(tx.TxIn)
	currentBalance = r.ConfirmedBalance()
	if currentBalance != startingBalance {
		t.Fatalf("current and starting balance should now match: "+
			"expected %v, got %v", startingBalance, currentBalance)
	}
}

var harnessTestCases = []HarnessTestCase{
	testSendOutputs,
	testConnectNode,
	testActiveHarnesses,
	testJoinBlocks,
	testJoinMempools, // Depends on results of testJoinBlocks
	testMemWalletReorg,
	testMemWalletLockedOutputs,
}

var mainHarness *Harness

func TestMain(m *testing.M) {
	var err error
	mainHarness, err = New(&chaincfg.SimNetParams, nil, nil)
	if err != nil {
		fmt.Println("unable to create main harness: ", err)
		os.Exit(1)
	}

	// Initialize the main mining node with a chain of length 42, providing
	// 25 mature coinbases to allow spending from for testing purposes.
	if err = mainHarness.SetUp(true, numMatureOutputs); err != nil {
		fmt.Println("unable to setup test chain: ", err)

		// Even though the harness was not fully setup, it still needs
		// to be torn down to ensure all resources such as temp
		// directories are cleaned up.  The error is intentionally
		// ignored since this is already an error path and nothing else
		// could be done about it anyways.
		_ = mainHarness.TearDown()
		os.Exit(1)
	}

	exitCode := m.Run()
	fmt.Println("exitCode:",exitCode)

	// Clean up any active harnesses that are still currently running.
	if len(ActiveHarnesses()) > 0 {
		if err := TearDownAll(); err != nil {
			fmt.Println("unable to tear down chain: ", err)
			os.Exit(1)
		}
	}

	os.Exit(exitCode)
}


func assertConnectedOut(t *testing.T, nodeA *Harness, nodeB *Harness) {
	nodeAPeers, err := nodeA.Node.GetPeerInfo()
	if err != nil {
		t.Fatalf("unable to get node's peer info")
	}

	nodeAddr := nodeB.node.config.listen
	addrFound := true
	for _, peerInfo := range nodeAPeers {
		if peerInfo.Addr == nodeAddr {
			addrFound = false
			break
		}
	}

	if addrFound {
		t.Fatal("nodeA not connected to node")
	}
}

func TestHarness(t *testing.T) {
	// We should have the expected amount of mature unspent outputs.
	expectedBalance := hcutil.Amount(numMatureOutputs * 300 * hcutil.AtomsPerCoin)
	fmt.Println("expectedBalance:",expectedBalance)
	harnessBalance := mainHarness.ConfirmedBalance()
	if harnessBalance != expectedBalance {
		t.Fatalf("expected wallet balance of %v instead have %v",
			expectedBalance, harnessBalance)
	}

	// Current tip should be at a height of numMatureOutputs plus the
	// required number of blocks for coinbase maturity plus an additional
	// block for the premine block.
	nodeInfo, err := mainHarness.Node.GetInfo()
	if err != nil {
		t.Fatalf("unable to execute getinfo on node: %v", err)
	}
	coinbaseMaturity := uint32(mainHarness.ActiveNet.CoinbaseMaturity)
	expectedChainHeight := numMatureOutputs + coinbaseMaturity + 1
	if uint32(nodeInfo.Blocks) != expectedChainHeight {
		t.Errorf("Chain height is %v, should be %v",
			nodeInfo.Blocks, expectedChainHeight)
	}

	for _, testCase := range harnessTestCases {
		testCase(mainHarness, t)
	}

	testTearDownAll(t)
}


func testNodesConnected(r *Harness, t *testing.T) {
	// Create a fresh test harness.
	harness, err := New(&chaincfg.SimNetParams, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := harness.SetUp(false, 0); err != nil {
		t.Fatalf("unable to complete rpctest setup: %v", err)
	}
	defer harness.TearDown()

	// Establish a p2p connection from our new local harness to the main
	// harness.
	if err := ConnectNode(harness, r); err != nil {
		t.Fatalf("unable to connect local to main harness: %v", err)
	}

	// Sanity check.
	assertConnectedTo(t, harness, r)

	// Ensure nodes are still connected.
	assertConnectedTo(t, harness, r)

	testCases := []struct {
		name         string
		allowReverse bool
		expected     bool
		from         *Harness
		to           *Harness
	}{
		// The existing connection is h->r.
		{"!allowReverse, h->r", false, true, harness, r},
		{"allowReverse, h->r", true, true, harness, r},
		{"!allowReverse, r->h", false, false, r, harness},
		{"allowReverse, r->h", true, true, r, harness},
	}

	for _, tc := range testCases {
		actual, err := NodesConnected(tc.from, tc.to, tc.allowReverse)
		if err != nil {
			t.Fatalf("unable to determine node connection: %v", err)
		}
		if actual != tc.expected {
			t.Fatalf("test case %s: actual result (%v) differs from expected "+
				"(%v)", tc.name, actual, tc.expected)
		}
	}

	// Disconnect the nodes.
	if err := RemoveNode(harness, r); err != nil {
		t.Fatalf("unable to disconnect local to main harness: %v", err)
	}
	// All test cases must return false now.
	for _, tc := range testCases {
		actual, err := NodesConnected(tc.from, tc.to, tc.allowReverse)
		if err != nil {
			t.Fatalf("unable to determine node connection: %v", err)
		}
		if actual {
			t.Fatalf("test case %s: nodes connected after commanded to "+
				"disconnect", tc.name)
		}
	}

}

