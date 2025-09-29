// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otel

import (
	"net/http"
	"slices"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats/opentelemetry"

	http_ "github.com/searKing/golang/go/net/http"
)

// Both of HTTP and gRPC transfer caller's trace by HTTP Header "traceparent".
//
// HTTP transfers injected tracer keys into HTTP Request's Header and Context by otelhttp.
// gRPC transfers injected tracer keys into HTTP Request's Header and Context by otelgrpc.
// injected tracer keys: "tracestate", "baggage", "traceparent" and so on.
//
// HTTP: Context
// gRPC: IncomingContext<Recv Req>、OutgoingContext<send Req>、ServerMetadataContext<Send or Recv Resp>
//
// middlewares:
// 1) otelhttp:
// 1.1) HTTP Req Header -> Context<send or recv Req>
// 2) otelgrpc:
// 2.1) gRPC Req Header -> IncomingContext<recv Req>
// 2.2) gRPC OutgoingContext<send Req> -> gRPC Req Header
// 3) otelhttp2grpc:
// 3.1) HTTP Req Context -> gRPC Req Header
// 4) gRPC Gateway:
// 4.1) otelhttp2grpc-endpoint:
//   HTTP Req Header -> gRPC OutgoingContext
// 4.2) otelhttp2grpc-no endpoint:
//   HTTP Req Header -> gRPC IncomingContext

func DialOptions() []grpc.DialOption {
	// 2.1) gRPC OutgoingContext<send Req> -> gRPC Req Header
	return []grpc.DialOption{grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		opentelemetry.DialOption(opentelemetry.Options{
			MetricsOptions: opentelemetry.MetricsOptions{
				MeterProvider: otel.GetMeterProvider(),
				Metrics:       opentelemetry.DefaultMetrics,
			},
		})}
}

func ServerOptions() []grpc.ServerOption {
	// 2.2) gRPC Req Header -> IncomingContext<recv Req>
	return []grpc.ServerOption{grpc.StatsHandler(otelgrpc.NewServerHandler()),
		opentelemetry.ServerOption(opentelemetry.Options{
			MetricsOptions: opentelemetry.MetricsOptions{
				MeterProvider: otel.GetMeterProvider(),
				Metrics:       opentelemetry.DefaultMetrics,
			}})}
}

func tracerHeaderMatcher(key string) (string, bool) {
	if alias := http.CanonicalHeaderKey(key); slices.ContainsFunc(
		otel.GetTextMapPropagator().Fields(), func(s string) bool { return http.CanonicalHeaderKey(s) == alias }) {
		return alias, true
	}
	return runtime.DefaultHeaderMatcher(key)
}

func ServeMuxOptions() []runtime.ServeMuxOption {
	// Trace Header Matcher From HTTP To gRPC!
	return []runtime.ServeMuxOption{
		// 2.1) gRPC Req Header -> IncomingContext<recv Req>
		runtime.WithIncomingHeaderMatcher(tracerHeaderMatcher),
		// 2.2) gRPC OutgoingContext<send Req> -> gRPC Req Header
		runtime.WithOutgoingHeaderMatcher(tracerHeaderMatcher)}
}

// HttpHandlerDecorators adds metrics and tracing to requests if the incoming request is sampled.
func HttpHandlerDecorators() []http_.HandlerDecorator {
	return []http_.HandlerDecorator{

		// 1.1) HTTP Req Header -> Context<send or recv Req>
		http_.HandlerDecoratorFunc(func(handler http.Handler) http.Handler {
			wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Add the http.target attribute to the otelhttp span
				// Workaround for https://github.com/open-telemetry/opentelemetry-go-contrib/issues/3743
				if r.URL != nil {
					trace.SpanFromContext(r.Context()).SetAttributes(semconv.HTTPTarget(r.URL.RequestURI()))
				}
				handler.ServeHTTP(w, r)
			})
			// With Noop TracerProvider, the otelhttp still handles context propagation.
			// See https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/examples/passthrough
			return otelhttp.NewHandler(wrappedHandler, "gRPC-Gateway",
				otelhttp.WithPublicEndpoint(),
				otelhttp.WithMeterProvider(otel.GetMeterProvider()),
				otelhttp.WithTracerProvider(otel.GetTracerProvider()),
				otelhttp.WithPropagators(otel.GetTextMapPropagator()),
				otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
					if operation == "" {
						operation = "HTTP"
					}
					return operation + " " + r.Method + " " + r.URL.Path
				}))
		}),

		// 3.1) HTTP Req Context -> gRPC Req Header
		http_.HandlerDecoratorFunc(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(r.Header))
				handler.ServeHTTP(w, r)
			})
		}),

		// Inject "trace_id" and "span_id" into logging fields!
		http_.HandlerDecoratorFunc(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fields := spanFields(trace.SpanFromContext(r.Context()))
				if len(fields) > 0 {
					r = r.WithContext(logging.InjectFields(r.Context(), fields))
				}
				handler.ServeHTTP(w, r)
			})
		}),
	}
}

func spanFields(span trace.Span) logging.Fields {
	var fields logging.Fields
	if span.SpanContext().HasTraceID() {
		fields = append(fields, "trace_id", span.SpanContext().TraceID())
	}
	if span.SpanContext().HasSpanID() {
		fields = append(fields, "span_id", span.SpanContext().SpanID())
	}
	return fields
}
