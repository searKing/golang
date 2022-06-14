// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps

// Set implements a non-thread safe Set
func Set[M map[K]struct{}, K comparable](ks ...K) (m M) {
	m = make(M, len(ks))
	for _, k := range ks {
		m[k] = struct{}{}
	}
	return m
}
