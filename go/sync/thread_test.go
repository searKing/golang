// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	_ = thread.Do(ctx, func(context.Context) { o.Increment() })
	c <- true
}

func TestThread(t *testing.T) {
	o := new(one)
	thread := new(sync_.Thread)
	defer thread.Shutdown(t.Context())
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

func threadPanicker(context.Context) { panic("failed") }

func TestThreadPanic(t *testing.T) {
	if os.Getenv("SYNC_TEST_THREAD_PANIC") == "1" {
		var thread sync_.Thread
		_ = thread.Do(context.Background(), threadPanicker)
		fmt.Println("Do returned")
		return
	}
	switch runtime.GOOS {
	case "js", "wasip1", "ios", "android":
		t.Skipf("cannot exec subprocess on %s", runtime.GOOS)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestThreadPanic$")
	cmd.Env = append(os.Environ(), "SYNC_TEST_THREAD_PANIC=1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("panic in f did not crash the program:\n%s", out)
	}
	// The crash must report where f panicked, not where it was rethrown.
	for _, want := range []string{"panic: failed", "sync_test.threadPanicker("} {
		if !bytes.Contains(out, []byte(want)) {
			t.Errorf("crash output does not contain %q:\n%s", want, out)
		}
	}
	if bytes.Contains(out, []byte("Do returned")) {
		t.Errorf("Do returned after f panicked:\n%s", out)
	}
}

func TestThreadCancel(t *testing.T) {
	var thread sync_.Thread
	var errThreadClosedByUser = errors.New("sync: Thread closed by user")
	{
		// Keep the Thread busy, so f below waits to be taken until canceled.
		busy := make(chan struct{})
		release := make(chan struct{})
		var busyWG sync.WaitGroup
		busyWG.Go(func() {
			_ = thread.Do(t.Context(), func(context.Context) {
				close(busy)
				<-release
			})
		})
		<-busy

		var wg sync.WaitGroup
		wg.Add(1)
		ctx, cancel := context.WithCancelCause(t.Context())

		go func() {
			defer wg.Done()
			err := thread.Do(ctx, func(context.Context) {
				t.Errorf("Thread.Do called f with ctx canceled before f was taken")
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
		close(release)
		busyWG.Wait()
	}

	var wg sync.WaitGroup
	release := make(chan struct{})
	wg.Go(func() {
		// Either rejected by Close, or run on the Thread until released.
		err := thread.Do(t.Context(), func(context.Context) {
			<-release
		})
		if err != nil && !errors.Is(err, sync_.ErrThreadClosed) {
			t.Errorf("Thread.Do = %v, want nil or %v", err, sync_.ErrThreadClosed)
		}
	})
	defer func() {
		close(release)
		wg.Wait()
		if err := thread.Shutdown(t.Context()); err != nil {
			t.Errorf("Shutdown after f released = %v, want nil", err)
		}
	}()

	thread.Close()

	{
		var do bool
		err := thread.Do(t.Context(), func(context.Context) {
			do = true
		})
		if do {
			t.Fatalf("Thread.Do called after Thread.Close")
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
	defer thread.Shutdown(b.Context())
	f := func(context.Context) {}
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
		_ = thread.Do(t.Context(), func(context.Context) {
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

	if err := thread.Shutdown(t.Context()); err != nil {
		t.Fatalf("Shutdown = %v, want nil", err)
	}
	if s := thread.Stats(); s.Running {
		t.Errorf("Stats after Shutdown = %+v, want not Running", s)
	}
}

func TestThread_PprofLabels(t *testing.T) {
	var thread sync_.Thread
	defer thread.Shutdown(t.Context())
	// Start the worker from a labeled goroutine, whose labels must not be inherited by the worker.
	pprof.Do(t.Context(), pprof.Labels("caller", "test"), func(ctx context.Context) {
		_ = thread.Do(ctx, func(context.Context) {})
	})

	var buf bytes.Buffer
	if err := pprof.Lookup("goroutine").WriteTo(&buf, 1); err != nil {
		t.Fatalf("write goroutine profile: %v", err)
	}
	var found bool
	for record := range strings.SplitSeq(buf.String(), "\n\n") {
		if !strings.Contains(record, "(*threadState).loop") {
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

func TestThread_Context(t *testing.T) {
	t.Run("caller canceled", func(t *testing.T) {
		var thread sync_.Thread
		defer thread.Shutdown(t.Context())
		errCaller := errors.New("caller canceled")
		ctx, cancel := context.WithCancelCause(context.WithValue(t.Context(), ctxKey{}, "v"))
		var value any
		var cause error
		err := thread.Do(ctx, func(ctx context.Context) {
			value = ctx.Value(ctxKey{})
			cancel(errCaller)
			<-ctx.Done()
			cause = context.Cause(ctx)
		})
		if err != nil {
			t.Fatalf("Do = %v, want nil", err)
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
			err := thread.Do(t.Context(), func(ctx context.Context) {
				close(started)
				<-ctx.Done()
				cause = context.Cause(ctx)
			})
			if err != nil {
				t.Errorf("Do = %v, want nil", err)
			}
		})
		<-started
		if err := thread.Shutdown(t.Context()); err != nil {
			t.Fatalf("Shutdown = %v, want nil", err)
		}
		wg.Wait()
		if !errors.Is(cause, sync_.ErrThreadClosed) {
			t.Errorf("cause of ctx passed to f = %v, want %v", cause, sync_.ErrThreadClosed)
		}
	})

	t.Run("close from within f", func(t *testing.T) {
		var thread sync_.Thread
		var cause error
		err := thread.Do(t.Context(), func(ctx context.Context) {
			thread.Close()
			<-ctx.Done()
			cause = context.Cause(ctx)
		})
		if err != nil {
			t.Fatalf("Do = %v, want nil", err)
		}
		if !errors.Is(cause, sync_.ErrThreadClosed) {
			t.Errorf("cause of ctx passed to f = %v, want %v", cause, sync_.ErrThreadClosed)
		}
		if err := thread.Shutdown(t.Context()); err != nil {
			t.Errorf("Shutdown = %v, want nil", err)
		}
	})

	t.Run("closed", func(t *testing.T) {
		var thread sync_.Thread
		thread.Close()
		var do bool
		err := thread.Do(t.Context(), func(context.Context) { do = true })
		if !errors.Is(err, sync_.ErrThreadClosed) {
			t.Errorf("Do on closed Thread = %v, want %v", err, sync_.ErrThreadClosed)
		}
		if do {
			t.Errorf("f called on closed Thread")
		}
		if s := thread.Stats(); s.Calls != 1 {
			t.Errorf("Stats.Calls = %d, want 1", s.Calls)
		}
	})
}

func TestThread_NestedDo(t *testing.T) {
	t.Run("inline", func(t *testing.T) {
		var thread sync_.Thread
		defer thread.Shutdown(t.Context())
		var outer, inner uint64
		err := thread.Do(t.Context(), func(ctx context.Context) {
			outer = sync_.CallerID()
			if err := thread.Do(ctx, func(context.Context) { inner = sync_.CallerID() }); err != nil {
				t.Errorf("nested Do = %v, want nil", err)
			}
		})
		if err != nil {
			t.Fatalf("Do = %v, want nil", err)
		}
		if inner == 0 || inner != outer {
			t.Errorf("nested f ran on %d, want %d of the Thread", inner, outer)
		}
	})

	t.Run("unrelated ctx", func(t *testing.T) {
		var thread sync_.Thread
		defer thread.Shutdown(t.Context())
		err := thread.Do(t.Context(), func(context.Context) {
			ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
			defer cancel()
			err := thread.Do(ctx, func(context.Context) {})
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("nested Do with unrelated ctx = %v, want %v", err, context.DeadlineExceeded)
			}
		})
		if err != nil {
			t.Fatalf("Do = %v, want nil", err)
		}
	})

	t.Run("panic", func(t *testing.T) {
		var thread sync_.Thread
		defer thread.Shutdown(t.Context())
		var r any
		err := thread.Do(t.Context(), func(ctx context.Context) {
			defer func() { r = recover() }()
			_ = thread.Do(ctx, threadPanicker)
		})
		if err != nil {
			t.Fatalf("Do = %v, want nil", err)
		}
		if r != "failed" {
			t.Errorf("outer f recovered %v from nested Do, want %v", r, "failed")
		}
	})

	t.Run("goroutine started by f", func(t *testing.T) {
		var thread sync_.Thread
		defer thread.Shutdown(t.Context())
		done := make(chan error, 1)
		err := thread.Do(t.Context(), func(ctx context.Context) {
			// Not on the Thread, so its f must wait for this f to return instead of running inline,
			// even with the ctx passed to f.
			ctx = context.WithoutCancel(ctx)
			go func() { done <- thread.Do(ctx, func(context.Context) {}) }()
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
}

func TestThread_Goexit(t *testing.T) {
	for _, nested := range []bool{false, true} {
		var thread sync_.Thread
		defer thread.Shutdown(t.Context())
		f := func(context.Context) { runtime.Goexit() }
		if nested {
			f = func(ctx context.Context) { _ = thread.Do(ctx, func(context.Context) { runtime.Goexit() }) }
		}
		// f has been called, so Do reports no error, as ErrThreadClosed means f was not called.
		if err := thread.Do(t.Context(), f); err != nil {
			t.Fatalf("Do of f calling runtime.Goexit = %v, want nil, nested=%v", err, nested)
		}
		if !sync_.ThreadClosed(&thread) {
			t.Fatalf("Thread not closed when Do of f calling runtime.Goexit returned, nested=%v", nested)
		}
		err := thread.Do(t.Context(), func(context.Context) {})
		if !errors.Is(err, sync_.ErrThreadClosed) {
			t.Fatalf("Do after runtime.Goexit = %v, want %v, nested=%v", err, sync_.ErrThreadClosed, nested)
		}
		if got, want := err.Error(), "sync: Thread closed: function passed to Do called runtime.Goexit"; got != want {
			t.Errorf("Do after runtime.Goexit = %q, want %q, nested=%v", got, want, nested)
		}
	}
}

func TestThread_ShutdownBeforeUse(t *testing.T) {
	// Count per Thread, as worker goroutines of Threads in other tests may still be starting.
	var mu sync.Mutex
	starts := make(map[any]int)
	defer sync_.SetTestHookThreadLoop(func(state any) {
		mu.Lock()
		defer mu.Unlock()
		starts[state]++
	})()
	started := func(thread *sync_.Thread) int {
		mu.Lock()
		defer mu.Unlock()
		return starts[sync_.ThreadState(thread)]
	}

	tests := []struct {
		name     string
		shutdown func(thread *sync_.Thread) error
	}{
		{"Close", func(thread *sync_.Thread) error { thread.Close(); return nil }},
		{"Shutdown", func(thread *sync_.Thread) error { return thread.Shutdown(t.Context()) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var thread sync_.Thread
			if err := tt.shutdown(&thread); err != nil {
				t.Fatalf("%s of unused Thread = %v, want nil", tt.name, err)
			}
			var do bool
			if err := thread.Do(t.Context(), func(context.Context) { do = true }); !errors.Is(err, sync_.ErrThreadClosed) {
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
	_ = thread.Do(t.Context(), func(context.Context) {})
	_ = thread.Do(t.Context(), func(context.Context) {})
	if err := thread.Shutdown(t.Context()); err != nil {
		t.Fatalf("Shutdown = %v, want nil", err)
	}
	if n := started(&thread); n != 1 {
		t.Errorf("worker goroutine of used Thread started %d times, want 1", n)
	}
}

func TestThread_Close(t *testing.T) {
	var thread sync_.Thread
	started := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		_ = thread.Do(t.Context(), func(context.Context) {
			close(started)
			<-release
		})
	})
	<-started

	thread.Close() // must return while f is still running
	if err := thread.Do(t.Context(), func(context.Context) {}); !errors.Is(err, sync_.ErrThreadClosed) {
		t.Errorf("Do after Close = %v, want %v", err, sync_.ErrThreadClosed)
	}
	if !thread.Stats().Running {
		t.Errorf("worker goroutine exited before running f returned")
	}

	close(release)
	wg.Wait()
	if err := thread.Shutdown(t.Context()); err != nil {
		t.Errorf("Shutdown after Close = %v, want nil", err)
	}
}

func TestThread_Shutdown(t *testing.T) {
	var thread sync_.Thread
	started := make(chan struct{})
	release := make(chan struct{})
	var finished atomic.Bool
	doErr := make(chan error, 1)
	go func() {
		doErr <- thread.Do(t.Context(), func(context.Context) {
			close(started)
			<-release
			finished.Store(true)
		})
	}()
	<-started

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	if err := thread.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown with running f = %v, want %v", err, context.DeadlineExceeded)
	}
	if err := thread.Do(t.Context(), func(context.Context) {}); !errors.Is(err, sync_.ErrThreadClosed) {
		t.Fatalf("Do after Shutdown = %v, want %v", err, sync_.ErrThreadClosed)
	}

	close(release)
	if err := thread.Shutdown(t.Context()); err != nil {
		t.Fatalf("Shutdown after f returned = %v, want nil", err)
	}
	if !finished.Load() {
		t.Fatalf("Shutdown returned before running f finished")
	}
	if err := <-doErr; err != nil {
		t.Fatalf("Do of running f = %v, want nil", err)
	}
}

func TestThread_NilFunc(t *testing.T) {
	var thread sync_.Thread
	defer thread.Shutdown(t.Context())
	defer func() {
		if r, want := recover(), "sync: Thread.Do called with nil function"; r != want {
			t.Errorf("Do with nil f panicked in the caller's goroutine with %v, want %q", r, want)
		}
		if s := thread.Stats(); s.Calls != 0 {
			t.Errorf("Stats.Calls after Do with nil f = %d, want 0", s.Calls)
		}
	}()
	_ = thread.Do(t.Context(), nil)
}

func TestThread_Unreferenced(t *testing.T) {
	done := func() <-chan struct{} {
		thread := new(sync_.Thread)
		if err := thread.Do(t.Context(), func(context.Context) {}); err != nil {
			t.Fatalf("Do = %v, want nil", err)
		}
		return sync_.ThreadDone(thread)
	}()

	for deadline := time.Now().Add(10 * time.Second); ; {
		runtime.GC()
		select {
		case <-done:
			return
		case <-time.After(time.Millisecond):
		}
		if time.Now().After(deadline) {
			t.Fatal("worker goroutine of an unreferenced Thread did not exit after garbage collection")
		}
	}
}

func TestThread_ReferencedOnlyByDo(t *testing.T) {
	// The Thread is referenced only by the running Do, which must keep it from being closed.
	err := new(sync_.Thread).Do(t.Context(), func(ctx context.Context) {
		for range 20 {
			runtime.GC()
			select {
			case <-ctx.Done():
				t.Errorf("ctx passed to f canceled while f runs: %v", context.Cause(ctx))
				return
			case <-time.After(time.Millisecond):
			}
		}
	})
	if err != nil {
		t.Fatalf("Do = %v, want nil", err)
	}
}
