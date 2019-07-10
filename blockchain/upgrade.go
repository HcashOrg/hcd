// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"github.com/HcashOrg/hcd/blockchain/aistake"
	"github.com/HcashOrg/hcd/blockchain/internal/progresslog"
	"github.com/HcashOrg/hcd/blockchain/stake"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/database"
)

// upgradeToVersion2 upgrades a version 1 blockchain to version 2, allowing
// use of the new on-disk ticket database.
func (b *BlockChain) upgradeToVersion2() error {
	log.Infof("Initializing upgrade to database version 2")
	best := b.BestSnapshot()
	progressLogger := progresslog.NewBlockProgressLogger("Upgraded", log)

	// The upgrade is atomic, so there is no need to set the flag that
	// the database is undergoing an upgrade here.  Get the stake node
	// for the genesis block, and then begin connecting stake nodes
	// incrementally.
	err := b.db.Update(func(dbTx database.Tx) error {
		bestStakeNode, errLocal := stake.InitDatabaseState(dbTx, b.chainParams)
		if errLocal != nil {
			return errLocal
		}

		parent, errLocal := dbFetchBlockByHeight(dbTx, 0)
		if errLocal != nil {
			return errLocal
		}

		for i := int64(1); i <= best.Height; i++ {
			block, errLocal := dbFetchBlockByHeight(dbTx, i)
			if errLocal != nil {
				return errLocal
			}

			// If we need the tickets, fetch them too.
			var newTickets []chainhash.Hash
			if i >= b.chainParams.StakeEnabledHeight {
				matureHeight := i - int64(b.chainParams.TicketMaturity)
				matureBlock, errLocal := dbFetchBlockByHeight(dbTx, matureHeight)
				if errLocal != nil {
					return errLocal
				}
				for _, stx := range matureBlock.MsgBlock().STransactions {
					if is, _ := stake.IsSStx(stx); is {
						h := stx.TxHash()
						newTickets = append(newTickets, h)
					}
				}
			}

			// Iteratively connect the stake nodes in memory.
			header := block.MsgBlock().Header
			tickets, _ := ticketsSpentInBlock(block)
			ticketsRv, _ := ticketsRevokedInBlock(block)

			bestStakeNode, errLocal = bestStakeNode.ConnectNode(header,
				tickets, ticketsRv, newTickets)
			if errLocal != nil {
				return errLocal
			}

			// Write the top block stake node to the database.
			errLocal = stake.WriteConnectedBestNode(dbTx, bestStakeNode,
				*best.Hash)
			if errLocal != nil {
				return errLocal
			}

			// Write the best block node when we reach it.
			if i == best.Height {
				tickets, _ := ticketsSpentInBlock(block)
				ticketsRv, _ := ticketsRevokedInBlock(block)

				b.bestNode.stakeNode = bestStakeNode
				b.bestNode.stakeUndoData = bestStakeNode.UndoData()
				b.bestNode.newTickets = newTickets
				b.bestNode.ticketsSpent = tickets
				b.bestNode.ticketsRevoked = ticketsRv

				//b.bestNode.aistakeNode = bestAiStakeNode
				//b.bestNode.aistakeUndoData = bestAiStakeNode.UndoData()
				//b.bestNode.newAiTickets = newAiTickets
				//b.bestNode.aiTicketsSpent = aiTickets
				//b.bestNode.aiTicketsRevoked = aiTicketsRv
			}

			progressLogger.LogBlockHeight(block.MsgBlock(), parent.MsgBlock())
			parent = block
		}

		// Write the new database version.
		b.dbInfo.version = 2
		return dbPutDatabaseInfo(dbTx, b.dbInfo)
	})
	if err != nil {
		return err
	}

	log.Infof("Upgrade to new stake database was successful!")

	return nil
}


// upgradeToVersion3 upgrades a version 2 blockchain to version 3, allowing
// use of the new on-disk ticket database.
func (b *BlockChain) upgradeToVersion3() error {
	log.Infof("Initializing upgrade to database version 3")
	best := b.BestSnapshot()
	progressLogger := progresslog.NewBlockProgressLogger("Upgraded", log)

	// The upgrade is atomic, so there is no need to set the flag that
	// the database is undergoing an upgrade here.  Get the stake node
	// for the genesis block, and then begin connecting stake nodes
	// incrementally.
	err := b.db.Update(func(dbTx database.Tx) error {
		bestAiStakeNode, errLocal := aistake.InitDatabaseState(dbTx, b.chainParams)
		if errLocal != nil {
			return errLocal
		}

		parent, errLocal := dbFetchBlockByHeight(dbTx, 0)
		if errLocal != nil {
			return errLocal
		}

		for i := int64(b.chainParams.AIUpdateHeight); i <= best.Height; i++ {
			block, errLocal := dbFetchBlockByHeight(dbTx, i)
			if errLocal != nil {
				return errLocal
			}

			// If we need the tickets, fetch them too.
			//var newTickets []chainhash.Hash
			var newAiTickets []chainhash.Hash
			if i >= b.chainParams.StakeEnabledHeight {
				matureHeight := i - int64(b.chainParams.AiTicketMaturity)
				matureBlock, errLocal := dbFetchBlockByHeight(dbTx, matureHeight)
				if errLocal != nil {
					return errLocal
				}
				for _, stx := range matureBlock.MsgBlock().STransactions {
					/*if is, _ := stake.IsSStx(stx); is {
						h := stx.TxHash()
						newTickets = append(newTickets, h)
					}else */if is, _ := stake.IsAiSStx(stx); is {
						h := stx.TxHash()
						newAiTickets = append(newAiTickets, h)
					}
				}
			}

			// Iteratively connect the stake nodes in memory.
			header := block.MsgBlock().Header
			_, aiT := ticketsSpentInBlock(block)
			_, aiRv := ticketsRevokedInBlock(block)

			bestAiStakeNode, errLocal = bestAiStakeNode.ConnectNode(header,
				aiT, aiRv, newAiTickets)
			if errLocal != nil {
				return errLocal
			}
			// Write the top block stake node to the database.
			errLocal = aistake.WriteConnectedBestNode(dbTx, bestAiStakeNode,
				*best.Hash)
			if errLocal != nil {
				return errLocal
			}

			// Write the best block node when we reach it.
			if i == best.Height {
				_, aiTickets := ticketsSpentInBlock(block)
				_, aiTicketsRv := ticketsRevokedInBlock(block)

//				b.bestNode.stakeNode = bestStakeNode
//				b.bestNode.stakeUndoData = bestStakeNode.UndoData()
//				b.bestNode.newTickets = newTickets
//				b.bestNode.ticketsSpent = tickets
//				b.bestNode.ticketsRevoked = ticketsRv

				b.bestNode.aistakeNode = bestAiStakeNode
				b.bestNode.aistakeUndoData = bestAiStakeNode.UndoData()
				b.bestNode.newAiTickets = newAiTickets
				b.bestNode.aiTicketsSpent = aiTickets
				b.bestNode.aiTicketsRevoked = aiTicketsRv

			}

			progressLogger.LogBlockHeight(block.MsgBlock(), parent.MsgBlock())
			parent = block
		}
		if uint64(best.Height) < b.chainParams.AIUpdateHeight {
			b.bestNode.aistakeNode = aistake.NullNode( b.chainParams, uint32(best.Height))
		}

		// Write the new database version.
		b.dbInfo.version = 3
		return dbPutDatabaseInfo(dbTx, b.dbInfo)
	})
	if err != nil {
		return err
	}

	log.Infof("Upgrade to new stake database was successful!")

	return nil
}


// upgrade applies all possible upgrades to the blockchain database iteratively,
// updating old clients to the newest version.
func (b *BlockChain) upgrade() error {
	if b.dbInfo.version == 1 {
		err := b.upgradeToVersion2()
		if err != nil {
			return err
		}
	}
	if b.dbInfo.version == 2 {
		err := b.upgradeToVersion3()
		if err != nil {
			return err
		}
	}

	return nil
}
