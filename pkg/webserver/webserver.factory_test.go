// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/searKing/golang/pkg/webserver"
	httptrace_ "github.com/searKing/golang/pkg/webserver/pkg/httptrace"
)

// processLogger is the slog logger at process start, before tests replace it.
// OtelHandling must not deadlock when wrapping this builtin handler.
var processLogger = slog.Default()

func TestOtelHandlingLogDoesNotDeadlock(t *testing.T) {
	prev := slog.Default()
	slog.SetDefault(processLogger)
	t.Cleanup(func() { slog.SetDefault(prev) })

	assertOtelPrepareRun(t)
}

func TestOtelHandlingWrappedBuiltinHandlerDoesNotDeadlock(t *testing.T) {
	prev := slog.Default()
	slog.SetDefault(processLogger)
	// Outer type is not *slog.defaultHandler, but Handle delegates to it.
	slog.SetDefault(slog.New(delegatingSlogHandler{next: slog.Default().Handler()}))
	t.Cleanup(func() { slog.SetDefault(prev) })

	assertOtelPrepareRun(t)
}

func TestOtelHandlingUsesCustomSlogHandler(t *testing.T) {
	prev := slog.Default()
	t.Cleanup(func() { slog.SetDefault(prev) })

	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})))

	srv, err := webserver.NewWebServer(webserver.FactoryConfig{
		BindAddress:  "127.0.0.1:0",
		OtelHandling: true,
	})
	if err != nil {
		t.Fatalf("create web server: %v", err)
	}
	if _, err := srv.PrepareRun(); err != nil {
		t.Fatalf("prepare web server: %v", err)
	}
	slog.Info("custom-slog-marker")
	log.Printf("custom-log-printf-marker")
	if !strings.Contains(buf.String(), "custom-slog-marker") {
		t.Fatalf("slog no longer uses the custom handler, logs=%s", buf.String())
	}
	if strings.Contains(buf.String(), "custom-log-printf-marker") {
		t.Fatalf("log.Printf re-entered the slog handler, logs=%s", buf.String())
	}
}

// delegatingSlogHandler forwards to the builtin handler under another type.
type delegatingSlogHandler struct {
	next slog.Handler
}

func (d delegatingSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return d.next.Enabled(ctx, level)
}

func (d delegatingSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	return d.next.Handle(ctx, r)
}

func (d delegatingSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return delegatingSlogHandler{next: d.next.WithAttrs(attrs)}
}

func (d delegatingSlogHandler) WithGroup(name string) slog.Handler {
	return delegatingSlogHandler{next: d.next.WithGroup(name)}
}

func assertOtelPrepareRun(t *testing.T) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		srv, err := webserver.NewWebServer(webserver.FactoryConfig{
			BindAddress:  "127.0.0.1:0",
			OtelHandling: true,
		})
		if err != nil {
			t.Errorf("create web server: %v", err)
			return
		}
		if _, err := srv.PrepareRun(); err != nil {
			t.Errorf("prepare web server: %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("deadlocked while logging with OtelHandling")
	}
}

func TestHTTPTraceLoggingSwitch(t *testing.T) {
	off, err := webserver.NewWebServer(webserver.FactoryConfig{BindAddress: "127.0.0.1:0"})
	if err != nil {
		t.Fatalf("create web server: %v", err)
	}
	on, err := webserver.NewWebServer(webserver.FactoryConfig{
		BindAddress:      "127.0.0.1:0",
		HTTPTraceLogging: true,
	})
	if err != nil {
		t.Fatalf("create web server with HTTPTraceLogging: %v", err)
	}

	got := len(on.HttpRoundTripDecorators())
	want := len(off.HttpRoundTripDecorators()) + 1
	if got != want {
		t.Fatalf("HTTPTraceLogging decorators = %d, want %d", got, want)
	}

	// Options alone do not install the decorator.
	optsOnly, err := webserver.NewWebServer(webserver.FactoryConfig{
		BindAddress: "127.0.0.1:0",
		HTTPTraceOptions: []httptrace_.Option{
			httptrace_.WithEvents(httptrace_.EventAll),
		},
	})
	if err != nil {
		t.Fatalf("create web server with HTTPTraceOptions only: %v", err)
	}
	if len(optsOnly.HttpRoundTripDecorators()) != len(off.HttpRoundTripDecorators()) {
		t.Fatal("HTTPTraceOptions installed a decorator while HTTPTraceLogging is false")
	}
}

func TestHTTPTraceOptionsSelectEvents(t *testing.T) {
	logs := traceLogs(t, httptrace_.WithEvents(httptrace_.EventGetConn))
	if !strings.Contains(logs, "getting connection") {
		t.Fatalf("WithEvents(EventGetConn) missing trace log, logs=%s", logs)
	}

	logs = traceLogs(t, httptrace_.WithEvents(0))
	if strings.Contains(logs, "getting connection") {
		t.Fatalf("WithEvents(0) logged trace events, logs=%s", logs)
	}
}

func TestHttpRoundTripDecoratorsClient(t *testing.T) {
	t.Run("with trace logging", func(t *testing.T) {
		status, logs := roundTripThroughDecorators(t, webserver.FactoryConfig{
			BindAddress:      "127.0.0.1:0",
			HTTPTraceLogging: true,
		})
		if status != http.StatusOK {
			t.Fatalf("status = %d, want 200", status)
		}
		for _, want := range []string{
			"http request sending",
			"finished http call with status code 200",
			"getting connection",
			"got connection",
		} {
			if !strings.Contains(logs, want) {
				t.Errorf("client log missing %q, logs=%s", want, logs)
			}
		}
	})

	t.Run("without trace logging", func(t *testing.T) {
		status, logs := roundTripThroughDecorators(t, webserver.FactoryConfig{
			BindAddress: "127.0.0.1:0",
		})
		if status != http.StatusOK {
			t.Fatalf("status = %d, want 200", status)
		}
		if !strings.Contains(logs, "http request sending") {
			t.Fatalf("client log missing request log, logs=%s", logs)
		}
		if strings.Contains(logs, "getting connection") {
			t.Fatalf("trace log present while HTTPTraceLogging is false, logs=%s", logs)
		}
	})
}

// traceLogs builds a server with HTTPTraceLogging and opts, issues one request
// through its round-trip decorators, and returns the captured slog output.
func traceLogs(t *testing.T, opts ...httptrace_.Option) string {
	t.Helper()
	_, logs := roundTripThroughDecorators(t, webserver.FactoryConfig{
		BindAddress:      "127.0.0.1:0",
		HTTPTraceLogging: true,
		HTTPTraceOptions: opts,
	})
	return logs
}

// roundTripThroughDecorators builds an http.Client from srv.HttpRoundTripDecorators
// and GETs an httptest server. It returns the status code and captured slog output.
func roundTripThroughDecorators(t *testing.T, fc webserver.FactoryConfig) (int, string) {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(prev)

	srv, err := webserver.NewWebServer(fc)
	if err != nil {
		t.Fatalf("create web server: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	base := ts.Client().Transport.(*http.Transport).Clone()
	client := &http.Client{Transport: srv.HttpRoundTripDecorators().WrapRoundTrip(base)}
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("GET %s: %v", ts.URL, err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return resp.StatusCode, buf.String()
}
