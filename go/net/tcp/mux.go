// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"bufio"
	"io"
	"net"
	"sync"
)

type ServeMux struct {
	mu sync.RWMutex
	h  Handler
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{}
}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux

func (mux *ServeMux) OnOpen(conn net.Conn) error {
	return mux.h.OnOpen(conn)
}

func (mux *ServeMux) OnMsgRead(r io.Reader) (req any, err error) {
	return mux.h.OnMsgRead(r)
}

func (mux *ServeMux) OnMsgHandle(w io.Writer, msg any) error {
	return mux.h.OnMsgHandle(w, msg)
}
func (mux *ServeMux) OnClose(w io.Writer, r io.Reader) error {
	return mux.h.OnClose(w, r)
}
func (mux *ServeMux) OnError(w io.Writer, r io.Reader, err error) error {
	return mux.h.OnError(w, r, err)
}
func (mux *ServeMux) Handle(handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if handler == nil {
		panic("tcp: nil handler")
	}
	mux.h = handler
}
func (mux *ServeMux) handle() Handler {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	if mux.h == nil {
		return NotFoundHandler()
	}
	return mux.h
}
func NotFoundHandler() Handler { return &NotFound{} }

var _ Handler = (*NotFound)(nil)

type NotFound struct {
	NopServer
}

func (notfound *NotFound) ReadMsg(b *bufio.Reader) (msg any, err error) {
	return nil, ErrNotFound
}
func (notfound *NotFound) HandleMsg(b *bufio.Writer, msg any) error {
	return ErrServerClosed
}
