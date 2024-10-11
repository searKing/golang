// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	http_ "github.com/searKing/golang/go/net/http"
	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"
)

type gatewayOption struct {
	// for grpc server
	grpcServerOpts struct {
		opts                  []grpc.ServerOption
		unaryInterceptors     []grpc.UnaryServerInterceptor
		streamInterceptors    []grpc.StreamServerInterceptor
		withReflectionService bool // registers the server reflection service on the given gRPC server.
	}

	// for http handler to call grpc service function directly.
	// as gRPC-Gateway does not support gRPC interceptors when call gRPC's service handler in process.
	// See: https://github.com/grpc-ecosystem/grpc-gateway/issues/1043
	httpNoForwardInterceptors grpc_.UnaryHandlerDecorators

	// for http client to redirect to grpc server
	grpcClientDialOpts []grpc.DialOption

	// for gateway
	srvMuxOpts       []runtime.ServeMuxOption
	httpInterceptors http_.HandlerInterceptorChain

	loggingOpts []logging.Option
}

func (opt *gatewayOption) ServerOptions() []grpc.ServerOption {
	var streamInterceptors []grpc.StreamServerInterceptor
	streamInterceptors = append(streamInterceptors, grpcrecovery.StreamServerInterceptor())
	streamInterceptors = append(streamInterceptors, opt.grpcServerOpts.streamInterceptors...)

	var unaryInterceptors []grpc.UnaryServerInterceptor
	unaryInterceptors = append(unaryInterceptors, grpcrecovery.UnaryServerInterceptor())
	unaryInterceptors = append(unaryInterceptors, opt.grpcServerOpts.unaryInterceptors...)
	return append(opt.grpcServerOpts.opts,
		grpc.ChainStreamInterceptor(streamInterceptors...),
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

func WithDefault() []GatewayOption {
	var opts []GatewayOption
	opts = append(opts, WithDefaultMarshalerOption()...)
	opts = append(opts, WithGrpcReflectionService(true))
	return opts
}
