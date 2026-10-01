// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package httptrace provides an http.RoundTripper decorator that logs
// client-side HTTP trace events (DNS, connection, TLS, TTFB, ...) via slog.
//
// It is implemented on top of net/http/httptrace and composes with any
// ClientTrace already attached to the request context: previously registered
// callbacks are preserved and will still be invoked.
package httptrace

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/http/httptrace"

	http_ "github.com/searKing/golang/go/net/http"
)

// options controls the behavior of the logging decorator.
type options struct {
	level  slog.Level
	events Event
}

// Option configures the logging decorator.
type Option func(*options)

// WithLevel sets the slog.Level used to emit trace events.
// If unset, slog.LevelDebug is used because trace events can be very noisy
// (GetConn/GotConn fire on every request, including keep-alive reuses).
func WithLevel(l slog.Level) Option {
	return func(o *options) { o.level = l }
}

// WithEvents selects which httptrace events are logged. Events can be
// OR-combined, e.g. WithEvents(EventDNS | EventTLS). If unset, EventDefault
// is used. Passing 0 disables all logging callbacks.
func WithEvents(e Event) Option {
	return func(o *options) { o.events = e }
}

func newOptions(opts ...Option) *options {
	o := &options{
		level:  slog.LevelDebug,
		events: EventDefault,
	}
	for _, f := range opts {
		if f != nil {
			f(o)
		}
	}
	return o
}

// newLoggingClientTrace builds a *httptrace.ClientTrace whose callbacks emit
// structured logs through the provided options.
//
// Note: the callbacks are bound to the request context captured at the time
// the trace is created; later ctx mutations by downstream code are not
// visible to the callbacks (net/http/httptrace does not pass a context to
// its callbacks, so this is the best we can do).
func newLoggingClientTrace(ctx context.Context, o *options) *httptrace.ClientTrace {
	log := func(event Event, msg string, attrs ...slog.Attr) {
		attrs = append([]slog.Attr{slog.String("event", event.String())}, attrs...)
		slog.LogAttrs(ctx, o.level, msg, attrs...)
	}
	t := &httptrace.ClientTrace{}
	if o.events.has(EventGetConn) {
		t.GetConn = func(hostPort string) {
			log(EventGetConn, "getting connection", slog.String("host_port", hostPort))
		}
	}
	if o.events.has(EventGotConn) {
		t.GotConn = func(info httptrace.GotConnInfo) {
			// Defensive: some custom RoundTrippers may pass a nil Conn.
			if info.Conn == nil {
				log(EventGotConn, "got connection",
					slog.Bool("reused", info.Reused),
					slog.Bool("was_idle", info.WasIdle),
					slog.Duration("idle_time", info.IdleTime),
				)
				return
			}
			log(EventGotConn, "got connection",
				slog.String("local_addr", info.Conn.LocalAddr().String()),
				slog.String("remote_addr", info.Conn.RemoteAddr().String()),
				slog.Bool("reused", info.Reused),
				slog.Bool("was_idle", info.WasIdle),
				slog.Duration("idle_time", info.IdleTime),
			)
		}
	}
	if o.events.has(EventDNSStart) {
		t.DNSStart = func(info httptrace.DNSStartInfo) {
			log(EventDNSStart, "resolving host", slog.String("host", info.Host))
		}
	}
	if o.events.has(EventDNSDone) {
		t.DNSDone = func(info httptrace.DNSDoneInfo) {
			addrs := make([]string, 0, len(info.Addrs))
			for _, a := range info.Addrs {
				addrs = append(addrs, a.String())
			}
			attrs := []slog.Attr{
				slog.Any("addrs", addrs),
				slog.Bool("coalesced", info.Coalesced),
			}
			if info.Err != nil {
				attrs = append(attrs, slog.Any("error", info.Err))
			}
			log(EventDNSDone, "finished dns lookup", attrs...)
		}
	}
	if o.events.has(EventConnectStart) {
		t.ConnectStart = func(network, addr string) {
			log(EventConnectStart, "connecting",
				slog.String("network", network),
				slog.String("addr", addr),
			)
		}
	}
	if o.events.has(EventConnectDone) {
		t.ConnectDone = func(network, addr string, err error) {
			attrs := []slog.Attr{
				slog.String("network", network),
				slog.String("addr", addr),
			}
			if err != nil {
				attrs = append(attrs, slog.Any("error", err))
			}
			log(EventConnectDone, "finished connecting", attrs...)
		}
	}
	if o.events.has(EventTLSHandshakeStart) {
		t.TLSHandshakeStart = func() {
			log(EventTLSHandshakeStart, "starting tls handshake")
		}
	}
	if o.events.has(EventTLSHandshakeDone) {
		t.TLSHandshakeDone = func(state tls.ConnectionState, err error) {
			attrs := []slog.Attr{
				slog.String("server_name", state.ServerName),
				slog.String("negotiated_protocol", state.NegotiatedProtocol),
				slog.Bool("did_resume", state.DidResume),
			}
			if err != nil {
				attrs = append(attrs, slog.Any("error", err))
			}
			log(EventTLSHandshakeDone, "finished tls handshake", attrs...)
		}
	}
	if o.events.has(EventWroteRequest) {
		t.WroteRequest = func(info httptrace.WroteRequestInfo) {
			if info.Err != nil {
				log(EventWroteRequest, "failed to write request", slog.Any("error", info.Err))
				return
			}
			log(EventWroteRequest, "wrote request")
		}
	}
	if o.events.has(EventGotFirstResponseByte) {
		t.GotFirstResponseByte = func() {
			log(EventGotFirstResponseByte, "got first response byte")
		}
	}
	if o.events.has(EventPutIdleConn) {
		t.PutIdleConn = func(err error) {
			if err != nil {
				log(EventPutIdleConn, "failed to return idle connection", slog.Any("error", err))
				return
			}
			log(EventPutIdleConn, "returned idle connection")
		}
	}
	return t
}

// LoggingDecorator returns a RoundTripDecorator that logs httptrace events
// via slog. It composes with any ClientTrace already attached to the request
// context; existing callbacks are preserved and will still be invoked.
func LoggingDecorator(opts ...Option) http_.RoundTripDecorator {
	o := newOptions(opts...)
	return http_.RoundTripDecoratorFunc(func(rt http.RoundTripper) http.RoundTripper {
		return http_.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()
			ctx = httptrace.WithClientTrace(ctx, newLoggingClientTrace(ctx, o))
			return rt.RoundTrip(req.WithContext(ctx))
		})
	})
}
