// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import (
	"testing"

	"github.com/searKing/golang/tools/common/ast"
)

type TokenizerTest struct {
	input  []rune
	output []ast.Token
}

var (
	tokenizerTests = []TokenizerTest{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{[]rune("NumValue<int, *encoding/json.Time>"), []ast.Token{{
			Type:  ast.TokenTypeName,
			Value: "NumValue",
		}, {
			Type:  ast.TokenTypeParen,
			Value: "<",
		}, {
			Type:  ast.TokenTypeName,
			Value: "int",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ",",
		}, {
			Type:  ast.TokenTypeParen,
			Value: "*",
		}, {
			Type:  ast.TokenTypeName,
			Value: "encoding/json.Time",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ">",
		}}},
		{[]rune("NumValue<int, string>, AnotherNumValue<[]int, interface{}>"), []ast.Token{{
			Type:  ast.TokenTypeName,
			Value: "NumValue",
		}, {
			Type:  ast.TokenTypeParen,
			Value: "<",
		}, {
			Type:  ast.TokenTypeName,
			Value: "int",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ",",
		}, {
			Type:  ast.TokenTypeName,
			Value: "string",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ">",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ",",
		}, {
			Type:  ast.TokenTypeName,
			Value: "AnotherNumValue",
		}, {
			Type:  ast.TokenTypeParen,
			Value: "<",
		}, {
			Type:  ast.TokenTypeParen,
			Value: "[",
		}, {
			Type:  ast.TokenTypeParen,
			Value: "]",
		}, {
			Type:  ast.TokenTypeName,
			Value: "int",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ",",
		}, {
			Type:  ast.TokenTypeName,
			Value: "interface{}",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ">",
		}}},
	}
)

func TestTokenizer(t *testing.T) {
Outer:
	for n, test := range tokenizerTests {
		runs := ast.Tokenizer(test.input)
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
