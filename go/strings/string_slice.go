// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"strings"
)

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

// SliceTrim returns a slice of the string ss with tt removed.
func SliceTrim(ss []string, tt ...string) []string {
	return SliceTrimFunc(ss, func(s string) bool {
		return SliceContains(tt, s)
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

// SliceContains  reports whether any t in tt is within ss.
func SliceContains(ss []string, tt ...string) bool {
	if len(tt) == 0 {
		return false
	}

	if len(tt) == 1 {
		for _, v := range ss {
			if v == tt[0] {
				return true
			}
		}
		return false
	}
	return SliceContains(ss, tt[1:]...)
}

// SliceUnique returns the given string slice with unique values.
func SliceUnique(s ...string) []string {
	if len(s) <= 0 {
		return nil
	}
	u := make([]string, 0, len(s))
	m := make(map[string]bool)

	for _, val := range s {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
