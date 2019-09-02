package ast_test

import (
	"testing"
)

type TokenizerTest struct {
	input  []rune
	output []Token
}

var (
	tokenizerTests = []TokenizerTest{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{[]rune("NumValue<int, *time.Time>"), []Token{{
			Type:  TokenTypeName,
			Value: "NumValue",
		}, {
			Type:  TokenTypeParen,
			Value: "<",
		}, {
			Type:  TokenTypeName,
			Value: "int",
		}, {
			Type:  TokenTypeParen,
			Value: ",",
		}, {
			Type:  TokenTypeParen,
			Value: "*",
		}, {
			Type:  TokenTypeName,
			Value: "time.Time",
		}, {
			Type:  TokenTypeParen,
			Value: ">",
		}}},
		{[]rune("NumValue<int, string>, AnotherNumValue<int, interface{}>"), []Token{{
			Type:  TokenTypeName,
			Value: "NumValue",
		}, {
			Type:  TokenTypeParen,
			Value: "<",
		}, {
			Type:  TokenTypeName,
			Value: "int",
		}, {
			Type:  TokenTypeParen,
			Value: ",",
		}, {
			Type:  TokenTypeName,
			Value: "string",
		}, {
			Type:  TokenTypeParen,
			Value: ">",
		}, {
			Type:  TokenTypeParen,
			Value: ",",
		}, {
			Type:  TokenTypeName,
			Value: "AnotherNumValue",
		}, {
			Type:  TokenTypeParen,
			Value: "<",
		}, {
			Type:  TokenTypeName,
			Value: "int",
		}, {
			Type:  TokenTypeParen,
			Value: ",",
		}, {
			Type:  TokenTypeName,
			Value: "interface{}",
		}, {
			Type:  TokenTypeParen,
			Value: ">",
		}}},
	}
)

func TestTokenizer(t *testing.T) {
Outer:
	for n, test := range tokenizerTests {
		runs := Tokenizer(test.input)
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
