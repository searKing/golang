// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package requestid

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// See https://http.dev/x-request-id
const requestId = "X-Request-ID"

// UnaryHandler returns a new unary server handler that performs recovering from a panic.
func UnaryHandler(handler grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req any) (_ any, err error) {
		newCtx, id := tagLoggingRequestId(ctx, req)
		resp, err := handler(newCtx, req)
		trySetRequestId(resp, id, true)
		// inject "X-Request-ID" into HTTP Header
		_ = grpc.SetHeader(ctx, metadata.Pairs(requestId, id))
		return resp, err
	}
}

// UnaryServerInterceptor returns a new unary server interceptors with logging tags in context with request_id.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return UnaryHandler(handler)(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptors with logging tags in context with request_id.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &requestIdServerStream{ServerStream: ss})
	}
}

// UnaryClientInterceptor returns a new unary client interceptors with logging tags in context with request_id.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx, id := tagLoggingRequestId(ctx, req)
		err := invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			return err
		}
		trySetRequestId(reply, id, true)
		return nil
	}
}

// StreamClientInterceptor returns a new stream client interceptors with tags in context with request_id.
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		newStream := &requestIdClientStream{ClientStream: clientStream}
		return newStream, err
	}
}
