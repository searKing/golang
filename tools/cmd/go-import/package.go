// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

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
	buildTags    []string
}

func (pkg Package) Package() string {
	if pkg.name != "" {
		return pkg.name
	}
	return path.Base(pkg.importPath)
}
