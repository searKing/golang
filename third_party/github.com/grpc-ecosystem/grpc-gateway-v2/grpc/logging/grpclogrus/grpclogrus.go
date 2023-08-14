// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpclogrus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	"github.com/sirupsen/logrus"
)

func WithLogrusLogger(logger *logrus.Logger) grpc.GatewayOption {
	return WithLogrusLoggerConfig(logger, gin.LoggerConfig{})
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

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func WithLogrusLoggerConfig(logger *logrus.Logger, conf gin.LoggerConfig) grpc.GatewayOption {
	return grpc.GatewayOptionFunc(func(gateway *grpc.Gateway) {
		if conf.Formatter == nil {
			conf.Formatter = gin_.LogFormatter("gRPC over HTTP")
		}
		if conf.Output == nil {
			conf.Output = logger.Writer()
		}
		loggerOpts := []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			logging.WithFieldsFromContext(grpc.FieldsFromContextWithForward),
			// Add any other option (check functions starting with logging.With).
		}

		var opts []grpc.GatewayOption
		opts = append(opts, grpc.WithGrpcStreamServerChain(logging.StreamServerInterceptor(interceptorLogrusLogger(logger), loggerOpts...)))
		opts = append(opts, grpc.WithGrpcUnaryServerChain(logging.UnaryServerInterceptor(interceptorLogrusLogger(logger), loggerOpts...)))
		opts = append(opts, grpc.WithGrpcStreamClientChain(logging.StreamClientInterceptor(interceptorLogrusLogger(logger), loggerOpts...)))
		opts = append(opts, grpc.WithGrpcUnaryClientChain(logging.UnaryClientInterceptor(interceptorLogrusLogger(logger), loggerOpts...)))
		opts = append(opts, grpc.WithHttpWrapper(func(handler http.Handler) http.Handler {
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
