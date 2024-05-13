// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

// A ConnState represents the state of a client connection to a server.
// It's used by the optional Server.ConnStateHook hook.
//
//go:generate go-enum -type ConnState -trimprefix=ConnState
type ConnState int

const (
	// ConnStateNew represents a new connection that is expected to
	// send a request immediately. Connections begin at this
	// state and then transition to either ConnStateActive or
	// ConnStateClosed.
	ConnStateNew ConnState = iota

	// ConnStateActive represents a connection that has read 1 or more
	// bytes of a request. The Server.ConnStateHook hook for
	// ConnStateActive fires before the request has entered a handler
	// and doesn't fire again until the request has been
	// handled. After the request is handled, the state
	// transitions to ConnStateClosed, ConnStateHijacked, or ConnStateIdle.
	// For HTTP/2, ConnStateActive fires on the transition from zero
	// to one active request, and only transitions away once all
	// active requests are complete. That means that ConnStateHook
	// cannot be used to do per-request work; ConnStateHook only notes
	// the overall state of the connection.
	ConnStateActive

	// ConnStateIdle represents a connection that has finished
	// handling a request and is in the keep-alive state, waiting
	// for a new request. Connections transition from ConnStateIdle
	// to either ConnStateActive or ConnStateClosed.
	ConnStateIdle

	// ConnStateHijacked represents a hijacked connection.
	// This is a terminal state. It does not transition to ConnStateClosed.
	ConnStateHijacked

	// ConnStateClosed represents a closed connection.
	// This is a terminal state. Hijacked connections do not
	// transition to ConnStateClosed.
	ConnStateClosed
)
