package strings

import (
	unicode_ "github.com/searKing/golang/go/unicode"
	"strings"
	"unicode"
	"unicode/utf8"
)

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]rune, 0, 32)
	sr, s_ := ExtractFirstRune(s)
	if sr == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		s = s_
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for s != "" {
		sr, s = ExtractFirstRune(s)
		if sr == '_' && s != "" {
			sr_, _ := ExtractFirstRune(s)
			if unicode_.IsASCIILower(sr_) {
				continue // Skip the underscore in s.
			}
		}
		if unicode_.IsASCIIDigit(sr) {
			t = append(t, sr)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if unicode_.IsASCIILower(sr) {
			sr = unicode.ToUpper(sr) // Make it a capital letter.
		}
		t = append(t, sr) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for s != "" {
			sr, s_ = ExtractFirstRune(s)
			if unicode_.IsASCIILower(sr) {
				s = s_
				t = append(t, sr)
				continue
			}
			break
		}
	}
	return string(t)
}

// SmallCamelCase returns the SmallCamelCased name.
// In short, _my_field_name_2 becomes xMyFieldName_2.
func SmallCamelCase(s string) string {
	s = CamelCase(s)
	sr, s := ExtractFirstRune(s)
	return string(unicode.ToLower(sr)) + s
}

func ExtractFirstRune(s string) (rune, string) {
	// Extract first rune from each string.
	var sr rune
	if s[0] < utf8.RuneSelf {
		sr, s = rune(s[0]), s[1:]
	} else {
		r, size := utf8.DecodeRuneInString(s)
		sr, s = r, s[size:]
	}
	return sr, s
}

// CamelCaseSlice is like CamelCase, but the argument is a slice of strings to
// be joined with "_".
func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

// SmallCamelCaseSlice is like SmallCamelCase, but the argument is a slice of strings to
// be joined with "_".
func SmallCamelCaseSlice(elem []string) string { return SmallCamelCase(strings.Join(elem, "_")) }

// DottedSlice turns a sliced name into a dotted name.
func DottedSlice(elem []string) string { return strings.Join(elem, ".") }
