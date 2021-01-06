// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	http_ "github.com/searKing/golang/pkg/net/http"
)

//go:generate go-option -type=Gateway
type Gateway struct {
	// options
	opt gatewayOption
	http.Server

	httpMuxToGrpc *runtime.ServeMux
	Handler       http.Handler

	// runtime
	grpcServer *grpc.Server

	once sync.Once
}

func NewGateway(addr string, opts ...GatewayOption) *Gateway {
	return NewGatewayTLS(addr, nil, opts...)
}

// TLSConfig optionally provides a TLS configuration for use
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

func (gateway *Gateway) lazyInit(opts ...GatewayOption) {
	gateway.once.Do(func() {
		gateway.ApplyOptions(opts...)

		tlsConfig := gateway.TLSConfig
		if tlsConfig != nil {
			// for grpc server
			gateway.opt.grpcServerOpts.opts = append(gateway.opt.grpcServerOpts.opts,
				grpc.Creds(credentials.NewTLS(tlsConfig)))
			// for grpc client to server
			gateway.opt.grpcClientDialOpts = append(gateway.opt.grpcClientDialOpts,
				grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		} else {
			if len(gateway.opt.grpcClientDialOpts) == 0 {
				// disables transport security
				gateway.opt.grpcClientDialOpts = append(gateway.opt.grpcClientDialOpts, grpc.WithInsecure())
			}
		}

		gateway.opt.srvMuxOpts = append(gateway.opt.srvMuxOpts,
			runtime.WithRoutingErrorHandler(
				func(ctx context.Context, mux *runtime.ServeMux,
					marshaler runtime.Marshaler,
					w http.ResponseWriter, r *http.Request, code int) {
					httpHandler := gateway.Handler
					if httpHandler == nil {
						httpHandler = http.DefaultServeMux
					}
					if code == http.StatusNotFound || code == http.StatusMethodNotAllowed {
						httpHandler.ServeHTTP(w, r)
						return
					}
					runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, code)
				}))

		gateway.grpcServer = grpc.NewServer(gateway.opt.ServerOptions()...)
		gateway.httpMuxToGrpc = runtime.NewServeMux(gateway.opt.srvMuxOpts...)
		gateway.Server.Handler = http_.GrpcOrDefaultHandler(gateway.grpcServer, &serverHandler{
			gateway: gateway,
		})
	})

}
func (gateway *Gateway) Serve(l net.Listener) error {
	gateway.lazyInit()
	gateway.registerGrpcReflection()
	return gateway.Server.Serve(l)
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
func (gateway *Gateway) ServeTLS(l net.Listener, certFile, keyFile string) error {
	gateway.lazyInit()
	gateway.registerGrpcReflection()
	return gateway.Server.ServeTLS(l, certFile, keyFile)
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// If srv.Addr is blank, ":http" is used.
//
// ListenAndServe always returns a non-nil error. After Shutdown or Close,
// the returned error is ErrServerClosed.
func (gateway *Gateway) ListenAndServe() error {
	gateway.lazyInit()
	gateway.registerGrpcReflection()
	return gateway.Server.ListenAndServe()
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
func (gateway *Gateway) ListenAndServeTLS(certFile, keyFile string) error {
	gateway.lazyInit()
	gateway.registerGrpcReflection()
	return gateway.Server.ListenAndServeTLS(certFile, keyFile)
}

// RegisterGRPCHandler registers grpc handler of the gateway
func (gateway *Gateway) RegisterGRPCHandler(handler GRPCHandler) {
	gateway.lazyInit()
	handler.Register(gateway.grpcServer)
}

// RegisterHTTPHandler registers http handler of the gateway
func (gateway *Gateway) RegisterHTTPHandler(ctx context.Context, handler HTTPHandler) error {
	gateway.lazyInit()
	//scheme://authority/endpoint
	return handler.Register(ctx, gateway.httpMuxToGrpc, "passthrough:///"+gateway.Server.Addr, gateway.opt.grpcClientDialOpts)
}

// RegisterGRPCFunc registers grpc handler of the gateway
func (gateway *Gateway) RegisterGRPCFunc(handler func(srv *grpc.Server)) {
	gateway.lazyInit()
	gateway.RegisterGRPCHandler(GRPCHandlerFunc(handler))
}

// RegisterHTTPFunc registers http handler of the gateway
func (gateway *Gateway) RegisterHTTPFunc(ctx context.Context, handler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error) error {
	gateway.lazyInit()
	return gateway.RegisterHTTPHandler(ctx, HTTPHandlerFunc(handler))
}

// registerGrpcReflection registers the server reflection service on the given gRPC server.
// can be called once, recommently to be called before Serve, ServeTLS, ListenAndServe or ListenAndServeTLS and so on.
func (gateway *Gateway) registerGrpcReflection() {
	if gateway.grpcServer == nil || !gateway.opt.grpcServerOpts.withReflectionService {
		return
	}
	// grpcurl -plaintext localhost:1234 list
	// -plaintext: avoid Failed to dial target host "localhost:1234": tls: first record does not look like a TLS handshake
	// avoid: Failed to list services: server does not support the reflection API
	// Register reflection service on gRPC server.
	reflection.Register(gateway.grpcServer)
}
