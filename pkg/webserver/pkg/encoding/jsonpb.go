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

func NewJSONPb() *JSONPb {
	return &JSONPb{
		JSONPb: runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}
}
