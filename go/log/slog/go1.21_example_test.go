// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.21 && !go1.22

package slog

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/searKing/golang/go/log/slog/internal/slogtest"
)

func TestTypedNil(t *testing.T) {
	var gots [4]bytes.Buffer
	tests := []struct {
		name    string
		handler slog.Handler
		want    string
	}{
		{
			name:    "text",
			handler: slog.NewTextHandler(&gots[0], &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}),
			want: `level=INFO msg=[slog/text] attr_typed_nil=<nil> args_typed_nil=<nil>
[slog/text] unexpected panic: value method github.com/searKing/golang/go/log/slog.text.MarshalText called using nil *text pointer
`,
		},
		{
			name:    "json",
			handler: slog.NewJSONHandler(&gots[1], &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}),
			want: `{"level":"INFO","msg":"[slog/json]","attr_typed_nil":null,"args_typed_nil":null}
[slog/json] unexpected panic: value method github.com/searKing/golang/go/log/slog.text.Error called using nil *text pointer
`,
		},
		{
			name:    "glog",
			handler: NewGlogHandler(&gots[2], &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}),
			want: `I 0] [slog/glog], attr_typed_nil=<nil>, args_typed_nil=<nil>
I 0] [slog/glog], attr_typed_nil=<nil>, args_typed_nil=<nil>
`,
		},
		{
			name:    "glog_human",
			handler: NewGlogHumanHandler(&gots[3], &slog.HandlerOptions{ReplaceAttr: slogtest.RemoveTime}),
			want: `[INFO ] [0] [slog/glog_human], attr_typed_nil=<nil>, args_typed_nil=<nil>
[INFO ] [0] [slog/glog_human], attr_typed_nil=<nil>, args_typed_nil=<nil>
`,
		},
	}

	getPid = func() int { return 0 } // 设置pid为0用于测试
	defer func() { getPid = os.Getpid }()

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slogNil(tt.name, tt.handler)
			got := gots[i].String()
			if got != tt.want {
				t.Errorf("#%d, got %q, want %q", i, got, tt.want)
			}
		})
	}
}
