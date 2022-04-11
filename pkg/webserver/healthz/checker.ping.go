// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import "net/http"

// PingHealthCheck returns true automatically when checked
var PingHealthCheck HealthCheck = ping{}

// ping implements the simplest possible healthz checker.
type ping struct{}

func (ping) Name() string {
	return "ping"
}

// Check is a health check that returns true.
func (ping) Check(_ *http.Request) error {
	return nil
}
