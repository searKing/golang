// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"encoding/json"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

func TestTruncateByLocation(t *testing.T) {
	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("malformed location: %s: %s", "Asia/Shanghai", err)
	}

	tests := []struct {
		t                  time.Time
		d                  time.Duration
		wantTimeByUTC      time.Time
		wantTimeByLocation time.Time
	}{
		{
			t:                  time.Date(2012, time.January, 1, 4, 15, 30, 0, time.UTC),
			d:                  time.Hour * 24,
			wantTimeByUTC:      time.Date(2012, time.January, 1, 0, 0, 0, 0, time.UTC),
			wantTimeByLocation: time.Date(2012, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			t:                  time.Date(2012, time.January, 1, 4, 15, 30, 0, cst),
			d:                  time.Hour * 24,
			wantTimeByUTC:      time.Date(2011, time.December, 31, 8, 0, 0, 0, cst),
			wantTimeByLocation: time.Date(2012, time.January, 1, 0, 0, 0, 0, cst),
		},
	}

	for i, tt := range tests {
		gotTime := tt.t.Truncate(tt.d)
		if tt.wantTimeByUTC != gotTime {
			t.Errorf("#%d: Truncate expected %q got %q", i, tt.wantTimeByUTC, gotTime)
		}
		gotTime = time_.TruncateByLocation(tt.t, tt.d)
		if tt.wantTimeByLocation != gotTime {
			t.Errorf("#%d: TruncateByLocation expected %q got %q", i, tt.wantTimeByLocation, gotTime)
		}
	}
}

var jsonTests = []struct {
	duration time_.Duration
	json     string
}{
	// simple
	{0, `"0s"`},
	{time_.Duration(5 * time.Second), `"5s"`},
	{time_.Duration(30 * time.Second), `"30s"`},
	{time_.Duration(1478 * time.Second), `"24m38s"`},
	// sign
	{time_.Duration(-5 * time.Second), `"-5s"`},
	{time_.Duration(5 * time.Second), `"5s"`},
	// decimal
	{time_.Duration(5*time.Second + 600*time.Millisecond), `"5.6s"`},
	{time_.Duration(5*time.Second + 4*time.Millisecond), `"5.004s"`},
	{time_.Duration(100*time.Second + 1*time.Millisecond), `"1m40.001s"`},
	// different units
	{time_.Duration(10 * time.Nanosecond), `"10ns"`},
	{time_.Duration(11 * time.Microsecond), `"11µs"`},
	{time_.Duration(13 * time.Millisecond), `"13ms"`},
	{time_.Duration(14 * time.Second), `"14s"`},
	{time_.Duration(15 * time.Minute), `"15m0s"`},
	{time_.Duration(16 * time.Hour), `"16h0m0s"`},
	// composite durations
	{time_.Duration(3*time.Hour + 30*time.Minute), `"3h30m0s"`},
	{time_.Duration(49*time.Minute + 48*time.Second + 372539827*time.Nanosecond), `"49m48.372539827s"`},
}

func TestDurationJSON(t *testing.T) {
	for _, tt := range jsonTests {
		var jsonDuration time_.Duration
		if jsonBytes, err := json.Marshal(tt.duration); err != nil {
			t.Errorf("%v json.Marshal error = %v, want nil", tt.duration, err)
		} else if string(jsonBytes) != tt.json {
			t.Errorf("%v JSON = %#q, want %#q", tt.duration, string(jsonBytes), tt.json)
		} else if err = json.Unmarshal(jsonBytes, &jsonDuration); err != nil {
			t.Errorf("%v json.Unmarshal error = %v, want nil", tt.duration, err)
		}
	}
}

func TestInvalidDurationJSON(t *testing.T) {
	var tt time_.Duration
	err := json.Unmarshal([]byte(`{"now is the time":"buddy"}`), &tt)
	_, isParseErr := err.(*time.ParseError)
	if !isParseErr {
		t.Errorf("expected *time.ParseError unmarshaling JSON, got %v", err)
	}
}
