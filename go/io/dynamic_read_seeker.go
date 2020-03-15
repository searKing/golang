// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"io"
	"sync"
)

// DynamicReadSeeker returns a ReadSeeker that reads from r got by getter at an offset.
// The underlying implementation is a *dynamicReadSeeker.
func DynamicReadSeeker(getter func(off int64) (io.Reader, error), size int64) io.ReadSeeker {
	return &dynamicReadSeeker{
		getter:    getter,
		totalSize: size,
	}
}

// A dynamicReadSeeker reads from r got by getter at an offset.
type dynamicReadSeeker struct {
	getter    func(off int64) (io.Reader, error)
	totalSize int64 // make no sense if rs implements io.Seeker

	rs         io.Reader // underlying readSeeker
	lastOffset int64     // make no sense if rs implements io.Seeker
	once       sync.Once
}

func (l *dynamicReadSeeker) lazyLoad() {
	l.once.Do(func() {
		if l.getter == nil {
			return
		}

		if l.rs == nil {
			l.rs, _ = l.getter(0)
			l.lastOffset = 0
		}
	})
}

func (l *dynamicReadSeeker) Read(p []byte) (n int, err error) {
	l.lazyLoad()

	if l.rs == nil {
		return 0, io.EOF
	}

	if l.lastOffset >= l.totalSize {
		return 0, io.EOF
	}

	n, err = l.rs.Read(p)
	if n >= 0 {
		l.lastOffset += int64(n)
	}

	return
}

func (l *dynamicReadSeeker) Seek(offset int64, whence int) (n int64, err error) {
	l.lazyLoad()
	if l.rs == nil {
		return 0, errSeeker
	}
	if seeker, ok := l.rs.(io.Seeker); ok {
		n, err = seeker.Seek(offset, whence)
		l.lastOffset = n
		return
	}

	// speed up
	if whence == io.SeekCurrent && offset == 0 {
		n = l.lastOffset
		return n, nil
	}

	if whence == io.SeekStart && offset == l.lastOffset {
		n = l.lastOffset
		return n, nil
	}

	switch whence {
	case io.SeekStart:
		break
	case io.SeekCurrent:
		offset += l.lastOffset
	case io.SeekEnd:
		offset += l.totalSize
	}

	if offset >= l.totalSize {
		l.lastOffset = offset
		return offset, nil
	}

	if err := l.Close(); err != nil {
		return 0, err
	}

	if l.getter == nil {
		return 0, errSeeker
	}

	l.rs, err = l.getter(offset)
	if err != nil {
		return 0, err
	}
	if l.rs == nil {
		return 0, errSeeker
	}

	l.lastOffset = offset
	return offset, nil
}

func (l *dynamicReadSeeker) Close() error {
	defer func() {
		l.lastOffset = 0
	}()

	if l.rs == nil {
		return nil
	}
	defer func() {
		l.rs = nil
	}()

	if closer, ok := l.rs.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
