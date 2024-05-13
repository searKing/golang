// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goroutine_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/searKing/golang/go/runtime/goroutine"
)

func TestLock(t *testing.T) {
	oldDebug := goroutine.DebugGoroutines
	goroutine.DebugGoroutines = true
	defer func() { goroutine.DebugGoroutines = oldDebug }()

	g := goroutine.NewLock()
	g.MustCheck()

	sawPanic := make(chan any)
	go func() {
		defer func() { sawPanic <- recover() }()
		g.MustCheck() // should panic
	}()
	e := <-sawPanic
	if e == nil {
		t.Fatal("did not see panic from check in other goroutine")
	}
	if !strings.Contains(fmt.Sprint(e), "wrong goroutine") {
		t.Errorf("expected on see panic about running on the wrong goroutine; got %v", e)
	}
}
