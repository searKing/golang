// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpclogrus

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	"github.com/sirupsen/logrus"
)

func WithLogrusLogger(logger *logrus.Logger) []grpc_.GatewayOption {
	return WithLogrusLoggerConfig(logger, nil)
}

// interceptorLogrusLogger adapts logrus logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogrusLogger(l logrus.FieldLogger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make(map[string]any, len(fields)/2)
		i := logging.Fields(fields).Iterator()
		if i.Next() {
			k, v := i.At()
			f[k] = v
		}
		l := l.WithFields(f)
		if lvl < logging.LevelInfo {
			l.Debug(msg)
		} else if lvl < logging.LevelWarn {
			l.Info(msg)
		} else if lvl < logging.LevelError {
			l.Warn(msg)
		} else {
			l.Error(msg)
		}
	})
}

func WithLogrusLoggerConfig(logger *logrus.Logger, slogOpts []logging.Option) []grpc_.GatewayOption {
	// interceptor's log below
	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithFieldsFromContext(grpc_.FieldsFromContextWithForward),
		// Add any other option (check functions starting with logging.With).
	}
	loggerOpts = append(loggerOpts, slogOpts...)
	l := interceptorLogrusLogger(logger)

	var opts []grpc_.GatewayOption
	opts = append(opts, grpc_.WithGrpcStreamServerChain(logging.StreamServerInterceptor(l, loggerOpts...)))
	opts = append(opts, grpc_.WithGrpcUnaryServerChain(logging.UnaryServerInterceptor(l, loggerOpts...)))
	opts = append(opts, grpc_.WithGrpcStreamClientChain(logging.StreamClientInterceptor(l, loggerOpts...)))
	opts = append(opts, grpc_.WithGrpcUnaryClientChain(logging.UnaryClientInterceptor(l, loggerOpts...)))
	opts = append(opts, grpc_.WithHttpWrapper(grpc_.HttpInterceptor(l)))
	return opts
}
