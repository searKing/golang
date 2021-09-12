// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

// gRPC tracing middleware
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/rpc.md
import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	grpccodes "google.golang.org/grpc/codes"
)

// spanInfo returns a span name and all appropriate attributes from the gRPC
// method and peer address.
func spanInfo(fullMethod, peerAddress, localAddress string, grpcType grpcType) (string, []attribute.KeyValue) {
	attrs := []attribute.KeyValue{semconv.RPCSystemKey.String("grpc")}
	name, mAttrs := parseFullMethod(fullMethod)
	attrs = append(attrs, mAttrs...)
	attrs = append(attrs, peerAttr(peerAddress)...)
	attrs = append(attrs, localAttr(localAddress)...)
	attrs = append(attrs, grpcTypeAttr(grpcType))
	return name, attrs
}

func fullMethodFromHttpRequest(req *http.Request) string {
	if req == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s", req.Method, req.URL.Path)
}

// peerAttr returns attributes about the peer address.
func peerAttr(addr string) []attribute.KeyValue {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return []attribute.KeyValue(nil)
	}

	if host == "" {
		host = "127.0.0.1"
	}

	return []attribute.KeyValue{
		semconv.NetPeerIPKey.String(host),
		semconv.NetPeerPortKey.String(port),
	}
}

// localAttr returns attributes about the local address.
func localAttr(addr string) []attribute.KeyValue {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return []attribute.KeyValue(nil)
	}

	if host == "" {
		host = "127.0.0.1"
	}

	return []attribute.KeyValue{
		semconv.NetHostIPKey.String(host),
		semconv.NetHostPortKey.String(port),
	}
}

// parseFullMethod returns a span name following the OpenTelemetry semantic
// conventions as well as all applicable span attribute.KeyValue attributes based
// on a gRPC's FullMethod.
func parseFullMethod(fullMethod string) (string, []attribute.KeyValue) {
	name := strings.TrimLeft(fullMethod, "/")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		// Invalid format, does not follow `/package.service/method`.
		return name, []attribute.KeyValue(nil)
	}

	var attrs []attribute.KeyValue
	if service := parts[0]; service != "" {
		attrs = append(attrs, semconv.RPCServiceKey.String(service))
	}
	if method := parts[1]; method != "" {
		attrs = append(attrs, semconv.RPCMethodKey.String(method))
	}
	return name, attrs
}

// statusCodeAttr returns status code attribute based on given gRPC code
func statusCodeAttr(c grpccodes.Code) attribute.KeyValue {
	return otelgrpc.GRPCStatusCodeKey.Int64(int64(c))
}

// grpcTypeAttr returns gRPC type attribute based on given gRPC type
func grpcTypeAttr(t grpcType) attribute.KeyValue {
	return otelgrpc.GRPCTypeKey.String(string(t))
}
