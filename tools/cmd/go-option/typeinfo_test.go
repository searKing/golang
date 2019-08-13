// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests for some of the internal functions.

package main

import (
	"testing"
)

type TokenizerTests struct {
	input  []rune
	output []_token
}

var (
	tokenizerTests = []TokenizerTests{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{[]rune("NumValue"), []_token{{
			typ:   tokenTypeName,
			value: "NumValue",
		}}},
		{[]rune("NumValue, AnotherNumValue"), []_token{{
			typ:   tokenTypeName,
			value: "NumValue",
		}, {
			typ:   tokenTypeParen,
			value: ",",
		}, {
			typ:   tokenTypeName,
			value: "AnotherNumValue",
		}}},
	}
)

func TestTokenizers(t *testing.T) {
Outer:
	for n, test := range tokenizerTests {
		runs := tokenizer(test.input)
		if len(runs) != len(test.output) {
			t.Errorf("#%d: %v: got %d runs; expected %d", n, string(test.input), len(runs), len(test.output))
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

type ParserTests struct {
	input  []_token
	output []typeInfo
}

var (
	parserTests = []ParserTests{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{[]_token{{
			typ:   tokenTypeName,
			value: "NumValue",
		}}, []typeInfo{{
			eleName:   "NumValue",
			eleImport: "",
		},
		}},
		{[]_token{{
			typ:   tokenTypeName,
			value: "NumValue",
		}, {
			typ:   tokenTypeParen,
			value: ",",
		}, {
			typ:   tokenTypeName,
			value: "AnotherNumValue",
		}}, []typeInfo{{
			eleName:   "NumValue",
			eleImport: "",
		}, {
			eleName:   "AnotherNumValue",
			eleImport: "",
		}}},
	}
)

func TestParserTests(t *testing.T) {
Outer:
	for n, test := range parserTests {
		runs := parser(test.input)
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
