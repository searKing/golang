// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	os_ "github.com/searKing/golang/go/os"
)

// countOpenFiles returns the number of open file descriptors of the process.
func countOpenFiles(t *testing.T) int {
	t.Helper()
	entries, err := os.ReadDir("/dev/fd")
	if err != nil {
		t.Skipf("count open files: %v", err)
	}
	return len(entries)
}

func TestRotateFile_FileLinkPath(t *testing.T) {
	tests := []struct {
		name         string
		prefix       func(dir string) string
		link         func(dir string) string
		wantRelative bool
	}{
		{
			name:         "relative in same dir",
			prefix:       func(string) string { return filepath.Join("log", "test.") },
			link:         func(string) string { return filepath.Join("log", "s.log") },
			wantRelative: true,
		},
		{
			name:         "relative in different dirs",
			prefix:       func(string) string { return filepath.Join("logs", "app.") },
			link:         func(string) string { return filepath.Join("links", "app.log") },
			wantRelative: true,
		},
		{
			name:   "absolute",
			prefix: func(dir string) string { return filepath.Join(dir, "logs", "app.") },
			link:   func(dir string) string { return filepath.Join(dir, "links", "app.log") },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			// the dir of link is not created in advance
			link := tt.link(dir)

			f := os_.NewRotateFile("2006-01-02")
			f.FilePathPrefix = tt.prefix(dir)
			f.FileLinkPath = link
			if _, err := f.WriteString("hello"); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}

			got, err := os.ReadFile(link)
			if err != nil {
				t.Fatalf("read through link: %v", err)
			}
			if string(got) != "hello" {
				t.Errorf("content through link = %q, want %q", got, "hello")
			}
			target, err := os.Readlink(link)
			if err != nil {
				t.Fatal(err)
			}
			if gotRelative := !filepath.IsAbs(target); gotRelative != tt.wantRelative {
				t.Errorf("link target %q is relative = %v, want %v", target, gotRelative, tt.wantRelative)
			}
		})
	}
}

func TestRotateFile_ForceNewFileOnStartupNoFdLeak(t *testing.T) {
	// files leaked are closed by finalizer on GC
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	dir := t.TempDir()
	before := countOpenFiles(t)
	for i := 0; i < 10; i++ {
		f := os_.NewRotateFile("2006-01-02")
		f.FilePathPrefix = filepath.Join(dir, "app"+strconv.Itoa(i)+".")
		f.ForceNewFileOnStartup = true
		if _, err := f.WriteString("a"); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if after := countOpenFiles(t); after != before {
		t.Errorf("open files = %d after 10 rotate files closed, want %d", after, before)
	}
}

// sleepToNextSecond sleeps until just after the next second boundary,
// so that following writes are in the same RotateInterval of one second.
func sleepToNextSecond() {
	now := time.Now()
	time.Sleep(now.Truncate(time.Second).Add(time.Second + 10*time.Millisecond).Sub(now))
}

func TestRotateFile_CopyTruncateRotateOncePerInterval(t *testing.T) {
	dir := t.TempDir()
	f := os_.NewRotateFile("2006-01-02_15-04-05")
	f.FilePathPrefix = filepath.Join(dir, "app.")
	f.RotateMode = os_.RotateModeCopyTruncate
	f.RotateInterval = time.Second
	var rotated []string
	f.PostRotateHandler = func(name string) { rotated = append(rotated, name) }
	defer f.Close()

	sleepToNextSecond()
	if _, err := f.WriteString("a"); err != nil {
		t.Fatal(err)
	}
	if len(rotated) != 1 {
		t.Fatalf("rotations on startup = %d, want 1", len(rotated))
	}
	writing := rotated[0]

	sleepToNextSecond()
	for i := 0; i < 20; i++ {
		if _, err := f.WriteString("b"); err != nil {
			t.Fatal(err)
		}
	}
	if len(rotated) != 2 {
		t.Fatalf("rotations after 20 writes in next interval = %d, want 2", len(rotated))
	}
	if rotated[1] != writing {
		t.Errorf("writing file after copytruncate = %q, want %q", rotated[1], writing)
	}

	got, err := os.ReadFile(writing)
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Repeat("b", 20); string(got) != want {
		t.Errorf("writing file content = %q, want %q", got, want)
	}
}
