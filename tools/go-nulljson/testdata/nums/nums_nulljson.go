// Code generated by "go-nulljson -type Nums<*string>"; DO NOT EDIT.

// Install go-nulljson by `go install github.com/searKing/golang/tools/go-nulljson@latest`
//
// Deprecated: Use [github.com/searKing/golang/go/exp/database/sql.NullJson[T]] instead.
// For more information, see:
// https://github.com/searKing/golang/blob/master/go/exp/database/sql/null_json.go
package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
)

func _() {
	// An "cannot convert Nums literal (type Nums) to type atomic.Value" compiler error signifies that the base type have changed.
	// Re-run the go-nulljson command to generate them again.
	var val Nums
	_ = (sql.Scanner)(&val)
	_ = (driver.Valuer)(&val)
}

// Nums represents an interface that may be null.
// Nums implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.
type Nums struct {
	Data *string

	Valid bool // Valid is true if Data is not NULL
}

// Scan implements the sql.Scanner interface.
func (nj *Nums) Scan(src any) error {
	if src == nil {
		nj.Data, nj.Valid = nil, false
		return nil
	}
	nj.Valid = true

	var err error
	switch src := src.(type) {
	case string:
		if len(src) > 0 {
			var v any = &nj.Data
			switch v := v.(type) {
			default:
				err = json.Unmarshal([]byte(src), v)
			}
		}
	case []byte:
		if len(src) > 0 {
			var v any = &nj.Data
			switch v := v.(type) {
			default:
				err = json.Unmarshal(src, v)
			}
		}
	case time.Time:
		srcBytes, _ := json.Marshal(src)
		var v any = &nj.Data
		switch v := v.(type) {
		case proto.Message:
		default:
			err = json.Unmarshal(srcBytes, v)
		}
	case nil:
		nj.Data = nil
		err = nil
	default:
		srcBytes, _ := json.Marshal(src)
		var v any = &nj.Data
		switch v := v.(type) {
		default:
			err = json.Unmarshal(srcBytes, v)
		}
	}
	if err == nil {
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T: %w", src, nj.Data, err)
}

// Value implements the driver.Valuer interface.
func (nj Nums) Value() (driver.Value, error) {
	if !nj.Valid {
		return nil, nil
	}
	var v any = &nj.Data
	switch v := v.(type) {
	default:
		return json.Marshal(v)
	}
}
