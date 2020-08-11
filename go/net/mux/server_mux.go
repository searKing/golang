// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/searKing/golang/go/error/multi"
	io_ "github.com/searKing/golang/go/io"
	net_ "github.com/searKing/golang/go/net"
)

// NotFound replies to the request with an HTTP 404 not found error.
func NotFound(c net.Conn) {
	_, _ = c.Write([]byte("404 page not found"))
}

// NotFoundHandler returns a simple request handler
// that replies to each request with a ``404 page not found'' reply.
func NotFoundHandler() HandlerConn { return HandlerConnFunc(NotFound) }

// MuxListener is a net.ServeMux that accepts only the connections that matched.
// goroutine unsafe
type ServeMux struct {
	// NotFound replies to the listener with a not found error.
	NotFoundHandler HandlerConn
	sniffTimeout    time.Duration

	mu sync.RWMutex
	m  []muxEntry
}

type muxEntry struct {
	h HandlerConn
	l *net_.NotifyListener

	pattern Matcher
}

func (e muxEntry) Serve(c net.Conn) {
	if e.h != nil {
		e.h.Serve(c)
		return
	}
	if e.l != nil {
		e.l.C <- c
		return
	}
	panic("mux_entry: nil handler")
}

type Matcher interface {
	Match(io.Writer, io.Reader) bool
}

// MatchWriter is a match that can also write response (say to do handshake).
type MatcherFunc func(io.Writer, io.Reader) bool

func (f MatcherFunc) Match(w io.Writer, r io.Reader) bool {
	return f(w, r)
}

func MatcherAny(matchers ...Matcher) Matcher {
	return MatcherFunc(func(w io.Writer, r io.Reader) bool {
		sniffReader := io_.SniffReader(r)
		for _, pattern := range matchers {
			sniffReader.Sniff(true)
			if pattern.Match(w, sniffReader) {
				sniffReader.Sniff(false)
				return true
			}
			sniffReader.Sniff(false)
		}
		return false
	})
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		NotFoundHandler: HandlerConnFunc(func(net.Conn) {}),
	}
}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux

// SetReadTimeout sets a timeout for the read of matchers
func (mux *ServeMux) SetReadTimeout(t time.Duration) {
	mux.sniffTimeout = t
}

func (mux *ServeMux) HandleListener(pattern Matcher) net.Listener {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == nil {
		panic("listener: invalid pattern")
	}

	e := muxEntry{l: net_.NewNotifyListener(), pattern: pattern}
	mux.m = append(mux.m, e)
	return e.l
}

func (mux *ServeMux) Handle(pattern Matcher, handler HandlerConn) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == nil {
		panic("listener: invalid pattern")
	}
	if handler == nil {
		panic("listener: nil handler")
	}

	e := muxEntry{h: handler, pattern: pattern}
	mux.m = append(mux.m, e)
	return
}

// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern Matcher, handler func(net.Conn)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerConnFunc(handler))
}

// Find a handler on a handler map.
func (mux *ServeMux) match(c *sniffConn) (h HandlerConn) {
	for _, e := range mux.m {
		c.startSniffing()
		if e.pattern.Match(c, c) {
			c.doneSniffing()
			return e
		}
		c.doneSniffing()
	}
	return nil

}

// handler is the main implementation of HandlerConn.
func (mux *ServeMux) Handler(c *sniffConn) (h HandlerConn) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	// set sniff timeout
	if mux.sniffTimeout > noTimeout {
		_ = c.SetReadDeadline(time.Now().Add(mux.sniffTimeout))
	}
	h = mux.match(c)

	// unset sniff timeout
	if mux.sniffTimeout > noTimeout {
		_ = c.SetReadDeadline(noTimeoutDeadline)
	}
	if h == nil {
		notFoundHandler := mux.NotFoundHandler
		if notFoundHandler == nil {
			notFoundHandler = NotFoundHandler()
		}
		h = muxEntry{
			h: notFoundHandler,
		}
	}
	return
}

func (mux *ServeMux) Serve(c net.Conn) {
	muxC, ok := c.(*sniffConn)
	if !ok {
		muxC = newMuxConn(c)
	}

	h := mux.Handler(muxC)
	h.Serve(c)
}

func (mux *ServeMux) Close() error {
	var errors []error
	for _, e := range mux.m {
		if l := e.l; l != nil {
			errors = append(errors, l.Close())
		}
	}
	return multi.New(errors...)
}

func HandleListener(pattern Matcher) net.Listener {
	return DefaultServeMux.HandleListener(pattern)
}

func Handle(pattern Matcher, handler HandlerConn) {
	DefaultServeMux.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern.
func HandleFunc(pattern Matcher, handler func(net.Conn)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
func SetReadTimeout(t time.Duration) {
	DefaultServeMux.SetReadTimeout(t)
}
