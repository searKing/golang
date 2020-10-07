// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"strings"
	"unicode"
)

// TagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type TagOptions string

// ParseTagOptions splits a struct field's json tag into its name and
// comma-separated options.
func ParseTagOptions(tag string) (tagName string, opts TagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], TagOptions(tag[idx+1:])
	}
	return tag, TagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o TagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

// IsValidTagKey checks if key of struct tag, like `json:"name,omitempty"`, is valid
func IsValidTagKey(tagKey string) bool {
	if tagKey == "" {
		return false
	}
	for _, c := range tagKey {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}
