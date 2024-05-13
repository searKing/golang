// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

// JSONPb is a Marshaler which marshals/unmarshals into/from JSON
// with the [proto.Message] by "google.golang.org/protobuf/encoding/protojson" marshaler and
// [any] by "encoding/json" marshaler.
// It supports the full functionality of protobuf unlike JSONBuiltin.
//
// The NewDecoder method returns a DecoderWrapper, so the underlying
// *json.Decoder methods can be used.
//
//go:generate go-option -type=JSONPb
type JSONPb struct {
	runtime.JSONPb
}

func (j *JSONPb) Marshal(v any) ([]byte, error) {
	// proto -> json
	if _, ok := v.(proto.Message); ok {
		return j.JSONPb.Marshal(v)
	}

	// interface{} -> json
	var buf bytes.Buffer
	if err := j.marshalTo(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal unmarshals JSON "data" into "v"
func (j *JSONPb) Unmarshal(data []byte, v any) error {
	return j.NewDecoder(bytes.NewReader(data)).Decode(v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (j *JSONPb) NewDecoder(r io.Reader) runtime.Decoder {
	return DecoderWrapper{
		decoderProto: j.JSONPb.NewDecoder(r),
		decoderJson:  json.NewDecoder(r),
	}
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONPb) NewEncoder(w io.Writer) runtime.Encoder {
	return j.JSONPb.NewEncoder(w)
}

// interface{} -> json
func (j *JSONPb) marshalTo(w io.Writer, v any) error {
	marshal := func() ([]byte, error) {
		if _, ok := v.(proto.Message); ok {
			return j.JSONPb.Marshal(v)
		}

		return json.Marshal(v)
	}
	buf, err := marshal()
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

// DecoderWrapper is a wrapper around a *json.Decoder that adds
// support for proto and json to the Decode method.
type DecoderWrapper struct {
	decoderProto runtime.Decoder // json -> proto
	decoderJson  *json.Decoder   // json -> interface{}
}

// Decode wraps the embedded decoder's Decode method to support
// protos using a jsonpb.Unmarshaler.
func (d DecoderWrapper) Decode(v any) error {
	if _, ok := v.(proto.Message); ok {
		return d.decoderProto.Decode(v)
	}
	return d.decoderJson.Decode(v)
}
