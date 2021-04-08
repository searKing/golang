// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"strings"
	"unicode"
	"unicode/utf8"

	bytes_ "github.com/searKing/golang/go/bytes"
	unicode_ "github.com/searKing/golang/go/unicode"
)

// ContainsRuneAnyFunc reports whether any of the Unicode code point r satisfying f(r) is within s.
func ContainsRuneAnyFunc(s string, f func(rune) bool) bool {
	if f == nil {
		return false
	}
	for _, r := range s {
		if f(r) {
			return true
		}
	}
	return false
}

// ContainsRuneOnlyFunc reports whether all of the Unicode code point r satisfying f(r) is within s.
func ContainsRuneOnlyFunc(s string, f func(rune) bool) bool {
	if f == nil {
		return true
	}
	for _, r := range s {
		if !f(r) {
			return false
		}
	}
	return true
}

// ContainsAnyRangeTable reports whether the string contains any rune in any of the specified table of ranges.
func ContainsAnyRangeTable(s string, rangeTabs ...*unicode.RangeTable) bool {
	if len(rangeTabs) == 0 {
		return ContainsRuneAnyFunc(s, nil)
	}
	return ContainsRuneAnyFunc(s, func(r rune) bool {
		for _, t := range rangeTabs {
			if t == nil {
				continue
			}
			if unicode.Is(t, r) {
				return true
			}
		}
		return false
	})
}

// ContainsOnlyRangeTable reports whether the string contains only rune in all of the specified table of ranges.
func ContainsOnlyRangeTable(s string, rangeTabs ...*unicode.RangeTable) bool {
	if len(rangeTabs) == 0 {
		return ContainsRuneOnlyFunc(s, nil)
	}
	return ContainsRuneOnlyFunc(s, func(r rune) bool {
		for _, t := range rangeTabs {
			if t == nil {
				continue
			}
			if !unicode.Is(t, r) {
				return false
			}
		}
		return true
	})
}

// ContainsAsciiVisual reports whether the string contains any rune in visual ascii code, that is [0x21, 0x7E].
func ContainsAsciiVisual(s string) bool {
	return ContainsAnyRangeTable(s, unicode_.AsciiVisual)
}

// ContainsAsciiVisual reports whether the string contains only rune in visual ascii code, that is [0x21, 0x7E].
func ContainsOnlyAsciiVisual(s string) bool {
	return ContainsOnlyRangeTable(s, unicode_.AsciiVisual)
}

// JoinRepeat behaves like strings.Join([]string{s,...,s}, sep)
func JoinRepeat(s string, sep string, n int) string {
	var b strings.Builder
	for i := 0; i < n-1; i++ {
		b.WriteString(s)
		b.WriteString(sep)
	}
	if n > 0 {
		b.WriteString(s)
	}
	return b.String()
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

// Truncate shrinks s's len to n at most
func Truncate(s string, n int) string {
	if n < 0 {
		n = 0
	}
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// PadLeft returns s padded to length n, padded left with repeated pad
// return s directly if pad is empty
// padding s with {{pad}} and spaces(less than len(pad)) as a prefix, as [pad]...[pad][space]...[space][s]
func PadLeft(s string, pad string, n int) string {
	if len(pad) == 0 {
		return s
	}

	pc, sc := ComputePad(s, pad, n)

	return strings.Repeat(pad, pc) + strings.Repeat(" ", sc) + s
}

// PadRight returns s padded to length n, padded right with repeated pad
// return s directly if pad is empty
// padding s with {{pad}} and spaces(less than len(pad))  as a suffix, as [s][space]...[space][pad]...[pad]
func PadRight(s string, pad string, n int) string {
	if len(pad) == 0 {
		return s
	}
	pc, sc := ComputePad(s, pad, n)

	return s + strings.Repeat(" ", sc) + strings.Repeat(pad, pc)
}

// ComputePad returns pad's count and space's count(less than len(pad)) will be need to pad s to len n
// padCount = (n-len(s))/len(pad)
// spaceCount = (n-len(s))%len(pad)
func ComputePad(s string, pad string, n int) (padCount, spaceCount int) {
	return bytes_.ComputePad([]byte(s), []byte(pad), n)
}

// ReverseByByte returns a string with the bytes of s in reverse order.
func ReverseByByte(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.String()
}

// ReverseByRune returns a string with the runes of s in reverse order.
// Invalid UTF-8 sequences, if any, will be reversed byte by byte.
func ReverseByRune(s string) string {
	res := make([]byte, len(s))
	prevPos, resPos := 0, len(s)
	for pos := range s {
		resPos -= pos - prevPos
		copy(res[resPos:], s[prevPos:pos])
		prevPos = pos
	}
	copy(res[0:], s[prevPos:])
	return string(res)
}

// CountPrefix counts the number of non-overlapping instances of continuous substr prefix in s.
// If substr is an empty string, CountPrefix returns 1 + the number of Unicode code points in s.
func CountPrefix(s, substr string) int {
	// special case
	if len(substr) == 0 {
		return utf8.RuneCountInString(s) + 1
	}
	n := 0
	gap := len(substr)
	for i := 0; i+gap <= len(s); i += gap {
		if s[i:i+gap] == substr {
			n++
			continue
		}
		break
	}
	return n
}

// CountPrefix counts the number of non-overlapping instances of continuous substr suffix in s.
// If substr is an empty string, CountPrefix returns 1 + the number of Unicode code points in s.
func CountSuffix(s, substr string) int {
	// special case
	if len(substr) == 0 {
		return utf8.RuneCountInString(s) + 1
	}
	n := 0
	gap := len(substr)
	for i := len(s); i-gap >= 0; i -= gap {
		if s[i-gap:i] == substr {
			n++
			continue
		}
		break
	}
	return n
}
