hcrpcclient
============

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

hcrpcclient implements a Websocket-enabled hc JSON-RPC client package
written in [Go](http://golang.org/).  It provides a robust and easy to use
client for interfacing with a hc RPC server that uses a
hcd/bitcoin core-like compatible hc JSON-RPC API.

## Status

This package is currently under active development and in a beta state.
There are still several RPCs left to implement and the API is not stable yet.

## Documentation

* [hcd Websockets Example](https://github.com/HcashOrg/hcrpcclient/blob/master/examples/hcwebsockets)  
  Connects to a hcd RPC server using TLS-secured websockets, registers for
  block connected and block disconnected notifications, and gets the current
  block count
* [hcwallet Websockets Example](https://github.com/HcashOrg/hcrpcclient/blob/master/examples/hcwalletwebsockets)  
  Connects to a hcwallet RPC server using TLS-secured websockets, registers for
  notifications about changes to account balances, and gets a list of unspent
  transaction outputs (utxos) the wallet can sign

## Major Features

* Supports Websockets (hcd/hcwallet) and HTTP POST mode (bitcoin core-like)
* Provides callback and registration functions for hcd/hcwallet notifications
* Supports hcd extensions
* Translates to and from higher-level and easier to use Go types
* Offers a synchronous (blocking) and asynchronous API
* When running in Websockets mode (the default):
  * Automatic reconnect handling (can be disabled)
  * Outstanding commands are automatically reissued
  * Registered notifications are automatically reregistered
  * Back-off support on reconnect attempts

## Installation

```bash
$ go get -u github.com/HcashOrg/hcrpcclient
```

## License

Package hcrpcclient, like hcrpcclient is licensed under the [copyfree](http://copyfree.org) ISC
License.
