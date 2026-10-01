// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptrace_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	httptrace_ "github.com/searKing/golang/pkg/webserver/pkg/httptrace"
)

// newJSONLogger returns a slog.Logger writing JSON lines into buf, honoring
// the given level threshold.
func newJSONLogger(buf *bytes.Buffer, level slog.Level) *slog.Logger {
	h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: level})
	return slog.New(h)
}

// useDefaultLogger installs logger as slog.Default for the rest of the test.
func useDefaultLogger(t *testing.T, logger *slog.Logger) {
	t.Helper()
	prev := slog.Default()
	slog.SetDefault(logger)
	t.Cleanup(func() { slog.SetDefault(prev) })
}

// collectMessages parses JSON log lines in buf and returns their "msg" fields.
func collectMessages(t *testing.T, buf *bytes.Buffer) []string {
	t.Helper()
	var msgs []string
	for line := range strings.SplitSeq(strings.TrimRight(buf.String(), "\n"), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("invalid json log line %q: %v", line, err)
		}
		if msg, ok := m["msg"].(string); ok {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

func hasMessage(msgs []string, want string) bool { return slices.Contains(msgs, want) }

// assertEventAttr checks that each logged message in want carries the expected
// event attribute. Messages that did not occur are ignored; callers assert
// presence separately.
func assertEventAttr(t *testing.T, buf *bytes.Buffer, want map[string]string) {
	t.Helper()
	seen := map[string]string{}
	for line := range strings.SplitSeq(strings.TrimRight(buf.String(), "\n"), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("invalid json log line %q: %v", line, err)
		}
		msg, _ := m["msg"].(string)
		event, _ := m["event"].(string)
		seen[msg] = event
	}
	for msg, event := range want {
		got, ok := seen[msg]
		if !ok {
			continue
		}
		if got != event {
			t.Errorf("message %q: event = %q, want %q", msg, got, event)
		}
	}
}

// doRequest issues one GET against ts using a client whose Transport is wrapped
// by the given decorator options.
func doRequest(t *testing.T, ts *httptest.Server, opts ...httptrace_.Option) {
	t.Helper()
	dec := httptrace_.LoggingDecorator(opts...)
	// Clone a transport that trusts the httptest TLS cert (if any).
	base := ts.Client().Transport.(*http.Transport).Clone()
	client := &http.Client{Transport: dec.WrapRoundTrip(base)}
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("GET %s: %v", ts.URL, err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

func TestLoggingDecorator_HTTPEmitsCoreEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))
	defer ts.Close()

	var buf bytes.Buffer
	useDefaultLogger(t, newJSONLogger(&buf, slog.LevelDebug))

	doRequest(t, ts)

	msgs := collectMessages(t, &buf)
	// EventDefault logs connection acquisition. DNS may not run for httptest.
	for _, want := range []string{
		"getting connection",
		"got connection",
	} {
		if !hasMessage(msgs, want) {
			t.Errorf("expected message %q, got %v", want, msgs)
		}
	}
	for _, unexpected := range []string{
		"wrote request",
		"got first response byte",
		"connecting",
		"returned idle connection",
	} {
		if hasMessage(msgs, unexpected) {
			t.Errorf("EventDefault logged %q, got %v", unexpected, msgs)
		}
	}
	assertEventAttr(t, &buf, map[string]string{
		"getting connection":  "GetConn",
		"got connection":      "GotConn",
		"finished dns lookup": "DNSDone",
	})
}

func TestLoggingDecorator_HTTPSEmitsTLSEvents(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))
	defer ts.Close()

	var buf bytes.Buffer
	useDefaultLogger(t, newJSONLogger(&buf, slog.LevelDebug))

	doRequest(t, ts, httptrace_.WithEvents(httptrace_.EventAll))

	msgs := collectMessages(t, &buf)
	for _, want := range []string{
		"starting tls handshake",
		"finished tls handshake",
		"connecting",
		"finished connecting",
	} {
		if !hasMessage(msgs, want) {
			t.Errorf("expected message %q, got %v", want, msgs)
		}
	}
	assertEventAttr(t, &buf, map[string]string{
		"starting tls handshake": "TLSHandshakeStart",
		"finished tls handshake": "TLSHandshakeDone",
		"connecting":             "ConnectStart",
		"finished connecting":    "ConnectDone",
	})
}

func TestLoggingDecorator_LevelFiltering(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer ts.Close()

	var buf bytes.Buffer
	// Handler threshold is Warn, but we emit at Info -> everything must be dropped.
	useDefaultLogger(t, newJSONLogger(&buf, slog.LevelWarn))

	doRequest(t, ts, httptrace_.WithLevel(slog.LevelInfo))

	if buf.Len() != 0 {
		t.Errorf("expected no log output when level is filtered, got: %s", buf.String())
	}
}

func TestLoggingDecorator_WithEventsOnlySelected(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer ts.Close()

	var buf bytes.Buffer
	useDefaultLogger(t, newJSONLogger(&buf, slog.LevelDebug))

	// Only enable DNS events; everything else must stay silent.
	doRequest(t, ts, httptrace_.WithEvents(httptrace_.EventDNS))

	msgs := collectMessages(t, &buf)
	for _, m := range msgs {
		if m != "resolving host" && m != "finished dns lookup" {
			t.Errorf("unexpected message %q with WithEvents(EventDNS); all=%v", m, msgs)
		}
	}
	// httptest uses 127.0.0.1 directly, so DNS callbacks may not fire. We
	// cannot assert they are present, only that no other events leaked.
}

func TestLoggingDecorator_WithEventsZeroDisablesAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer ts.Close()

	var buf bytes.Buffer
	useDefaultLogger(t, newJSONLogger(&buf, slog.LevelDebug))

	doRequest(t, ts, httptrace_.WithEvents(0))

	if buf.Len() != 0 {
		t.Errorf("expected no log output with WithEvents(0), got: %s", buf.String())
	}
}

func TestLoggingDecorator_ComposesWithExistingTrace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer ts.Close()

	var buf bytes.Buffer
	useDefaultLogger(t, newJSONLogger(&buf, slog.LevelDebug))

	var upstreamGetConn atomic.Int32
	dec := httptrace_.LoggingDecorator()
	base := ts.Client().Transport.(*http.Transport).Clone()
	client := &http.Client{Transport: dec.WrapRoundTrip(base)}

	// Attach an upstream ClientTrace BEFORE the decorator sees the request.
	upstream := &httptrace.ClientTrace{
		GetConn: func(string) { upstreamGetConn.Add(1) },
	}
	ctx := httptrace.WithClientTrace(context.Background(), upstream)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	if got := upstreamGetConn.Load(); got == 0 {
		t.Error("expected upstream GetConn callback to still fire after decorator wrapping, got 0")
	}
	if !hasMessage(collectMessages(t, &buf), "getting connection") {
		t.Error("expected decorator's own GetConn log to be present")
	}
}

func TestLoggingDecorator_KeepAliveReusedConn(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer ts.Close()

	var buf bytes.Buffer
	useDefaultLogger(t, newJSONLogger(&buf, slog.LevelDebug))

	dec := httptrace_.LoggingDecorator(httptrace_.WithEvents(httptrace_.EventGotConn))
	base := ts.Client().Transport.(*http.Transport).Clone()
	client := &http.Client{Transport: dec.WrapRoundTrip(base)}

	for i := range 2 {
		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatalf("GET #%d: %v", i, err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}

	// Parse all lines and look for a got-conn record with reused=true.
	var sawReused bool
	for line := range strings.SplitSeq(strings.TrimRight(buf.String(), "\n"), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("invalid json: %v", err)
		}
		if m["msg"] != "got connection" {
			continue
		}
		if r, ok := m["reused"].(bool); ok && r {
			sawReused = true
			break
		}
	}
	if !sawReused {
		t.Errorf("expected a 'got conn' event with reused=true across two sequential requests, logs=%s", buf.String())
	}
}

// fakeAddr implements net.Addr for the nil-Conn test, but we actually test the
// nil-Conn branch by invoking the ClientTrace callback directly via a custom
// RoundTripper rather than relying on real transports.
type recordingHandler struct {
	mu      sync.Mutex
	records []map[string]any
}

func (h *recordingHandler) Enabled(context.Context, slog.Level) bool { return true }
func (h *recordingHandler) WithAttrs([]slog.Attr) slog.Handler       { return h }
func (h *recordingHandler) WithGroup(string) slog.Handler            { return h }
func (h *recordingHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	entry := map[string]any{"msg": r.Message}
	r.Attrs(func(a slog.Attr) bool {
		entry[a.Key] = a.Value.Any()
		return true
	})
	h.records = append(h.records, entry)
	return nil
}

// nilConnRoundTripper synthesises a GotConn callback with a nil Conn, which
// some fake transports may do. The decorator must not panic.
type nilConnRoundTripper struct{}

func (nilConnRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if trace := httptrace.ContextClientTrace(req.Context()); trace != nil && trace.GotConn != nil {
		trace.GotConn(httptrace.GotConnInfo{Conn: nil, Reused: true, WasIdle: true})
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    req,
		Header:     make(http.Header),
	}, nil
}

func TestLoggingDecorator_GotConnNilConnDoesNotPanic(t *testing.T) {
	rec := &recordingHandler{}
	useDefaultLogger(t, slog.New(rec))

	dec := httptrace_.LoggingDecorator()
	rt := dec.WrapRoundTrip(nilConnRoundTripper{})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.invalid/", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic on nil Conn: %v", r)
		}
	}()
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip: %v", err)
	}
	_ = resp.Body.Close()

	// Verify a got-conn record was emitted and did not try to log local/remote.
	rec.mu.Lock()
	defer rec.mu.Unlock()
	var sawGotConn bool
	for _, e := range rec.records {
		if e["msg"] == "got connection" {
			sawGotConn = true
			if _, hasLocal := e["local_addr"]; hasLocal {
				t.Errorf("expected no 'local_addr' attr on nil-Conn GotConn, got %v", e)
			}
			if _, hasRemote := e["remote_addr"]; hasRemote {
				t.Errorf("expected no 'remote_addr' attr on nil-Conn GotConn, got %v", e)
			}
			if e["event"] != "GotConn" {
				t.Errorf("event = %v, want GotConn", e["event"])
			}
		}
	}
	if !sawGotConn {
		t.Errorf("expected a 'got conn' record, got %v", rec.records)
	}
}

// failingTraceRoundTripper invokes trace callbacks with errors so tests can
// assert failure messages and the "error" attribute without a real transport.
type failingTraceRoundTripper struct{}

func (failingTraceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	trace := httptrace.ContextClientTrace(req.Context())
	if trace == nil {
		return nil, errors.New("missing client trace")
	}
	if trace.DNSDone != nil {
		trace.DNSDone(httptrace.DNSDoneInfo{Err: errors.New("dns failed")})
	}
	if trace.ConnectDone != nil {
		trace.ConnectDone("tcp", "127.0.0.1:443", errors.New("connect failed"))
	}
	if trace.TLSHandshakeDone != nil {
		trace.TLSHandshakeDone(tls.ConnectionState{}, errors.New("tls failed"))
	}
	if trace.WroteRequest != nil {
		trace.WroteRequest(httptrace.WroteRequestInfo{Err: errors.New("write failed")})
	}
	if trace.PutIdleConn != nil {
		trace.PutIdleConn(errors.New("idle failed"))
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    req,
		Header:     make(http.Header),
	}, nil
}

func TestLoggingDecorator_FailureMessagesUseErrorAttr(t *testing.T) {
	rec := &recordingHandler{}
	useDefaultLogger(t, slog.New(rec))

	dec := httptrace_.LoggingDecorator(httptrace_.WithEvents(httptrace_.EventAll))
	rt := dec.WrapRoundTrip(failingTraceRoundTripper{})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.invalid/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip: %v", err)
	}
	_ = resp.Body.Close()

	want := map[string]struct {
		event   string
		errText string
	}{
		"finished dns lookup":              {"DNSDone", "dns failed"},
		"finished connecting":              {"ConnectDone", "connect failed"},
		"finished tls handshake":           {"TLSHandshakeDone", "tls failed"},
		"failed to write request":          {"WroteRequest", "write failed"},
		"failed to return idle connection": {"PutIdleConn", "idle failed"},
	}

	rec.mu.Lock()
	defer rec.mu.Unlock()
	got := map[string]string{}
	for _, e := range rec.records {
		msg, _ := e["msg"].(string)
		w, ok := want[msg]
		if !ok {
			continue
		}
		if e["event"] != w.event {
			t.Errorf("message %q: event = %v, want %s", msg, e["event"], w.event)
		}
		if _, hasErr := e["err"]; hasErr {
			t.Errorf("message %q logged error under key \"err\", record=%v", msg, e)
		}
		errVal, ok := e["error"].(error)
		if !ok {
			t.Errorf("message %q missing error attr, record=%v", msg, e)
			continue
		}
		got[msg] = errVal.Error()
	}
	for msg, w := range want {
		if got[msg] != w.errText {
			t.Errorf("message %q: error = %q, want %q", msg, got[msg], w.errText)
		}
	}
}

// Compile-time check that TLS state value is accepted (defensive against future
// signature changes in net/http/httptrace).
var _ = func() tls.ConnectionState { return tls.ConnectionState{} }
