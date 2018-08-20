// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hcjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/HcashOrg/hcd/hcjson"
)

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hcjson.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &hcjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hcjson.NewAddMultisigAddressCmd(2, keys, hcjson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &hcjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   hcjson.String("test"),
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hcjson.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &hcjson.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return hcjson.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &hcjson.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &hcjson.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &hcjson.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &hcjson.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &hcjson.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &hcjson.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &hcjson.GetBalanceCmd{
				Account: nil,
				MinConf: hcjson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBalanceCmd(hcjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &hcjson.GetBalanceCmd{
				Account: hcjson.String("acct"),
				MinConf: hcjson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBalanceCmd(hcjson.String("acct"), hcjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &hcjson.GetBalanceCmd{
				Account: hcjson.String("acct"),
				MinConf: hcjson.Int(6),
			},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetNewAddressCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &hcjson.GetNewAddressCmd{
				Account:   nil,
				GapPolicy: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getnewaddress", "acct", "ignore")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetNewAddressCmd(hcjson.String("acct"), hcjson.String("ignore"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct","ignore"],"id":1}`,
			unmarshalled: &hcjson.GetNewAddressCmd{
				Account:   hcjson.String("acct"),
				GapPolicy: hcjson.String("ignore"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &hcjson.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetRawChangeAddressCmd(hcjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &hcjson.GetRawChangeAddressCmd{
				Account: hcjson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &hcjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: hcjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetReceivedByAccountCmd("acct", hcjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &hcjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: hcjson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &hcjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: hcjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetReceivedByAddressCmd("1Address", hcjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &hcjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: hcjson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &hcjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetTransactionCmd("123", hcjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &hcjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: hcjson.Bool(true),
			},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return hcjson.NewImportPrivKeyCmd("abc", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &hcjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  hcjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return hcjson.NewImportPrivKeyCmd("abc", hcjson.String("label"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &hcjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   hcjson.String("label"),
				Rescan:  hcjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return hcjson.NewImportPrivKeyCmd("abc", hcjson.String("label"), hcjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &hcjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   hcjson.String("label"),
				Rescan:  hcjson.Bool(false),
			},
		},
		{
			name: "importprivkey optional3",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("importprivkey", "abc", "label", false, 12345)
			},
			staticCmd: func() interface{} {
				return hcjson.NewImportPrivKeyCmd("abc", hcjson.String("label"), hcjson.Bool(false), hcjson.Int(12345))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false,12345],"id":1}`,
			unmarshalled: &hcjson.ImportPrivKeyCmd{
				PrivKey:  "abc",
				Label:    hcjson.String("label"),
				Rescan:   hcjson.Bool(false),
				ScanFrom: hcjson.Int(12345),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return hcjson.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &hcjson.KeyPoolRefillCmd{
				NewSize: hcjson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return hcjson.NewKeyPoolRefillCmd(hcjson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &hcjson.KeyPoolRefillCmd{
				NewSize: hcjson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &hcjson.ListAccountsCmd{
				MinConf: hcjson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListAccountsCmd(hcjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &hcjson.ListAccountsCmd{
				MinConf: hcjson.Int(6),
			},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &hcjson.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &hcjson.ListReceivedByAccountCmd{
				MinConf:          hcjson.Int(1),
				IncludeEmpty:     hcjson.Bool(false),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListReceivedByAccountCmd(hcjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &hcjson.ListReceivedByAccountCmd{
				MinConf:          hcjson.Int(6),
				IncludeEmpty:     hcjson.Bool(false),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListReceivedByAccountCmd(hcjson.Int(6), hcjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &hcjson.ListReceivedByAccountCmd{
				MinConf:          hcjson.Int(6),
				IncludeEmpty:     hcjson.Bool(true),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListReceivedByAccountCmd(hcjson.Int(6), hcjson.Bool(true), hcjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &hcjson.ListReceivedByAccountCmd{
				MinConf:          hcjson.Int(6),
				IncludeEmpty:     hcjson.Bool(true),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &hcjson.ListReceivedByAddressCmd{
				MinConf:          hcjson.Int(1),
				IncludeEmpty:     hcjson.Bool(false),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListReceivedByAddressCmd(hcjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &hcjson.ListReceivedByAddressCmd{
				MinConf:          hcjson.Int(6),
				IncludeEmpty:     hcjson.Bool(false),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListReceivedByAddressCmd(hcjson.Int(6), hcjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &hcjson.ListReceivedByAddressCmd{
				MinConf:          hcjson.Int(6),
				IncludeEmpty:     hcjson.Bool(true),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListReceivedByAddressCmd(hcjson.Int(6), hcjson.Bool(true), hcjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &hcjson.ListReceivedByAddressCmd{
				MinConf:          hcjson.Int(6),
				IncludeEmpty:     hcjson.Bool(true),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &hcjson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: hcjson.Int(1),
				IncludeWatchOnly:    hcjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListSinceBlockCmd(hcjson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &hcjson.ListSinceBlockCmd{
				BlockHash:           hcjson.String("123"),
				TargetConfirmations: hcjson.Int(1),
				IncludeWatchOnly:    hcjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListSinceBlockCmd(hcjson.String("123"), hcjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &hcjson.ListSinceBlockCmd{
				BlockHash:           hcjson.String("123"),
				TargetConfirmations: hcjson.Int(6),
				IncludeWatchOnly:    hcjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListSinceBlockCmd(hcjson.String("123"), hcjson.Int(6), hcjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &hcjson.ListSinceBlockCmd{
				BlockHash:           hcjson.String("123"),
				TargetConfirmations: hcjson.Int(6),
				IncludeWatchOnly:    hcjson.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &hcjson.ListTransactionsCmd{
				Account:          nil,
				Count:            hcjson.Int(10),
				From:             hcjson.Int(0),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListTransactionsCmd(hcjson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &hcjson.ListTransactionsCmd{
				Account:          hcjson.String("acct"),
				Count:            hcjson.Int(10),
				From:             hcjson.Int(0),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListTransactionsCmd(hcjson.String("acct"), hcjson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &hcjson.ListTransactionsCmd{
				Account:          hcjson.String("acct"),
				Count:            hcjson.Int(20),
				From:             hcjson.Int(0),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListTransactionsCmd(hcjson.String("acct"), hcjson.Int(20),
					hcjson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &hcjson.ListTransactionsCmd{
				Account:          hcjson.String("acct"),
				Count:            hcjson.Int(20),
				From:             hcjson.Int(1),
				IncludeWatchOnly: hcjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListTransactionsCmd(hcjson.String("acct"), hcjson.Int(20),
					hcjson.Int(1), hcjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &hcjson.ListTransactionsCmd{
				Account:          hcjson.String("acct"),
				Count:            hcjson.Int(20),
				From:             hcjson.Int(1),
				IncludeWatchOnly: hcjson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return hcjson.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &hcjson.ListUnspentCmd{
				MinConf:   hcjson.Int(1),
				MaxConf:   hcjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListUnspentCmd(hcjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &hcjson.ListUnspentCmd{
				MinConf:   hcjson.Int(6),
				MaxConf:   hcjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return hcjson.NewListUnspentCmd(hcjson.Int(6), hcjson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &hcjson.ListUnspentCmd{
				MinConf:   hcjson.Int(6),
				MaxConf:   hcjson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return hcjson.NewListUnspentCmd(hcjson.Int(6), hcjson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &hcjson.ListUnspentCmd{
				MinConf:   hcjson.Int(6),
				MaxConf:   hcjson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hcjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return hcjson.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1,"tree":0}]],"id":1}`,
			unmarshalled: &hcjson.LockUnspentCmd{
				Unlock: true,
				Transactions: []hcjson.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &hcjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hcjson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSendFromCmd("from", "1Address", 0.5, hcjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &hcjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hcjson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return hcjson.NewSendFromCmd("from", "1Address", 0.5, hcjson.Int(6),
					hcjson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &hcjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hcjson.Int(6),
				Comment:     hcjson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return hcjson.NewSendFromCmd("from", "1Address", 0.5, hcjson.Int(6),
					hcjson.String("comment"), hcjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &hcjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hcjson.Int(6),
				Comment:     hcjson.String("comment"),
				CommentTo:   hcjson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hcjson.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &hcjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hcjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hcjson.NewSendManyCmd("from", amounts, hcjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &hcjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hcjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hcjson.NewSendManyCmd("from", amounts, hcjson.Int(6), hcjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &hcjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hcjson.Int(6),
				Comment:     hcjson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &hcjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return hcjson.NewSendToAddressCmd("1Address", 0.5, hcjson.String("comment"),
					hcjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &hcjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   hcjson.String("comment"),
				CommentTo: hcjson.String("commentto"),
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &hcjson.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return hcjson.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &hcjson.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return hcjson.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &hcjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    hcjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hcjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return hcjson.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &hcjson.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]hcjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    hcjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hcjson.RawTxInput{}
				privKeys := []string{"abc"}
				return hcjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &hcjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]hcjson.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    hcjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []hcjson.RawTxInput{}
				privKeys := []string{}
				return hcjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					hcjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &hcjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]hcjson.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    hcjson.String("ALL"),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return hcjson.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &hcjson.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return hcjson.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &hcjson.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return hcjson.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &hcjson.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := hcjson.MarshalCmd(testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = hcjson.MarshalCmd(testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request hcjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = hcjson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}
