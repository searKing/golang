// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// A Gateway defines parameters for running an HTTP+gRPC-Gateway+gRPC server.
// The zero value for Gateway is a valid configuration.
//
//go:generate go-option -type=Gateway
type Gateway struct {
	// options
	opt         gatewayOption `option:"-"`
	http.Server `option:"-"`  // Gateway is a HTTP server, actually.

	httpMuxToGrpc *runtime.ServeMux `option:"-"` // gRPC to JSON proxy generator following the gRPC HTTP spec
	Handler       http.Handler      `option:"-"` // Not Found HTTP Handler

	// runtime
	grpcServer *grpc.Server `option:"-"` // a gRPC server to serve RPC requests.
	listenAddr *net.TCPAddr `option:"-"` // addr actually listen, useful for :0 or port not specified.
}

func NewGateway(addr string, opts ...GatewayOption) *Gateway {
	return NewGatewayTLS(addr, nil, opts...)
}

// NewGatewayTLS TLSConfig optionally provides a TLS configuration for use
// by ServeTLS and ListenAndServeTLS. Note that this value is
// cloned by ServeTLS and ListenAndServeTLS, so it's not
// possible to modify the configuration with methods like
// tls.Config.SetSessionTicketKeys. To use
// SetSessionTicketKeys, use Server.Serve with a TLS Listener
// instead.
func NewGatewayTLS(addr string, tlsConfig *tls.Config, opts ...GatewayOption) *Gateway {
	server := &Gateway{
		Server: http.Server{
			Addr:      addr,
			TLSConfig: tlsConfig,
		},
	}
	return server.ApplyOptions(opts...)
}

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler http.Handler, opts ...GatewayOption) error {
	server := &Gateway{
		Server: http.Server{
			Addr: addr,
		},
		Handler: handler,
	}
	return server.ApplyOptions(opts...).ListenAndServe()
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects HTTPS connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler, opts ...GatewayOption) error {
	server := &Gateway{
		Server: http.Server{
			Addr: addr,
		},
		Handler: handler,
	}
	return server.ApplyOptions(opts...).ListenAndServeTLS(certFile, keyFile)
}

// Serve accepts incoming HTTP connections on the listener l,
// creating a new service goroutine for each. The service goroutines
// read requests and then call handler to reply to them.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// HTTP/2 support is only enabled if the Listener returns *tls.Conn
// connections and they were configured with "h2" in the TLS
// Config.NextProtos.
//
// Serve always returns a non-nil error.
func Serve(l net.Listener, handler http.Handler, opts ...GatewayOption) error {
	srv := &Gateway{Handler: handler}
	return srv.ApplyOptions(opts...).Serve(l)
}

// ServeTLS accepts incoming HTTPS connections on the listener l,
// creating a new service goroutine for each. The service goroutines
// read requests and then call handler to reply to them.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// Additionally, files containing a certificate and matching private key
// for the server must be provided. If the certificate is signed by a
// certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
//
// ServeTLS always returns a non-nil error.
func ServeTLS(l net.Listener, handler http.Handler, certFile, keyFile string, opts ...GatewayOption) error {
	srv := &Gateway{Handler: handler}
	return srv.ApplyOptions(opts...).ServeTLS(l, certFile, keyFile)
}

func (gw *Gateway) Serve(l net.Listener) error {
	gw.preServe()
	return gw.Server.Serve(l)
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
func (gw *Gateway) ServeTLS(l net.Listener, certFile, keyFile string) error {
	gw.preServe()
	return gw.Server.ServeTLS(l, certFile, keyFile)
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// If srv.Addr is blank, ":http" is used.
//
// ListenAndServe always returns a non-nil error. After Shutdown or Close,
// the returned error is ErrServerClosed.
func (gw *Gateway) ListenAndServe() error {
	gw.preServe()
	return gw.Server.ListenAndServe()
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
func (gw *Gateway) ListenAndServeTLS(certFile, keyFile string) error {
	gw.preServe()
	return gw.Server.ListenAndServeTLS(certFile, keyFile)
}

// BindAddr returns actual addr bind after a net.Listener on
func (gw *Gateway) BindAddr() string {
	if gw.listenAddr != nil {
		return gw.listenAddr.String()
	}
	return gw.Server.Addr
}

// RegisterGRPCHandler registers grpc handler of the gateway
func (gw *Gateway) RegisterGRPCHandler(handler GRPCHandler) {
	handler.Register(gw.grpcServer)
}

// RegisterHTTPHandler registers http handler of the gateway
func (gw *Gateway) RegisterHTTPHandler(ctx context.Context, handler HTTPHandler) error {
	//scheme://authority/endpoint
	return handler.Register(ctx, gw.httpMuxToGrpc, "passthrough:///"+gw.Server.Addr, gw.opt.ClientDialOpts())
}

// RegisterGRPCFunc registers grpc handler of the gateway
func (gw *Gateway) RegisterGRPCFunc(handler func(srv *grpc.Server)) {
	gw.RegisterGRPCHandler(GRPCHandlerFunc(handler))
}

// RegisterHTTPFunc registers http handler of the gateway
func (gw *Gateway) RegisterHTTPFunc(ctx context.Context, handler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error) error {
	return gw.RegisterHTTPHandler(ctx, HTTPHandlerFunc(handler))
}

func (gw *Gateway) preServe(opts ...GatewayOption) {
	gw.ApplyOptions(opts...)

	gw.resolveListenAddrLater()

	gw.installTlsConfig()

	// gRPC to JSON proxy generator following the gRPC HTTP spec
	gw.installNotFoundHandler() // --> srvMuxOpts
	gw.httpMuxToGrpc = runtime.NewServeMux(gw.opt.srvMuxOpts...)

	// a gRPC server to serve RPC requests.
	gw.grpcServer = grpc.NewServer(gw.opt.ServerOptions()...)
	// register the server reflection service on the given gRPC server.
	if gw.opt.grpcServerOpts.withReflectionService {
		gw.registerGrpcReflection()
	}

	// delegate to grpcServer on incoming gRPC connections or HTTP connections otherwise.
	gw.Server.Handler = grpc_.GrpcOrDefaultHandler(gw.grpcServer, &serverHandler{gateway: gw})
}

// registerGrpcReflection registers the server reflection service on the given gRPC server.
// can be called once, recommend being called before Serve, ServeTLS, ListenAndServe or ListenAndServeTLS and so on.
func (gw *Gateway) registerGrpcReflection() {
	// grpcurl -plaintext localhost:1234 list
	// -plaintext: avoid Failed to dial target host "localhost:1234": tls: first record does not look like a TLS handshake
	// avoid: Failed to list services: server does not support the reflection API
	// Register reflection service on gRPC server.
	reflection.Register(gw.grpcServer)
}

func (gw *Gateway) installTlsConfig() {
	tlsConfig := gw.TLSConfig
	if tlsConfig != nil {
		creds := credentials.NewTLS(tlsConfig)
		// for grpc server
		gw.opt.grpcServerOpts.opts = append([]grpc.ServerOption{grpc.Creds(creds)}, gw.opt.grpcServerOpts.opts...)
		// for grpc client to server
		gw.opt.grpcClientDialOpts = append([]grpc.DialOption{grpc.WithTransportCredentials(creds)}, gw.opt.grpcClientDialOpts...)
	} else {
		// disables transport security
		gw.opt.grpcClientDialOpts = append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, gw.opt.grpcClientDialOpts...)
	}
}

func (gw *Gateway) installNotFoundHandler() {
	// Not Found Handler
	gw.opt.srvMuxOpts = append(gw.opt.srvMuxOpts,
		runtime.WithRoutingErrorHandler(
			func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
				if httpStatus == http.StatusNotFound || httpStatus == http.StatusMethodNotAllowed {
					notFoundHandler := gw.Handler
					if notFoundHandler == nil {
						notFoundHandler = http.DefaultServeMux
					}
					notFoundHandler.ServeHTTP(w, r)
					return
				}
				runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, httpStatus)
			}))
}

func (gw *Gateway) resolveListenAddrLater() {
	ctx := gw.Server.BaseContext
	gw.Server.BaseContext = func(lis net.Listener) context.Context {
		slog.Info(fmt.Sprintf("Serve() passed a net.Listener on %s", lis.Addr().String()))
		if addr, ok := lis.Addr().(*net.TCPAddr); !ok {
			slog.Warn(fmt.Sprintf("GatewayServer expects listener to return a net.TCPAddr. Got %T", lis.Addr()))
		} else {
			gw.listenAddr = addr
		}
		if ctx == nil {
			return context.Background()
		}
		return ctx(lis)
	}
}
