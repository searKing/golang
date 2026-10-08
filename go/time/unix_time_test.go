// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"encoding/json"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

type unixTime interface {
	json.Marshaler
	String() string
}

func TestUnixTimeMarshalJSON(t *testing.T) {
	const days = 20000 // 2024-10-04
	tm := time.Unix(days*24*60*60, 0).UTC()
	tests := []struct {
		t    unixTime
		ptr  json.Unmarshaler
		want string
	}{
		{time_.UnixTimeNanosecond{Time: tm}, &time_.UnixTimeNanosecond{}, "1728000000000000000"},
		{time_.UnixTimeMicrosecond{Time: tm}, &time_.UnixTimeMicrosecond{}, "1728000000000000"},
		{time_.UnixTimeMillisecond{Time: tm}, &time_.UnixTimeMillisecond{}, "1728000000000"},
		{time_.UnixSecondTime{Time: tm}, &time_.UnixSecondTime{}, "1728000000"},
		{time_.UnixTimeMinute{Time: tm}, &time_.UnixTimeMinute{}, "28800000"},
		{time_.UnixTimeHour{Time: tm}, &time_.UnixTimeHour{}, "480000"},
		{time_.UnixTimeDay{Time: tm}, &time_.UnixTimeDay{}, "20000"},
	}
	for _, tt := range tests {
		data, err := json.Marshal(tt.t)
		if err != nil {
			t.Fatalf("json.Marshal(%T) error: %v", tt.t, err)
		}
		if got := string(data); got != tt.want {
			t.Errorf("json.Marshal(%T) = %s; want %s", tt.t, got, tt.want)
		}
		if got := tt.t.String(); got != tt.want {
			t.Errorf("%T.String() = %s; want %s", tt.t, got, tt.want)
		}
		if err := json.Unmarshal(data, tt.ptr); err != nil {
			t.Fatalf("json.Unmarshal(%s, %T) error: %v", data, tt.ptr, err)
		}
		if got, want := tt.ptr.(unixTime).String(), tt.want; got != want {
			t.Errorf("json.Unmarshal(%s, %T) = %s; want %s", data, tt.ptr, got, want)
		}
	}
}
