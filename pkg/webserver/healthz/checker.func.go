// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import (
	"net/http"

	"github.com/searKing/golang/go/exp/maps"
)

// healthzCheck implements HealthChecker on an arbitrary name and check function.
type healthzCheck struct {
	name  string
	check func(r *http.Request) error
}

var _ HealthChecker = &healthzCheck{}

func (c *healthzCheck) Name() string {
	return c.name
}

func (c *healthzCheck) Check(r *http.Request) error {
	return c.check(r)
}

// getExcludedChecks extracts the health check names to be excluded from the query param
func getExcludedChecks(r *http.Request) map[string]struct{} {
	checks, found := r.URL.Query()["exclude"]
	if found {
		return maps.Set(checks...)
	}
	return nil
}
