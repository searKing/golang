// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math_test

import (
	"testing"

	"github.com/searKing/golang/go/math"
)

var vf = []int64{
	5,
	7,
	-3,
	-5,
	10,
	-8,
	0,
}

var fabs = []int64{
	5,
	7,
	3,
	5,
	10,
	8,
	0,
}

func TestAbs(t *testing.T) {
	for i := 0; i < len(vf); i++ {
		if f := math.AbsInt64(vf[i]); fabs[i] != f {
			t.Errorf("AbsInt64(%d) = %d, want %d", vf[i], f, fabs[i])
		}
	}
}
