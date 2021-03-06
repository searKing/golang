package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixTimeMillisecond struct {
	time.Time
}

func (t UnixTimeMillisecond) unit() time.Duration {
	return time.Millisecond
}

func (t UnixTimeMillisecond) String() string {
	ratio, divide := RatioFrom(time.Nanosecond, t.unit())
	if divide {
		return fmt.Sprintf("%d", t.UnixNano()/int64(ratio))
	}
	return fmt.Sprintf("%d", t.UnixNano()*int64(ratio))
}

func (t UnixTimeMillisecond) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UnixNano() * int64(t.unit()))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *UnixTimeMillisecond) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	t.Time = UnixWithUnit(timestamp, t.unit())

	var s = t.Time.String()
	_ = s
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in Unix Seconds, with sub-second precision added if present.
func (t UnixTimeMillisecond) MarshalText() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (t *UnixTimeMillisecond) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}
