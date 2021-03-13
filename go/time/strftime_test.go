// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"strings"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

// http://cpp.sh/6gqrg
func TestLayoutTimeToSimilarStrftime(t *testing.T) {
	tests := []struct {
		layout     string
		wantLayout string
	}{
		{
			layout:     time.ANSIC,             // Fri Mar 12 19:49:05 2021
			wantLayout: "%a %b %e %H:%M:%S %Y", // Fri Mar 12 09:42:04 2021
		},
		{
			layout:     time.UnixDate,             // Fri Mar 12 19:49:05 CST 2021
			wantLayout: "%a %b %e %H:%M:%S %Z %Y", // Fri Mar 12 11:57:19 UTC 2021
		},
		{
			layout:     time.RubyDate,             // Fri Mar 12 19:49:05 +0800 2021
			wantLayout: "%a %b %d %H:%M:%S %z %Y", // Fri Mar 12 11:57:19 +0000 2021
		},
		{
			layout:     time.RFC822,         // 12 Mar 21 19:49 CST
			wantLayout: "%d %b %y %H:%M %Z", // 12 Mar 21 11:57 UTC
		},
		{
			layout:     time.RFC822Z,        // 12 Mar 21 19:49 +0800
			wantLayout: "%d %b %y %H:%M %z", // 12 Mar 21 11:57 +0000
		},
		{
			layout:     time.RFC850,                // Friday, 12-Mar-21 19:49:05 CST
			wantLayout: "%A, %d-%b-%y %H:%M:%S %Z", // Friday, 12-Mar-21 11:57:19 UTC
		},
		{
			layout:     time.RFC1123,               // Fri, 12 Mar 2021 19:49:05 CST
			wantLayout: "%a, %d %b %Y %H:%M:%S %Z", // Fri, 12 Mar 2021 11:57:19 UTC
		},
		{
			layout:     time.RFC1123Z,              // Fri, 12 Mar 2021 19:49:05 +0800
			wantLayout: "%a, %d %b %Y %H:%M:%S %z", // Fri, 12 Mar 2021 11:57:19 +0000
		},
		{
			layout:     time.RFC3339,          // 2021-03-12T19:49:05+08:00
			wantLayout: "%Y-%m-%dT%H:%M:%S%z", // 2021-03-12T17:25:56+0000
		},
		{
			layout:     time.RFC3339Nano,      // 2021-03-12T19:49:05.236589+08:00
			wantLayout: "%Y-%m-%dT%H:%M:%S%z", // 2021-03-12T11:57:19Z07:00
		},
		{
			layout:     time.Kitchen, // 7:49PM
			wantLayout: "%I:%M%p",    // 05:27PM
		},
		{
			layout:     time.Stamp,       // Mar 12 19:49:05
			wantLayout: "%b %e %H:%M:%S", // Mar 12 11:57:19
		},
		{
			layout:     time.StampMilli,  // Mar 12 19:49:05.236
			wantLayout: "%b %e %H:%M:%S", // Mar 12 11:57:19
		},
		{
			layout:     time.StampMicro,  // Mar 12 19:49:05.236607
			wantLayout: "%b %e %H:%M:%S", // Mar 12 11:57:19
		},
		{
			layout:     time.StampNano,   // Mar 12 19:49:05.236611000
			wantLayout: "%b %e %H:%M:%S", // Mar 12 11:57:19
		},
	}
	for i, tt := range tests {
		gotLayout := time_.LayoutTimeToSimilarStrftime(tt.layout)
		if tt.wantLayout != gotLayout {
			t.Logf("layout: %s, time: %s", tt.layout, time.Now().Format(tt.layout))
			t.Errorf("#%d: LayoutTimeToStrftime expected %q got %q", i, tt.wantLayout, gotLayout)
		}
	}
}

func TestLayoutStrftimeToSimilarTime(t *testing.T) {
	tests := []struct {
		layout     string
		wantLayout string
	}{
		{
			layout:     "%a %b %e %H:%M:%S %Y", // Fri Mar 12 09:42:04 2021
			wantLayout: time.ANSIC,             // Fri Mar 12 19:49:05 2021
		},
		{
			layout:     "%a %b %e %H:%M:%S %Z %Y", // Fri Mar 12 11:57:19 UTC 2021
			wantLayout: time.UnixDate,             // Fri Mar 12 19:49:05 CST 2021
		},
		{
			layout:     "%a %b %d %H:%M:%S %z %Y",                            // Fri Mar 12 11:57:19 +0000 2021
			wantLayout: strings.ReplaceAll(time.RubyDate, "-0700", "Z07:00"), // Fri Mar 12 19:49:05 +0800 2021
		},
		{
			layout:     "%d %b %y %H:%M %Z", // 12 Mar 21 11:57 UTC
			wantLayout: time.RFC822,         // 12 Mar 21 19:49 CST
		},
		{
			layout:     "%d %b %y %H:%M %z",                                 // 12 Mar 21 11:57 +0000
			wantLayout: strings.ReplaceAll(time.RFC822Z, "-0700", "Z07:00"), // 12 Mar 21 19:49 +0800
		},
		{
			layout:     "%A, %d-%b-%y %H:%M:%S %Z", // Friday, 12-Mar-21 11:57:19 UTC
			wantLayout: time.RFC850,                // Friday, 12-Mar-21 19:49:05 CST
		},
		{
			layout:     "%a, %d %b %Y %H:%M:%S %Z", // Fri, 12 Mar 2021 11:57:19 UTC
			wantLayout: time.RFC1123,               // Fri, 12 Mar 2021 19:49:05 CST
		},
		{
			layout:     "%a, %d %b %Y %H:%M:%S %z",                           // Fri, 12 Mar 2021 11:57:19 +0000
			wantLayout: strings.ReplaceAll(time.RFC1123Z, "-0700", "Z07:00"), // Fri, 12 Mar 2021 19:49:05 +0800
		},
		{
			layout:     "%Y-%m-%dT%H:%M:%S%z", // 2021-03-12T17:25:56+0000
			wantLayout: time.RFC3339,          // 2021-03-12T19:49:05+08:00
		},
		{
			layout:     "%Y-%m-%dT%H:%M:%S%z",                                  // 2021-03-12T11:57:19Z07:00
			wantLayout: strings.ReplaceAll(time.RFC3339Nano, ".999999999", ""), // 2021-03-12T19:49:05.236589+08:00
		},
		{
			layout:     "%I:%M%p",                                   // 05:27PM
			wantLayout: strings.ReplaceAll(time.Kitchen, "3", "03"), // 7:49PM
		},
		{
			layout:     "%b %e %H:%M:%S", // Mar 12 11:57:19
			wantLayout: time.Stamp,       // Mar 12 19:49:05
		},
		{
			layout:     "%b %e %H:%M:%S",                                    // Mar 12 11:57:19
			wantLayout: strings.ReplaceAll(time.StampMilli, "05.000", "05"), // Mar 12 19:49:05.236
		},
		{
			layout:     "%b %e %H:%M:%S",                                       // Mar 12 11:57:19
			wantLayout: strings.ReplaceAll(time.StampMicro, "05.000000", "05"), // Mar 12 19:49:05.236607
		},
		{
			layout:     "%b %e %H:%M:%S",                                         // Mar 12 11:57:19
			wantLayout: strings.ReplaceAll(time.StampNano, "05.000000000", "05"), // Mar 12 19:49:05.236611000
		},
	}
	for i, tt := range tests {
		gotLayout := time_.LayoutStrftimeToSimilarTime(tt.layout)
		if tt.wantLayout != gotLayout {
			t.Logf("layout: %s, time: %s", tt.layout, time.Now().Format(tt.wantLayout))
			t.Errorf("#%d: LayoutTimeToStrftime expected %q got %q", i, tt.wantLayout, gotLayout)
		}
	}

}
