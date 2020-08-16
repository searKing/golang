// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

const (
	// baseGrpcContentType is the base content-type for gRPC.  This is a valid
	// content-type on it's own, but can also include a content-subtype such as
	// "proto" as a suffix after "+" or ";".  See
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests
	// for more details.
	baseGrpcContentType = "application/grpc"
)

// GrpcOrDefaultHandler returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func GrpcOrDefaultHandler(grpcServer *grpc.Server, defaultHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/blob/68098483a7afa91b353453641408e3968ad92738/internal/transport/handler_server.go#L51
		contentType := r.Header.Get("Content-Type")
		// TODO: do we assume contentType is lowercase? we did before
		_, validGrpcContentType := contentSubtype(contentType)
		var h http.Handler
		// This is a partial recreation of gRPC's example code https://github.com/grpc/grpc-go/blob/68098483a7afa91b353453641408e3968ad92738/server.go#L862
		if r.ProtoMajor == 2 && validGrpcContentType {
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

// contentSubtype returns the content-subtype for the given content-type.  The
// given content-type must be a valid content-type that starts with
// "application/grpc". A content-subtype will follow "application/grpc" after a
// "+" or ";", format as "application/grpc" [("+proto" / "+json" / {custom})]. See
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests for
// more details.
//
// If contentType is not a valid content-type for gRPC, the boolean
// will be false, otherwise true. If content-type == "application/grpc",
// "application/grpc+", or "application/grpc;", the boolean will be true,
// but no content-subtype will be returned.
//
// contentType is assumed to be lowercase already.
func contentSubtype(contentType string) (string, bool) {
	if contentType == baseGrpcContentType {
		return "", true
	}
	if !strings.HasPrefix(contentType, baseGrpcContentType) {
		return "", false
	}
	// guaranteed since != baseGrpcContentType and has baseGrpcContentType prefix
	switch contentType[len(baseGrpcContentType)] {
	case '+', ';':
		// this will return true for "application/grpc+" or "application/grpc;"
		// which the previous validContentType function tested to be valid, so we
		// just say that no content-subtype is specified in this case
		return contentType[len(baseGrpcContentType)+1:], true
	default: // custom
		return "", false
	}
}

// contentSubtype is assumed to be lowercase
func contentType(contentSubtype string) string {
	if contentSubtype == "" {
		return baseGrpcContentType
	}
	return baseGrpcContentType + "+" + contentSubtype
}
