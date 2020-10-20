// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"strings"
	"unicode"
)

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

// SliceCombine combine elements to a new slice.
func SliceCombine(ss ...[]string) []string {
	var total int
	for _, s := range ss {
		total += len(s)
	}
	if total == 0 {
		return nil
	}
	var tt = make([]string, 0, total)
	for _, s := range ss {
		tt = append(tt, s...)
	}
	return tt
}

// SliceEqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, which is a more general
// form of case-sensitivity.
func SliceEqual(s, t []string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != t[i] {
			return false
		}
	}
	return true
}

// SliceEqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, which is a more general
// form of case-insensitivity.
func SliceEqualFold(s, t []string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !strings.EqualFold(s[i], t[i]) {
			return false
		}
	}
	return true
}

// SliceTrimEmpty trim empty columns
func SliceTrimEmpty(ss ...string) []string {
	return SliceTrimFunc(ss, func(s string) bool {
		return s == ""
	})
}

// SliceTrimFunc returns a slice of the string ss satisfying f(c) removed.
func SliceTrimFunc(ss []string, f func(s string) bool) []string {
	var trimmed []string
	for _, s := range ss {
		if f(s) {
			continue
		}
		trimmed = append(trimmed, s)
	}
	return trimmed
}

// SliceContains  reports whether s is within ss.
func SliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// SliceUnique returns the given string slice with unique values.
func SliceUnique(i []string) []string {
	u := make([]string, 0, len(i))
	m := make(map[string]bool)

	for _, val := range i {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
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
