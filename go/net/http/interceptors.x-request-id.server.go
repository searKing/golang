// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// XRequestIdServerInterceptor returns a new server interceptor with x-request-id in context and response's Header.
// keys is context's key
func XRequestIdServerInterceptor(next http.Handler, keys ...interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(newContextForHandleServerRequestID(w, r, keys...))
		next.ServeHTTP(w, r)
	})
}

// key is RequestID within Context if have
func newContextForHandleServerRequestID(w http.ResponseWriter, r *http.Request, keys ...interface{}) context.Context {
	requestID := fromHTTPContext(r, keys...)
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return setInOutMetadata(r.Context(), w, r, requestID)
}
