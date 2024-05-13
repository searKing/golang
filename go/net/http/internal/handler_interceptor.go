// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import "net/http"

var DefaultPreHandler = func(w http.ResponseWriter, r *http.Request) error { return nil }
var DefaultWrapHandler = func(h http.Handler) http.Handler { return h }
var DefaultPostHandler = func(w http.ResponseWriter, r *http.Request) {}
var DefaultAfterCompletion = func(w http.ResponseWriter, r *http.Request, err any) {
	if err != nil {
		panic(err)
	}
}

type HandlerInterceptor struct {
	// Intercept the execution of a handler.
	// The default implementation returns true.
	// Parameters:
	// request - current HTTP request
	// response - current HTTP response
	// handler - chosen handler to execute, for type and/or instance evaluation
	// Returns:
	// true if the execution chain should proceed with the next interceptor or the handler itself.
	// Else, DispatcherServlet assumes that this interceptor has already dealt with the response itself.
	PreHandleFunc func(w http.ResponseWriter, r *http.Request) error

	// Intercept the execution of a handler.
	// The default implementation is empty.
	// Parameters:
	// handler - current HTTP handler
	// Returns:
	// handler - wrapped HTTP handler
	WrapHandleFunc func(h http.Handler) http.Handler

	// Intercept the execution of a handler.
	// The default implementation is empty.
	// Parameters:
	// request - current HTTP request
	// response - current HTTP response
	// handler - the handler (or HandlerMethod) that started asynchronous execution, for type and/or instance examination
	PostHandleFunc func(w http.ResponseWriter, r *http.Request)
	// Callback after completion of request processing, that is, after rendering the view.
	// The default implementation is empty.
	// Parameters:
	// request - current HTTP request
	// response - current HTTP response
	// handler - the handler (or HandlerMethod) that started asynchronous execution, for type and/or instance examination
	// ex - any exception thrown on handler execution, if any; this does not include exceptions that have been handled through an exception resolver
	AfterCompletionFunc func(w http.ResponseWriter, r *http.Request, err any)
}

func (filter HandlerInterceptor) PreHandle(w http.ResponseWriter, r *http.Request) error {
	if filter.PreHandleFunc == nil {
		return DefaultPreHandler(w, r)
	}
	return filter.PreHandleFunc(w, r)
}

func (filter HandlerInterceptor) WrapHandle(h http.Handler) http.Handler {
	if filter.WrapHandleFunc == nil {
		return DefaultWrapHandler(h)
	}
	return filter.WrapHandleFunc(h)
}

func (filter HandlerInterceptor) PostHandle(w http.ResponseWriter, r *http.Request) {
	if filter.PostHandleFunc == nil {
		DefaultPostHandler(w, r)
		return
	}
	filter.PostHandleFunc(w, r)
}

func (filter HandlerInterceptor) AfterCompletion(w http.ResponseWriter, r *http.Request, err any) {
	if filter.AfterCompletionFunc == nil {
		DefaultAfterCompletion(w, r, err)
		return
	}

	filter.AfterCompletionFunc(w, r, err)
}
