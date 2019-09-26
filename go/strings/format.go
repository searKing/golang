package strings

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// CamelCase returns the CamelCased name.
// If there is an interior split rune such as an underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string, seps ...rune) string {
	return TransformCase(s, func(r string) string {
		r = withDefault(leadingString, strings.ToLower)(r)
		return ToUpperLeading(r)
	}, seps...)
}

// SmallCamelCase returns the SmallCamelCased name.
// In short, _my_field_name_2 becomes xMyFieldName_2.
func SmallCamelCase(s string, seps ...rune) string {
	s = CamelCase(s, seps...)
	sr, s := ExtractFirstRune(s)
	return ToLowerLeading(string(sr)) + s
}

// SnakeCase returns the SnakeCased name.
// In short, _my_field_name_2 becomes x_my_field_name_2.
func SnakeCase(s string, seps ...rune) string {
	return TransformCase(s, MapAndJoin("_", withDefault(leadingString, strings.ToLower)), seps...)
}

// KebabCase returns the KebabCased name.
// In short, _my_field_name_2 becomes x-my-field-name-2.
func KebabCase(s string, seps ...rune) string {
	return TransformCase(s, MapAndJoin("-", withDefault(leadingString, strings.ToLower)), seps...)
}

// DotCase returns the KebabCased name.
// In short, _my_field_name_2 becomes x.my.field.name.2.
func DotCase(s string, seps ...rune) string {
	return TransformCase(s, MapAndJoin(".", withDefault(leadingString, strings.ToLower)), seps...)
}

func MapAndJoin(sep string, mapping func(r string) string) func(r string) string {
	var written bool
	return func(r string) string {
		r = mapping(r)
		if written {
			r = sep + r
		}
		written = true

		return r
	}
}

func ExtractFirstRune(s string) (rune, string) {
	if s == "" {
		return -1, s
	}
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
func CamelCaseSlice(elem ...string) string { return CamelCase(strings.Join(elem, "_"), '_') }

// SmallCamelCaseSlice is like SmallCamelCase, but the argument is a slice of strings to
// be joined with "_".
func SmallCamelCaseSlice(elem ...string) string { return SmallCamelCase(strings.Join(elem, "_"), '_') }

// DottedSlice turns a sliced name into a dotted name.
func DottedSlice(elem ...string) string { return strings.Join(elem, ".") }

const (
	leadingString = "X"
)

// TransformCase Splits and apply map on every splits
func TransformCase(s string, f func(r string) string, seps ...rune) string {
	var out strings.Builder
	for _, sub := range splits(s, seps...) {
		out.WriteString(f(sub))
	}
	return out.String()
}

// split s into sub strings
// meet seps, split
// meet Capital rune, split
// if s is leading with seps, filled with "" at first
func splits(s string, seps ...rune) []string {
	if s == "" {
		return nil
	}
	var splited []string
	sr, s_ := ExtractFirstRune(s)
	if strings.ContainsRune(string(seps), sr) {
		// Need a non sep letter; drop the split rune, such as '_'.
		splited = append(splited, "")
		s = s_
	}

	var ele strings.Builder
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for s != "" {
		ele.Reset()
		sr, s = ExtractFirstRune(s)
		if strings.ContainsRune(string(seps), sr) && s != "" {
			for s != "" {
				sr, s = ExtractFirstRune(s)
				if strings.ContainsRune(string(seps), sr) {
					// ignore seps that follows
					continue
				}
				break
			}
			// EOF
			if strings.ContainsRune(string(seps), sr) {
				break
			}
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		ele.WriteRune(sr)
		// Accept lower case sequence that follows.
		for s != "" {
			sr, s_ = ExtractFirstRune(s)
			if !strings.ContainsRune(string(seps), sr) && unicode.IsLower(sr) {
				s = s_
				ele.WriteRune(sr)
				continue
			}
			break
		}
		splited = append(splited, ele.String())
	}
	return splited
}

func withDefault(def string, f func(s string) string) func(s string) string {
	return func(s string) string {
		if s == "" {
			s = def
		}
		return f(s)
	}
}
