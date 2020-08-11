// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import "net"

type HandlerConn interface {
	Serve(net.Conn)
}

// HandlerConnFunc is a match that can also write response (say to do handshake).
type HandlerConnFunc func(net.Conn)

func (f HandlerConnFunc) Serve(c net.Conn) {
	f(c)
}

var ignoreErrorHandler = ErrorHandlerFunc(func(_ error) bool { return true })

type ErrorHandler interface {
	Continue(error) bool
}

// ErrorHandler handles an error and returns whether
// the mux should continue serving the listener.
type ErrorHandlerFunc func(error) bool

func (f ErrorHandlerFunc) Continue(err error) bool {
	return f(err)
}
