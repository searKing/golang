// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixTime = UnixSecondTime
type UnixSecondTime struct {
	time.Time
}

func (t UnixSecondTime) unit() time.Duration {
	return time.Second
}

func (t UnixSecondTime) String() string {
	return fmt.Sprintf("%d", t.Unix())
}

func (t UnixSecondTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a number of Unix seconds.
func (t *UnixSecondTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	t.Time = time.Unix(timestamp, 0)
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in Unix seconds.
func (t UnixSecondTime) MarshalText() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be a number of Unix seconds.
func (t *UnixSecondTime) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}
