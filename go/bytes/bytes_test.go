// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytes_test

import (
	"testing"

	bytes_ "github.com/searKing/golang/go/bytes"
)

func TestPadLeft(t *testing.T) {
	table := []struct {
		Q   string
		pad string
		n   int
		R   string
	}{
		{
			Q:   "a",
			pad: "*",
			n:   -1,
			R:   "a",
		},
		{
			Q:   "a",
			pad: "*",
			n:   10,
			R:   "*********a",
		},
		{
			Q:   "a",
			pad: "*^",
			n:   5,
			R:   "*^*^a",
		},
		{
			Q:   "a",
			pad: "*^",
			n:   6,
			R:   "*^*^ a",
		},
	}

	for i, test := range table {
		qr := bytes_.PadLeft([]byte(test.Q), []byte(test.pad), test.n)
		if string(qr) != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestPadRight(t *testing.T) {
	table := []struct {
		Q   string
		pad string
		n   int
		R   string
	}{
		{
			Q:   "a",
			pad: "*",
			n:   -1,
			R:   "a",
		},
		{
			Q:   "a",
			pad: "*",
			n:   1,
			R:   "a",
		},
		{
			Q:   "a",
			pad: "*",
			n:   10,
			R:   "a*********",
		},
		{
			Q:   "a",
			pad: "*^",
			n:   5,
			R:   "a*^*^",
		},
		{
			Q:   "a",
			pad: "*^",
			n:   6,
			R:   "a *^*^",
		},
	}

	for i, test := range table {
		qr := bytes_.PadRight([]byte(test.Q), []byte(test.pad), test.n)
		if string(qr) != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestReverse(t *testing.T) {
	table := []struct {
		Q string
		R string
	}{
		{
			Q: "abc123",
			R: "321cba",
		},
		{
			Q: "Hello, 世界",
			R: "\x8c\x95疸\xe4 ,olleH",
		},
	}

	for i, test := range table {
		qr := bytes_.Reverse([]byte(test.Q))
		if string(qr) != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}
