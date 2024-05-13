// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

// The PrintfFunc type is an adapter to allow the use of
// ordinary functions as Printf handlers. If f is a function
// with the appropriate signature, PrintfFunc(f) is a
// Handler that calls f.
type PrintfFunc func(format string, a ...any)

// Write calls f(p).
func (f PrintfFunc) Write(p []byte) (n int, err error) {
	f("%s", string(p))
	return len(p), nil
}
