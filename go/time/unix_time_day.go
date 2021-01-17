package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixTimeDay struct {
	time.Time
}

func (t UnixTimeDay) unit() time.Duration {
	return time.Hour * 24
}

func (t UnixTimeDay) String() string {
	ratio, divide := RatioFrom(t.unit(), time.Second)
	if divide {
		return fmt.Sprintf("%d", t.Unix()/int64(ratio))
	}
	return fmt.Sprintf("%d", t.Unix()*int64(ratio))
}

func (t UnixTimeDay) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *UnixTimeDay) UnmarshalJSON(data []byte) error {
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
func (t UnixTimeDay) MarshalText() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (t *UnixTimeDay) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}
