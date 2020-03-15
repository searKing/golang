// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiple_prefix_test

import (
	"testing"

	"github.com/searKing/golang/go/format/multiple_prefix"
)

type BinaryFormatFloatCaseTest struct {
	input     float64
	precision int
	output    string
}

var (
	binaryFormatFloatCaseTests = []BinaryFormatFloatCaseTest{
		{
			input:     1234.567890,
			precision: 1,
			output:    "1.2Ki",
		}, {
			input:     2000.567890,
			precision: 2,
			output:    "1.95Ki",
		}, {
			input:     1999.567890,
			precision: 4,
			output:    "1.9527Ki",
		}, {
			input:     1234.567890,
			precision: 1,
			output:    "1.2Ki",
		}, {
			input:     2048.567890,
			precision: 2,
			output:    "2Ki",
		}, {
			input:     1999.567890,
			precision: 2,
			output:    "1.95Ki",
		}, {
			input:     123.45,
			precision: 2,
			output:    "123.45",
		}, {
			input:     0.12345,
			precision: 2,
			output:    "0.12",
		}, {
			input:     -0.12345,
			precision: 2,
			output:    "-0.12",
		}, {
			input:     -0.00012345,
			precision: 5,
			output:    "-0.00012",
		}, {
			input:     -0.0001,
			precision: 2,
			output:    "-0",
		},
	}
)

func TestBinaryFormatFloat(t *testing.T) {
	for n, test := range binaryFormatFloatCaseTests {
		if got := multiple_prefix.BinaryFormatFloat(test.input, test.precision); got != test.output {
			t.Errorf("#%d: BinaryFormatFloat(%g,%d) = %s, want %s", n, test.input, test.precision,
				got, test.output)
		}
	}
}

type SplitBinaryCaseTest struct {
	input              string
	outputNumber       string
	outputPrefixSymbol string
	outputUnparsed     string
}

var (
	splitBinaryCaseTests = []SplitBinaryCaseTest{
		{
			input:              "1234.567890HelloWorld",
			outputNumber:       "1234.567890",
			outputPrefixSymbol: "",
			outputUnparsed:     "HelloWorld",
		}, {
			input:              "+1234.567890KiB",
			outputNumber:       "+1234.567890",
			outputPrefixSymbol: "Ki",
			outputUnparsed:     "B",
		}, {
			input:              "0xFFKiB",
			outputNumber:       "0xFF",
			outputPrefixSymbol: "Ki",
			outputUnparsed:     "B",
		}, {
			input:              "0xFFkiB",
			outputNumber:       "0xFF",
			outputPrefixSymbol: "",
			outputUnparsed:     "kiB",
		},
	}
)

func TestSplitBinary(t *testing.T) {
	for n, test := range splitBinaryCaseTests {
		gotNumber, gotPrefix, gotUnparsed := multiple_prefix.SplitBinary(test.input)
		if gotPrefix == nil {
			gotPrefix = multiple_prefix.BinaryMultiplePrefixTODO.Copy()
		}
		if gotNumber != test.outputNumber || gotPrefix.Symbol() != test.outputPrefixSymbol || gotUnparsed != test.outputUnparsed {
			t.Errorf("#%d: BinaryFormatFloat(%s) = (%s, %s, %s), want (%s, %s, %s)", n, test.input,
				gotNumber, gotPrefix.Symbol(), gotUnparsed,
				test.outputNumber, test.outputPrefixSymbol, test.outputUnparsed)
		}
	}
}
