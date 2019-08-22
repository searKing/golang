// Copyright 2016 The CMux Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package cmux

import (
	"context"
	"crypto/tls"
	net_ "github.com/searKing/golang/go/net"
	"github.com/searKing/golang/go/strings"
	"log"
	"net"
	"sync"
	"time"
)

// for readability of sniffTimeout
var noTimeout time.Duration
var noTimeoutDeadline time.Time

// New instantiates a new connection multiplexer.
func New(parent context.Context) CMux {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &cMux{
		ctx:          ctx,
		cancel:       cancel,
		maxIdleConns: 1024,
		errHandler:   ignoreErrorHandler,
	}
}

// CMux is a multiplexer for network connections.
type CMux interface {
	// Match returns a net.Listener that accepts only the
	// connections that matched by at least of the matcher writers.
	//
	// The order used to call Match determines the priority of matchers.
	Match(matchers ...Matcher) net.Listener

	// MatchAndGoServe starts a server with a listener that accepts only the
	// connections that matched by at least of the matcher writers.
	//
	// The order used to call Match determines the priority of matchers.
	MatchAndGoServe(handler Handler)

	// Serve starts multiplexing the listener. Serve blocks and perhaps
	// should be invoked concurrently within a go routine.
	Serve(l net.Listener) error

	// ListenAndServe listens on the TCP network address addr and then calls
	// Serve with handler to handle requests on incoming connections.
	// Accepted connections are configured to enable TCP keep-alives.
	//
	// The handler is typically nil, in which case the DefaultServeMux is used.
	//
	// ListenAndServe always returns a non-nil error.
	ListenAndServe(addr string) error

	// ListenAndServeTLS listens on the TCP network address srv.Addr and
	// then calls ServeTLS to handle requests on incoming TLS connections.
	// Accepted connections are configured to enable TCP keep-alives.
	//
	// Filenames containing a certificate and matching private key for the
	// server must be provided if neither the Server's TLSConfig.Certificates
	// nor TLSConfig.GetCertificate are populated. If the certificate is
	// signed by a certificate authority, the certFile should be the
	// concatenation of the server's certificate, any intermediates, and
	// the CA's certificate.
	//
	// If srv.Addr is blank, ":https" is used.
	//
	// ListenAndServeTLS always returns a non-nil error. After Shutdown or
	// Close, the returned error is ErrServerClosed.
	ListenAndServeTLS(addr string, tlsConfig *tls.Config, certFile, keyFile string) error

	// Close immediately closes all active net.Listeners and any
	// connections in state StateNew, StateActive, or StateIdle. For a
	// graceful shutdown, use Shutdown.
	//
	// Close does not attempt to close (and does not even know about)
	// any hijacked connections, such as WebSockets.
	//
	// Close returns any error returned from closing the Server's
	// underlying Listener(s).
	Close() error

	// Shutdown gracefully shuts down the server without interrupting any
	// active connections. Shutdown works by first closing all open
	// listeners, then closing all idle connections, and then waiting
	// indefinitely for connections to return to idle and then shut down.
	// If the provided context expires before the shutdown is complete,
	// Shutdown returns the context's error, otherwise it returns any
	// error returned from closing the Server's underlying Listener(s).
	//
	// When Shutdown is called, Serve, ListenAndServe, and
	// ListenAndServeTLS immediately return ErrServerClosed. Make sure the
	// program doesn't exit and waits instead for Shutdown to return.
	//
	// Shutdown does not attempt to close nor wait for hijacked
	// connections such as WebSockets. The caller of Shutdown should
	// separately notify such long-lived connections of shutdown and wait
	// for them to close, if desired. See RegisterOnShutdown for a way to
	// register shutdown notification functions.
	//
	// Once Shutdown has been called on a server, it may not be reused;
	// future calls to methods such as Serve will return ErrServerClosed.
	Shutdown(ctx context.Context) error
	// HandleError registers an error handler that handles listener errors.
	HandleError(ErrorHandler)
	// sets a timeout for the read of matchers
	SetReadTimeout(time.Duration)
}

type matchersListener struct {
	matchers []Matcher
	l        *muxListener
}

type cMux struct {
	root              net.Listener
	maxIdleConns      int
	errHandler        ErrorHandler
	matchersListeners []matchersListener
	sniffTimeout      time.Duration

	// ConnStateHook specifies an optional callback function that is
	// called when a client connection changes state. See the
	// ConnStateHook type and associated constants for details.
	ConnStateHook func(net.Conn, ConnState)
	// ErrorLog specifies an optional logger for errors accepting
	// connections, unexpected behavior from handlers, and
	// underlying FileSystem errors.
	// If nil, logging is done via the log package's standard logger.
	errorLog *log.Logger
	ctx      context.Context
	cancel   func()

	mu         sync.Mutex
	listeners  map[*net.Listener]struct{}
	activeConn map[*conn]struct{}
	onShutdown []func()
}

func (m *cMux) Match(matchers ...Matcher) net.Listener {
	m.mu.Lock()
	defer m.mu.Unlock()
	ml := newMuxListener(m.ctx, m.listeners, m.maxIdleConns)
	m.matchersListeners = append(m.matchersListeners, matchersListener{matchers: matchers, l: ml})
	return ml
}

func (m *cMux) MatchAndGoServe(handler Handler) {
	if handler == nil {
		return
	}
	go handler.Serve(m.Match(handler))
}

// Context returns the request's context. To change the context, use
// WithContext.
//
// The returned context is always non-nil; it defaults to the
// background context.
func (m *cMux) Context() context.Context {
	return m.ctx
}

func (m *cMux) SetReadTimeout(t time.Duration) {
	m.sniffTimeout = t
}

func (m *cMux) logf(format string, args ...interface{}) {
	if m.errorLog != nil {
		m.errorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// Create new connection from rwc.
func (m *cMux) newConn(rwc net.Conn) *conn {
	return &conn{
		server: m,
		muc:    newMuxConn(rwc),
	}
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines read requests and
// then call srv.Handler to reply to them.
func (m *cMux) Serve(l net.Listener) error {
	if m.shuttingDown() {
		return ErrServerClosed
	}

	l = net_.OnceCloseListener(l)
	defer l.Close()

	if !m.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer m.trackListener(&l, false)

	var tempDelay time.Duration // how long to sleep on accept failure
	ctx := context.WithValue(m.Context(), ServerContextKey, m)

	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return ErrServerClosed
			default:
			}
			if !m.handleErr(err) {
				return err
			}

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				m.logf("cmux: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		c := m.newConn(rw)
		c.setState(c.muc, ConnStateNew) // before Serve can return

		go c.serve(ctx)

	}
}

// ServeTLS accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines perform TLS
// setup and then read requests, calling srv.Handler to reply to them.
//
// Files containing a certificate and matching private key for the
// server must be provided if neither the Server's
// TLSConfig.Certificates nor TLSConfig.GetCertificate are populated.
// If the certificate is signed by a certificate authority, the
// certFile should be the concatenation of the server's certificate,
// any intermediates, and the CA's certificate.
//
// ServeTLS always returns a non-nil error. After Shutdown or Close, the
// returned error is ErrServerClosed.
func (m *cMux) ServeTLS(l net.Listener, tLSConfig *tls.Config, certFile, keyFile string) error {
	config := cloneTLSConfig(tLSConfig)
	if !strings.SliceContains(config.NextProtos, "http/1.1") {
		config.NextProtos = append(config.NextProtos, "http/1.1")
	}

	configHasCert := len(config.Certificates) > 0 || config.GetCertificate != nil
	if !configHasCert || certFile != "" || keyFile != "" {
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
	}

	tlsListener := tls.NewListener(l, config)
	return m.Serve(tlsListener)
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// If srv.Addr is blank, ":http" is used.
//
// ListenAndServe always returns a non-nil error. After Shutdown or Close,
// the returned error is ErrServerClosed.
func (m *cMux) ListenAndServe(addr string) error {
	if m.shuttingDown() {
		return ErrServerClosed
	}
	if addr == "" {
		addr = ":tcp"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return m.Serve(ln)
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls ServeTLS to handle requests on incoming TLS connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Filenames containing a certificate and matching private key for the
// server must be provided if neither the Server's TLSConfig.Certificates
// nor TLSConfig.GetCertificate are populated. If the certificate is
// signed by a certificate authority, the certFile should be the
// concatenation of the server's certificate, any intermediates, and
// the CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
//
// ListenAndServeTLS always returns a non-nil error. After Shutdown or
// Close, the returned error is ErrServerClosed.
func (m *cMux) ListenAndServeTLS(addr string, tlsConfig *tls.Config, certFile, keyFile string) error {
	if m.shuttingDown() {
		return ErrServerClosed
	}
	if addr == "" {
		addr = ":https"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return m.ServeTLS(net_.TcpKeepAliveListener(ln.(*net.TCPListener), 3*time.Minute), tlsConfig, certFile, keyFile)
}

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(ctx context.Context, addr string, handler Handler) error {
	server := New(ctx)
	server.MatchAndGoServe(handler)
	return server.ListenAndServe(addr)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects HTTPS connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
func ListenAndServeTLS(ctx context.Context, addr string, handler Handler, tlsConfig *tls.Config, certFile, keyFile string) error {
	server := New(ctx)
	server.MatchAndGoServe(handler)
	return server.ListenAndServeTLS(addr, tlsConfig, certFile, keyFile)
}

// Close immediately closes all active net.Listeners and any
// connections in state StateNew, StateActive, or StateIdle. For a
// graceful shutdown, use Shutdown.
//
// Close does not attempt to close (and does not even know about)
// any hijacked connections, such as WebSockets.
//
// Close returns any error returned from closing the Server's
// underlying Listener(s).
func (m *cMux) Close() error {
	if m.shuttingDown() {
		return ErrServerClosed
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cancel()

	err := m.closeListenersLocked()
	for c := range m.activeConn {
		c.close()
		delete(m.activeConn, c)
	}
	return err
}

// shutdownPollInterval is how often we poll for quiescence
// during Server.Shutdown. This is lower during tests, to
// speed up tests.
// Ideally we could find a solution that doesn't involve polling,
// but which also doesn't have a high runtime cost (and doesn't
// involve any contentious mutexes), but that is left as an
// exercise for the reader.
var shutdownPollInterval = 500 * time.Millisecond

// Shutdown gracefully shuts down the server without interrupting any
// active connections. Shutdown works by first closing all open
// listeners, then closing all idle connections, and then waiting
// indefinitely for connections to return to idle and then shut down.
// If the provided context expires before the shutdown is complete,
// Shutdown returns the context's error, otherwise it returns any
// error returned from closing the Server's underlying Listener(s).
//
// When Shutdown is called, Serve, ListenAndServe, and
// ListenAndServeTLS immediately return ErrServerClosed. Make sure the
// program doesn't exit and waits instead for Shutdown to return.
//
// Shutdown does not attempt to close nor wait for hijacked
// connections such as WebSockets. The caller of Shutdown should
// separately notify such long-lived connections of shutdown and wait
// for them to close, if desired. See RegisterOnShutdown for a way to
// register shutdown notification functions.
//
// Once Shutdown has been called on a server, it may not be reused;
// future calls to methods such as Serve will return ErrServerClosed.
func (m *cMux) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	lnerr := m.closeListenersLocked()
	m.cancel()

	for _, f := range m.onShutdown {
		go f()
	}
	m.mu.Unlock()

	ticker := time.NewTicker(shutdownPollInterval)
	defer ticker.Stop()
	for {
		if m.closeIdleConns() {
			return lnerr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// RegisterOnShutdown registers a function to call on Shutdown.
// This can be used to gracefully shutdown connections that have
// undergone NPN/ALPN protocol upgrade or that have been hijacked.
// This function should start protocol-specific graceful shutdown,
// but should not wait for shutdown to complete.
func (m *cMux) RegisterOnShutdown(f func()) {
	m.mu.Lock()
	m.onShutdown = append(m.onShutdown, f)
	m.mu.Unlock()
}

func (m *cMux) HandleError(h ErrorHandler) {
	m.errHandler = h
}

func (m *cMux) handleErr(err error) bool {
	if m.errHandler == nil {
		return true
	}
	if !m.errHandler.Continue(err) {
		return false
	}

	if ne, ok := err.(net.Error); ok {
		return ne.Temporary()
	}

	return false
}

// trackListener adds or removes a net.Listener to the set of tracked
// listeners.
//
// We store a pointer to interface in the map set, in case the
// net.Listener is not comparable. This is safe because we only call
// trackListener via Serve and can track+defer untrack the same
// pointer to local variable there. We never need to compare a
// Listener from another caller.
//
// It reports whether the server is still up (not Shutdown or Closed).
func (m *cMux) trackListener(ln *net.Listener, add bool) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.listeners == nil {
		m.listeners = make(map[*net.Listener]struct{})
	}
	if add {
		if m.shuttingDown() {
			return false
		}
		m.listeners[ln] = struct{}{}
	} else {
		delete(m.listeners, ln)
	}
	return true
}

func (m *cMux) trackConn(c *conn, add bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeConn == nil {
		m.activeConn = make(map[*conn]struct{})
	}
	if add {
		m.activeConn[c] = struct{}{}
	} else {
		delete(m.activeConn, c)
	}
}

func (m *cMux) shuttingDown() bool {
	select {
	case <-m.ctx.Done():
		return true
	default:
		return false
	}
}

// closeIdleConns closes all idle connections and reports whether the
// server is quiescent.
func (m *cMux) closeIdleConns() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	quiescent := true
	for c := range m.activeConn {
		st, unixSec := c.getState()
		// Issue 22682: treat StateNew connections as if
		// they're idle if we haven't read the first request's
		// header in over 5 seconds.
		if st == ConnStateNew && unixSec < time.Now().Unix()-5 {
			st = ConnStateIdle
		}
		if st != ConnStateIdle || unixSec == 0 {
			// Assume unixSec == 0 means it's a very new
			// connection, without state set yet.
			quiescent = false
			continue
		}
		c.close()
		delete(m.activeConn, c)
	}
	return quiescent
}

func (m *cMux) closeListenersLocked() error {
	var err error
	for ln := range m.listeners {
		if cerr := (*ln).Close(); cerr != nil && err == nil {
			err = cerr
		}
		delete(m.listeners, ln)
	}
	return err
}
