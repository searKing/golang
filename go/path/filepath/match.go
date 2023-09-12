// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filepath

import (
	"errors"
	"path/filepath"
)

// WalkGlobFunc is the type of the function called by WalkGlob to visit
// all files matching pattern.
//
// The error result returned by the function controls how WalkGlob
// continues. If the function returns the special value [filepath.SkipAll],
// WalkGlob skips all files matching pattern satisfying f(c). Otherwise,
// if the function returns a non-nil error, WalkGlob stops entirely and
// returns that error.
type WalkGlobFunc func(path string) error

// WalkGlob returns the names of all files matching pattern satisfying f(c) or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
func WalkGlob(pattern string, fn WalkGlobFunc) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		err = fn(match)
		if err != nil && errors.Is(err, filepath.SkipAll) {
			// Successfully skipped all remaining files and directories.
			err = nil
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// GlobFunc returns the names of all files matching pattern satisfying f(c) or nil
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
