// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson

// WithEncOptsTruncate sets truncate in encOpts.
func WithEncOptsTruncate(v int) EncOptsOption {
	return EncOptsOptionFunc(func(o *encOpts) {
		o.truncateBytes = v
		o.truncateString = v
		o.truncateMap = v
		o.truncateSliceOrArray = v
	})
}
