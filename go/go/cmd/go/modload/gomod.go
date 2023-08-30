// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modload

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ImportPackage finds the package and returns
// its source directory (an element of {dir}'s parent and forefathers in order).
func ImportPackage(importPath string, dirs ...string) (srcdir string, modName string, err error) {
	return importPackage(importPath, dirs...)
}

// ImportFile finds the package containing importPath, and returns
// its source directory (an element of $GOPATH) and its import path
// relative to it.
func ImportFile(file string) (srcdir, importPath string, err error) {
	return importFile(file)
}

func importPackage(importPath string, dirs ...string) (srcdir string, modName string, err error) {
	for _, dir := range dirs {
		modRoot := findModuleRoot(dir)
		modName, err := FindModuleName(modRoot)
		if err != nil {
			return "", "", err
		}
		segmentedAbsFileDir := segments(importPath)
		d := prefixLen(segments(modName), segmentedAbsFileDir)
		if d < 0 {
			continue
		}
		resolvedFilename := filepath.Join(modRoot, filepath.Join(segmentedAbsFileDir[len(segmentedAbsFileDir)-d:]...))
		if _, err := os.Stat(resolvedFilename); !os.IsNotExist(err) {
			return modRoot, modName, nil
		}
	}
	return "", "", fmt.Errorf("import path %s is not beneath any of these directories: %s",
		importPath, strings.Join(dirs, ", "))
}

func importFile(file string) (srcdir, importPath string, err error) {
	fi, err := os.Lstat(file)
	if err != nil {
		return "", "", fmt.Errorf("can't evaluate stat of %s: %v", file, err)
	}
	absFile, err := filepath.Abs(file)
	if err != nil {
		return "", "", fmt.Errorf("can't form the absolute path of %s: %v", file, err)
	}
	var absFileDir string
	if fi.IsDir() {
		absFileDir = absFile
	} else {
		absFileDir = filepath.Dir(absFile)
	}
	resolvedAbsFileDir, err := filepath.EvalSymlinks(absFileDir)
	if err != nil {
		return "", "", fmt.Errorf("can't evaluate symlinks of %s: %v", absFileDir, err)
	}

	modRoot := FindModuleRoot(resolvedAbsFileDir)
	modName, err := FindModuleName(modRoot)
	if err != nil {
		return "", "", fmt.Errorf("can't find module name of %s: %v", modRoot, err)
	}

	importPath, err = filepath.Rel(modRoot, resolvedAbsFileDir)
	if err != nil {
		return "", "", fmt.Errorf("can't evaluate rel of %s base on %s: %v",
			resolvedAbsFileDir, modRoot, err)
	}

	if modName == "" {
		// This happens for e.g. $GOPATH/src/a.go, but
		// "" is not a valid path for (*go/build).Import.
		return "", "", fmt.Errorf("cannot load package in root of source directory %s", srcdir)
	}

	srcdir = modRoot
	importPath = filepath.Join(modName, importPath)

	return srcdir, filepath.ToSlash(importPath), nil
}

// FindModuleName Extract module name from {modRoot}/go.mod
func FindModuleName(modRoot string) (name string, err error) {
	gomod := filepath.Join(modRoot, "go.mod")
	f, err := os.Open(gomod)
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) < 0 {
			continue
		}
		if fields[0] != "module" {
			continue
		}
		if len(fields) != 2 {
			return "", fmt.Errorf("malformed module declaration in %s", gomod)
		}
		return fields[1], scanner.Err()
	}
	return "", scanner.Err()

}

// FindGoMod returns the nearest go.mod by walk dir and i's parent and forefathers in order
func FindGoMod(dir string) string {
	modRoot := findModuleRoot(dir)
	if modRoot == "" {
		return ""
	}

	// InDir checks whether path is in the file tree rooted at dir.
	// If so, InDir returns an equivalent path relative to dir.
	// If not, InDir returns an empty string.
	// InDir makes some effort to succeed even in the presence of symbolic links.
	rel, _ := filepath.Rel(os.TempDir(), modRoot)
	if rel == "." {
		// If you create /tmp/go.mod for experimenting,
		// then any tests that create work directories under /tmp
		// will find it and get modules when they're not expecting them.
		// It's a bit of a peculiar thing to disallow but quite mysterious
		// when it happens. See golang.org/issue/26708.
		return ""
	}
	return filepath.Join(dir, "go.mod")
}

// FindModuleRoot returns the nearest dir of go.mod by walk dir and dir's parent and forefathers in order
func FindModuleRoot(dir string) (root string) {
	return findModuleRoot(dir)
}

// borrow from golang src code
// //go:linkname findModuleRoot cmd/go/internal/modload.findModuleRoot
func findModuleRoot(dir string) (root string) {
	if dir == "" {
		panic("dir not set")
	}
	dir = filepath.Clean(dir)

	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}
	return ""
}

func segments(path string) []string {
	return strings.Split(path, string(os.PathSeparator))
}

// prefixLen returns the length of the remainder of y if x is a prefix
// of y, a negative number otherwise.
func prefixLen(prefix, path []string) int {
	d := len(path) - len(prefix)
	if d >= 0 {
		for i := range prefix {
			if path[i] != prefix[i] {
				return -1 // not a prefix
			}
		}
	}
	return d
}
