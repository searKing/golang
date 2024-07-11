// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// WithMarshaler Whether to render enum values as integers, as opposed to string values.
func WithMarshaler(marshaler runtime.Marshaler) HTTPBodyPbOption {
	return HTTPBodyPbOptionFunc(func(pb *HTTPBodyPb) {
		pb.Marshaler = marshaler
	})
}

func NewHTTPBodyJsonMarshaler(options ...JSONPbOption) runtime.Marshaler {
	// github.com/grpc-ecosystem/grpc-gateway/runtime/Handler.go
	// fulfill if httpBodyMarshaler, ok := marshaler.(*HTTPBodyMarshaler); ok
	return &runtime.HTTPBodyMarshaler{Marshaler: &runtime.JSONPb{}}
}

func NewHTTPBodyProtoMarshaler() runtime.Marshaler {
	// github.com/grpc-ecosystem/grpc-gateway/runtime/Handler.go
	// fulfill if httpBodyMarshaler, ok := marshaler.(*HTTPBodyMarshaler); ok
	return &runtime.HTTPBodyMarshaler{Marshaler: &runtime.ProtoMarshaller{}}
}

func NewHTTPBodyYamlMarshaler() runtime.Marshaler {
	// github.com/grpc-ecosystem/grpc-gateway/runtime/Handler.go
	// fulfill if httpBodyMarshaler, ok := marshaler.(*HTTPBodyMarshaler); ok
	return &runtime.HTTPBodyMarshaler{Marshaler: &YamlMarshaller{}}
}
