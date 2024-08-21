// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryHandlerGeneric is a generic version of [grpc.UnaryHandler]
// as gRPC-Gateway does not support gRPC interceptors when call gRPC's service handler in process.
// See: https://github.com/grpc-ecosystem/grpc-gateway/issues/1043
type UnaryHandlerGeneric[REQ any, RESP any] func(ctx context.Context, req REQ) (resp RESP, err error)

func (f UnaryHandlerGeneric[REQ, RESP]) UnaryHandler() grpc.UnaryHandler {
	return func(ctx context.Context, req any) (any, error) { return f(ctx, req.(REQ)) }
}

func NewUnaryHandlerGeneric[REQ any, RESP any](handler grpc.UnaryHandler) UnaryHandlerGeneric[REQ, RESP] {
	return func(ctx context.Context, req REQ) (resp RESP, err error) {
		resp_, err := handler(ctx, req)
		resp, _ = resp_.(RESP)
		return resp, err
	}
}
