// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	http_ "github.com/searKing/golang/go/net/http"
	time_ "github.com/searKing/golang/go/time"
	grpclog_ "github.com/searKing/golang/third_party/google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/grpclog"
)

func WithSlogLogger(logger slog.Handler) GatewayOption {
	return WithSlogLoggerConfig(logger, nil)
}

// interceptorSlogLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorSlogLogger(h slog.Handler) logging.Logger {
	l := slog.New(h)
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func WithSlogLoggerConfig(logger slog.Handler, slogOpts []logging.Option) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		l := logger.WithGroup("grpc-gateway")
		grpclog.SetLoggerV2(grpclog_.NewSlogger(l))

		// interceptor's log below
		loggerOpts := []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			logging.WithFieldsFromContext(FieldsFromContextWithForward),
			// Add any other option (check functions starting with logging.With).
		}
		loggerOpts = append(loggerOpts, slogOpts...)

		var opts []GatewayOption
		opts = append(opts, WithGrpcStreamServerChain(logging.StreamServerInterceptor(interceptorSlogLogger(l), loggerOpts...)))
		opts = append(opts, WithGrpcUnaryServerChain(logging.UnaryServerInterceptor(interceptorSlogLogger(l), loggerOpts...)))
		opts = append(opts, WithGrpcStreamClientChain(logging.StreamClientInterceptor(interceptorSlogLogger(l), loggerOpts...)))
		opts = append(opts, WithGrpcUnaryClientChain(logging.UnaryClientInterceptor(interceptorSlogLogger(l), loggerOpts...)))
		opts = append(opts, WithHttpWrapper(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var cost time_.Cost
				cost.Start()

				rw := http_.NewRecordResponseWriter(w)
				handler.ServeHTTP(rw, r)
				attrs := logging.ExtractFields(r.Context())
				attrs = append(attrs, slog.String("remote_addr", r.RemoteAddr))
				ip := http_.ClientIP(r)
				if ip != "" && !strings.HasPrefix(r.RemoteAddr, ip) {
					attrs = append(attrs, slog.String("client_ip", http_.ClientIP(r)))
				}

				attrs = append(attrs, slog.String("method", r.Method), slog.String("request_uri", r.RequestURI))

				absRequestURI := strings.HasPrefix(r.RequestURI, "http://") || strings.HasPrefix(r.RequestURI, "https://")
				if !absRequestURI {
					host := r.Host
					if host == "" && r.URL != nil {
						host = r.URL.Host
					}
					if host != "" {
						attrs = append(attrs, slog.String("host", host))
					}
				}
				attrs = append(attrs, slog.String("status_code", http.StatusText(rw.Status())),
					slog.Int64("request_body_size", r.ContentLength),
					slog.Int("response_body_size", rw.Size()))
				slog.With(attrs...).Info(fmt.Sprintf("finished http call with code %d", rw.Status()))
			})
		}))
		gateway.ApplyOptions(opts...)
	})
}
