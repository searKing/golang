// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmux

import (
	"context"
	"net"
	"strings"
)

// muxListener receives conn if matched
// goroutine unsafe
type muxListener struct {
	listeners multiAddrs
	connC     chan net.Conn

	ctx    context.Context
	cancel func()
}

func newMuxListener(parent context.Context, ls map[*net.Listener]struct{}, maxIdleConns int) *muxListener {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	return &muxListener{
		ctx:       ctx,
		cancel:    cancel,
		listeners: ls,
		connC:     make(chan net.Conn, maxIdleConns),
	}
}

func (ml *muxListener) Context() context.Context {
	if ml.ctx != nil {
		return ml.ctx
	}
	return context.Background()
}

// Notify transfers conn to listener.Accept's callers.
// drops conn when error non nil
func (ml *muxListener) Notify(ctx context.Context, conn net.Conn) error {
	// exit as early as possible
	if ml.shuttingDown() {
		return ErrListenerClosed
	}
	select {
	case <-ml.Context().Done():
		// Already closed. Don't Notify again.
		return ErrListenerClosed
	case <-ctx.Done():
		// Already closed. Don't Notify again.
		return ErrListenerClosed
	case ml.connC <- conn:
		return nil
	}
}

// Addr returns the listener's network address.
func (ml *muxListener) Addr() net.Addr {
	return ml.listeners
}

func (ml *muxListener) Accept() (net.Conn, error) {
	// Already closed. Don't close again.
	if ml.shuttingDown() {
		return nil, ErrListenerClosed
	}
	select {
	case <-ml.Context().Done():
		// Already closed. Don't Accept again.
		return nil, ErrListenerClosed
	case c, ok := <-ml.connC:
		if !ok {
			// Already closed. Don't Accept again.
			return nil, ErrListenerClosed
		}
		return c, nil
	}
}

func (ml *muxListener) Close() error {
	// Already closed. Don't close again.
	if ml.shuttingDown() {
		return ErrListenerClosed
	}
	ml.cancel()
	// Drain the connections enqueued for the listener.
L:
	for {
		select {
		case c, ok := <-ml.connC:
			if !ok {
				// already closed
				return nil
			}
			_ = c.Close()
		default:
			break L
		}
	}
	select {
	case <-ml.connC:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(ml.connC)
	}

	return nil
}

func (ml *muxListener) shuttingDown() bool {
	select {
	case <-ml.ctx.Done():
		return true
	default:
		return false
	}
}

type multiAddrs map[*net.Listener]struct{}

func (ls multiAddrs) Network() string {
	var networkStrs []string
	for l, _ := range ls {
		networkStrs = append(networkStrs, (*l).Addr().Network())
	}
	return strings.Join(networkStrs, ",")
}

func (ls multiAddrs) String() string {
	var addrStrs []string
	for l, _ := range ls {
		addrStrs = append(addrStrs, (*l).Addr().String())
	}
	return strings.Join(addrStrs, ",")
}
