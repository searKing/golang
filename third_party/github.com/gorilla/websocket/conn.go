// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/searKing/golang/go/x/dispatch"
)

// maxInt64 is the effective "infinite" value for the Server and
// Transport's byte-limiting readers.
const maxInt64 = 1<<63 - 1

// A conn represents the server side of an HTTP connection.
type conn struct {
	// server is the server on which the connection arrived.
	// Immutable; never nil.
	server *Server

	// cancelCtx cancels the connection-level context.
	cancelCtx context.CancelFunc

	// rwc is the underlying network connection.
	// This is never wrapped by other types and is the value given out
	// to CloseNotifier callers. It is usually of type *net.TCPConn or
	// *tls.Conn.
	rwc *WebSocketConn

	// remoteAddr is rwc.RemoteAddr().String(). It is not populated synchronously
	// inside the Listener's Accept goroutine, as some implementations block.
	// It is populated immediately inside the (*conn).serve goroutine.
	// This is the value of a onMsgHandleHandler's (*Request).RemoteAddr.
	remoteAddr string

	// werr is set to the first write error to rwc.
	// It is set via checkConnErrorWriter{w}, where bufw writes.
	werr error

	curState struct{ atomic uint64 } // packed (unixtime<<8|uint8(ConnState))
}

// Close the connection.
func (c *conn) close() error {
	err := c.server.onCloseHandler.OnClose(checkConnErrorWebSocket{c: c})
	c.rwc.Close()
	return err
}

func (c *conn) setState(nc *WebSocketConn, state ConnState) {
	srv := c.server
	if state > 0xff || state < 0 {
		panic("internal error")
	}
	packedState := uint64(time.Now().Unix()<<8) | uint64(state)
	atomic.StoreUint64(&c.curState.atomic, packedState)
	if hook := srv.ConnState; hook != nil {
		hook(nc, state)
	}
}

func (c *conn) getState() (state ConnState, unixSec int64) {
	packedState := atomic.LoadUint64(&c.curState.atomic)
	return ConnState(packedState & 0xff), int64(packedState >> 8)
}

// ErrAbortHandler is a sentinel panic value to abort a handler.
// While any panic from OnHandshake aborts the response to the client,
// panicking with ErrAbortHandler also suppresses logging of a stack
// trace to the server's error log.
var ErrAbortHandler = errors.New("net/websocket: abort onMsgHandleHandler")
var errTooLarge = errors.New("websocket: read too large")

// isCommonNetReadError reports whether err is a common error
// encountered during reading a request off the network when the
// client has gone away or had its read fail somehow. This is used to
// determine which logs are interesting enough to log about.
func isCommonNetReadError(err error) bool {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		return true
	}
	return false
}

// onMsgReadHandler next request from connection.
func (c *conn) readRequest(ctx context.Context) (req any, err error) {

	var (
		wholeReqDeadline time.Time // or zero if none
		hdrDeadline      time.Time // or zero if none
	)
	t0 := time.Now()
	if d := c.server.readTimeout(); d != 0 {
		hdrDeadline = t0.Add(d)
	}
	if d := c.server.ReadTimeout; d != 0 {
		wholeReqDeadline = t0.Add(d)
	}
	c.rwc.SetReadDeadline(hdrDeadline)
	if d := c.server.WriteTimeout; d != 0 {
		defer func() {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}()
	}
	c.rwc.SetReadLimit(c.server.initialReadLimitSize())
	req, err = c.server.onMsgReadHandler.OnMsgRead(checkConnErrorWebSocket{c: c})
	if err != nil {
		if err == websocket.ErrReadLimit {
			return nil, errTooLarge
		}
		return nil, err
	}

	//c.rwc.setInfiniteReadLimit()
	c.rwc.SetReadLimit(maxInt64)

	// Adjust the read deadline if necessary.
	if !hdrDeadline.Equal(wholeReqDeadline) {
		c.rwc.SetReadDeadline(wholeReqDeadline)
	}

	return req, nil
}

// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	// handle close
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.server.logf("websocket: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		c.server.CheckError(c.rwc, c.close())
		c.setState(c.rwc, StateClosed)
	}()

	// WebSocket from here on.
	// cancel after this connection is handled
	ctx, cancelCtx := context.WithCancel(ctx)
	c.cancelCtx = cancelCtx
	defer cancelCtx()

	// read and handle the msg
	dispatch.NewDispatch(dispatch.ReaderFunc(func(ctx context.Context) (any, error) {
		msg, err := c.readRequest(ctx)
		if err = c.server.CheckError(c.rwc, err); err != nil {
			if isCommonNetReadError(err) {
				return nil, err // don't reply
			}
			return nil, err
		}
		c.setState(c.rwc, StateActive)
		return msg, nil
	}), dispatch.HandlerFunc(func(ctx context.Context, msg any) error {
		return c.server.CheckError(c.rwc, c.server.onMsgHandleHandler.OnMsgHandle(checkConnErrorWebSocket{c: c}, msg))
	})).WithContext(ctx).Start()
	return
}

type WebSocketReader interface {
	ReadMessage() (messageType int, p []byte, err error)
	//ReadJSON(v interface{}) error
}
type WebSocketWriter interface {
	//WriteControl(messageType int, data []byte, deadline time.Time) error
	//WritePreparedMessage(pm *websocket.PreparedMessage) error
	//WriteMessage(messageType int, data []byte) error
	WriteJSON(v any) error
}
type WebSocketCloser interface {
	Close() error
}

type WebSocketReadWriteCloser interface {
	WebSocketReader
	WebSocketWriter
	WebSocketCloser
}

// checkConnErrorWriter writes to c.rwc and records any write errors to c.werr.
// It only contains one field (and a pointer field at that), so it
// fits in an interface value without an extra allocation.
type checkConnErrorWebSocket struct {
	c *conn
	WebSocketReadWriteCloser
}

func (w checkConnErrorWebSocket) WriteControl(messageType int, data []byte, deadline time.Time) error {
	err := w.c.rwc.WriteControl(messageType, data, deadline)
	if err != nil && w.c.werr == nil {
		w.c.werr = err
		w.c.cancelCtx()
	}
	return err
}
func (w checkConnErrorWebSocket) WritePreparedMessage(pm *websocket.PreparedMessage) error {
	err := w.c.rwc.WritePreparedMessage(pm)
	if err != nil && w.c.werr == nil {
		w.c.werr = err
		w.c.cancelCtx()
	}
	return err
}
func (w checkConnErrorWebSocket) WriteMessage(messageType int, data []byte) error {
	err := w.c.rwc.WriteMessage(messageType, data)
	if err != nil && w.c.werr == nil {
		w.c.werr = err
		w.c.cancelCtx()
	}
	return err
}
func (w checkConnErrorWebSocket) WriteJSON(v any) (err error) {
	err = w.c.rwc.WriteJSON(v)
	if err != nil && w.c.werr == nil {
		w.c.werr = err
		w.c.cancelCtx()
	}
	return
}
func (w checkConnErrorWebSocket) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = w.c.rwc.ReadMessage()
	if err != nil && w.c.werr == nil {
		w.c.werr = err
		w.c.cancelCtx()
	}
	return messageType, p, err
}
func (w checkConnErrorWebSocket) ReadJSON(v any) error {
	err := w.c.rwc.ReadJSON(v)
	if err != nil && w.c.werr == nil {
		w.c.werr = err
		w.c.cancelCtx()
	}
	return err
}

func (w checkConnErrorWebSocket) Close() error {
	err := w.c.rwc.Close()
	if err != nil && w.c.werr == nil {
		w.c.werr = err
		w.c.cancelCtx()
	}
	return err
}
