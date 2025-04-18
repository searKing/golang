// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type HTTPHandler interface {
	Register(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
}
type HTTPHandlerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

func (f HTTPHandlerFunc) Register(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return f(ctx, mux, endpoint, opts)
}
