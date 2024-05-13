// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package any

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	any_ "github.com/golang/protobuf/ptypes/any"
	struct_ "github.com/searKing/golang/third_party/github.com/golang/protobuf/ptypes/struct"
)

// ToProtoAny converts v, which must marshal into a JSON object,
// into a Google Any proto.
// Deprecated: use anypb.ToProtoAny instead.
func ToProtoAny(data any) (*any_.Any, error) {
	if data == nil {
		return &any_.Any{}, nil
	}
	var datapb proto.Message
	switch data.(type) {
	case proto.Message:
		datapb = data.(proto.Message)
	default:
		dataStructpb, err := struct_.ToProtoStruct(data)
		if err != nil {
			return nil, err
		}
		datapb = dataStructpb
	}
	return ptypes.MarshalAny(datapb)
}
