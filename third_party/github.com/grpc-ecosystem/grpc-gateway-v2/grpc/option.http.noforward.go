// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"

// WithHttpNoForwardHandlerInterceptor sets gRPC-Gateway server middlewares for all handlers
// to call grpc service function directly.
// This is useful as an alternative to gRPC interceptors when using the direct-to-implementation registration methods
// and can rely on gRPC interceptors.
// as gRPC-Gateway does not support gRPC interceptors when call gRPC's service handler in process.
// See: https://github.com/grpc-ecosystem/grpc-gateway/issues/1043
func WithHttpNoForwardHandlerInterceptor(interceptors ...grpc_.UnaryHandlerDecorator) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.httpNoForwardInterceptors = append(gateway.opt.httpNoForwardInterceptors, interceptors...)
	})
}
