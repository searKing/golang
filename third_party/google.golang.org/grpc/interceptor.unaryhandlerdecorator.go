// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import "google.golang.org/grpc"

// UnaryHandlerDecorator is an interface representing the ability to decorate or wrap a grpc.UnaryHandler.
type UnaryHandlerDecorator interface {
	WrapUnaryHandler(rt grpc.UnaryHandler) grpc.UnaryHandler
}

// The UnaryHandlerDecoratorFunc type is an adapter to allow the use of
// ordinary functions as gRPC handler decorators. If f is a function
// with the appropriate signature, UnaryHandlerDecoratorFunc(f) is a
// [UnaryHandlerDecorator] that calls f.
type UnaryHandlerDecoratorFunc func(next grpc.UnaryHandler) grpc.UnaryHandler

// WrapUnaryHandler calls f(rt).
func (f UnaryHandlerDecoratorFunc) WrapUnaryHandler(next grpc.UnaryHandler) grpc.UnaryHandler {
	return f(next)
}

// UnaryHandlerDecorators defines a UnaryHandlerDecorator slice.
type UnaryHandlerDecorators []UnaryHandlerDecorator

func (chain UnaryHandlerDecorators) WrapUnaryHandler(next grpc.UnaryHandler) grpc.UnaryHandler {
	for i := range chain {
		h := chain[len(chain)-1-i]
		if h != nil {
			next = h.WrapUnaryHandler(next)
		}
	}
	return next
}
