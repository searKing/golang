// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protojson

import (
	"bytes"
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

//go:generate go-option -type=Marshaler

// Marshaler is a configurable object for converting between
// protocol buffer objects and a JSON representation for them.
type Marshaler struct {
	protojson.MarshalOptions
}

// WithMarshalMultiline specifies whether the marshaler should format the output in
// indented-form with every textual element on a new line.
// If Indent is an empty string, then an arbitrary indent is chosen.
func WithMarshalMultiline(multiline bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.Multiline = multiline
	})
}

// WithMarshalIndent specifies the set of indentation characters to use in a multiline
// formatted output such that every entry is preceded by Indent and
// terminated by a newline. If non-empty, then Multiline is treated as true.
// Indent can only be composed of space or tab characters.
func WithMarshalIndent(indent string) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.Indent = indent
	})
}

// WithMarshalAllowPartial allows messages that have missing required fields to marshal
// without returning an error. If AllowPartial is false (the default),
// Marshal will return error if there are any missed required fields.
func WithMarshalAllowPartial(allowPartial bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.AllowPartial = allowPartial
	})
}

// WithMarshalUseProtoNames uses proto field name instead of lowerCamelCase name in JSON
// field names.
func WithMarshalUseProtoNames(useProtoNames bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.UseProtoNames = useProtoNames
	})
}

// WithMarshalUseEnumNumbers emits enum values as numbers.
func WithMarshalUseEnumNumbers(useEnumNumbers bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.UseEnumNumbers = useEnumNumbers
	})
}

// WithMarshalEmitUnpopulated specifies whether to emit unpopulated fields. It does not
// emit unpopulated oneof fields or unpopulated extension fields.
// The JSON value emitted for unpopulated fields are as follows:
//  ╔═══════╤════════════════════════════╗
//  ║ JSON  │ Protobuf field             ║
//  ╠═══════╪════════════════════════════╣
//  ║ false │ proto3 boolean fields      ║
//  ║ 0     │ proto3 numeric fields      ║
//  ║ ""    │ proto3 string/bytes fields ║
//  ║ null  │ proto2 scalar fields       ║
//  ║ null  │ message fields             ║
//  ║ []    │ list fields                ║
//  ║ {}    │ map fields                 ║
//  ╚═══════╧════════════════════════════╝
func WithMarshalEmitUnpopulated(emitUnpopulated bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.EmitUnpopulated = emitUnpopulated
	})
}

// WithMarshalResolver is used for looking up types when expanding google.protobuf.Any
// messages. If nil, this defaults to using protoregistry.GlobalTypes.
func WithMarshalResolver(resolver interface {
	protoregistry.ExtensionTypeResolver
	protoregistry.MessageTypeResolver
}) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.Resolver = resolver
	})
}

// Marshal marshals the given proto.Message in the JSON format using options in
// MarshalOptions. Do not depend on the output being stable. It may change over
// time across different versions of the program.
func Marshal(pb proto.Message, options ...MarshalerOption) ([]byte, error) {
	m := Marshaler{
		MarshalOptions: protojson.MarshalOptions{
			AllowPartial: true,
		},
	}
	m.ApplyOptions(options...)
	return m.Marshal(pb)
}

// MarshalIndent is like Marshal but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(pb proto.Message, prefix, indent string, options ...MarshalerOption) ([]byte, error) {
	b, err := Marshal(pb, options...)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
