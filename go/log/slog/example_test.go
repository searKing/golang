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

	"github.com/searKing/golang/go/exp/types"
	"github.com/searKing/golang/go/log/slog/internal/slogtest"
)

func slogNil(name string, h slog.Handler) {
	defer func() {
		if err := recover(); err != nil {
			// in go1.21, [slog.TextHandler] and [slog.JSONHandler] will panic
			// https://github.com/golang/go/commit/73667209c1c83bd48fe7338c3b4caaa05c073202
			fmt.Printf("[slog/%s] unexpected panic: %v\n", name, err)
		}
	}()
	logger := slog.New(h)
	{
		var typedNil *text
		logger.With("attr_typed_nil", types.Any(typedNil)).Info(fmt.Sprintf("[slog/%s]", name), "args_typed_nil", types.Any(typedNil))
		logger.With("attr_typed_nil", typedNil).Info(fmt.Sprintf("[slog/%s]", name), "args_typed_nil", typedNil)
	}
}

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
		fmt.Printf("----json----\n")
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
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
	// ----json----
	// {"level":"INFO","msg":"finished","req":{"method":"GET","url":"localhost"},"status":200,"duration":1000000000,"error":"test ExampleGroup: invalid argument"}
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

func ExampleMultiHandler() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	// ...
	{
		fmt.Printf("----text----\n")
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("text message")
	}
	// ...
	{
		fmt.Printf("----json----\n")
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("json message")
	}
	// ...
	{
		fmt.Printf("----glog----\n")
		logger := slog.New(NewGlogHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("glog message")
	}
	// ...
	{
		fmt.Printf("----glog_human----\n")
		logger := slog.New(NewGlogHumanHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
		logger.Info("glog_human message")
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
		logger.Info("multi[text-json-glog-glog_human] message")
	}

	// Output:
	// ----text----
	// level=INFO msg="text message"
	// ----json----
	// {"level":"INFO","msg":"json message"}
	// ----glog----
	// I 0] glog message
	// ----glog_human----
	// [INFO ] [0] glog_human message
	// ----multi[text-json-glog-glog_human]----
	// level=INFO msg="multi[text-json-glog-glog_human] message"
	// {"level":"INFO","msg":"multi[text-json-glog-glog_human] message"}
	// I 0] multi[text-json-glog-glog_human] message
	// [INFO ] [0] multi[text-json-glog-glog_human] message
}
