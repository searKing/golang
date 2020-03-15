// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

// AbsInt64 returns the absolute value of x.
func AbsInt64(x int64) int64 {
	y := x >> 63       // y <- x>> 63
	return (x ^ y) - y // (x XOR y) - y
}
