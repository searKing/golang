// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goroutine_test

import (
	"fmt"
	"sync"

	"github.com/searKing/golang/go/runtime/goroutine"
)

func ExampleID() {
	fmt.Printf("%d\n", goroutine.ID())
	// Output:
	// 1
}

func ExampleNewLock() {
	oldDebug := goroutine.DebugGoroutines
	goroutine.DebugGoroutines = true
	defer func() { goroutine.DebugGoroutines = oldDebug }()

	g := goroutine.NewLock()
	g.MustCheck()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic recovered: %v\n", r)
			}
		}()
		g.MustCheck() // should panic
	}()
	wg.Wait()
	// Output:
	// panic recovered: running on the wrong goroutine
}
