// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tags

import (
	"context"
	"path"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	context_ "github.com/searKing/golang/go/context"
)

var (
	// SystemField is used in every statement made through value. Can be overwritten before any initialization code.
	SystemField = "system"

	// KindField indicates whether this is a server or a client interceptor.
	KindField = "grpc.kind"

	// ServiceField indicates rpc's service name
	ServiceField = "grpc.service"
	// MethodField indicates rpc's service method
	MethodField = "grpc.method"
)

// UnaryServerInterceptor returns a new unary server interceptors with tags in context.
func UnaryServerInterceptor(key any, values map[string]any) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		newCtx := newContextTagsForCall(ctx, KindServer, info.FullMethod, key, values)
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor with tags in context.
func StreamServerInterceptor(key any, values map[string]any) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := newContextTagsForCall(stream.Context(), KindServer, info.FullMethod, key, values)
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

func newContextTagsForCall(ctx context.Context, kind Kind, fullMethodString string, key any, values map[string]any) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	tags := context_.NewMapTags(context_.WithMapTagsMimeKey())
	tags.Set(SystemField, "grpc")
	tags.Set(KindField, kind)
	tags.Set(ServiceField, service)
	tags.Set(MethodField, method)
	for key, val := range values {
		tags.Set(key, val)
	}
	return context_.WithTags(ctx, key, tags)
}
