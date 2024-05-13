// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tags

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a new unary client interceptor with tags in context.
func UnaryClientInterceptor(key any, values map[string]any) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := newContextTagsForCall(ctx, KindClient, method, key, values)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

// StreamClientInterceptor returns a new streaming client interceptor with tags in context.
func StreamClientInterceptor(key any, values map[string]any) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newCtx := newContextTagsForCall(ctx, KindClient, method, key, values)
		return streamer(newCtx, desc, cc, method, opts...)
	}
}
