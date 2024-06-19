// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/searKing/golang/go/runtime"
	time_ "github.com/searKing/golang/go/time"
)

// LogHealthCheck returns true if logging is not blocked
var LogHealthCheck HealthChecker = &log{}

type log struct {
	startOnce    sync.Once
	lastVerified atomic.Value
}

func (l *log) Name() string { return "log" }

func (l *log) Check(_ *http.Request) error {
	l.startOnce.Do(func() {
		l.lastVerified.Store(time.Now())
		go time_.Forever(func() {
			defer runtime.NeverPanicButLog.Recover()
			l.lastVerified.Store(time.Now())
		}, time.Minute)
	})

	lastVerified := l.lastVerified.Load().(time.Time)
	if time.Since(lastVerified) < (2 * time.Minute) {
		return nil
	}
	return fmt.Errorf("logging blocked")
}
