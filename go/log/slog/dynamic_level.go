// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"
	"log/slog"
)

var _ slog.Handler = (*dynamicLevelHandler)(nil)

type dynamicLevelHandler struct {
	// getLevel defines an optional func to returns the level for the given context.
	// If getLevel is nil, the level is always enabled.
	getLevel func(ctx context.Context) slog.Level

	handler slog.Handler
}

func (t dynamicLevelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return (t.getLevel == nil || (t.getLevel != nil && level >= t.getLevel(ctx))) && t.handler.Enabled(ctx, level)
}

func (t dynamicLevelHandler) Handle(ctx context.Context, record slog.Record) error {
	if !t.Enabled(ctx, record.Level) {
		return nil
	}
	return t.handler.Handle(ctx, record)
}

func (t dynamicLevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return DynamicLevelHandler(t.getLevel, t.handler.WithAttrs(attrs))
}

func (t dynamicLevelHandler) WithGroup(name string) slog.Handler {
	return DynamicLevelHandler(t.getLevel, t.handler.WithGroup(name))
}

// DynamicLevelHandler creates a slog.Handler that changes the level of the handler dynamically.
func DynamicLevelHandler(getLevel func(ctx context.Context) slog.Level, handler slog.Handler) slog.Handler {
	return dynamicLevelHandler{
		getLevel: getLevel,
		handler:  handler,
	}
}
