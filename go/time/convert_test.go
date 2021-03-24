// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

func TestConvertTimestamp(t *testing.T) {
	now := time.Now()
	tests := []struct {
		timestamp     int64
		fromUnit      time.Duration
		toUnit        time.Duration
		wantTimestamp int64
	}{
		{
			timestamp:     now.Unix(),
			fromUnit:      time.Second,
			toUnit:        time.Nanosecond,
			wantTimestamp: now.Unix() * int64(time.Second/time.Nanosecond),
		},
		{
			timestamp:     now.UnixNano(),
			fromUnit:      time.Nanosecond,
			toUnit:        time.Second,
			wantTimestamp: now.Unix(),
		},
	}
	for i, tt := range tests {
		gotTimestamp := time_.ConvertTimestamp(tt.timestamp, tt.fromUnit, tt.toUnit)
		if tt.wantTimestamp != gotTimestamp {
			t.Errorf("#%d: ConvertTimestamp expected %q got %q", i, tt.wantTimestamp, gotTimestamp)
		}
	}
}

func TestTimestamp(t *testing.T) {
	now := time.Now()
	tests := []struct {
		timestamp     time.Time
		unit          time.Duration
		wantTimestamp int64
	}{
		{
			timestamp:     now,
			unit:          time.Second,
			wantTimestamp: now.Unix(),
		},
		{
			timestamp:     now,
			unit:          time.Millisecond,
			wantTimestamp: now.UnixNano() * int64(time.Nanosecond) / int64(time.Millisecond),
		},
		{
			timestamp:     now,
			unit:          time.Nanosecond,
			wantTimestamp: now.UnixNano(),
		},
		{
			timestamp:     now,
			unit:          time.Second,
			wantTimestamp: now.Unix(),
		},
	}
	for i, tt := range tests {
		gotTimestamp := time_.Timestamp(tt.timestamp, tt.unit)
		if tt.wantTimestamp != gotTimestamp {
			t.Errorf("#%d: ConvertTimestamp expected %d got %d", i, tt.wantTimestamp, gotTimestamp)
		}
	}
}

func TestUnixWithUnit(t *testing.T) {
	now := time.Now()
	tests := []struct {
		timestamp int64
		unit      time.Duration
		wantTime  time.Time
	}{
		{
			timestamp: now.Unix(),
			unit:      time.Second,
			wantTime:  time.Unix(now.Unix(), 0),
		},
		{
			timestamp: now.UnixNano(),
			unit:      time.Nanosecond,
			wantTime:  time.Unix(0, now.UnixNano()),
		},
	}
	for i, tt := range tests {
		gotTimestamp := time_.UnixWithUnit(tt.timestamp, tt.unit)
		if tt.wantTime != gotTimestamp {
			t.Errorf("#%d: ConvertTimestamp expected %q got %q", i, tt.wantTime, gotTimestamp)
		}
	}
}
