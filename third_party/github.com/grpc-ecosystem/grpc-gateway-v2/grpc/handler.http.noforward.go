// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"
)

type HTTPNoForwardHandler interface {
	Register(ctx context.Context, mux *runtime.ServeMux, decorators ...grpc_.UnaryHandlerDecorator) error
}
type HTTPNoForwardHandlerFunc func(ctx context.Context, mux *runtime.ServeMux, decorators ...grpc_.UnaryHandlerDecorator) error

func (f HTTPNoForwardHandlerFunc) Register(ctx context.Context, mux *runtime.ServeMux, decorators ...grpc_.UnaryHandlerDecorator) error {
	return f(ctx, mux, decorators...)
}
