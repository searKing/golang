// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"encoding/json"

	"github.com/searKing/golang/third_party/google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ToGolangMap converts v into a Golang map proto.
func ToGolangMap(pb proto.Message, options ...protojson.MarshalerOption) (map[string]any, error) {
	if pb == nil {
		return nil, nil
	}

	data, err := protojson.Marshal(pb, options...)
	if err != nil {
		return nil, err
	}
	var anyJson map[string]any
	err = json.Unmarshal(data, &anyJson)
	if err != nil {
		return nil, err
	}
	return anyJson, nil
}
