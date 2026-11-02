// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

// SetTestHookThreadLoop sets the hook called by the worker goroutine of a Thread when it starts,
// with the ThreadState of the Thread, and returns a function restoring no hook.
func SetTestHookThreadLoop(hook func(state any)) (restore func()) {
	f := func(s *threadState) { hook(s) }
	testHookThreadLoop.Store(&f)
	return func() { testHookThreadLoop.Store(nil) }
}

// SetTestHookOnStartupThread sets the hook replacing OnStartupThread for the worker goroutine of a Thread,
// with the ThreadState of the Thread, and returns a function restoring no hook.
func SetTestHookOnStartupThread(hook func(state any) bool) (restore func()) {
	f := func(s *threadState) bool { return hook(s) }
	testHookOnStartupThread.Store(&f)
	return func() { testHookOnStartupThread.Store(nil) }
}

var (
	CallerID        = callerID
	OnStartupThread = onStartupThread
)

// ThreadState returns the state of t shared with its worker goroutine, or nil if t is unused.
func ThreadState(t *Thread) any {
	if s := t.state.Load(); s != nil {
		return s
	}
	return nil
}

// ThreadClosed reports whether t is closed, without waiting.
func ThreadClosed(t *Thread) bool {
	s := t.state.Load()
	return s != nil && s.ctx.Err() != nil
}

// ThreadDone returns a channel closed when the worker goroutine of t has exited.
// It does not refer to t, so t can still be garbage collected.
func ThreadDone(t *Thread) <-chan struct{} { return t.state.Load().done }
