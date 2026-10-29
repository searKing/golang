// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

// SetTestHookThreadLoop sets the hook called by the worker goroutine of a Thread when it starts,
// and returns a function restoring no hook.
func SetTestHookThreadLoop(hook func(t *Thread)) (restore func()) {
	testHookThreadLoop.Store(&hook)
	return func() { testHookThreadLoop.Store(nil) }
}
