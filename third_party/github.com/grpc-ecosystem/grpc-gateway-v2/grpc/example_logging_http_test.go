// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

// exampleLogger logs to stdout only the message and fields of headers, to keep the output stable.
func exampleLogger() logging.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey || strings.HasPrefix(a.Key, "http.request.header.") {
				return a
			}
			return slog.Attr{}
		},
	}))
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		logger.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// authScheme returns the auth scheme of Authorization, such as Bearer, without credentials.
func authScheme(authorization string) string {
	scheme, _, _ := strings.Cut(authorization, " ")
	return scheme
}

// This example logs headers chosen by the caller, rather than all headers with sensitive ones redacted
// by HTTP_GO_LOG_HTTP_HEADER, which is left unset. A middleware before HttpInterceptor injects the headers
// wanted into the request context by logging.InjectFields, which are logged with the request.
func ExampleHttpInterceptor_customHeaders() {
	injectHeaders := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logging.InjectFields(r.Context(), logging.Fields{
				"http.request.header.X-Request-Id", r.Header.Get("X-Request-Id"),
				"http.request.header.Authorization", authScheme(r.Header.Get("Authorization")),
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("ok")) })
	h := injectHeaders(grpc_.HttpInterceptor(exampleLogger())(handler))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-Request-Id", "req-1")
	r.Header.Set("Authorization", "Bearer secret")
	h.ServeHTTP(httptest.NewRecorder(), r)

	// Output:
	// msg="http request received" http.request.header.X-Request-Id=req-1 http.request.header.Authorization=Bearer
	// msg="finished http call with status code 200" http.request.header.X-Request-Id=req-1 http.request.header.Authorization=Bearer
}

// This example logs headers chosen by the caller of a http client, injecting them into the request context
// by logging.InjectFields before sending, as ExampleHttpInterceptor_customHeaders does for a http server.
func ExampleHttpRoundTripDecorator_customHeaders() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()
	client := &http.Client{Transport: grpc_.HttpRoundTripDecorator(exampleLogger()).WrapRoundTrip(http.DefaultTransport)}

	r, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	r.Header.Set("X-Request-Id", "req-1")
	r = r.WithContext(logging.InjectFields(r.Context(), logging.Fields{
		"http.request.header.X-Request-Id", r.Header.Get("X-Request-Id"),
	}))
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	_ = resp.Body.Close()

	// Output:
	// msg="http request sending" http.request.header.X-Request-Id=req-1
	// msg="finished http call with status code 200" http.request.header.X-Request-Id=req-1
}
