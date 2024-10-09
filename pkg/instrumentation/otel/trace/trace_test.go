// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr/funcr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"

	trace_ "github.com/searKing/golang/pkg/instrumentation/otel/trace"

	_ "github.com/searKing/golang/pkg/instrumentation/otel/trace/otlptrace/otlptracegrpc" // for otel-grpc
	_ "github.com/searKing/golang/pkg/instrumentation/otel/trace/otlptrace/otlptracehttp" // for otel-http
	_ "github.com/searKing/golang/pkg/instrumentation/otel/trace/stdouttrace"             // for stdout
)

const testLoop = 1
const testInterval = 0 * time.Minute
const testForcePush = false

func TestNewTracerProvider(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	otel.SetLogger(funcr.New(func(prefix, args string) { t.Logf("otel: %s", fmt.Sprint(prefix, args)) }, funcr.Options{Verbosity: 1}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) { t.Errorf("otel: handler returned an error: %s", err.Error()) }))
	//otel.SetLogger(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)))
	//otel.SetLogger(stdr.New(slog.NewLogLogger(slog.Default().Handler(), slog.LevelWarn)))
	tp, err := trace_.NewTracerProvider(ctx, trace_.WithOptionExporterEndpoints(
		"stdout://localhost?allow_stdout&pretty_print&no_timestamps",
		//`otlp-http://some_endpoint/some_path?compression=gzip&insecure`,
		//`otlp-grpc://some_endpoint/some_path?compression=gzip&insecure`,
	), trace_.WithOptionResourceAttrs(
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String("demo-client"), // 必填，服务名称
	), trace_.WithOptionTracerProviderOptions(sdktrace.WithSampler(sdktrace.AlwaysSample())))
	if err != nil {
		t.Fatalf("create meter provider failed: %s", err.Error())
		return
	}
	otel.SetTracerProvider(tp)
	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			t.Fatalf("shutdown meter provider failed: %s", err.Error())
			return
		}
	}()
	defer func() {
		err := tp.ForceFlush(context.Background())
		if err != nil {
			t.Fatalf("force plush meter provider failed: %s", err.Error())
			return
		}
	}()
	ctx, span := otel.Tracer("").Start(context.Background(), "TestNewTracerProvider")
	defer span.End()
	t.Logf("trace_id: %s", span.SpanContext().TraceID())
	t.Logf("span_id: %s", span.SpanContext().SpanID())
	for range testLoop {
		for range testLoop {
			testAllTraces(ctx, t)
		}
		if testForcePush {
			err := tp.ForceFlush(context.Background())
			if err != nil {
				t.Fatalf("force plush trace provider failed: %s", err.Error())
				return
			}
		}
		t.Logf("sleep %s", testInterval)
		time.Sleep(testInterval)
	}
}

func testAllTraces(ctx context.Context, t *testing.T) {
	ctx, span := otel.Tracer("").Start(context.Background(), "testAllTraces", trace.WithNewRoot())
	defer span.End()
	v := multiply(ctx, 2, 2)
	v = multiply(ctx, v, 10)
	add(ctx, v, 2)
}

func add(ctx context.Context, x, y int64) int64 {
	var span trace.Span
	ctx, span = otel.Tracer("").Start(ctx, "Addition", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	v := x + y
	span.SetAttributes(attribute.Int64("x", x), attribute.Int64("y", y), attribute.Int64("x+y", v))
	time.Sleep(10 * time.Millisecond)
	return v
}

func multiply(ctx context.Context, x, y int64) int64 {
	var span trace.Span
	ctx, span = otel.Tracer("").Start(ctx, "Multiplication", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	v := x * y
	span.SetAttributes(attribute.Int64("x", x), attribute.Int64("y", y), attribute.Int64("x*y", v))
	time.Sleep(20 * time.Millisecond)
	return v
}
