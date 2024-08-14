// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	http_ "github.com/searKing/golang/go/net/http"
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

func WithDefault() []GatewayOption {
	var opts []GatewayOption
	opts = append(opts, WithDefaultMarshalerOption()...)
	opts = append(opts, WithGrpcReflectionService(true))
	return opts
}
