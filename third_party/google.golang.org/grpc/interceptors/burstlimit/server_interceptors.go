// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package burstlimit

import (
	"context"
	"time"

	"github.com/searKing/golang/go/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptors that performs request burst limiting.
// b 令牌桶大小, <0 无限制
// timeout 获取令牌超时返回时间, <0 无限制
func UnaryServerInterceptor(b int, timeout time.Duration) grpc.UnaryServerInterceptor {
	limiter := rate.NewFullBurstLimiter(b)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var limiterCtx = ctx
		var cancel context.CancelFunc
		if timeout >= 0 {
			limiterCtx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		if b >= 0 {
			err := limiter.Wait(limiterCtx)
			if err != nil {
				return nil, status.Errorf(codes.ResourceExhausted,
					"%s is rejected by burstlimit unary server middleware, please retry later: %w", info.FullMethod, err)
			}
			defer limiter.PutToken()
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that performs burst limiting on the request.
// b 令牌桶大小, <0 无限制
// timeout 获取令牌超时返回时间, <0 无限制
func StreamServerInterceptor(b int, timeout time.Duration) grpc.StreamServerInterceptor {
	limiter := rate.NewFullBurstLimiter(b)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var limiterCtx = stream.Context()
		var cancel context.CancelFunc
		if timeout >= 0 {
			limiterCtx, cancel = context.WithTimeout(limiterCtx, timeout)
			defer cancel()
		}
		if b >= 0 {
			err := limiter.Wait(limiterCtx)
			if err != nil {
				return status.Errorf(codes.ResourceExhausted,
					"%s is rejected by burstlimit stream server middleware, please retry later: %w", info.FullMethod, err)
			}
			defer limiter.PutToken()
		}
		return handler(srv, stream)
	}
}
