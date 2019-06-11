package mempool

import (
	"fmt"
	"github.com/HcashOrg/hcd/blockchain"
	"github.com/HcashOrg/hcd/blockchain/stake"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcd/txscript"
	"github.com/HcashOrg/hcd/wire"
)

const (
	defaultConfirmNum = 24
	defaultBehindNums = 10
)

type TxLockDesc struct {
	Tx *hcutil.Tx
	// Height is the block height when the entry was added to the the source
	// pool.
	AddHeight int64

	MineHeight int64 //
}

type lockPool struct {
	txLockPool    map[chainhash.Hash]*TxLockDesc //for instantsend lock tx pool
	lockOutpoints map[wire.OutPoint]*hcutil.Tx
}

//we will update tx state according the mined height
func (mp *TxPool) modifyLockTransaction(tx *hcutil.Tx, height int64) {
	msgTx := tx.MsgTx()
	isLockTx := false
	for _, txOut := range msgTx.TxOut {
		if txscript.IsLockTx(txOut.PkScript) {
			isLockTx = true
			break
		}
	}
	if !isLockTx {
		return
	}

	if desc, exist := mp.txLockPool[*tx.Hash()]; exist {
		desc.MineHeight = height
	}
}

func (mp *TxPool) ModifyLockTransaction(tx *hcutil.Tx, height int64) {
	// Protect concurrent access.
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	mp.modifyLockTransaction(tx, height)
}

func (mp *TxPool) RemoveConfirmedLockTransaction(height int64) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	for hash, desc := range mp.txLockPool {
		if desc.MineHeight != 0 && desc.MineHeight < height-defaultConfirmNum {
			delete(mp.txLockPool, hash)

			for _, txIn := range desc.Tx.MsgTx().TxIn {
				delete(mp.lockOutpoints, txIn.PreviousOutPoint)
			}
		}
	}
}

//Is tx in  locked?
func (mp *TxPool) isTxLockExist(hash *chainhash.Hash) bool {
	if _, exists := mp.txLockPool[*hash]; exists {
		return true
	}
	return false
}

//Is txVin  in locked?
func (mp *TxPool) isTxLockInExist(outPoint *wire.OutPoint) (*hcutil.Tx, bool) {
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
		ret[hash.String()] = &hcjson.TxLockInfo{AddHeight: desc.AddHeight, MineHeight: desc.MineHeight}
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
	for _, desc := range mp.txLockPool {
		if desc.MineHeight != 0 && desc.AddHeight < minHeight {
			bts, err := desc.Tx.MsgTx().Bytes()
			if err == nil {
				retMsgTx = append(retMsgTx, bts)
			}
		}
	}

	return retMsgTx

}

//check block transactions is conflict with lockPool .we can reject the conflict block by not notify the winningTicket
// to wallet
func (mp *TxPool) CheckConflictWithTxLockPool(block *hcutil.Block) (bool, error) {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()

	for _, tx := range block.Transactions() {
		if !mp.isTxLockExist(tx.Hash()) {
			for _, txIn := range tx.MsgTx().TxIn {
				if _, exist := mp.isTxLockInExist(&txIn.PreviousOutPoint); exist {
					return false, fmt.Errorf("lock transaction conflict")
				}
			}
		}
	}
	return true, nil
}

//remove txlock which is conflict with tx
func (mp *TxPool) RemoveTxLockDoubleSpends(tx *hcutil.Tx) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	//if is the same tx ,just return
	if mp.isTxLockExist(tx.Hash()) {
		return
	}

	//if tx in is conflict with txlock ,just remove txlock and lockOutpoint
	for _, invalue := range tx.MsgTx().TxIn {
		if txLock, exist := mp.isTxLockInExist(&invalue.PreviousOutPoint); exist {
			delete(mp.txLockPool, *txLock.Hash())

			for _, txIn := range txLock.MsgTx().TxIn {
				delete(mp.lockOutpoints, txIn.PreviousOutPoint)
			}
		}
	}

}

//this is called after insert to mempool
func (mp *TxPool) maybeAddtoLockPool(utxoView *blockchain.UtxoViewpoint,
	tx *hcutil.Tx, txType stake.TxType, height int64, fee int64) {

	//if exist just return ,or will rewrite the state of this txlock
	if mp.isTxLockExist(tx.Hash()) {
		return
	}

	msgTx := tx.MsgTx()
	isLockTx := false
	for _, txOut := range msgTx.TxOut {
		if txscript.IsLockTx(txOut.PkScript) {
			isLockTx = true
			break
		}
	}

	if isLockTx {
		mp.txLockPool[*tx.Hash()] = &TxLockDesc{Tx: tx, AddHeight: height, MineHeight: 0}

		for _, txIn := range msgTx.TxIn {
			mp.lockOutpoints[txIn.PreviousOutPoint] = tx
		}
	}
}
