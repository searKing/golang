// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"bytes"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gopkg.in/yaml.v3"
)

var _ runtime.Marshaler = (*YamlMarshaller)(nil)

// YamlMarshaller is a Marshaler which marshals/unmarshals into/from YAML
// with the "gopkg.in/yaml.v3" marshaler.
// It supports the full functionality of protobuf unlike JSONBuiltin.
//
// The NewDecoder method returns a DecoderWrapper, so the underlying
// *yaml.Decoder methods can be used.
type YamlMarshaller struct{}

// ContentType returns the Content-Type which this marshaler is responsible for.
func (*YamlMarshaller) ContentType(_ any) string {
	return "application/yaml"
}

// Marshal marshals "v" into byte sequence.
func (*YamlMarshaller) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal unmarshal "data" into "v".
// "v" must be a pointer value.
func (y *YamlMarshaller) Unmarshal(data []byte, v any) error {
	return y.NewDecoder(bytes.NewReader(data)).Decode(v)
}

// NewDecoder returns a Decoder which reads byte sequence from "r".
func (*YamlMarshaller) NewDecoder(r io.Reader) runtime.Decoder {
	return yaml.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes bytes sequence into "w".
func (*YamlMarshaller) NewEncoder(w io.Writer) runtime.Encoder {
	return yaml.NewEncoder(w)
}

// YamlDecoderWrapper is a wrapper around a *json.Decoder that adds
// support for proto and json to the Decode method.
type YamlDecoderWrapper struct {
	decoderYaml *yaml.Decoder // json -> interface{}
}

// Decode wraps the embedded decoder's Decode method to support
// protos using a jsonpb.Unmarshaler.
func (d YamlDecoderWrapper) Decode(v any) error {
	return d.decoderYaml.Decode(v)
}
