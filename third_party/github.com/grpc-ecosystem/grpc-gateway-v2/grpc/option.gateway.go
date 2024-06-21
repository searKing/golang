// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	http_ "github.com/searKing/golang/go/net/http"
	runtime_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/runtime"
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
	loggingOpts  []logging.Option
}

func (opt *gatewayOption) ServerOptions() []grpc.ServerOption {
	var streamInterceptors []grpc.StreamServerInterceptor
	streamInterceptors = append(streamInterceptors, grpcrecovery.StreamServerInterceptor())
	streamInterceptors = append(streamInterceptors, opt.grpcServerOpts.streamInterceptors...)

	var unaryInterceptors []grpc.UnaryServerInterceptor
	unaryInterceptors = append(unaryInterceptors, grpcrecovery.UnaryServerInterceptor())
	unaryInterceptors = append(unaryInterceptors, opt.grpcServerOpts.unaryInterceptors...)
	return append(opt.grpcServerOpts.opts, grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...))
}

func (opt *gatewayOption) ClientDialOpts() []grpc.DialOption {
	var streamInterceptors []grpc.StreamClientInterceptor
	streamInterceptors = append(streamInterceptors)

	var unaryInterceptors []grpc.UnaryClientInterceptor
	unaryInterceptors = append(unaryInterceptors)
	return append(opt.grpcClientDialOpts, grpc.WithChainStreamInterceptor(streamInterceptors...),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...))
}

func WithLoggingOption(opts ...logging.Option) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.loggingOpts = append(gateway.opt.loggingOpts, opts...)
	})
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

func WithStreamErrorHandler(fn runtime.StreamErrorHandlerFunc) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		WithGrpcServeMuxOption(runtime.WithStreamErrorHandler(fn)).apply(gateway)
	})
}

// WithHTTPErrorHandler replies to the request with the error.
// You can set a custom function to this variable to customize error format.
func WithHTTPErrorHandler(fn HTTPErrorHandler) GatewayOption {
	return WithGrpcServeMuxOption(runtime.WithErrorHandler(fn.HandleHTTPError))
}

func WithMarshalerOption(mime string, marshaler runtime.Marshaler) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		WithGrpcServeMuxOption(runtime.WithMarshalerOption(mime, marshaler)).apply(gateway)
	})
}

func WithDefaultMarshalerOption() []GatewayOption {
	return []GatewayOption{
		WithMarshalerOption(runtime.MIMEWildcard, runtime_.NewHTTPBodyJsonMarshaler()),
		WithMarshalerOption(binding.MIMEJSON, runtime_.NewHTTPBodyJsonMarshaler()),
		WithMarshalerOption(binding.MIMEPROTOBUF, runtime_.NewHTTPBodyProtoMarshaler()),
		WithMarshalerOption(binding.MIMEYAML, runtime_.NewHTTPBodyYamlMarshaler()),
	}
}

// Deprecated: Use WithDefaultMarshalerOption instead.
func WithDefaultMarsherOption() []GatewayOption {
	return WithDefaultMarshalerOption()
}

func WithDefault() []GatewayOption {
	var opts []GatewayOption
	opts = append(opts, WithDefaultMarshalerOption()...)
	opts = append(opts, WithGrpcReflectionService(true))
	return opts
}

//func WithForwardResponseMessageHandler(fn ForwardResponseMessageHandler) GatewayOption {
//	return GatewayOptionFunc(func(gateway *Gateway) {
//		runtime.WithForwardResponseOption()
//		runtime.ForwardResponseMessage = nil
//	})
//}

func WithForwardResponseMessageHandler(fn ForwardResponseOptionHandler) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		WithGrpcServeMuxOption(runtime.WithForwardResponseOption(fn.ForwardResponseOption))
	})
}

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
	afterCompletion func(w http.ResponseWriter, r *http.Request, err any)) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(nil, nil, nil, afterCompletion))
}

// Deprecated: Use WithHttpPreHandler instead.
func WithHttpRewriter(rewriter func(w http.ResponseWriter, r *http.Request) error) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(rewriter, nil, nil, nil))
}

// ExtractLoggingOptions extract all [logging.Option] from the given options.
func ExtractLoggingOptions(options ...GatewayOption) []logging.Option {
	var g Gateway
	(&g).ApplyOptions(options...)
	return g.opt.loggingOpts
}
