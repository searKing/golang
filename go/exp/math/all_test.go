// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math_test

import (
	"testing"

	math_ "github.com/searKing/golang/go/exp/math"
)

var vf = []int{
	1,
	2,
	3,
}
var ceil = []int{
	1,
	4,
	6,
}
var floor = []int{
	1,
	-2,
	-3,
}

var fdim = []int{
	1,
	2,
	3,
}

func TestDim(t *testing.T) {
	for i := 0; i < len(vf); i++ {
		if f := math_.Dim(vf[i], 0); fdim[i] != f {
			t.Errorf("Dim(%d, %d) = %d, want %d", vf[i], 0, f, fdim[i])
		}
	}
}

func TestMax(t *testing.T) {
	for i := 0; i < len(vf); i++ {
		if f := math_.Max(vf[i], ceil[i]); ceil[i] != f {
			t.Errorf("Max(%d, %d) = %d, want %d", vf[i], ceil[i], f, ceil[i])
		}
	}
}

func TestMin(t *testing.T) {
	for i := 0; i < len(vf); i++ {
		if f := math_.Min(vf[i], floor[i]); floor[i] != f {
			t.Errorf("Min(%d, %d) = %d, want %d", vf[i], floor[i], f, floor[i])
		}
	}
}
