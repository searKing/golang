// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixTimeMinute struct {
	time.Time
}

func (t UnixTimeMinute) unit() time.Duration {
	return time.Minute
}

func (t UnixTimeMinute) String() string {
	ratio, divide := RatioFrom(time.Second, t.unit())
	if divide {
		return fmt.Sprintf("%d", t.Unix()/int64(ratio))
	}
	return fmt.Sprintf("%d", t.Unix()*int64(ratio))
}

func (t UnixTimeMinute) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a number of Unix minutes.
func (t *UnixTimeMinute) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	t.Time = UnixWithUnit(timestamp, t.unit())
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in Unix minutes.
func (t UnixTimeMinute) MarshalText() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be a number of Unix minutes.
func (t *UnixTimeMinute) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}
