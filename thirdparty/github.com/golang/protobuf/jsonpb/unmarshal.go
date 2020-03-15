// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonpb

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

//go:generate go-option -type=Unmarshaler

// Unmarshaler is a configurable object for converting from a JSON
// representation to a protocol buffer object.
type Unmarshaler struct {
	jsonpb.Unmarshaler
}

func WithUnmarshalAllowUnknownFields(allowUnknownFields bool) UnmarshalerOption {
	return UnmarshalerOptionFunc(func(m *Unmarshaler) {
		m.AllowUnknownFields = allowUnknownFields
	})
}

// A custom URL resolver to use when marshaling Any messages to JSON.
// If unset, the default resolution strategy is to extract the
// fully-qualified type name from the type URL and pass that to
// proto.MessageType(string).
func WithUnmarshalAnyResolver(anyResolver jsonpb.AnyResolver) UnmarshalerOption {
	return UnmarshalerOptionFunc(func(m *Unmarshaler) {
		m.AnyResolver = anyResolver
	})
}

// Unmarshal unmarshals a JSON object stream into a protocol
// buffer. This function is lenient and will decode any options
// permutations of the related Marshaler.
func Unmarshal(data []byte, pb proto.Message, options ...UnmarshalerOption) error {
	m := Unmarshaler{
		Unmarshaler: jsonpb.Unmarshaler{AllowUnknownFields: true},
	}

	m.ApplyOptions(options...)

	return m.Unmarshal(bytes.NewReader(data), pb)
}
