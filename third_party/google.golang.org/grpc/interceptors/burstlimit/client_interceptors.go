// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package burstlimit

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryClientInterceptor returns a new unary client interceptor that performs request burst limiting.
// b bucket size, take effect if b > 0
// timeout ResourceExhausted if cost more than timeout to get a token, take effect if timeout > 0
func UnaryClientInterceptor(b int, timeout time.Duration) grpc.UnaryClientInterceptor {
	limiter := make(chan struct{}, b)
	for i := 0; i < b; i++ {
		limiter <- struct{}{}
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if b > 0 {
			var limiterCtx = ctx
			var cancel context.CancelFunc
			if timeout > 0 {
				limiterCtx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}
			select {
			case <-limiter:
				defer func() {
					limiter <- struct{}{}
				}()
			case <-limiterCtx.Done():
				return status.Errorf(codes.ResourceExhausted,
					"%s is rejected by burstlimit unary client middleware, please retry later: %s", method, limiterCtx.Err())
			}
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClientInterceptor returns a new streaming client interceptor that performs burst limiting on the request.
// b bucket size, take effect if b > 0
// timeout ResourceExhausted if cost more than timeout to get a token, take effect if timeout > 0
func StreamClientInterceptor(b int, timeout time.Duration) grpc.StreamClientInterceptor {
	limiter := make(chan struct{}, b)
	for i := 0; i < b; i++ {
		limiter <- struct{}{}
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if b > 0 {
			var limiterCtx = ctx
			var cancel context.CancelFunc
			if timeout > 0 {
				limiterCtx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}
			select {
			case <-limiter:
				defer func() {
					limiter <- struct{}{}
				}()
			case <-limiterCtx.Done():
				return nil, status.Errorf(codes.ResourceExhausted,
					"%s is rejected by burstlimit stream client middleware, please retry later: %s", method, limiterCtx.Err())
			}
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}
