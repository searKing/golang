// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// HTTPBodyPb is a Marshaler which supports marshaling of a
// google.api.HttpBody message as the full response body if it is
// the actual message used as the response. If not, then this will
// simply fallback to the Marshaler specified as its default Marshaler.
//go:generate go-option -type=HTTPBodyPb
type HTTPBodyPb runtime.HTTPBodyMarshaler
