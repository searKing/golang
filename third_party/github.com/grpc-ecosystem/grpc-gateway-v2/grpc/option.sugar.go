// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	runtime_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/runtime"
)

// Grammar Sugar

func WithStreamErrorHandler(fn runtime.StreamErrorHandlerFunc) GatewayOption {
	return WithGrpcServeMuxOption(runtime.WithStreamErrorHandler(fn))
}

// WithHTTPErrorHandler replies to the request with the error.
// You can set a custom function to this variable to customize error format.
func WithHTTPErrorHandler(fn HTTPErrorHandler) GatewayOption {
	return WithGrpcServeMuxOption(runtime.WithErrorHandler(fn.HandleHTTPError))
}

// WithHttpMiddlewares sets gRPC-Gateway server middleware for all registered handlers.
// This is useful as an alternative to gRPC interceptors when using the direct-to-implementation registration methods
// and cannot rely on gRPC interceptors. It's recommended to use gRPC interceptors instead if possible.
func WithHttpMiddlewares(m ...runtime.Middleware) GatewayOption {
	return WithGrpcServeMuxOption(runtime.WithMiddlewares(m...))
}

func WithMarshalerOption(mime string, marshaler runtime.Marshaler) GatewayOption {
	return WithGrpcServeMuxOption(runtime.WithMarshalerOption(mime, marshaler))
}

func WithDefaultMarshalerOption() []GatewayOption {
	return []GatewayOption{
		WithMarshalerOption(runtime.MIMEWildcard, runtime_.NewHTTPBodyJsonMarshaler()),
		WithMarshalerOption(binding.MIMEJSON, runtime_.NewHTTPBodyJsonMarshaler()),
		WithMarshalerOption(binding.MIMEPROTOBUF, runtime_.NewHTTPBodyProtoMarshaler()),
		WithMarshalerOption(binding.MIMEYAML, runtime_.NewHTTPBodyYamlMarshaler()),
	}
}

// Deprecated: Use WithDefaultMarshalerOption instead.
func WithDefaultMarsherOption() []GatewayOption {
	return WithDefaultMarshalerOption()
}

//func WithForwardResponseMessageHandler(fn ForwardResponseMessageHandler) GatewayOption {
//	return GatewayOptionFunc(func(gateway *Gateway) {
//		runtime.WithForwardResponseOption()
//		runtime.ForwardResponseMessage = nil
//	})
//}

func WithForwardResponseOptionHandler(fn ForwardResponseOptionHandler) GatewayOption {
	return WithGrpcServeMuxOption(runtime.WithForwardResponseOption(fn.ForwardResponseOption))
}

// Deprecated: Use WithForwardResponseOptionHandler instead.
func WithForwardResponseMessageHandler(fn ForwardResponseOptionHandler) GatewayOption {
	return WithGrpcServeMuxOption(runtime.WithForwardResponseOption(fn.ForwardResponseOption))
}
