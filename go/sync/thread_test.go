// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"testing"

	sync_ "github.com/searKing/golang/go/sync"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(thread *sync_.Thread, o *one, c chan bool) {
	thread.Do(func() { o.Increment() })
	c <- true
}

func TestThread(t *testing.T) {
	o := new(one)
	thread := new(sync_.Thread)
	defer thread.Shutdown()
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(thread, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != N {
		t.Errorf("once failed outside run: %d is not %d", *o, N)
	}
}

func TestThreadPanic(t *testing.T) {
	var thread sync_.Thread
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Thread.Do did not panic")
			}
		}()
		thread.Do(func() {
			panic("failed")
		})
	}()

	{
		var do bool
		thread.Do(func() {
			do = true
		})
		if !do {
			t.Fatalf("Thread.Do did not called")
		}
	}
	thread.Shutdown()

	{
		var do bool
		thread.Do(func() {
			do = true
		})
		if do {
			t.Fatalf("Thread.Do called after Thread.Shutdown")
		}
	}
}

func BenchmarkThread(b *testing.B) {
	var thread sync_.Thread
	defer thread.Shutdown()
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			thread.Do(f)
		}
	})
}
