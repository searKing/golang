// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logruslogger defines logrus-based logging for grpc.
// Importing this package will install logrus as the logger used by grpclog.
// Attention, Info -> Debug to mute verbose messages.
package logruslogger

import (
	runtime_ "github.com/searKing/golang/go/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

const d = 3

func init() {
	grpclog.SetLoggerV2(NewLogger(logrus.NewEntry(logrus.StandardLogger())))
}

// Logger adapts logrus's Logger to be compatible with [grpclog.LoggerV2], the experimental [grpclog.DepthLoggerV2] and the deprecated [grpclog.Logger].
//
//go:generate go-option -type "Logger"
type Logger struct {
	Entry *logrus.Entry

	verbose         int
	LevelTranslator func(level logrus.Level) logrus.Level
}

func NewLogger(logger *logrus.Entry, opts ...LoggerOption) *Logger {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	var log = Logger{Entry: logger}
	log.LevelTranslator = func(level logrus.Level) logrus.Level {
		//if level == logrus.InfoLevel {
		//	return logrus.DebugLevel
		//}
		return level
	}
	log.ApplyOptions(opts...)
	return &log
}

func (g *Logger) log(level logrus.Level, args ...any) {
	level = g.LevelTranslator(level)
	if g.Entry.Logger.IsLevelEnabled(level) && g.Entry.Logger.ReportCaller && g.Entry.Caller == nil {
		g.Entry.Caller = runtime_.GetCallerFrame(d + 1)
	}
	g.Entry.Log(g.LevelTranslator(level), args...)
}

func (g *Logger) logln(level logrus.Level, args ...any) {
	level = g.LevelTranslator(level)
	if g.Entry.Logger.IsLevelEnabled(level) && g.Entry.Logger.ReportCaller && g.Entry.Caller == nil {
		g.Entry.Caller = runtime_.GetCallerFrame(d + 1)
	}
	g.Entry.Logln(g.LevelTranslator(level), args...)
}

func (g *Logger) logf(level logrus.Level, format string, args ...any) {
	level = g.LevelTranslator(level)
	if g.Entry.Logger.IsLevelEnabled(level) && g.Entry.Logger.ReportCaller && g.Entry.Caller == nil {
		g.Entry.Caller = runtime_.GetCallerFrame(d + 1)
	}
	g.Entry.Logf(g.LevelTranslator(level), format, args...)
}

func (g *Logger) logDepth(level logrus.Level, depth int, args ...any) {
	level = g.LevelTranslator(level)
	if g.Entry.Logger.IsLevelEnabled(level) && g.Entry.Logger.ReportCaller && g.Entry.Caller == nil {
		g.Entry.Caller = runtime_.GetCallerFrame(depth + d)
	}
	g.Entry.Log(g.LevelTranslator(level), args...)
}

// Info implements grpclog.Entry.LoggerV2's Info.
func (g *Logger) Info(args ...any) {
	g.log(logrus.InfoLevel, args...)
}

// Infoln implements grpclog.Entry.LoggerV2's Infoln.
func (g *Logger) Infoln(args ...any) {
	g.logln(logrus.InfoLevel, args...)
}

// Infof implements grpclog.Entry.LoggerV2's Infof.
func (g *Logger) Infof(format string, args ...any) {
	g.logf(logrus.InfoLevel, format, args...)
}

// InfoDepth implements grpclog.Entry.LoggerV2's DebugDepth.
// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth(0, "msg") is the same as Info("msg").
func (g *Logger) InfoDepth(depth int, args ...any) {
	g.logDepth(logrus.InfoLevel, depth+1, args...)
}

// Warning implements grpclog.Entry.LoggerV2's Warn.
func (g *Logger) Warning(args ...any) {
	g.log(logrus.WarnLevel, args...)
}

// Warningln implements grpclog.Entry.LoggerV2's Warnln.
func (g *Logger) Warningln(args ...any) {
	g.logln(logrus.WarnLevel, args...)
}

// Warningf implements grpclog.Entry.LoggerV2's Warnf.
func (g *Logger) Warningf(format string, args ...any) {
	g.logf(logrus.WarnLevel, format, args...)
}

// WarningDepth acts as Warn but uses depth to determine which call frame to log.
// WarningDepth(0, "msg") is the same as Warn("msg").
func (g *Logger) WarningDepth(depth int, args ...any) {
	g.logDepth(logrus.WarnLevel, depth+1, args...)
}

// Error implements grpclog.Entry.LoggerV2's Error.
func (g *Logger) Error(args ...any) {
	g.log(logrus.ErrorLevel, args...)
}

// Errorln implements grpclog.Entry.LoggerV2's Errorln.
func (g *Logger) Errorln(args ...any) {
	g.logln(logrus.ErrorLevel, args...)
}

// Errorf implements grpclog.Entry.LoggerV2's Errorf.
func (g *Logger) Errorf(format string, args ...any) {
	g.logf(logrus.ErrorLevel, format, args...)
}

// ErrorDepth acts as Warn but uses depth to determine which call frame to log.
// ErrorDepth(0, "msg") is the same as Error("msg").
func (g *Logger) ErrorDepth(depth int, args ...any) {
	g.logDepth(logrus.ErrorLevel, depth+1, args...)
}

// Fatal implements grpclog.Entry.LoggerV2's Fatal.
func (g *Logger) Fatal(args ...any) {
	g.log(logrus.FatalLevel, args...)
}

// Fatalln implements grpclog.Entry.LoggerV2's Fatalln.
func (g *Logger) Fatalln(args ...any) {
	g.logln(logrus.FatalLevel, args...)
}

// Fatalf implements grpclog.Entry.LoggerV2's Fatalf.
func (g *Logger) Fatalf(format string, args ...any) {
	g.logf(logrus.FatalLevel, format, args...)
}

// FatalDepth acts as Warn but uses depth to determine which call frame to log.
// FatalDepth(0, "msg") is the same as Fatal("msg").
func (g *Logger) FatalDepth(depth int, args ...any) {
	g.logDepth(logrus.FatalLevel, depth+1, args...)
}

// Log implements grpclog.Entry.LoggerV2's Log.
func (g *Logger) Log(level logrus.Level, args ...any) {
	g.log(level, args...)
}

// Logln implements grpclog.Entry.LoggerV2's Logln.
func (g *Logger) Logln(level logrus.Level, args ...any) {
	g.logln(level, args...)
}

// Logf implements grpclog.Entry.LoggerV2's Logf.
func (g *Logger) Logf(level logrus.Level, format string, args ...any) {
	g.logf(level, format, args...)
}

// LogDepth acts as Warn but uses depth to determine which call frame to log.
// LogDepth(0, "msg") is the same as Log("msg").
func (g *Logger) LogDepth(level logrus.Level, depth int, args ...any) {
	g.logDepth(level, depth+1, args...)
}

// V implements grpclog.Entry.LoggerV2.
func (g *Logger) V(l int) bool {
	return l <= g.verbose
}
