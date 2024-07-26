// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	slices_ "github.com/searKing/golang/go/exp/slices"
	http_ "github.com/searKing/golang/go/net/http"
	time_ "github.com/searKing/golang/go/time"
)

var (
	// SystemTag is tag representing an event inside gRPC call.
	SystemTag = []string{"protocol", "http"}
)

// HttpInterceptor returns a new unary http interceptors that optionally logs endpoint handling.
// Logger will read existing and write new logging.Fields available in current context.
// See `ExtractFields` and `InjectFields` for details.
func HttpInterceptor(l logging.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var cost time_.Cost
			cost.Start()
			var attrs []any
			attrs = append(attrs, slog.String(SystemTag[0], SystemTag[1]))
			attrs = append(attrs, slog.Time("http.start_time", time.Now()))
			attrs = append(attrs, extractLoggingFieldsFromHttpRequest(r)...)
			l.Log(r.Context(), logging.LevelInfo, fmt.Sprintf("http request received"), attrs...)

			rw := http_.NewRecordResponseWriter(w)
			handler.ServeHTTP(rw, r)

			attrs = append(attrs, slog.String("http.status_code", slices_.FirstOrZero(http.StatusText(rw.Status()), "CODE("+strconv.FormatInt(int64(rw.Status()), 10)+")")),
				slog.Duration("cost", cost.Elapse()),
				slog.Int64("http.request_body_size", r.ContentLength),
				slog.Int("http.response_body_size", rw.Size()))
			l.Log(r.Context(), logging.LevelInfo, fmt.Sprintf("finished http call with code %d", rw.Status()), attrs...)
		})
	}
}

func extractLoggingFieldsFromHttpRequest(r *http.Request) []any {
	attrs := logging.ExtractFields(r.Context())
	if slog.Default().Enabled(r.Context(), slog.LevelDebug) {
		if d, ok := r.Context().Deadline(); ok {
			attrs = append(attrs, slog.Time("http.request.deadline", d))
		}
	}
	attrs = append(attrs, slog.String("http.remote_addr", r.RemoteAddr))
	ip := http_.ClientIP(r)
	if ip != "" && !strings.HasPrefix(r.RemoteAddr, ip) {
		attrs = append(attrs, slog.String("http.client_ip", http_.ClientIP(r)))
	}

	attrs = append(attrs, slog.String("http.method", r.Method), slog.String("http.request.uri", r.RequestURI))

	absRequestURI := strings.HasPrefix(r.RequestURI, "http://") || strings.HasPrefix(r.RequestURI, "https://")
	if !absRequestURI {
		host := r.Host
		if host == "" && r.URL != nil {
			host = r.URL.Host
		}
		if host != "" {
			attrs = append(attrs, slog.String("http.host", host))
		}
	}
	return attrs
}
