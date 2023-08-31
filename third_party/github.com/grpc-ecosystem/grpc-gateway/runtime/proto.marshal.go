// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// ProtoMarshaller is a Marshaller which marshals/unmarshals into/from serialize proto bytes,
// but ContentType() always returns "application/x-protobuf".
// []byte -> proto|interface{}
type ProtoMarshaller struct {
	runtime.ProtoMarshaller
}

// ContentType always returns "application/x-protobuf".
func (*ProtoMarshaller) ContentType() string {
	return binding.MIMEPROTOBUF
}
