// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin/binding"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gopkg.in/yaml.v3"
)

// []byte -> proto|interface{}
type YamlMarshaller struct {
	runtime.ProtoMarshaller
}

// Marshal marshals "v" into byte sequence.
func (*YamlMarshaller) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal unmarshals "data" into "v".
// "v" must be a pointer value.
func (y *YamlMarshaller) Unmarshal(data []byte, v interface{}) error {
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

// ContentType returns the Content-Type which this marshaler is responsible for.
func (*YamlMarshaller) ContentType(_ interface{}) string {
	return binding.MIMEYAML
}

// DecoderWrapper is a wrapper around a *json.Decoder that adds
// support for proto and json to the Decode method.
type YamlDecoderWrapper struct {
	decoderYaml *yaml.Decoder // json -> interface{}
}

// Decode wraps the embedded decoder's Decode method to support
// protos using a jsonpb.Unmarshaler.
func (d YamlDecoderWrapper) Decode(v interface{}) error {
	return d.decoderYaml.Decode(v)
}
