// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

// TestCheckPubKeyEncoding ensures the internal checkPubKeyEncoding function
// works as expected.
func TstCheckPubKeyEncoding(pubKey []byte, flag ScriptFlags) error {
	eg := Engine{flags: flag}
	return eg.checkPubKeyEncoding(pubKey)
}