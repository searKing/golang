// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	slices_ "github.com/searKing/golang/go/exp/slices"
	"github.com/searKing/golang/go/sync/atomic"
)

type ServerHandler interface {
	OnHandshakeHandler
	OnOpenHandler
	OnMsgReadHandler   //Block
	OnMsgHandleHandler //Unblock
	OnCloseHandler
	OnErrorHandler
}

func NewServerFunc(onHandshake OnHandshakeHandler,
	onOpen OnOpenHandler,
	onMsgRead OnMsgReadHandler,
	onMsgHandle OnMsgHandleHandler,
	onClose OnCloseHandler,
	onError OnErrorHandler) *Server {
	return &Server{
		onHandshakeHandler: slices_.FirstOrZero[OnHandshakeHandler](onHandshake, NopOnHandshakeHandler),
		onOpenHandler:      slices_.FirstOrZero[OnOpenHandler](onOpen, NopOnOpenHandler),
		onMsgReadHandler:   slices_.FirstOrZero[OnMsgReadHandler](onMsgRead, NopOnMsgReadHandler),
		onMsgHandleHandler: slices_.FirstOrZero[OnMsgHandleHandler](onMsgHandle, NopOnMsgHandleHandler),
		onCloseHandler:     slices_.FirstOrZero[OnCloseHandler](onClose, NopOnCloseHandler),
		onErrorHandler:     slices_.FirstOrZero[OnErrorHandler](onError, NopOnErrorHandler),
	}
}
func NewServer(h ServerHandler) *Server {
	return NewServerFunc(h, h, h, h, h, h)
}

type Server struct {
	upgrader           websocket.Upgrader // use default options
	onHandshakeHandler OnHandshakeHandler
	onOpenHandler      OnOpenHandler
	onMsgReadHandler   OnMsgReadHandler
	onMsgHandleHandler OnMsgHandleHandler
	onCloseHandler     OnCloseHandler
	onErrorHandler     OnErrorHandler

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	MaxBytes     int

	ErrorLog *log.Logger

	mu         sync.Mutex
	activeConn *conn
	onShutdown []func()

	// server state
	disableKeepAlives atomic.Bool // accessed atomically.
	inShutdown        atomic.Bool
	// ConnState specifies an optional callback function that is
	// called when a client connection changes state. See the
	// ConnState type and associated constants for details.
	ConnState func(*WebSocketConn, ConnState)
}

func (srv *Server) CheckError(conn WebSocketReadWriteCloser, err error) error {
	if err == nil {
		return nil
	}
	return srv.onErrorHandler.OnError(conn, err)
}

// OnHandshake takes over the http handler
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	// transfer http to websocket
	srv.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := srv.upgrader.Upgrade(w, r, nil)
	if srv.CheckError(nil, err) != nil {
		return err
	}
	defer ws.Close()
	ctx := context.WithValue(context.Background(), ServerContextKey, srv)
	// Handle HTTP Handshake
	err = srv.onHandshakeHandler.OnHandshake(w, r)
	if srv.CheckError(nil, err) != nil {
		return err
	}
	// takeover the connect
	c := srv.newConn(ws)
	// Handle websocket On
	err = srv.onOpenHandler.OnOpen(c.rwc)
	if err = srv.CheckError(c.rwc, err); err != nil {
		c.close()
		return err
	}
	c.setState(c.rwc, StateNew) // before Serve can return

	c.serve(ctx)
	return nil
}

// Create new connection from rwc.
func (srv *Server) newConn(wc *websocket.Conn) *conn {
	c := &conn{
		server: srv,
		rwc: &WebSocketConn{
			Conn: wc,
		},
	}
	return c
}
func (srv *Server) logf(format string, args ...any) {
	if srv.ErrorLog != nil {
		srv.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
