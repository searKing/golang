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

// UnaryServerInterceptor returns a new unary server interceptors that performs request burst limiting.
// b bucket size, take effect if b > 0
// timeout ResourceExhausted if cost more than timeout to get a token, take effect if timeout > 0
func UnaryServerInterceptor(b int, timeout time.Duration) grpc.UnaryServerInterceptor {
	limiter := make(chan struct{}, b)
	for i := 0; i < b; i++ {
		limiter <- struct{}{}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
					"%s is rejected by burstlimit unary server middleware, please retry later: %s", info.FullMethod, limiterCtx.Err())
			}
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that performs burst limiting on the request.
// b bucket size, take effect if b > 0
// timeout ResourceExhausted if cost more than timeout to get a token, take effect if timeout > 0
func StreamServerInterceptor(b int, timeout time.Duration) grpc.StreamServerInterceptor {
	limiter := make(chan struct{}, b)
	for i := 0; i < b; i++ {
		limiter <- struct{}{}
	}
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if b > 0 {
			var limiterCtx = stream.Context()
			var cancel context.CancelFunc
			if timeout > 0 {
				limiterCtx, cancel = context.WithTimeout(limiterCtx, timeout)
				defer cancel()
			}
			select {
			case <-limiter:
				defer func() {
					limiter <- struct{}{}
				}()
			case <-limiterCtx.Done():
				return status.Errorf(codes.ResourceExhausted,
					"%s is rejected by burstlimit stream server middleware, please retry later: %s", info.FullMethod, limiterCtx.Err())
			}
		}
		return handler(srv, stream)
	}
}
