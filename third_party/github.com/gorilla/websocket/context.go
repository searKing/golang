// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

var (
	// ServerContextKey is a context key. It can be used in Websocket
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of
	// type *Server.
	ServerContextKey = &contextKey{"websocket-server"}
	// ClientContextKey is a context key. It can be used in Websocket
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of
	// type *Client.
	ClientContextKey = &contextKey{"websocket-client"}

	// LocalAddrContextKey is a context key. It can be used in
	// TCP handlers with context.WithValue to access the local
	// address the connection arrived on.
	// The associated value will be of type net.Addr.
	LocalAddrContextKey = &contextKey{"local-addr"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "net/http context value " + k.name }
