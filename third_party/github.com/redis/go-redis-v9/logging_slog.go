// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/redis/go-redis/v9"
)

// WithSlogLogger returns a redis.Hook that logs dials, commands and pipelines of go-redis clients to h,
// customized by opts, and logs internal logs of go-redis to h, by redis.SetLogger globally.
func WithSlogLogger(h slog.Handler, opts ...LoggingOption) redis.Hook {
	l := slog.New(h)
	redis.SetLogger(InternalLogger(l))
	return Hook(l, opts...)
}

// PrintfFunc is an adapter to allow the use of ordinary functions as the logger of go-redis set by redis.SetLogger.
type PrintfFunc func(ctx context.Context, format string, v ...any)

func (f PrintfFunc) Printf(ctx context.Context, format string, v ...any) { f(ctx, format, v...) }

// InternalLogger returns the logger of go-redis to be set by redis.SetLogger, which logs internal logs of go-redis,
// such as failures of the connection pool and reconnections, to l, slog.Default() if nil,
// at warn level, with fields of SystemTag.
func InternalLogger(l *slog.Logger) PrintfFunc {
	return func(ctx context.Context, format string, v ...any) {
		logger := l
		if logger == nil {
			logger = slog.Default()
		}
		logger.Log(ctx, slog.LevelWarn, strings.TrimSpace(fmt.Sprintf(format, v...)), SystemTag[0], SystemTag[1])
	}
}
