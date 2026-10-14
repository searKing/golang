// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leakcheck

// it's borrowed from https://github.com/grpc/grpc-go/blob/master/internal/leakcheck/leakcheck_test.go

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type testErrorfer struct {
	errorCount int
	errors     []string
}

func (e *testErrorfer) Errorf(format string, args ...any) {
	e.errors = append(e.errors, fmt.Sprintf(format, args...))
	e.errorCount++
}

func testDurations(t *testing.T) (leak, ignored, detect, cleanup time.Duration) {
	t.Helper()
	if testing.Short() {
		return 200 * time.Millisecond, 300 * time.Millisecond, 100 * time.Millisecond, 500 * time.Millisecond
	}
	return 2 * time.Second, 3 * time.Second, time.Second, 3 * time.Second
}

func TestCheck(t *testing.T) {
	leakDuration, _, detectTimeout, cleanupTimeout := testDurations(t)
	const leakCount = 3
	for i := 0; i < leakCount; i++ {
		go func() { time.Sleep(leakDuration) }()
	}
	if ig := interestingGoroutines(); len(ig) == 0 {
		t.Error("blah")
	}
	e := &testErrorfer{}
	check(e, detectTimeout)
	if e.errorCount != leakCount {
		t.Errorf("check found %v leaks, want %v leaks", e.errorCount, leakCount)
		t.Logf("leaked goroutines:\n%v", strings.Join(e.errors, "\n"))
	}
	check(t, cleanupTimeout)
}

func ignoredTestingLeak(d time.Duration) {
	time.Sleep(d)
}

func TestCheckRegisterIgnore(t *testing.T) {
	RegisterIgnoreGoroutine("ignoredTestingLeak")
	leakDuration, ignoredDuration, detectTimeout, cleanupTimeout := testDurations(t)
	const leakCount = 3
	for i := 0; i < leakCount; i++ {
		go func() { time.Sleep(leakDuration) }()
	}
	go func() { ignoredTestingLeak(ignoredDuration) }()
	if ig := interestingGoroutines(); len(ig) == 0 {
		t.Error("blah")
	}
	e := &testErrorfer{}
	check(e, detectTimeout)
	if e.errorCount != leakCount {
		t.Errorf("check found %v leaks, want %v leaks", e.errorCount, leakCount)
		t.Logf("leaked goroutines:\n%v", strings.Join(e.errors, "\n"))
	}
	check(t, cleanupTimeout)
}
