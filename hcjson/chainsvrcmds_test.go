// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2016 The Decred developers
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

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal
// into valid results include handling of optional fields being omitted in the
// marshalled command, while optional fields with defaults have the default
// assigned on unmarshalled commands.
func TestChainSvrCmds(t *testing.T) {
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
			name: "addnode",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("addnode", "127.0.0.1", hcjson.ANRemove)
			},
			staticCmd: func() interface{} {
				return hcjson.NewAddNodeCmd("127.0.0.1", hcjson.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &hcjson.AddNodeCmd{Addr: "127.0.0.1", SubCmd: hcjson.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []hcjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return hcjson.NewCreateRawTransactionCmd(txInputs, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1,"tree":0}],{"456":0.0123}],"id":1}`,
			unmarshalled: &hcjson.CreateRawTransactionCmd{
				Inputs:  []hcjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1,"tree":0}]`,
					`{"456":0.0123}`, int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []hcjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return hcjson.NewCreateRawTransactionCmd(txInputs, amounts, hcjson.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1,"tree":0}],{"456":0.0123},12312333333],"id":1}`,
			unmarshalled: &hcjson.CreateRawTransactionCmd{
				Inputs:   []hcjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: hcjson.Int64(12312333333),
			},
		},
		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return hcjson.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &hcjson.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return hcjson.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &hcjson.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &hcjson.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetAddedNodeInfoCmd(true, hcjson.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &hcjson.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: hcjson.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &hcjson.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockCmd("123", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &hcjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hcjson.Bool(true),
				VerboseTx: hcjson.Bool(false),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {
				// Intentionally use a source param that is
				// more pointers than the destination to
				// exercise that path.
				verbosePtr := hcjson.Bool(true)
				return hcjson.NewCmd("getblock", "123", &verbosePtr)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockCmd("123", hcjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true],"id":1}`,
			unmarshalled: &hcjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hcjson.Bool(true),
				VerboseTx: hcjson.Bool(false),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblock", "123", true, true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockCmd("123", hcjson.Bool(true), hcjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true,true],"id":1}`,
			unmarshalled: &hcjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hcjson.Bool(true),
				VerboseTx: hcjson.Bool(true),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &hcjson.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &hcjson.GetBlockCountCmd{},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &hcjson.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &hcjson.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: hcjson.Bool(true),
			},
		},
		{
			name: "getblocksubsidy",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblocksubsidy", 123, 256)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockSubsidyCmd(123, 256)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocksubsidy","params":[123,256],"id":1}`,
			unmarshalled: &hcjson.GetBlockSubsidyCmd{
				Height: 123,
				Voters: 256,
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &hcjson.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := hcjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return hcjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &hcjson.GetBlockTemplateCmd{
				Request: &hcjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := hcjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return hcjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &hcjson.GetBlockTemplateCmd{
				Request: &hcjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := hcjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return hcjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &hcjson.GetBlockTemplateCmd{
				Request: &hcjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &hcjson.GetChainTipsCmd{},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &hcjson.GetConnectionCountCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getdifficulty")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":[],"id":1}`,
			unmarshalled: &hcjson.GetDifficultyCmd{},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &hcjson.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &hcjson.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &hcjson.GetInfoCmd{},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &hcjson.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &hcjson.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &hcjson.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &hcjson.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &hcjson.GetNetworkHashPSCmd{
				Blocks: hcjson.Int(120),
				Height: hcjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetNetworkHashPSCmd(hcjson.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &hcjson.GetNetworkHashPSCmd{
				Blocks: hcjson.Int(200),
				Height: hcjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetNetworkHashPSCmd(hcjson.Int(200), hcjson.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &hcjson.GetNetworkHashPSCmd{
				Blocks: hcjson.Int(200),
				Height: hcjson.Int(123),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &hcjson.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetRawMempoolCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &hcjson.GetRawMempoolCmd{
				Verbose: hcjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetRawMempoolCmd(hcjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &hcjson.GetRawMempoolCmd{
				Verbose: hcjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional 2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getrawmempool", false, "all")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetRawMempoolCmd(hcjson.Bool(false), hcjson.String("all"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false,"all"],"id":1}`,
			unmarshalled: &hcjson.GetRawMempoolCmd{
				Verbose: hcjson.Bool(false),
				TxType:  hcjson.String("all"),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &hcjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: hcjson.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetRawTransactionCmd("123", hcjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &hcjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: hcjson.Int(1),
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &hcjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: hcjson.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetTxOutCmd("123", 1, hcjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &hcjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: hcjson.Bool(true),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &hcjson.GetTxOutSetInfoCmd{},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &hcjson.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return hcjson.NewGetWorkCmd(hcjson.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &hcjson.GetWorkCmd{
				Data: hcjson.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return hcjson.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &hcjson.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return hcjson.NewHelpCmd(hcjson.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &hcjson.HelpCmd{
				Command: hcjson.String("getblock"),
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return hcjson.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &hcjson.PingCmd{},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return hcjson.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &hcjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcjson.Int(1),
				Skip:        hcjson.Int(0),
				Count:       hcjson.Int(100),
				VinExtra:    hcjson.Int(0),
				Reverse:     hcjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSearchRawTransactionsCmd("1Address",
					hcjson.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &hcjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcjson.Int(0),
				Skip:        hcjson.Int(0),
				Count:       hcjson.Int(100),
				VinExtra:    hcjson.Int(0),
				Reverse:     hcjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSearchRawTransactionsCmd("1Address",
					hcjson.Int(0), hcjson.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &hcjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcjson.Int(0),
				Skip:        hcjson.Int(5),
				Count:       hcjson.Int(100),
				VinExtra:    hcjson.Int(0),
				Reverse:     hcjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSearchRawTransactionsCmd("1Address",
					hcjson.Int(0), hcjson.Int(5), hcjson.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &hcjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcjson.Int(0),
				Skip:        hcjson.Int(5),
				Count:       hcjson.Int(10),
				VinExtra:    hcjson.Int(0),
				Reverse:     hcjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSearchRawTransactionsCmd("1Address",
					hcjson.Int(0), hcjson.Int(5), hcjson.Int(10), hcjson.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &hcjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcjson.Int(0),
				Skip:        hcjson.Int(5),
				Count:       hcjson.Int(10),
				VinExtra:    hcjson.Int(1),
				Reverse:     hcjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSearchRawTransactionsCmd("1Address",
					hcjson.Int(0), hcjson.Int(5), hcjson.Int(10),
					hcjson.Int(1), hcjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &hcjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcjson.Int(0),
				Skip:        hcjson.Int(5),
				Count:       hcjson.Int(10),
				VinExtra:    hcjson.Int(1),
				Reverse:     hcjson.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return hcjson.NewSearchRawTransactionsCmd("1Address",
					hcjson.Int(0), hcjson.Int(5), hcjson.Int(10),
					hcjson.Int(1), hcjson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &hcjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hcjson.Int(0),
				Skip:        hcjson.Int(5),
				Count:       hcjson.Int(10),
				VinExtra:    hcjson.Int(1),
				Reverse:     hcjson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendrawtransaction", "1122")
			},
			staticCmd: func() interface{} {
				return hcjson.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
			unmarshalled: &hcjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: hcjson.Bool(false),
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("sendrawtransaction", "1122", false)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSendRawTransactionCmd("1122", hcjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &hcjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: hcjson.Bool(false),
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSetGenerateCmd(true, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &hcjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: hcjson.Int(-1),
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("setgenerate", true, 6)
			},
			staticCmd: func() interface{} {
				return hcjson.NewSetGenerateCmd(true, hcjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6],"id":1}`,
			unmarshalled: &hcjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: hcjson.Int(6),
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return hcjson.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &hcjson.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return hcjson.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &hcjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := hcjson.SubmitBlockOptions{
					WorkID: "12345",
				}
				return hcjson.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &hcjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &hcjson.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return hcjson.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &hcjson.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return hcjson.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &hcjson.VerifyChainCmd{
				CheckLevel: hcjson.Int64(3),
				CheckDepth: hcjson.Int64(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return hcjson.NewVerifyChainCmd(hcjson.Int64(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &hcjson.VerifyChainCmd{
				CheckLevel: hcjson.Int64(2),
				CheckDepth: hcjson.Int64(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return hcjson.NewVerifyChainCmd(hcjson.Int64(2), hcjson.Int64(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &hcjson.VerifyChainCmd{
				CheckLevel: hcjson.Int64(2),
				CheckDepth: hcjson.Int64(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {
				return hcjson.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return hcjson.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &hcjson.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
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
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
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

// TestChainSvrCmdErrors ensures any errors that occur in the command during
// custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &hcjson.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &hcjson.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        hcjson.Error{Code: hcjson.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &hcjson.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        hcjson.Error{Code: hcjson.ErrInvalidType},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error type - got `%T` (%v), got `%T`",
				i, test.name, err, err, test.err)
			continue
		}

		if terr, ok := test.err.(hcjson.Error); ok {
			gotErrorCode := err.(hcjson.Error).Code
			if gotErrorCode != terr.Code {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.Code)
				continue
			}
		}
	}
}
