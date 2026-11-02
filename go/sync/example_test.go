// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"fmt"
	"time"

	"github.com/searKing/golang/go/sync"
)

func ExampleWalk() {

	// chan WalkInfo
	walkChan := make(chan any, 0)

	p := sync.Walk{
		Burst: 1,
	}
	defer p.Wait()

	p.Walk(context.Background(), walkChan, func(name any) error {
		fmt.Printf("%s\n", name)
		return nil
	})

	for i := range 5 {
		walkChan <- fmt.Sprintf("%d", i)
	}
	close(walkChan)
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleThread() {
	var thread sync.Thread
	defer thread.Shutdown(context.Background())

	// Every f runs on the same OS thread, one at a time,
	// so calls depending on per-thread state can be split across calls to Do.
	for _, step := range []string{"init", "use", "release"} {
		err := thread.Do(context.Background(), func(ctx context.Context) {
			fmt.Println(step)
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// Output:
	// init
	// use
	// release
}

func ExampleThread_Do() {
	var thread sync.Thread
	started := make(chan struct{})
	done := make(chan error)
	go func() {
		done <- thread.Do(context.Background(), func(ctx context.Context) {
			close(started)
			// A long-running f returns early when the Thread shuts down.
			<-ctx.Done()
			fmt.Println("f:", context.Cause(ctx))
		})
	}()
	<-started

	if err := thread.Shutdown(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Do:", <-done)
	// Output:
	// f: sync: Thread closed
	// Do: <nil>
}

func ExampleThread_Shutdown() {
	var thread sync.Thread
	started := make(chan struct{})
	release := make(chan struct{})
	go thread.Do(context.Background(), func(ctx context.Context) {
		close(started)
		<-release
	})
	<-started

	// Shutdown waits for the running f, which is still blocked.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	fmt.Println(thread.Shutdown(ctx))

	// The Thread no longer accepts f, even though it has not exited yet.
	fmt.Println(thread.Do(context.Background(), func(ctx context.Context) {}))

	close(release)
	fmt.Println(thread.Shutdown(context.Background()))
	// Output:
	// context deadline exceeded
	// sync: Thread closed
	// <nil>
}

func ExampleThread_Stats() {
	var thread sync.Thread
	defer thread.Shutdown(context.Background())

	// To publish the stats, use expvar:
	//
	//	expvar.Publish("thread", expvar.Func(func() any { return thread.Stats() }))
	_ = thread.Do(context.Background(), func(ctx context.Context) {})
	stats := thread.Stats()
	fmt.Println(stats.Running, stats.Active, stats.Calls)
	// Output:
	// true 0 1
}
