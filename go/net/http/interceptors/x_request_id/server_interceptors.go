// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x_request_id

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// ServerInterceptor returns a new server interceptors with x-request-id in context and response's Header.
func ServerInterceptor(next http.Handler, keys ...interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(newContextForHandleServerRequestID(w, r, keys...))
		next.ServeHTTP(w, r)
	})
}

// ServerChainedInterceptor returns a new server interceptors with x-request-id chain in context and response's Header.
func ServerChainedInterceptor(next http.Handler, keys ...interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(newContextForHandleServerRequestIDChain(w, r, keys...))
		next.ServeHTTP(w, r)
	})
}

// key is RequestID within Context if have
func newContextForHandleServerRequestID(w http.ResponseWriter, r *http.Request, keys ...interface{}) context.Context {
	requestIDs, ok := fromHTTPContext(r)
	if !ok || len(requestIDs) == 0 {
		return appendInOutMetadata(r.Context(), w, newRequestID(r.Context(), keys...)...)
	}
	return appendInOutMetadata(r.Context(), w, requestIDs...)
}

// to chain multiple request ids by generating new request id for each request and concatenating it to original request ids.
// key is RequestID within Context if have
func newContextForHandleServerRequestIDChain(w http.ResponseWriter, r *http.Request, keys ...interface{}) context.Context {
	requestIDs, ok := fromHTTPContext(r)
	if !ok || len(requestIDs) == 0 {
		return appendInOutMetadata(r.Context(), w, newRequestIDChain(r.Context(), keys...)...)
	}
	requestIDs = append(requestIDs, uuid.New().String())
	return appendInOutMetadata(r.Context(), w, requestIDs...)
}
