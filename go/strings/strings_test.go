package strings_test

import (
	"github.com/searKing/golang/go/strings"
	"testing"
)

type SliceContainsTest struct {
	inputSS []string
	inputS  string
	output  bool
}

var (
	sliceContainsTests = []SliceContainsTest{
		{
			[]string{"A", "B", "C", "D"},
			"A",
			true,
		},
		{
			[]string{"A", "B", "C", "D"},
			"E",
			false,
		},
	}
)

func TestSliceContains(t *testing.T) {
	for n, test := range sliceContainsTests {
		out := strings.SliceContains(test.inputSS, test.inputS)
		if out != test.output {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}
