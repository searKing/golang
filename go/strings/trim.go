// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// SplitPrefixNumber slices s into number prefix and unparsed and
// returns a slice of those substrings.
// If s does not start with number, SplitPrefixNumber returns
// a slice of length 1 whose only element is s.
// If s is with number only, SplitPrefixNumber returns
// a slice of length 1 whose only element is s.
func SplitPrefixNumber(s string) []string {
	unparsed := TrimPrefixNumber(s)
	if unparsed == "" || len(unparsed) == len(s) {
		return []string{s}
	}
	return []string{s[:len(s)-len(unparsed)], unparsed}
}

// TrimPrefixNumber returns s without the leading number prefix string.
// If s doesn't start with number prefix, s is returned unchanged.
func TrimPrefixNumber(s string) string {
	unparsedFloat := TrimPrefixFloat(s)
	unparsedInt := TrimPrefixInteger(s)
	if len(unparsedFloat) < len(unparsedInt) {
		return unparsedFloat
	}
	return unparsedInt
}

// TrimPrefixFloat returns s without the leading float prefix string.
// If s doesn't start with float prefix, s is returned unchanged.
func TrimPrefixFloat(s string) string {
	var value float64
	var unparsed string
	// Scanf, Fscanf, and Sscanf parse the arguments according to a format string, analogous to that of Printf.
	// In the text that follows, 'space' means any Unicode whitespace character except newline.
	// Input processed by verbs is implicitly space-delimited:
	// the implementation of every verb except %c starts by discarding leading spaces from the remaining input,
	// and the %s verb (and %v reading into a string) stops consuming input at the first space or newline character.
	// see https://golang.org/pkg/fmt/#hdr-Scanning
	space := strings.IndexFunc(s, func(r rune) bool {
		if r == '\n' {
			return false
		}
		return unicode.IsSpace(r)
	})
	if space < 0 {
		space = len(s)
	}
	count, err := fmt.Sscanf(s[:space], `%v%s`, &value, &unparsed)

	if (err != nil && err != io.EOF) || (count == 0) {
		return s
	}
	return s[space-len(unparsed):]
}

// TrimPrefixInteger returns s without the leading integer prefix string.
// If s doesn't start with integer prefix, s is returned unchanged.
func TrimPrefixInteger(s string) string {
	var value int64
	var unparsed string
	// Scanf, Fscanf, and Sscanf parse the arguments according to a format string, analogous to that of Printf.
	// In the text that follows, 'space' means any Unicode whitespace character except newline.
	// Input processed by verbs is implicitly space-delimited:
	// the implementation of every verb except %c starts by discarding leading spaces from the remaining input,
	// and the %s verb (and %v reading into a string) stops consuming input at the first space or newline character.
	// see https://golang.org/pkg/fmt/#hdr-Scanning
	space := strings.IndexFunc(s, func(r rune) bool {
		if r == '\n' {
			return false
		}
		return unicode.IsSpace(r)
	})
	if space < 0 {
		space = len(s)
	}
	count, err := fmt.Sscanf(s[:space], `%v%s`, &value, &unparsed)

	if (err != nil && err != io.EOF) || (count == 0) {
		return s
	}
	return s[space-len(unparsed):]
}

// TrimPrefixComplex returns s without the leading complex prefix string.
// If s doesn't start with complex prefix, s is returned unchanged.
func TrimPrefixComplex(s string) string {
	var value complex128
	var unparsed string
	// Scanf, Fscanf, and Sscanf parse the arguments according to a format string, analogous to that of Printf.
	// In the text that follows, 'space' means any Unicode whitespace character except newline.
	// Input processed by verbs is implicitly space-delimited:
	// the implementation of every verb except %c starts by discarding leading spaces from the remaining input,
	// and the %s verb (and %v reading into a string) stops consuming input at the first space or newline character.
	// see https://golang.org/pkg/fmt/#hdr-Scanning
	space := strings.IndexFunc(s, func(r rune) bool {
		if r == '\n' {
			return false
		}
		return unicode.IsSpace(r)
	})
	if space < 0 {
		space = len(s)
	}
	count, err := fmt.Sscanf(s[:space], `%v%s`, &value, &unparsed)

	if (err != nil && err != io.EOF) || (count == 0) {
		return s
	}
	return s[space-len(unparsed):]
}
