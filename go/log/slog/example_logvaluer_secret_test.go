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

// A token is a secret value that grants permissions.
type Token string

// LogValue implements slog.LogValuer.
// It avoids revealing the token.
func (Token) LogValue() slog.Value {
	return slog.StringValue("REDACTED_TOKEN")
}

// This example demonstrates a Value that replaces itself
// with an alternative representation to avoid revealing secrets.
func ExampleLogValuer_secret() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()
	t := Token("shhhh!")
	{
		fmt.Printf("----text----\n")
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("permission granted", "user", "Perry", "token", t)
	}
	{
		fmt.Printf("----glog----\n")
		logger := slog.New(NewGlogHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("permission granted", "user", "Perry", "token", t)
	}
	{
		fmt.Printf("----glog_human----\n")
		logger := slog.New(NewGlogHumanHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("permission granted", "user", "Perry", "token", t)
	}

	// Output:
	// ----text----
	// level=INFO msg="permission granted" user=Perry token=REDACTED_TOKEN
	// ----glog----
	// I 0] permission granted, user=Perry, token=REDACTED_TOKEN
	// ----glog_human----
	// [INFO ] [0] permission granted, user=Perry, token=REDACTED_TOKEN
}
