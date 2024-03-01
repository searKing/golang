// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.22

package slog

import (
	"log/slog"
	"os"

	"github.com/searKing/golang/go/log/slog/internal/slogtest"
)

func ExampleTypedNil() {
	getPid = func() int { return 0 } // set pid to zero for test

	defer func() { getPid = os.Getpid }()
	slogNil("text", slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
	slogNil("json", slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
	slogNil("glog", NewGlogHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))
	slogNil("glog_human", NewGlogHumanHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}))

	// Output:
	// level=INFO msg=[slog/text] attr_typed_nil=<nil> args_typed_nil=<nil>
	// level=INFO msg=[slog/text] attr_typed_nil=<nil> args_typed_nil=<nil>
	// {"level":"INFO","msg":"[slog/json]","attr_typed_nil":null,"args_typed_nil":null}
	// {"level":"INFO","msg":"[slog/json]","attr_typed_nil":"<nil>","args_typed_nil":"<nil>"}
	// I 0] [slog/glog], attr_typed_nil=<nil>, args_typed_nil=<nil>
	// I 0] [slog/glog], attr_typed_nil=<nil>, args_typed_nil=<nil>
	// [INFO ] [0] [slog/glog_human], attr_typed_nil=<nil>, args_typed_nil=<nil>
	// [INFO ] [0] [slog/glog_human], attr_typed_nil=<nil>, args_typed_nil=<nil>
}
