// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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

type ctxKey struct{}

func TestThread_DoContext(t *testing.T) {
	t.Run("caller canceled", func(t *testing.T) {
		var thread sync_.Thread
		defer thread.ShutdownContext(t.Context())
		errCaller := errors.New("caller canceled")
		ctx, cancel := context.WithCancelCause(context.WithValue(t.Context(), ctxKey{}, "v"))
		var value any
		var cause error
		err := thread.DoContext(ctx, func(ctx context.Context) {
			value = ctx.Value(ctxKey{})
			cancel(errCaller)
			<-ctx.Done()
			cause = context.Cause(ctx)
		})
		if err != nil {
			t.Fatalf("DoContext = %v, want nil", err)
		}
		if value != "v" {
			t.Errorf("value of ctx passed to f = %v, want %v", value, "v")
		}
		if !errors.Is(cause, errCaller) {
			t.Errorf("cause of ctx passed to f = %v, want %v", cause, errCaller)
		}
	})

	t.Run("shutdown", func(t *testing.T) {
		var thread sync_.Thread
		started := make(chan struct{})
		var cause error
		var wg sync.WaitGroup
		wg.Go(func() {
			err := thread.DoContext(t.Context(), func(ctx context.Context) {
				close(started)
				<-ctx.Done()
				cause = context.Cause(ctx)
			})
			if err != nil {
				t.Errorf("DoContext = %v, want nil", err)
			}
		})
		<-started
		if err := thread.ShutdownContext(t.Context()); err != nil {
			t.Fatalf("ShutdownContext = %v, want nil", err)
		}
		wg.Wait()
		if !errors.Is(cause, sync_.ErrThreadClosed) {
			t.Errorf("cause of ctx passed to f = %v, want %v", cause, sync_.ErrThreadClosed)
		}
	})

	t.Run("closed", func(t *testing.T) {
		var thread sync_.Thread
		thread.Shutdown()
		var do bool
		err := thread.DoContext(t.Context(), func(ctx context.Context) { do = true })
		if !errors.Is(err, sync_.ErrThreadClosed) {
			t.Errorf("DoContext on closed Thread = %v, want %v", err, sync_.ErrThreadClosed)
		}
		if do {
			t.Errorf("f called on closed Thread")
		}
		if s := thread.Stats(); s.Calls != 1 {
			t.Errorf("Stats.Calls = %d, want 1", s.Calls)
		}
	})

	t.Run("escape thread", func(t *testing.T) {
		var thread sync_.Thread
		var ok bool
		err := thread.DoContext(t.Context(), func(ctx context.Context) {
			// Shutdown from within f is safe with escape, as f does not run on the Thread.
			_ = thread.ShutdownContext(ctx)
			<-ctx.Done()
			ok = errors.Is(context.Cause(ctx), sync_.ErrThreadClosed)
		}, sync_.WithThreadDoEscapeThread(true))
		if err != nil {
			t.Fatalf("DoContext = %v, want nil", err)
		}
		if !ok {
			t.Errorf("ctx passed to f is not canceled with %v by Shutdown", sync_.ErrThreadClosed)
		}
	})
}

func TestThread_NestedDo(t *testing.T) {
	skipIfNoOSThreadID := func(t *testing.T) {
		if _, ok := sync_.OSThreadID(); !ok {
			t.Skipf("OS thread ID is not supported on %s", runtime.GOOS)
		}
	}

	t.Run("inline", func(t *testing.T) {
		skipIfNoOSThreadID(t)
		var thread sync_.Thread
		defer thread.ShutdownContext(t.Context())
		var outer, inner uint64
		err := thread.Do(t.Context(), func() {
			outer, _ = sync_.OSThreadID()
			if err := thread.Do(t.Context(), func() { inner, _ = sync_.OSThreadID() }); err != nil {
				t.Errorf("nested Do = %v, want nil", err)
			}
		})
		if err != nil {
			t.Fatalf("Do = %v, want nil", err)
		}
		if inner == 0 || inner != outer {
			t.Errorf("nested f ran on OS thread %d, want %d of the Thread", inner, outer)
		}
	})

	t.Run("goroutine started by f", func(t *testing.T) {
		skipIfNoOSThreadID(t)
		var thread sync_.Thread
		defer thread.ShutdownContext(t.Context())
		done := make(chan error, 1)
		err := thread.Do(t.Context(), func() {
			// Not on the Thread, so its f must wait for this f to return instead of running inline.
			go func() { done <- thread.Do(t.Context(), func() {}) }()
			select {
			case <-done:
				t.Errorf("Do from a goroutine started by f returned while f is running")
			case <-time.After(20 * time.Millisecond):
			}
		})
		if err != nil {
			t.Fatalf("Do = %v, want nil", err)
		}
		if err := <-done; err != nil {
			t.Errorf("Do from a goroutine started by f = %v, want nil", err)
		}
	})

	t.Run("GoRoutine", func(t *testing.T) {
		thread := sync_.Thread{GoRoutine: true}
		defer thread.ShutdownContext(t.Context())
		err := thread.Do(t.Context(), func() {
			// Without a locked OS thread, the caller cannot be told apart, so nested Do still blocks.
			ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
			defer cancel()
			if err := thread.Do(ctx, func() {}); !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("nested Do with GoRoutine = %v, want %v", err, context.DeadlineExceeded)
			}
		})
		if err != nil {
			t.Fatalf("Do = %v, want nil", err)
		}
	})
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
	// Count per Thread, as worker goroutines of Threads in other tests may still be starting.
	var mu sync.Mutex
	starts := make(map[*sync_.Thread]int)
	defer sync_.SetTestHookThreadLoop(func(thread *sync_.Thread) {
		mu.Lock()
		defer mu.Unlock()
		starts[thread]++
	})()
	started := func(thread *sync_.Thread) int {
		mu.Lock()
		defer mu.Unlock()
		return starts[thread]
	}

	tests := []struct {
		name     string
		shutdown func(thread *sync_.Thread) error
	}{
		{"Shutdown", func(thread *sync_.Thread) error { thread.Shutdown(); return nil }},
		{"ShutdownContext", func(thread *sync_.Thread) error { return thread.ShutdownContext(t.Context()) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var thread sync_.Thread
			if err := tt.shutdown(&thread); err != nil {
				t.Fatalf("%s of unused Thread = %v, want nil", tt.name, err)
			}
			var do bool
			if err := thread.Do(t.Context(), func() { do = true }); !errors.Is(err, sync_.ErrThreadClosed) {
				t.Fatalf("Do after %s = %v, want %v", tt.name, err, sync_.ErrThreadClosed)
			}
			if do {
				t.Fatalf("Thread.Do called after Thread.%s", tt.name)
			}
			if n := started(&thread); n != 0 {
				t.Errorf("worker goroutine started %d times, want 0", n)
			}
		})
	}

	// The hook itself works: a used Thread starts the worker goroutine once.
	var thread sync_.Thread
	_ = thread.Do(t.Context(), func() {})
	_ = thread.Do(t.Context(), func() {})
	if err := thread.ShutdownContext(t.Context()); err != nil {
		t.Fatalf("ShutdownContext = %v, want nil", err)
	}
	if n := started(&thread); n != 1 {
		t.Errorf("worker goroutine of used Thread started %d times, want 1", n)
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
