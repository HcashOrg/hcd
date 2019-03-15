// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/HcashOrg/hcd/blockchain/stake"
	"github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/txscript"
	"github.com/HcashOrg/hcd/wire"
	"github.com/HcashOrg/hcd/hcutil"
	"math"
)

// The number of values to precalculate on initialization of the subsidy
// cache.
const subsidyCacheInitWidth = 4

// SubsidyCache is a structure that caches calculated values of subsidy so that
// they're not constantly recalculated. The blockchain struct itself possesses a
// pointer to a preinitialized SubsidyCache.
type SubsidyCache struct {
	subsidyCache     map[uint64]int64
	subsidyCacheLock sync.RWMutex

	params *chaincfg.Params
}

// NewSubsidyCache initializes a new subsidy cache for a given height. It
// precalculates the values of the subsidy that are most likely to be seen by
// the client when it connects to the network.
func NewSubsidyCache(height int64, params *chaincfg.Params) *SubsidyCache {
	scm := make(map[uint64]int64)
	sc := SubsidyCache{
		subsidyCache: scm,
		params:       params,
	}

	iteration := uint64(height / params.SubsidyReductionInterval)
	if iteration < subsidyCacheInitWidth {
		return &sc
	}

	for i := iteration - 4; i <= iteration; i++ {
		sc.CalcBlockSubsidy(int64(iteration) * params.SubsidyReductionInterval)
	}

	return &sc
}

// CalcBlockSubsidy returns the subsidy amount a block at the provided height
// should have. This is mainly used for determining how much the coinbase for
// newly generated blocks awards as well as validating the coinbase for blocks
// has the expected value.
//
// Subsidy calculation for exponential reductions:
// 0 for i in range (0, height / SubsidyReductionInterval):
// 1     subsidy *= MulSubsidy
// 2     subsidy /= DivSubsidy
//
// Safe for concurrent access.
func (s *SubsidyCache) CalcBlockSubsidy(height int64) int64 {
	// Block height 1 subsidy is 'special' and used to
	// distribute initial tokens, if any.
	if height == 1 {
		return s.params.BlockOneSubsidy()
	}
	if uint64(height) >= s.params.UpdateHeight {
		return s.CalcBlockSubsidyV2(height)
	}

	iteration := uint64(height / s.params.SubsidyReductionInterval)

	if iteration == 0 {
		return s.params.BaseSubsidy
	}

	// First, check the cache.
	s.subsidyCacheLock.RLock()
	cachedValue, existsInCache := s.subsidyCache[iteration]
	s.subsidyCacheLock.RUnlock()
	if existsInCache {
		return cachedValue
	}

	// Is the previous one in the cache? If so, calculate
	// the subsidy from the previous known value and store
	// it in the database and the cache.
	//A(n) = (a1+(n-1)d)q^(n-1)
	//S(n) = a1(1-q^n)/(1-q) + d[q(1-q^(n-1))/((1-q)^2) - (n-1)q^n/(1-q)]
	//A(n) = A(n-1) *q + d*q^(n-1)

	var q float64 = float64(s.params.MulSubsidy)/float64(s.params.DivSubsidy)
	var temp float64 = 0.0

	if iteration < 1682 {
		temp = float64(s.params.BaseSubsidy) * (1.0 - float64(iteration) * 59363.0 / 100000000.0) * math.Pow(q,float64(iteration))
	}else{//after 99 years
		temp = 100000000.0/float64(s.params.SubsidyReductionInterval) * math.Pow(0.1, float64(float64(iteration)-1681.0))
	}
	subsidy := int64(temp)
	s.subsidyCacheLock.Lock()
	s.subsidyCache[iteration] = subsidy
	s.subsidyCacheLock.Unlock()
	return subsidy
}

func (s *SubsidyCache) CalcBlockSubsidyV2(height int64) int64 {
	// Block height 1 subsidy is 'special' and used to
	// distribute initial tokens, if any.
	if height == 1 {
		return s.params.BlockOneSubsidy()
	}

	iteration := uint64(height / s.params.SubsidyReductionInterval)

	if iteration == 0 {
		return s.params.BaseSubsidy
	}

	// First, check the cache.
	s.subsidyCacheLock.RLock()
	cachedValue, existsInCache := s.subsidyCache[iteration]
	s.subsidyCacheLock.RUnlock()
	if existsInCache {
		return cachedValue
	}

	// Is the previous one in the cache? If so, calculate
	// the subsidy from the previous known value and store
	// it in the database and the cache.
	//A(n) = (a1+(n-1)d)q^(n-1)
	//S(n) = a1(1-q^n)/(1-q) + d[q(1-q^(n-1))/((1-q)^2) - (n-1)q^n/(1-q)]
	//A(n) = A(n-1) *q + d*q^(n-1)

	var q float64 = float64(s.params.MulSubsidyV2)/float64(s.params.DivSubsidy)
	var temp float64 = 0.0

	if iteration < 4205 {
		temp = float64(s.params.BaseSubsidyV2) * (1.0 + float64(iteration) * 0.00331) * math.Pow(q,float64(iteration))
	}else{//after 99 years
		temp = 100000000.0/float64(s.params.SubsidyReductionInterval) * math.Pow(0.1, float64(float64(iteration)-4204.0))
	}
	subsidy := int64(temp)
	s.subsidyCacheLock.Lock()
	s.subsidyCache[iteration] = subsidy
	s.subsidyCacheLock.Unlock()
	return subsidy
}

// CalcBlockWorkSubsidy calculates the proof of work subsidy for a block as a
// proportion of the total subsidy.
func CalcBlockWorkSubsidy(subsidyCache *SubsidyCache, height int64,
	voters uint16, params *chaincfg.Params) int64 {
	subsidy := subsidyCache.CalcBlockSubsidy(height)

	proportionWork := int64(params.WorkRewardProportion)
	proportions := int64(params.TotalSubsidyProportions())
	subsidy *= proportionWork
	subsidy /= proportions

	// Ignore the voters field of the header before we're at a point
	// where there are any voters.
	if height < params.StakeValidationHeight {
		return subsidy
	}

	// If there are no voters, subsidy is 0. The block will fail later anyway.
	if voters == 0 {
		return 0
	}

	// Adjust for the number of voters. This shouldn't ever overflow if you start
	// with 50 * 10^8 Atoms and voters and potentialVoters are uint16.
	potentialVoters := params.TicketsPerBlock
	actual := (int64(voters) * subsidy) / int64(potentialVoters)

	return actual
}

// CalcStakeVoteSubsidy calculates the subsidy for a stake vote based on the height
// of its input SStx.
//
// Safe for concurrent access.
func CalcStakeVoteSubsidy(subsidyCache *SubsidyCache, height int64,
	params *chaincfg.Params) int64 {
	// Calculate the actual reward for this block, then further reduce reward
	// proportional to StakeRewardProportion.
	// Note that voters/potential voters is 1, so that vote reward is calculated
	// irrespective of block reward.
	subsidy := subsidyCache.CalcBlockSubsidy(height)

	proportionStake := int64(params.StakeRewardProportion)
	proportions := int64(params.TotalSubsidyProportions())
	subsidy *= proportionStake
	subsidy /= (proportions * int64(params.TicketsPerBlock))

	return subsidy
}

// CalcBlockTaxSubsidy calculates the subsidy for the organization address in the
// coinbase.
//
// Safe for concurrent access.
func CalcBlockTaxSubsidy(subsidyCache *SubsidyCache, height int64, voters uint16,
	params *chaincfg.Params) int64 {
	if params.BlockTaxProportion == 0 {
		return 0
	}

	subsidy := subsidyCache.CalcBlockSubsidy(height)

	proportionTax := int64(params.BlockTaxProportion)
	proportions := int64(params.TotalSubsidyProportions())
	subsidy *= proportionTax
	subsidy /= proportions

	// Assume all voters 'present' before stake voting is turned on.
	if height < params.StakeValidationHeight {
		voters = 5
	}

	// If there are no voters, subsidy is 0. The block will fail later anyway.
	if voters == 0 && height >= params.StakeValidationHeight {
		return 0
	}

	// Adjust for the number of voters. This shouldn't ever overflow if you start
	// with 50 * 10^8 Atoms and voters and potentialVoters are uint16.
	potentialVoters := params.TicketsPerBlock
	adjusted := (int64(voters) * subsidy) / int64(potentialVoters)

	return adjusted
}

// BlockOneCoinbasePaysTokens checks to see if the first block coinbase pays
// out to the network initial token ledger.
func BlockOneCoinbasePaysTokens(tx *hcutil.Tx,
	params *chaincfg.Params) error {
	// If no ledger is specified, just return true.
	if len(params.BlockOneLedger) == 0 {
		return nil
	}

	if tx.MsgTx().LockTime != 0 {
		errStr := fmt.Sprintf("block 1 coinbase has invalid locktime")
		return ruleError(ErrBlockOneTx, errStr)
	}

	if tx.MsgTx().Expiry != wire.NoExpiryValue {
		errStr := fmt.Sprintf("block 1 coinbase has invalid expiry")
		return ruleError(ErrBlockOneTx, errStr)
	}

	if tx.MsgTx().TxIn[0].Sequence != wire.MaxTxInSequenceNum {
		errStr := fmt.Sprintf("block 1 coinbase not finalized")
		return ruleError(ErrBlockOneInputs, errStr)
	}

	if len(tx.MsgTx().TxOut) == 0 {
		errStr := fmt.Sprintf("coinbase outputs empty in block 1")
		return ruleError(ErrBlockOneOutputs, errStr)
	}

	ledger := params.BlockOneLedger
	if len(ledger) != len(tx.MsgTx().TxOut) {
		errStr := fmt.Sprintf("wrong number of outputs in block 1 coinbase; "+
			"got %v, expected %v", len(tx.MsgTx().TxOut), len(ledger))
		return ruleError(ErrBlockOneOutputs, errStr)
	}

	// Check the addresses and output amounts against those in the ledger.
	for i, txout := range tx.MsgTx().TxOut {
		if txout.Version != txscript.DefaultScriptVersion {
			errStr := fmt.Sprintf("bad block one output version; want %v, got %v",
				txscript.DefaultScriptVersion, txout.Version)
			return ruleError(ErrBlockOneOutputs, errStr)
		}

		// There should only be one address.
		_, addrs, _, err :=
			txscript.ExtractPkScriptAddrs(txout.Version, txout.PkScript, params)
		if err != nil {
			return ruleError(ErrBlockOneOutputs, err.Error())
		}
		if len(addrs) != 1 {
			errStr := fmt.Sprintf("too many addresses in output")
			return ruleError(ErrBlockOneOutputs, errStr)
		}

		addrLedger, err := hcutil.DecodeAddress(ledger[i].Address)
		if err != nil {
			return err
		}

		if !bytes.Equal(addrs[0].ScriptAddress(), addrLedger.ScriptAddress()) {
			errStr := fmt.Sprintf("address in output %v has non matching "+
				"address; got %v (hash160 %x), want %v (hash160 %x)",
				i,
				addrs[0].EncodeAddress(),
				addrs[0].ScriptAddress(),
				addrLedger.EncodeAddress(),
				addrLedger.ScriptAddress())
			return ruleError(ErrBlockOneOutputs, errStr)
		}

		if txout.Value != ledger[i].Amount {
			errStr := fmt.Sprintf("address in output %v has non matching "+
				"amount; got %v, want %v", i, txout.Value, ledger[i].Amount)
			return ruleError(ErrBlockOneOutputs, errStr)
		}
	}

	return nil
}

// CoinbasePaysTax checks to see if a given block's coinbase correctly pays
// tax to the developer organization.
func CoinbasePaysTax(subsidyCache *SubsidyCache, tx *hcutil.Tx, height uint32,
	voters uint16, params *chaincfg.Params) error {
	// Taxes only apply from block 2 onwards.
	if height <= 1 {
		return nil
	}

	// Tax is disabled.
	if params.BlockTaxProportion == 0 {
		return nil
	}

	if len(tx.MsgTx().TxOut) == 0 {
		errStr := fmt.Sprintf("invalid coinbase (no outputs)")
		return ruleError(ErrNoTxOutputs, errStr)
	}

	taxOutput := tx.MsgTx().TxOut[0]
	if taxOutput.Version != params.OrganizationPkScriptVersion {
		return ruleError(ErrNoTax,
			"coinbase tax output uses incorrect script version")
	}
	if !bytes.Equal(taxOutput.PkScript, params.OrganizationPkScript) {
		return ruleError(ErrNoTax,
			"coinbase tax output script does not match the "+
				"required script")
	}

	// Get the amount of subsidy that should have been paid out to
	// the organization, then check it.
	orgSubsidy := CalcBlockTaxSubsidy(subsidyCache, int64(height), voters, params)
	if orgSubsidy != taxOutput.Value {
		errStr := fmt.Sprintf("amount in output 0 has non matching org "+
			"calculated amount; got %v, want %v", taxOutput.Value,
			orgSubsidy)
		return ruleError(ErrNoTax, errStr)
	}

	return nil
}

// CalculateAddedSubsidy calculates the amount of subsidy added by a block
// and its parent. The blocks passed to this function MUST be valid blocks
// that have already been confirmed to abide by the consensus rules of the
// network, or the function might panic.
func CalculateAddedSubsidy(block, parent *hcutil.Block) int64 {
	var subsidy int64

	regularTxTreeValid := hcutil.IsFlagSet16(block.MsgBlock().Header.VoteBits,
		hcutil.BlockValid)
	if regularTxTreeValid {
		subsidy += parent.MsgBlock().Transactions[0].TxIn[0].ValueIn
	}

	for _, stx := range block.MsgBlock().STransactions {
		if isSSGen, _ := stake.IsSSGen(stx); isSSGen {
			subsidy += stx.TxIn[0].ValueIn
		}
	}

	return subsidy
}
