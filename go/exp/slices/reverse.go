// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import "slices"

// Reverse reorder a slice of any ordered type in reverse order.
// Reverse modifies the contents of the slice s; it does not create a new slice.
// Deprecated: Use slices.Reverse instead since go1.21.
func Reverse[S ~[]E, E any](s S) {
	slices.Reverse(s)
}
