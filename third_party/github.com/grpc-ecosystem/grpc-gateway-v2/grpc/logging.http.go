// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	slices_ "github.com/searKing/golang/go/exp/slices"
	http_ "github.com/searKing/golang/go/net/http"
	"google.golang.org/grpc/codes"
)

var (
	// SystemTag is tag representing an event inside gRPC call.
	SystemTag = []string{"protocol", "http"}
	// ComponentFieldKey is a tag representing the client/server that is calling,
	// valued as logging.KindServerFieldValue or logging.KindClientFieldValue.
	ComponentFieldKey = "http.component"
)

// HttpInterceptor returns a new unary http interceptors that optionally logs endpoint handling.
// Logger will read existing and write new logging.Fields available in current context.
// See `ExtractFields` and `InjectFields` for details.
// Logs and fields are consistent with logging.UnaryServerInterceptor of gRPC, such as levels by
// logging.DefaultServerCodeToLevel of the gRPC code mapped from the http status code by DefaultHttpToCode,
// customized by opts as logging.Option does.
// Headers are logged if HTTP_GO_LOG_HTTP_HEADER is true, with values of sensitive ones redacted,
// see ExampleHttpInterceptor_customHeaders to log headers chosen instead.
func HttpInterceptor(l logging.Logger, opts ...HttpLoggingOption) func(handler http.Handler) http.Handler {
	o := evaluateHttpLoggingOptions(logging.DefaultServerCodeToLevel, opts)
	logHttpHeader := logHttpHeaderFromEnv()
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx, fields := injectLoggingFields(r, logging.KindServerFieldValue, start, o)
			r = r.WithContext(ctx)

			if o.has(logging.StartCall) {
				reqFields := fields
				if logHttpHeader {
					reqFields = append(reqFields, httpHeaderToFields(r.Header, "http.request.header")...)
				}
				l.Log(ctx, o.levelFunc(codes.OK), "started call", reqFields...)
			}

			rw := http_.NewResponseWriterDelegator(w)
			handler.ServeHTTP(rw, r)
			if !o.has(logging.FinishCall) {
				return
			}
			status := rw.Status()
			if status == 0 {
				// net/http replies 200 if the handler writes nothing
				status = http.StatusOK
			}

			fields = fields.WithUnique(logging.ExtractFields(ctx))
			fields = fields.AppendUnique(logging.Fields{
				"http.status_code", httpStatusText(status),
				"http.request_body_size", r.ContentLength,
				"http.response_body_size", rw.Written(),
			})
			if o.fieldsFromRequestFn != nil {
				// fieldsFromRequestFn dups override the existing fields.
				fields = o.fieldsFromRequestFn(ctx, r).AppendUnique(fields)
			}
			fields = fields.AppendUnique(o.durationFieldFunc(time.Since(start)))
			if logHttpHeader {
				fields = append(fields, httpHeaderToFields(rw.Header(), "http.response.header")...)
			}
			l.Log(ctx, o.levelFunc(o.codeFunc(status, nil)), "finished call", fields...)
		})
	}
}

// HttpRoundTripDecorator returns a new http RoundTripDecorator that optionally logs outgoing HTTP requests.
// Logs and fields are consistent with logging.UnaryClientInterceptor of gRPC, such as levels by
// logging.DefaultClientCodeToLevel of the gRPC code mapped from the http status code or the error by
// DefaultHttpToCode, customized by opts as logging.Option does.
// Headers are logged if HTTP_GO_LOG_HTTP_HEADER is true, with values of sensitive ones redacted,
// see ExampleHttpRoundTripDecorator_customHeaders to log headers chosen instead.
func HttpRoundTripDecorator(l logging.Logger, opts ...HttpLoggingOption) http_.RoundTripDecorator {
	o := evaluateHttpLoggingOptions(logging.DefaultClientCodeToLevel, opts)
	logHttpHeader := logHttpHeaderFromEnv()
	return http_.RoundTripDecoratorFunc(func(rt http.RoundTripper) http.RoundTripper {
		return http_.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			start := time.Now()
			ctx, fields := injectLoggingFields(r, logging.KindClientFieldValue, start, o)

			if o.has(logging.StartCall) {
				reqFields := fields
				if logHttpHeader {
					reqFields = append(reqFields, httpHeaderToFields(r.Header, "http.request.header")...)
				}
				l.Log(ctx, o.levelFunc(codes.OK), "started call", reqFields...)
			}

			resp, err := rt.RoundTrip(r.WithContext(ctx))
			if !o.has(logging.FinishCall) {
				return resp, err
			}

			fields = fields.WithUnique(logging.ExtractFields(ctx))
			fields = fields.AppendUnique(logging.Fields{"http.request_body_size", r.ContentLength})
			var code codes.Code
			if err != nil {
				code = o.codeFunc(0, err)
				fields = fields.AppendUnique(logging.Fields{"http.error", fmt.Sprintf("%v", err)})
				if o.errorToFieldsFunc != nil {
					fields = fields.AppendUnique(o.errorToFieldsFunc(err))
				}
			} else {
				code = o.codeFunc(resp.StatusCode, nil)
				fields = fields.AppendUnique(logging.Fields{
					"http.status_code", httpStatusText(resp.StatusCode),
					"http.response_body_size", resp.ContentLength,
				})
			}
			if o.fieldsFromRequestFn != nil {
				// fieldsFromRequestFn dups override the existing fields.
				fields = o.fieldsFromRequestFn(ctx, r).AppendUnique(fields)
			}
			fields = fields.AppendUnique(o.durationFieldFunc(time.Since(start)))
			if logHttpHeader && resp != nil {
				fields = append(fields, httpHeaderToFields(resp.Header, "http.response.header")...)
			}
			l.Log(ctx, o.levelFunc(code), "finished call", fields...)
			return resp, err
		})
	})
}

func logHttpHeaderFromEnv() bool {
	v, _ := strconv.ParseBool(os.Getenv("HTTP_GO_LOG_HTTP_HEADER"))
	return v
}

// injectLoggingFields injects common fields of r into its context, as logging interceptors of gRPC do,
// returning the context and the fields to log, plus fields of start time and deadline used only once.
func injectLoggingFields(r *http.Request, kind string, start time.Time, o *httpLoggingOptions) (context.Context, logging.Fields) {
	ctx := r.Context()
	fields := logging.Fields{SystemTag[0], SystemTag[1], ComponentFieldKey, kind}

	absRequestURI := strings.HasPrefix(r.RequestURI, "http://") || strings.HasPrefix(r.RequestURI, "https://")
	uri := r.RequestURI
	if (uri == "" || absRequestURI) && r.URL != nil {
		// RequestURI is unset for client requests, and may contain a password if absolute
		uri = r.URL.Redacted()
	}
	fields = append(fields, "http.method", r.Method, "http.request.uri", uri)
	if !absRequestURI {
		host := r.Host
		if host == "" && r.URL != nil {
			host = r.URL.Host
		}
		if host != "" {
			fields = append(fields, "http.host", host)
		}
	}
	if kind == logging.KindServerFieldValue {
		fields = append(fields, "http.remote_addr", r.RemoteAddr)
		if ip := http_.ClientIP(r); ip != "" && !strings.HasPrefix(r.RemoteAddr, ip) {
			fields = append(fields, "http.client_ip", ip)
		}
	}
	for _, key := range o.disableHttpLogFields {
		fields.Delete(key)
	}
	fields = fields.WithUnique(logging.ExtractFields(ctx))
	if o.fieldsFromRequestFn != nil {
		// fieldsFromRequestFn dups override the existing fields.
		fields = o.fieldsFromRequestFn(ctx, r).AppendUnique(fields)
	}
	ctx = logging.InjectFields(ctx, fields)

	singleUseFields := logging.Fields{"http.start_time", start.Format(o.timestampFormat)}
	if d, ok := ctx.Deadline(); ok {
		singleUseFields = singleUseFields.AppendUnique(logging.Fields{"http.request.deadline", d.Format(o.timestampFormat)})
	}
	return ctx, fields.WithUnique(singleUseFields)
}

// httpStatusText returns the text of the http status code, as "OK" of grpc.code for gRPC.
func httpStatusText(code int) string {
	return slices_.FirstOrZero(http.StatusText(code), "CODE("+strconv.Itoa(code)+")")
}

// httpStatusToCode maps the http status code to the gRPC code, reversing runtime.HTTPStatusFromCode of grpc-gateway,
// and the less severe code if ambiguous.
func httpStatusToCode(status int) codes.Code {
	switch status {
	case 499:
		return codes.Canceled
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusInternalServerError:
		return codes.Internal
	}
	switch {
	case status < http.StatusBadRequest:
		return codes.OK
	case status < http.StatusInternalServerError:
		return codes.InvalidArgument
	default:
		return codes.Unknown
	}
}

// httpErrorToCode maps the error of a http round trip to the gRPC code,
// as gRPC reports Unavailable if the transport fails.
func httpErrorToCode(err error) codes.Code {
	switch {
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return codes.DeadlineExceeded
	}
	return codes.Unavailable
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

// httpHeaderToFields groups headers h as k, with values of sensitive headers redacted.
func httpHeaderToFields(h http.Header, k string) logging.Fields {
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
	return logging.Fields{k, slog.GroupValue(attrs...)}
}
