// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"time"
)

// parseRetryAfter parses the Retry-After header and returns the
// delay duration according to the spec: https://httpwg.org/specs/rfc7231.html#header.retry-after
// The bool returned will be true if the header was successfully parsed.
// Otherwise, the header was either not present, or was not parseable according to the spec.
//
// Retry-After headers come in two flavors: Seconds or HTTP-Date
//
// Examples:
// * Retry-After: Fri, 31 Dec 1999 23:59:59 GMT
// * Retry-After: 120
func parseRetryAfter(text string) (time.Duration, bool) {
	if text == "" {
		return 0, false
	}
	// Retry-After: 120
	// A non-negative decimal integer indicating the seconds to delay after the response is received.
	if delay, err := time.ParseDuration(text + "s"); err == nil {
		if delay < 0 { // a negative sleep doesn't make sense
			return 0, false
		}
		return delay, true
	}

	// Retry-After: Fri, 31 Dec 1999 23:59:59 GMT
	// A date after which to retry.
	retryTime, err := time.Parse(time.RFC1123, text)
	if err != nil {
		return 0, false
	}
	if until := retryTime.Sub(time.Now()); until > 0 {
		return until, true
	}
	// date is in the past
	return 0, true
}
