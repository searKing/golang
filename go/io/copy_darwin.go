// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import "os"

func copyFileClone(srcFile, dstFile *os.File) error {
	return ErrNotImplemented
}

func copyFileRange(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int, err error) {
	return 0, ErrNotImplemented
}
