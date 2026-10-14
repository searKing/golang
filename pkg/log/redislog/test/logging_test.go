// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"

	"github.com/searKing/golang/pkg/log/redislog"
)

type entry struct {
	level  slog.Level
	msg    string
	fields map[string]any
	source runtime.Frame
}

// recorder is a slog.Handler recording logs, with fields of groups keyed by "group.key".
type recorder struct {
	mu      sync.Mutex
	entries []entry
}

func (r *recorder) Enabled(context.Context, slog.Level) bool { return true }
func (r *recorder) WithAttrs([]slog.Attr) slog.Handler       { return r }
func (r *recorder) WithGroup(string) slog.Handler            { return r }

func (r *recorder) Handle(_ context.Context, record slog.Record) error {
	m := map[string]any{}
	record.Attrs(func(a slog.Attr) bool {
		if a.Value.Kind() == slog.KindGroup {
			for _, g := range a.Value.Group() {
				m[a.Key+"."+g.Key] = g.Value.Any()
			}
			return true
		}
		m[a.Key] = a.Value.Any()
		return true
	})
	r.mu.Lock()
	defer r.mu.Unlock()
	source, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
	r.entries = append(r.entries, entry{level: record.Level, msg: record.Message, fields: m, source: source})
	return nil
}

// executed returns entries of commands and pipelines, skipping dials and the handshake by HELLO.
func (r *recorder) executed() []entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	var es []entry
	for _, e := range r.entries {
		if e.msg == "redis dialed" || strings.HasPrefix(fmt.Sprint(e.fields["redis.cmd"]), "hello ") {
			continue
		}
		es = append(es, e)
	}
	return es
}

func newClient(t *testing.T, config redislog.Config) (*redis.Client, *recorder) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{
		Addr:                     mr.Addr(),
		DisableIdentity:          true,
		MaintNotificationsConfig: &maintnotifications.Config{Mode: maintnotifications.ModeDisabled},
	})
	t.Cleanup(func() { _ = rdb.Close() })
	r := &recorder{}
	rdb.AddHook(redislog.NewHook(slog.New(r), config))
	return rdb, r
}

func TestHook(t *testing.T) {
	rdb, r := newClient(t, redislog.DefaultConfig)
	ctx := context.Background()
	if err := rdb.Set(ctx, "k1", "hello world", 0).Err(); err != nil {
		t.Fatal(err)
	}
	if err := rdb.Get(ctx, "missing").Err(); !errors.Is(err, redis.Nil) {
		t.Fatalf("get missing: %v", err)
	}
	if _, err := rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.Get(ctx, "k1")
		p.Get(ctx, "missing")
		return nil
	}); !errors.Is(err, redis.Nil) {
		t.Fatalf("pipeline: %v", err)
	}
	if err := rdb.HSet(ctx, "h", "f", "v").Err(); err != nil {
		t.Fatal(err)
	}
	_ = rdb.Incr(ctx, "h").Err()
	_ = rdb.EvalSha(ctx, "0000000000000000000000000000000000000000", nil).Err()
	// AUTH fails on miniredis without passwords, which is enough to log its arguments.
	_ = rdb.Do(ctx, "auth", "user", "secret").Err()

	type want struct {
		msg   string
		level slog.Level
		cmd   string
		error string // "" if left out.
	}
	wants := []want{
		{"redis executed", slog.LevelDebug, `set k1 "hello world"`, ""},
		{"redis executed", slog.LevelDebug, "get missing", "redis: nil"},
		{"redis pipeline executed", slog.LevelDebug, "[get k1 get missing]", "redis: nil"},
		{"redis executed", slog.LevelDebug, "hset h f v", ""},
		{"redis executed", slog.LevelWarn, "incr h", "WRONGTYPE Operation against a key holding the wrong kind of value"},
		{"redis executed", slog.LevelDebug, "evalsha 0000000000000000000000000000000000000000 0", "NOSCRIPT No matching script. Please use EVAL."},
		{"redis executed", slog.LevelError, "auth user xxxxx", "WRONGPASS invalid username-password pair"},
	}
	got := r.executed()
	if len(got) != len(wants) {
		t.Fatalf("got %d calls want %d: %+v", len(got), len(wants), got)
	}
	for i, w := range wants {
		e := got[i]
		errField, hasErr := e.fields["redis.error"]
		if e.msg != w.msg || e.level != w.level || fmt.Sprint(e.fields["redis.cmd"]) != w.cmd ||
			hasErr != (w.error != "") || (hasErr && errField != w.error) {
			t.Errorf("call %d: got %+v want %+v", i, e, w)
		}
	}
	for _, e := range r.entries {
		if strings.Contains(fmt.Sprint(e.fields), "secret") {
			t.Errorf("credential leaked into log: %+v", e)
		}
	}
}

func TestHookParameterized(t *testing.T) {
	rdb, r := newClient(t, redislog.Config{Parameterized: true})
	if err := rdb.Set(context.Background(), "k1", "hello world", 0).Err(); err != nil {
		t.Fatal(err)
	}
	got := r.executed()
	if len(got) != 1 || got[0].fields["redis.cmd"] != "set k1 ?" {
		t.Errorf("got %+v", got)
	}
}

// TestHookSource checks sources of logs, which are callers of redis rather than the hook.
func TestHookSource(t *testing.T) {
	rdb, r := newClient(t, redislog.Config{})
	ctx := context.Background()
	_, _, line, _ := runtime.Caller(0)
	_ = rdb.Set(ctx, "k1", "v1", 0).Err() // The first command dials and runs HELLO nested.
	_, _ = rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.Get(ctx, "k1")
		return nil
	})
	// Dials have no source, as go-redis dials in goroutines of connection pools.
	// Pipelines dial again, for go-redis runs them on a pool of their own.
	type call struct {
		name string
		line int
	}
	want := []call{
		{"dial", 0}, {"hello", line + 1}, {"set", line + 1},
		{"dial", 0}, {"hello", line + 2}, {"pipeline", line + 2},
	}
	if len(r.entries) != len(want) {
		t.Fatalf("got %d entries want %d: %+v", len(r.entries), len(want), r.entries)
	}
	for i, e := range r.entries {
		name, _, _ := strings.Cut(fmt.Sprint(e.fields["redis.cmd"]), " ")
		switch e.msg {
		case "redis dialed":
			name = "dial"
		case "redis pipeline executed":
			name = "pipeline"
		}
		if name != want[i].name || e.source.Line != want[i].line ||
			(want[i].line != 0 && !strings.HasSuffix(e.source.File, "/test/logging_test.go")) {
			t.Errorf("entry %d: got %s at %s:%d (%s), want %s at line %d of this file",
				i, name, e.source.File, e.source.Line, e.source.Function, want[i].name, want[i].line)
		}
	}
}

// TestInternalLoggerSource checks sources of internal logs, which are where go-redis logs, the caller of Printf.
func TestInternalLoggerSource(t *testing.T) {
	r := &recorder{}
	_, _, line, _ := runtime.Caller(0)
	redislog.InternalLogger(slog.New(r)).Printf(context.Background(), "redis: connection pool: %s", "failed")
	if len(r.entries) != 1 {
		t.Fatalf("got %d entries want 1", len(r.entries))
	}
	if source := r.entries[0].source; !strings.HasSuffix(source.File, "/test/logging_test.go") || source.Line != line+1 {
		t.Errorf("got source %s:%d (%s), want line %d of this file", source.File, source.Line, source.Function, line+1)
	}
}
