// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	runtime_ "github.com/searKing/golang/go/runtime"
	"github.com/searKing/golang/pkg/webserver/healthz"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	"github.com/sirupsen/logrus"
)

const (
	defaultKeepAlivePeriod = 3 * time.Minute
)

type WebHandler interface {
	SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway)
}

type WebServer struct {
	// Name is name of web server, optional
	Name string
	// BindAddress is the host name to use for bind (local internet) facing URLs (e.g. Loopback)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	BindAddress string
	// ExternalAddress is the host name to use for external (public internet) facing URLs (e.g. Swagger)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	ExternalAddress string

	GinBackend  *gin.Engine
	GrpcBackend *grpc.Gateway

	// PostStartHooks are each called after the server has started listening, in a separate go func for each
	// with no guarantee of ordering between them.  The map key is a name used for error reporting.
	// It may kill the process with a panic if it wishes to by returning an error.
	postStartHookLock    sync.Mutex
	postStartHooks       map[string]postStartHookEntry
	postStartHooksCalled bool

	preShutdownHookLock    sync.Mutex
	preShutdownHooks       map[string]preShutdownHookEntry
	preShutdownHooksCalled bool

	// healthz checks
	healthzLock            sync.Mutex
	healthzChecks          []healthz.HealthCheck
	healthzChecksInstalled bool
	// livez checks
	livezLock            sync.Mutex
	livezChecks          []healthz.HealthCheck
	livezChecksInstalled bool
	// readyz checks
	readyzLock            sync.Mutex
	readyzChecks          []healthz.HealthCheck
	readyzChecksInstalled bool
	livezGracePeriod      time.Duration

	// the readiness stop channel is used to signal that the apiserver has initiated a shutdown sequence, this
	// will cause readyz to return unhealthy.
	readinessStopCh chan struct{}

	delayedStopOrDrainedCh chan struct{}

	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration

	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration

	// ShutdownSendRetryAfter dictates when to initiate shutdown of the HTTP
	// Server during the graceful termination of the apiserver. If true, we wait
	// for non longrunning requests in flight to be drained and then initiate a
	// shutdown of the HTTP Server. If false, we initiate a shutdown of the HTTP
	// Server as soon as ShutdownDelayDuration has elapsed.
	// If enabled, after ShutdownDelayDuration elapses, any incoming request is
	// rejected with a 429 status code and a 'Retry-After' response.
	ShutdownSendRetryAfter bool
}

// preparedWebServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedWebServer struct {
	*WebServer
}

// PrepareRun does post API installation setup steps. It calls recursively the same function of the delegates.
func (s *WebServer) PrepareRun() (preparedWebServer, error) {
	if s.GrpcBackend != nil {
		s.GrpcBackend.Handler = s.GinBackend
	}

	s.installHealthz()
	s.installLivez()
	err := s.addReadyzShutdownCheck(s.readinessStopCh)
	if err != nil {
		logrus.Errorf("Failed to parseViper readyz shutdown check %s", err)
		return preparedWebServer{}, err
	}
	s.installReadyz()

	// Register audit backend preShutdownHook.
	return preparedWebServer{s}, nil
}

// Run spawns the secure http server. It only returns if stopCh is closed
// or the secure port cannot be listened on initially.
func (s preparedWebServer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		<-ctx.Done()

		// As soon as shutdown is initiated, /readyz should start returning failure.
		// This gives the load balancer a window defined by ShutdownDelayDuration to detect that /readyz is red
		// and stop sending traffic to this server.
		close(s.readinessStopCh)
		logrus.Info("[graceful-termination] /readyz is red")
		time.Sleep(s.ShutdownDelayDuration)
		logrus.Info("[graceful-termination] ready to stop sending traffic to this server")
	}()

	shutdownTimeout := s.ShutdownTimeout
	if s.ShutdownSendRetryAfter {
		// when this mode is enabled, we do the following:
		// - the server will continue to listen until all existing requests in flight
		//   (not including active long runnning requests) have been drained.
		// - once drained, http Server Shutdown is invoked with a timeout of 2s,
		//   net/http waits for 1s for the peer to respond to a GO_AWAY frame, so
		//   we should wait for a minimum of 2s
		shutdownTimeout = 2 * time.Second
		logrus.WithField("ShutdownTimeout", shutdownTimeout).Infof("[graceful-termination] using HTTP Server shutdown timeout")
	}

	// pre-shutdown hooks need to finish before we stop the http server
	preShutdownHooksHasStoppedCh := make(chan struct{})
	stopHttpServerCtx, stopHttpServerCancel := context.WithCancel(ctx)
	go func() {
		defer stopHttpServerCancel()
		<-preShutdownHooksHasStoppedCh
	}()

	// close socket after delayed stopCh
	stoppedCh, listenerStoppedCh, err := s.NonBlockingRun(stopHttpServerCtx, shutdownTimeout)
	if err != nil {
		return err
	}
	go func() {
		<-listenerStoppedCh
		logrus.WithField("name", s.Name).Info("[graceful-termination] shutdown web server")
	}()

	logrus.Info("[graceful-termination] waiting for shutdown to be initiated")
	<-ctx.Done() // we can pre shutdown web server

	// run shutdown hooks directly.
	func() {
		defer close(preShutdownHooksHasStoppedCh)
		err = s.RunPreShutdownHooks()
	}()
	if err != nil {
		logrus.WithError(err).Error("[graceful-termination] RunPreShutdownHooks has completed")
		return err
	}
	logrus.Info("[graceful-termination] RunPreShutdownHooks has completed")

	// wait for the delayed stopCh before closing the handler chain (it rejects everything after Wait has been called).
	wg.Wait()
	// wait for stoppedCh that is closed when the graceful termination (server.Shutdown) is finished.
	<-stoppedCh
	logrus.Info("[graceful-termination] web server is exiting")
	return nil
}

// NonBlockingRun spawns the secure http|grpc server. An error is
// returned if the secure port cannot be listened on.
// The returned context is done when the (asynchronous) termination is finished.
// Serve runs the secure http server. It fails only if certificates cannot be loaded or the initial listen call fails.
// The actual server loop (stoppable by closing stopCh) runs in a go routine, i.e. Serve does not block.
// It returns a stoppedCh that is closed when all non-hijacked active requests have been processed.
// It returns a listenerStoppedCh that is closed when the underlying http Server has stopped listening.
func (s preparedWebServer) NonBlockingRun(ctx context.Context, shutdownTimeout time.Duration) (<-chan struct{}, <-chan struct{}, error) {
	// Start the server backend before any request comes in. This means we must call Backend.Run
	// before http server start serving. Otherwise, the Backend.ProcessEvents call might block.

	logrus.Infof("Serving securely on %s", s.GrpcBackend.Addr)
	addr := s.GrpcBackend.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.WithError(err).WithField("addr", addr).Errorf("failed to create listener")
		return nil, nil, err
	}

	var stoppedCh <-chan struct{}
	var listenerStoppedCh <-chan struct{}
	{
		var err error
		stoppedCh, listenerStoppedCh, err = RungGRPCServer(ctx, s.GrpcBackend, ln, shutdownTimeout)
		if err != nil {
			return nil, nil, err
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.RunPostStartHooks(ctx)

	return stoppedCh, listenerStoppedCh, nil
}

func (s *WebServer) InstallWebHandlers(handlers ...WebHandler) {
	for _, h := range handlers {
		if h == nil {
			continue
		}
		h.SetRoutes(s.GinBackend, s.GrpcBackend)
	}
}

// RungGRPCServer spawns a go-routine continuously serving until the stopCh is
// closed.
// It returns a stoppedCh that is closed when all non-hijacked active requests
// have been processed.
// This function does not block
func RungGRPCServer(ctx context.Context, server *grpc.Gateway, ln net.Listener, shutdownTimeout time.Duration) (
	<-chan struct{}, <-chan struct{}, error) {
	if ln == nil {
		return nil, nil, fmt.Errorf("listener must not be nil")
	}

	// Shutdown server gracefully.
	serverShutdownCh, listenerStoppedCh := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(serverShutdownCh)
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		logrus.Infof("Shutting down http server on %s", server.Addr)
		err := server.Shutdown(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("Have not shut down http server on %s", server.Addr)
			return
		}
		logrus.Infof("Have Shut down http server on %s", server.Addr)
	}()

	go func() {
		defer runtime_.NeverPanicButLog.Recover()
		defer close(listenerStoppedCh)

		var listener net.Listener
		listener = tcpKeepAliveListener{ln}
		if server.TLSConfig != nil {
			listener = tls.NewListener(listener, server.TLSConfig)
		}

		err := server.Serve(listener)

		msg := fmt.Sprintf("Stopped listening on %s", ln.Addr().String())
		select {
		case <-ctx.Done():
			logrus.Infof(msg)
		default:
			panic(fmt.Sprintf("%s due to error: %v", msg, err))
		}
	}()

	return serverShutdownCh, listenerStoppedCh, nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// Copied from Go 1.7.2 net/http/server.go
type tcpKeepAliveListener struct {
	net.Listener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	c, err := ln.Listener.Accept()
	if err != nil {
		return nil, err
	}
	if tc, ok := c.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(defaultKeepAlivePeriod)
	}
	return c, nil
}

// tlsHandshakeErrorWriter writes TLS handshake errors to log with
// trace level - V(5), to avoid flooding of tls handshake errors.
type tlsHandshakeErrorWriter struct {
	out io.Writer
}

const tlsHandshakeErrorPrefix = "http: TLS handshake error"

func (w *tlsHandshakeErrorWriter) Write(p []byte) (int, error) {
	if strings.Contains(string(p), tlsHandshakeErrorPrefix) {
		logrus.Infof(string(p))
		return len(p), nil
	}

	// for non tls handshake error, log it as usual
	return w.out.Write(p)
}
