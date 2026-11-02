// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || windows

package sync

// startupThreadID is the callerID of the startup thread, which runs all package initialization
// (see runtime.LockOSThread), and which the runtime never terminates.
var startupThreadID = callerID()

// onStartupThread reports whether the calling goroutine, locked to its OS thread, is on the startup thread.
func onStartupThread() bool { return callerID() == startupThreadID }
