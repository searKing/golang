// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"context"
	"io"
	"os/exec"
	"time"
)

func CommandWithTimeoutHandler(timeout time.Duration, name string, arg ...string) (data []byte, err error) {
	ctx, stop := context.WithTimeout(context.Background(), timeout)
	go func() {
		data, err = exec.Command(name, arg...).CombinedOutput()
		stop()
	}()
	select {
	case <-ctx.Done():
		if ctx.Err() == context.Canceled {
			return data, nil
		}
		return nil, ctx.Err()
	}
}

func CommandWithTimeout(handle func(io.Reader), timeout time.Duration, name string, arg ...string) (err error) {
	cs, err := newCommandServerWithTimeout(handle, timeout, name, arg...)
	if err != nil {
		return err
	}
	err = cs.wait()
	if err != nil {
		cs.Stop()
		return err
	}
	return nil
}
func newCommandServerWithTimeout(handle func(io.Reader), timeout time.Duration, name string, args ...string) (*commandServer, error) {
	ctx, stop := context.WithTimeout(context.Background(), timeout)
	return newCommandServer(ctx, stop, handle, name, args...)
}
