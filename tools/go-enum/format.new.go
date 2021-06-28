// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var newImportPackages = []string{}

// Arguments to format are:
//	[1]: type name
const newTemplate = `
// New returns a pointer to a new addr filled with the %[1]s value passed in.
func (i %[1]s) New () *%[1]s {
	clone := i
	return &clone
}
`
