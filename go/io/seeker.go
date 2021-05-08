// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"errors"
	"io"
)

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")
var errSeeker = errors.New("Seek: can't seek")

// SeekerLen returns the length of the file and an error, if any.
func SeekerLen(s io.Seeker) (int64, error) {
	curOffset, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	endOffset, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = s.Seek(curOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return endOffset - curOffset, nil
}

// SniffCopy copies the seekable reader to an io.Writer
func SniffCopy(dst io.Writer, src io.ReadSeeker) (int64, error) {
	curPos, err := src.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// copy errors may be assumed to be from the body.
	n, err := io.Copy(dst, src)
	if err != nil {
		return n, err
	}

	// seek back to the first position after reading to reset
	// the body for transmission.
	_, err = src.Seek(curPos, io.SeekStart)
	if err != nil {
		return n, err
	}

	return n, nil
}

// SniffRead reads up to len(p) bytes into p.
func SniffRead(p []byte, src io.ReadSeeker) (int, error) {
	curPos, err := src.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// copy errors may be assumed to be from the body.
	n, err := src.Read(p)
	if err != nil {
		return n, err
	}

	// seek back to the first position after reading to reset
	// the body for transmission.
	_, err = src.Seek(curPos, io.SeekStart)
	if err != nil {
		return n, err
	}

	return n, nil
}

// LimitReadSeeker returns a Reader that reads from r
// but stops with EOF after n bytes.
// The underlying implementation is a *LimitedReader.
func LimitReadSeeker(r io.ReadSeeker, n int64) io.ReadSeeker { return &LimitedReadSeeker{r, n} }

// LimitedReadSeeker A LimitSeekable reads from R but limits the size of the file N bytes.
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
