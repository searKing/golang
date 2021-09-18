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

// UnaryClientInterceptor returns a new unary client interceptor that performs request burst limiting.
// b 令牌桶大小, <0 无限制
// timeout 获取令牌超时返回时间, <0 无限制
func UnaryClientInterceptor(b int, timeout time.Duration) grpc.UnaryClientInterceptor {
	limiter := rate.NewFullBurstLimiter(b)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var limiterCtx = ctx
		var cancel context.CancelFunc
		if timeout >= 0 {
			limiterCtx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		if b >= 0 {
			err := limiter.Wait(limiterCtx)
			if err != nil {
				return status.Errorf(codes.ResourceExhausted,
					"%s is rejected by burstlimit unary client middleware, please retry later: %w", method, err)
			}
			defer limiter.PutToken()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClientInterceptor returns a new streaming client interceptor that performs burst limiting on the request.
// b 令牌桶大小, <0 无限制
// timeout 获取令牌超时返回时间, <0 无限制
func StreamClientInterceptor(b int, timeout time.Duration) grpc.StreamClientInterceptor {
	limiter := rate.NewFullBurstLimiter(b)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
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
					"%s is rejected by burstlimit stream client middleware, please retry later: %w", method, err)
			}
			defer limiter.PutToken()
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}
