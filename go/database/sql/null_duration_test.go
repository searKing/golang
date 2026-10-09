// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql_test

import (
	"testing"
	"time"

	"github.com/searKing/golang/go/database/sql"
)

func TestNullDurationScan(t *testing.T) {
	tests := []struct {
		src     any
		want    sql.NullDuration
		wantErr bool
	}{
		{nil, sql.NullDuration{}, false},
		{"1.5s", sql.NullDuration{Data: 1500 * time.Millisecond, Valid: true}, false},
		{[]byte("2s"), sql.NullDuration{Data: 2 * time.Second, Valid: true}, false},
		{time.Second, sql.NullDuration{Data: time.Second, Valid: true}, false},
		{int64(time.Minute), sql.NullDuration{Data: time.Minute, Valid: true}, false},
		{"bad", sql.NullDuration{}, true},
		{1.5, sql.NullDuration{}, true},
	}
	for _, tt := range tests {
		got := sql.NullDuration{Data: time.Hour, Valid: true}
		err := got.Scan(tt.src)
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("Scan(%T(%v)) = %+v, %v; want %+v, error %t", tt.src, tt.src, got, err, tt.want, tt.wantErr)
		}
	}
}
