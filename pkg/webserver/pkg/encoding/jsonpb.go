// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ encoding.Codec = (*JSONPb)(nil)
var _ runtime.Marshaler = (*JSONPb)(nil)

// Name is the name registered for the json codec, used as the content-subtype of grpc,
// that is "application/grpc+json".
const Name = "json"

// JSONPb is a json codec of protobuf messages by protojson, which serves both as an
// encoding.Codec of grpc, registered by encoding.RegisterCodec, and as a runtime.Marshaler
// of grpc-gateway. Values not of proto.Message are marshaled by encoding/json.
//
// Use NewJSONPb to create a JSONPb, which emits no unpopulated fields and discards unknown
// fields by default, as Tencent Cloud API 3.0 does.
type JSONPb struct {
	runtime.JSONPb
}

// Name returns the name of the codec, implementing encoding.Codec.
func (j *JSONPb) Name() string {
	return Name
}

// A JSONPbOption sets options of JSONPb.
type JSONPbOption interface {
	apply(*JSONPb)
}

// JSONPbOptionFunc wraps a function that modifies JSONPb into an
// implementation of the JSONPbOption interface.
type JSONPbOptionFunc func(*JSONPb)

func (f JSONPbOptionFunc) apply(do *JSONPb) {
	f(do)
}

// WithMarshalOptions replaces MarshalOptions, which defaults to protojson.MarshalOptions{EmitUnpopulated: false}.
func WithMarshalOptions(o protojson.MarshalOptions) JSONPbOption {
	return JSONPbOptionFunc(func(j *JSONPb) {
		j.MarshalOptions = o
	})
}

// WithUnmarshalOptions replaces UnmarshalOptions, which defaults to protojson.UnmarshalOptions{DiscardUnknown: true}.
// Set DiscardUnknown explicitly to keep ignoring unknown fields.
func WithUnmarshalOptions(o protojson.UnmarshalOptions) JSONPbOption {
	return JSONPbOptionFunc(func(j *JSONPb) {
		j.UnmarshalOptions = o
	})
}

// NewJSONPb returns a JSONPb, which emits no unpopulated fields and discards unknown fields by default.
func NewJSONPb(opts ...JSONPbOption) *JSONPb {
	j := &JSONPb{
		JSONPb: runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.apply(j)
	}
	return j
}
