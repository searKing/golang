// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"

	"github.com/searKing/golang/go/net/http/internal"
)

func WithHandlerInterceptor(
	// Intercept the execution of a handler.
	// The default implementation returns true.
	// Parameters:
	// request - current HTTP request
	// response - current HTTP response
	// Returns:
	// true if the execution chain should proceed with the next interceptor or the handler itself.
	// Else, DispatcherServlet assumes that this interceptor has already dealt with the response itself.
	preHandle func(w http.ResponseWriter, r *http.Request) error,

	// Intercept the execution of a handler.
	// The default implementation is empty.
	// Parameters:
	// handler - current HTTP handler
	// Returns:
	// handler - wrapped HTTP handler
	wrapHandle func(h http.Handler) http.Handler,

	// Intercept the execution of a handler.
	// The default implementation is empty.
	// Parameters:
	// request - current HTTP request
	// response - current HTTP response
	postHandle func(w http.ResponseWriter, r *http.Request),
	// Callback after completion of request processing, that is, after rendering the view.
	// The default implementation is empty.
	// Parameters:
	// request - current HTTP request
	// response - current HTTP response
	// ex - any exception thrown on handler execution, if any; this does not include exceptions that have been handled through an exception resolverreturns a new server interceptors with recovery from panic.
	afterCompletion func(w http.ResponseWriter, r *http.Request, err any),
) HandlerInterceptorChainOption {
	return HandlerInterceptorChainOptionFunc(func(chain *HandlerInterceptorChain) {
		chain.interceptors = append(chain.interceptors, internal.HandlerInterceptor{
			PreHandleFunc:       preHandle,
			WrapHandleFunc:      wrapHandle,
			PostHandleFunc:      postHandle,
			AfterCompletionFunc: afterCompletion,
		})
	})
}
