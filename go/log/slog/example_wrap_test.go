// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Infof is an example of a user-defined logging function that wraps slog.
// The log record contains the source position of the caller of Infof.
func Infof(logger *slog.Logger, format string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

func Example_wrapping() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove time.
		if a.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}
	{
		fmt.Printf("----text----\n")
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}))
		Infof(logger, "message, %s", "formatted")
	}
	{
		fmt.Printf("----glog----\n")
		logger := slog.New(NewGlogHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}))
		Infof(logger, "message, %s", "formatted")
	}
	{
		fmt.Printf("----glog_human----\n")
		logger := slog.New(NewGlogHumanHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}))
		Infof(logger, "message, %s", "formatted")
	}

	// Output:
	// ----text----
	// level=INFO source=example_wrap_test.go:47 msg="message, formatted"
	// ----glog----
	// I 0 example_wrap_test.go:52] message, formatted
	// ----glog_human----
	// [INFO ] [0] [example_wrap_test.go:57](Example_wrapping) message, formatted

}
