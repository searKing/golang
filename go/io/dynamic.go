// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import "io"

// DynamicReaderFunc returns a Reader that's from
// the provided input reader getter function.
type DynamicReaderFunc func() io.Reader

func (f DynamicReaderFunc) Read(p []byte) (n int, err error) {
	var r io.Reader
	if f != nil {
		r = f()
	}
	if r == nil {
		r = EOFReader()
	}
	return r.Read(p)
}
