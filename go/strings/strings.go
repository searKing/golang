package strings

import "unicode"

// SliceContains  reports whether s is within ss.
func SliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// MapLeading returns a copy of the string s with its first characters modified
// according to the mapping function. If mapping returns a negative value, the character is
// dropped from the string with no replacement.
func MapLeading(mapping func(rune) rune, s string) string {
	if s == "" {
		return s
	}
	rLeading, sRight := ExtractFirstRune(s)
	srMapped := mapping(rLeading)
	if srMapped < 0 {
		return sRight
	}

	// Fast path for unchanged input
	if rLeading == srMapped {
		return s
	}
	return string(srMapped) + sRight
}

// ToLowerLeading returns s with it's first Unicode letter mapped to their lower case.
func ToLowerLeading(s string) string {
	return MapLeading(unicode.ToLower, s)
}

// ToUpperLeading returns s with it's first Unicode letter mapped to their upper case.
func ToUpperLeading(s string) string {
	return MapLeading(unicode.ToUpper, s)
}
