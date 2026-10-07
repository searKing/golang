// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"context"
	"io"
	"os/exec"
	"syscall"
	"time"
)

// waitDelay is how long to wait for the process to exit after SIGTERM is sent, before it is killed.
const waitDelay = time.Second

// command returns the Cmd struct to execute the named program with the given arguments.
// The process is sent SIGTERM if ctx is done before the command completes on its own,
// and is killed if it does not exit within waitDelay.
// Signals other than Kill are not supported on Windows, so the process is killed after waitDelay there.
func command(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Cancel = func() error { return cmd.Process.Signal(syscall.SIGTERM) }
	cmd.WaitDelay = waitDelay
	return cmd
}

// commandContext starts the named program, feeds its stdout to handle, and waits for it to exit.
// ctx.Err() is returned if ctx is done before the command completes on its own.
func commandContext(ctx context.Context, handle func(reader io.Reader), name string, arg ...string) error {
	cmd := command(ctx, name, arg...)
	// Use a pipe other than cmd.StdoutPipe, so that Wait closes stdout after WaitDelay,
	// even if a subprocess holds stdout open.
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	if err := cmd.Start(); err != nil {
		return err
	}

	handled := make(chan struct{})
	go func() {
		defer close(handled)
		if handle != nil {
			handle(pr)
		}
		// drain stdout left, or the process may block on writing stdout
		_, _ = io.Copy(io.Discard, pr)
	}()

	err := cmd.Wait()
	_ = pw.Close()
	<-handled
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}
