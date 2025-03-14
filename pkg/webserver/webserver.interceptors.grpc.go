// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"github.com/searKing/golang/pkg/webserver/pkg/requestid"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/tags"
	"google.golang.org/grpc"

	"github.com/searKing/golang/pkg/webserver/pkg/recovery"
	validator_ "github.com/searKing/golang/pkg/webserver/pkg/validator"
	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/burstlimit"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/timeoutlimit"
)

type tagsCtxMarker struct{}

var (
	// tagsCtxMarkerKey is the Context value marker that is used by logging middleware to read and write logging fields into context.
	tagsCtxMarkerKey = &tagsCtxMarker{}
)

// UnaryHandlers returns new unary server handlers.
//
// gRPC-Gateway does not support gRPC interceptors when call gRPC's service handler in process.
// See: https://github.com/grpc-ecosystem/grpc-gateway/issues/1043
func (f *Factory) UnaryHandlers(handlers ...grpc_.UnaryHandlerDecorator) []grpc_.UnaryHandlerDecorator {
	// recover
	handlers = append(handlers, grpc_.UnaryHandlerDecoratorFunc(recovery.UnaryHandler))
	// validate
	if v := f.fc.Validator; v != nil {
		handlers = append(handlers, validator_.UnaryHandlerDecorator(v))
	}
	// request id
	if f.fc.FillRequestId {
		handlers = append(handlers, grpc_.UnaryHandlerDecoratorFunc(requestid.UnaryHandler))
	}
	return handlers
}

func (f *Factory) UnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	// recover
	interceptors = append(interceptors, recovery.UnaryServerInterceptor())
	// handle request timeout
	interceptors = append(interceptors, timeoutlimit.UnaryServerInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	interceptors = append(interceptors, burstlimit.UnaryServerInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	// map tags in context, for request local storage
	interceptors = append(interceptors, tags.UnaryServerInterceptor(tagsCtxMarkerKey, map[string]any{}))
	// validate
	if v := f.fc.Validator; v != nil {
		interceptors = append(interceptors, validator_.UnaryServerInterceptor(v))
	}
	// request id
	if f.fc.FillRequestId {
		interceptors = append(interceptors, requestid.UnaryServerInterceptor())
	}
	return interceptors
}

func (f *Factory) StreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) []grpc.StreamServerInterceptor {
	// recover
	interceptors = append(interceptors, recovery.StreamServerInterceptor())
	// handle request timeout
	interceptors = append(interceptors, timeoutlimit.StreamServerInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	interceptors = append(interceptors, burstlimit.StreamServerInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	// map tags in context, for request local storage
	interceptors = append(interceptors, tags.StreamServerInterceptor(tagsCtxMarkerKey, map[string]any{}))
	// validate
	if v := f.fc.Validator; v != nil {
		interceptors = append(interceptors, validator_.StreamServerInterceptor(v))
	}
	// request id
	if f.fc.FillRequestId {
		interceptors = append(interceptors, requestid.StreamServerInterceptor())
	}
	return interceptors
}

func (f *Factory) UnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) []grpc.UnaryClientInterceptor {
	// handle request timeout
	interceptors = append(interceptors, timeoutlimit.UnaryClientInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	interceptors = append(interceptors, burstlimit.UnaryClientInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	// map tags in context, for request local storage
	interceptors = append(interceptors, tags.UnaryClientInterceptor(tagsCtxMarkerKey, map[string]any{}))
	// request id
	if f.fc.FillRequestId {
		interceptors = append(interceptors, requestid.UnaryClientInterceptor())
	}

	return interceptors
}

func (f *Factory) StreamClientInterceptors(interceptors ...grpc.StreamClientInterceptor) []grpc.StreamClientInterceptor {
	// handle request timeout
	interceptors = append(interceptors, timeoutlimit.StreamClientInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	interceptors = append(interceptors, burstlimit.StreamClientInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	// map tags in context, for request local storage
	interceptors = append(interceptors, tags.StreamClientInterceptor(tagsCtxMarkerKey, map[string]any{}))
	// request id
	if f.fc.FillRequestId {
		interceptors = append(interceptors, requestid.StreamClientInterceptor())
	}
	return interceptors
}
