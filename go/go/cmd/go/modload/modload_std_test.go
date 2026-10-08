// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modload

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// goModRoot returns the module root of dir found by the go command.
func goModRoot(t *testing.T, gobin, dir string) string {
	cmd := exec.Command(gobin, "env", "GOMOD")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GO111MODULE=on", "GOWORK=off", "GOFLAGS=")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go env GOMOD in %s: %v", dir, err)
	}
	gomod := strings.TrimSpace(string(out))
	if gomod == "" || gomod == os.DevNull {
		return ""
	}
	return filepath.Dir(gomod)
}

// TestFindModuleRootMatchesGo detects drift between the copied findModuleRoot and the go command.
func TestFindModuleRootMatchesGo(t *testing.T) {
	gobin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go command not found")
	}

	// The go command reports paths with symlinks evaluated, e.g. /private/var rather than /var on macOS.
	tempDir := func() string {
		dir, err := filepath.EvalSymlinks(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		return dir
	}
	root := tempDir()
	for _, d := range []string{"a/b/c", "nested/x", "dirmod/go.mod", "dirmod/y"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"go.mod", "nested/go.mod"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("module example.com/m\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	outside := tempDir()

	for _, dir := range []string{
		root,
		filepath.Join(root, "a/b/c"),
		filepath.Join(root, "nested"),
		filepath.Join(root, "nested/x"),
		filepath.Join(root, "dirmod/y"), // go.mod as a directory is not a module
		outside,
	} {
		if got, want := findModuleRoot(dir), goModRoot(t, gobin, dir); got != want {
			t.Errorf("findModuleRoot(%q) = %q, go command found %q", dir, got, want)
		}
	}
}
