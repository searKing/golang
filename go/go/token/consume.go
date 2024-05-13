// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
)

// A mode value is a set of flags (or 0).
// They control scanner behavior.
type Mode uint

const (
	ModeCaseSensitive Mode = 1 << iota
	ModeRegexpPerl
	ModeRegexpPosix
)

func ConsumeIdentifier(inputs []rune, current int, runeType Type) (token Token, next int) {
	posBegin := current
	if current < 0 {
		current = 0
	}

	if current >= len(inputs) {
		return Token{
			Typ:   TypeEOF,
			Value: "",
		}, len(inputs)
	}

	char := inputs[current]
	var value bytes.Buffer

	// identifier = letter { letter | unicode_digit } .
	// letter        = unicode_letter | "_" .
	// decimal_digit = "0" … "9" .
	// octal_digit   = "0" … "7" .
	// hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
	// newline        = /* the Unicode code point U+000A */ .
	// unicode_char   = /* an arbitrary Unicode code point except newline */ .
	// unicode_letter = /* a Unicode code point classified as "Letter" */ .
	// unicode_digit  = /* a Unicode code point classified as "Number, decimal digit" */ .
	if unicode.IsLetter(char) || char == '_' {
		for unicode.IsLetter(char) || char == '_' || unicode.IsNumber(char) || char == '.' {
			value.WriteRune(char)
			current++
			if current >= len(inputs) {
				break
			}
			char = inputs[current]
		}

		return Token{
			Typ:   runeType,
			Value: value.String(),
		}, current
	}
	// restore pos
	return Token{Typ: TypeILLEGAL}, posBegin
}

func ComsumeRunesAny(inputs []rune, current int, runeType Type, expectRunes ...rune) (token Token, next int) {
	posBegin := current
	if current < 0 {
		current = 0
	}

	if current >= len(inputs) {
		return Token{
			Typ:   TypeEOF,
			Value: "",
		}, len(inputs)
	}

	char := inputs[current]
	current++

	for _, expect := range expectRunes {
		if char == expect {
			return Token{
				Typ:   runeType,
				Value: "",
			}, current
		}
	}
	// restore pos
	return Token{Typ: TypeILLEGAL}, posBegin
}

func ComsumeStringsAny(inputs []rune, current int, runeType Type, mode Mode, expectStrs ...string) (token Token, next int) {
	posBegin := current
	if current < 0 {
		current = 0
	}

	if current >= len(inputs) {
		return Token{
			Typ:   TypeEOF,
			Value: "",
		}, len(inputs)
	}

	// regex mode
	if mode&(ModeRegexpPerl|ModeRegexpPosix) != 0 {
		for _, expect := range expectStrs {
			var reg *regexp.Regexp
			if mode&ModeRegexpPosix != 0 {
				reg = regexp.MustCompilePOSIX(expect)
			} else {
				reg = regexp.MustCompile(expect)
			}

			matches := reg.FindStringSubmatch(string(inputs[current:]))
			if len(matches) == 0 {
				continue
			}

			current = current + len(matches[0])
			return Token{
				Typ:   runeType,
				Value: string(matches[0]),
			}, current
		}
		// restore pos
		return Token{Typ: TypeILLEGAL}, posBegin
	}

	// none regexp
	for _, expect := range expectStrs {

		endPos := current + len(expect)
		if endPos > len(inputs) {
			continue
		}
		selected := inputs[current:endPos]

		if ((mode&ModeCaseSensitive != 0) && strings.EqualFold(string(selected), expect)) ||
			string(selected) == expect {
			return Token{
				Typ:   runeType,
				Value: string(selected),
			}, endPos
		}
	}
	// restore pos
	return Token{Typ: TypeILLEGAL}, posBegin
}
