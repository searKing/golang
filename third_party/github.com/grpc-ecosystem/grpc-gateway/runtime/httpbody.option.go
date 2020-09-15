// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "github.com/grpc-ecosystem/grpc-gateway/runtime"

// Whether to render enum values as integers, as opposed to string values.
func WithMarshaler(marshaler runtime.Marshaler) HTTPBodyPbOption {
	return HTTPBodyPbOptionFunc(func(pb *HTTPBodyPb) {
		pb.Marshaler = marshaler
	})
}

func NewHTTPBodyJsonMarshaler(options ...JSONPbOption) runtime.Marshaler {
	// github.com/grpc-ecosystem/grpc-gateway/runtime/Handler.go
	// fulfill if httpBodyMarshaler, ok := marshaler.(*HTTPBodyMarshaler); ok
	var o = []JSONPbOption{WithOrigName(true), WithEmitDefaults(true), WithIndent("\t")}
	o = append(o, options...)
	return (*runtime.HTTPBodyMarshaler)((&HTTPBodyPb{}).ApplyOptions(
		WithMarshaler((&JSONPb{}).ApplyOptions(o...))))
}

func NewHTTPBodyProtoMarshaler() runtime.Marshaler {
	// github.com/grpc-ecosystem/grpc-gateway/runtime/Handler.go
	// fulfill if httpBodyMarshaler, ok := marshaler.(*HTTPBodyMarshaler); ok
	return (*runtime.HTTPBodyMarshaler)((&HTTPBodyPb{}).ApplyOptions(
		WithMarshaler(&ProtoMarshaller{})))
}

func NewHTTPBodyYamlMarshaler() runtime.Marshaler {
	// github.com/grpc-ecosystem/grpc-gateway/runtime/Handler.go
	// fulfill if httpBodyMarshaler, ok := marshaler.(*HTTPBodyMarshaler); ok
	return (*runtime.HTTPBodyMarshaler)((&HTTPBodyPb{}).ApplyOptions(
		WithMarshaler(&YamlMarshaller{})))
}
