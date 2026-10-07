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

// Name is the name registered for the proto compressor.
const Name = "json"

// JSONPb 遵循云API3.0标准协议的JSONPb
type JSONPb struct {
	runtime.JSONPb
}

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
