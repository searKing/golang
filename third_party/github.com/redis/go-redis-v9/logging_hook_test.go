// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type entry struct {
	level  slog.Level
	msg    string
	fields map[string]any
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
	r.entries = append(r.entries, entry{level: record.Level, msg: record.Message, fields: m})
	return nil
}

// replyError is an error replied by the redis server.
type replyError string

func (e replyError) Error() string { return string(e) }
func (replyError) RedisError()     {}

func TestHookProcess(t *testing.T) {
	ctx := context.Background()
	r := &recorder{}
	process := NewHook(slog.New(r), Config{}).ProcessHook(func(context.Context, redis.Cmder) error { return nil })
	if err := process(ctx, redis.NewCmd(ctx, "set", "k1", "hello world")); err != nil {
		t.Fatal(err)
	}
	process = NewHook(slog.New(r), Config{}).ProcessHook(func(context.Context, redis.Cmder) error { return replyError("WRONGTYPE wrong") })
	_ = process(ctx, redis.NewCmd(ctx, "incr", "k1"))
	if len(r.entries) != 2 {
		t.Fatalf("got %d entries want 2: %+v", len(r.entries), r.entries)
	}
	if e := r.entries[0]; e.msg != "redis executed" || e.level != slog.LevelDebug ||
		e.fields["redis.cmd"] != `set k1 "hello world"` || !strings.HasSuffix(fmt.Sprint(e.fields["redis.duration"]), "ms") {
		t.Errorf("got %+v", e)
	}
	if _, ok := r.entries[0].fields["redis.error"]; ok {
		t.Errorf("got error of a succeeded call: %+v", r.entries[0])
	}
	if e := r.entries[1]; e.level != slog.LevelWarn || e.fields["redis.cmd"] != "incr k1" || e.fields["redis.error"] != "WRONGTYPE wrong" {
		t.Errorf("got %+v", e)
	}
}

func TestHookLazy(t *testing.T) {
	ctx := context.Background()
	called := false
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}))
	h := NewHook(logger, Config{FieldsFromContext: func(context.Context) []any {
		called = true
		return nil
	}})
	if err := h.ProcessHook(func(context.Context, redis.Cmder) error { return nil })(ctx, redis.NewCmd(ctx, "get", "k1")); err != nil {
		t.Fatal(err)
	}
	if called {
		t.Errorf("fields are built for logs of disabled levels")
	}
}

func TestHookConfig(t *testing.T) {
	ctx := context.Background()
	r := &recorder{}
	h := NewHook(slog.New(r), Config{
		Parameterized:     true,
		Levels:            func(error) slog.Level { return slog.LevelError },
		FieldsFromContext: func(context.Context) []any { return []any{"tenant", "t1"} },
	})
	_ = h.ProcessHook(func(context.Context, redis.Cmder) error { return nil })(ctx, redis.NewCmd(ctx, "set", "k1", "v1"))
	if len(r.entries) != 1 {
		t.Fatalf("got %d entries want 1: %+v", len(r.entries), r.entries)
	}
	if e := r.entries[0]; e.level != slog.LevelError || e.fields["tenant"] != "t1" || e.fields["redis.cmd"] != "set k1 ?" {
		t.Errorf("got %+v", e)
	}
}

func TestHookDialError(t *testing.T) {
	r := &recorder{}
	dial := NewHook(slog.New(r), Config{}).DialHook(func(context.Context, string, string) (net.Conn, error) {
		return nil, &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}
	})
	if _, err := dial(context.Background(), "tcp", "127.0.0.1:0"); err == nil {
		t.Fatal("expected dial error")
	}
	if len(r.entries) != 1 {
		t.Fatalf("got %d entries want 1: %+v", len(r.entries), r.entries)
	}
	if e := r.entries[0]; e.msg != "redis dialed" || e.level != slog.LevelWarn || e.fields["redis.network"] != "tcp" ||
		e.fields["redis.addr"] != "127.0.0.1:0" || e.fields["redis.error"] != "dial tcp: connection refused" {
		t.Errorf("got %+v", e)
	}
}

func TestHookPipeline(t *testing.T) {
	ctx := context.Background()
	var cmds []redis.Cmder
	for i := range 40 {
		cmds = append(cmds, redis.NewCmd(ctx, "get", fmt.Sprintf("k%d", i)))
	}
	r := &recorder{}
	process := NewHook(slog.New(r), Config{}).ProcessPipelineHook(func(context.Context, []redis.Cmder) error { return redis.Nil })
	if err := process(ctx, cmds); !errors.Is(err, redis.Nil) {
		t.Fatal(err)
	}
	if len(r.entries) != 1 {
		t.Fatalf("got %d entries want 1: %+v", len(r.entries), r.entries)
	}
	e := r.entries[0]
	if e.msg != "redis pipeline executed" || e.level != slog.LevelDebug || e.fields["redis.error"] != "redis: nil" {
		t.Errorf("got %+v", e)
	}
	if v, _ := e.fields["redis.cmd"].(commandLines); len(v) != maxLoggedElems || v[0] != "get k0" ||
		v[maxLoggedElems-1] != "... (9 more commands)" {
		t.Errorf("redis.cmd: got %v", e.fields["redis.cmd"])
	}
}

func TestInternalLogger(t *testing.T) {
	r := &recorder{}
	InternalLogger(slog.New(r)).Printf(context.Background(), "redis: connection pool: %s\n", "failed")
	if len(r.entries) != 1 {
		t.Fatalf("got %d entries want 1", len(r.entries))
	}
	if e := r.entries[0]; e.level != slog.LevelWarn || e.msg != "redis: connection pool: failed" {
		t.Errorf("got %+v", e)
	}
}

func TestDefaultErrorToLevel(t *testing.T) {
	tests := []struct {
		err  error
		want slog.Level
	}{
		{nil, slog.LevelDebug},
		{redis.Nil, slog.LevelDebug},
		{context.Canceled, slog.LevelDebug},
		{replyError("NOSCRIPT No matching script"), slog.LevelDebug},
		{replyError("NOAUTH Authentication required."), slog.LevelError},
		{replyError("WRONGPASS invalid username-password pair"), slog.LevelError},
		{replyError("NOPERM this user has no permissions"), slog.LevelError},
		{replyError("OOM command not allowed when used memory > 'maxmemory'."), slog.LevelError},
		{replyError("ERR NOAUTH Authentication required."), slog.LevelError},
		{replyError("ERR invalid password, unauthenticated"), slog.LevelError},
		{redis.ErrNoScript, slog.LevelDebug},
		{context.DeadlineExceeded, slog.LevelWarn},
		{replyError("WRONGTYPE Operation against a key holding the wrong kind of value"), slog.LevelWarn},
		{replyError("BUSYGROUP Consumer Group name already exists"), slog.LevelWarn},
		{errors.New("redis: can't marshal"), slog.LevelWarn},
	}
	for _, tt := range tests {
		if got := DefaultErrorToLevel(tt.err); got != tt.want {
			t.Errorf("DefaultErrorToLevel(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestSlowThreshold(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		config    Config
		cmds      [][]any
		err       error
		wantLevel slog.Level
		wantSlow  bool
	}{
		{Config{}, [][]any{{"get", "k1"}}, nil, slog.LevelDebug, false},
		{DefaultConfig, [][]any{{"get", "k1"}}, nil, slog.LevelDebug, false},
		{Config{SlowThreshold: time.Millisecond}, [][]any{{"get", "k1"}}, nil, slog.LevelWarn, true},
		{Config{SlowThreshold: time.Millisecond}, [][]any{{"get", "k1"}}, replyError("NOAUTH auth"), slog.LevelError, true},
		{Config{SlowThreshold: time.Millisecond}, [][]any{{"blpop", "k1", 1}}, nil, slog.LevelDebug, false},
		{Config{SlowThreshold: time.Millisecond}, [][]any{{"xread", "count", 1, "block", 10, "streams", "k1", "$"}}, nil, slog.LevelDebug, false},
		{Config{SlowThreshold: time.Millisecond}, [][]any{{"xread", "streams", "block", "$"}}, nil, slog.LevelWarn, true},
		{Config{SlowThreshold: time.Millisecond}, [][]any{{"get", "k1"}, {"get", "k2"}}, nil, slog.LevelWarn, true},
		{Config{SlowThreshold: time.Millisecond}, [][]any{{"get", "k1"}, {"brpop", "k2", 1}}, nil, slog.LevelDebug, false},
	}
	for i, tt := range tests {
		r := &recorder{}
		h := NewHook(slog.New(r), tt.config)
		cmds := make([]redis.Cmder, 0, len(tt.cmds))
		for _, args := range tt.cmds {
			cmds = append(cmds, redis.NewCmd(ctx, args...))
		}
		if len(cmds) == 1 {
			_ = h.ProcessHook(func(context.Context, redis.Cmder) error {
				time.Sleep(2 * time.Millisecond)
				return tt.err
			})(ctx, cmds[0])
		} else {
			_ = h.ProcessPipelineHook(func(context.Context, []redis.Cmder) error {
				time.Sleep(2 * time.Millisecond)
				return tt.err
			})(ctx, cmds)
		}
		if len(r.entries) != 1 {
			t.Fatalf("case %d: got %d entries want 1: %+v", i, len(r.entries), r.entries)
		}
		e := r.entries[0]
		if _, slow := e.fields["redis.slow"]; e.level != tt.wantLevel || slow != tt.wantSlow {
			t.Errorf("case %d: got %+v, want level %v and slow %v", i, e, tt.wantLevel, tt.wantSlow)
		}
	}
}

func TestCommandArgsRedacted(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		args []any
		want string
	}{
		{[]any{"get", "k1"}, "get k1"},
		{[]any{"config", "get", "requirepass"}, "config get xxxxx"},
		{[]any{"hello", 3, "auth", "user", "secret"}, "hello 3 xxxxx xxxxx xxxxx"},
		{[]any{"acl", "setuser", "user", ">secret"}, "acl setuser xxxxx xxxxx"},
		{[]any{"migrate", "host", 6379, "k1", 0, 1000, "auth", "secret"}, "migrate xxxxx xxxxx xxxxx xxxxx xxxxx xxxxx xxxxx"},
	}
	for _, tt := range tests {
		if got := commandArgs(redis.NewCmd(ctx, tt.args...)); got != tt.want {
			t.Errorf("commandArgs(%v) = %s, want %s", tt.args, got, tt.want)
		}
	}
}

func TestCommandArgsFormat(t *testing.T) {
	ctx := context.Background()
	s := "v1"
	tests := []struct {
		args []any
		want string
	}{
		{[]any{"set", "k1", "hello world", "ex", 10}, `set k1 "hello world" ex 10`},
		{[]any{"set", "k1", ""}, `set k1 ""`},
		{[]any{"set", "k1", []byte("v1")}, `set k1 v1`},
		{[]any{"set", "k1", `a"b'c\d`}, `set k1 "a\"b'c\\d"`},
		{[]any{"set", "k1", "a\nb\r\t\x00\xff中"}, `set k1 "a\nb\r\t\x00\xff\xe4\xb8\xad"`},
		{[]any{"set", "k1", true, false, 1.5, float32(0.5), 2 * time.Second, &s, (*string)(nil), nil},
			`set k1 1 0 1.5 0.5 2000000000 v1 "" ""`},
		{[]any{"set", "k1", time.Date(2026, 10, 14, 3, 0, 0, 0, time.UTC)}, `set k1 2026-10-14T03:00:00Z`},
	}
	for _, tt := range tests {
		if got := commandArgs(redis.NewCmd(ctx, tt.args...)); got != tt.want {
			t.Errorf("commandArgs(%q) = %s, want %s", tt.args, got, tt.want)
		}
	}
	if got, want := commandStatement(redis.NewCmd(ctx, "set", "my key", "v1")), `set "my key" ?`; got != want {
		t.Errorf("commandStatement = %s, want %s", got, want)
	}

	cmds := []redis.Cmder{redis.NewCmd(ctx, "get", "k1"), redis.NewCmd(ctx, "set", "k1", "hello world")}
	for _, tt := range []struct {
		handler func(io.Writer) slog.Handler
		want    string
	}{
		{func(w io.Writer) slog.Handler { return slog.NewTextHandler(w, nil) },
			`redis.cmd="get k1\nset k1 \"hello world\""`},
		{func(w io.Writer) slog.Handler { return slog.NewJSONHandler(w, nil) },
			`"redis":{"cmd":["get k1","set k1 \"hello world\""]`},
	} {
		var b strings.Builder
		_ = NewHook(slog.New(tt.handler(&b)), Config{Levels: func(error) slog.Level { return slog.LevelInfo }}).
			ProcessPipelineHook(func(context.Context, []redis.Cmder) error { return nil })(ctx, cmds)
		if !strings.Contains(b.String(), tt.want) {
			t.Errorf("got %s, want %s", b.String(), tt.want)
		}
	}
}

func TestCommandKeys(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		args []any
		want []int
	}{
		{[]any{"get", "k1"}, []int{1}},
		{[]any{"get"}, nil},
		{[]any{"set", "k1", "v1", "ex", 10}, []int{1}},
		{[]any{"mget", "k1", "k2", "k3"}, []int{1, 2, 3}},
		{[]any{"mset", "k1", "v1", "k2", "v2"}, []int{1, 3}},
		{[]any{"blpop", "l1", "l2", 0}, []int{1, 2}},
		{[]any{"bitop", "and", "dest", "k1", "k2"}, []int{2, 3, 4}},
		{[]any{"rename", "k1", "k2"}, []int{1, 2}},
		{[]any{"zunionstore", "dest", 2, "z1", "z2", "weights", 1, 2}, []int{1, 3, 4}},
		{[]any{"eval", "return 1", 2, "k1", "k2", "arg"}, []int{3, 4}},
		{[]any{"evalsha", "sha", 0, "arg"}, nil},
		{[]any{"fcall", "fn", 2, "k1", "k2"}, []int{3, 4}},
		{[]any{"xread", "count", 1, "streams", "s1", "s2", "0", "0"}, []int{4, 5}},
		{[]any{"xreadgroup", "group", "g1", "c1", "STREAMS", "s1", ">"}, []int{5}},
		{[]any{"zinter", 2, "z1", "z2"}, []int{2, 3}},
		{[]any{"blmpop", 0, 1, "l1", "left"}, []int{3}},
		{[]any{"memory", "usage", "k1"}, []int{2}},
		{[]any{"memory", "stats"}, nil},
		{[]any{"object", "encoding", "k1"}, []int{2}},
		{[]any{"xinfo", "stream", "s1"}, []int{2}},
		{[]any{"xgroup", "create", "s1", "g1", "$"}, []int{2}},
		{[]any{"migrate", "host", 6379, "k1", 0, 1000}, []int{3}},
		{[]any{"migrate", "host", 6379, "", 0, 1000, "keys", "k1"}, nil},
		{[]any{"config", "get", "maxmemory"}, nil},
		{[]any{"client", "setname", "app"}, nil},
		{[]any{"publish", "channel", "message"}, nil},
		{[]any{"ping"}, nil},
		{[]any{"auth", "user", "secret"}, nil},
		{[]any{"hello", 3, "auth", "user", "secret"}, nil},
	}
	for _, tt := range tests {
		if got := commandKeys(redis.NewCmd(ctx, tt.args...)); !slices.Equal(got, tt.want) {
			t.Errorf("commandKeys(%v) = %v, want %v", tt.args, got, tt.want)
		}
	}
}

func TestCommandStatement(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		args []any
		want string
	}{
		{[]any{"get", "k1"}, "get k1"},
		{[]any{"set", "k1", "v1", "ex", 10}, "set k1 ? ? ?"},
		{[]any{"mset", "k1", "v1", "k2", "v2"}, "mset k1 ? k2 ?"},
		{[]any{"hset", "h", "f1", "v1"}, "hset h ? ?"},
		{[]any{"eval", "return 1", 1, "k1", "arg"}, "eval ? ? k1 ?"},
		{[]any{"config", "set", "requirepass", "secret"}, "config set ? ?"},
		{[]any{"auth", "user", "secret"}, "auth ? ?"},
		{[]any{"hello", 3, "auth", "user", "secret"}, "hello ? ? ? ?"},
		{[]any{"object", "encoding", "k1"}, "object encoding k1"},
		{[]any{"ping"}, "ping"},
	}
	for _, tt := range tests {
		if got := commandStatement(redis.NewCmd(ctx, tt.args...)); got != tt.want {
			t.Errorf("commandStatement(%v) = %q, want %q", tt.args, got, tt.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	ctx := context.Background()
	long := strings.Repeat("v", maxLoggedBytes+10)
	if got, want := commandArgs(redis.NewCmd(ctx, "set", "k1", long)),
		`set k1 "`+long[:maxLoggedBytes]+`... (10 more bytes)"`; got != want {
		t.Errorf("long argument: got %s want %s", got, want)
	}

	args := []any{"mget"}
	for i := range 40 {
		args = append(args, fmt.Sprintf("k%d", i))
	}
	if got := strings.Fields(commandArgs(redis.NewCmd(ctx, args...))); len(got) != maxLoggedElems+3 ||
		strings.Join(got[maxLoggedElems-2:], " ") != "k29 ... (10 more arguments)" {
		t.Errorf("many arguments: got %v", got)
	}
	if statement := commandStatement(redis.NewCmd(ctx, args...)); !strings.HasSuffix(statement, "k29 ... (10 more arguments)") {
		t.Errorf("many arguments of statement: got %s", statement)
	}
}
