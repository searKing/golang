// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmux

import (
	"io"

	"github.com/searKing/golang/go/container/trie_tree/ternary_search_tree"
)

// Any is a Matcher that matches any connection.
func Any() MatcherFunc {
	return func(w io.Writer, r io.Reader) bool { return true }
}

// PrefixMatcher returns a matcher that matches a connection if it
// starts with any of the strings in strs.
func PrefixMatcher(strs ...string) MatcherFunc {
	tree := ternary_search_tree.New(strs...)
	return func(w io.Writer, r io.Reader) bool {
		buf := make([]byte, tree.Depth())
		n, _ := io.ReadFull(r, buf)
		_, _, ok := tree.Follow(string(buf[:n]))
		return ok
	}
}

func PrefixByteMatcher(list ...[]byte) MatcherFunc {
	tree := ternary_search_tree.NewWithBytes(list...)
	return func(w io.Writer, r io.Reader) bool {
		buf := make([]byte, tree.Depth())
		n, _ := io.ReadFull(r, buf)
		_, _, ok := tree.Follow(string(buf[:n]))
		return ok
	}
}
