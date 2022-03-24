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

// UnaryServerInterceptor returns a new unary server interceptors that performs request rate limiting.
func UnaryServerInterceptor(r rate.Limit, b int) grpc.UnaryServerInterceptor {
	limiter := rate.NewLimiter(r, b)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Allow() {
			return handler(ctx, req)
		}
		return nil, status.Errorf(codes.ResourceExhausted,
			"%s is rejected by timeratelimit unary server middleware, please retry later", info.FullMethod)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that performs rate limiting on the request.
func StreamServerInterceptor(r rate.Limit, b int) grpc.StreamServerInterceptor {
	limiter := rate.NewLimiter(r, b)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if limiter.Allow() {
			return handler(srv, stream)
		}
		return status.Errorf(codes.ResourceExhausted,
			"%s is rejected by timeratelimit stream server middleware, please retry later", info.FullMethod)
	}
}
