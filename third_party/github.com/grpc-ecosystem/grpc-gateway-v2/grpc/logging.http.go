// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
// Headers are logged if HTTP_GO_LOG_HTTP_HEADER is true, with values of sensitive ones redacted,
// see ExampleHttpInterceptor_customHeaders to log headers chosen instead.
func HttpInterceptor(l logging.Logger) func(handler http.Handler) http.Handler {
	var logHttpHeader bool
	{
		vHeader := os.Getenv("HTTP_GO_LOG_HTTP_HEADER")
		if vh, err := strconv.ParseBool(vHeader); err == nil {
			logHttpHeader = vh
		}
	}

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var cost time_.Cost
			cost.Start()
			var attrs []any
			attrs = append(attrs, slog.String(SystemTag[0], SystemTag[1]))
			attrs = append(attrs, slog.Time("http.start_time", time.Now()))
			attrs = append(attrs, extractLoggingFieldsFromHttpRequest(r)...)

			var reqAttrs = attrs
			if logHttpHeader {
				reqAttrs = append(reqAttrs, httpHeaderToAttr(r.Header, "http.request.header"))
			}
			l.Log(r.Context(), logging.LevelInfo, fmt.Sprintf("http request received"), reqAttrs...)

			rw := http_.NewResponseWriterDelegator(w)
			handler.ServeHTTP(rw, r)

			attrs = append(attrs,
				slog.String("http.status_code", slices_.FirstOrZero(http.StatusText(rw.Status()), "CODE("+strconv.FormatInt(int64(rw.Status()), 10)+")")),
				slog.Duration("cost", cost.Elapse()),
				slog.Int64("http.request_body_size", r.ContentLength),
				slog.Int64("http.response_body_size", rw.Written()))

			var respAttrs = attrs
			if logHttpHeader {
				respAttrs = append(respAttrs, httpHeaderToAttr(rw.Header(), "http.response.header"))
			}
			l.Log(r.Context(), logging.LevelInfo, fmt.Sprintf("finished http call with status code %d", rw.Status()),
				respAttrs...)
		})
	}
}

// HttpRoundTripDecorator returns a new http RoundTripDecorator that optionally logs outgoing HTTP requests.
// Headers are logged if HTTP_GO_LOG_HTTP_HEADER is true, with values of sensitive ones redacted,
// see ExampleHttpRoundTripDecorator_customHeaders to log headers chosen instead.
func HttpRoundTripDecorator(l logging.Logger) http_.RoundTripDecorator {
	var logHttpHeader bool
	{
		vHeader := os.Getenv("HTTP_GO_LOG_HTTP_HEADER")
		if vh, err := strconv.ParseBool(vHeader); err == nil {
			logHttpHeader = vh
		}
	}

	return http_.RoundTripDecoratorFunc(func(rt http.RoundTripper) http.RoundTripper {
		return http_.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			var cost time_.Cost
			cost.Start()
			var attrs []any
			attrs = append(attrs, slog.String(SystemTag[0], SystemTag[1]))
			attrs = append(attrs, slog.Time("http.start_time", time.Now()))
			attrs = append(attrs, extractLoggingFieldsFromHttpRequest(r)...)

			var reqAttrs = attrs
			if logHttpHeader {
				reqAttrs = append(reqAttrs, httpHeaderToAttr(r.Header, "http.request.header"))
			}
			l.Log(r.Context(), logging.LevelInfo, "http request sending", reqAttrs...)

			resp, err := rt.RoundTrip(r)

			attrs = append(attrs,
				slog.Duration("cost", cost.Elapse()),
				slog.Int64("http.request_body_size", r.ContentLength))

			if resp != nil {
				attrs = append(attrs,
					slog.String("http.status_code", slices_.FirstOrZero(http.StatusText(resp.StatusCode), "CODE("+strconv.FormatInt(int64(resp.StatusCode), 10)+")")),
					slog.Int64("http.response_body_size", resp.ContentLength))
			}
			if err != nil {
				attrs = append(attrs, slog.String("http.error", err.Error()))
			}

			var respAttrs = attrs
			if logHttpHeader && resp != nil {
				respAttrs = append(respAttrs, httpHeaderToAttr(resp.Header, "http.response.header"))
			}

			if err != nil {
				l.Log(r.Context(), logging.LevelError, "finished http call with error", respAttrs...)
			} else {
				l.Log(r.Context(), logging.LevelInfo, fmt.Sprintf("finished http call with status code %d", resp.StatusCode), respAttrs...)
			}
			return resp, err
		})
	})
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

// redacted replaces values of sensitive headers in logs, as [net/url.URL.Redacted] does for passwords.
const redacted = "xxxxx"

// isSensitiveHeader reports whether the header name carries credentials, whose values are redacted in logs.
// Sensitive headers refer to the Go standard library, which strips them on redirect to other domains,
// see (*Client).makeHeadersCopier in net/http/client.go, plus Set-Cookie of responses.
//
// To customize headers logged, unset HTTP_GO_LOG_HTTP_HEADER, and inject the headers wanted by
// logging.InjectFields into the request context, in a middleware before HttpInterceptor,
// or before sending by HttpRoundTripDecorator, as fields of the context are logged.
func isSensitiveHeader(name string) bool {
	switch http.CanonicalHeaderKey(name) {
	case "Authorization", "Www-Authenticate", "Cookie", "Cookie2",
		"Proxy-Authorization", "Proxy-Authenticate", "Set-Cookie":
		return true
	}
	return false
}

// redactHeaderValue redacts v of the sensitive header name,
// keeping the auth scheme of Authorization and alike, such as "Bearer xxxxx".
func redactHeaderValue(name, v string) string {
	switch http.CanonicalHeaderKey(name) {
	case "Authorization", "Www-Authenticate", "Proxy-Authorization", "Proxy-Authenticate":
		if scheme, _, ok := strings.Cut(v, " "); ok && scheme != "" {
			return scheme + " " + redacted
		}
	}
	return redacted
}

// httpHeaderToAttr groups headers h as k, with values of sensitive headers redacted.
func httpHeaderToAttr(h http.Header, k string) slog.Attr {
	var attrs []slog.Attr
	for name, vs := range h {
		if isSensitiveHeader(name) {
			redactedVs := make([]string, len(vs))
			for i, v := range vs {
				redactedVs[i] = redactHeaderValue(name, v)
			}
			vs = redactedVs
		}
		attrs = append(attrs, slog.Any(name, vs))
	}
	return slog.GroupAttrs(k, attrs...)
}
