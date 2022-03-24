// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timeratelimit

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryClientInterceptor returns a new unary client interceptor that performs request rate limiting.
func UnaryClientInterceptor(r rate.Limit, b int) grpc.UnaryClientInterceptor {
	limiter := rate.NewLimiter(r, b)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if limiter.Allow() {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		return status.Errorf(codes.ResourceExhausted,
			"%s is rejected by timeratelimit unary client middleware, please retry later", method)
	}
}

// StreamClientInterceptor returns a new streaming client interceptor that performs rate limiting on the request.
func StreamClientInterceptor(r rate.Limit, b int) grpc.StreamClientInterceptor {
	limiter := rate.NewLimiter(r, b)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if limiter.Allow() {
			return streamer(ctx, desc, cc, method, opts...)
		}
		return nil, status.Errorf(codes.ResourceExhausted,
			"%s is rejected by timeratelimit stream client middleware, please retry later", method)
	}
}
