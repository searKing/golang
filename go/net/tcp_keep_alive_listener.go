// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"net"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
	d time.Duration
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if ln.d == 0 {
		_ = tc.SetKeepAlive(false)
		return tc, nil
	}
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(ln.d)
	return tc, nil
}

// TcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
func TcpKeepAliveListener(l *net.TCPListener, d time.Duration) net.Listener {
	return &tcpKeepAliveListener{l, d}
}
