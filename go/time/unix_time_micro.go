// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixTimeMicrosecond struct {
	time.Time
}

func (t UnixTimeMicrosecond) unit() time.Duration {
	return time.Microsecond
}

func (t UnixTimeMicrosecond) String() string {
	ratio, divide := RatioFrom(time.Nanosecond, t.unit())
	if divide {
		return fmt.Sprintf("%d", t.UnixNano()/int64(ratio))
	}
	return fmt.Sprintf("%d", t.UnixNano()*int64(ratio))
}

func (t UnixTimeMicrosecond) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a number of Unix microseconds.
func (t *UnixTimeMicrosecond) UnmarshalJSON(data []byte) error {
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
// The time is formatted in Unix microseconds.
func (t UnixTimeMicrosecond) MarshalText() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be a number of Unix microseconds.
func (t *UnixTimeMicrosecond) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}
