// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"bytes"
	"io"
)

// ReadSniffer is the interface that groups the basic Read and Sniff methods.
type ReadSniffer interface {
	io.Reader
	Sniff(sniffing bool)
}

type sniffReader struct {
	source io.Reader
	buffer bytes.Buffer

	selectorF DynamicReaderFunc

	sniffing bool
}

func (sr *sniffReader) Sniff(sniffing bool) {
	if sr.sniffing == sniffing {
		return
	}
	sr.sniffing = sniffing
	if sniffing {
		// We don't need the buffer anymore.
		// Reset it to release the internal slice.
		sr.buffer = bytes.Buffer{}
		sr.selectorF = func() io.Reader {
			return io.TeeReader(sr.source, &sr.buffer)
		}
		return
	}
	sr.resetSelector()
}

func (sr *sniffReader) resetSelector() {
	sr.selectorF = func() io.Reader {
		// clear if EOF meet
		bufferReader := WatchReader(&sr.buffer, WatcherFunc(func(p []byte, n int, err error) (int, error) {
			if err == io.EOF {
				sr.buffer = bytes.Buffer{} // recycle memory
			}
			return n, err
		}))

		return io.MultiReader(bufferReader, sr.source)
	}
}

func (sr *sniffReader) Read(p []byte) (n int, err error) {
	return sr.selectorF.Read(p)
}

// SniffReader returns a Reader that allows sniff and read from
// the provided input reader.
// data is buffered if Sniff(true) is called.
// buffered data is taken first, if Sniff(false) is called.
func SniffReader(r io.Reader) ReadSniffer {
	sr := &sniffReader{
		source: r,
	}
	sr.resetSelector()
	return sr
}
