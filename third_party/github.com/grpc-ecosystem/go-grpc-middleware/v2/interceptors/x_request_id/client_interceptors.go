// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x_request_id

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a new unary client interceptor with tags in context.
func UnaryClientInterceptor(key interface{}, chained bool) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var newCtx context.Context
		if chained {
			newCtx = newContextForHandleRequestIDChain(ctx, key)
		} else {
			newCtx = newContextForHandleRequestID(ctx, key)
		}
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

// StreamServerInterceptor returns a new streaming client interceptor with tags in context.
func StreamClientInterceptor(key interface{}, chained bool) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var newCtx context.Context
		if chained {
			newCtx = newContextForHandleRequestIDChain(ctx, key)
		} else {
			newCtx = newContextForHandleRequestID(ctx, key)
		}
		return streamer(newCtx, desc, cc, method, opts...)
	}
}
