// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmux

import (
	"context"
	"errors"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/searKing/golang/go/sync/atomic"
)

var (
	// ErrHijacked is returned by ResponseWriter.Write calls when
	// the underlying connection has been hijacked using the
	// Hijacker interface. A zero-byte write on a hijacked
	// connection will return ErrHijacked without any other side
	// effects.
	ErrHijacked = errors.New("cmux: connection has been hijacked")
)

type conn struct {
	// server is the server on which the connection arrived.
	// Immutable; never nil.
	server *cMux

	// cancelCtx cancels the connection-level context.
	cancelCtx context.CancelFunc

	// rwc is the underlying network connection.
	muc *muxConn

	// remoteAddr is rwc.RemoteAddr().String(). It is not populated synchronously
	// inside the Listener's Accept goroutine, as some implementations block.
	// It is populated immediately inside the (*conn).serve goroutine.
	// This is the value of a Handler's (*Request).RemoteAddr.
	remoteAddr string

	curPacketState atomic.Uint64 // packed (unixtime<<8|uint8(ConnStateHook))

	// mu guards hijackedv
	mu sync.Mutex

	// hijackedv is whether this connection has been hijacked
	// by a Handler with the Hijacker interface.
	// It is guarded by mu.
	hijackedv bool
}

func (c *conn) hijacked() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hijackedv
}

// c.mu must be held.
func (c *conn) hijackLocked() (rwc net.Conn, err error) {
	if c.hijackedv {
		return nil, ErrHijacked
	}

	c.hijackedv = true
	rwc = c.muc
	_ = rwc.SetDeadline(noTimeoutDeadline)
	c.muc.doneSniffing()

	//buf = bufio.NewReadWriter(bufio.NewReader(c.muc), bufio.NewWriter(rwc))
	c.setState(rwc, ConnStateHijacked)
	return
}

// Close the connection.
func (c *conn) close() {
	_ = c.muc.Close()
}

// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.muc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.muc.LocalAddr())
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.server.logf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		if !c.hijacked() {
			c.close()
			c.setState(c.muc, ConnStateClosed)
		}
	}()

	ctx, cancelCtx := context.WithCancel(ctx)
	c.cancelCtx = cancelCtx
	defer cancelCtx()

	c.setState(c.muc, ConnStateActive)

	// set sniff timeout
	if c.server.sniffTimeout > noTimeout {
		_ = c.muc.SetReadDeadline(time.Now().Add(c.server.sniffTimeout))
	}
	for _, ml := range c.server.matchersListeners {
		// skip if listener is closed
		if ml.l.shuttingDown() {
			continue
		}
		for _, s := range ml.matchers {
			matched := s.Match(c.muc, c.muc.startSniffing())
			c.muc.doneSniffing()
			if matched {
				// unset sniff timeout
				if c.server.sniffTimeout > noTimeout {
					_ = c.muc.SetReadDeadline(noTimeoutDeadline)
				}
				rwc, err := c.hijackLocked()
				if err != nil {
					return
				}
				if err := ml.l.Notify(ctx, rwc); err != nil {
					_ = rwc.Close()
				}
				return
			}

			// resniff
			c.muc = newMuxConn(c.muc)
		}
	}

	_ = c.close
	err := ErrNotMatched{c: c.muc}
	c.server.handleErr(err)

	c.setState(c.muc, ConnStateIdle)
}

func (c *conn) setState(nc net.Conn, state ConnState) {
	srv := c.server
	switch state {
	case ConnStateNew:
		srv.trackConn(c, true)
	case ConnStateHijacked, ConnStateClosed:
		srv.trackConn(c, false)
	}
	if state > 0xff || state < 0 {
		panic("internal error")
	}
	packedState := uint64(time.Now().Unix()<<8) | uint64(state)
	c.curPacketState.Store(packedState)
	if hook := srv.ConnStateHook; hook != nil {
		hook(nc, state)
	}
}

func (c *conn) getState() (state ConnState, unixSec int64) {
	packedState := c.curPacketState.Load()
	return ConnState(packedState & 0xff), int64(packedState >> 8)
}
