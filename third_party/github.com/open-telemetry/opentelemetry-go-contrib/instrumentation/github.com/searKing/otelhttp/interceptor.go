// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

// gRPC tracing middleware
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/rpc.md
import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
)

// AttrsFromRequest generates attributes as specified by the OpenTelemetry specification
// for a span.
func AttrsFromRequest(req *http.Request, localAddress string) []attribute.KeyValue {
	attrs := []attribute.KeyValue{semconv.RPCSystemKey.String("http")}
	attrs = append(attrs, httpconv.ClientRequest(req)...)
	attrs = append(attrs, semconv.HTTPClientIPKey.String(localAddress))
	return attrs
}

// AttrsFromResponse generates attributes as specified by the OpenTelemetry specification
// for a span.
func AttrsFromResponse(resp *http.Response) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	attrs = append(attrs, httpconv.ClientResponse(resp)...)
	return attrs
}
