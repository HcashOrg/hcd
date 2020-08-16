package rpcserver

import (
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcd/internal/mempool"
)

type TxMempooler interface {
	// HaveTransactions returns whether or not the passed transactions
	// already exist in the main pool or in the orphan pool.
	HaveTransactions(hashes []*chainhash.Hash) []bool

	// TxDescs returns a slice of descriptors for all the transactions in
	// the pool. The descriptors must be treated as read only.
	TxDescs() []*mempool.TxDesc


	// Count returns the number of transactions in the main pool. It does
	// not include the orphan pool.
	Count() int

	// FetchTransaction returns the requested transaction from the
	// transaction pool. This only fetches from the main transaction pool
	// and does not include orphans.
	FetchTransaction(txHash *chainhash.Hash) (*hcutil.Tx, error)
}
