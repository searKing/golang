package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixTime time.Time

func (t UnixTime) Time() time.Time {
	return time.Time(t)
}

func (t UnixTime) String() string {
	return fmt.Sprintf("%d", time.Time(t).Unix())
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Unix())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *UnixTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	*t = UnixTime(time.Unix(timestamp, 0))
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in Unix Seconds, with sub-second precision added if present.
func (t UnixTime) MarshalText() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (t *UnixTime) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (t UnixTime) MarshalBinary() ([]byte, error) {
	return time.Time(t).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *UnixTime) UnmarshalBinary(data []byte) error {
	return ((*time.Time)(t)).UnmarshalBinary(data)
}
