// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

/*
Package txscript implements the HC transaction script language.

This package provides data structures and functions to parse and execute
HC transaction scripts.

Script Overview

Hcd transaction scripts are written in a stack-base, FORTH-like language.

The HC script language consists of a number of opcodes which fall into
several categories such pushing and popping data to and from the stack,
performing basic and bitwise arithmetic, conditional branching, comparing
hashes, and checking cryptographic signatures.  Scripts are processed from left
to right and intentionally do not provide loops.

The vast majority of Hcd scripts at the time of this writing are of several
standard forms which consist of a spender providing a public key and a signature
which proves the spender owns the associated private key.  This information
is used to prove the the spender is authorized to perform the transaction.

One benefit of using a scripting language is added flexibility in specifying
what conditions must be met in order to spend HCs.

Errors

Errors returned by this package are of type txscript.Error.  This allows the
caller to programmatically determine the specific error by examining the
ErrorCode field of the type asserted txscript.Error while still providing rich
error messages with contextual information.
*/
package txscript
