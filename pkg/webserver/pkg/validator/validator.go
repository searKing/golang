// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptor verified with validator.
func UnaryServerInterceptor(validator *validator.Validate) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return unaryHandler(validator, handler)(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor verified with validator.
func StreamServerInterceptor(validator *validator.Validate) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &validatorServerStream{ServerStream: ss, validator: validator})
	}
}

// UnaryHandlerDecorator returns a new unary server handler decorator that performs data verified with validator.
func UnaryHandlerDecorator(validator *validator.Validate) grpc_.UnaryHandlerDecorator {
	return grpc_.UnaryHandlerDecoratorFunc(func(next grpc.UnaryHandler) grpc.UnaryHandler {
		return unaryHandler(validator, next)
	})
}

func unaryHandler(validator *validator.Validate, handler func(ctx context.Context, req any) (any, error)) func(
	ctx context.Context, req any) (any, error) {
	return func(ctx context.Context, req any) (any, error) {
		if v := validator; v != nil {
			if err := v.StructCtx(ctx, req); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
		}
		// DON'T CHECK RESPONSE
		return handler(ctx, req)
	}
}
