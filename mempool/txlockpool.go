package mempool

import (
	"errors"
	"fmt"
	"github.com/HcashOrg/hcd/blockchain"
	"github.com/HcashOrg/hcd/blockchain/stake"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcd/txscript"
	"github.com/HcashOrg/hcd/wire"
	"math"
	"time"
)

const (
	defaultConfirmNum = 6
	defaultBehindNums = 10
)

type InstantTxDesc struct {
	Tx *hcutil.InstantTx
	// Height is the block height when the entry was added to the source
	// pool.
	AddHeight int64
	Votes     []*hcutil.InstantTxVote
	Send      bool

	MineHeight int64 //
}

type lockPool struct {
	txLockPool     map[chainhash.Hash]*InstantTxDesc //for instantsend lock tx pool
	lockOutpoints  map[wire.OutPoint]*hcutil.InstantTx
	instantTxVotes map[chainhash.Hash]*hcutil.InstantTxVote
}

//update inistant tx state according the mined height
func (mp *TxPool) modifyInstantTxHeight(tx *hcutil.Tx, height int64) {
	if desc, exist := mp.txLockPool[*tx.Hash()]; exist {
		desc.MineHeight = height
	}
}

func (mp *TxPool) AppendInstantTxVote(hash *chainhash.Hash, vote *hcutil.InstantTxVote) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	mp.appendInstantTxVote(hash, vote)
}

func (mp *TxPool) appendInstantTxVote(hash *chainhash.Hash, vote *hcutil.InstantTxVote) {
	if desc, exist := mp.txLockPool[*hash]; exist && vote != nil {
		desc.Votes = append(desc.Votes, vote)

		mp.instantTxVotes[*vote.Hash()] = vote
	}
}

func (mp *TxPool) GetInstantTxDesc(hash *chainhash.Hash) (desc *InstantTxDesc, exist bool) {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()

	return mp.getInstantTxDesc(hash)
}

func (mp *TxPool) getInstantTxDesc(hash *chainhash.Hash) (desc *InstantTxDesc, exist bool) {
	desc, exist = mp.txLockPool[*hash]
	return
}

func (mp *TxPool) ModifyInstantTxHeight(tx *hcutil.Tx, height int64) {
	// Protect concurrent access.
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	mp.modifyInstantTxHeight(tx, height)
}

func (mp *TxPool) RemoveConfirmedInstantTx(height int64) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	for hash, desc := range mp.txLockPool {
		if desc.MineHeight != 0 && desc.MineHeight < height-defaultConfirmNum {
			//remove vote index
			for _, vote := range desc.Votes {
				delete(mp.instantTxVotes, *vote.Hash())
			}

			//remove instantTx
			delete(mp.txLockPool, hash)

			//remove tx output index
			for _, txIn := range desc.Tx.MsgTx().TxIn {
				delete(mp.lockOutpoints, txIn.PreviousOutPoint)
			}
		}
	}
}

func (mp *TxPool) IsInstantTxExist(hash *chainhash.Hash) bool {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()
	return mp.isInstantTxExist(hash)
}

//Is tx in  locked?
func (mp *TxPool) isInstantTxExist(hash *chainhash.Hash) bool {
	if _, exists := mp.txLockPool[*hash]; exists {
		return true
	}
	return false
}

//Is txVin  in locked?
func (mp *TxPool) isInstantTxInputExist(outPoint *wire.OutPoint) (*hcutil.InstantTx, bool) {
	if txLock, exists := mp.lockOutpoints[*outPoint]; exists {
		return txLock, true
	}
	return nil, false
}

func (mp *TxPool) TxLockPoolInfo() map[string]*hcjson.TxLockInfo {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()

	ret := make(map[string]*hcjson.TxLockInfo, len(mp.txLockPool))

	for hash, desc := range mp.txLockPool {
		votesHash := make([]string, 0, 5)
		for _, vote := range desc.Votes {
			votesHash = append(votesHash, vote.Hash().String()+"-"+vote.MsgInstantTxVote().TicketHash.String())
		}

		ret[hash.String()] = &hcjson.TxLockInfo{AddHeight: desc.AddHeight, MineHeight: desc.MineHeight, Votes: votesHash, Send: desc.Send}
	}

	return ret
}

func (mp *TxPool) FetchPendingLockTx(behindNums int64) [][]byte {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()

	if behindNums <= 0 {
		behindNums = defaultBehindNums
	}
	bestHeight := mp.cfg.BestHeight()
	minHeight := bestHeight - behindNums

	retMsgTx := make([][]byte, 0)
	for hash, desc := range mp.txLockPool {
		if desc.MineHeight == 0 && desc.AddHeight < minHeight {
			if desc.Send {//voted but not be mine,it will be resend by wallet
				bts, err := desc.Tx.MsgTx().Bytes()
				if err == nil {
					retMsgTx = append(retMsgTx, bts)
				}
			}else{// remove from txlockpool,because havn`t be voted for a long time

				//remove vote index
				for _, vote := range desc.Votes {
					delete(mp.instantTxVotes, *vote.Hash())
				}

				//remove instantTx
				delete(mp.txLockPool, hash)

				//remove tx output index
				for _, txIn := range desc.Tx.MsgTx().TxIn {
					delete(mp.lockOutpoints, txIn.PreviousOutPoint)
				}

			}
		}

	}

	return retMsgTx
}

//check block transactions is conflict with lockPool
func (mp *TxPool) CheckBlkConflictWithTxLockPool(block *hcutil.Block) (bool, error) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	for _, tx := range block.Transactions() {
		err := mp.checkTxWithLockPool(tx)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

//check the input double spent
func (mp *TxPool) checkTxWithLockPool(tx *hcutil.Tx) error {
	if !mp.isInstantTxExist(tx.Hash()) {
		for _, txIn := range tx.MsgTx().TxIn {
			if _, exist := mp.isInstantTxInputExist(&txIn.PreviousOutPoint); exist {
				return fmt.Errorf("tx %v conflict with lock pool", tx.Hash())
			}
		}
	}
	return nil
}

//remove txlock which is conflict with tx
func (mp *TxPool) RemoveInstantTxDoubleSpends(tx *hcutil.Tx) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	//if is the same tx ,just return
	if mp.isInstantTxExist(tx.Hash()) {
		return
	}

	//if tx in is conflict with txlock ,just remove txlock and lockOutpoint
	for _, invalue := range tx.MsgTx().TxIn {
		if txLock, exist := mp.isInstantTxInputExist(&invalue.PreviousOutPoint); exist {
			instantTxdesc := mp.txLockPool[*txLock.Hash()]

			for _, vote := range instantTxdesc.Votes {
				delete(mp.instantTxVotes, *vote.Hash())
			}

			delete(mp.txLockPool, *txLock.Hash())

			for _, txIn := range txLock.MsgTx().TxIn {
				delete(mp.lockOutpoints, txIn.PreviousOutPoint)
			}

		}
	}

}

func (mp *TxPool) MayBeAddToLockPool(tx *hcutil.InstantTx, isNew, rateLimit, allowHighFees bool) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	mp.maybeAddtoLockPool(tx, isNew, rateLimit, allowHighFees)
}

//this is called before inserting to mempool,must be called with lock
func (mp *TxPool) maybeAddtoLockPool(instantTx *hcutil.InstantTx, isNew, rateLimit, allowHighFees bool) {
	//if exist just return ,or will rewrite the state of this txlock
	if mp.isInstantTxExist(instantTx.Hash()) {
		return
	}
	//check with lockpool
	tx := instantTx.Tx
	err := mp.checkTxWithLockPool(&tx)
	if err != nil {
		log.Tracef("instant Transaction %v is conflict with lockpool : %v", instantTx.Hash(),
			err)
		return
	}
	//check with mempool
	_, err = mp.checkInstantTxWithMem(instantTx, isNew, rateLimit, allowHighFees)
	if err != nil {
		log.Tracef("instant Transaction %v is conflict with mempool : %v", instantTx.Hash(),
			err)
		return
	}

	//check instant tag
	msgTx := instantTx.MsgTx()
	_, isInstantTx := txscript.IsInstantTx(msgTx)
	if !isInstantTx {
		log.Tracef("Transaction %v is not instant instantTx ", instantTx.Hash())
		return
	}
	bestHeight := mp.cfg.BestHeight()
	mp.txLockPool[*instantTx.Hash()] = &InstantTxDesc{
		Tx:         instantTx,
		AddHeight:  bestHeight,
		MineHeight: 0,
		Send:       false,
		Votes:      make([]*hcutil.InstantTxVote, 0, 5)}

	for _, txIn := range msgTx.TxIn {
		mp.lockOutpoints[txIn.PreviousOutPoint] = instantTx
	}
}

func (mp *TxPool) checkInstantTxWithMem(instantTx *hcutil.InstantTx, isNew, rateLimit, allowHighFees bool) ([]*chainhash.Hash, error) {
	tx := &instantTx.Tx
	msgTx := tx.MsgTx()
	txHash := tx.Hash()
	// Don't accept the transaction if it already exists in the pool.  This
	// applies to orphan transactions as well.  This check is intended to
	// be a quick check to weed out duplicates.
	if mp.haveTransaction(txHash) {
		str := fmt.Sprintf("already have transaction %v", txHash)
		return nil, txRuleError(wire.RejectDuplicate, str)
	}

	// Perform preliminary sanity checks on the transaction.  This makes
	// use of chain which contains the invariant rules for what
	// transactions are allowed into blocks.
	err := blockchain.CheckTransactionSanity(msgTx, mp.cfg.ChainParams)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	// A standalone transaction must not be a coinbase transaction.
	if blockchain.IsCoinBase(tx) {
		str := fmt.Sprintf("transaction %v is an individual coinbase",
			txHash)
		return nil, txRuleError(wire.RejectInvalid, str)
	}

	// Don't accept transactions with a lock time after the maximum int32
	// value for now.  This is an artifact of older bitcoind clients which
	// treated this field as an int32 and would treat anything larger
	// incorrectly (as negative).
	// 	if msgTx.LockTime > math.MaxInt32 {
	// 		str := fmt.Sprintf("transaction %v has a lock time after "+
	// 			"2038 which is not accepted yet", txHash)
	// 		return nil, txRuleError(wire.RejectNonstandard, str)
	// 	}

	// Get the current height of the main chain.  A standalone transaction
	// will be mined into the next block at best, so its height is at least
	// one more than the current height.
	bestHeight := mp.cfg.BestHeight()
	nextBlockHeight := bestHeight + 1

	// Determine what type of transaction we're dealing with (regular or stake).
	// Then, be sure to set the tx tree correctly as it's possible a use submitted
	// it to the network with TxTreeUnknown.
	txType := stake.DetermineTxType(msgTx)
	if txType == stake.TxTypeRegular {
		tx.SetTree(wire.TxTreeRegular)
	} else {
		return nil, txRuleError(wire.RejectNonstandard, "instant transaction  type must be regular")
	}

	// Don't allow non-standard transactions if the network parameters
	// forbid their relaying.
	medianTime := mp.cfg.PastMedianTime()
	if !mp.cfg.Policy.RelayNonStd {
		err := checkTransactionStandard(tx, txType, nextBlockHeight,
			medianTime, mp.cfg.Policy.MinRelayTxFee,
			mp.cfg.Policy.MaxTxVersion)
		if err != nil {
			// Attempt to extract a reject code from the error so
			// it can be retained.  When not possible, fall back to
			// a non standard error.
			rejectCode, found := extractRejectCode(err)
			if !found {
				rejectCode = wire.RejectNonstandard
			}
			str := fmt.Sprintf("transaction %v is not standard: %v",
				txHash, err)
			return nil, txRuleError(rejectCode, str)
		}
	}

	// The transaction may not use any of the same outputs as other
	// transactions already in the pool as that would ultimately result in a
	// double spend.  This check is intended to be quick and therefore only
	// detects double spends within the transaction pool itself.  The
	// transaction could still be double spending coins from the main chain
	// at this point.  There is a more in-depth check that happens later
	// after fetching the referenced transaction inputs from the main chain
	// which examines the actual spend data and prevents double spends.
	err = mp.checkPoolDoubleSpend(tx, txType)
	if err != nil {
		return nil, err
	}

	// Fetch all of the unspent transaction outputs referenced by the inputs
	// to this transaction.  This function also attempts to fetch the
	// transaction itself to be used for detecting a duplicate transaction
	// without needing to do a separate lookup.
	utxoView, err := mp.fetchInputUtxos(tx)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	// Don't allow the transaction if it exists in the main chain and is not
	// not already fully spent.
	txEntry := utxoView.LookupEntry(txHash)
	if txEntry != nil && !txEntry.IsFullySpent() {
		return nil, txRuleError(wire.RejectDuplicate,
			"transaction already exists")
	}
	delete(utxoView.Entries(), *txHash)

	// Transaction is an orphan if any of the inputs don't exist.
	var missingParents []*chainhash.Hash
	for i, txIn := range msgTx.TxIn {
		if i == 0 && (txType == stake.TxTypeSSGen || txType == stake.TxTypeAiSSGen) {
			continue
		}

		entry := utxoView.LookupEntry(&txIn.PreviousOutPoint.Hash)
		if entry == nil || entry.IsFullySpent() {
			// Must make a copy of the hash here since the iterator
			// is replaced and taking its address directly would
			// result in all of the entries pointing to the same
			// memory location and thus all be the final hash.
			hashCopy := txIn.PreviousOutPoint.Hash
			missingParents = append(missingParents, &hashCopy)

			// Prevent a panic in the logger by continuing here if the
			// transaction input is nil.
			if entry == nil {
				log.Tracef("Transaction %v uses unknown input %v "+
					"and will be considered an orphan", txHash,
					txIn.PreviousOutPoint.Hash)
				continue
			}
			if entry.IsFullySpent() {
				log.Tracef("Transaction %v uses full spent input %v "+
					"and will be considered an orphan", txHash,
					txIn.PreviousOutPoint.Hash)
			}
		}
	}

	//instant tx don`t allow missing parents
	if len(missingParents) > 0 {
		return missingParents, txRuleError(wire.RejectNonstandard, "instant transaction inputs missing parents")
	}

	// Don't allow the transaction into the mempool unless its sequence
	// lock is active, meaning that it'll be allowed into the next block
	// with respect to its defined relative lock times.
	seqLock, err := mp.cfg.CalcSequenceLock(tx, utxoView)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}
	if !blockchain.SequenceLockActive(seqLock, nextBlockHeight, medianTime) {
		return nil, txRuleError(wire.RejectNonstandard,
			"transaction sequence locks on inputs not met")
	}

	// Perform several checks on the transaction inputs using the invariant
	// rules in chain for what transactions are allowed into blocks.
	// Also returns the fees associated with the transaction which will be
	// used later.  The fraud proof is not checked because it will be
	// filled in by the miner.
	txFee, err := blockchain.CheckTransactionInputs(mp.cfg.SubsidyCache,
		tx, nextBlockHeight, utxoView, false, mp.cfg.ChainParams)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	// Don't allow transactions with non-standard inputs if the network
	// parameters forbid their relaying.
	if !mp.cfg.Policy.RelayNonStd {
		err := checkInputsStandard(tx, txType, utxoView)
		if err != nil {
			// Attempt to extract a reject code from the error so
			// it can be retained.  When not possible, fall back to
			// a non standard error.
			rejectCode, found := extractRejectCode(err)
			if !found {
				rejectCode = wire.RejectNonstandard
			}
			str := fmt.Sprintf("transaction %v has a non-standard "+
				"input: %v", txHash, err)
			return nil, txRuleError(rejectCode, str)
		}
	}

	// NOTE: if you modify this code to accept non-standard transactions,
	// you should add code here to check that the transaction does a
	// reasonable number of ECDSA signature verifications.

	// Don't allow transactions with an excessive number of signature
	// operations which would result in making it impossible to mine.  Since
	// the coinbase address itself can contain signature operations, the
	// maximum allowed signature operations per transaction is less than
	// the maximum allowed signature operations per block.
	numSigOps, err := blockchain.CountP2SHSigOps(tx, false,
		(txType == stake.TxTypeSSGen || txType == stake.TxTypeAiSSGen), utxoView)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	numSigOps += blockchain.CountSigOps(tx, false, (txType == stake.TxTypeSSGen || txType == stake.TxTypeAiSSGen))
	if numSigOps > mp.cfg.Policy.MaxSigOpsPerTx {
		str := fmt.Sprintf("transaction %v has too many sigops: %d > %d",
			txHash, numSigOps, mp.cfg.Policy.MaxSigOpsPerTx)
		return nil, txRuleError(wire.RejectNonstandard, str)
	}

	// Don't allow transactions with fees too low to get into a mined block.
	//
	// Most miners allow a free transaction area in blocks they mine to go
	// alongside the area used for high-priority transactions as well as
	// transactions with fees.  A transaction size of up to 1000 bytes is
	// considered safe to go into this section.  Further, the minimum fee
	// calculated below on its own would encourage several small
	// transactions to avoid fees rather than one single larger transaction
	// which is more desirable.  Therefore, as long as the size of the
	// transaction does not exceeed 1000 less than the reserved space for
	// high-priority transactions, don't require a fee for it.
	// This applies to non-stake transactions only.
	serializedSize := int64(msgTx.SerializeSize())
	minFee := calcMinRequiredTxRelayFee(serializedSize,
		mp.cfg.Policy.MinRelayTxFee)
	if txType == stake.TxTypeRegular { // Non-stake only
		if serializedSize >= (DefaultBlockPrioritySize-1000) &&
			txFee < minFee {

			str := fmt.Sprintf("transaction %v has %v fees which "+
				"is under the required amount of %v", txHash,
				txFee, minFee)
			return nil, txRuleError(wire.RejectInsufficientFee, str)
		}
	}

	// Require that free transactions have sufficient priority to be mined
	// in the next block.  Transactions which are being added back to the
	// memory pool from blocks that have been disconnected during a reorg
	// are exempted.
	//
	// This applies to non-stake transactions only.
	if isNew && !mp.cfg.Policy.DisableRelayPriority && txFee < minFee &&
		txType == stake.TxTypeRegular {

		currentPriority := CalcPriority(msgTx, utxoView,
			nextBlockHeight)
		if currentPriority <= MinHighPriority {
			str := fmt.Sprintf("transaction %v has insufficient priority (%g <= %g)", txHash,
				currentPriority, MinHighPriority)
			return nil, txRuleError(wire.RejectInsufficientFee, str)
		}
	}

	// Free-to-relay transactions are rate limited here to prevent
	// penny-flooding with tiny transactions as a form of attack.
	// This applies to non-stake transactions only.
	if rateLimit && txFee < minFee && txType == stake.TxTypeRegular {
		nowUnix := time.Now().Unix()
		// Decay passed data with an exponentially decaying ~10 minute
		// window.
		mp.pennyTotal *= math.Pow(1.0-1.0/600.0,
			float64(nowUnix-mp.lastPennyUnix))
		mp.lastPennyUnix = nowUnix

		// Are we still over the limit?
		if mp.pennyTotal >= mp.cfg.Policy.FreeTxRelayLimit*10*1000 {
			str := fmt.Sprintf("transaction %v has been rejected "+
				"by the rate limiter due to low fees", txHash)
			return nil, txRuleError(wire.RejectInsufficientFee, str)
		}
		oldTotal := mp.pennyTotal

		mp.pennyTotal += float64(serializedSize)
		log.Tracef("rate limit: curTotal %v, nextTotal: %v, "+
			"limit %v", oldTotal, mp.pennyTotal,
			mp.cfg.Policy.FreeTxRelayLimit*10*1000)
	}

	// Check whether allowHighFees is set to false (default), if so, then make
	// sure the current fee is sensible.  If people would like to avoid this
	// check then they can AllowHighFees = true
	if !allowHighFees {
		maxFee := calcMinRequiredTxRelayFee(serializedSize*maxRelayFeeMultiplier,
			mp.cfg.Policy.MinRelayTxFee)
		if txFee > maxFee {
			err = fmt.Errorf("transaction %v has %v fee which is above the "+
				"allowHighFee check threshold amount of %v", txHash,
				txFee, maxFee)
			return nil, err
		}
	}

	// Verify crypto signatures for each input and reject the transaction if
	// any don't verify.
	flags, err := mp.cfg.Policy.StandardVerifyFlags()
	if err != nil {
		return nil, err
	}
	err = blockchain.ValidateTransactionScripts(tx, utxoView, flags,
		mp.cfg.SigCache)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	return nil, nil
}

func (mp *TxPool) FetchInstantTx(txHash *chainhash.Hash, includeRecentBlock bool) (*hcutil.InstantTx, error) {
	// Protect concurrent access.
	mp.mtx.RLock()
	txDesc, exists := mp.txLockPool[*txHash]
	mp.mtx.RUnlock()

	if exists {
		return txDesc.Tx, nil
	}

	tx, err := mp.FetchTransaction(txHash, includeRecentBlock)
	if err != nil {
		return nil, err
	}
	msgInstantTx := wire.NewMsgInstantTx()
	msgInstantTx.MsgTx = *tx.MsgTx()
	instantTx := hcutil.NewInstantTx(msgInstantTx)
	instantTx.SetTree(tx.Tree())
	instantTx.SetIndex(tx.Index())

	return instantTx, nil
}

func (mp *TxPool) FetchInstantTxVote(txVoteHash *chainhash.Hash) (*hcutil.InstantTxVote, error) {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()
	return mp.fetchInstantTxVote(txVoteHash)
}

func (mp *TxPool) fetchInstantTxVote(txVoteHash *chainhash.Hash) (*hcutil.InstantTxVote, error) {
	if instantTxVote, exist := mp.instantTxVotes[*txVoteHash]; exist {
		return instantTxVote, nil
	}
	return nil, errors.New("instantTxVote not exist ")
}
