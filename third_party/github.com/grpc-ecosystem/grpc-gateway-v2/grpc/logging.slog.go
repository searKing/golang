// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	runtime_ "github.com/searKing/golang/go/runtime"
	grpclog_ "github.com/searKing/golang/third_party/google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/grpclog"
)

const d = 3

func WithSlogLogger(logger slog.Handler) []GatewayOption {
	return WithSlogLoggerConfig(logger, nil)
}

// interceptorSlogLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorSlogLogger(h slog.Handler) logging.Logger {
	l := slog.New(h)
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		//l.Log(ctx, slog.Level(lvl), msg, fields...)
		if l.Enabled(ctx, slog.Level(lvl)) {
			pc := runtime_.GetCallerFrame(d).PC
			r := slog.NewRecord(time.Now(), slog.Level(lvl), msg, pc)
			r.Add(fields...)
			_ = h.Handle(ctx, r)
		}
	})
}

func WithSlogLoggerConfig(h slog.Handler, slogOpts []logging.Option) []GatewayOption {
	grpclog.SetLoggerV2(grpclog_.NewSlogger(h))

	// interceptor's log below
	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithFieldsFromContext(FieldsFromContextWithForward),
		// Add any other option (check functions starting with logging.With).
	}
	loggerOpts = append(loggerOpts, slogOpts...)
	l := interceptorSlogLogger(h)

	var opts []GatewayOption
	opts = append(opts, WithHttpWrapper(HttpInterceptor(l)))
	opts = append(opts, WithGrpcStreamServerChain(logging.StreamServerInterceptor(l, loggerOpts...)))
	opts = append(opts, WithGrpcUnaryServerChain(logging.UnaryServerInterceptor(l, loggerOpts...)))
	opts = append(opts, WithGrpcStreamClientChain(logging.StreamClientInterceptor(l, loggerOpts...)))
	opts = append(opts, WithGrpcUnaryClientChain(logging.UnaryClientInterceptor(l, loggerOpts...)))
	return opts
}
