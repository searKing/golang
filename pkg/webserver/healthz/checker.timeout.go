// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import (
	"net/http"
	"time"
)

// NamedTimeoutCheck returns a healthz checker for the given name , timeout and function.
func NamedTimeoutCheck(name string, timeout time.Duration, check func(r *http.Request, timeout time.Duration) error) HealthChecker {
	return NamedCheck(name, func(r *http.Request) error {
		return check(r, timeout)
	})
}

// NamedDeadlineCheck returns a healthz checker for the given name, deadline and function.
func NamedDeadlineCheck(name string, d time.Time, check func(r *http.Request, d time.Time) error) HealthChecker {
	return NamedCheck(name, func(r *http.Request) error {
		return check(r, d)
	})
}
