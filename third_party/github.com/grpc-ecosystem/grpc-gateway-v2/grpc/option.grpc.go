// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

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
