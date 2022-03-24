// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anypb

import (
	"github.com/searKing/golang/third_party/google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// ToProtoAny converts v, which must marshal into a JSON object,
// into a Google Any proto.
func ToProtoAny(data interface{}) (*anypb.Any, error) {
	if data == nil {
		return &anypb.Any{}, nil
	}
	var datapb proto.Message
	switch data.(type) {
	case proto.Message:
		datapb = data.(proto.Message)
	default:
		dataStructpb, err := structpb.ToProtoStruct(data)
		if err != nil {
			return nil, err
		}
		datapb = dataStructpb
	}
	return anypb.New(datapb)
}
