// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmux

import (
	"io"
	"net"
)

// Serve will be called in another goroutine
type Handler interface {
	Matcher
	Serve(l net.Listener)
}

type Matcher interface {
	Match(io.Writer, io.Reader) bool
}

// MatchWriter is a match that can also write response (say to do handshake).
type MatcherFunc func(io.Writer, io.Reader) bool

func (f MatcherFunc) Match(w io.Writer, r io.Reader) bool {
	return f(w, r)
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
