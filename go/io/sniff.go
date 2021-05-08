// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"io"
)

// ReadSniffer is the interface that groups the basic Read and Sniff methods.
type ReadSniffer interface {
	io.Reader
	Sniff(sniffing bool) ReadSniffer
}

// SniffReader returns a Reader that allows sniff and read from
// the provided input reader.
// data is buffered if Sniff(true) is called.
// buffered data is taken first, if Sniff(false) is called.
func SniffReader(r io.Reader) ReadSniffer {
	if r, ok := r.(io.ReadSeeker); ok {
		return newSniffReadSeeker(r)
	}
	return newSniffReader(r)
}
