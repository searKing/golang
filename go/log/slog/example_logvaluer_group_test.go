// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/searKing/golang/go/log/slog/internal/slogtest"
)

type Name struct {
	First, Last string
}

// LogValue implements slog.LogValuer.
// It returns a group containing the fields of
// the Name, so that they appear together in the log output.
func (n Name) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("first", n.First),
		slog.String("last", n.Last))
}

func ExampleLogValuer_group() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()
	n := Name{"Perry", "Platypus"}

	{
		fmt.Printf("----text----\n")
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("mission accomplished", "agent", n)
	}
	{
		fmt.Printf("----glog----\n")
		logger := slog.New(NewGlogHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("mission accomplished", "agent", n)
	}
	{
		fmt.Printf("----glog_human----\n")
		logger := slog.New(NewGlogHumanHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("mission accomplished", "agent", n)
	}

	// Output:
	// ----text----
	// level=INFO msg="mission accomplished" agent.first=Perry agent.last=Platypus
	// ----glog----
	// I 0] mission accomplished, agent.first=Perry, agent.last=Platypus
	// ----glog_human----
	// [INFO ] [0] mission accomplished, agent.first=Perry, agent.last=Platypus
}
