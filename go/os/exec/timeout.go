// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"context"
	"errors"
	"io"
	"time"
)

// CommandWithTimeoutHandler runs the named program and returns its combined stdout and stderr.
// The process is sent SIGTERM if it does not complete within timeout, and killed if it does not exit
// within one more second, in which case nil and [context.DeadlineExceeded] are returned.
func CommandWithTimeoutHandler(timeout time.Duration, name string, arg ...string) (data []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	data, err = command(ctx, name, arg...).CombinedOutput()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, ctx.Err()
	}
	return data, err
}

// CommandWithTimeout runs the named program, feeding its stdout to handle.
// The process is sent SIGTERM if it does not complete within timeout, and killed if it does not exit
// within one more second, in which case [context.DeadlineExceeded] is returned.
func CommandWithTimeout(handle func(io.Reader), timeout time.Duration, name string, arg ...string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return commandContext(ctx, handle, name, arg...)
}
