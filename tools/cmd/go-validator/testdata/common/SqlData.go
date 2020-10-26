// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"time"
)

type SqlData struct {
	Id        uint      `db:"id" json:"sql_data_id,omitempty" validate:"gte=0,lte=130"`
	CreatedAt time.Time `db:"created_at" json:"sql_data_created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"sql_data_updated_at,omitempty"`

	IsDeleted bool         `json:"sql_data_is_deleted,omitempty" db:"is_deleted"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"sql_data_deleted_at,omitempty"`

	Version uint `db:"version" json:"sql_data_version,omitempty"`
} // sql_data

type JsonData struct {
	Id uint `db:"id" json:"sql_data_id,omitempty" validate:"gte=0,lte=130"`
} // sql_data

func main() {

}
