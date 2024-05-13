// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"io"
	"os"
)

var _ io.Closer = CloserFunc(nil)

// Stater is the interface that wraps the basic Stat method.
// Stat returns the FileInfo structure describing file.
type Stater interface {
	Stat() (os.FileInfo, error)
}

// The CloserFunc type is an adapter to allow the use of
// ordinary functions as io.Closer handlers. If f is a function
// with the appropriate signature, CloserFunc(f) is a
// Handler that calls f.
type CloserFunc func() error

// Close calls f(w, r).
func (f CloserFunc) Close() error {
	if f == nil {
		return nil
	}
	return f()
}

// CloseIf call Close() if arg implements io.Closer
func CloseIf(c any) error {
	if c, ok := c.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}
