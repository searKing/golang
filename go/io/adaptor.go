// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

type WriterFunc func(p []byte) (n int, err error)

func (f WriterFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

type WriterFuncPrintfLike func(format string, args ...any)

func (f WriterFuncPrintfLike) Write(p []byte) (n int, err error) {
	f("%s", string(p))
	return len(p), nil
}
