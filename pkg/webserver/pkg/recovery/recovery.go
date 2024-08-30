// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/searKing/golang/pkg/webserver/pkg/logging"
)

// UnaryServerInterceptor returns a new unary server interceptor that performs recovering from a panic.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandlerContext(recoveryLogHandler))
}

// StreamServerInterceptor returns a new stream server interceptor that performs recovering from a panic.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandlerContext(recoveryLogHandler))
}

// UnaryHandler returns a new unary server handler that performs recovering from a panic.
func UnaryHandler(handler grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req any) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoveryLogHandler(ctx, r)
			}
		}()
		resp, err := handler(ctx, req)
		return resp, err
	}
}

func recoveryLogHandler(ctx context.Context, p any) (err error) {
	{
		_, _ = os.Stderr.Write([]byte(fmt.Sprintf("panic: %s", p)))
		debug.PrintStack()
		_, _ = os.Stderr.Write([]byte(" [recovered]\n"))
	}
	logger := slog.With(logging.Attrs[any](ctx)...)
	logger.With(slog_.Error(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack()))).Error("recovered in grpc")
	return status.Errorf(codes.Internal, "%s", p)
}
