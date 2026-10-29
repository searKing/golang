// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"bytes"
	"context"
	"errors"
	"expvar"
	"fmt"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	expvar_ "github.com/searKing/golang/go/expvar"
	sync_ "github.com/searKing/golang/go/sync"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(ctx context.Context, thread *sync_.Thread, o *one, c chan bool) {
	_ = thread.Do(ctx, func() { o.Increment() })
	c <- true
}

func TestThread(t *testing.T) {
	o := new(one)
	thread := new(sync_.Thread)
	defer thread.ShutdownContext(t.Context())
	c := make(chan bool)
	const N = 10
	for range N {
		go run(t.Context(), thread, o, c)
	}
	for range N {
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
			r := recover()
			if r == nil {
				t.Fatalf("Thread.Do did not panic")
			}
			if r != "failed" {
				t.Fatalf("Thread.Do panic with %v, want %v", r, "failed")
			}
		}()
		_ = thread.Do(t.Context(), func() {
			panic("failed")
		})
	}()

	{
		var do bool
		_ = thread.Do(t.Context(), func() {
			do = true
		})
		if !do {
			t.Fatalf("Thread.Do did not called")
		}
	}
	thread.Shutdown()

	{
		var do bool
		_ = thread.Do(t.Context(), func() {
			do = true
		})
		if do {
			t.Fatalf("Thread.Do called after Thread.Shutdown")
		}
	}
}

func TestThreadCancel(t *testing.T) {
	var thread sync_.Thread
	var errThreadClosedByUser = errors.New("sync: Thread closed by user")
	{
		var wg sync.WaitGroup
		wg.Add(1)
		ctx, cancel := context.WithCancelCause(t.Context())

		go func() {
			defer wg.Done()
			// Block until canceled
			err := thread.Do(ctx, func() {
				select {}
			})
			if err == nil {
				t.Errorf("Thread.Do did not return error")
				return
			}

			if !errors.Is(err, errThreadClosedByUser) {
				t.Errorf("Thread.Do did not return errThreadClosedByUser")
				return
			}
		}()
		cancel(errThreadClosedByUser)
		wg.Wait()
	}

	var wg sync.WaitGroup
	release := make(chan struct{})
	wg.Go(func() {
		// Either rejected by Shutdown, or run on the Thread until released.
		err := thread.Do(t.Context(), func() {
			<-release
		})
		if err != nil && !errors.Is(err, sync_.ErrThreadClosed) {
			t.Errorf("Thread.Do = %v, want nil or %v", err, sync_.ErrThreadClosed)
		}
	})
	defer func() {
		close(release)
		wg.Wait()
		if err := thread.ShutdownContext(t.Context()); err != nil {
			t.Errorf("ShutdownContext after f released = %v, want nil", err)
		}
	}()

	thread.Shutdown()

	{
		var do bool
		err := thread.Do(t.Context(), func() {
			do = true
		})
		if do {
			t.Fatalf("Thread.Do called after Thread.Shutdown")
		}
		if err == nil {
			t.Errorf("Thread.Do did not return error")
			return
		}

		if !errors.Is(err, sync_.ErrThreadClosed) {
			t.Errorf("Thread.Do did not return ErrThreadClosed")
			return
		}
	}
}

func BenchmarkThread(b *testing.B) {
	var thread sync_.Thread
	defer thread.ShutdownContext(b.Context())
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = thread.Do(b.Context(), f)
		}
	})
}

var osThreadLeak = expvar_.NewLeak("os_thread_leak")
var goroutineLeak = expvar_.NewLeak("goroutine_leak")
var handlerLeak = expvar_.NewLeak("handler_leak")

func BenchmarkThreadWithLeak(b *testing.B) {
	var thread sync_.Thread
	thread.OSThreadLeak = osThreadLeak
	thread.GoroutineLeak = goroutineLeak
	thread.HandlerLeak = handlerLeak
	defer b.Logf("%s %s %s ", osThreadLeak, goroutineLeak, handlerLeak)
	defer thread.ShutdownContext(b.Context())
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = thread.Do(b.Context(), f)
		}
	})
}

func TestThreadWithLeak(t *testing.T) {
	o := new(one)
	thread := new(sync_.Thread)
	thread.OSThreadLeak = osThreadLeak
	thread.GoroutineLeak = goroutineLeak
	thread.HandlerLeak = handlerLeak

	defer thread.ShutdownContext(t.Context())
	c := make(chan bool)
	const N = 10
	for range N {
		go run(t.Context(), thread, o, c)
	}
	for range N {
		<-c
	}
	if *o != N {
		t.Errorf("once failed outside run: %d is not %d", *o, N)
	}
	if got := thread.OSThreadLeak.String(); got != "[1 1]" {
		t.Errorf("OSThreadLeak.String() = %q, want %q", got, "[1 1]")
	}
	if got := thread.GoroutineLeak.String(); got != "[1 1]" {
		t.Errorf("GoroutineLeak.String() = %q, want %q", got, "[1 1]")
	}

	want := fmt.Sprintf("[0 %d]", N)
	if got := thread.HandlerLeak.String(); got != want {
		t.Errorf("HandlerLeak.String() = %q, want %q", got, want)
	}
}

func TestThread_WatchStats(t *testing.T) {
	o := new(one)
	thread := new(sync_.Thread)
	thread.WatchStats()

	defer thread.ShutdownContext(t.Context())
	c := make(chan bool)
	const N = 10
	for range N {
		go run(t.Context(), thread, o, c)
	}
	for range N {
		<-c
	}
	if *o != N {
		t.Errorf("once failed outside run: %d is not %d", *o, N)
	}
	if got := thread.OSThreadLeak.String(); got != "[1 1]" {
		t.Errorf("OSThreadLeak.String() = %q, want %q", got, "[1 1]")
	}
	if got := thread.GoroutineLeak.String(); got != "[1 1]" {
		t.Errorf("GoroutineLeak.String() = %q, want %q", got, "[1 1]")
	}

	want := fmt.Sprintf("[0 %d]", N)
	if got := thread.HandlerLeak.String(); got != want {
		t.Errorf("HandlerLeak.String() = %q, want %q", got, want)
	}
	t.Logf("sync.Thread: %s", expvar.Get("sync.Thread"))
}

func TestThread_Stats(t *testing.T) {
	var thread sync_.Thread
	if got, want := thread.Stats(), (sync_.ThreadStats{}); got != want {
		t.Fatalf("Stats of unused Thread = %+v, want %+v", got, want)
	}

	started := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		_ = thread.Do(t.Context(), func() {
			close(started)
			<-release
		})
	})
	<-started
	if s := thread.Stats(); !s.Running || s.Active != 1 || s.Calls != 1 {
		t.Errorf("Stats with f running = %+v, want Running, Active 1, Calls 1", s)
	}
	close(release)
	wg.Wait()
	if s := thread.Stats(); !s.Running || s.Active != 0 || s.Calls != 1 || s.WaitDuration <= 0 {
		t.Errorf("Stats after f returned = %+v, want Running, Active 0, Calls 1, WaitDuration > 0", s)
	}

	if err := thread.ShutdownContext(t.Context()); err != nil {
		t.Fatalf("ShutdownContext = %v, want nil", err)
	}
	if s := thread.Stats(); s.Running {
		t.Errorf("Stats after Shutdown = %+v, want not Running", s)
	}
}

func TestThread_PprofLabels(t *testing.T) {
	var thread sync_.Thread
	defer thread.ShutdownContext(t.Context())
	// Start the worker from a labeled goroutine, whose labels must not be inherited by the worker.
	pprof.Do(t.Context(), pprof.Labels("caller", "test"), func(ctx context.Context) {
		_ = thread.Do(ctx, func() {})
	})

	var buf bytes.Buffer
	if err := pprof.Lookup("goroutine").WriteTo(&buf, 1); err != nil {
		t.Fatalf("write goroutine profile: %v", err)
	}
	var found bool
	for record := range strings.SplitSeq(buf.String(), "\n\n") {
		if !strings.Contains(record, "(*Thread).loop") {
			continue
		}
		if !strings.Contains(record, `# labels: {"sync.Thread":"worker"}`) {
			t.Errorf("worker goroutine labels, want only \"sync.Thread\":\"worker\":\n%s", record)
		}
		found = true
	}
	if !found {
		t.Fatalf("worker goroutine not found in goroutine profile:\n%s", buf.String())
	}
}

func TestThread_Goexit(t *testing.T) {
	var thread sync_.Thread
	defer thread.ShutdownContext(t.Context())
	// Escape first, as runtime.Goexit on the Thread closes it.
	for _, escape := range []bool{true, false} {
		var returned bool
		var wg sync.WaitGroup
		wg.Go(func() {
			_ = thread.Do(t.Context(), runtime.Goexit, sync_.WithThreadDoEscapeThread(escape))
			returned = true
		})
		wg.Wait()
		if returned {
			t.Fatalf("Do of f calling runtime.Goexit returned in caller's goroutine, escape=%v", escape)
		}
	}
	if err := thread.Do(t.Context(), func() {}); !errors.Is(err, sync_.ErrThreadClosed) {
		t.Fatalf("Do after runtime.Goexit = %v, want %v", err, sync_.ErrThreadClosed)
	}
}

func TestThread_ShutdownBeforeUse(t *testing.T) {
	var thread sync_.Thread
	thread.GoroutineLeak = new(expvar_.Leak)
	thread.OSThreadLeak = new(expvar_.Leak)
	if err := thread.ShutdownContext(t.Context()); err != nil {
		t.Fatalf("ShutdownContext of unused Thread = %v, want nil", err)
	}
	var do bool
	if err := thread.Do(t.Context(), func() { do = true }); !errors.Is(err, sync_.ErrThreadClosed) {
		t.Fatalf("Do after Shutdown = %v, want %v", err, sync_.ErrThreadClosed)
	}
	if do {
		t.Fatalf("Thread.Do called after Thread.Shutdown")
	}
	if n := thread.GoroutineLeak.New.Value(); n != 0 {
		t.Errorf("worker goroutine started %d times, want 0", n)
	}
	if n := thread.OSThreadLeak.New.Value(); n != 0 {
		t.Errorf("OS thread locked %d times, want 0", n)
	}
}

func TestThread_ShutdownNoWait(t *testing.T) {
	var thread sync_.Thread
	started := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		_ = thread.Do(t.Context(), func() {
			close(started)
			<-release
		})
	})
	<-started

	thread.Shutdown() // must return while f is still running
	if err := thread.Do(t.Context(), func() {}); !errors.Is(err, sync_.ErrThreadClosed) {
		t.Errorf("Do after Shutdown = %v, want %v", err, sync_.ErrThreadClosed)
	}
	if !thread.Stats().Running {
		t.Errorf("worker goroutine exited before running f returned")
	}

	close(release)
	wg.Wait()
	if err := thread.ShutdownContext(t.Context()); err != nil {
		t.Errorf("ShutdownContext after Shutdown = %v, want nil", err)
	}
}

func TestThread_ShutdownGraceful(t *testing.T) {
	var thread sync_.Thread
	started := make(chan struct{})
	release := make(chan struct{})
	var finished atomic.Bool
	doErr := make(chan error, 1)
	go func() {
		doErr <- thread.Do(t.Context(), func() {
			close(started)
			<-release
			finished.Store(true)
		})
	}()
	<-started

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	if err := thread.ShutdownContext(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("ShutdownContext with running f = %v, want %v", err, context.DeadlineExceeded)
	}
	if err := thread.Do(t.Context(), func() {}); !errors.Is(err, sync_.ErrThreadClosed) {
		t.Fatalf("Do after Shutdown = %v, want %v", err, sync_.ErrThreadClosed)
	}

	close(release)
	if err := thread.ShutdownContext(t.Context()); err != nil {
		t.Fatalf("ShutdownContext after f returned = %v, want nil", err)
	}
	if !finished.Load() {
		t.Fatalf("ShutdownContext returned before running f finished")
	}
	if err := <-doErr; err != nil {
		t.Fatalf("Do of running f = %v, want nil", err)
	}
}
