// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
	"time"

	os_ "github.com/searKing/golang/go/os"
)

func TestUnlinkOldestFiles(t *testing.T) {
	tests := []struct {
		name     string
		files    int
		symlink  bool // a.link links to the newest file
		maxCount int
		wantLeft []string
	}{
		{name: "under max count", files: 2, maxCount: 5, wantLeft: []string{"a.0", "a.1"}},
		{name: "equal max count", files: 3, maxCount: 3, wantLeft: []string{"a.0", "a.1", "a.2"}},
		{name: "exceed max count", files: 5, maxCount: 2, wantLeft: []string{"a.3", "a.4"}},
		{name: "symlink not counted", files: 3, symlink: true, maxCount: 2, wantLeft: []string{"a.1", "a.2", "a.link"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			now := time.Now()
			for i := 0; i < tt.files; i++ {
				name := filepath.Join(dir, "a."+strconv.Itoa(i))
				if err := os.WriteFile(name, []byte("x"), 0644); err != nil {
					t.Fatal(err)
				}
				// a.0 is the oldest
				mt := now.Add(time.Duration(i-tt.files) * time.Minute)
				if err := os.Chtimes(name, mt, mt); err != nil {
					t.Fatal(err)
				}
			}
			if tt.symlink {
				newest := filepath.Join(dir, "a."+strconv.Itoa(tt.files-1))
				if err := os.Symlink(newest, filepath.Join(dir, "a.link")); err != nil {
					t.Fatal(err)
				}
			}

			err := os_.UnlinkOldestFiles(filepath.Join(dir, "a.*"), os_.DiskQuota{MaxCount: tt.maxCount})
			if err != nil {
				t.Fatalf("UnlinkOldestFiles() error = %v", err)
			}

			matches, _ := filepath.Glob(filepath.Join(dir, "a.*"))
			var left []string
			for _, m := range matches {
				left = append(left, filepath.Base(m))
			}
			slices.Sort(left)
			if !slices.Equal(left, tt.wantLeft) {
				t.Errorf("files left = %v, want %v", left, tt.wantLeft)
			}
		})
	}
}

func TestUnlinkOldestFiles_SameModTime(t *testing.T) {
	dir := t.TempDir()
	mt := time.Now().Add(-time.Minute)
	// ordered as a.log.1, a.log.2, a.log.9, a.log.10, a.log if ModTime equal
	for _, name := range []string{"a.log", "a.log.1", "a.log.2", "a.log.9", "a.log.10"} {
		name = filepath.Join(dir, name)
		if err := os.WriteFile(name, []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(name, mt, mt); err != nil {
			t.Fatal(err)
		}
	}

	if err := os_.UnlinkOldestFiles(filepath.Join(dir, "a.log*"), os_.DiskQuota{MaxCount: 2}); err != nil {
		t.Fatalf("UnlinkOldestFiles() error = %v", err)
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "a.log*"))
	var left []string
	for _, m := range matches {
		left = append(left, filepath.Base(m))
	}
	slices.Sort(left)
	if want := []string{"a.log", "a.log.10"}; !slices.Equal(left, want) {
		t.Errorf("files left = %v, want %v", left, want)
	}
}
