// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests for some of the internal functions.

package enum

import (
	"testing"
)

type ParserTest struct {
	input  string
	output []typeInfo
}

var (
	parserTests = []ParserTest{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{"NumMap", []typeInfo{{
			Name:   "NumMap",
			Import: "",
		}}},
	}
)

func TestParserTests(t *testing.T) {
Outer:
	for n, test := range parserTests {
		runs := newTypeInfo(test.input)
		if len(runs) != len(test.output) {
			t.Errorf("#%d: %v: got %d runs; expected %d", n, test.input, len(runs), len(test.output))
			continue
		}
		for i, run := range runs {
			if run != test.output[i] {
				t.Errorf("#%d: got %v; expected %v", n, runs, test.output)
				continue Outer
			}
		}
	}
}
