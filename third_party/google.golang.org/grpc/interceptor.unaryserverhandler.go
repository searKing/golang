// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryServerInterceptorGeneric is a generic version of [grpc.UnaryServerInterceptor]
// as gRPC-Gateway does not support gRPC interceptors when call gRPC's service handler in process.
// See: https://github.com/grpc-ecosystem/grpc-gateway/issues/1043
type UnaryServerInterceptorGeneric[REQ any, RESP any] func(ctx context.Context, req REQ, info *grpc.UnaryServerInfo, handler UnaryHandlerGeneric[REQ, RESP]) (resp RESP, err error)

func (f UnaryServerInterceptorGeneric[REQ, RESP]) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return f(ctx, req.(REQ), info, func(ctx context.Context, req REQ) (resp RESP, err error) {
			resp_, err := handler(ctx, req)
			resp, _ = resp_.(RESP)
			return resp, err
		})
	}
}

func NewUnaryServerInterceptorGeneric[REQ any, RESP any](interceptor grpc.UnaryServerInterceptor) UnaryServerInterceptorGeneric[REQ, RESP] {
	return func(ctx context.Context, req REQ, info *grpc.UnaryServerInfo, handler UnaryHandlerGeneric[REQ, RESP]) (resp RESP, err error) {
		resp_, err := interceptor(ctx, req, info, handler.UnaryHandler())
		resp, _ = resp_.(RESP)
		return resp, err
	}
}
