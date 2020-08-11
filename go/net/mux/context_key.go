// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

var (
	// ServerContextKey is a context key. It can be used in cmux
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of
	// type CMuxer.
	ServerContextKey = &contextKey{"cmux-server"}

	// LocalAddrContextKey is a context key. It can be used in
	// cmux handlers with context.WithValue to access the local
	// address the connection arrived on.
	// The associated value will be of type net.Addr.
	LocalAddrContextKey = &contextKey{"local-addr"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "net/mux context value " + k.name }
