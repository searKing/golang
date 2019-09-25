package sql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// NullJson represents an interface that may be null.
// NullJson implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullJson struct {
	Data interface{} // must be set with a pointer to zero value of expect type

	Valid bool // Valid is true if Data is not NULL
}

// Scan implements the sql.Scanner interface.
func (nj *NullJson) Scan(value interface{}) error {
	if value == nil {
		nj.Data, nj.Valid = nil, false
		return nil
	}
	nj.Valid = true

	var err error
	switch src := value.(type) {
	case string:
		err = json.Unmarshal([]byte(src), &nj.Data)
	case []byte:
		err = json.Unmarshal(src, &nj.Data)
	case time.Time:
		srcBytes, _ := json.Marshal(src)
		err = json.Unmarshal(srcBytes, &nj.Data)
	case nil:
		nj.Data = nil
		err = nil
	default:
		srcBytes, _ := json.Marshal(src)
		err = json.Unmarshal(srcBytes, &nj.Data)
	}
	if err == nil {
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T : %w", value, nj.Data, err)
}

// Value implements the driver.Valuer interface.
func (nj NullJson) Value() (driver.Value, error) {
	if !nj.Valid {
		return nil, nil
	}
	return json.Marshal(nj.Data)
}
