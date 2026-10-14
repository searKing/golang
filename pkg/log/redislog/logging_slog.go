// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redislog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// WithSlogLogger returns a redis.Hook that logs dials, commands and pipelines of go-redis clients to h,
// by NewHook with config, and logs internal logs of go-redis to h, by redis.SetLogger globally.
func WithSlogLogger(h slog.Handler, config Config) redis.Hook {
	l := slog.New(h)
	redis.SetLogger(InternalLogger(l))
	return NewHook(l, config)
}

// PrintfFunc is an adapter to allow the use of ordinary functions as the logger of go-redis set by redis.SetLogger.
type PrintfFunc func(ctx context.Context, format string, v ...any)

func (f PrintfFunc) Printf(ctx context.Context, format string, v ...any) { f(ctx, format, v...) }

// InternalLogger returns the logger of go-redis to be set by redis.SetLogger, which logs internal logs of go-redis,
// such as failures of the connection pool and reconnections, to l, slog.Default() if nil,
// at warn level, sourced from go-redis where it logs.
// go-redis logs without levels, but filters internal logs by its own level, error by default,
// which is set by logging.SetLogLevel of github.com/redis/go-redis/v9/logging.
func InternalLogger(l *slog.Logger) PrintfFunc {
	return func(ctx context.Context, format string, v ...any) {
		logger := l
		if logger == nil {
			logger = slog.Default()
		}
		if ctx == nil {
			ctx = context.Background()
		}
		if !logger.Enabled(ctx, slog.LevelWarn) {
			return
		}
		r := slog.NewRecord(time.Now(), slog.LevelWarn, strings.TrimSpace(fmt.Sprintf(format, v...)), callerPC(isOfThisPackage))
		_ = logger.Handler().Handle(ctx, r)
	}
}
