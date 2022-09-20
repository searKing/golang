// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"errors"
	"io"

	"github.com/gin-gonic/gin/binding"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

// ProtoMarshaller []byte -> proto|interface{}
type ProtoMarshaller struct {
	proto.MarshalOptions
	proto.UnmarshalOptions
	// runtime.ProtoMarshaller
}

// ContentType always returns "application/x-protobuf".
func (*ProtoMarshaller) ContentType(_ interface{}) string {
	return binding.MIMEPROTOBUF
}

// Marshal marshals "value" into Proto
func (marshaller *ProtoMarshaller) Marshal(value interface{}) ([]byte, error) {
	message, ok := value.(proto.Message)
	if !ok {
		return nil, errors.New("unable to marshal non proto field")
	}
	return marshaller.MarshalOptions.Marshal(message)
}

// Unmarshal unmarshals proto "data" into "value"
func (marshaller *ProtoMarshaller) Unmarshal(data []byte, value interface{}) error {
	message, ok := value.(proto.Message)
	if !ok {
		return errors.New("unable to unmarshal non proto field")
	}
	return marshaller.UnmarshalOptions.Unmarshal(data, message)
}

// NewDecoder returns a Decoder which reads proto stream from "reader".
func (marshaller *ProtoMarshaller) NewDecoder(reader io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(value interface{}) error {
		buffer, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		return marshaller.Unmarshal(buffer, value)
	})
}

// NewEncoder returns an Encoder which writes proto stream into "writer".
func (marshaller *ProtoMarshaller) NewEncoder(writer io.Writer) runtime.Encoder {
	return runtime.EncoderFunc(func(value interface{}) error {
		buffer, err := marshaller.Marshal(value)
		if err != nil {
			return err
		}
		_, err = writer.Write(buffer)
		if err != nil {
			return err
		}
		return nil
	})
}
