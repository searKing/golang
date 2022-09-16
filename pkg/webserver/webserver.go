// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	runtime_ "github.com/searKing/golang/go/runtime"
	"github.com/searKing/golang/pkg/webserver/healthz"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	"github.com/sirupsen/logrus"
)

type WebHandler interface {
	SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway)
}

type WebServer struct {
	Name string
	// BindAddress is the host name to use for bind (local internet) facing URLs (e.g. Loopback)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	BindAddress string
	// ExternalAddress is the host name to use for external (public internet) facing URLs (e.g. Swagger)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	ExternalAddress string

	ginBackend  *gin.Engine
	grpcBackend *grpc.Gateway

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
	healthzChecks          []healthz.HealthChecker
	healthzChecksInstalled bool
	// livez checks
	livezLock            sync.Mutex
	livezChecks          []healthz.HealthChecker
	livezChecksInstalled bool
	// readyz checks
	readyzLock            sync.Mutex
	readyzChecks          []healthz.HealthChecker
	readyzChecksInstalled bool
	livezGracePeriod      time.Duration

	// the readiness stop channel is used to signal that the apiserver has initiated a shutdown sequence, this
	// will cause readyz to return unhealthy.
	readinessStopCh chan struct{}

	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration

	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration
}

func NewWebServer(fc FactoryConfig, configs ...FactoryConfigFunc) (*WebServer, error) {
	f, err := NewFactory(fc, configs...)
	if err != nil {
		return nil, err
	}
	return f.New()
}

// preparedWebServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedWebServer struct {
	*WebServer
}

// PrepareRun does post API installation setup steps. It calls recursively the same function of the delegates.
func (s *WebServer) PrepareRun() (preparedWebServer, error) {
	if s.grpcBackend != nil {
		s.grpcBackend.Handler = s.ginBackend
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
	logrus.Infof("Serving securely on %s", s.grpcBackend.Addr)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stoppedHttpServerCtx, stopHttpServer := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer stopHttpServer()
		defer logrus.Info("[graceful-termination] shutdown executed")
		<-ctx.Done()

		// As soon as shutdown is initiated, /readyz should start returning failure.
		// This gives the load balancer a window defined by ShutdownDelayDuration to detect that /readyz is red
		// and stop sending traffic to this server.
		close(s.readinessStopCh)
		logrus.Infof("[graceful-termination] shutdown is initiated and delayed after %d", s.ShutdownDelayDuration)
		time.Sleep(s.ShutdownDelayDuration)
	}()

	// close socket after delayed stopCh
	stopHttpServerCtx, stoppedHttpServerCtx, err := s.NonBlockingRun(stoppedHttpServerCtx)
	if err != nil {
		cancel()
		return err
	}

	logrus.Info("[graceful-termination] waiting for shutdown to be initiated")
	// wait for stoppedCh that is closed when the graceful termination (server.Shutdown) is finished.
	<-stopHttpServerCtx.Done()
	// run shutdown hooks directly. This includes deregistering from the kubernetes endpoint in case of web server.
	func() {
		defer cancel()
		defer func() {
			logrus.WithField("name", s.Name).Infof("[graceful-termination] pre-shutdown hooks completed")
		}()
		err = s.RunPreShutdownHooks()
	}()
	if err != nil {
		return err
	}

	// wait for the delayed stopCh before closing the handler chain (it rejects everything after Wait has been called).
	logrus.Info("[graceful-termination] waiting for http server to be stopped")
	<-stoppedHttpServerCtx.Done()
	logrus.Info("[graceful-termination] waiting for http server to be shutdown executed")

	wg.Wait()
	logrus.Info("[graceful-termination] webserver is exiting")
	return nil
}

// NonBlockingRun spawns the secure http|grpc server. An error is
// returned if the secure port cannot be listened on.
// The returned context is done when the (asynchronous) termination is finished.
func (s preparedWebServer) NonBlockingRun(ctx context.Context) (stopCtx, stoppedCtx context.Context, err error) {
	// Shutdown server gracefully.
	stopCtx, stop := context.WithCancel(ctx)
	stoppedCtx, stopped := context.WithCancel(context.Background())
	// Start the shutdown daemon before any request comes in.
	go func() {
		defer stopped()
		select {
		case <-stopCtx.Done():
		}
		// Now that listener have bound successfully, it is the
		// responsibility of the caller to close the provided channel to
		// ensure cleanup.
		if s.ShutdownTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), s.ShutdownTimeout)
			defer cancel()
		}
		err := s.grpcBackend.Shutdown(ctx)
		msg := fmt.Sprintf("Have shutdown http server on %s", s.grpcBackend.Addr)
		if err != nil {
			logrus.WithError(err).Errorf(msg)
		}
		logrus.Infof(msg)
	}()

	go func() {
		defer runtime_.LogPanic.Recover()
		defer stop()
		var err error
		err = s.grpcBackend.ListenAndServe()
		msg := fmt.Sprintf("Stopped listening on %s", s.grpcBackend.Addr)
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			logrus.Info(msg)
			return
		}
		select {
		case <-stoppedCtx.Done():
			logrus.Info(msg)
		default: // not caused by Shutdown
			logrus.WithError(err).Errorf(msg)
		}
		return
	}()

	// Now that listener have bound successfully, it is the
	// responsibility of the caller to close the provided channel to
	// ensure cleanup.
	go func() {
		defer runtime_.LogPanic.Recover()

		var err error
		defer func() {
			if err != nil {
				stop()
			}
		}()
		err = s.RunPostStartHooks(stopCtx)
		msg := fmt.Sprintf("RunPostStartHooks on %s", s.grpcBackend.Addr)
		if err == nil {
			logrus.Info(msg)
			return
		}
		if errors.Is(err, context.Canceled) {
			logrus.WithError(err).Errorf(msg)
			return
		}
		logrus.WithError(err).Errorf("%s due to error: %v", msg, err)
	}()

	return stopCtx, stoppedCtx, nil
}

func (s *WebServer) InstallWebHandlers(handlers ...WebHandler) {
	for _, h := range handlers {
		if h == nil {
			continue
		}
		h.SetRoutes(s.ginBackend, s.grpcBackend)
	}
}
