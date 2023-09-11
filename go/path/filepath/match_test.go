// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filepath_test

import (
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	filepath_ "github.com/searKing/golang/go/path/filepath"
)

// contains reports whether vector contains the string s.
func contains(vector []string, s string) bool {
	for _, elem := range vector {
		if elem == s {
			return true
		}
	}
	return false
}

var globTests = []struct {
	pattern, result string
}{
	{"match.go", "match.go"},
	{"mat?h.go", "match.go"},
	{"./mat?h.go", "match.go"},
	{"*", "match.go"},
	{"../*/match.go", "../filepath/match.go"},
	{"./*", "match.go"},
	{"../*/match.go", "../filepath/match.go"},
	{"../../*/*/match.go", "../../path/filepath/match.go"},

	// ** for zero or more directories Not Support, https://github.com/golang/go/issues/11862
	//{"../../**/match.go", "../filepath/match.go"},

	// no magic characters recognized by [filepath.Match].
	{"./match.go", "./match.go"}, // return if no
	{"../filepath/match.go", "../filepath/match.go"},
}

func TestWalkGlobDir(t *testing.T) {
	for _, tt := range globTests {
		pattern := tt.pattern
		result := tt.result
		if runtime.GOOS == "windows" {
			pattern = filepath.Clean(pattern)
			result = filepath.Clean(result)
		}
		var matches []string
		err := filepath_.WalkGlob(pattern, func(path string) error {
			matches = append(matches, path)
			return nil
		})
		if err != nil {
			t.Errorf("WalkGlob error for %q: %s", pattern, err)
			continue
		}
		if !contains(matches, result) {
			t.Errorf("WalkGlob(%#q, ) = %#v want %v", pattern, matches, result)
		}
		expectMatches, expectErr := filepath.Glob(pattern)
		if expectErr != nil {
			t.Errorf("filepath.Glob error for %q: %s", pattern, err)
			continue
		}
		if !slices.Equal(matches, expectMatches) {
			t.Errorf("WalkGlob(%#q, ) = %#v want %v", pattern, matches, expectMatches)
		}
	}
	for _, pattern := range []string{"no_match", "../*/no_match"} {
		var matches []string
		err := filepath_.WalkGlob(pattern, func(path string) error {
			matches = append(matches, path)
			return nil
		})
		if err != nil {
			t.Errorf("WalkGlob error for %q: %s", pattern, err)
			continue
		}
		if len(matches) != 0 {
			t.Errorf("WalkGlob(%#q, ) = %#v want []", pattern, matches)
		}
		expectMatches, expectErr := filepath.Glob(pattern)
		if expectErr != nil {
			t.Errorf("filepath.Glob error for %q: %s", pattern, err)
			continue
		}
		if !slices.Equal(matches, expectMatches) {
			t.Errorf("WalkGlob(%#q, ) = %#v want %v", pattern, matches, expectMatches)
		}
	}
}
