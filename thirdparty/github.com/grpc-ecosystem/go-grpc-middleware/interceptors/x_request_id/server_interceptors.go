// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x_request_id

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors with tags in context.
// key is RequestID within Context if have
// chained to chain multiple request ids by generating new request id for each request and concatenating it to original request ids.
func UnaryServerInterceptor(key interface{}, chained bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		if chained {
			newCtx = newContextForHandleRequestIDChain(ctx, key)
		} else {
			newCtx = newContextForHandleRequestID(ctx, key)
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor with tags in context.
// key is RequestID within Context if have
// chained to chain multiple request ids by generating new request id for each request and concatenating it to original request ids.
func StreamServerInterceptor(key interface{}, chained bool) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		if chained {
			newCtx = newContextForHandleRequestIDChain(stream.Context(), key)
		} else {
			newCtx = newContextForHandleRequestID(stream.Context(), key)
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
