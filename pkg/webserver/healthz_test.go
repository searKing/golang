// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"testing"
	"time"
)

func TestDelayedHealthCheck(t *testing.T) {
	t.Run("test that liveness check returns true until the delay has elapsed", func(t *testing.T) {
		doneCh := make(chan struct{})

		healthCheck := delayedHealthCheck(postStartHookHealthz{name: "test", done: doneCh}, 10*time.Millisecond)
		err := healthCheck.Check(nil)
		if err != nil {
			t.Errorf("Got %v, expected no error", err)
		}
		err = healthCheck.Check(nil)
		if err != nil {
			t.Errorf("Got %v, expected no error", err)
		}
		time.Sleep(15 * time.Millisecond)
		err = healthCheck.Check(nil)
		if err == nil || err.Error() != "not finished" {
			t.Errorf("Got '%v', but expected error to be 'not finished'", err)
		}
		close(doneCh)
		err = healthCheck.Check(nil)
		if err != nil {
			t.Errorf("Got %v, expected no error", err)
		}
	})
	t.Run("test that liveness check does not toggle false even if done channel is closed early", func(t *testing.T) {
		doneCh := make(chan struct{})

		healthCheck := delayedHealthCheck(postStartHookHealthz{name: "test", done: doneCh}, 10*time.Millisecond)
		err := healthCheck.Check(nil)
		if err != nil {
			t.Errorf("Got %v, expected no error", err)
		}
		close(doneCh)
		err = healthCheck.Check(nil)
		if err != nil {
			t.Errorf("Got %v, expected no error", err)
		}
		time.Sleep(15 * time.Millisecond)
		err = healthCheck.Check(nil)
		if err != nil {
			t.Errorf("Got %v, expected no error", err)
		}
	})

}
