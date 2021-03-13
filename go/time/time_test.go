// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
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
