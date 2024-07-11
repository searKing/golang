// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var _ runtime.Marshaler = (*JSONPb)(nil)

// JSONPb is a Marshaler which marshals/unmarshals into/from JSON
// with the [proto.Message] by "google.golang.org/protobuf/encoding/protojson" marshaler and
// [any] by "encoding/json" marshaler.
// It supports the full functionality of protobuf unlike JSONBuiltin.
//
// The NewDecoder method returns a DecoderWrapper, so the underlying
// *json.Decoder methods can be used.
//
// Deprecated: Use runtime.JSONPb instead.
//
//go:generate go-option -type=JSONPb
type JSONPb struct {
	runtime.JSONPb
}
