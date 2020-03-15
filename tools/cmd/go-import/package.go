// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "path"

// Package holds a single parsed package and associated files and ast files.
type Package struct {
	// Name is the package name as it appears in the package source code.
	name string

	importPath string

	// configs
	importPrefix string
	lineComment  bool
	globImport   string
	buildTag     string
}

func (pkg Package) Package() string {
	if pkg.name != "" {
		return pkg.name
	}
	return path.Base(pkg.importPath)
}
