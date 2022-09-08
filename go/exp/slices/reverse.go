// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// Reverse reorder a slice of any ordered type in reverse order.
// Reverse modifies the contents of the slice s; it does not create a new slice.
func Reverse[S ~[]E, E any](x S) {
	for i := 0; i < len(x)>>1; i++ {
		t := x[i]
		x[i] = x[len(x)-1-i]
		x[len(x)-1-i] = t
	}
}
