package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixTimeNanosecond struct {
	time.Time
}

func (t UnixTimeNanosecond) unit() time.Duration {
	return time.Nanosecond
}

func (t UnixTimeNanosecond) String() string {
	return fmt.Sprintf("%d", t.UnixNano())
}

func (t UnixTimeNanosecond) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UnixNano())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *UnixTimeNanosecond) UnmarshalJSON(data []byte) error {
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
// The time is formatted in Unix Seconds, with sub-second precision added if present.
func (t UnixTimeNanosecond) MarshalText() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (t *UnixTimeNanosecond) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}
