// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketConn makes websocket concurrent safe
// see https://godoc.org/github.com/gorilla/websocket#hdr-Concurrency
type WebSocketConn struct {
	*websocket.Conn
	muRead  sync.Mutex
	muWrite sync.Mutex
}

func NewWebSocketConn(rw *websocket.Conn) *WebSocketConn {
	if rw == nil {
		panic("nil WebSocketConn")
	}
	return &WebSocketConn{
		Conn: rw,
	}
}
func (c *WebSocketConn) ReadMessage() (messageType int, p []byte, err error) {
	c.muRead.Lock()
	defer c.muRead.Unlock()
	return c.Conn.ReadMessage()
}
func (c *WebSocketConn) WriteJSON(v any) error {
	c.muWrite.Lock()
	defer c.muWrite.Unlock()
	return c.Conn.WriteJSON(v)
}
func (c *WebSocketConn) Close() error {
	c.muWrite.Lock()
	defer c.muWrite.Unlock()
	return c.Conn.Close()
}
