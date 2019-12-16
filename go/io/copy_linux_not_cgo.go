// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.
//+build linux,!cgo

package io

import (
	"os"
)

func copyFileClone(srcFile, dstFile *os.File) error {
	return ErrNotImplemented
}