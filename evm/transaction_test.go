package evm

import (
	"github.com/HcashOrg/hcd/evm/common"
	"github.com/HcashOrg/hcd/evm/rawdb"
	"github.com/HcashOrg/hcd/evm/state"
	"github.com/HcashOrg/hcd/evm/types"
	"math/big"
	"testing"
)

type chainContextTest struct{
}


func (c chainContextTest)GetHeader(common.Hash, uint64) *types.Header  {
	return &types.Header{}
}

func TestCreateAccount(t *testing.T) {
	db,err:=rawdb.NewLevelDBDatabase("testdata",5,5,"world")
	if err!=nil{
		t.Error(err)
	}
	stateDb, err := state.New(common.HexToHash("532dab88a3fbdcaed934a58a6d2fddb0761276853a0dfd74a614098ead7e09fa"), state.NewDatabaseWithCache(db, 0))

	stateDb.SetBalance(common.BytesToAddress([]byte("bob")),big.NewInt(50000000))

	root:=stateDb.IntermediateRoot(false)
	t.Log("root",root.String())
	stateDb.Commit(true)
	stateDb.Database().TrieDB().Commit(root,true)
}


