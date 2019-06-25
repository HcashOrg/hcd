package mempool

import (
	"github.com/HcashOrg/hcd/chaincfg"
	"testing"
)

func TestTxLockPool(t *testing.T) {
	t.Parallel()
	var txLen = 10
	harness, spendableOuts, err := newLockPoolHarness(&chaincfg.MainNetParams)
	for _, v := range spendableOuts {
		t.Log(v.outPoint.String())
	}

	if err != nil {
		t.Fatalf("unable to create test pool: %v", err)
	}

	// Create a chain of transactions rooted with the first spendable output
	// provided by the harness.
	chainedTxns, err := harness.CreateLockTxChain(spendableOuts[0], uint32(txLen))
	if err != nil {
		t.Fatalf("unable to create transaction chain: %v", err)
	}

	// Ensure orphans are rejected when the allow orphans flag is not set.
	for _, tx := range chainedTxns[:] {

		harness.txPool.MayBeAddToLockPool(tx, 888, true,
			false, false)
	}

	if len(harness.txPool.txLockPool) != txLen {
		t.Fatalf("maybeAddtoLockPool fail,lockpool len:%v",len(harness.txPool.txLockPool))
	}
	harness.chain.currentHeight=45888
	t.Log(harness.txPool.FetchPendingLockTx(1))

	t.Log(harness.txPool.TxLockPoolInfo())

	for _, tx := range chainedTxns[:] {
		harness.txPool.ModifyLockTransaction(tx, 45668)
	}

	for _, desc := range harness.txPool.txLockPool {
		if len(harness.txPool.txLockPool) != txLen || desc.MineHeight != 45668 {
			t.Fatalf("ModifyLockTransaction 45668 err")
		}
	}

	t.Log(harness.txPool.TxLockPoolInfo())
	for _, tx := range chainedTxns[:] {
		harness.txPool.ModifyLockTransaction(tx, 0)

	}
	for _, desc := range harness.txPool.txLockPool {
		if len(harness.txPool.txLockPool) != txLen || desc.MineHeight != 0 {
			t.Fatalf("ModifyLockTransaction 0 err")
		}
	}


	t.Log(harness.txPool.TxLockPoolInfo())

	for _, tx := range chainedTxns[:] {
		harness.txPool.ModifyLockTransaction(tx, 45668)
	}
	for _, desc := range harness.txPool.txLockPool {
		if len(harness.txPool.txLockPool) != txLen || desc.MineHeight != 45668 {
			t.Fatalf("ModifyLockTransaction 45668 err")
		}
	}

	t.Log(harness.txPool.TxLockPoolInfo())

	harness.txPool.RemoveConfirmedLockTransaction(45768)

	if len(harness.txPool.txLockPool) != 0||len(harness.txPool.lockOutpoints)!=0 {
		t.Fatalf("RemoveConfirmedLockTransaction err")
	}

	//t.Log(harness.txPool.TxLockPoolInfo())

	for _, tx := range chainedTxns[:] {
		//t.Log(tx.MsgTx().TxIn[0].PreviousOutPoint.String())
		harness.txPool.MayBeAddToLockPool(tx,888, true,
			false, false)
	}


	if len(harness.txPool.txLockPool) != txLen {
		t.Fatalf("maybeAddtoLockPool err")
	}
	t.Log(harness.txPool.TxLockPoolInfo())

	for _, tx := range chainedTxns[:] {

		chainedTxns2, _ := harness.CreateTxChain(spendableOutput{tx.MsgTx().TxIn[0].PreviousOutPoint, 0}, 1)

		harness.txPool.RemoveTxLockDoubleSpends(chainedTxns2[0])
		//t.Log(harness.txPool.TxLockPoolInfo())
	}
	if len(harness.txPool.txLockPool) != 0 ||len(harness.txPool.lockOutpoints)!=0{
		t.Fatalf("RemoveTxLockDoubleSpends err")
	}

	for _, tx := range chainedTxns[:] {
		//t.Log(tx.MsgTx().TxIn[0].PreviousOutPoint.String())
		harness.txPool.MayBeAddToLockPool(tx,888, true,
			false, false)
	}

	if len(harness.txPool.txLockPool) != txLen {
		t.Fatalf("maybeAddtoLockPool err")
	}
	t.Log(harness.txPool.TxLockPoolInfo())

	for _, tx := range chainedTxns[:txLen/2] {
		harness.txPool.ModifyLockTransaction(tx, 45668)
	}


	t.Log(harness.txPool.TxLockPoolInfo())
	for _, tx := range chainedTxns[:] {

		chainedTxns2, _ := harness.CreateTxChain(spendableOutput{tx.MsgTx().TxIn[0].PreviousOutPoint, 0}, 1)

		harness.txPool.RemoveTxLockDoubleSpends(chainedTxns2[0])
		//t.Log(harness.txPool.TxLockPoolInfo())
	}

	if len(harness.txPool.txLockPool) != 0 || len(harness.txPool.lockOutpoints) != 0 {
		t.Fatalf("RemoveTxLockDoubleSpends err")
	}

	t.Log(harness.txPool.TxLockPoolInfo())
}

