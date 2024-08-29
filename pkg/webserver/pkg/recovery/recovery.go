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
	slices_ "github.com/searKing/golang/go/exp/slices"
	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/searKing/golang/pkg/webserver/pkg/logging"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptor that performs recovering from a panic.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandlerContext(recoveryLogHandler))
}

// StreamServerInterceptor returns a new stream server interceptor that performs recovering from a panic.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandlerContext(recoveryLogHandler))
}

// WrapRecovery returns a new unary server interceptor that performs recovering from a panic.
func WrapRecovery[REQ any, RESP any](handler grpc_.UnaryHandler[REQ, RESP]) grpc_.UnaryHandler[REQ, RESP] {
	return func(ctx context.Context, req REQ) (_ RESP, err error) {
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
	logger := slog.With(slices_.MapFunc(logging.Attrs(ctx), func(e slog.Attr) any { return e })...)
	logger.With(slog_.Error(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack()))).Error("recovered in grpc")
	return status.Errorf(codes.Internal, "%s", p)
}
