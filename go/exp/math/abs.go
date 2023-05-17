// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import "github.com/searKing/golang/go/exp/constraints"

// Abs returns the absolute value of last.
func Abs[T constraints.Number](x T) T {
	if x >= 0 {
		return x
	}
	return -x
}
