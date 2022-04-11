// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import (
	"net/http"
)

// HealthCheck is a named healthz checker.
type HealthCheck interface {
	Name() string
	Check(req *http.Request) error
}
