// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"
	"testing"
	"time"

	sync_ "github.com/searKing/golang/go/sync"
)

func TestThread_OSThreadTerminated(t *testing.T) {
	var thread sync_.Thread
	var tid int
	if err := thread.Do(t.Context(), func() { tid = syscall.Gettid() }); err != nil {
		t.Fatalf("Do = %v, want nil", err)
	}
	if err := thread.ShutdownContext(t.Context()); err != nil {
		t.Fatalf("ShutdownContext = %v, want nil", err)
	}

	// The runtime terminates the OS thread shortly after the worker goroutine exits.
	task := fmt.Sprintf("/proc/self/task/%d", tid)
	for deadline := time.Now().Add(5 * time.Second); ; time.Sleep(time.Millisecond) {
		_, err := os.Stat(task)
		if errors.Is(err, fs.ErrNotExist) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("OS thread %d of the Thread is not terminated after Shutdown: stat %s: %v", tid, task, err)
		}
	}
}
