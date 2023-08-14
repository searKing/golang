// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
)

func WithSlogLogger(logger slog.Handler) GatewayOption {
	return WithSlogLoggerConfig(logger, gin.LoggerConfig{})
}

// interceptorSlogLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorSlogLogger(h slog.Handler) logging.Logger {
	l := slog.New(h)
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func WithSlogLoggerConfig(logger slog.Handler, conf gin.LoggerConfig) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		l := logger.WithGroup("grpc-gateway")
		if conf.Formatter == nil {
			conf.Formatter = gin_.LogFormatter("HTTP")
		}
		if conf.Output == nil {
			conf.Output = slog.NewLogLogger(logger, slog.LevelInfo).Writer()
		}

		loggerOpts := []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			logging.WithFieldsFromContext(FieldsFromContextWithForward),
			// Add any other option (check functions starting with logging.With).
		}

		var opts []GatewayOption
		opts = append(opts, WithGrpcStreamServerChain(logging.StreamServerInterceptor(interceptorSlogLogger(l), loggerOpts...)))
		opts = append(opts, WithGrpcUnaryServerChain(logging.UnaryServerInterceptor(interceptorSlogLogger(l), loggerOpts...)))
		opts = append(opts, WithGrpcStreamClientChain(logging.StreamClientInterceptor(interceptorSlogLogger(l), loggerOpts...)))
		opts = append(opts, WithGrpcUnaryClientChain(logging.UnaryClientInterceptor(interceptorSlogLogger(l), loggerOpts...)))
		opts = append(opts, WithHttpWrapper(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				g := gin.New()
				g.Use(gin.LoggerWithConfig(conf))
				g.Any("/*path", gin.WrapH(handler))
				g.ServeHTTP(w, r)
			})
		}))
		gateway.ApplyOptions(opts...)
	})
}
