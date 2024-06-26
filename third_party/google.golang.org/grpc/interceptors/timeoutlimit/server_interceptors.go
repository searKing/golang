// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timeoutlimit

import (
	"context"
	"time"

	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors with timeout limit of handle.
// take effect if timeout > 0
func UnaryServerInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptors with timeout limit of handle.
// take effect if timeout > 0
func StreamServerInterceptor(timeout time.Duration) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(ss.Context(), timeout)
			defer cancel()
			wrapped := interceptors.WrapServerStream(ss)
			wrapped.WrappedContext = ctx
			ss = wrapped
		}
		return handler(srv, ss)
	}
}
