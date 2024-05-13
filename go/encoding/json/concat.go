// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Join(bufa []byte, bufbs ...[]byte) ([]byte, error) {
	if len(bufbs) == 0 {
		return bufa, nil
	}
	// append bs to b.
	bufb, err := Join(bufbs[0], bufbs[1:]...)
	if err != nil {
		return nil, err
	}

	// Concatenate public and private claim JSON objects.
	if !bytes.HasSuffix(bufa, []byte{'}'}) {
		return nil, fmt.Errorf("json: invalid JSON %s", bufa)
	}
	if !bytes.HasPrefix(bufb, []byte{'{'}) {
		return nil, fmt.Errorf("json: invalid JSON %s", bufb)
	}
	bufa[len(bufa)-1] = ','          // Replace closing curly brace with a comma.
	bufa = append(bufa, bufb[1:]...) // Append vb after va.
	return bufa, nil
}

// MarshalConcat returns the JSON encoding of va, vbs...,
// ignore conflict keys of json if meet later.
func MarshalConcat(va any, vbs ...any) ([]byte, error) {
	unique, err := marshalConcat(nil, va, vbs...)
	if err != nil {
		return nil, err
	}
	return json.Marshal(unique)
}

func marshalConcat(unique map[string]any, va any, vbs ...any) (map[string]any, error) {
	bufa, err := json.Marshal(va)
	if err != nil {
		return unique, err
	}

	// Marshal vbs and then append it to uniqueMap.
	var mapa map[string]any
	if err := json.Unmarshal(bufa, &mapa); err != nil {
		return unique, err
	}
	if unique == nil {
		unique = map[string]any{}
	}
	// unique
	for k, v := range mapa {
		if _, ok := unique[k]; ok {
			continue
		}
		unique[k] = v
	}

	if len(vbs) == 0 {
		return unique, nil
	}
	if len(vbs) == 1 {
		return marshalConcat(unique, vbs[0])
	}
	return marshalConcat(unique, vbs[0], vbs[1:])
}

func MarshalIndentConcat(va any, prefix, indent string, vbs ...any) ([]byte, error) {
	b, err := MarshalConcat(va, vbs)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v. If v is nil or not a pointer,
// Unmarshal returns an InvalidUnmarshalError.
// ignore conflict keys of json if meet later.
func UnmarshalConcat(data []byte, va any, vbs ...any) error {
	var unique map[string]any
	err := json.Unmarshal(data, &unique)
	if err != nil {
		return err
	}
	return unmarshalConcat(unique, va, vbs...)
}

func unmarshalConcat(unique map[string]any, va any, vbs ...any) error {
	data, err := json.Marshal(unique)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, va); err != nil {
		return err
	}
	if len(vbs) == 0 {
		return nil
	}
	dataa, err := json.Marshal(va)
	if err != nil {
		return err
	}

	var mapa map[string]any
	if err := json.Unmarshal(dataa, &mapa); err != nil {
		return err
	}

	// unique
	for k, _ := range mapa {
		delete(unique, k)
	}

	if len(vbs) == 0 {
		return nil
	}
	if len(vbs) == 1 {
		return unmarshalConcat(unique, vbs[0])
	}
	return unmarshalConcat(unique, vbs[0], vbs[1:])
}
