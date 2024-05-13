// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structpb

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// ToProtoStruct converts v, which must marshal into a JSON object,
// into a Google Struct proto.
func ToProtoStruct(v any) (*structpb.Struct, error) {
	if v == nil {
		return &structpb.Struct{}, nil
	}

	// Fast path: if v is already a *structpb.Struct, nothing to do.
	if s, ok := v.(*structpb.Struct); ok {
		return s, nil
	}

	var jb []byte
	switch v.(type) {
	case []byte:
		jb = v.([]byte)
	case *[]byte:
		jb = *(v.(*[]byte))
	case string:
		jb = []byte(v.(string))
	case *string:
		jb = []byte(*(v.(*string)))
	case proto.Message:
		// v is a Go struct that supports JSON marshalling. We want a Struct
		// protobuf. Some day we may have a more direct way to get there, but right
		// now the only way is to marshal the Go struct to JSON, unmarshal into a
		// map, and then build the Struct proto from the map.
		m := protojson.MarshalOptions{EmitUnpopulated: true}
		dataStr, err := m.Marshal(v.(proto.Message))
		if err != nil {
			return nil, fmt.Errorf("jsonpb.Marshal: %v", err)
		}
		jb = dataStr
	default:
		var err error
		jb, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("jsonpb.Marshal: %v", err)
		}
	}

	var dataStructpb structpb.Struct

	if err := protojson.Unmarshal(jb, &dataStructpb); err != nil {
		return nil, err
	}
	return &dataStructpb, nil
}
