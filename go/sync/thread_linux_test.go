// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"testing"
	"time"

	sync_ "github.com/searKing/golang/go/sync"
)

// serveMainThreadEnv makes the test binary serve a Thread from function main, see serveMainThread.
const serveMainThreadEnv = "SYNC_TEST_THREAD_SERVE_MAIN_THREAD"

func init() {
	if os.Getenv(serveMainThreadEnv) == "1" {
		runtime.LockOSThread()
	}
}

func TestMain(m *testing.M) {
	if os.Getenv(serveMainThreadEnv) == "1" {
		serveMainThread()
	}
	os.Exit(m.Run())
}

// serveMainThread does what function main does to run f on the main thread:
// TestMain runs on the main goroutine, which init has locked to the main thread.
// It exits the process with status 0 if f ran on the main thread.
func serveMainThread() {
	thread := sync_.Thread{ServeOnly: true}
	var tid int
	go func() {
		defer thread.Close()
		_ = thread.Do(context.Background(), func(context.Context) { tid = syscall.Gettid() })
	}()
	err := thread.Serve()
	fmt.Printf("Serve = %v; f ran on thread %d of process %d\n", err, tid, os.Getpid())
	if !errors.Is(err, sync_.ErrThreadClosed) || tid != os.Getpid() {
		os.Exit(1)
	}
	os.Exit(0)
}

func TestThread_ServeMainThread(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = append(os.Environ(), serveMainThreadEnv+"=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Serve from function main, locked to the main thread by init, did not run f on it: %v\n%s", err, out)
	}
}

func TestThread_OSThreadTerminated(t *testing.T) {
	var thread sync_.Thread
	var tid int
	if err := thread.Do(t.Context(), func(context.Context) { tid = syscall.Gettid() }); err != nil {
		t.Fatalf("Do = %v, want nil", err)
	}
	if tid == os.Getpid() {
		t.Fatalf("f ran on the main thread %d, which the runtime never terminates", tid)
	}
	if err := thread.Shutdown(t.Context()); err != nil {
		t.Fatalf("Shutdown = %v, want nil", err)
	}

	// The runtime terminates the OS thread shortly after the worker goroutine exits.
	task := fmt.Sprintf("/proc/self/task/%d", tid)
	for deadline := time.Now().Add(5 * time.Second); ; time.Sleep(time.Millisecond) {
		_, err := os.Stat(task)
		if errors.Is(err, fs.ErrNotExist) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("OS thread %d of the Thread is not terminated after Shutdown (process %d, main thread: %v): stat %s: %v",
				tid, os.Getpid(), tid == os.Getpid(), task, err)
		}
	}
}
