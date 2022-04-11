// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import "net/http"

// NamedCheck returns a healthz checker for the given name and function.
func NamedCheck(name string, check func(r *http.Request) error) HealthCheck {
	return &namedHealthChecker{name, check}
}

// namedHealthChecker implements HealthCheck on an arbitrary name and check function.
type namedHealthChecker struct {
	name  string
	check func(r *http.Request) error
}

var _ HealthCheck = &namedHealthChecker{}

func (c *namedHealthChecker) Name() string {
	return c.name
}

func (c *namedHealthChecker) Check(r *http.Request) error {
	return c.check(r)
}
