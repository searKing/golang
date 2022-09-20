// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MessageShortName returns the short name of m.
// If m is nil, it returns an empty string.
func MessageShortName(m proto.Message) protoreflect.Name {
	if m == nil {
		return ""
	}
	return m.ProtoReflect().Descriptor().Name()
}
