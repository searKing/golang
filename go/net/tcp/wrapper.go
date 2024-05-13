// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"bufio"
	"sync"
)

type TCPConn struct {
	*bufio.ReadWriter
	muRead  sync.Mutex
	muWrite sync.Mutex
}

func NewTCPConn(rw *bufio.ReadWriter) *TCPConn {
	if rw == nil {
		panic("nil io.ReadWriter")
	}
	return &TCPConn{
		ReadWriter: rw,
	}
}
