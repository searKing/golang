// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	http_ "github.com/searKing/golang/go/net/http"
	runtime_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type gatewayOption struct {
	// for grpc server
	grpcServerOpts struct {
		opts                  []grpc.ServerOption
		unaryInterceptors     []grpc.UnaryServerInterceptor
		streamInterceptors    []grpc.StreamServerInterceptor
		withReflectionService bool // registers the server reflection service on the given gRPC server.
	}

	// for http client to redirect to grpc server
	grpcClientDialOpts []grpc.DialOption

	srvMuxOpts []runtime.ServeMuxOption

	interceptors http_.HandlerInterceptorChain
}

func (opt *gatewayOption) ServerOptions() []grpc.ServerOption {
	var streamInterceptors []grpc.StreamServerInterceptor
	streamInterceptors = append(streamInterceptors, grpc_ctxtags.StreamServerInterceptor(),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor())
	streamInterceptors = append(streamInterceptors, opt.grpcServerOpts.streamInterceptors...)

	var unaryInterceptors []grpc.UnaryServerInterceptor
	unaryInterceptors = append(unaryInterceptors, grpc_ctxtags.UnaryServerInterceptor(),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor())
	unaryInterceptors = append(unaryInterceptors, opt.grpcServerOpts.unaryInterceptors...)

	return append(opt.grpcServerOpts.opts,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)))
}

func (opt *gatewayOption) ClientDialOpts() []grpc.DialOption {
	var clientInterceptors []grpc.StreamClientInterceptor
	clientInterceptors = append(clientInterceptors,
		grpc_opentracing.StreamClientInterceptor(),
		grpc_prometheus.StreamClientInterceptor)

	var unaryInterceptors []grpc.UnaryClientInterceptor
	unaryInterceptors = append(unaryInterceptors,
		grpc_opentracing.UnaryClientInterceptor(),
		grpc_prometheus.UnaryClientInterceptor)

	return append(opt.grpcClientDialOpts,
		grpc.WithChainStreamInterceptor(grpc_middleware.ChainStreamClient(clientInterceptors...)),
		grpc.WithChainUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)))
}

func WithGrpcUnaryServerChain(interceptors ...grpc.UnaryServerInterceptor) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.grpcServerOpts.unaryInterceptors = append(gateway.opt.grpcServerOpts.unaryInterceptors, interceptors...)
	})
}

func WithGrpcStreamServerChain(interceptors ...grpc.StreamServerInterceptor) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.grpcServerOpts.streamInterceptors = append(gateway.opt.grpcServerOpts.streamInterceptors, interceptors...)
	})
}

func WithGrpcUnaryClientChain(interceptors ...grpc.UnaryClientInterceptor) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.grpcClientDialOpts = append(gateway.opt.grpcClientDialOpts, grpc.WithChainUnaryInterceptor(interceptors...))
	})
}

func WithGrpcStreamClientChain(interceptors ...grpc.StreamClientInterceptor) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.grpcClientDialOpts = append(gateway.opt.grpcClientDialOpts, grpc.WithChainStreamInterceptor(interceptors...))
	})
}

func WithGrpcReflectionService(autoRegistered bool) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.grpcServerOpts.withReflectionService = autoRegistered
	})
}

func WithGrpcServerOption(opts ...grpc.ServerOption) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.grpcServerOpts.opts = append(gateway.opt.grpcServerOpts.opts, opts...)
	})
}

func WithGrpcServeMuxOption(opts ...runtime.ServeMuxOption) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.srvMuxOpts = append(gateway.opt.srvMuxOpts, opts...)
	})
}

func WithGrpcDialOption(opts ...grpc.DialOption) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.grpcClientDialOpts = append(gateway.opt.grpcClientDialOpts, opts...)
	})
}

// helper below

func WithLogrusLogger(logger *logrus.Logger) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		l := logrus.NewEntry(logger)
		WithGrpcStreamServerChain(grpc_logrus.StreamServerInterceptor(l)).apply(gateway)
		WithGrpcUnaryServerChain(grpc_logrus.UnaryServerInterceptor(l)).apply(gateway)
	})
}

func WithStreamErrorHandler(fn runtime.StreamErrorHandlerFunc) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		WithGrpcServeMuxOption(runtime.WithStreamErrorHandler(fn)).apply(gateway)
	})
}

// WithHTTPErrorHandler replies to the request with the error.
// You can set a custom function to this variable to customize error format.
func WithHTTPErrorHandler(fn HTTPErrorHandler) GatewayOption {
	return WithGrpcServeMuxOption(runtime.WithErrorHandler(fn.HandleHTTPError))
	//return GatewayOptionFunc(func(gateway *Gateway) {
	//	runtime.HTTPError = fn.HandleHTTPError
	//})
}

func WithMarshalerOption(mime string, marshaler runtime.Marshaler) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		WithGrpcServeMuxOption(runtime.WithMarshalerOption(mime, marshaler)).apply(gateway)
	})
}

func WithDefaultMarsherOption() []GatewayOption {
	return []GatewayOption{
		WithMarshalerOption(runtime.MIMEWildcard, runtime_.NewHTTPBodyJsonMarshaler()),
		WithMarshalerOption(binding.MIMEJSON, runtime_.NewHTTPBodyJsonMarshaler()),
		WithMarshalerOption(binding.MIMEPROTOBUF, runtime_.NewHTTPBodyProtoMarshaler()),
		WithMarshalerOption(binding.MIMEYAML, runtime_.NewHTTPBodyYamlMarshaler()),
		WithGrpcReflectionService(true),
	}

}

func WithDefault() []GatewayOption {
	var opts []GatewayOption
	opts = append(opts, WithDefaultMarsherOption()...)
	opts = append(opts, WithGrpcReflectionService(true))
	return opts
}

//func WithForwardResponseMessageHandler(fn ForwardResponseMessageHandler) GatewayOption {
//	return GatewayOptionFunc(func(gateway *Gateway) {
//		runtime.WithForwardResponseOption()
//		runtime.ForwardResponseMessage = nil
//	})
//}

func WithHttpHandlerInterceptor(opts ...http_.HandlerInterceptorChainOption) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.interceptors.ApplyOptions(opts...)
	})
}

func WithHttpPreHandler(preHandle func(w http.ResponseWriter, r *http.Request) error) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(preHandle, nil, nil, nil))
}

// WithHttpWrapper is a decorator or middleware of http.Handler
func WithHttpWrapper(wrapper func(http.Handler) http.Handler) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(nil, wrapper, nil, nil))
}

func WithHttpPostHandler(
	postHandle func(w http.ResponseWriter, r *http.Request)) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(nil, nil, postHandle, nil))
}

func WithHttpAfterCompletion(
	afterCompletion func(w http.ResponseWriter, r *http.Request, err interface{})) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(nil, nil, nil, afterCompletion))
}

// Deprecated: Use WithHttpPreHandler instead.
func WithHttpRewriter(rewriter func(w http.ResponseWriter, r *http.Request) error) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(rewriter, nil, nil, nil))
}
