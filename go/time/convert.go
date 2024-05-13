// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import "time"

// RatioFrom Example:
// ratio, divide := RatioFrom(fromUnit, toUnit)
//
//	if divide {
//	  toDuration = fromDuration / ratio
//	  fromDuration = toDuration * ratio
//	} else {
//
//	  toDuration = fromDuration * ratio
//	  fromDuration = toDuration / ratio
//	}
func RatioFrom(from time.Duration, to time.Duration) (ratio time.Duration, divide bool) {
	if from >= to {
		return from / to, false
	}

	return to / from, true
}

// ConvertTimestamp convert timestamp from one unit to another unit
func ConvertTimestamp(timestamp int64, from, to time.Duration) int64 {
	ratio, divide := RatioFrom(from, to)
	if divide {
		return int64(time.Duration(timestamp) / ratio)
	}
	return int64(time.Duration(timestamp) * ratio)
}

// Timestamp returns a Unix timestamp in unit from "January 1, 1970 UTC".
// The result is undefined if the Unix time cannot be represented by an int64.
// Which includes calling TimeUnixMilli on a zero Time is undefined.
//
// See Go stdlib https://golang.org/pkg/time/#Time.UnixNano for more information.
func Timestamp(t time.Time, unit time.Duration) int64 {
	return ConvertTimestamp(t.UnixNano(), time.Nanosecond, unit)
}

// UnixWithUnit returns the local Time corresponding to the given Unix time by unit
func UnixWithUnit(timestamp int64, unit time.Duration) time.Time {
	sec := ConvertTimestamp(timestamp, unit, time.Second)
	nsec := time.Duration(timestamp)*unit - time.Duration(sec)*time.Second
	return time.Unix(sec, int64(nsec))
}
