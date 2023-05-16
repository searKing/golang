// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"golang.org/x/exp/constraints"
)

// RingRem returns the remainder of x looped by +ny until in y > 0 ? [0, y) : (y, 0].
// RingRem panics for y == 0 (division by zero).
// y > 0, then ∈ [0, y)
// y < 0, then ∈ (y, 0]
func RingRem[T constraints.Integer](x, y T) T {
	if y == 0 {
		panic("division by zero")
	}
	return ((x % y) + y) % y
}
