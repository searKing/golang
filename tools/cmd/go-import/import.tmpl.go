// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type ImportTemplateInfo struct {
	GoImportToolName string
	GoImportToolArgs []string
	ModuleName       string
	ImportPaths      []string
	BuildTag         string
}
