// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otel_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	otel_ "github.com/searKing/golang/pkg/webserver/pkg/otel"
)

func TestGrpcOrHttpHandlerDecorator(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	prevProvider := otel.GetTracerProvider()
	prevPropagator := otel.GetTextMapPropagator()
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() {
		otel.SetTracerProvider(prevProvider)
		otel.SetTextMapPropagator(prevPropagator)
	})

	parentCtx, parent := provider.Tracer("test").Start(t.Context(), "caller")
	parent.End()

	req := httptest.NewRequest(http.MethodPost, "/pkg.Service/Method", nil)
	otel.GetTextMapPropagator().Inject(parentCtx, propagation.HeaderCarrier(req.Header))
	rec := httptest.NewRecorder()

	var gotWriter http.ResponseWriter
	var gotSpan trace.SpanContext
	otel_.GrpcOrHttpHandlerDecorator().WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotWriter = w
		gotSpan = trace.SpanFromContext(r.Context()).SpanContext()
		w.WriteHeader(http.StatusNotFound)
	})).ServeHTTP(rec, req)

	if gotWriter != rec {
		t.Fatal("decorator wrapped ResponseWriter")
	}
	if gotSpan.TraceID() != parent.SpanContext().TraceID() {
		t.Fatalf("trace id = %s, want %s", gotSpan.TraceID(), parent.SpanContext().TraceID())
	}
	var server sdktrace.ReadOnlySpan
	for _, span := range recorder.Ended() {
		if span.SpanContext().SpanID() == gotSpan.SpanID() {
			server = span
		}
	}
	if server == nil {
		t.Fatal("server span was not recorded")
	}
	if server.Parent().SpanID() != parent.SpanContext().SpanID() {
		t.Fatalf("parent span id = %s, want %s", server.Parent().SpanID(), parent.SpanContext().SpanID())
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	injected := propagation.TraceContext{}.Extract(t.Context(), propagation.HeaderCarrier(req.Header))
	if sc := trace.SpanContextFromContext(injected); sc.SpanID() != gotSpan.SpanID() {
		t.Fatalf("injected traceparent span = %s, want %s", sc.SpanID(), gotSpan.SpanID())
	}
}

func TestGrpcOrHttpHandlerDecoratorStartsRootWithoutTraceparent(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	prevProvider := otel.GetTracerProvider()
	prevPropagator := otel.GetTextMapPropagator()
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() {
		otel.SetTracerProvider(prevProvider)
		otel.SetTextMapPropagator(prevPropagator)
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	var got trace.SpanContext
	var gotWriter http.ResponseWriter
	otel_.GrpcOrHttpHandlerDecorator().WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = trace.SpanFromContext(r.Context()).SpanContext()
		gotWriter = w
		w.WriteHeader(http.StatusInternalServerError)
		w.(http.Flusher).Flush()
	})).ServeHTTP(rec, req)

	if !got.IsValid() {
		t.Fatal("downstream context has no span")
	}
	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("ended spans = %d, want 1", len(spans))
	}
	server := spans[0]
	if server.Parent().IsValid() {
		t.Fatal("span has a parent")
	}
	if server.SpanContext().SpanID() != got.SpanID() {
		t.Fatal("downstream span is not the root span")
	}
	if server.SpanKind() != trace.SpanKindServer {
		t.Fatalf("span kind = %v, want server", server.SpanKind())
	}
	if server.Name() != "GET /healthz" {
		t.Fatalf("span name = %q", server.Name())
	}
	if gotWriter != rec {
		t.Fatal("decorator wrapped ResponseWriter")
	}
	if rec.Code != http.StatusInternalServerError || !rec.Flushed {
		t.Fatalf("response status = %d flushed = %v, want 500 flushed", rec.Code, rec.Flushed)
	}
	injected := propagation.TraceContext{}.Extract(t.Context(), propagation.HeaderCarrier(req.Header))
	if sc := trace.SpanContextFromContext(injected); sc.SpanID() != got.SpanID() {
		t.Fatalf("injected traceparent span = %s, want %s", sc.SpanID(), got.SpanID())
	}
}

func TestGrpcOrHttpHandlerDecoratorUsesPropagatorsOption(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	prevProvider := otel.GetTracerProvider()
	prevPropagator := otel.GetTextMapPropagator()
	otel.SetTracerProvider(provider)
	// Global propagator does not understand X-Trace-Id.
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() {
		otel.SetTracerProvider(prevProvider)
		otel.SetTextMapPropagator(prevPropagator)
	})

	const traceID = "0123456789abcdef0123456789abcdef"
	req := httptest.NewRequest(http.MethodPost, "/pkg.Service/Method", nil)
	req.Header.Set("X-Trace-Id", traceID)
	var got trace.TraceID
	otel_.GrpcOrHttpHandlerDecorator(otelhttp.WithPropagators(customTracePropagator{})).
		WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got = trace.SpanFromContext(r.Context()).SpanContext().TraceID()
		})).ServeHTTP(httptest.NewRecorder(), req)

	want, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("trace id = %s, want %s", got, want)
	}
	injected := propagation.TraceContext{}.Extract(t.Context(), propagation.HeaderCarrier(req.Header))
	if sc := trace.SpanContextFromContext(injected); sc.TraceID() != want {
		t.Fatalf("injected traceparent trace id = %s, want %s", sc.TraceID(), want)
	}
}

type customTracePropagator struct{}

func (customTracePropagator) Fields() []string { return []string{"X-Trace-Id"} }

func (customTracePropagator) Inject(context.Context, propagation.TextMapCarrier) {}

func (customTracePropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	tid, err := trace.TraceIDFromHex(carrier.Get("X-Trace-Id"))
	if err != nil || !tid.IsValid() {
		return ctx
	}
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     trace.SpanID{1},
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})
	return trace.ContextWithRemoteSpanContext(ctx, sc)
}
