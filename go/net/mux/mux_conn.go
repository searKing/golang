// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"io"
	"net"

	io_ "github.com/searKing/golang/go/io"
)

// sniffConn wraps a net.Conn and provides transparent sniffing of connection data.
type sniffConn struct {
	net.Conn
	sniffer io_.ReadSniffer
}

func newMuxConn(c net.Conn) *sniffConn {
	return &sniffConn{
		Conn:    c,
		sniffer: io_.SniffReader(c),
	}
}

// From the io.Reader documentation:
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read.  It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil.  The next Read should
// return 0, EOF.
func (m *sniffConn) Read(p []byte) (int, error) {
	return m.sniffer.Read(p)
}

func (m *sniffConn) startSniffing() io.Reader {
	m.sniffer.Sniff(true)
	return m.sniffer
}

func (m *sniffConn) doneSniffing() {
	m.sniffer.Sniff(false)
}
