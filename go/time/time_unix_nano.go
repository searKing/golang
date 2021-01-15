package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type UnixNanoTime time.Time

func (t UnixNanoTime) Time() time.Time {
	return time.Time(t)
}

func (t UnixNanoTime) String() string {
	return fmt.Sprintf("%d", time.Time(t).UnixNano())
}

func (t UnixNanoTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UnixNano())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *UnixNanoTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	*t = UnixNanoTime(time.Unix(0, timestamp))
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in Unix Seconds, with sub-second precision added if present.
func (t UnixNanoTime) MarshalText() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (t *UnixNanoTime) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (t UnixNanoTime) MarshalBinary() ([]byte, error) {
	return time.Time(t).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *UnixNanoTime) UnmarshalBinary(data []byte) error {
	return ((*time.Time)(t)).UnmarshalBinary(data)
}
