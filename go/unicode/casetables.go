// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unicode

import (
	"unicode"
)

var Vowels = []rune{
	'A', 'E', 'I', 'O', 'U',
	'a', 'e', 'i', 'o', 'u',
}

var Consonants = []rune{
	'B', 'C', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'V', 'W', 'X', 'Y', 'Z',
	'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'y', 'z',
}

var VowelCase = func(toUpper, toLower, toTitle func(r rune) rune) unicode.SpecialCase {
	return SpecialCaseBuilder(toUpper, toLower, toTitle, Vowels...)
}

var ConsonantCase = func(toUpper, toLower, toTitle func(r rune) rune) unicode.SpecialCase {
	return SpecialCaseBuilder(toUpper, toLower, toTitle, Consonants...)
}

func SpecialCaseBuilder(toUpper, toLower, toTitle func(r rune) rune, points ...rune) unicode.SpecialCase {
	apply := func(to func(r rune) rune, r rune) rune {
		if to == nil {
			return 0
		}
		return to(r) - r
	}
	var cases unicode.SpecialCase
	for _, point := range points {
		cases = append(cases, unicode.CaseRange{
			Lo: uint32(point),
			Hi: uint32(point),
			Delta: [unicode.MaxCase]rune{
				apply(toUpper, point),
				apply(toLower, point),
				apply(toTitle, point),
			}})
	}
	return cases
}

var _AsciiVisual = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0021, 0x007E, 1},
	},
	LatinOffset: 1,
}
var (
	AsciiVisual = _AsciiVisual // AsciiVisual is the set of Unicode characters with visual character of ascii.
)
