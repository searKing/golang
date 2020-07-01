// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux,!cgo

package io

import (
	"os"
)

func copyFileClone(srcFile, dstFile *os.File) error {
	return ErrNotImplemented
}
