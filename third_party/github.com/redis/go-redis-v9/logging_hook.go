// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ redis.Hook = (*hook)(nil)

// NewHook returns a redis.Hook that logs dials, commands and pipelines of go-redis clients to l,
// slog.Default() if nil, as the slog logger of gorm logs SQL, with a command in place of SQL.
//
// Each call is logged in a single line, leveled by Config.Levels, with fields grouped by "redis":
//   - cmd: the command as a line of redis-cli, such as `set k1 "hello world"`, which can be run by redis-cli as is,
//     or commands of pipelines as lines, with arguments which may carry credentials redacted,
//     such as passwords of AUTH, and truncated as SLOWLOG of redis does;
//   - network and addr: the network and address dialed, of dials only;
//   - duration: the duration of the call, such as 0.123ms;
//   - slow: true if the call is slower than Config.SlowThreshold, left out otherwise;
//   - error: the error of the call, left out if succeeded.
func NewHook(l *slog.Logger, config Config) redis.Hook {
	if config.Levels == nil {
		config.Levels = DefaultErrorToLevel
	}
	return &hook{logger: l, config: config}
}

type hook struct {
	logger *slog.Logger
	config Config
}

func (h *hook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		start := time.Now()
		conn, err := next(ctx, network, addr)
		h.log(ctx, "redis dialed", start, false, err, func() []slog.Attr {
			return []slog.Attr{slog.String("network", network), slog.String("addr", addr)}
		})
		return conn, err
	}
}

func (h *hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		h.log(ctx, "redis executed", start, isBlocking(cmd), err, func() []slog.Attr {
			return []slog.Attr{slog.String("cmd", h.command(cmd))}
		})
		return err
	}
}

func (h *hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		// The pipeline returns the first error of its commands, such as redis.Nil if any command misses its key,
		// while the rest are executed as usual.
		err := next(ctx, cmds)
		h.log(ctx, "redis pipeline executed", start, slices.ContainsFunc(cmds, isBlocking), err, func() []slog.Attr {
			lines := make([]string, 0, len(cmds))
			for _, cmd := range cmds {
				lines = append(lines, h.command(cmd))
			}
			return []slog.Attr{slog.Any("cmd", truncateLines(lines))}
		})
		return err
	}
}

// command returns cmd as a line of redis-cli, parameterized if Config.Parameterized.
func (h *hook) command(cmd redis.Cmder) string {
	if h.config.Parameterized {
		return commandStatement(cmd)
	}
	return commandArgs(cmd)
}

// log logs the call started at start with err, and fields of attrs, which are not built if the level is disabled.
// Calls blocking by design, such as BLPOP, are never slow.
func (h *hook) log(ctx context.Context, msg string, start time.Time, blocking bool, err error, attrs func() []slog.Attr) {
	duration := time.Since(start)
	slow := h.config.SlowThreshold > 0 && duration > h.config.SlowThreshold && !blocking
	level := h.config.Levels(err)
	if slow {
		level = max(level, slog.LevelWarn)
	}
	logger := h.logger
	if logger == nil {
		logger = slog.Default()
	}
	if !logger.Enabled(ctx, level) {
		return
	}

	fields := append(attrs(), slog.String("duration", fmt.Sprintf("%.3fms", float64(duration.Nanoseconds())/1e6)))
	if slow {
		fields = append(fields, slog.Bool("slow", true))
	}
	if err != nil {
		fields = append(fields, slog.String("error", err.Error()))
	}
	record := slog.NewRecord(time.Now(), level, msg, callerPC(isOfRedis))
	if h.config.FieldsFromContext != nil {
		record.Add(h.config.FieldsFromContext(ctx)...)
	}
	record.AddAttrs(slog.Attr{Key: "redis", Value: slog.GroupValue(fields...)})
	_ = logger.Handler().Handle(ctx, record)
}

// callerPC returns the program counter of the first caller outside frames skipped by skip,
// for the source of logs, as Logger.Log does for its caller,
// or 0 if none, such as dials of connection pools of go-redis in goroutines of their own.
func callerPC(skip func(function string) bool) uintptr {
	var pcs [64]uintptr
	// Skip runtime.Callers, callerPC and the caller of callerPC.
	frames := runtime.CallersFrames(pcs[:runtime.Callers(3, pcs[:])])
	for {
		f, more := frames.Next()
		if !skip(f.Function) && !strings.HasPrefix(f.Function, "runtime.") {
			// The PC of the frame is of the call, while the record takes the return address as runtime.Callers,
			// which works for frames of inlined calls too, of which the PC is of the inlined call.
			return f.PC + 1
		}
		if !more {
			return 0
		}
	}
}

// isOfRedis reports whether the function is of go-redis, including its extra packages, such as redisotel,
// or of this package, skipped to locate the caller of redis, such as the caller of rdb.Get,
// as redisotel does for code.function of spans.
func isOfRedis(function string) bool {
	return strings.HasPrefix(function, "github.com/redis/go-redis/") || isOfThisPackage(function)
}

// isOfThisPackage reports whether the function is of this package, excluding its subpackages.
func isOfThisPackage(function string) bool {
	rest, ok := strings.CutPrefix(function, thisPackage)
	return ok && strings.HasPrefix(rest, ".")
}

// thisPackage is the import path of this package.
var thisPackage = reflect.TypeFor[hook]().PkgPath()
