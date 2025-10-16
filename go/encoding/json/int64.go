// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"encoding/json"
	"strconv"
)

var _ json.Unmarshaler = (*Int64)(nil)
var _ json.Marshaler = Int64(0)

type Int64 int64

func (n *Int64) UnmarshalJSON(data []byte) error {
	if string(data) == `""` {
		*n = 0
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}
	value, err := num.Int64()
	if err != nil {
		return err
	}
	*n = Int64(value)
	return nil
}

func (n Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(n), 10))
}
