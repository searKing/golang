// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiple_prefix_test

import (
	"testing"

	"github.com/searKing/golang/go/format/multiple_prefix"
)

type DecimalFormatFloatCaseTest struct {
	input     float64
	precision int
	output    string
}

var (
	decimalFormatFloatCaseTests = []DecimalFormatFloatCaseTest{
		{
			input:     1234.567890,
			precision: 1,
			output:    "1.2k",
		}, {
			input:     2000.567890,
			precision: 2,
			output:    "2k",
		}, {
			input:     1999.567890,
			precision: 4,
			output:    "1.9996k",
		}, {
			input:     1234.567890,
			precision: 1,
			output:    "1.2k",
		}, {
			input:     2048.567890,
			precision: 2,
			output:    "2.05k",
		}, {
			input:     1999.567890,
			precision: 2,
			output:    "2k",
		}, {
			input:     123.45,
			precision: 2,
			output:    "123.45",
		}, {
			input:     0.12345,
			precision: 2,
			output:    "123.45m",
		}, {
			input:     -0.12345,
			precision: 2,
			output:    "-123.45m",
		}, {
			input:     -0.00012345,
			precision: 2,
			output:    "-123.45μ",
		}, {
			input:     -0.0001,
			precision: 2,
			output:    "-100μ",
		},
	}
)

func TestDecimalFormatFloat(t *testing.T) {
	for n, test := range decimalFormatFloatCaseTests {
		if got := multiple_prefix.DecimalFormatFloat(test.input, test.precision); got != test.output {
			t.Errorf("#%d: DecimalFormatFloat(%g,%d) = %s, want %s", n, test.input, test.precision,
				got, test.output)
		}
	}
}

type SplitDecimalCaseTest struct {
	input              string
	outputNumber       string
	outputPrefixSymbol string
	outputUnparsed     string
}

var (
	splitDecimalCaseTests = []SplitDecimalCaseTest{
		{
			input:              "1234.567890HelloWorld",
			outputNumber:       "1234.567890",
			outputPrefixSymbol: "",
			outputUnparsed:     "HelloWorld",
		}, {
			input:              "+1234.567890\tkB",
			outputNumber:       "+1234.567890",
			outputPrefixSymbol: "k",
			outputUnparsed:     "B",
		}, {
			input:              "0xFFkB",
			outputNumber:       "0xFF",
			outputPrefixSymbol: "k",
			outputUnparsed:     "B",
		}, {
			input:              "0xFFKB",
			outputNumber:       "0xFF",
			outputPrefixSymbol: "",
			outputUnparsed:     "KB",
		},
	}
)

func TestSplitDecimal(t *testing.T) {
	for n, test := range splitDecimalCaseTests {
		gotNumber, gotPrefix, gotUnparsed := multiple_prefix.SplitDecimal(test.input)
		if gotPrefix == nil {
			gotPrefix = multiple_prefix.DecimalMultiplePrefixTODO.Copy()
		}
		if gotNumber != test.outputNumber || gotPrefix.Symbol() != test.outputPrefixSymbol || gotUnparsed != test.outputUnparsed {
			t.Errorf("#%d: DecimalFormatFloat(%s) = (%s, %s, %s), want (%s, %s, %s)", n, test.input,
				gotNumber, gotPrefix.Symbol(), gotUnparsed,
				test.outputNumber, test.outputPrefixSymbol, test.outputUnparsed)
		}
	}
}
