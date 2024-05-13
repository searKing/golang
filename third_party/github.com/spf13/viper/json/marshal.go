// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Marshal returns the JSON encoding of v.
// enhance json.Marshal with map's key type, not only string, like interface{}
func Marshal(v any) ([]byte, error) {
	{
		data, err := json.Marshal(v)
		if err == nil {
			return data, err
		}
	}

	// If key of a map is not string, like interface{}, which should implement encoding.TextMarshaler(),
	// or json.Marshal will complain about `error decoding '': json: unsupported type: map[interface {}]interface {}"`
	// recover this complaint by yaml to transcode.
	// map[interface{}]interface{} -> map[string]interface{}
	dataBytes, err := yaml.Marshal(v)
	if err != nil {
		return nil, err
	}

	var d any
	err = yaml.Unmarshal(dataBytes, &d)
	if err != nil {
		return nil, err
	}

	return json.Marshal(d)
}

// ValueMarshalable translate key of map from interface{} to string
func ValueMarshalable(v any) any {
	if _, err := json.Marshal(v); err == nil {
		return v
	}
	// If key of a map is not string, like interface{}, which should implement encoding.TextMarshaler(),
	// or json.Marshal will complain about `error decoding '': json: unsupported type: map[interface {}]interface {}"`
	// recover this complaint by yaml to transcode.
	// map[interface{}]interface{} -> map[string]interface{}
	dataBytes, err := yaml.Marshal(v)
	if err != nil {
		return v
	}

	var d any
	err = yaml.Unmarshal(dataBytes, &d)
	if err != nil {
		return v
	}
	return d
}
