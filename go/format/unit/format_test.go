package unit_test

import (
	"testing"

	"github.com/searKing/golang/go/format/unit"
)

type FormatFloatCaseTest struct {
	input     float64
	precision int
	output    string
}

var (
	formatFloatCaseTests = []FormatFloatCaseTest{
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

func TestFormatFloat(t *testing.T) {
	for n, test := range formatFloatCaseTests {
		if got := unit.FormatFloat(test.input, test.precision); got != test.output {
			t.Errorf("#%d: FormatFloat(%g,%d) = %s, want %s", n, test.input, test.precision,
				got, test.output)
		}
	}
}
