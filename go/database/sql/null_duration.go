// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

func _() {
	// An "cannot convert NullDuration literal (type NullDuration) to type atomic.Value" compiler error signifies that the base type have changed.
	// Re-run the go-nulljson command to generate them again.
	_ = (sql.Scanner)(&NullDuration{})
	_ = (driver.Valuer)(&NullDuration{})
}

var nilTimeDurationValue = func() (val time.Duration) { return }()

// NullDuration represents an interface that may be null.
// NullDuration implements the Scanner interface so it can be used as a scan destination, similar to sql.NullString.
type NullDuration struct {
	Data time.Duration

	Valid bool // Valid is true if Data is not NULL
}

// Scan implements the sql.Scanner interface.
func (nj *NullDuration) Scan(src any) error {
	if src == nil {
		nj.Data, nj.Valid = nilTimeDurationValue, false
		return nil
	}
	nj.Valid = true

	var err error
	switch src := src.(type) {
	case string:
		nj.Data, err = time.ParseDuration(src)
	case []byte:
		nj.Data, err = time.ParseDuration(string(src))
	case time.Duration:
		nj.Data, err = time.ParseDuration(src.String())
	case nil:
		nj.Data = nilTimeDurationValue
		err = nil
	default:
		nj.Data, err = time.ParseDuration(fmt.Sprintf("%s", src))
	}
	if err == nil {
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T : %w", src, nj.Data, err)
}

// Value implements the driver.Valuer interface.
func (nj NullDuration) Value() (driver.Value, error) {
	if !nj.Valid {
		return nil, nil
	}
	return nj.Data.String(), nil
}
