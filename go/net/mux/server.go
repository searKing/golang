// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
	"time"

	net_ "github.com/searKing/golang/go/net"
	http_ "github.com/searKing/golang/go/net/http"
	"github.com/searKing/golang/go/strings"
	"github.com/searKing/golang/go/sync/atomic"
)

// for readability of sniffTimeout
var noTimeout time.Duration
var noTimeoutDeadline time.Time

// NewServer wraps NewServerWithContext using the background context.
//
//go:generate go-option -type "Server"
func NewServer() *Server {
	return NewServerWithContext(context.Background())
}

// NewServerWithContext returns a new Server.
func NewServerWithContext(ctx context.Context) *Server {
	return &Server{
		ctx:          ctx,
		maxIdleConns: 1024,
		errHandler:   ignoreErrorHandler,
	}
}

// Server is a multiplexer for network connections.
type Server struct {
	Handler HandlerConn // handler to invoke, mux.DefaultServeMux if nil

	// maxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	maxIdleConns int
	errHandler   ErrorHandler

	// ConnStateHook specifies an optional callback function that is
	// called when a client connection changes state. See the
	// ConnStateHook type and associated constants for details.
	ConnStateHook func(net.Conn, ConnState)
	// ErrorLog specifies an optional logger for errors accepting
	// connections, unexpected behavior from handlers, and
	// underlying FileSystem errors.
	// If nil, logging is done via the log package's standard logger.
	errorLog   *log.Logger
	ctx        context.Context
	inShutdown atomic.Bool // accessed atomically (non-zero means we're in Shutdown)

	mu         sync.Mutex
	listeners  map[*net.Listener]struct{}
	activeConn map[*conn]struct{}
	doneChan   chan struct{}
	onShutdown []func()
}

func (srv *Server) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.getDoneChanLocked()
}

func (srv *Server) getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

func (srv *Server) closeDoneChanLocked() {
	ch := srv.getDoneChanLocked()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(ch)
	}
}

// Context returns the request's context. To change the context, use
// WithContext.
//
// The returned context is always non-nil; it defaults to the
// background context.
func (srv *Server) Context() context.Context {
	if srv.ctx != nil {
		return srv.ctx
	}
	return context.Background()
}

func (srv *Server) logf(format string, args ...any) {
	if srv.errorLog != nil {
		srv.errorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) *conn {
	return &conn{
		server: srv,
		muc:    newMuxConn(rwc),
	}
}

// Serve starts multiplexing the listener. Serve blocks and perhaps
// should be invoked concurrently within a go routine.
// Serve accepts incoming connections on the ServeMux l, creating a
// new service goroutine for each. The service goroutines read requests and
// then call srv.HandlerConn to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	l = net_.OnceCloseListener(l)
	defer l.Close()

	if srv.shuttingDown() {
		return ErrServerClosed
	}

	if !srv.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer srv.trackListener(&l, false)

	var tempDelay time.Duration // how long to sleep on accept failure
	ctx := context.WithValue(srv.Context(), ServerContextKey, srv)

	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return ErrServerClosed
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			if !srv.handleErr(err) {
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
				srv.logf("cmux: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		c := srv.newConn(rw)
		c.setState(c.muc, ConnStateNew) // before Serve can return

		go c.serve(ctx)

	}
}

// ServeTLS accepts incoming connections on the ServeMux l, creating a
// new service goroutine for each. The service goroutines perform TLS
// setup and then read requests, calling srv.HandlerConn to reply to them.
//
// Files containing a certificate and matching private key for the
// server must be provided if neither the ServeMux's
// TLSConfig.Certificates nor TLSConfig.GetCertificate are populated.
// If the certificate is signed by a certificate authority, the
// certFile should be the concatenation of the server's certificate,
// any intermediates, and the CA's certificate.
//
// ServeTLS always returns a non-nil error. After Shutdown or Close, the
// returned error is ErrServerClosed.
func (srv *Server) ServeTLS(l net.Listener, tLSConfig *tls.Config, certFile, keyFile string) error {
	config := http_.CloneTLSConfig(tLSConfig)
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
	return srv.Serve(tlsListener)
}

// serverHandler delegates to either the server's HandlerConn or
// DefaultServeMux and also handles "OPTIONS *" requests.
type serverHandler struct {
	srv *Server
}

func (sh serverHandler) Handler(c *sniffConn) HandlerConn {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	return handler
}

func (sh serverHandler) handler() HandlerConn {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	return handler
}
func (sh serverHandler) Serve(c net.Conn) {
	sh.handler().Serve(c)
}

func (sh serverHandler) Close() error {
	if closer, ok := sh.handler().(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// If srv.Addr is blank, ":http" is used.
//
// ListenAndServe always returns a non-nil error. After Shutdown or Close,
// the returned error is ErrServerClosed.
func (srv *Server) ListenAndServe(addr string) error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	if addr == "" {
		addr = ":tcp"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls ServeTLS to handle requests on incoming TLS connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Filenames containing a certificate and matching private key for the
// server must be provided if neither the ServeMux's TLSConfig.Certificates
// nor TLSConfig.GetCertificate are populated. If the certificate is
// signed by a certificate authority, the certFile should be the
// concatenation of the server's certificate, any intermediates, and
// the CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
//
// ListenAndServeTLS always returns a non-nil error. After Shutdown or
// Close, the returned error is ErrServerClosed.
func (srv *Server) ListenAndServeTLS(addr string, tlsConfig *tls.Config, certFile, keyFile string) error {
	if srv.shuttingDown() {
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

	return srv.ServeTLS(net_.TcpKeepAliveListener(ln.(*net.TCPListener), 3*time.Minute), tlsConfig, certFile, keyFile)
}

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler HandlerConn) error {
	server := NewServer()
	server.Handler = handler
	return server.ListenAndServe(addr)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects HTTPS connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
func ListenAndServeTLS(addr string, handler HandlerConn, tlsConfig *tls.Config, certFile, keyFile string) error {
	server := NewServer()
	server.Handler = handler
	return server.ListenAndServeTLS(addr, tlsConfig, certFile, keyFile)
}

// Close immediately closes all active net.Listeners and any
// connections in state StateNew, StateActive, or StateIdle. For a
// graceful shutdown, use Shutdown.
//
// Close does not attempt to close (and does not even know about)
// any hijacked connections, such as WebSockets.
//
// Close returns any error returned from closing the ServeMux's
// underlying ServeMux(s).
func (srv *Server) Close() error {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.closeDoneChanLocked()

	err := srv.closeListenersLocked()
	for c := range srv.activeConn {
		c.close()
		delete(srv.activeConn, c)
	}
	if srv.Handler == nil {
		DefaultServeMux.Close()
	} else {
		if closer, ok := srv.Handler.(io.Closer); ok {
			closer.Close()
		}
	}
	return err
}

// shutdownPollInterval is how often we poll for quiescence
// during ServeMux.Shutdown. This is lower during tests, to
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
// error returned from closing the ServeMux's underlying ServeMux(s).
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
func (srv *Server) Shutdown(ctx context.Context) error {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	lnerr := srv.closeListenersLocked()
	srv.closeDoneChanLocked()

	for _, f := range srv.onShutdown {
		go f()
	}
	srv.mu.Unlock()

	ticker := time.NewTicker(shutdownPollInterval)
	defer ticker.Stop()
	for {
		if srv.closeIdleConns() {
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
func (srv *Server) RegisterOnShutdown(f func()) {
	srv.mu.Lock()
	srv.onShutdown = append(srv.onShutdown, f)
	srv.mu.Unlock()
}

// HandleError registers an error handler that handles listener errors.
func (srv *Server) HandleError(h ErrorHandler) {
	srv.errHandler = h
}

func (srv *Server) handleErr(err error) bool {
	if srv.errHandler == nil {
		return true
	}
	if !srv.errHandler.Continue(err) {
		return false
	}

	if ne, ok := err.(net.Error); ok {
		return ne.Temporary()
	}

	return false
}

// trackListener adds or removes a net.ServeMux to the set of tracked listeners.
//
// We store a pointer to interface in the map set, in case the
// net.ServeMux is not comparable. This is safe because we only call
// trackListener via Serve and can track+defer untrack the same
// pointer to local variable there. We never need to compare a
// ServeMux from another caller.
//
// It reports whether the server is still up (not Shutdown or Closed).
func (srv *Server) trackListener(ln *net.Listener, add bool) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if add {
		if srv.shuttingDown() {
			return false
		}
		if srv.listeners == nil {
			srv.listeners = make(map[*net.Listener]struct{})
		}
		srv.listeners[ln] = struct{}{}
		return true
	}
	delete(srv.listeners, ln)
	return true
}

func (srv *Server) trackConn(c *conn, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if add {
		if srv.activeConn == nil {
			srv.activeConn = make(map[*conn]struct{})
		}
		srv.activeConn[c] = struct{}{}
	} else {
		delete(srv.activeConn, c)
	}
}

func (srv *Server) shuttingDown() bool {
	return srv.inShutdown.Load()
}

// closeIdleConns closes all idle connections and reports whether the
// server is quiescent.
func (srv *Server) closeIdleConns() bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	quiescent := true
	for c := range srv.activeConn {
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
		delete(srv.activeConn, c)
	}
	return quiescent
}

func (srv *Server) closeListenersLocked() error {
	var err error
	for ln := range srv.listeners {
		if cerr := (*ln).Close(); cerr != nil && err == nil {
			err = cerr
		}
		delete(srv.listeners, ln)
	}
	return err
}
