// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"errors"
	"fmt"
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

// ErrThreadClosed is returned by the Thread's Do method, without calling f,
// after a call to Close or Shutdown, or after an earlier f called runtime.Goexit.
// In the latter case it is wrapped to report why, so use errors.Is(err, ErrThreadClosed).
var ErrThreadClosed = errors.New("sync: Thread closed")

// errThreadGoexit closes a Thread whose f called runtime.Goexit.
var errThreadGoexit = fmt.Errorf("%w: function passed to Do called runtime.Goexit", ErrThreadClosed)

// Thread runs functions on a single dedicated goroutine locked to an OS thread
// (see runtime.LockOSThread), for calling OS services or non-Go libraries
// that depend on per-thread state.
//
// The OS thread is never used by other goroutines: when the Thread exits,
// the thread is terminated, so f may change its state permanently, e.g. by setns(2) on Linux.
// The runtime never terminates the main thread of the process, so on Linux and Windows,
// a Thread never uses it. Elsewhere a Thread may happen to; when the Thread exits,
// the main thread is then blocked forever instead, still never used by other goroutines.
//
// Thus a Thread cannot be relied on to run f on the main thread, which some APIs require,
// e.g. Cocoa on macOS; to run code there, call runtime.LockOSThread in an init function
// of package main, and run that code from function main.
//
// The zero value is ready to use. A Thread must not be copied after first use.
//
// A Thread that has been used holds its goroutine and OS thread until it is closed.
// The garbage collector can close a Thread some time after it is no longer referenced,
// so calling Close is not required; but, as with an os.File, it is recommended
// to call Close or Shutdown once the Thread is no longer needed,
// to release the OS thread at once rather than whenever the garbage collector runs.
// The worker goroutine carries the pprof label "sync.Thread"="worker",
// so live Threads can be found in goroutine profiles.
type Thread struct {
	once   sync.Once
	closed atomic.Bool                 // true once Close has been called
	state  atomic.Pointer[threadState] // set by initOnce
}

// threadState is the state of a Thread shared with its worker goroutine.
// It does not refer to the Thread, so the worker goroutine does not keep
// an unreferenced Thread from being garbage collected, and thus closed.
type threadState struct {
	worker atomic.Uint64 // callerID of the worker goroutine, 0 if not running

	// stats
	running      atomic.Bool
	active       atomic.Int64
	calls        atomic.Int64
	waitDuration atomic.Int64 // in nanoseconds

	callCh chan *threadCall // delivers calls to the worker goroutine
	done   chan struct{}    // closed when the worker goroutine has exited

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
	s := t.state.Load()
	if s == nil {
		return ThreadStats{}
	}
	return ThreadStats{
		Running:      s.running.Load(),
		Active:       int(s.active.Load()),
		Calls:        s.calls.Load(),
		WaitDuration: time.Duration(s.waitDuration.Load()),
	}
}

// Do calls f on the Thread's goroutine and waits for f to return.
// Called from another goroutine, Do runs f as if f were started there by a go statement;
// called from within f, as described below, Do calls f as a plain function call.
//
// The context passed to f is derived from ctx, and is also canceled
// with ErrThreadClosed by Close or Shutdown, so a long-running f can return early.
// The context passed to f is canceled when Do returns.
//
// Do waits for the Thread to finish running any other f before running f.
// If the Thread is closed, or ctx is done, before f starts running,
// Do returns an error matching ErrThreadClosed, or the cause of ctx,
// and f is never called.
// Once f starts running, Do waits until f returns, even if ctx is done.
//
// As on any goroutine, a panic in f that is not recovered on the Thread's goroutine
// crashes the program. If f calls runtime.Goexit, which exits the Thread's goroutine,
// the Thread is closed before Do returns nil, as f has been called.
//
// Calling Do from within f with the context passed to f, or one derived from it,
// calls the inner f directly, as the caller is already on the Thread;
// so a panic or runtime.Goexit in the inner f unwinds the outer f, as with any function call.
// Calling it with any other context blocks until that context is done,
// as the Thread is busy running f.
//
// Do panics if f is nil.
func (t *Thread) Do(ctx context.Context, f func(ctx context.Context)) error {
	if f == nil {
		panic("sync: Thread.Do called with nil function")
	}
	s := t.initOnce()
	s.calls.Add(1)
	s.active.Add(1)
	defer s.active.Add(-1)
	// Check done first, since select chooses randomly among ready cases;
	// a closed Thread then fails fast.
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-s.ctx.Done():
		return context.Cause(s.ctx)
	default:
	}

	fctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	stop := context.AfterFunc(s.ctx, func() { cancel(context.Cause(s.ctx)) })
	defer stop()

	if s.onThread(ctx) {
		f(fctx)
		return nil
	}
	c := &threadCall{ctx: fctx, f: f, done: make(chan struct{})}
	if err := s.enqueue(ctx, c); err != nil {
		return err
	}
	<-c.done
	// Keep t reachable until f returns, so the garbage collector does not close it while f runs.
	runtime.KeepAlive(t)
	return nil
}

// threadCtxKey is the context key marking the contexts passed to f by the worker goroutine of a Thread.
type threadCtxKey struct{}

// onThread reports whether the caller is the worker goroutine, i.e. calling Do from within f with ctx.
func (s *threadState) onThread(ctx context.Context) bool {
	// ctx alone is not enough, as f may pass it to other goroutines.
	// Check it first, as callerID is slow on some platforms.
	if ctx.Value(threadCtxKey{}) != s {
		return false
	}
	worker := s.worker.Load()
	return worker != 0 && callerID() == worker
}

// enqueue blocks until the worker goroutine takes c, ctx is done or the Thread is closed.
func (s *threadState) enqueue(ctx context.Context, c *threadCall) error {
	waitStart := time.Now()
	defer func() { s.waitDuration.Add(int64(time.Since(waitStart))) }()
	select {
	case s.callCh <- c:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-s.ctx.Done():
		return context.Cause(s.ctx)
	}
}

// Close stops the Thread from accepting new f, and cancels the contexts passed to f
// with ErrThreadClosed. Close does not wait for the running f to return
// and the worker goroutine to exit; for that, use Shutdown.
//
// Once Close or Shutdown has been called on a Thread, it may not be reused;
// future calls to Do will return ErrThreadClosed.
// Close may be called multiple times and concurrently.
func (t *Thread) Close() {
	t.closed.Store(true) // must precede initOnce, so a never used Thread does not start the worker
	t.initOnce().cancel(ErrThreadClosed)
}

// Shutdown gracefully shuts down the Thread without interrupting a running f.
// Shutdown works by first calling Close, then waiting for the running f to return
// and the worker goroutine to exit, which terminates the OS thread.
// If the provided context expires before the shutdown is complete,
// Shutdown returns the context's cause.
//
// Shutdown may be called multiple times and concurrently.
// Calling Shutdown from within f blocks until ctx is done, as the Thread is busy running f.
func (t *Thread) Shutdown(ctx context.Context) error {
	t.Close()
	select {
	case <-t.state.Load().done:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}

func (t *Thread) initOnce() *threadState {
	t.once.Do(func() {
		s := &threadState{
			callCh: make(chan *threadCall),
			done:   make(chan struct{}),
		}
		s.ctx, s.cancel = context.WithCancelCause(context.Background())
		t.state.Store(s)
		if t.closed.Load() {
			s.cancel(ErrThreadClosed)
			close(s.done)
			return
		}
		runtime.AddCleanup(t, func(s *threadState) { s.cancel(ErrThreadClosed) }, s)
		go s.loop()
	})
	return t.state.Load()
}

// testHookThreadLoop, if set, is called by the worker goroutine of s when it starts.
// It is atomic, as worker goroutines of other Threads may still be starting when it is set.
var testHookThreadLoop atomic.Pointer[func(s *threadState)]

// testHookOnStartupThread, if set, replaces onStartupThread for the worker goroutine of s.
var testHookOnStartupThread atomic.Pointer[func(s *threadState) bool]

// loop locks the worker goroutine to an OS thread, and serves calls on it until the Thread is closed.
func (s *threadState) loop() {
	if hook := testHookThreadLoop.Load(); hook != nil {
		(*hook)(s)
	}
	runtime.LockOSThread()
	onStartup := onStartupThread()
	if hook := testHookOnStartupThread.Load(); hook != nil {
		onStartup = (*hook)(s)
	}
	if onStartup {
		// The runtime never terminates the startup thread, but blocks it forever
		// once its locked goroutine exits, so the state f leaves on it would stay,
		// e.g. visible process-wide through /proc/self on Linux.
		// Serve on a new goroutine instead, which cannot run on this thread while it is locked,
		// and then release this thread, which no f has run on.
		locked := make(chan struct{})
		go func() {
			runtime.LockOSThread()
			close(locked)
			s.serve()
		}()
		<-locked
		runtime.UnlockOSThread()
		return
	}
	s.serve()
}

// serve runs calls received from Do on the worker goroutine,
// which is locked to its OS thread, until the Thread is closed.
// It never unlocks the thread: f may change its state permanently,
// so the thread is terminated with the worker goroutine instead of running other goroutines.
func (s *threadState) serve() {
	// Replace the labels inherited from the goroutine calling Do first,
	// so the worker can be found by its own label in goroutine profiles.
	pprof.SetGoroutineLabels(pprof.WithLabels(context.Background(), pprof.Labels(threadLabelKey, threadLabelWorker)))
	defer close(s.done)
	s.running.Store(true)
	defer s.running.Store(false)
	s.worker.Store(callerID())
	// Clear before the thread is terminated, as its ID may be reused by a new thread.
	defer s.worker.Store(0)
	for {
		// Check done first, since select chooses randomly among ready cases;
		// a closed Thread then exits instead of taking new f.
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		select {
		case c := <-s.callCh:
			c.ctx = context.WithValue(c.ctx, threadCtxKey{}, s)
			s.run(c)
		case <-s.ctx.Done():
			return
		}
	}
}

// threadCall is a single f submitted to Do.
type threadCall struct {
	ctx  context.Context // passed to f
	f    func(ctx context.Context)
	done chan struct{} // closed when f has returned or called runtime.Goexit
}

// run calls c.f on the worker goroutine.
func (s *threadState) run(c *threadCall) {
	returned := false
	defer func() {
		if x := recover(); x != nil {
			// f panicked, which is fatal, as on any goroutine.
			// Leave done open, so Do does not return and race with the crash.
			panic(x)
		}
		if !returned {
			// f called runtime.Goexit, which exits the worker goroutine;
			// close the Thread before Do returns.
			s.cancel(errThreadGoexit)
		}
		close(c.done)
	}()
	c.f(c.ctx)
	returned = true
}
