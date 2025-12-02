// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	slog_ "github.com/searKing/golang/go/log/slog"
)

func TestDynamicLevelHandler(t *testing.T) {
	tests := []struct {
		name         string
		initialLevel slog.Level
		logLevel     slog.Level
		message      string
		shouldLog    bool
	}{
		{
			name:         "log info when level is info",
			initialLevel: slog.LevelInfo,
			logLevel:     slog.LevelInfo,
			message:      "info message",
			shouldLog:    true,
		},
		{
			name:         "log warn when level is info",
			initialLevel: slog.LevelInfo,
			logLevel:     slog.LevelWarn,
			message:      "warn message",
			shouldLog:    true,
		},
		{
			name:         "log error when level is info",
			initialLevel: slog.LevelInfo,
			logLevel:     slog.LevelError,
			message:      "error message",
			shouldLog:    true,
		},
		{
			name:         "not log debug when level is info",
			initialLevel: slog.LevelInfo,
			logLevel:     slog.LevelDebug,
			message:      "debug message",
			shouldLog:    false,
		},
		{
			name:         "log debug when level is debug",
			initialLevel: slog.LevelDebug,
			logLevel:     slog.LevelDebug,
			message:      "debug message",
			shouldLog:    true,
		},
		{
			name:         "not log info when level is warn",
			initialLevel: slog.LevelWarn,
			logLevel:     slog.LevelInfo,
			message:      "info message",
			shouldLog:    false,
		},
		{
			name:         "not log warn when level is error",
			initialLevel: slog.LevelError,
			logLevel:     slog.LevelWarn,
			message:      "warn message",
			shouldLog:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
				Level: slog.LevelDebug, // 基础 handler 设置为最低级别，让 DynamicLevelHandler 控制
			})

			// 使用函数返回固定级别
			getLevel := func(ctx context.Context) slog.Level {
				return tt.initialLevel
			}
			handler := slog_.DynamicLevelHandler(getLevel, baseHandler)
			logger := slog.New(handler)

			// 根据日志级别记录日志
			switch tt.logLevel {
			case slog.LevelDebug:
				logger.Debug(tt.message)
			case slog.LevelInfo:
				logger.Info(tt.message)
			case slog.LevelWarn:
				logger.Warn(tt.message)
			case slog.LevelError:
				logger.Error(tt.message)
			}

			output := buf.String()
			hasOutput := len(output) > 0 && bytes.Contains(buf.Bytes(), []byte(tt.message))

			if hasOutput != tt.shouldLog {
				t.Errorf("expected shouldLog=%v, but got hasOutput=%v, output: %s", tt.shouldLog, hasOutput, output)
			}
		})
	}
}

func TestDynamicLevelHandler_Enabled(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	getLevel := func(ctx context.Context) slog.Level {
		return slog.LevelInfo
	}
	handler := slog_.DynamicLevelHandler(getLevel, baseHandler)
	ctx := context.Background()

	tests := []struct {
		name    string
		level   slog.Level
		enabled bool
	}{
		{"debug should be disabled", slog.LevelDebug, false},
		{"info should be enabled", slog.LevelInfo, true},
		{"warn should be enabled", slog.LevelWarn, true},
		{"error should be enabled", slog.LevelError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enabled := handler.Enabled(ctx, tt.level)
			if enabled != tt.enabled {
				t.Errorf("Enabled(%v) = %v, want %v", tt.level, enabled, tt.enabled)
			}
		})
	}
}

func TestDynamicLevelHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	getLevel := func(ctx context.Context) slog.Level {
		return slog.LevelInfo
	}
	handler := slog_.DynamicLevelHandler(getLevel, baseHandler)
	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("key", "value"),
	})

	logger := slog.New(handlerWithAttrs)
	logger.Info("test message")

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("key=value")) {
		t.Errorf("expected output to contain 'key=value', got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("test message")) {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
}

func TestDynamicLevelHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	getLevel := func(ctx context.Context) slog.Level {
		return slog.LevelInfo
	}
	handler := slog_.DynamicLevelHandler(getLevel, baseHandler)
	handlerWithGroup := handler.WithGroup("mygroup")

	logger := slog.New(handlerWithGroup)
	logger.Info("test message", "key", "value")

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("mygroup")) {
		t.Errorf("expected output to contain 'mygroup', got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("test message")) {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
}

func TestDynamicLevelHandler_Handle(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	getLevel := func(ctx context.Context) slog.Level {
		return slog.LevelWarn
	}
	handler := slog_.DynamicLevelHandler(getLevel, baseHandler)
	ctx := context.Background()

	// 测试低于阈值的日志不会被处理
	infoRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "info message", 0)
	err := handler.Handle(ctx, infoRecord)
	if err != nil {
		t.Errorf("Handle() returned error: %v", err)
	}
	if buf.Len() > 0 {
		t.Errorf("expected no output for info level, got: %s", buf.String())
	}

	// 测试达到阈值的日志会被处理
	buf.Reset()
	warnRecord := slog.NewRecord(time.Time{}, slog.LevelWarn, "warn message", 0)
	err = handler.Handle(ctx, warnRecord)
	if err != nil {
		t.Errorf("Handle() returned error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("warn message")) {
		t.Errorf("expected output to contain 'warn message', got: %s", buf.String())
	}
}

// TestDynamicLevelHandler_DynamicChange 测试动态改变日志级别
func TestDynamicLevelHandler_DynamicChange(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// 使用变量来动态控制级别
	currentLevel := slog.LevelInfo
	getLevel := func(ctx context.Context) slog.Level {
		return currentLevel
	}
	handler := slog_.DynamicLevelHandler(getLevel, baseHandler)
	logger := slog.New(handler)

	// 初始级别为 Info，Debug 日志不应该输出
	logger.Debug("debug message 1")
	if buf.Len() > 0 {
		t.Errorf("expected no output for debug level when level is Info, got: %s", buf.String())
	}

	// Info 日志应该输出
	buf.Reset()
	logger.Info("info message 1")
	if !bytes.Contains(buf.Bytes(), []byte("info message 1")) {
		t.Errorf("expected output to contain 'info message 1', got: %s", buf.String())
	}

	// 动态改变级别为 Debug
	currentLevel = slog.LevelDebug
	buf.Reset()
	logger.Debug("debug message 2")
	if !bytes.Contains(buf.Bytes(), []byte("debug message 2")) {
		t.Errorf("expected output to contain 'debug message 2' after level change, got: %s", buf.String())
	}

	// 动态改变级别为 Error
	currentLevel = slog.LevelError
	buf.Reset()
	logger.Info("info message 2")
	if buf.Len() > 0 {
		t.Errorf("expected no output for info level when level is Error, got: %s", buf.String())
	}

	// Error 日志应该输出
	buf.Reset()
	logger.Error("error message 1")
	if !bytes.Contains(buf.Bytes(), []byte("error message 1")) {
		t.Errorf("expected output to contain 'error message 1', got: %s", buf.String())
	}
}

// TestDynamicLevelHandler_ContextBasedLevel 测试基于 context 的动态级别
func TestDynamicLevelHandler_ContextBasedLevel(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	type contextKey string
	const levelKey contextKey = "log_level"

	// 从 context 中获取日志级别
	getLevel := func(ctx context.Context) slog.Level {
		if level, ok := ctx.Value(levelKey).(slog.Level); ok {
			return level
		}
		return slog.LevelInfo // 默认级别
	}
	handler := slog_.DynamicLevelHandler(getLevel, baseHandler)
	logger := slog.New(handler)

	// 使用默认级别 (Info)
	ctx1 := context.Background()
	logger.DebugContext(ctx1, "debug message 1")
	if buf.Len() > 0 {
		t.Errorf("expected no output for debug level with default context, got: %s", buf.String())
	}

	buf.Reset()
	logger.InfoContext(ctx1, "info message 1")
	if !bytes.Contains(buf.Bytes(), []byte("info message 1")) {
		t.Errorf("expected output to contain 'info message 1', got: %s", buf.String())
	}

	// 使用 Debug 级别的 context
	ctx2 := context.WithValue(context.Background(), levelKey, slog.LevelDebug)
	buf.Reset()
	logger.DebugContext(ctx2, "debug message 2")
	if !bytes.Contains(buf.Bytes(), []byte("debug message 2")) {
		t.Errorf("expected output to contain 'debug message 2' with debug context, got: %s", buf.String())
	}

	// 使用 Error 级别的 context
	ctx3 := context.WithValue(context.Background(), levelKey, slog.LevelError)
	buf.Reset()
	logger.InfoContext(ctx3, "info message 2")
	if buf.Len() > 0 {
		t.Errorf("expected no output for info level with error context, got: %s", buf.String())
	}

	buf.Reset()
	logger.ErrorContext(ctx3, "error message 1")
	if !bytes.Contains(buf.Bytes(), []byte("error message 1")) {
		t.Errorf("expected output to contain 'error message 1', got: %s", buf.String())
	}
}

// TestDynamicLevelHandler_NilGetLevel 测试 getLevel 为 nil 的情况
func TestDynamicLevelHandler_NilGetLevel(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// getLevel 为 nil，应该完全依赖 baseHandler 的级别控制
	handler := slog_.DynamicLevelHandler(nil, baseHandler)
	logger := slog.New(handler)

	// Debug 日志不应该输出（baseHandler 的级别是 Info）
	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Errorf("expected no output for debug level, got: %s", buf.String())
	}

	// Info 日志应该输出
	buf.Reset()
	logger.Info("info message")
	if !bytes.Contains(buf.Bytes(), []byte("info message")) {
		t.Errorf("expected output to contain 'info message', got: %s", buf.String())
	}
}
