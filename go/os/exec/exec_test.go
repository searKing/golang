// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

package exec_test

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"testing"
	"time"

	exec_ "github.com/searKing/golang/go/os/exec"
)

func TestCommandWithTimeout(t *testing.T) {
	var got string
	err := exec_.CommandWithTimeout(func(r io.Reader) {
		b, _ := io.ReadAll(r)
		got = string(b)
	}, time.Second, "echo", "hi")
	if err != nil {
		t.Fatalf("CommandWithTimeout() error = %v", err)
	}
	if got != "hi\n" {
		t.Errorf("stdout = %q, want %q", got, "hi\n")
	}
}

func TestCommandWithTimeout_Timeout(t *testing.T) {
	start := time.Now()
	err := exec_.CommandWithTimeout(nil, 100*time.Millisecond, "sleep", "10")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("CommandWithTimeout() error = %v, want %v", err, context.DeadlineExceeded)
	}
	if cost := time.Since(start); cost > 5*time.Second {
		t.Errorf("CommandWithTimeout() returned after %v, process not killed on timeout", cost)
	}
}

func TestCommandWithTimeout_SIGTERM(t *testing.T) {
	var got string
	start := time.Now()
	// sleep in background holds stdout after sh exits
	err := exec_.CommandWithTimeout(func(r io.Reader) {
		b, _ := io.ReadAll(r)
		got = string(b)
	}, 200*time.Millisecond, "sh", "-c", `trap 'echo term; exit 0' TERM; sleep 10 & wait`)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("CommandWithTimeout() error = %v, want %v", err, context.DeadlineExceeded)
	}
	if got != "term\n" {
		t.Errorf("stdout = %q, want %q", got, "term\n")
	}
	if cost := time.Since(start); cost > 5*time.Second {
		t.Errorf("CommandWithTimeout() returned after %v, stdout not closed after WaitDelay", cost)
	}
}

func TestCommandWithTimeout_IgnoreSIGTERM(t *testing.T) {
	start := time.Now()
	err := exec_.CommandWithTimeout(nil, 200*time.Millisecond, "sh", "-c", `trap '' TERM; sleep 10`)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("CommandWithTimeout() error = %v, want %v", err, context.DeadlineExceeded)
	}
	if cost := time.Since(start); cost > 5*time.Second {
		t.Errorf("CommandWithTimeout() returned after %v, process not killed after WaitDelay", cost)
	}
}

func TestCommandWithCancel(t *testing.T) {
	var got string
	err := exec_.CommandWithCancel(func(r io.Reader) {
		b, _ := io.ReadAll(r)
		got = string(b)
	}, "echo", "hi")
	if err != nil {
		t.Fatalf("CommandWithCancel() error = %v", err)
	}
	if got != "hi\n" {
		t.Errorf("stdout = %q, want %q", got, "hi\n")
	}

	err = exec_.CommandWithCancel(nil, "false")
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("CommandWithCancel() error = %v, want *exec.ExitError", err)
	}
}

func TestCommandWithTimeoutHandler(t *testing.T) {
	data, err := exec_.CommandWithTimeoutHandler(time.Second, "sh", "-c", "echo out; echo err >&2")
	if err != nil {
		t.Fatalf("CommandWithTimeoutHandler() error = %v", err)
	}
	if string(data) != "out\nerr\n" {
		t.Errorf("output = %q, want %q", data, "out\nerr\n")
	}

	_, err = exec_.CommandWithTimeoutHandler(time.Second, "false")
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("CommandWithTimeoutHandler() error = %v, want *exec.ExitError", err)
	}

	start := time.Now()
	data, err = exec_.CommandWithTimeoutHandler(100*time.Millisecond, "sleep", "10")
	if !errors.Is(err, context.DeadlineExceeded) || data != nil {
		t.Errorf("CommandWithTimeoutHandler() = %q, %v, want nil, %v", data, err, context.DeadlineExceeded)
	}
	if cost := time.Since(start); cost > 5*time.Second {
		t.Errorf("CommandWithTimeoutHandler() returned after %v, process not killed on timeout", cost)
	}
}

func TestCommandNotFound(t *testing.T) {
	if err := exec_.CommandWithCancel(nil, "searking-command-not-found"); err == nil {
		t.Error("CommandWithCancel() error = nil, want not found")
	}
}
