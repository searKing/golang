// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"errors"
	"io"
	"os"
)

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")
var errSeeker = errors.New("Seek: can't seek")

// Stater is the interface that wraps the basic Stat method.
// Stat returns the FileInfo structure describing file.
type Stater interface {
	Stat() (os.FileInfo, error)
}

// LimitReader returns a Reader that reads from r
// but stops with EOF after n bytes.
// The underlying implementation is a *LimitedReader.
func LimitReadSeeker(r io.ReadSeeker, n int64) io.ReadSeeker { return &LimitedReadSeeker{r, n} }

// A LimitSeekable reads from R but limits the size of the file N bytes.
// Read returns EOF when N <= 0 or when the underlying R returns EOF.
type LimitedReadSeeker struct {
	rs    io.ReadSeeker // underlying readSeeker
	limit int64         // max bytes remaining
}

func (l *LimitedReadSeeker) Read(p []byte) (n int, err error) {
	// speedup
	if l.limit <= 0 {
		return 0, io.EOF
	}

	offset, err := l.rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, errOffset
	}

	readLimit := l.limit - offset

	if readLimit <= 0 {
		return 0, io.EOF
	}

	if int64(len(p)) > readLimit {
		p = p[0:readLimit]
	}
	n, err = l.rs.Read(p)
	return
}

func (l *LimitedReadSeeker) Seek(offset int64, whence int) (int64, error) {
	lastPos, err := l.rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, errOffset
	}

	size, err := l.rs.Seek(offset, whence)
	if err != nil {
		return 0, errSeeker
	}
	if size > l.limit {
		// recover if overflow
		_, _ = l.rs.Seek(lastPos, io.SeekStart)
		return 0, errOffset
	}
	return size, nil
}
