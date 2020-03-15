// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test_test

// The actual test functions are in non-_test.go files
// so that they can use cgo (import "C").
// These wrappers are here for gotest to find.
func ExampleGoStringArray() { ExampleGoStringArray() }
func ExampleCStringArray()  { ExampleCStringArray() }
