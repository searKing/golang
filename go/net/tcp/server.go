// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	slices_ "github.com/searKing/golang/go/exp/slices"
	time_ "github.com/searKing/golang/go/time"
)

type Handler interface {
	OnOpenHandler
	OnMsgReadHandler
	OnMsgHandleHandler
	OnCloseHandler
	OnErrorHandler
}

func NewServerFunc(
	onOpen OnOpenHandler,
	onMsgRead OnMsgReadHandler,
	onMsgHandle OnMsgHandleHandler,
	onClose OnCloseHandler,
	onError OnErrorHandler) *Server {
	return &Server{
		onOpenHandler:      slices_.FirstOrZero[OnOpenHandler](onOpen, NopOnOpenHandler),
		onMsgReadHandler:   slices_.FirstOrZero[OnMsgReadHandler](onMsgRead, NopOnMsgReadHandler),
		onMsgHandleHandler: slices_.FirstOrZero[OnMsgHandleHandler](onMsgHandle, NopOnMsgHandleHandler),
		onCloseHandler:     slices_.FirstOrZero[OnCloseHandler](onClose, NopOnCloseHandler),
		onErrorHandler:     slices_.FirstOrZero[OnErrorHandler](onError, NopOnErrorHandler),
	}
}
func NewServer(h Handler) *Server {
	return NewServerFunc(h, h, h, h, h)
}

type Server struct {
	Addr               string // TCP address to listen on, ":tcp" if empty
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
	listeners  map[*net.Listener]struct{}
	activeConn map[*conn]struct{}
	doneChan   chan struct{}
	onShutdown []func()

	// server state
	inShutdown atomic.Bool

	// ConnState specifies an optional callback function that is
	// called when a client connection changes state. See the
	// ConnState type and associated constants for details.
	ConnState func(net.Conn, ConnState)
}

func (srv *Server) CheckError(w io.Writer, r io.Reader, err error) error {
	if err == nil {
		return nil
	}
	return srv.onErrorHandler.OnError(w, r, err)
}

func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return srv.CheckError(nil, nil, ErrServerClosed)
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":tcp"
	}
	ln, err := net.Listen("tcp", addr)
	if srv.CheckError(nil, nil, err) != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

func (srv *Server) Serve(l net.Listener) error {
	l = &onceCloseListener{Listener: l}
	defer l.Close()

	// how long to sleep on accept failure
	var tempDelay = time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionInitialInterval(5*time.Millisecond),
		time_.WithExponentialBackOffOptionMaxInterval(time.Second),
		time_.WithExponentialBackOffOptionMaxElapsedDuration(-1),
		time_.WithExponentialBackOffOptionMaxElapsedCount(-1))
	ctx := context.WithValue(context.Background(), ServerContextKey, srv)
	for {
		rw, e := l.Accept()
		if e != nil {
			// return if server is cancaled, means normally close
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			// retry if it's recoverable
			var ne net.Error
			if errors.As(e, &ne) && ne.Temporary() {
				delay, ok := tempDelay.NextBackOff()
				if !ok {
					srv.logf("http: Accept error: %v; retried canceled as time exceed(%v)", e, tempDelay.GetMaxElapsedDuration())
					// return if timeout
					return srv.CheckError(nil, nil, e)
				}
				srv.logf("http: Accept error: %v; retrying in %v", e, delay)
				time.Sleep(delay)
				continue
			}
			// return otherwise
			return srv.CheckError(nil, nil, e)
		}
		tempDelay.Reset()

		// takeover this connect
		c := srv.newConn(rw)
		// Handle websocket On
		err := srv.onOpenHandler.OnOpen(c.rwc)
		if err = srv.CheckError(c.w, c.r, err); err != nil {
			c.close()
			return err
		}
		c.setState(c.rwc, StateNew) // before Serve can return
		go c.serve(ctx)
	}
}

func (srv *Server) trackConn(c *conn, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.activeConn == nil {
		srv.activeConn = make(map[*conn]struct{})
	}
	if add {
		srv.activeConn[c] = struct{}{}
	} else {
		delete(srv.activeConn, c)
	}
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
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

func ListenAndServe(addr string, readMsg OnMsgReadHandler, handleMsg OnMsgHandleHandler) error {
	server := &Server{Addr: addr, onMsgReadHandler: readMsg, onMsgHandleHandler: handleMsg}
	return server.ListenAndServe()
}
