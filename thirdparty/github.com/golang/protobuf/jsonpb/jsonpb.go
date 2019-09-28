package jsonpb

import (
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

//go:generate go-option -type=Marshaler
type Marshaler struct {
	jsonpb.Marshaler
}

// Whether to render enum values as integers, as opposed to string values.
func WithEnumsAsInts(enumsAsInts bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.EnumsAsInts = enumsAsInts
	})
}

// Whether to render fields with zero values.
func WithEmitDefaults(emitDefaults bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.EmitDefaults = emitDefaults
	})
}

// A string to indent each level by. The presence of this field will
// also cause a space to appear between the field separator and
// value, and for newlines to be appear between fields and array
// elements.
func WithIndent(indent string) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.Indent = indent
	})
}

// Whether to use the original (.proto) name for fields.
func WithOrigName(origName bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.OrigName = origName
	})
}

// A custom URL resolver to use when marshaling Any messages to JSON.
// If unset, the default resolution strategy is to extract the
// fully-qualified type name from the type URL and pass that to
// proto.MessageType(string).
func WithAnyResolver(anyResolver jsonpb.AnyResolver) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.AnyResolver = anyResolver
	})
}

// Marshal returns the JSON encoding of v.
func Marshal(pb proto.Message, options ...MarshalerOption) ([]byte, error) {
	var buf bytes.Buffer
	m := Marshaler{
		Marshaler: jsonpb.Marshaler{EmitDefaults: false, OrigName: true},
	}
	m.ApplyOptions(options...)

	if err := m.Marshal(&buf, pb); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
