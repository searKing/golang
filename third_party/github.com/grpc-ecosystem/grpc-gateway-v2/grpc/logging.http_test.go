// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	http_ "github.com/searKing/golang/go/net/http"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

// headerLogger records header groups logged, keyed by group name.
type headerLogger map[string]map[string][]string

func (hl headerLogger) Log(_ context.Context, _ logging.Level, _ string, fields ...any) {
	i := logging.Fields(fields).Iterator()
	for i.Next() {
		k, v := i.At()
		gv, ok := v.(slog.Value)
		if !ok || gv.Kind() != slog.KindGroup {
			continue
		}
		h := map[string][]string{}
		for _, ga := range gv.Group() {
			h[ga.Key], _ = ga.Value.Any().([]string)
		}
		hl[k] = h
	}
}

func checkHeader(t *testing.T, got map[string][]string, want map[string][]string) {
	t.Helper()
	for k, wv := range want {
		if gv := got[k]; !slices.Equal(gv, wv) {
			t.Errorf("header %s = %q, want %q", k, gv, wv)
		}
	}
}

func newSensitiveRequest(url string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, url, nil)
	r.Header.Set("Authorization", "TC3-HMAC-SHA256 Credential=AKID/2026-10-08/cvm/tc3_request, Signature=abc")
	r.Header.Add("Cookie", "sid=1")
	r.Header.Add("Cookie", "uid=2")
	r.Header.Set("Cookie2", "$Version=1")
	r.Header.Set("Proxy-Authorization", "basic")
	r.Header.Set("X-TC-Token", "token")
	r.Header.Set("X-Request-Id", "req-1")
	return r
}

// sensitive headers as net/http strips on redirect are redacted, and others are kept
var wantRequestHeader = map[string][]string{
	"Authorization":       {"TC3-HMAC-SHA256 xxxxx"},
	"Cookie":              {"xxxxx", "xxxxx"},
	"Cookie2":             {"xxxxx"},
	"Proxy-Authorization": {"xxxxx"},
	"X-Tc-Token":          {"token"},
	"X-Request-Id":        {"req-1"},
}

func setSensitiveResponseHeader(h http.Header) {
	h.Set("Set-Cookie", "sid=1")
	h.Set("Www-Authenticate", `Bearer realm="api", error="invalid_token"`)
	h.Set("X-Trace-Id", "trace-1")
}

var wantResponseHeader = map[string][]string{
	"Set-Cookie":       {"xxxxx"},
	"Www-Authenticate": {"Bearer xxxxx"},
	"X-Trace-Id":       {"trace-1"},
}

func TestHttpInterceptor_RedactHeader(t *testing.T) {
	t.Setenv("HTTP_GO_LOG_HTTP_HEADER", "true")
	hl := headerLogger{}
	h := grpc_.HttpInterceptor(hl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setSensitiveResponseHeader(w.Header())
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(httptest.NewRecorder(), newSensitiveRequest("/"))

	checkHeader(t, hl["http.request.header"], wantRequestHeader)
	checkHeader(t, hl["http.response.header"], wantResponseHeader)
	if _, ok := hl["http.response.header"]["X-Request-Id"]; ok {
		t.Errorf("response header logged request header X-Request-Id")
	}
}

func TestHttpRoundTripDecorator_RedactHeader(t *testing.T) {
	t.Setenv("HTTP_GO_LOG_HTTP_HEADER", "true")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setSensitiveResponseHeader(w.Header())
	}))
	defer srv.Close()

	hl := headerLogger{}
	rt := grpc_.HttpRoundTripDecorator(hl).WrapRoundTrip(http.DefaultTransport)
	r := newSensitiveRequest(srv.URL)
	r.RequestURI = ""
	resp, err := rt.RoundTrip(r)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	checkHeader(t, hl["http.request.header"], wantRequestHeader)
	checkHeader(t, hl["http.response.header"], wantResponseHeader)
}

// logRecord is a log recorded by recordLogger, with fields not grouped, keyed by name.
type logRecord struct {
	level  logging.Level
	msg    string
	fields map[string]string
}

// recordLogger records logs, marking fields logged more than once as "DUPLICATE".
type recordLogger []logRecord

func (rl *recordLogger) Log(_ context.Context, level logging.Level, msg string, fields ...any) {
	m := map[string]string{}
	i := logging.Fields(fields).Iterator()
	for i.Next() {
		k, v := i.At()
		if _, dup := m[k]; dup {
			m[k] = "DUPLICATE"
			continue
		}
		m[k] = fmt.Sprint(v)
	}
	*rl = append(*rl, logRecord{level: level, msg: msg, fields: m})
}

// anyValue wants a field logged with any value.
const anyValue = "<any>"

// checkFields checks fields got, where fields wanted as "" are not logged.
func checkFields(t *testing.T, name string, got map[string]string, want map[string]string) {
	t.Helper()
	for k, wv := range want {
		gv, ok := got[k]
		switch {
		case wv == "":
			if ok {
				t.Errorf("%s: field %s = %q, want not logged", name, k, gv)
			}
		case !ok:
			t.Errorf("%s: field %s not logged, want %q", name, k, wv)
		case wv != anyValue && gv != wv:
			t.Errorf("%s: field %s = %q, want %q", name, k, gv, wv)
		}
	}
}

func checkLog(t *testing.T, rec logRecord, level logging.Level, msg string, want map[string]string) {
	t.Helper()
	if rec.level != level || rec.msg != msg {
		t.Errorf("log = %v %q, want %v %q", rec.level, rec.msg, level, msg)
	}
	checkFields(t, msg, rec.fields, want)
}

func TestHttpInterceptor_Fields(t *testing.T) {
	var rl recordLogger
	var injected map[string]string
	h := grpc_.HttpInterceptor(&rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var hl recordLogger
		hl.Log(r.Context(), logging.LevelInfo, "", logging.ExtractFields(r.Context())...)
		injected = hl[0].fields
		_, _ = w.Write([]byte("ok"))
	}))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r := httptest.NewRequestWithContext(ctx, http.MethodPost, "/a?b=1", strings.NewReader("hello"))
	r.RemoteAddr = "192.0.2.1:1234"
	h.ServeHTTP(httptest.NewRecorder(), r)

	common := map[string]string{
		"protocol":              "http",
		"http.component":        "server",
		"http.method":           "POST",
		"http.request.uri":      "/a?b=1",
		"http.host":             "example.com",
		"http.remote_addr":      "192.0.2.1:1234",
		"http.start_time":       "",
		"http.status_code":      "",
		"http.time_ms":          "",
		"http.client_ip":        "",
		"http.request.deadline": "",
	}
	checkFields(t, "fields injected into the handler context", injected, common)

	if len(rl) != 2 {
		t.Fatalf("logged %d, want 2", len(rl))
	}
	once := map[string]string{"http.start_time": anyValue, "http.request.deadline": anyValue}
	checkLog(t, rl[0], logging.LevelInfo, "started call", map[string]string{
		"protocol": "http", "http.component": "server", "http.method": "POST", "http.request.uri": "/a?b=1",
		"http.host": "example.com", "http.remote_addr": "192.0.2.1:1234",
		"http.start_time": anyValue, "http.request.deadline": anyValue, "http.status_code": "",
	})
	checkLog(t, rl[1], logging.LevelInfo, "finished call", map[string]string{
		"protocol": "http", "http.component": "server", "http.method": "POST", "http.request.uri": "/a?b=1",
		"http.status_code": "OK", "http.request_body_size": "5", "http.response_body_size": "2",
		"http.time_ms": anyValue, "cost": "",
	})
	checkFields(t, "finished call", rl[1].fields, once)
	for _, rec := range rl {
		for k, v := range rec.fields {
			if v == "DUPLICATE" {
				t.Errorf("%s: field %s logged more than once", rec.msg, k)
			}
		}
	}
}

func TestHttpInterceptor_StatusWithoutWrite(t *testing.T) {
	var rl recordLogger
	// net/http replies 200 if the handler writes nothing
	h := grpc_.HttpInterceptor(&rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status replied = %d, want %d", rec.Code, http.StatusOK)
	}
	checkLog(t, rl[1], logging.LevelInfo, "finished call", map[string]string{"http.status_code": "OK"})
}

func TestHttpInterceptor_Level(t *testing.T) {
	tests := []struct {
		status int
		want   logging.Level
	}{
		{http.StatusOK, logging.LevelInfo},
		{http.StatusFound, logging.LevelInfo},
		{http.StatusBadRequest, logging.LevelInfo},
		{http.StatusNotFound, logging.LevelInfo},
		{http.StatusMethodNotAllowed, logging.LevelInfo},
		{http.StatusForbidden, logging.LevelWarn},
		{http.StatusTooManyRequests, logging.LevelWarn},
		{http.StatusServiceUnavailable, logging.LevelWarn},
		{http.StatusGatewayTimeout, logging.LevelWarn},
		{http.StatusInternalServerError, logging.LevelError},
		{http.StatusNotImplemented, logging.LevelError},
		{http.StatusBadGateway, logging.LevelError},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.status), func(t *testing.T) {
			var rl recordLogger
			h := grpc_.HttpInterceptor(&rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))
			h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
			checkLog(t, rl[1], tt.want, "finished call", map[string]string{"http.status_code": http.StatusText(tt.status)})
		})
	}
}
func TestHttpInterceptor_RedactURI(t *testing.T) {
	tests := []struct {
		target string
		want   string
	}{
		{"/a?b=1", "/a?b=1"},
		{"http://user:pass@example.com/a?b=1", "http://user:xxxxx@example.com/a?b=1"},
	}
	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			var rl recordLogger
			h := grpc_.HttpInterceptor(&rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, tt.target, nil))
			for _, rec := range rl {
				checkFields(t, rec.msg, rec.fields, map[string]string{"http.request.uri": tt.want})
			}
		})
	}
}

func TestHttpRoundTripDecorator_Fields(t *testing.T) {
	var rl recordLogger
	var injected map[string]string
	rt := grpc_.HttpRoundTripDecorator(&rl).WrapRoundTrip(http_.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		var hl recordLogger
		hl.Log(r.Context(), logging.LevelInfo, "", logging.ExtractFields(r.Context())...)
		injected = hl[0].fields
		return &http.Response{StatusCode: http.StatusOK, ContentLength: 2, Body: http.NoBody, Request: r}, nil
	}))
	r, err := http.NewRequest(http.MethodPost, "http://user:pass@example.com/a?b=1", strings.NewReader("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := rt.RoundTrip(r); err != nil {
		t.Fatal(err)
	}

	common := map[string]string{
		"protocol":         "http",
		"http.component":   "client",
		"http.method":      "POST",
		"http.request.uri": "http://user:xxxxx@example.com/a?b=1",
		"http.host":        "example.com",
		"http.remote_addr": "",
	}
	checkFields(t, "fields injected into the request context", injected, common)

	if len(rl) != 2 {
		t.Fatalf("logged %d, want 2", len(rl))
	}
	checkLog(t, rl[0], logging.LevelDebug, "started call", common)
	checkFields(t, "started call", rl[0].fields, map[string]string{"http.start_time": anyValue})
	checkLog(t, rl[1], logging.LevelDebug, "finished call", common)
	checkFields(t, "finished call", rl[1].fields, map[string]string{
		"http.status_code": "OK", "http.request_body_size": "5", "http.response_body_size": "2",
		"http.time_ms": anyValue, "http.error": "",
	})
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "i/o timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestHttpRoundTripDecorator_Level(t *testing.T) {
	refused := &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}
	tests := []struct {
		name       string
		status     int
		err        error
		want       logging.Level
		wantFields map[string]string
	}{
		{"200", http.StatusOK, nil, logging.LevelDebug, map[string]string{"http.status_code": "OK"}},
		{"404", http.StatusNotFound, nil, logging.LevelDebug, map[string]string{"http.status_code": "Not Found"}},
		{"401", http.StatusUnauthorized, nil, logging.LevelInfo, map[string]string{"http.status_code": "Unauthorized"}},
		{"500", http.StatusInternalServerError, nil, logging.LevelWarn, map[string]string{"http.status_code": "Internal Server Error"}},
		{"503", http.StatusServiceUnavailable, nil, logging.LevelWarn, map[string]string{"http.status_code": "Service Unavailable"}},
		{"refused", 0, refused, logging.LevelWarn, map[string]string{"http.error": refused.Error(), "http.status_code": ""}},
		{"timeout", 0, &net.OpError{Op: "read", Net: "tcp", Err: timeoutError{}}, logging.LevelInfo, map[string]string{"http.error": anyValue}},
		{"canceled", 0, context.Canceled, logging.LevelDebug, map[string]string{"http.error": "context canceled"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rl recordLogger
			rt := grpc_.HttpRoundTripDecorator(&rl).WrapRoundTrip(http_.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
				if tt.err != nil {
					return nil, tt.err
				}
				return &http.Response{StatusCode: tt.status, Body: http.NoBody, Request: r}, nil
			}))
			r, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := rt.RoundTrip(r); err != tt.err {
				t.Fatalf("RoundTrip() error = %v, want %v", err, tt.err)
			}
			checkLog(t, rl[1], tt.want, "finished call", tt.wantFields)
		})
	}
}
