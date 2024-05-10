// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"bytes"
	"io"
)

type sniffReader struct {
	// sniff start: read from [historyBuffers..., source] and buffered in buffer
	// sniff stop: read from [buffer, historyBuffers..., source] and clean buffer and historyBuffers if meet EOF
	source io.Reader

	// virtual reader: buffer, historyBuffers..., source
	buffer         *bytes.Buffer // latest read data
	historyBuffers []io.Reader   // new, old, older read data

	selectorF DynamicReaderFunc

	sniffing bool
}

func newSniffReader(r io.Reader) *sniffReader {
	sr := &sniffReader{
		source: r,
	}
	sr.stopSniff()
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
		sr.startSniff()
		return sr
	}
	sr.stopSniff()
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
					sr.historyBuffers = sr.historyBuffers[1:] // remove head to recover space
				}
				return n, err
			}))
			sr.historyBuffers = append([]io.Reader{bufferReader}, sr.historyBuffers...)
		}
		sr.buffer = nil
	}
}

// startSniff starts sniff and return a TeeReader that writes to buffer while reads from history buffers and source
func (sr *sniffReader) startSniff() {
	sr.shrinkToHistory()
	// We don't need the buffer anymore.
	// Reset it to release the internal slice.
	sr.buffer = &bytes.Buffer{}

	readers := append(sr.historyBuffers, sr.source)
	reader := io.TeeReader(io.MultiReader(readers...), sr.buffer)
	sr.selectorF = func() io.Reader { return reader }
}

// stopSniff stops sniff and return a MultiReader of history buffers and source
func (sr *sniffReader) stopSniff() {
	sr.shrinkToHistory()
	readers := append(sr.historyBuffers, sr.source)
	reader := io.MultiReader(readers...)
	sr.selectorF = func() io.Reader { return reader }
}

func (sr *sniffReader) Read(p []byte) (n int, err error) { return sr.selectorF.Read(p) }
