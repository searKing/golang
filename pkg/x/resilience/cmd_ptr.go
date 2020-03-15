// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resilience

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

type Command struct {
	*exec.Cmd
	preHandles  []func() error
	postHandles []func(err error) error
}

func NewCommand(cmd *exec.Cmd) *Command {
	return &Command{
		Cmd: cmd,
	}
}

func (r *Command) AppendPreHandles(h ...func() error) {
	r.preHandles = append(r.preHandles, h...)
}

func (r *Command) AppendPostHandles(h ...func(err error) error) {
	r.postHandles = append(r.postHandles, h...)
}

func (r *Command) Value() interface{} {
	return r
}

func (r *Command) Ready() error {
	if r == nil || r.Cmd == nil {
		return fmt.Errorf("command: empty value")
	}
	if r.Cmd.Process == nil {
		return errors.New("command: not started")
	}

	if r.Cmd.ProcessState != nil && r.Cmd.ProcessState.Exited() {
		return fmt.Errorf("command: exited already %s", r.Cmd.ProcessState.String())
	}
	return nil
}

func (r *Command) Close() {
	if r == nil || r.Cmd == nil {
		return
	}
	proc := r.Cmd.Process
	if proc != nil {
		_ = proc.Signal(syscall.SIGTERM)
		proc.Wait()
		// proc.Kill()
		// no need to close attached log file.
		// see "Wait releases any resources associated with the cmd."
		// if closer, ok := cmd.Stdout.(io.Closer); ok {
		// 	closer.Close()
		// 	logger.Printf("process:%v Stdout closed.", proc)
		// }
	}
	r.Cmd = nil
}
func (r *Command) Run() error {
	if r == nil || r.Cmd == nil {
		return fmt.Errorf("command: empty value")
	}
	for _, h := range r.preHandles {
		if h == nil {
			continue
		}
		if err := h(); err != nil {
			return err
		}
	}

	runErr := r.Cmd.Run()
	for _, h := range r.postHandles {
		if h == nil {
			continue
		}
		if err := h(runErr); err != nil {
			return err
		}
	}
	return runErr
}
