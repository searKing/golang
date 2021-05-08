// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"io"
)

type sniffReadSeeker struct {
	source   io.ReadSeeker
	pos      int64 // history pos
	sniffing bool
}

func newSniffReadSeeker(r io.ReadSeeker) *sniffReadSeeker {
	sr := &sniffReadSeeker{
		source: r,
	}
	sr.init()
	return sr
}

func (sr *sniffReadSeeker) init() ReadSniffer {
	curPos, err := sr.source.Seek(0, io.SeekCurrent)
	if err != nil {
		// in case of source with io.ReadSeeker but not seekable actually, that is FakeReadSeeker
		return newSniffReader(sr.source)
	}
	sr.pos = curPos
	return sr
}

// Sniff starts or stops sniffing, restarts if stop and start called one by one
// true to start sniffing all data unread actually
// false to return a multi reader with all data sniff buffered and source
func (sr *sniffReadSeeker) Sniff(sniffing bool) ReadSniffer {
	if sr.sniffing == sniffing {
		return sr
	}
	sr.sniffing = sniffing
	if !sniffing {
		_, err := sr.source.Seek(sr.pos, io.SeekStart)
		if err != nil {
			// in case of source with io.ReadSeeker but not seekable actually, that is FakeReadSeeker
			return newSniffReader(sr.source)
		}
	}
	return sr
}

func (sr *sniffReadSeeker) Read(p []byte) (n int, err error) {
	n, err = sr.source.Read(p)
	if err != nil {
		return
	}
	if !sr.sniffing {
		sr.pos = sr.pos + int64(n)
	}
	return
}
