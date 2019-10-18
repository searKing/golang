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

type TrimBinaryMultiplePrefixCaseTest struct {
	input  string
	output string
}

var (
	trimBinaryMultiplePrefixCaseTests = []TrimBinaryMultiplePrefixCaseTest{
		{
			input:  "+1234.567890KiB",
			output: "B",
		},	{
			input:  "0xFFKiB",
			output: "B",
		},
	}
)

func TestTrimBinaryMultiplePrefix(t *testing.T) {
	for n, test := range trimBinaryMultiplePrefixCaseTests {
		if got := multiple_prefix.TrimBinaryMultiplePrefix(test.input); got != test.output {
			t.Errorf("#%d: TrimBinaryMultiplePrefix(%s) = %s, want %s", n, test.input,
				got, test.output)
		}
	}
}
