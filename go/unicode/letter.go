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
