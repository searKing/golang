// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filepath

import (
	"path/filepath"
)

// Glob returns the names of all files matching pattern satisfying f(c) or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
func GlobFunc(pattern string, handler func(name string) bool) (matches []string, err error) {
	matches, err = filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if handler == nil {
		return matches, err
	}

	var a []string
	for _, match := range matches {
		if handler(match) {
			a = append(a, match)
		}
	}
	return a, nil
}
