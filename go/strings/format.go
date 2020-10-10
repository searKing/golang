// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"strings"
	"unicode"
	"unicode/utf8"

	unicode_ "github.com/searKing/golang/go/unicode"
)

// see https://en.wikipedia.org/wiki/Camel_case
// Camel case (stylized as camelCase; also known as camel caps or more formally as medial capitals) is the practice of
// writing phrases such that each word or abbreviation in the middle of the phrase begins with a capital letter,
// with no intervening spaces or punctuation.
// Common examples include "iPhone" and "eBay".

var (
	PascalCase           = UpperCamelCase
	CapitalizedWordsCase = UpperCamelCase
	CapWordsCase         = UpperCamelCase
	CapitalizedWords     = UpperCamelCase

	// SentenceCase is a mixed-case style in which the first word of the sentence is capitalised,
	// as well as proper nouns and other words as required by a more specific rule.
	// This is generally equivalent to the baseline universal standard of formal English orthography.
	// https://en.wikipedia.org/wiki/Letter_case#Sentence_Case
	// "The quick brown fox jumps over the lazy dog"
	SentenceCase = strings.Title

	// TitleCase capitalises all words but retains the spaces between them
	// https://en.wikipedia.org/wiki/Letter_case#Title_Case
	// "The Quick Brown Fox Jumps over the Lazy Dog"
	TitleCase = strings.ToTitle

	// AllCapsCase is an unicase style with capital letters only.
	// https://en.wikipedia.org/wiki/Letter_case#All_caps
	// "THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG"
	AllCapsCase  = strings.ToUpper
	AllUpperCase = AllCapsCase

	// AllLowercase is an unicase style with no capital letters.
	// https://en.wikipedia.org/wiki/Letter_case#All_lowercase
	// "the quick brown fox jumps over the lazy dog"
	AllLowercase = strings.ToLower
)

var (
	DromedaryCase = LowerCamelCase
	// Some people and organizations, notably Microsoft, use the term camel case only for lower camel case.
	// Pascal case means only upper camel case.
	CamelCase = LowerCamelCase

	// MixedCase for lower camel case in Python
	MixedCase = LowerCamelCase
)

var (
	// lowercase
	LowerCase = strings.ToLower
)

// UpperCamelCase returns the CamelCased name by initial uppercase letter.
// If there is an interior split rune such as an underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ JoinGenerator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
// "TheQuickBrownFoxJumpsOverTheLazyDog"
func UpperCamelCase(s string, seps ...rune) string {
	return TransformCase(s, func(r string) string {
		r = withDefault(leadingString, strings.ToLower)(r)
		return ToUpperLeading(r)
	}, seps...)
}

// LowerCamelCase returns the CamelCased name by lowercase uppercase letter.
// In short, _my_field_name_2 becomes xMyFieldName_2.
// "theQuickBrownFoxJumpsOverTheLazyDog"
func LowerCamelCase(s string, seps ...rune) string {
	s = UpperCamelCase(s, seps...)
	sr, s := ExtractFirstRune(s)
	return ToLowerLeading(string(sr)) + s
}

// SnakeCase returns the SnakeCased name.
// In short, _my_field_name_2 becomes x_my_field_name_2.
// seps will append '_' if len(seps) == 0
// "the_quick_brown_fox_jumps_over_the_lazy_dog"
func SnakeCase(s string, seps ...rune) string {
	if len(seps) == 0 {
		seps = append(seps, '_')
	}
	return TransformCase(s, JoinGenerator("_", withDefault(leadingString, strings.ToLower)), seps...)
}

// DarwinCase returns the DarwinCased name.
// Darwin case uses underscores between words with initial uppercase letters, as in "Sample_Type"
// In short, _my_field_name_2 becomes X_My_Field_Name_2.
// see https://en.wikipedia.org/wiki/Camel_case
func DarwinCase(s string, seps ...rune) string {
	return TransformCase(s, JoinGenerator("_", withDefault(leadingString, func(s string) string {
		return strings.Title(strings.ToLower(s))
	})), seps...)
}

// KebabCase returns the KebabCased name.
// In short, _my_field_name_2 becomes x-my-field-name-2.
// "the-quick-brown-fox-jumps-over-the-lazy-dog"
func KebabCase(s string, seps ...rune) string {
	return TransformCase(s, JoinGenerator("-", withDefault(leadingString, strings.ToLower)), seps...)
}

// DotCase returns the KebabCased name.
// In short, _my_field_name_2 becomes x.my.field.name.2.
func DotCase(s string, seps ...rune) string {
	return TransformCase(s, JoinGenerator(".", withDefault(leadingString, strings.ToLower)), seps...)
}

// Studly caps is a form of text notation in which the capitalization of letters varies by some pattern, or arbitrarily,
// usually also omitting spaces between words and often omitting some letters, for example, StUdLyCaPs or STuDLyCaPS.
// Such patterns are identified by many users, ambiguously, as camel case.
// The typical alternative is to just replace spaces with underscores (as in snake case).
// Messages may be hidden in the capital and lower-case letters such as "ShoEboX" which spells
// "SEX" in capitals and "hobo" in lower-case.
// https://en.wikipedia.org/wiki/Studly_caps
// "tHeqUicKBrOWnFoXJUmpsoVeRThElAzydOG"
// "THiS iS aN eXCePTioNaLLy eLiTe SeNTeNCe"
func StudlyCapsCase(upperCase unicode.SpecialCase, s string) string {
	return strings.ToLowerSpecial(upperCase, s)
}

// "thEqUIckbrOwnfOxjUmpsOvErthElAzydOg"
func StudlyCapsVowelUpperCase(s string) string {
	return strings.ToLowerSpecial(unicode_.VowelCase(nil, func(r rune) rune {
		return unicode.ToUpper(r)
	}, nil), s)
}

// "THeQuiCKBRoWNFoXJuMPSoVeRTHeLaZYDoG"
func StudlyCapsConsonantUpperCase(s string) string {
	return strings.ToLowerSpecial(unicode_.ConsonantCase(nil, func(r rune) rune {
		return unicode.ToUpper(r)
	}, nil), s)
}

// lower_case_with_underscores
func LowerCaseWithUnderscores(s string, seps ...rune) string {
	return func(s_ string) string {
		return strings.ToLower(SnakeCase(s_, seps...))
	}(s)
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

// UpperCamelCaseSlice is like UpperCamelCase, but the argument is a slice of strings to
// be joined with "_".
func UpperCamelCaseSlice(elem ...string) string { return UpperCamelCase(strings.Join(elem, "_"), '_') }

// LowerCamelCaseSlice is like LowerCamelCase, but the argument is a slice of strings to
// be joined with "_".
func LowerCamelCaseSlice(elem ...string) string { return LowerCamelCase(strings.Join(elem, "_"), '_') }

// DottedSlice turns a sliced name into a dotted name.
func DottedSlice(elem ...string) string { return strings.Join(elem, ".") }

const (
	leadingString = "X"
)

// TransformCase Splits and apply map on every splits
func TransformCase(s string, join func(r string) string, seps ...rune) string {
	var out strings.Builder
	for _, sub := range splits(s, seps...) {
		out.WriteString(join(sub))
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
