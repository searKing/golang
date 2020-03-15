// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

// take in a character set and a length and will generate a random string using that character set.
func StringMathWithCharset(length int, charset string) string {

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRandMath.Intn(len(charset))]
	}
	return string(b)
}

// only take in a length, and will use a default characters set to generate a random string
func StringMath(length int) string {
	return StringMathWithCharset(length, CharsetAlphaNum)
}
