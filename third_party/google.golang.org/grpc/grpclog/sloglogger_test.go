// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpclog

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"regexp"
	"testing"

	"google.golang.org/grpc/grpclog"
)

func TestWithSloggerHandler(t *testing.T) {
	var buf bytes.Buffer
	grpclog.SetLoggerV2(NewSlogger(slog.New(slog.NewTextHandler(&buf, nil)).Handler()))

	grpclog.Info(slog.LevelInfo.String())
	grpclog.Warning(slog.LevelWarn.String())
	grpclog.Error(slog.LevelError.String())
	s := bufio.NewScanner(&buf)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		// The content of info buffer should be something like:
		// time=2023-08-14T16:09:23.360+08:00 level=INFO msg=INFO
		// time=2023-08-14T16:09:23.360+08:00 level=WARN msg=WARN
		// time=2023-08-14T16:09:23.360+08:00 level=ERROR msg=ERROR
		if err := checkLogForSeverity(s.Bytes()); err != nil {
			t.Fatal(err)
		}
	}
}

// check if b is in the format of:
//
//	2017/04/07 14:55:42 WARNING: WARNING
func checkLogForSeverity(b []byte) error {
	expected := regexp.MustCompile(`^time=[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3}\+[0-9]{2}:[0-9]{2} level=(INFO|WARN|ERROR) msg=(INFO|WARN|ERROR)$`)
	if m := expected.Match(b); !m {
		return fmt.Errorf("got: %q, want string in format of: %q", b, "2023/08/14 15:41:14 INFO INFO")
	}
	return nil
}
