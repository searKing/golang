// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// ToGolangMap converts v into a Golang map proto.
// Deprecated: use proto.ToGolangMap
// in github.com/searKing/golang/third_party/google.golang.org/protobuf/encoding/proto instead.
func ToGolangMap(pb proto.Message) (map[string]any, error) {
	if pb == nil {
		return nil, nil
	}

	m := jsonpb.Marshaler{EmitDefaults: false, Indent: "\t", OrigName: true}
	pbStr, err := m.MarshalToString(pb)
	if err != nil {
		return nil, err
	}
	var anyJson map[string]any
	err = json.Unmarshal([]byte(pbStr), &anyJson)
	if err != nil {
		return nil, err
	}
	return anyJson, nil
}
