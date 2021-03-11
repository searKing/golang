// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import "time"

// TruncateByLocation only happens in local semantics, apparently.
// observed values for truncating given time with 24 Hour:
// before truncation: 2012-01-01 04:15:30 +0800 CST
// after  truncation: 2012-01-01 00:00:00 +0800 CST
//
// time.Truncate only happens in UTC semantics, apparently.
// observed values for truncating given time with 24 Hour:
//
// before truncation: 2012-01-01 04:15:30 +0800 CST
// after  truncation: 2011-12-31 08:00:00 +0800 CST
//
// This is really annoying when we want to truncate in local time
// we take the apparent local time in the local zone, and pretend
// that it's in UTC. do our math, and put it back to the local zone
func TruncateByLocation(t time.Time, d time.Duration) time.Time {
	if t.Location() == time.UTC {
		return t.Truncate(d)
	}

	utc := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	utc = utc.Truncate(d)
	return time.Date(utc.Year(), utc.Month(), utc.Day(), utc.Hour(), utc.Minute(), utc.Second(), utc.Nanosecond(),
		t.Location())
}
