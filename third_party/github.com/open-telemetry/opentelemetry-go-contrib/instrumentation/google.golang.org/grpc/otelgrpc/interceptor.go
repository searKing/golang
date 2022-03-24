// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

// gRPC tracing middleware
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/rpc.md
import (
	"context"
	"net"
	"strings"

	otelgrpc_ "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

// spanInfo returns a span name and all appropriate attributes from the gRPC
// method and peer address.
func spanInfo(fullMethod, peerAddress, localAddress string, grpcType grpcType, client bool) (string, []attribute.KeyValue) {
	attrs := []attribute.KeyValue{otelgrpc_.RPCSystemGRPC}
	name, mAttrs := parseFullMethod(fullMethod)
	attrs = append(attrs, mAttrs...)
	attrs = append(attrs, peerAttr(peerAddress, client)...)
	attrs = append(attrs, localAttr(localAddress, client)...)
	attrs = append(attrs, grpcTypeAttr(grpcType))
	return name, attrs
}

// peerAttr returns attributes about the peer address.
func peerAttr(addr string, client bool) []attribute.KeyValue {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return []attribute.KeyValue(nil)
	}

	var attrs []attribute.KeyValue
	if host != "" {
		attrs = append(attrs, semconv.NetPeerIPKey.String(host))
	}
	if port != "" && port != "0" && client { // avoid bombs of various client's peer port
		attrs = append(attrs, semconv.NetPeerPortKey.String(port))
	}
	return attrs
}

// localAttr returns attributes about the local address.
func localAttr(addr string, client bool) []attribute.KeyValue {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return []attribute.KeyValue(nil)
	}
	var attrs []attribute.KeyValue
	if host != "" {
		attrs = append(attrs, semconv.NetHostIPKey.String(host))
	}
	if port != "" && port != "0" && !client { // avoid bombs of various server's local port
		attrs = append(attrs, semconv.NetHostPortKey.String(port))
	}
	return attrs
}

// peerFromCtx returns a peer address from a context, if one exists.
func peerFromCtx(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	return p.Addr.String()
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
	return semconv.RPCGRPCStatusCodeKey.Int64(int64(c))
}

// grpcTypeAttr returns gRPC type attribute based on given gRPC type
func grpcTypeAttr(t grpcType) attribute.KeyValue {
	return GRPCTypeKey.String(string(t))
}
