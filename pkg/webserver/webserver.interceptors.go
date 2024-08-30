// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"github.com/rs/cors"
	"google.golang.org/grpc"

	http_ "github.com/searKing/golang/go/net/http"
	"github.com/searKing/golang/pkg/webserver/pkg/recovery"
	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/burstlimit"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/timeoutlimit"
)

// UnaryHandler returns a new unary server handler.
//
// gRPC-Gateway does not support gRPC interceptors when call gRPC's service handler in process.
// See: https://github.com/grpc-ecosystem/grpc-gateway/issues/1043
func (f *Factory) UnaryHandler(handlers ...grpc_.UnaryHandlerDecorator) []grpc_.UnaryHandlerDecorator {
	// recover
	handlers = append(handlers, grpc_.UnaryHandlerDecoratorFunc(recovery.UnaryHandler))
	return handlers
}

func (f *Factory) UnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	// recover
	interceptors = append(interceptors, recovery.UnaryServerInterceptor())
	// handle request timeout
	interceptors = append(interceptors, timeoutlimit.UnaryServerInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	interceptors = append(interceptors, burstlimit.UnaryServerInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	return interceptors
}

func (f *Factory) StreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) []grpc.StreamServerInterceptor {
	// recover
	interceptors = append(interceptors, recovery.StreamServerInterceptor())
	// handle request timeout
	interceptors = append(interceptors, timeoutlimit.StreamServerInterceptor(f.fc.HandledTimeoutUnary))
	// burst limit
	interceptors = append(interceptors, burstlimit.StreamServerInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary))
	return interceptors
}

func (f *Factory) HttpServerInterceptors(decorators ...http_.HandlerDecorator) []http_.HandlerDecorator {
	// cors
	decorators = append(decorators, http_.HandlerDecoratorFunc(cors.New(f.fc.Cors).Handler))
	return decorators
}
