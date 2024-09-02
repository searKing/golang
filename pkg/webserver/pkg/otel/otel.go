// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otel

import (
	"net/http"
	"slices"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	http_ "github.com/searKing/golang/go/net/http"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func DialOptions() []grpc.DialOption {
	return []grpc.DialOption{grpc.WithStatsHandler(otelgrpc.NewClientHandler())}
}

func ServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{grpc.StatsHandler(otelgrpc.NewServerHandler())}
}

func ServeMuxOptions() []runtime.ServeMuxOption {
	// Trace Header Matcher From HTTP To gRPC!
	return []runtime.ServeMuxOption{runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		// HTTP Server -> gRPC Server
		if alias := http.CanonicalHeaderKey(key); slices.ContainsFunc(
			otel.GetTextMapPropagator().Fields(), func(s string) bool { return http.CanonicalHeaderKey(s) == alias }) {
			return alias, true
		}
		return runtime.DefaultHeaderMatcher(key)
	})}
}

func HttpHandlerDecorators() []http_.HandlerDecorator {
	return []http_.HandlerDecorator{
		// Root Span of HTTP!
		http_.HandlerDecoratorFunc(
			func(handler http.Handler) http.Handler {
				return otelhttp.NewHandler(handler, "gRPC-Gateway",
					otelhttp.WithMeterProvider(otel.GetMeterProvider()),
					otelhttp.WithTracerProvider(otel.GetTracerProvider()),
					otelhttp.WithPropagators(otel.GetTextMapPropagator()),
					otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
						if operation == "" {
							operation = "HTTP"
						}
						return operation + "." + r.Method + "." + r.URL.Path
					}))
			}),
		// Trace Header From HTTP To gRPC!
		http_.HandlerDecoratorFunc(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Inject set cross-cutting concerns from the Context into the carrier.
				// HTTP Server -> [gRPC Client] -> gRPC Server
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
