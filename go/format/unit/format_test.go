package unit_test

import (
	"testing"

	"github.com/searKing/golang/go/format/unit"
)

type FormatFloatCaseTest struct {
	input      float64
	baseFormat unit.BaseFormat
	precision  int
	output     string
}

var (
	formatFloatCaseTests = []FormatFloatCaseTest{
		{
			input:      1234.567890,
			baseFormat: unit.BaseFormatBinary,
			precision:  1,
			output:     "1.2K",
		},
		{
			input:      2048.567890,
			baseFormat: unit.BaseFormatBinary,
			precision:  2,
			output:     "2K",
		},
		{
			input:      1999.567890,
			baseFormat: unit.BaseFormatBinary,
			precision:  2,
			output:     "1.95K",
		},
		{
			input:      1234.567890,
			baseFormat: unit.BaseFormatDecimal,
			precision:  1,
			output:     "1.2K",
		},
		{
			input:      2048.567890,
			baseFormat: unit.BaseFormatDecimal,
			precision:  2,
			output:     "2.05K",
		},
		{
			input:      1999.567890,
			baseFormat: unit.BaseFormatDecimal,
			precision:  2,
			output:     "2K",
		},
		{
			input:      123.45,
			baseFormat: unit.BaseFormat(1),
			precision:  2,
			output:     "123.45",
		},
		{
			input:      123.45,
			baseFormat: unit.BaseFormat(10),
			precision:  2,
			output:     "1.23M",
		},
	}
)

func TestFormatFloat(t *testing.T) {
	for n, test := range formatFloatCaseTests {
		if got := unit.FormatFloat(test.input, test.baseFormat, test.precision); got != test.output {
			t.Errorf("#%d: FormatFloat(%g,%d,%d) = %s, want %s", n, test.input, test.baseFormat, test.precision,
				got, test.output)
		}
	}
}
