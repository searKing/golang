package time

import "time"

// Example:
// ratio, divide := RatioFrom(fromUnit, toUnit)
// if divide {
//   toDuration = fromDuration / ratio
//   fromDuration = toDuration * ratio
// } else {
//   toDuration = fromDuration * ratio
//   fromDuration = toDuration / ratio
// }
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

// Timestamp convert time to timestamp reprensent in unit
func Timestamp(t time.Time, unit time.Duration) int64 {
	return ConvertTimestamp(t.UnixNano(), time.Nanosecond, unit)
}

// UnixWithUnit returns the local Time corresponding to the given Unix time by unit
func UnixWithUnit(timestamp int64, unit time.Duration) time.Time {
	sec := ConvertTimestamp(timestamp, unit, time.Second)
	nsec := time.Duration(timestamp)*unit - time.Duration(sec)*time.Second
	return time.Unix(sec, int64(nsec))
}
