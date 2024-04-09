// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bufio

import (
	"bytes"
	"strings"
	"testing"
)

var pairTests = []string{
	"",
	"{abc}",
	"{abc}{def}",
	"{abc{def}hij}",
	"sss{abc}",
	"{abc}eee",
	"sss{abc}eee",
	"{a[b]c}",
	"}{abc}",
}
var expectedPairTests = []string{
	"",
	"{abc}",
	"{abc}",
	"{abc{def}hij}",
	"{abc}",
	"{abc}",
	"{abc}",
	"{a[b]c}",
	"{abc}",
}

func TestPairScanner(t *testing.T) {
	for i, test := range pairTests {
		s := NewPairScanner(strings.NewReader(test)).SetDiscardLeading(true)
		p, err := s.ScanDelimiters("{}")
		if err != nil && i != 0 {
			t.Errorf("#%d: Scan error:%v\n", i, err)
			continue
		}

		if !bytes.Equal(p, []byte(expectedPairTests[i])) {
			t.Errorf("#%d: expected %q got %q", i, test, string(p))
		}
	}
}
