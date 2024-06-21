// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

type HTTPErrorHandler interface {
	HandleHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error)
}
type HTTPErrorHandlerFunc func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error)

func (f HTTPErrorHandlerFunc) HandleHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	f(ctx, mux, marshaler, w, r, err)
}

type ForwardResponseMessageHandler interface {
	ForwardResponseMessage(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error)
}

type ForwardResponseMessageHandlerFunc func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error)

func (f ForwardResponseMessageHandlerFunc) ForwardResponseMessage(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error) {
	f(ctx, mux, marshaler, w, req, resp, opts...)
}

type ForwardResponseOptionHandler interface {
	ForwardResponseOption(context.Context, http.ResponseWriter, proto.Message) error
}

type ForwardResponseHandlerFunc func(context.Context, http.ResponseWriter, proto.Message) error

func (f ForwardResponseHandlerFunc) ForwardResponseOption(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
	return f(ctx, w, m)
}
