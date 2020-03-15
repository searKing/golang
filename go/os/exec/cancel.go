// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"context"
	"io"
)

func CommandWithCancel(handle func(reader io.Reader), name string, arg ...string) (err error) {
	cs, err := newCommandServerWithCancel(handle, name, arg...)
	err = cs.wait()
	if err != nil {
		cs.Stop()
		return err
	}
	return nil
}
func newCommandServerWithCancel(handle func(reader io.Reader), name string, arg ...string) (*commandServer, error) {
	ctx, stop := context.WithCancel(context.Background())
	return newCommandServer(ctx, stop, handle, name, arg...)
}
