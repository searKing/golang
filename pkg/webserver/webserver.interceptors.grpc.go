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
	var s []grpc_.UnaryHandlerDecorator
	// recover
	s = append(s, grpc_.UnaryHandlerDecoratorFunc(recovery.UnaryHandler))
	// validate
	if v := f.fc.Validator; v != nil {
		s = append(s, validator_.UnaryHandlerDecorator(v))
	}
	// request id
	if f.fc.FillRequestId {
		s = append(s, grpc_.UnaryHandlerDecoratorFunc(requestid.UnaryHandler))
	}
	return append(s, handlers...)
}

func (f *Factory) UnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	var s []grpc.UnaryServerInterceptor
	// recover
	s = append(s, recovery.UnaryServerInterceptor())
	// handle request timeout
	s = append(s, timeoutlimit.UnaryServerInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	s = append(s, burstlimit.UnaryServerInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	// map tags in context, for request local storage
	s = append(s, tags.UnaryServerInterceptor(tagsCtxMarkerKey, map[string]any{}))
	// validate
	if v := f.fc.Validator; v != nil {
		s = append(s, validator_.UnaryServerInterceptor(v))
	}
	// request id
	if f.fc.FillRequestId {
		s = append(s, requestid.UnaryServerInterceptor())
	}
	return append(s, interceptors...)
}

func (f *Factory) StreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) []grpc.StreamServerInterceptor {
	var s []grpc.StreamServerInterceptor
	// recover
	s = append(s, recovery.StreamServerInterceptor())
	// handle request timeout
	s = append(s, timeoutlimit.StreamServerInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	s = append(s, burstlimit.StreamServerInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	// map tags in context, for request local storage
	s = append(s, tags.StreamServerInterceptor(tagsCtxMarkerKey, map[string]any{}))
	// validate
	if v := f.fc.Validator; v != nil {
		s = append(s, validator_.StreamServerInterceptor(v))
	}
	// request id
	if f.fc.FillRequestId {
		s = append(s, requestid.StreamServerInterceptor())
	}
	return append(s, interceptors...)
}

func (f *Factory) UnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) []grpc.UnaryClientInterceptor {
	var s []grpc.UnaryClientInterceptor
	// handle request timeout
	s = append(s, timeoutlimit.UnaryClientInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	s = append(s, burstlimit.UnaryClientInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	// map tags in context, for request local storage
	s = append(s, tags.UnaryClientInterceptor(tagsCtxMarkerKey, map[string]any{}))
	// request id
	if f.fc.FillRequestId {
		s = append(s, requestid.UnaryClientInterceptor())
	}
	return append(s, interceptors...)
}

func (f *Factory) StreamClientInterceptors(interceptors ...grpc.StreamClientInterceptor) []grpc.StreamClientInterceptor {
	var s []grpc.StreamClientInterceptor
	// handle request timeout
	s = append(s, timeoutlimit.StreamClientInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	s = append(s, burstlimit.StreamClientInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	// map tags in context, for request local storage
	s = append(s, tags.StreamClientInterceptor(tagsCtxMarkerKey, map[string]any{}))
	// request id
	if f.fc.FillRequestId {
		s = append(s, requestid.StreamClientInterceptor())
	}
	return append(s, interceptors...)
}
