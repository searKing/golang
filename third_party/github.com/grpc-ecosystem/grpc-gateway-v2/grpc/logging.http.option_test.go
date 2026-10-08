// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/codes"

	http_ "github.com/searKing/golang/go/net/http"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

// serveHttp serves a request with headers h by HttpInterceptor customized by opts, replying status,
// returning logs and fields injected into the handler context.
func serveHttp(t *testing.T, status int, h http.Header, opts ...grpc_.HttpLoggingOption) (recordLogger, map[string]string) {
	t.Helper()
	var rl recordLogger
	var injected map[string]string
	handler := grpc_.HttpInterceptor(&rl, opts...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var hl recordLogger
		hl.Log(r.Context(), logging.LevelInfo, "", logging.ExtractFields(r.Context())...)
		injected = hl[0].fields
		w.WriteHeader(status)
	}))
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	for k, vs := range h {
		r.Header[k] = vs
	}
	handler.ServeHTTP(httptest.NewRecorder(), r)
	return rl, injected
}

func TestHttpInterceptor_WithHttpLogOnEvents(t *testing.T) {
	rl, _ := serveHttp(t, http.StatusOK, nil, grpc_.WithHttpLogOnEvents(logging.FinishCall))
	if len(rl) != 1 {
		t.Fatalf("logged %d, want 1", len(rl))
	}
	checkLog(t, rl[0], logging.LevelInfo, "finished call", map[string]string{"http.status_code": "OK"})

	rl, _ = serveHttp(t, http.StatusOK, nil, grpc_.WithHttpLogOnEvents(logging.StartCall))
	if len(rl) != 1 {
		t.Fatalf("logged %d, want 1", len(rl))
	}
	checkLog(t, rl[0], logging.LevelInfo, "started call", nil)
}

func TestHttpInterceptor_WithHttpLevels(t *testing.T) {
	var gotCodes []codes.Code
	rl, _ := serveHttp(t, http.StatusNotFound, nil, grpc_.WithHttpLevels(func(code codes.Code) logging.Level {
		gotCodes = append(gotCodes, code)
		return logging.LevelError
	}))
	checkLog(t, rl[0], logging.LevelError, "started call", nil)
	checkLog(t, rl[1], logging.LevelError, "finished call", nil)
	if want := []codes.Code{codes.OK, codes.NotFound}; len(gotCodes) != 2 || gotCodes[0] != want[0] || gotCodes[1] != want[1] {
		t.Errorf("codes = %v, want %v", gotCodes, want)
	}
}

func TestHttpInterceptor_WithHttpCodes(t *testing.T) {
	rl, _ := serveHttp(t, http.StatusOK, nil, grpc_.WithHttpCodes(func(status int, err error) codes.Code {
		if status == http.StatusOK && err == nil {
			return codes.Internal
		}
		return codes.OK
	}))
	checkLog(t, rl[1], logging.LevelError, "finished call", map[string]string{"http.status_code": "OK"})
}

func TestHttpInterceptor_WithHttpDurationField(t *testing.T) {
	rl, _ := serveHttp(t, http.StatusOK, nil, grpc_.WithHttpDurationField(grpc_.HttpDurationToDurationField))
	checkLog(t, rl[1], logging.LevelInfo, "finished call", map[string]string{"http.duration": anyValue, "http.time_ms": ""})
}

func TestHttpInterceptor_WithHttpTimestampFormat(t *testing.T) {
	rl, _ := serveHttp(t, http.StatusOK, nil, grpc_.WithHttpTimestampFormat("2006"))
	if got := rl[0].fields["http.start_time"]; !regexp.MustCompile(`^\d{4}$`).MatchString(got) {
		t.Errorf("http.start_time = %q, want in format 2006", got)
	}
}

func TestHttpInterceptor_WithHttpDisableLoggingFields(t *testing.T) {
	disabled := map[string]string{"http.method": "", "http.remote_addr": "", "protocol": "http"}
	rl, injected := serveHttp(t, http.StatusOK, nil, grpc_.WithHttpDisableLoggingFields("http.method", "http.remote_addr"))
	checkFields(t, "fields injected into the handler context", injected, disabled)
	checkLog(t, rl[0], logging.LevelInfo, "started call", disabled)
	checkLog(t, rl[1], logging.LevelInfo, "finished call", disabled)
}

func TestHttpInterceptor_WithHttpFieldsFromContextAndRequest(t *testing.T) {
	h := http.Header{"X-Forwarded-For": {"203.0.113.1"}, "X-Forwarded-Host": {"example.org"}}
	forwarded := map[string]string{"http.x-forwarded-for": "[203.0.113.1]", "http.x-forwarded-host": "[example.org]"}
	rl, injected := serveHttp(t, http.StatusOK, h, grpc_.WithHttpFieldsFromContextAndRequest(grpc_.HttpFieldsFromRequestWithForward))
	checkFields(t, "fields injected into the handler context", injected, forwarded)
	checkLog(t, rl[0], logging.LevelInfo, "started call", forwarded)
	checkLog(t, rl[1], logging.LevelInfo, "finished call", forwarded)
}

func TestHttpInterceptor_WithHttpFieldsFromContext(t *testing.T) {
	// fields from context override the existing fields
	rl, _ := serveHttp(t, http.StatusOK, nil, grpc_.WithHttpFieldsFromContext(func(ctx context.Context) logging.Fields {
		return logging.Fields{"http.method", "OVERRIDDEN", "tenant", "t1"}
	}))
	want := map[string]string{"http.method": "OVERRIDDEN", "tenant": "t1"}
	checkLog(t, rl[0], logging.LevelInfo, "started call", want)
	checkLog(t, rl[1], logging.LevelInfo, "finished call", want)
}

func TestHttpRoundTripDecorator_WithHttpErrorFields(t *testing.T) {
	var rl recordLogger
	wantErr := errors.New("refused")
	rt := grpc_.HttpRoundTripDecorator(&rl,
		grpc_.WithHttpErrorFields(func(err error) logging.Fields { return logging.Fields{"error.kind", "dial"} }),
		grpc_.WithHttpCodes(func(status int, err error) codes.Code { return codes.Internal }),
	).WrapRoundTrip(http_.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, wantErr
	}))
	r, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := rt.RoundTrip(r); err != wantErr {
		t.Fatalf("RoundTrip() error = %v, want %v", err, wantErr)
	}
	// logging.DefaultClientCodeToLevel maps Internal to Warn
	checkLog(t, rl[1], logging.LevelWarn, "finished call", map[string]string{"http.error": "refused", "error.kind": "dial"})
}

func TestExtractHttpLoggingOptions(t *testing.T) {
	opts := grpc_.ExtractHttpLoggingOptions(
		grpc_.WithHttpLoggingOption(grpc_.WithHttpLogOnEvents(logging.FinishCall)),
		grpc_.WithHttpLoggingOption(grpc_.WithHttpDurationField(grpc_.HttpDurationToDurationField)),
	)
	rl, _ := serveHttp(t, http.StatusOK, nil, opts...)
	if len(rl) != 1 {
		t.Fatalf("logged %d, want 1", len(rl))
	}
	checkLog(t, rl[0], logging.LevelInfo, "finished call", map[string]string{"http.duration": anyValue})
}
