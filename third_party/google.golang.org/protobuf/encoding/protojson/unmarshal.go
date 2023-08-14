// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protojson

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Unmarshaler is a configurable object for converting from a JSON
// representation to a protocol buffer object.
//
//go:generate go-option -type=Unmarshaler
type Unmarshaler struct {
	protojson.UnmarshalOptions
}

// WithUnmarshalAllowPartial will not return an error if input for messages that
// will result in missing required fields.
func WithUnmarshalAllowPartial(allowPartial bool) UnmarshalerOption {
	return UnmarshalerOptionFunc(func(m *Unmarshaler) {
		m.AllowPartial = allowPartial
	})
}

// WithUnmarshalDiscardUnknown ignore unknown fields.
func WithUnmarshalDiscardUnknown(discardUnknown bool) UnmarshalerOption {
	return UnmarshalerOptionFunc(func(m *Unmarshaler) {
		m.DiscardUnknown = discardUnknown
	})
}

// WithUnmarshalResolver is used for looking up types when unmarshaling
// google.protobuf.Any messages or extension fields.
// If nil, this defaults to using protoregistry.GlobalTypes.
func WithUnmarshalResolver(resolver interface {
	protoregistry.ExtensionTypeResolver
	protoregistry.MessageTypeResolver
}) UnmarshalerOption {
	return UnmarshalerOptionFunc(func(m *Unmarshaler) {
		m.Resolver = resolver
	})
}

// Unmarshal reads the given []byte and populates the given proto.Message
// using options in the UnmarshalOptions object.
// It will clear the message first before setting the fields.
// If it returns an error, the given message may be partially set.
// The provided message must be mutable (e.g., a non-nil pointer to a message).
func Unmarshal(data []byte, pb proto.Message, options ...UnmarshalerOption) error {
	m := Unmarshaler{
		UnmarshalOptions: protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true},
	}

	m.ApplyOptions(options...)

	return m.Unmarshal(data, pb)
}
