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

func (e *testErrorfer) Errorf(format string, args ...interface{}) {
	e.errors = append(e.errors, fmt.Sprintf(format, args...))
	e.errorCount++
}

func TestCheck(t *testing.T) {
	const leakCount = 3
	for i := 0; i < leakCount; i++ {
		go func() { time.Sleep(2 * time.Second) }()
	}
	if ig := interestingGoroutines(); len(ig) == 0 {
		t.Error("blah")
	}
	e := &testErrorfer{}
	check(e, time.Second)
	if e.errorCount != leakCount {
		t.Errorf("check found %v leaks, want %v leaks", e.errorCount, leakCount)
		t.Logf("leaked goroutines:\n%v", strings.Join(e.errors, "\n"))
	}
	check(t, 3*time.Second)
}

func ignoredTestingLeak(d time.Duration) {
	time.Sleep(d)
}

func TestCheckRegisterIgnore(t *testing.T) {
	RegisterIgnoreGoroutine("ignoredTestingLeak")
	const leakCount = 3
	for i := 0; i < leakCount; i++ {
		go func() { time.Sleep(2 * time.Second) }()
	}
	go func() { ignoredTestingLeak(3 * time.Second) }()
	if ig := interestingGoroutines(); len(ig) == 0 {
		t.Error("blah")
	}
	e := &testErrorfer{}
	check(e, time.Second)
	if e.errorCount != leakCount {
		t.Errorf("check found %v leaks, want %v leaks", e.errorCount, leakCount)
		t.Logf("leaked goroutines:\n%v", strings.Join(e.errors, "\n"))
	}
	check(t, 3*time.Second)
}
