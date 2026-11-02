// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux && !windows

package sync

import (
	"bytes"
	"runtime"
	"strconv"
)

// callerID returns the ID of the calling goroutine, or 0 if unknown.
// It is slow, as runtime.Stack serializes concurrent callers.
func callerID() uint64 {
	var buf [64]byte
	// The stack trace starts with "goroutine 123 [running]:".
	b := bytes.TrimPrefix(buf[:runtime.Stack(buf[:], false)], []byte("goroutine "))
	if i := bytes.IndexByte(b, ' '); i >= 0 {
		b = b[:i]
	}
	id, _ := strconv.ParseUint(string(b), 10, 64)
	return id
}
