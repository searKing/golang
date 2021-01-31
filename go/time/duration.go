package time

import (
	"math/rand"
	"time"
)

// InfDuration is the duration returned by Delay when a Reservation is not OK.
const InfDuration = time.Duration(1<<63 - 1)

// Jitter returns a time.Duration between
// [duration - maxFactor*duration, duration + maxFactor*duration].
//
// This allows clients to avoid converging on periodic behavior.
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	delta := time.Duration(maxFactor) * duration
	minInterval := duration - delta
	maxInterval := duration + delta
	// Get a random value from the range [minInterval, maxInterval].
	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	return minInterval + time.Duration(rand.Float64()*float64(maxInterval-minInterval+1))
}
