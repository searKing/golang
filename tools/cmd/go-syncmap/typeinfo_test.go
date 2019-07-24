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
		{[]rune("NumMap<int, string>"), []_token{{
			typ:   tokenTypeName,
			value: "NumMap",
		}, {
			typ:   tokenTypeParen,
			value: "<",
		}, {
			typ:   tokenTypeName,
			value: "int",
		}, {
			typ:   tokenTypeParen,
			value: ",",
		}, {
			typ:   tokenTypeName,
			value: "string",
		}, {
			typ:   tokenTypeParen,
			value: ">",
		}}},
		{[]rune("NumMap<int, string>, AnotherNumMap<int, Time>"), []_token{{
			typ:   tokenTypeName,
			value: "NumMap",
		}, {
			typ:   tokenTypeParen,
			value: "<",
		}, {
			typ:   tokenTypeName,
			value: "int",
		}, {
			typ:   tokenTypeParen,
			value: ",",
		}, {
			typ:   tokenTypeName,
			value: "string",
		}, {
			typ:   tokenTypeParen,
			value: ">",
		}, {
			typ:   tokenTypeParen,
			value: ",",
		}, {
			typ:   tokenTypeName,
			value: "AnotherNumMap",
		}, {
			typ:   tokenTypeParen,
			value: "<",
		}, {
			typ:   tokenTypeName,
			value: "int",
		}, {
			typ:   tokenTypeParen,
			value: ",",
		}, {
			typ:   tokenTypeName,
			value: "Time",
		}, {
			typ:   tokenTypeParen,
			value: ">",
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
			value: "NumMap",
		}, {
			typ:   tokenTypeParen,
			value: "<",
		}, {
			typ:   tokenTypeName,
			value: "int",
		}, {
			typ:   tokenTypeParen,
			value: ",",
		}, {
			typ:   tokenTypeName,
			value: "string",
		}, {
			typ:   tokenTypeParen,
			value: ">",
		}}, []typeInfo{{
			mapName:     "NumMap",
			mapImport:   "",
			keyType:     "int",
			keyImport:   "",
			valueType:   "string",
			valueImport: "",
		},
		}},
		{[]_token{{
			typ:   tokenTypeName,
			value: "NumMap",
		}, {
			typ:   tokenTypeParen,
			value: "<",
		}, {
			typ:   tokenTypeName,
			value: "a.b",
		}, {
			typ:   tokenTypeParen,
			value: ",",
		}, {
			typ:   tokenTypeName,
			value: "a.b.c",
		}, {
			typ:   tokenTypeParen,
			value: ">",
		}}, []typeInfo{{
			mapName:     "NumMap",
			mapImport:   "",
			keyType:     "a.b",
			keyImport:   "a",
			valueType:   "b.c",
			valueImport: "a.b",
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
