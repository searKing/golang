// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package redis logs dials, commands and pipelines of go-redis v9 clients by log/slog.
//
// Deprecated: Use github.com/searKing/golang/pkg/log/redislog instead, of which this package is an alias.
package redis

import (
	"log/slog"

	"github.com/redis/go-redis/v9"

	"github.com/searKing/golang/pkg/log/redislog"
)

// DefaultSlowThreshold is the threshold of slow calls of DefaultConfig.
//
// Deprecated: Use redislog.DefaultSlowThreshold instead.
const DefaultSlowThreshold = redislog.DefaultSlowThreshold

// DefaultConfig logs calls slower than DefaultSlowThreshold as slow.
// It is a copy of redislog.DefaultConfig, so setting either does not affect the other.
//
// Deprecated: Use redislog.DefaultConfig instead.
var DefaultConfig = redislog.DefaultConfig

// Config configures NewHook.
//
// Deprecated: Use redislog.Config instead.
type Config = redislog.Config

// ErrorToLevel maps errors of calls to levels of logs.
//
// Deprecated: Use redislog.ErrorToLevel instead.
type ErrorToLevel = redislog.ErrorToLevel

// PrintfFunc logs internal logs of go-redis, as redis.SetLogger requires.
//
// Deprecated: Use redislog.PrintfFunc instead.
type PrintfFunc = redislog.PrintfFunc

// DefaultErrorToLevel is the default ErrorToLevel.
//
// Deprecated: Use redislog.DefaultErrorToLevel instead.
func DefaultErrorToLevel(err error) slog.Level { return redislog.DefaultErrorToLevel(err) }

// NewHook returns a hook of go-redis logging by l.
//
// Deprecated: Use redislog.NewHook instead.
func NewHook(l *slog.Logger, config Config) redis.Hook { return redislog.NewHook(l, config) }

// WithSlogLogger returns a hook of go-redis logging by h.
//
// Deprecated: Use redislog.WithSlogLogger instead.
func WithSlogLogger(h slog.Handler, config Config) redis.Hook { return redislog.WithSlogLogger(h, config) }

// InternalLogger returns a PrintfFunc logging internal logs of go-redis by l.
//
// Deprecated: Use redislog.InternalLogger instead.
func InternalLogger(l *slog.Logger) PrintfFunc { return redislog.InternalLogger(l) }
