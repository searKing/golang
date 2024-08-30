// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"net/http"

	slices_ "github.com/searKing/golang/go/exp/slices"
	http_ "github.com/searKing/golang/go/net/http"
)

// WithHttpHandlerDecorators sets gRPC-Gateway server middlewares for all handlers.
// This is useful as an alternative to gRPC interceptors when using the direct-to-implementation registration methods
// and can rely on gRPC interceptors.
func WithHttpHandlerDecorators(decorators ...http_.HandlerDecorator) GatewayOption {
	return WithHttpHandlerInterceptor(
		slices_.MapFunc(decorators, func(e http_.HandlerDecorator) http_.HandlerInterceptorChainOption {
			return http_.WithHandlerInterceptor(nil, e.WrapHandler, nil, nil)
		})...)
}

// below here for advance usage only

// WithHttpHandlerInterceptor sets gRPC-Gateway server middleware for all handlers.
// This is useful as an alternative to gRPC interceptors when using the direct-to-implementation registration methods
// and can rely on gRPC interceptors.
func WithHttpHandlerInterceptor(opts ...http_.HandlerInterceptorChainOption) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.interceptors.ApplyOptions(opts...)
	})
}

func WithHttpPreHandler(preHandle func(w http.ResponseWriter, r *http.Request) error) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(preHandle, nil, nil, nil))
}

// WithHttpWrapper is a decorator or middleware of http.Handler
func WithHttpWrapper(wrapper func(http.Handler) http.Handler) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(nil, wrapper, nil, nil))
}

func WithHttpPostHandler(
	postHandle func(w http.ResponseWriter, r *http.Request)) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(nil, nil, postHandle, nil))
}

func WithHttpAfterCompletion(
	afterCompletion func(w http.ResponseWriter, r *http.Request, err any)) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(nil, nil, nil, afterCompletion))
}

// Deprecated: Use WithHttpPreHandler instead.
func WithHttpRewriter(rewriter func(w http.ResponseWriter, r *http.Request) error) GatewayOption {
	return WithHttpHandlerInterceptor(
		http_.WithHandlerInterceptor(rewriter, nil, nil, nil))
}
