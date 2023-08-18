// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	http_ "github.com/searKing/golang/go/net/http"
	time_ "github.com/searKing/golang/go/time"
)

// HttpInterceptor returns a new unary http interceptors that optionally logs endpoint handling.
// Logger will read existing and write new logging.Fields available in current context.
// See `ExtractFields` and `InjectFields` for details.
func HttpInterceptor(l logging.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
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
				slog.Duration("cost", cost.Elapse()),
				slog.Int64("request_body_size", r.ContentLength),
				slog.Int("response_body_size", rw.Size()))
			l.Log(r.Context(), logging.LevelInfo, fmt.Sprintf("finished http call with code %d", rw.Status()), attrs...)
		})
	}
}
