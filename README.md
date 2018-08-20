hcd
====

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

## What is HyperCash(HC)?
hcd is a Hc full node implementation written in Go (golang).

This acts as a chain daemon for the Hc cryptocurrency.
hcd maintains the entire past transactional ledger of Hc and allows
 relaying of quantum resistant transactions to other Hc nodes across the world.

Note: To send or receive funds and join Proof-of-Stake mining, you will also need
[hcwallet](https://github.com/HcashOrg/hcwallet).

HC is forked from [Decred](https://github.com/HcashOrg) and [btcd](https://github.com/btcsuite/btcd) which are full node implementations written in Go. Both projects are ongoing and under active development. Since hcd is synced and will merge with upstream commits from hcd and btcd, it will get the benefit of both hcd and btcd's ongoing upgrades to staking, voting, peer and connection handling, database optimization and other blockchain related technology improvements. Advances made by hcd can also be pulled back upstream to hcd and btcd including quantum resistant signature schemes and more.

## Development Process
The master branch is regularly built and tested, but is not guaranteed to be completely stable. Tags are created regularly to indicate new official, stable release versions of hcd.

## Requirements

[Go](http://golang.org) 1.7 or newer.

## Getting Started

- hcd (and utilities) will now be installed in either ```$GOROOT/bin``` or
  ```$GOPATH/bin``` depending on your configuration.  If you did not already
  add the bin directory to your system path during Go installation, we
  recommend you do so now.

## Installing

#### Build from Source

- **Dep**

  Dep is used to manage project dependencies and provide reproducible builds.
  To install:

  `go get -u github.com/golang/dep/cmd/dep`

Unfortunately, the use of `dep` prevents a handy tool such as `go get` from
automatically downloading, building, and installing the source in a single
command.  Instead, the latest project and dependency sources must be first
obtained manually with `git` and `dep`, and then `go` is used to build and
install the project.


**Getting the source**:

For a first time installation, the project and dependency sources can be
obtained manually with `git` and `dep` (create directories as needed):

```
git clone https://github.com/HcashOrg/hcd $GOPATH/src/github.com/HcashOrg/hcd
cd $GOPATH/src/github.com/HcashOrg/hcd
dep ensure
go install . ./cmd/...
```

To update an existing source tree, pull the latest changes and install the
matching dependencies:

```
cd $GOPATH/src/github.com/HcashOrg/hcd
git pull
dep ensure -update
go install . ./cmd/...
```

## Running

Make sure you are working the correct GOPATH and run the following in your terminal:

```
hcd -u YOURUNIQUERPCUSERNAME -P YOURUNIQUERPCPASSWORD --testnet
```

To use your node for mining add the miningaddr flag when running hcd:

```
hcd -u YOURUNIQUERPCUSERNAME -P YOURUNIQUERPCPASSWORD --testnet --miningaddr=YOURTESTNETADDRESS
```

To generate a testnet mining address you must install [hcwallet](https://github.com/HcashOrg/hcwallet)

To begin CPU mining after hcd is already running you can run the following in your terminal:

```
hcctl -u YOURUNIQUEUSERNAME -P YOURUNIQUEPASSWORD --testnet setgenerate true
```

## REINSTALL

#### Linux

```
$ rm -rf $HOME/.hcd
$ rm -rf $HOME/.hcwallet
$ cd $HOME/go/src/github.com/HcashOrg/hcd
$ git pull
$ go install $(glide nv)
```

#### Mac/OSX
```
$ rm -rf $HOME/Library/Application\ Support/Hcwallet
$ rm -rf $HOME/Library/Application\ Support/Hcd
$ cd $HOME/go/src/github.com/HcashOrg/hcd
$ git pull
$ go install $(glide nv)
```

## Issue Tracker

The [integrated github issue tracker](https://github.com/HcashOrg/hcd/issues)
is used for this project.

## Documentation

The documentation is a work-in-progress.  It is located in the [docs](https://github.com/HcashOrg/hcd/tree/master/docs) folder.

## License

hcd is licensed under the [copyfree](http://copyfree.org) ISC License.
