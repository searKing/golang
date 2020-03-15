// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x_request_id

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	http_ "github.com/searKing/golang/go/net/http"
)

// ClientInterceptor returns a new client interceptors with x-request-id in context and request's Header.
func ClientInterceptor(next http_.RoundTripHandler, keys ...interface{}) http_.RoundTripHandler {
	return http_.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		req = req.WithContext(newContextForHandleClientRequestID(req, keys...))
		return next.RoundTrip(req)
	})
}

// ClientChainedInterceptor returns a new client interceptors with x-request-id chain in context and request's Header.
func ClientChainedInterceptor(next http_.RoundTripHandler, keys ...interface{}) http_.RoundTripHandler {
	return http_.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		req = req.WithContext(newContextForHandleClientRequestIDChain(req, keys...))
		return next.RoundTrip(req)
	})
}

// key is RequestID within Context if have
func newContextForHandleClientRequestID(r *http.Request, keys ...interface{}) context.Context {
	requestIDs, ok := fromHTTPContext(r)
	if !ok || len(requestIDs) == 0 {
		return appendInOutMetadata(r.Context(), r, newRequestID(r.Context(), keys...)...)
	}
	return appendInOutMetadata(r.Context(), r, requestIDs...)
}

// to chain multiple request ids by generating new request id for each request and concatenating it to original request ids.
// key is RequestID within Context if have
func newContextForHandleClientRequestIDChain(r *http.Request, keys ...interface{}) context.Context {
	requestIDs, ok := fromHTTPContext(r)
	if !ok || len(requestIDs) == 0 {
		return appendInOutMetadata(r.Context(), r, newRequestIDChain(r.Context(), keys...)...)
	}
	requestIDs = append(requestIDs, uuid.New().String())
	return appendInOutMetadata(r.Context(), r, requestIDs...)
}
