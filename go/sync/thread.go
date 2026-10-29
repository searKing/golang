// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"errors"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"
)

// pprof label of the Thread's worker goroutine.
const (
	threadLabelKey    = "sync.Thread"
	threadLabelWorker = "worker"
)

// ErrThreadClosed is returned by the Thread's Do and DoContext methods
// after a call to Shutdown or ShutdownContext,
// or after f called runtime.Goexit on the Thread.
var ErrThreadClosed = errors.New("sync: Thread closed")

//go:generate go-option -type=threadDo
type threadDo struct {
	// EscapeThread runs f directly in the caller's goroutine instead of on the Thread,
	// e.g. when calling Do from within f. f is then neither run on the Thread's OS thread,
	// nor serialized with other f; so outside of f, use it only for functions safe
	// to call concurrently from any thread.
	EscapeThread bool
}

// Thread runs functions on a single dedicated goroutine locked to an OS thread
// (see runtime.LockOSThread), for calling OS services or non-Go libraries
// that depend on per-thread state.
//
// Unless GoRoutine is true, the OS thread is never used by other goroutines: when the Thread exits,
// the thread is terminated, so f may change its state permanently, e.g. by setns(2) on Linux.
//
// A Thread cannot be relied on to run f on the main thread of the process, which some APIs require,
// e.g. Cocoa on macOS; to run code there, call runtime.LockOSThread in an init function
// of package main, and run that code from function main.
//
// The zero value is ready to use. GoRoutine must be set before the first call
// to Do, DoContext, Shutdown or ShutdownContext, and must not be changed afterwards.
//
// The worker goroutine carries the pprof label "sync.Thread"="worker",
// so live Threads, e.g. ones never shut down, can be found in goroutine profiles.
type Thread struct {
	// GoRoutine, if true, runs f on a dedicated goroutine
	// without locking it to an OS thread.
	GoRoutine bool

	once       sync.Once
	inShutdown atomic.Bool   // true when Thread is in shutdown
	osThread   atomic.Uint64 // ID of the OS thread locked by the worker goroutine, 0 if unknown

	// stats
	running      atomic.Bool
	active       atomic.Int64
	calls        atomic.Int64
	waitDuration atomic.Int64 // in nanoseconds

	fCh  chan func()   // delivers f to the worker goroutine
	done chan struct{} // closed when the worker goroutine has exited

	ctx    context.Context // canceled with ErrThreadClosed when Thread is closed
	cancel context.CancelCauseFunc
}

// ThreadStats contains Thread statistics.
type ThreadStats struct {
	// Thread Status
	Running bool // Whether the worker goroutine is running.
	Active  int  // The number of calls to Do that have not returned.

	// Counters
	Calls        int64         // The total number of calls to Do.
	WaitDuration time.Duration // The total time blocked waiting for the Thread to take f.
}

// Stats returns Thread statistics.
//
// To publish them by expvar, use expvar.Func:
//
//	expvar.Publish("thread", expvar.Func(func() any { return t.Stats() }))
func (t *Thread) Stats() ThreadStats {
	return ThreadStats{
		Running:      t.running.Load(),
		Active:       int(t.active.Load()),
		Calls:        t.calls.Load(),
		WaitDuration: time.Duration(t.waitDuration.Load()),
	}
}

// Do calls f on the Thread, or in the caller's goroutine if WithThreadDoEscapeThread(true) is given,
// and waits for f to return.
//
// Do returns ErrThreadClosed without calling f if the Thread is closed,
// or the cause of ctx if ctx is done before f starts running.
// Once f starts running, ctx is ignored and Do waits until f returns.
//
// If f panics, Do panics in the caller's goroutine with the same value.
// If f calls runtime.Goexit, Do calls runtime.Goexit in the caller's goroutine;
// if f was running on the Thread, the Thread is closed as well.
//
// Calling Do from within f calls the inner f directly on Linux and Windows unless GoRoutine is true,
// as the caller is already on the Thread. Elsewhere it blocks until ctx is done,
// as the Thread is busy running f; pass WithThreadDoEscapeThread(true) there.
func (t *Thread) Do(ctx context.Context, f func(), opts ...ThreadDoOption) error {
	var opt threadDo
	opt.ApplyOptions(opts...)

	t.initOnce()
	t.calls.Add(1)
	t.active.Add(1)
	defer t.active.Add(-1)
	// Check done first, since select chooses randomly among ready cases;
	// a closed Thread then fails fast, and never calls f with WithThreadDoEscapeThread(true).
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-t.ctx.Done():
		return context.Cause(t.ctx)
	default:
	}

	c := &threadCall{f: f, done: make(chan struct{})}
	if opt.EscapeThread || t.onThread() {
		c.run()
		c.result()
		return nil
	}
	if err := t.enqueue(ctx, c.run); err != nil {
		return err
	}
	<-c.done
	c.result()
	return nil
}

// DoContext is like Do, but passes f a context derived from ctx, which is also canceled
// with ErrThreadClosed when Shutdown or ShutdownContext is called, so a long-running f can return early.
// The context is canceled when DoContext returns.
func (t *Thread) DoContext(ctx context.Context, f func(ctx context.Context), opts ...ThreadDoOption) error {
	t.initOnce()
	fctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	// Do rejects f on a closed Thread, so skip AfterFunc, which would start a goroutine at once.
	if t.ctx.Err() == nil {
		stop := context.AfterFunc(t.ctx, func() { cancel(context.Cause(t.ctx)) })
		defer stop()
	}
	return t.Do(ctx, func() { f(fctx) }, opts...)
}

// onThread reports whether the caller is the worker goroutine, i.e. calling Do from within f.
// It is sound, as no other goroutine runs on the OS thread locked by the worker goroutine.
func (t *Thread) onThread() bool {
	worker := t.osThread.Load()
	if worker == 0 {
		return false
	}
	id, ok := osThreadID()
	return ok && id == worker
}

// enqueue blocks until the worker goroutine takes f, ctx is done or the Thread is closed.
func (t *Thread) enqueue(ctx context.Context, f func()) error {
	waitStart := time.Now()
	defer func() { t.waitDuration.Add(int64(time.Since(waitStart))) }()
	select {
	case t.fCh <- f:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-t.ctx.Done():
		return context.Cause(t.ctx)
	}
}

// Shutdown is like ShutdownContext, but returns without waiting for the running f to return
// and the worker goroutine to exit.
func (t *Thread) Shutdown() {
	t.close()
}

// ShutdownContext gracefully shuts down the Thread without interrupting a running f.
// ShutdownContext works by first stopping acceptance of new calls to Do,
// and canceling the contexts passed to f by DoContext, then waiting
// for the f running on the Thread to return and the worker goroutine to exit,
// which terminates the OS thread. f running in the caller's goroutine with
// WithThreadDoEscapeThread(true) is not waited for.
// If the provided context expires before the shutdown is complete,
// ShutdownContext returns the context's cause.
//
// Once Shutdown or ShutdownContext has been called on a Thread, it may not be reused;
// future calls to Do will return ErrThreadClosed.
// ShutdownContext may be called multiple times and concurrently.
// Calling ShutdownContext from within f blocks until ctx is done, as the Thread is busy running f.
func (t *Thread) ShutdownContext(ctx context.Context) error {
	t.close()
	select {
	case <-t.done:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}

// close stops the Thread from accepting new f, and signals the worker goroutine to exit.
func (t *Thread) close() {
	t.inShutdown.Store(true) // must precede initOnce, so a never used Thread does not start the worker
	t.initOnce()
	t.cancel(ErrThreadClosed)
}

func (t *Thread) initOnce() {
	t.once.Do(func() {
		t.ctx, t.cancel = context.WithCancelCause(context.Background())
		t.fCh = make(chan func())
		t.done = make(chan struct{})
		if t.inShutdown.Load() {
			t.cancel(ErrThreadClosed)
			close(t.done)
			return
		}
		go t.loop()
	})
}

// testHookThreadLoop, if set, is called by the worker goroutine of t when it starts.
// It is atomic, as worker goroutines of other Threads may still be starting when it is set.
var testHookThreadLoop atomic.Pointer[func(t *Thread)]

// loop runs f received from Do on the worker goroutine until the Thread is closed.
func (t *Thread) loop() {
	if hook := testHookThreadLoop.Load(); hook != nil {
		(*hook)(t)
	}
	// Replace the labels inherited from the goroutine calling Do first,
	// so the worker can be found by its own label in goroutine profiles.
	pprof.SetGoroutineLabels(pprof.WithLabels(context.Background(), pprof.Labels(threadLabelKey, threadLabelWorker)))
	defer close(t.done)
	defer t.cancel(ErrThreadClosed) // also closes the Thread if f calls runtime.Goexit
	t.running.Store(true)
	defer t.running.Store(false)
	if !t.GoRoutine {
		// Never unlock: f may change the state of the OS thread permanently,
		// so the thread is terminated with the worker goroutine instead of running other goroutines.
		runtime.LockOSThread()
		if id, ok := osThreadID(); ok {
			t.osThread.Store(id)
			// Clear before the thread is terminated, as its ID may be reused by a new thread.
			defer t.osThread.Store(0)
		}
	}
	for {
		// Check done first, since select chooses randomly among ready cases;
		// a closed Thread then exits instead of taking new f.
		select {
		case <-t.ctx.Done():
			return
		default:
		}
		select {
		case f := <-t.fCh:
			f()
		case <-t.ctx.Done():
			return
		}
	}
}

// threadCall is a single f submitted to Do.
type threadCall struct {
	f          func()
	done       chan struct{} // closed when f has returned, panicked or called runtime.Goexit
	returned   bool          // f returned normally
	panicValue any           // value f panicked with
}

func (c *threadCall) run() {
	defer close(c.done)
	defer func() {
		if !c.returned {
			c.panicValue = recover() // nil if f called runtime.Goexit
		}
	}()
	c.f()
	c.returned = true
}

// result replays in the caller's goroutine how f ended: panic, runtime.Goexit or normal return.
func (c *threadCall) result() {
	if c.panicValue != nil {
		panic(c.panicValue)
	}
	if !c.returned {
		runtime.Goexit()
	}
}
