// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resilience

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type CommandSharedPtr struct {
	*SharedPtr
	preHandles  []func() error
	postHandles []func(err error) error
}

func NewCommandSharedPtr(ctx context.Context, cmd func() *exec.Cmd,
	preHandles []func() error, postHandles []func(err error) error,
	l logrus.FieldLogger) *CommandSharedPtr {
	resilienceSharedPtr := &CommandSharedPtr{
		SharedPtr: NewSharedPtr(ctx, func() (Ptr, error) {
			if cmd == nil {
				return nil, fmt.Errorf("resillence cmd: empty value")
			}
			cmder := NewCommand(cmd())
			cmder.AppendPreHandles(preHandles...)
			cmder.AppendPostHandles(postHandles...)
			return cmder, nil
		}, l),
	}
	return resilienceSharedPtr
}

func (g *CommandSharedPtr) AppendPreHandles(h ...func() error) {
	g.preHandles = append(g.preHandles, h...)
}

func (g *CommandSharedPtr) AppendPostHandles(h ...func() error) {
	g.preHandles = append(g.preHandles, h...)
}

func (g *CommandSharedPtr) GetUntilReady() (*Command, error) {
	x, err := g.SharedPtr.GetUntilReady()
	if err != nil {
		return nil, err
	}
	ffmpeg, ok := x.Value().(*Command)
	if ok {
		return ffmpeg, nil
	}
	return nil, fmt.Errorf("unexpected type %T", x)
}
func (g *CommandSharedPtr) GetWithRetry() (*Command, error) {
	x, err := g.SharedPtr.GetWithRetry()
	if err != nil {
		return nil, err
	}
	cmd, ok := x.Value().(*Command)
	if ok {
		return cmd, nil
	}
	return nil, fmt.Errorf("unexpected type %T", x)
}
func (g *CommandSharedPtr) Get() (*Command, error) {
	x := g.SharedPtr.Get()
	if x == nil {
		return nil, nil
	}
	ffmpeg, ok := x.Value().(*Command)
	if ok {
		return ffmpeg, nil
	}
	return nil, fmt.Errorf("unexpected type %T", x)
}

func (g *CommandSharedPtr) Run() error {
	cmd, err := g.Get()
	if err != nil {
		return err
	}
	return cmd.Run()
}
