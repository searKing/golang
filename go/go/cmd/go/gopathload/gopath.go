// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopathload

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ImportPackage finds the package and returns
// its source directory (an element of $GOPATH).
func ImportPackage(importPath string) (srcdir string, err error) {
	return importPackage(importPath, &build.Default)
}

// ImportFile finds the package containing importPath, and returns
// its source directory (an element of $GOPATH) and its import path
// relative to it.
func ImportFile(file string) (srcdir, importPath string, err error) {
	return importFile(file, &build.Default)
}

func importPackage(importPath string, buildContext *build.Context) (srcdir string, err error) {
	srcDirs := build.Default.SrcDirs()
	for _, src := range srcDirs {
		resolvedFilename := filepath.Join(src, importPath)
		if _, err := os.Stat(resolvedFilename); !os.IsNotExist(err) {
			return src, nil
		}
	}
	return "", fmt.Errorf("import path %s is not beneath any of these GOROOT/GOPATH directories: %s",
		importPath, strings.Join(buildContext.SrcDirs(), ", "))
}

func importFile(file string, buildContext *build.Context) (srcdir, importPath string, err error) {
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

	segmentedAbsFileDir := segments(resolvedAbsFileDir)
	// Find the innermost directory in $GOPATH that encloses importPath.
	minD := 1024
	for _, srcDirInGoPath := range buildContext.SrcDirs() {
		absDir, err := filepath.Abs(srcDirInGoPath)
		if err != nil {
			continue // e.g. non-existent dir on $GOPATH
		}
		resolvedAbsDir, err := filepath.EvalSymlinks(absDir)
		if err != nil {
			continue // e.g. non-existent dir on $GOPATH
		}

		d := prefixLen(segments(resolvedAbsDir), segmentedAbsFileDir)
		// If there are multiple matches,
		// prefer the innermost enclosing directory
		// (smallest d).
		if d >= 0 && d < minD {
			minD = d
			srcdir = srcDirInGoPath
			importPath = path.Join(segmentedAbsFileDir[len(segmentedAbsFileDir)-minD:]...)
		}
	}
	if srcdir == "" {
		return "", "", fmt.Errorf("directory %s is not beneath any of these GOROOT/GOPATH directories: %s",
			absFileDir, strings.Join(buildContext.SrcDirs(), ", "))
	}
	if importPath == "" {
		// This happens for e.g. $GOPATH/src/a.go, but
		// "" is not a valid path for (*go/build).Import.
		return "", "", fmt.Errorf("cannot load package in root of source directory %s", srcdir)
	}
	return srcdir, filepath.ToSlash(importPath), nil
}

func segments(path string) []string {
	return strings.Split(path, string(os.PathSeparator))
}

// prefixLen returns the length of the remainder of y if x is a prefix
// of y, a negative number otherwise.
func prefixLen(x, y []string) int {
	d := len(y) - len(x)
	if d >= 0 {
		for i := range x {
			if y[i] != x[i] {
				return -1 // not a prefix
			}
		}
	}
	return d
}
