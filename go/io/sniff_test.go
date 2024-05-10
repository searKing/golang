// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io_test

import (
	"io"
	"strings"
	"testing"

	io_ "github.com/searKing/golang/go/io"
)

func shrinkSnifferMayBeWithHoleInHistoryBuffers(r io_.ReadSniffer, sniffing bool) {
	// stop sniffing and start sniffing, to move buffer to history buffers
	r.Sniff(!sniffing).Sniff(sniffing)
}

func TestSniffReaderSeeker(t *testing.T) {
	r := strings.NewReader("HEADER BODY TAILER")
	sniff := io_.SniffReader(r)

	// start sniffing
	sniff.Sniff(true)
	// ["HEADER BODY TAILER"]

	{ // sniff "HEADER"
		b := make([]byte, len("HEADER"))
		n, err := sniff.Read(b)
		if err != nil {
			t.Errorf("Error reading header: %v", err)
		}
		if string(b[:n]) != "HEADER" {
			t.Errorf("expected %q, got %q", "HEADER", b[:n])
		}
		shrinkSnifferMayBeWithHoleInHistoryBuffers(sniff, true)
		// ["HEADER BODY TAILER"]
	}

	{ // sniff "HEAD"
		b := make([]byte, len("HEAD"))
		n, err := sniff.Read(b)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
		}
		if string(b[:n]) != "HEAD" {
			t.Errorf("expected %q, got %q", " BODY", b[:n])
		}
		shrinkSnifferMayBeWithHoleInHistoryBuffers(sniff, true)
		// ["HEADER BODY TAILER"]
	}

	// stop sniffing
	sniff.Sniff(false)

	{ // sniff "HEADER"
		b := make([]byte, len("HEADER"))
		n, err := sniff.Read(b)
		if err != nil {
			t.Errorf("Error reading header: %v", err)
		}
		if string(b[:n]) != "HEADER" {
			t.Errorf("expected %q, got %q", "HEAD", b[:n])
		}
	}
	{
		b, err := io.ReadAll(sniff)
		if err != nil {
			t.Errorf("Error reading all: %v", err)
		}
		if string(b) != " BODY TAILER" {
			t.Errorf("expected %q, got %q", " BODY TAIL", b)
		}
	}
}

type reader struct {
	r io.Reader
}

func (r *reader) Read(b []byte) (int, error) {
	return r.r.Read(b)
}

func TestSniffReaderNotSeeker(t *testing.T) {
	r := &reader{r: strings.NewReader("HEADER BODY TAILER")}
	sniff := io_.SniffReader(r)

	// start sniffing
	sniff.Sniff(true)
	// ["HEADER BODY TAILER"]

	{ // sniff "HEADER"
		b := make([]byte, len("HEADER"))
		n, err := sniff.Read(b)
		if err != nil {
			t.Errorf("Error reading header: %v", err)
		}
		if string(b[:n]) != "HEADER" {
			t.Errorf("expected %q, got %q", "HEADER", b[:n])
		}
		shrinkSnifferMayBeWithHoleInHistoryBuffers(sniff, true)
		// ["HEADER", " BODY TAILER"]
	}

	{ // sniff "HEAD"
		b := make([]byte, len("HEAD"))
		n, err := sniff.Read(b)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
		}
		if string(b[:n]) != "HEAD" {
			t.Errorf("expected %q, got %q", " BODY", b[:n])
		}
		shrinkSnifferMayBeWithHoleInHistoryBuffers(sniff, true)
		// ["HEAD", "ER", " BODY TAILER"]
	}

	// stop sniffing
	sniff.Sniff(false)

	{ // sniff "HEADER"
		b := make([]byte, len("HEADER"))
		n, err := sniff.Read(b)
		if err != nil {
			t.Errorf("Error reading header: %v", err)
		}
		if string(b[:n]) != "HEAD" {
			t.Errorf("expected %q, got %q", "HEAD", b[:n])
		}
		n, err = sniff.Read(b)
		if err != nil {
			t.Errorf("Error reading header: %v", err)
		}
		if string(b[:n]) != "ER" {
			t.Errorf("expected %q, got %q", "ER", b[:n])
		}
	}
	{
		b, err := io.ReadAll(sniff)
		if err != nil {
			t.Errorf("Error reading all: %v", err)
		}
		if string(b) != " BODY TAILER" {
			t.Errorf("expected %q, got %q", " BODY TAIL", b)
		}
	}
}
