// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type commandServer struct {
	cmd    *exec.Cmd
	handle func(reader io.Reader)
	ctx    context.Context
	done   context.CancelFunc
}

func newCommandServer(parent context.Context, stop context.CancelFunc, handle func(reader io.Reader), name string, args ...string) (*commandServer, error) {
	if parent == nil {
		parent = context.Background()
	}
	if stop == nil {
		stop = func() {}
	}
	if handle == nil {
		handle = func(reader io.Reader) {}
	}

	cs := &commandServer{
		cmd:    exec.Command(name, args...),
		handle: handle,
		ctx:    parent,
		done:   stop,
	}

	r, err := cs.cmd.StdoutPipe()
	if err != nil {
		cs.Stop()
		return nil, err
	}
	go cs.watch(r)
	return cs, nil
}

func (cs *commandServer) wait() error {
	select {
	case <-cs.ctx.Done():
		return cs.ctx.Err()
	}
	return nil
}

func (cs *commandServer) watch(r io.Reader) {
	cs.handle(r)
	cs.cmd.Wait()
	cs.done()
}
func (cs *commandServer) Stop() {
	cs.cmd.Process.Signal(os.Interrupt)
	cs.done()
}
