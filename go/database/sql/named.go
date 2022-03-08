// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import "strings"

// CompliantName returns a compliant id name
// that can be used for a bind or as var.
// replace special runes with '_'
// a.b -> a_b
func CompliantName(in string) string {
	var buf strings.Builder
	for i, c := range in {
		if !isLetter(uint16(c)) {
			if i == 0 || !isDigit(uint16(c)) {
				buf.WriteByte('_')
				continue
			}
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

func isDigit(ch uint16) bool {
	return '0' <= ch && ch <= '9'
}
func isLetter(ch uint16) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}
