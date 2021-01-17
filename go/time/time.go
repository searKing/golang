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

// UnixWithUnit returns the local Time corresponding to the given Unix time by unit
func UnixWithUnit(timestamp int64, unit time.Duration) time.Time {
	var sec, nsec int64
	ratio, divide := RatioFrom(unit, time.Second)
	if divide {
		sec = int64(time.Duration(timestamp) / ratio)
	} else {
		sec = int64(time.Duration(timestamp) * ratio)
	}
	nsec = int64((time.Duration(timestamp) - time.Duration(sec)*ratio) * unit / time.Nanosecond)
	return time.Unix(sec, nsec)
}
