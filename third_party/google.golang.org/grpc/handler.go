// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"net/http"

	"github.com/searKing/golang/third_party/google.golang.org/grpc/internal/grpcutil"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

// GrpcOrDefaultHandler returns a http.Handler that delegates to grpcServer on incoming gRPC
// connections or defaultHandler otherwise. Copied from cockroachdb.
func GrpcOrDefaultHandler(grpcServer *grpc.Server, defaultHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/blob/v1.64.0/internal/transport/handler_server.go#L53
		contentType := r.Header.Get("Content-Type")
		// TODO: do we assume contentType is lowercase? we did before
		_, validContentType := grpcutil.ContentSubtype(contentType)

		var h http.Handler
		if _, ok := w.(http.Flusher); !ok { // gRPC requires a ResponseWriter supporting http.Flusher
			h = defaultHandler
		} else if r.ProtoMajor == 2 && validContentType {
			// This is a partial recreation of gRPC's example code https://github.com/grpc/grpc-go/blob/v1.64.0/server.go#L1055
			h = grpcServer
		} else {
			h = defaultHandler
		}
		if h == nil {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}), &http2.Server{})
}
