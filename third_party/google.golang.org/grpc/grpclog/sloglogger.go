// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package grpclog defines slog-based logging for grpc.
// Importing this package will install slog as the logger used by grpclog.
package grpclog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	runtime_ "github.com/searKing/golang/go/runtime"
	"google.golang.org/grpc/grpclog"
)

var _ grpclog.Logger = (*slogger)(nil)
var _ grpclog.LoggerV2 = (*slogger)(nil)
var _ grpclog.DepthLoggerV2 = (*slogger)(nil)

const d = 3
const LevelFatal slog.Level = slog.LevelError + 4

func init() {
	grpclog.SetLoggerV2(DefaultSlogLogger())
}

// slogger adapts slog's Logger to be compatible with
// [grpclog.LoggerV2], the experimental [grpclog.DepthLoggerV2] and the deprecated [grpclog.Logger].
//
//go:generate go-option -type "slogger"
type slogger struct {
	Handler slog.Handler
}

// DefaultSlogLogger wraps [slog.Default]'s [slog.Handler] to be compatible with
// [grpclog.LoggerV2], the experimental [grpclog.DepthLoggerV2] and the deprecated [grpclog.Logger].
func DefaultSlogLogger() *slogger {
	return NewSlogger(nil)
}

// NewSlogger wraps a slog's Logger to be compatible with
// [grpclog.LoggerV2], the experimental [grpclog.DepthLoggerV2] and the deprecated [grpclog.Logger].
func NewSlogger(logger slog.Handler, opts ...SloggerOption) *slogger {
	var log = slogger{Handler: logger}
	log.ApplyOptions(opts...)
	return &log
}

func (g *slogger) log(ctx context.Context, level slog.Level, args ...any) {
	h := g.Handler
	if h == nil {
		h = slog.Default().Handler()
	}
	if h.Enabled(ctx, level) {
		pc := runtime_.GetCallerFrame(d + 1).PC
		r := slog.NewRecord(time.Now(), level, fmt.Sprint(args...), pc)
		_ = h.WithGroup("grpc").Handle(ctx, r)
	}
}

func (g *slogger) logln(ctx context.Context, level slog.Level, args ...any) {
	h := g.Handler
	if h == nil {
		h = slog.Default().Handler()
	}
	if h.Enabled(ctx, level) {
		pc := runtime_.GetCallerFrame(d + 1).PC
		r := slog.NewRecord(time.Now(), level, fmt.Sprintln(args...), pc)
		_ = h.WithGroup("grpc").Handle(ctx, r)
	}
}

func (g *slogger) logf(ctx context.Context, level slog.Level, format string, args ...any) {
	h := g.Handler
	if h == nil {
		h = slog.Default().Handler()
	}
	if h.Enabled(ctx, level) {
		pc := runtime_.GetCallerFrame(d + 1).PC
		r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pc)
		_ = h.WithGroup("grpc").Handle(ctx, r)
	}
}

func (g *slogger) logDepth(ctx context.Context, level slog.Level, depth int, args ...any) {
	h := g.Handler
	if h == nil {
		h = slog.Default().Handler()
	}
	if h.Enabled(ctx, level) {
		pc := runtime_.GetCallerFrame(d + depth).PC
		r := slog.NewRecord(time.Now(), level, fmt.Sprint(args...), pc)
		_ = h.WithGroup("grpc").Handle(ctx, r)
	}
}

// Info implements grpclog.Entry.LoggerV2's Info.
func (g *slogger) Info(args ...any) {
	g.log(context.Background(), slog.LevelInfo, args...)
}

// Infoln implements grpclog.Entry.LoggerV2's Infoln.
func (g *slogger) Infoln(args ...any) {
	g.logln(context.Background(), slog.LevelInfo, args...)
}

// Infof implements grpclog.Entry.LoggerV2's Infof.
func (g *slogger) Infof(format string, args ...any) {
	g.logf(context.Background(), slog.LevelInfo, format, args...)
}

// InfoDepth implements grpclog.Entry.LoggerV2's DebugDepth.
// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth(0, "msg") is the same as Info("msg").
func (g *slogger) InfoDepth(depth int, args ...any) {
	g.logDepth(context.Background(), slog.LevelInfo, depth+1, args...)
}

// Warning implements grpclog.Entry.LoggerV2's Warn.
func (g *slogger) Warning(args ...any) {
	g.log(context.Background(), slog.LevelWarn, args...)
}

// Warningln implements grpclog.Entry.LoggerV2's Warnln.
func (g *slogger) Warningln(args ...any) {
	g.logln(context.Background(), slog.LevelWarn, args...)
}

// Warningf implements grpclog.Entry.LoggerV2's Warnf.
func (g *slogger) Warningf(format string, args ...any) {
	g.logf(context.Background(), slog.LevelWarn, format, args...)
}

// WarningDepth acts as Warn but uses depth to determine which call frame to log.
// WarningDepth(0, "msg") is the same as Warn("msg").
func (g *slogger) WarningDepth(depth int, args ...any) {
	g.logDepth(context.Background(), slog.LevelWarn, depth+1, args...)
}

// Error implements grpclog.Entry.LoggerV2's Error.
func (g *slogger) Error(args ...any) {
	g.log(context.Background(), slog.LevelError, args...)
}

// Errorln implements grpclog.Entry.LoggerV2's Errorln.
func (g *slogger) Errorln(args ...any) {
	g.logln(context.Background(), slog.LevelError, args...)
}

// Errorf implements grpclog.Entry.LoggerV2's Errorf.
func (g *slogger) Errorf(format string, args ...any) {
	g.logf(context.Background(), slog.LevelError, format, args...)
}

// ErrorDepth acts as Warn but uses depth to determine which call frame to log.
// ErrorDepth(0, "msg") is the same as Error("msg").
func (g *slogger) ErrorDepth(depth int, args ...any) {
	g.logDepth(context.Background(), slog.LevelError, depth+1, args...)
}

// Fatal implements grpclog.Entry.LoggerV2's Fatal.
func (g *slogger) Fatal(args ...any) {
	g.log(context.Background(), LevelFatal, args...)
	os.Exit(1)
}

// Fatalln implements grpclog.Entry.LoggerV2's Fatalln.
func (g *slogger) Fatalln(args ...any) {
	g.logln(context.Background(), LevelFatal, args...)
	os.Exit(1)
}

// Fatalf implements grpclog.Entry.LoggerV2's Fatalf.
func (g *slogger) Fatalf(format string, args ...any) {
	g.logf(context.Background(), LevelFatal, format, args...)
	os.Exit(1)
}

// FatalDepth acts as Warn but uses depth to determine which call frame to log.
// FatalDepth(0, "msg") is the same as Fatal("msg").
func (g *slogger) FatalDepth(depth int, args ...any) {
	g.logDepth(context.Background(), LevelFatal, depth+1, args...)
	os.Exit(1)
}

// V implements grpclog.LoggerV2.
func (g *slogger) V(l int) bool {
	h := g.Handler
	if h == nil {
		h = slog.Default().Handler()
	}
	// grpc's verbose level is always 2 thus far.
	return h.Enabled(context.Background(), slog.LevelDebug+slog.Level(l))
}

// Print implements grpclog.Logger's Print.
func (g *slogger) Print(args ...any) {
	g.log(context.Background(), slog.LevelInfo, args...)
}

// Printf implements grpclog.Logger's Printf.
func (g *slogger) Printf(format string, args ...any) {
	g.logf(context.Background(), slog.LevelInfo, format, args...)
}

// Println implements grpclog.Logger's Println.
func (g *slogger) Println(args ...any) {
	g.logln(context.Background(), slog.LevelInfo, args...)
}
