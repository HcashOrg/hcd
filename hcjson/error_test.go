// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hcjson_test

import (
	"testing"

	"github.com/HcashOrg/hcd/hcjson"
)

// TestErrorCodeStringer tests the stringized output for the ErrorCode type.
func TestErrorCodeStringer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   hcjson.ErrorCode
		want string
	}{
		{hcjson.ErrDuplicateMethod, "ErrDuplicateMethod"},
		{hcjson.ErrInvalidUsageFlags, "ErrInvalidUsageFlags"},
		{hcjson.ErrInvalidType, "ErrInvalidType"},
		{hcjson.ErrEmbeddedType, "ErrEmbeddedType"},
		{hcjson.ErrUnexportedField, "ErrUnexportedField"},
		{hcjson.ErrUnsupportedFieldType, "ErrUnsupportedFieldType"},
		{hcjson.ErrNonOptionalField, "ErrNonOptionalField"},
		{hcjson.ErrNonOptionalDefault, "ErrNonOptionalDefault"},
		{hcjson.ErrMismatchedDefault, "ErrMismatchedDefault"},
		{hcjson.ErrUnregisteredMethod, "ErrUnregisteredMethod"},
		{hcjson.ErrNumParams, "ErrNumParams"},
		{hcjson.ErrMissingDescription, "ErrMissingDescription"},
		{0xffff, "Unknown ErrorCode (65535)"},
	}

	// Detect additional error codes that don't have the stringer added.
	if len(tests)-1 != int(hcjson.TstNumErrorCodes) {
		t.Errorf("It appears an error code was added without adding an " +
			"associated stringer test")
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestError tests the error output for the Error type.
func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   hcjson.Error
		want string
	}{
		{
			hcjson.Error{Message: "some error"},
			"some error",
		},
		{
			hcjson.Error{Message: "human-readable error"},
			"human-readable error",
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.Error()
		if result != test.want {
			t.Errorf("Error #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}
