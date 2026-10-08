// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

// headerLogger records header groups logged, keyed by group name.
type headerLogger map[string]map[string][]string

func (hl headerLogger) Log(_ context.Context, _ logging.Level, _ string, fields ...any) {
	for _, f := range fields {
		a, ok := f.(slog.Attr)
		if !ok || a.Value.Kind() != slog.KindGroup {
			continue
		}
		h := map[string][]string{}
		for _, ga := range a.Value.Group() {
			h[ga.Key], _ = ga.Value.Any().([]string)
		}
		hl[a.Key] = h
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
