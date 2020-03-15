// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytes

import "bytes"

func NewIndent(dst *bytes.Buffer, prefix, indent string, depth int) {
	dst.WriteString(prefix)
	for i := 0; i < depth; i++ {
		dst.WriteString(indent)
	}
}

func NewLine(dst *bytes.Buffer, prefix, indent string, depth int) {
	dst.WriteByte('\n')
	NewIndent(dst, prefix, indent, depth)
}
