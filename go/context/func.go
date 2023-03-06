// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"context"
)

// WrapFunc returns a function literal that uses ctx to
// wait for the f parameter to be called and executed done
// or return immediately with f not terminate if ctx done.
// f will not be terminated if it's called.
func WrapFunc(ctx context.Context, f func() error) func() error {
	return func() (err error) {
		doneC := make(chan struct{})
		go func() {
			defer close(doneC)
			err = f()
		}()

		select {
		case <-doneC:
			return nil
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}
