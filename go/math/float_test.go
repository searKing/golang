// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math_test

import (
	math_ "math"
	"testing"

	"github.com/searKing/golang/go/math"
)

type TruncPrecisionCaseTest struct {
	input  float64
	n      int
	output float64
}

var (
	truncPrecisionCaseTests = []TruncPrecisionCaseTest{
		{
			input:  -100.00001,
			n:      2,
			output: -100,
		}, {
			input:  1254.567890,
			n:      -2,
			output: 1300,
		}, {
			input:  1234.567890,
			n:      -1,
			output: 1230,
		}, {
			input:  1234.567890,
			n:      0,
			output: 1235,
		}, {
			input:  1234.567890,
			n:      1,
			output: 1234.6,
		}, {
			input:  1234.567890,
			n:      6,
			output: 1234.56789,
		}, {
			input:  1234.567890,
			n:      10,
			output: 1234.56789,
		}, {
			input:  math_.Inf(-1),
			n:      1,
			output: math_.Inf(-1),
		}, {
			input:  math_.Copysign(0, -1),
			n:      1,
			output: math_.Copysign(0, -1),
		}, {
			input:  0,
			n:      1,
			output: 0,
		}, {
			input:  math_.Inf(1),
			n:      1,
			output: math_.Inf(1),
		}, {
			input:  math_.NaN(),
			n:      1,
			output: math_.NaN(),
		},
	}
)

func TestTruncPrecision(t *testing.T) {
	for n, test := range truncPrecisionCaseTests {
		if got := math.TruncPrecision(test.input, test.n); !math.Alike(got, test.output) {
			t.Errorf("#%d: TruncPrecision(%g,%d) = %g, want %g", n, test.input, test.n, got, test.output)
		}
	}
}
