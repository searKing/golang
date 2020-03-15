// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unicode

import "unicode"

// IsASCII reports whether the rune is a ASCII character.
func IsASCII(s rune) bool {
	return s < unicode.MaxASCII
}

// IsLatin1 reports whether the rune is a latin1 (ISO-8859-1) character.
func IsLatin1(s rune) bool {
	return s < unicode.MaxLatin1
}

// IsASCIIUpper reports whether the rune is an ASCII and upper case letter.
func IsASCIIUpper(r rune) bool {
	return IsASCII(r) && unicode.IsUpper(r)
}

// IsASCIILower reports whether the rune is an ASCII and lower case letter.
func IsASCIILower(r rune) bool {
	return IsASCII(r) && unicode.IsLower(r)
}

// IsASCIIDigit reports whether the rune is an ASCII and decimal digit.
func IsASCIIDigit(r rune) bool {
	return IsASCII(r) && unicode.IsDigit(r)
}

// IsVowel reports whether the rune is an ASCII and vowel case letter.
func IsVowel(s rune) bool {
	switch unicode.ToUpper(s) {
	case 'A', 'E', 'I', 'O', 'U':
		return true
	default:
		return false
	}
}

// IsConsonant reports whether the rune is an ASCII and consonant case letter.
func IsConsonant(s rune) bool {
	switch unicode.ToUpper(s) {
	case 'B', 'C', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'V', 'W', 'X', 'Y', 'Z':
		return true
	default:
		return false
	}
}
