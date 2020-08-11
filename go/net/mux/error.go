// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"errors"
	"net"

	net_ "github.com/searKing/golang/go/net"
)

// type check
var _ net.Error = ErrListenerClosed

// ErrListenerClosed is returned from MuxListener.Accept when the underlying
// listener is closed.
var ErrListenerClosed = net_.ErrListenerClosed

// ErrServerClosed is returned by the Server's Serve, ServeTLS, ListenAndServe,
// and ListenAndServeTLS methods after a call to Shutdown or Close.
var ErrServerClosed = errors.New("net/mux: Server closed")

// ErrAbortHandler is a sentinel panic value to abort a handler.
// While any panic from ServeHTTP aborts the response to the client,
// panicking with ErrAbortHandler also suppresses logging of a stack
// trace to the server's error log.
var ErrAbortHandler = errors.New("net/mux: abort HandlerConn")
