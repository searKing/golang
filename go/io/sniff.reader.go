// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"bytes"
	"io"
)

type sniffReader struct {
	source         io.Reader
	buffer         *bytes.Buffer
	historyBuffers []io.Reader

	selectorF DynamicReaderFunc

	sniffing bool
}

func newSniffReader(r io.Reader) *sniffReader {
	sr := &sniffReader{
		source: r,
	}
	sr.resetSelector()
	return sr
}

// Sniff starts or stops sniffing, restarts if stop and start called one by one
// true to start sniffing all data unread actually
// false to return a multi reader with all data sniff buffered and source
func (sr *sniffReader) Sniff(sniffing bool) ReadSniffer {
	if sr.sniffing == sniffing {
		return sr
	}
	sr.sniffing = sniffing
	if sniffing {
		sr.shrinkToHistory()
		// We don't need the buffer anymore.
		// Reset it to release the internal slice.
		sr.buffer = &bytes.Buffer{}

		readers := sr.historyBuffers
		readers = append(readers, sr.source)
		reader := io.TeeReader(io.MultiReader(readers...), sr.buffer)
		sr.selectorF = func() io.Reader {
			return reader
		}
		return sr
	}
	sr.resetSelector()
	return sr
}

// shrinkToHistory shrink buffer to history buffers
func (sr *sniffReader) shrinkToHistory() {
	if sr.buffer != nil {
		if sr.buffer.Len() > 0 {
			// clear if EOF meet
			bufferReader := WatchReader(bytes.NewBuffer(sr.buffer.Bytes()), WatcherFunc(func(p []byte, n int, err error) (int, error) {
				if err == io.EOF {
					// historyBuffers is consumed head first, so can be cleared from head
					sr.historyBuffers = sr.historyBuffers[1:] // recycle memory
				}
				return n, err
			}))
			var rs []io.Reader
			rs = append(rs, bufferReader)
			sr.historyBuffers = append(rs, sr.historyBuffers...)
		}
		sr.buffer = nil
	}
}

// resetSelector stops sniff and return a MultiReader of history buffers and source
func (sr *sniffReader) resetSelector() {
	sr.shrinkToHistory()
	readers := append(sr.historyBuffers, sr.source)
	reader := io.MultiReader(readers...)
	sr.selectorF = func() io.Reader {
		return reader
	}
}

func (sr *sniffReader) Read(p []byte) (n int, err error) {
	return sr.selectorF.Read(p)
}
