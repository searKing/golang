// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests for some of the internal functions.

package main

import (
	"testing"
)

type ParserTests struct {
	input  string
	output []typeInfo
}

var (
	parserTests = []ParserTests{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{"NumValue<**[]*[]*int>", []typeInfo{{
			Name:            "NumValue",
			Import:          "",
			valueType:       "int",
			valueImport:     "",
			valueIsPointer:  true,
			valueTypePrefix: "*[]*[]*",
		},
		}},
		{"NumValue<a.b>", []typeInfo{{
			Name:        "NumValue",
			Import:      "",
			valueType:   "a.b",
			valueImport: "a",
		}}},
		{"NumValue<a.b/c.d>", []typeInfo{{
			Name:        "NumValue",
			Import:      "",
			valueType:   "c.d",
			valueImport: "a.b/c",
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
