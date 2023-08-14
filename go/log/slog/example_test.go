// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/searKing/golang/go/log/slog/internal/slogtest"
)

func ExampleGroup() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	r, _ := http.NewRequest("GET", "localhost", nil)
	err := &os.PathError{
		Op:   "test",
		Path: "ExampleGroup",
		Err:  os.ErrInvalid,
	}
	// ...
	{
		fmt.Printf("----text----\n")
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("finished",
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String())),
			slog.Int("status", http.StatusOK),
			slog.Duration("duration", time.Second),
			Error(err))
	}
	// ...
	{
		fmt.Printf("----glog----\n")
		logger := slog.New(NewGlogHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("finished",
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String())),
			slog.Int("status", http.StatusOK),
			slog.Duration("duration", time.Second),
			Error(err))
	}
	// ...
	{
		fmt.Printf("----glog_human----\n")
		logger := slog.New(NewGlogHumanHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("finished",
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String())),
			slog.Int("status", http.StatusOK),
			slog.Duration("duration", time.Second),
			Error(err))
	}
	// ...
	{
		fmt.Printf("----multi[text-json-glog-glog_human]----\n")
		logger := slog.New(MultiHandler(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}),
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}),
			NewGlogHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}),
			NewGlogHumanHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}),
		))
		logger.Info("finished",
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String())),
			slog.Int("status", http.StatusOK),
			slog.Duration("duration", time.Second),
			Error(err))
	}

	// Output:
	// ----text----
	// level=INFO msg=finished req.method=GET req.url=localhost status=200 duration=1s error="test ExampleGroup: invalid argument"
	// ----glog----
	// I 0] finished, req.method=GET, req.url=localhost, status=200, duration=1s, error=test ExampleGroup: invalid argument
	// ----glog_human----
	// [INFO ] [0] finished, req.method=GET, req.url=localhost, status=200, duration=1s, error=test ExampleGroup: invalid argument
	// ----multi[text-json-glog-glog_human]----
	// level=INFO msg=finished req.method=GET req.url=localhost status=200 duration=1s error="test ExampleGroup: invalid argument"
	// {"level":"INFO","msg":"finished","req":{"method":"GET","url":"localhost"},"status":200,"duration":1000000000,"error":"test ExampleGroup: invalid argument"}
	// I 0] finished, req.method=GET, req.url=localhost, status=200, duration=1s, error=test ExampleGroup: invalid argument
	// [INFO ] [0] finished, req.method=GET, req.url=localhost, status=200, duration=1s, error=test ExampleGroup: invalid argument
}
