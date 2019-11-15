package evm

import (
	"github.com/HcashOrg/hcd/evm/common"
	"github.com/HcashOrg/hcd/evm/params"
	"github.com/HcashOrg/hcd/evm/rawdb"
	"github.com/HcashOrg/hcd/evm/state"
	"github.com/HcashOrg/hcd/evm/types"
	"github.com/HcashOrg/hcd/evm/vm"
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


/**
```
pragma solidity ^0.5.12;
contract DataStore{
    uint256 data;
    event posit(
        uint256 _value
    );
    function set(uint256 x) public{
        data = x;
    }
    function get() public  returns (uint256){
        emit posit(data);
        return data;
    }
}
```
**/


func TestDeploy(t *testing.T) {
	db,err:=rawdb.NewLevelDBDatabase("testdata",5,5,"world")
	if err!=nil{
		t.Error(err)
	}
	//root :=common.Hash{222 ,136 ,59 ,116, 50 ,27, 206, 138, 8, 57, 53, 52, 85, 228, 131, 215, 182, 208, 127, 142, 223, 26 ,103 ,19, 163, 185, 115, 235, 232, 26, 166 ,122}
	stateDb, err := state.New(common.HexToHash("532dab88a3fbdcaed934a58a6d2fddb0761276853a0dfd74a614098ead7e09fa"), state.NewDatabaseWithCache(db, 0))

	balance:=stateDb.GetBalance(common.BytesToAddress([]byte("contract")))

	t.Log("balance",balance)

	chainConfig := params.ChainConfig{
	}

	author:=common.BytesToAddress([]byte("author"))
	header:=types.Header{Number:big.NewInt(5),Difficulty:big.NewInt(5),GasLimit:50000000}


	tx:=types.NewContractCreation(6,
		common.BytesToAddress([]byte("contract")),
		big.NewInt(0),
		uint64(6000000),
		big.NewInt(1),
		common.Hex2Bytes("608060405234801561001057600080fd5b50610100806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806360fe47b11460375780636d4ce63c146062575b600080fd5b606060048036036020811015604b57600080fd5b8101908080359060200190929190505050607e565b005b60686088565b6040518082815260200191505060405180910390f35b8060008190555050565b60007f81e9ebcb91d869eb6085b37d6c6ebddfc79c22572cd44d73a68476e1ac7305a46000546040518082815260200191505060405180910390a160005490509056fea265627a7a723158204b85598e5a9e406055a4705a6480ba3da9293871b69d963eefbe5ab4b27aab9164736f6c634300050c0032"),)

	gasPool := new(GasPool).AddGas(header.GasLimit)
	vmConfig:=vm.Config{}
	if err!=nil{
		t.Error(err)
	}

	chainContext:= chainContextTest{}

	usedGas:=uint64(50)
	receipt,err:=ApplyTransaction(&chainConfig,chainContext,&author,gasPool,stateDb,&header,tx,&usedGas,vmConfig)
	if err!=nil{
		t.Fatal(err)
	}
	root:=stateDb.IntermediateRoot(false)
	t.Log("root",root.String())
	stateDb.Commit(false)
	stateDb.Database().TrieDB().Commit(root,true)
	t.Log(receipt)
	t.Log("contractAddress",receipt.ContractAddress.String())


}


func TestCall(t *testing.T) {
	db,err:=rawdb.NewLevelDBDatabase("testdata",5,5,"world")
	if err!=nil{
		t.Error(err)
	}
	stateDb, err := state.New(common.HexToHash("f3fcb0147e3619a53dd99f62842a627430400266a59e6816431ffcae8c5b10b8"), state.NewDatabaseWithCache(db, 0))

	balance:=stateDb.GetBalance(common.BytesToAddress([]byte("bob")))

	t.Log("balance",balance)

	chainConfig := params.ChainConfig{
	}

	author:=common.BytesToAddress([]byte("author"))
	header:=types.Header{Number:big.NewInt(5),Difficulty:big.NewInt(5),GasLimit:50000000}

	tx:=types.NewTransaction(
		0,
		common.BytesToAddress([]byte("bob")),
		common.HexToAddress("68656B2f7AC1b9f532101b0dA5348518e680D762"),
		big.NewInt(0),
		uint64(6000000),
		big.NewInt(1),
		common.Hex2Bytes("6d4ce63c"),
	)


	gasPool := new(GasPool).AddGas(header.GasLimit)
	vmConfig:=vm.Config{}
	if err!=nil{
		t.Error(err)
	}

	chainContext:= chainContextTest{}

	usedGas:=uint64(50)

	stateDb.Prepare(tx.Hash(), common.Hash{}, 0) //记录当前处理的tx信息，addlog 的时候会用到
	receipt,err:=ApplyTransaction(&chainConfig,chainContext,&author,gasPool,stateDb,&header,tx,&usedGas,vmConfig)
	if err!=nil{
		t.Fatal(err)
	}
	root:=stateDb.IntermediateRoot(false)
	t.Log("root",root.String())
	stateDb.Commit(false)
	stateDb.Database().TrieDB().Commit(root,true)
	t.Log(receipt)
	t.Log(receipt.ContractAddress.String())

}