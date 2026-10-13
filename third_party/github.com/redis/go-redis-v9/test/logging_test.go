// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"

	redis_ "github.com/searKing/golang/third_party/github.com/redis/go-redis-v9"
)

type entry struct {
	level  slog.Level
	msg    string
	fields map[string]any
}

// recorder is a slog.Handler recording logs.
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
		m[a.Key] = a.Value.Any()
		return true
	})
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, entry{level: record.Level, msg: record.Message, fields: m})
	return nil
}

// logged returns entries of msg, skipping dials and the handshake by HELLO.
func (r *recorder) logged(msg string) []entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	var es []entry
	for _, e := range r.entries {
		if e.msg != msg || e.fields[redis_.CallTypeFieldKey] == string(redis_.Dial) || e.fields[redis_.CommandFieldKey] == "hello" {
			continue
		}
		es = append(es, e)
	}
	return es
}

func (r *recorder) finished() []entry { return r.logged("finished call") }

func newRedisClient(t *testing.T) *redis.Client {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{
		Addr:                     mr.Addr(),
		DisableIdentity:          true,
		MaintNotificationsConfig: &maintnotifications.Config{Mode: maintnotifications.ModeDisabled},
	})
	t.Cleanup(func() { _ = rdb.Close() })
	return rdb
}

func newClient(t *testing.T, opts ...redis_.LoggingOption) (*redis.Client, *recorder) {
	rdb := newRedisClient(t)
	r := &recorder{}
	rdb.AddHook(redis_.Hook(slog.New(r), opts...))
	return rdb, r
}

func TestHook(t *testing.T) {
	rdb, r := newClient(t)
	ctx := context.Background()
	const value = "feature-blob"
	if err := rdb.Set(ctx, "k1", value, 0).Err(); err != nil {
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
	if err := rdb.HSet(ctx, "k2", "f", "v").Err(); err != nil {
		t.Fatal(err)
	}
	if err := rdb.Incr(ctx, "k2").Err(); err == nil {
		t.Fatal("expected WRONGTYPE error")
	}

	for _, e := range r.entries {
		if strings.Contains(fmt.Sprint(e.fields), value) {
			t.Fatalf("value leaked into log: %+v", e)
		}
		// HELLO runs nested in the first command, with its context.
		if _, ok := e.fields[redis_.KeyFieldKey]; ok && e.fields[redis_.CommandFieldKey] == "hello" {
			t.Errorf("hello is logged with the key of the outer command: %+v", e)
		}
	}

	type want struct {
		level    slog.Level
		callType redis_.CallType
		command  string
		key      string
		keys     string
		code     string
		error    bool
	}
	wants := []want{
		{level: slog.LevelDebug, callType: redis_.Command, command: "set", key: "k1", code: "OK"},
		{level: slog.LevelDebug, callType: redis_.Command, command: "get", key: "missing", code: "NIL", error: true},
		{level: slog.LevelDebug, callType: redis_.Pipeline, keys: "[k1 missing]", code: "NIL", error: true},
		{level: slog.LevelDebug, callType: redis_.Command, command: "hset", key: "k2", code: "OK"},
		{level: slog.LevelWarn, callType: redis_.Command, command: "incr", key: "k2", code: "WRONGTYPE", error: true},
	}
	got := r.finished()
	if len(got) != len(wants) {
		t.Fatalf("got %d finished calls want %d: %+v", len(got), len(wants), got)
	}
	for i, w := range wants {
		e := got[i]
		_, hasErr := e.fields["redis.error"]
		if e.level != w.level || e.fields["redis.code"] != w.code || hasErr != w.error ||
			e.fields[redis_.CallTypeFieldKey] != string(w.callType) ||
			e.fields["protocol"] != "redis" || e.fields[redis_.ComponentFieldKey] != redis_.KindClientFieldValue {
			t.Errorf("call %d: got %+v want %+v", i, e, w)
		}
		if w.keys != "" {
			if fmt.Sprint(e.fields[redis_.KeysFieldKey]) != w.keys {
				t.Errorf("call %d: got keys %v want %s", i, e.fields[redis_.KeysFieldKey], w.keys)
			}
			continue
		}
		if e.fields[redis_.CommandFieldKey] != w.command || e.fields[redis_.KeyFieldKey] != w.key {
			t.Errorf("call %d: got %+v want %+v", i, e, w)
		}
	}
	for _, e := range r.entries {
		if e.msg != "finished call" {
			t.Errorf("got %q, want finished calls only by default", e.msg)
		}
		if _, ok := e.fields["redis.request.content"]; ok {
			t.Errorf("payloads are logged by default: %+v", e)
		}
	}
}

func TestHookPayload(t *testing.T) {
	rdb, r := newClient(t,
		redis_.WithLogOnEvents(redis_.StartCall, redis_.FinishCall),
		redis_.WithPayloads(redis_.RequestPayload, redis_.ResponsePayload))
	ctx := context.Background()
	if err := rdb.Set(ctx, "k1", "v1", 0).Err(); err != nil {
		t.Fatal(err)
	}
	if err := rdb.Get(ctx, "k1").Err(); err != nil {
		t.Fatal(err)
	}
	if _, err := rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.Get(ctx, "k1")
		p.Incr(ctx, "n")
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	// Replies of a pipeline are logged even if some commands failed.
	if _, err := rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.Get(ctx, "k1")
		p.Get(ctx, "missing")
		return nil
	}); !errors.Is(err, redis.Nil) {
		t.Fatalf("pipeline: %v", err)
	}
	// AUTH fails on miniredis without passwords, which is enough to log its arguments.
	_ = rdb.Do(ctx, "auth", "user", "secret").Err()

	type want struct {
		request  string
		response string // "" if left out.
	}
	wants := []want{
		{"[set k1 v1]", "OK"},
		{"[get k1]", "v1"},
		{"[[get k1] [incr n]]", "[v1 1]"},
		{"[[get k1] [get missing]]", "[v1 redis: nil]"},
		// AUTH failed, of which the reply is left out.
		{"[auth xxxxx xxxxx]", ""},
	}
	started, finished := r.logged("started call"), r.finished()
	if len(started) != len(wants) || len(finished) != len(wants) {
		t.Fatalf("got %d started and %d finished calls, want %d: %+v", len(started), len(finished), len(wants), r.entries)
	}
	for i, w := range wants {
		if got := fmt.Sprint(started[i].fields["redis.request.content"]); got != w.request {
			t.Errorf("started call %d: got request %s want %s", i, got, w.request)
		}
		if _, ok := started[i].fields["redis.response.content"]; ok {
			t.Errorf("started call %d: got response %+v", i, started[i])
		}
		e := finished[i]
		if got := fmt.Sprint(e.fields["redis.request.content"]); got != w.request {
			t.Errorf("finished call %d: got request %s want %s", i, got, w.request)
		}
		response, ok := e.fields["redis.response.content"]
		if (w.response == "") == ok || (ok && fmt.Sprint(response) != w.response) {
			t.Errorf("finished call %d: got response %v want %q", i, response, w.response)
		}
	}
	for _, e := range r.entries {
		if strings.Contains(fmt.Sprint(e.fields), "secret") {
			t.Errorf("credential leaked into log: %+v", e)
		}
	}
}

func TestHookCallMeta(t *testing.T) {
	var mu sync.Mutex
	var metas []redis_.CallMeta
	rdb, _ := newClient(t, redis_.WithFieldsFromContextAndCallMeta(func(_ context.Context, c redis_.CallMeta) redis_.Fields {
		mu.Lock()
		defer mu.Unlock()
		metas = append(metas, c)
		return nil
	}))
	if err := rdb.Set(context.Background(), "k1", "v1", 0).Err(); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	var dialed, commanded bool
	for _, c := range metas {
		switch c.Typ {
		case redis_.Dial:
			dialed = c.Network == "tcp" && c.Addr != "" && c.Cmds == nil
		case redis_.Command:
			commanded = commanded || (len(c.Cmds) == 1 && c.Cmds[0].Name() == "set")
		}
	}
	if !dialed || !commanded {
		t.Errorf("got call metas %+v, want a dial and a set command", metas)
	}
}

func TestHookOptions(t *testing.T) {
	rdb, r := newClient(t,
		redis_.WithLevels(func(error) slog.Level { return slog.LevelError }),
		redis_.WithDurationField(redis_.DurationToDurationField),
		redis_.WithDisableLoggingFields("redis.key"),
		redis_.WithFieldsFromContext(func(context.Context) redis_.Fields { return redis_.Fields{"tenant", "t1"} }),
	)
	if err := rdb.Set(context.Background(), "k1", "v", 0).Err(); err != nil {
		t.Fatal(err)
	}
	got := r.finished()
	if len(got) != 1 {
		t.Fatalf("got %d finished calls want 1: %+v", len(got), got)
	}
	e := got[0]
	if _, ok := e.fields[redis_.KeyFieldKey]; ok {
		t.Errorf("redis.key is not disabled: %+v", e)
	}
	if _, ok := e.fields["redis.duration"]; !ok {
		t.Errorf("redis.duration is not logged: %+v", e)
	}
	if e.level != slog.LevelError || e.fields["tenant"] != "t1" {
		t.Errorf("got %+v, want level error and tenant t1", e)
	}
}

// TestHookRedisError checks codes and levels of errors replied by a redis server.
func TestHookRedisError(t *testing.T) {
	rdb, r := newClient(t)
	ctx := context.Background()
	if err := rdb.HSet(ctx, "h", "f", "v").Err(); err != nil {
		t.Fatal(err)
	}
	_ = rdb.Incr(ctx, "h").Err()
	_ = rdb.EvalSha(ctx, "0000000000000000000000000000000000000000", nil).Err()
	_ = rdb.Do(ctx, "nosuchcommand").Err()

	type want struct {
		code  string
		level slog.Level
	}
	wants := []want{{"OK", slog.LevelDebug}, {"WRONGTYPE", slog.LevelWarn}, {"NOSCRIPT", slog.LevelDebug}, {"ERR", slog.LevelWarn}}
	got := r.finished()
	if len(got) != len(wants) {
		t.Fatalf("got %d finished calls want %d: %+v", len(got), len(wants), got)
	}
	for i, w := range wants {
		if e := got[i]; e.fields["redis.code"] != w.code || e.level != w.level {
			t.Errorf("call %d: got %+v want %+v", i, e, w)
		}
	}
}

func TestInstrument(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{
		Addr:                     mr.Addr(),
		DB:                       2,
		DisableIdentity:          true,
		MaintNotificationsConfig: &maintnotifications.Config{Mode: maintnotifications.ModeDisabled},
	})
	t.Cleanup(func() { _ = rdb.Close() })
	r := &recorder{}
	redis_.Instrument(rdb, slog.New(r), redis_.WithPayloads(redis_.StatementPayload))
	if err := rdb.Set(context.Background(), "k1", "v1", 0).Err(); err != nil {
		t.Fatal(err)
	}
	// The handshake selects the database by a pipeline of SELECT, before SET.
	i := slices.IndexFunc(r.finished(), func(e entry) bool { return e.fields[redis_.CommandFieldKey] == "set" })
	if i < 0 {
		t.Fatalf("set is not logged: %+v", r.entries)
	}
	if e := r.finished()[i]; e.fields[redis_.DBFieldKey] != int64(2) || e.fields[redis_.AddrFieldKey] != mr.Addr() ||
		e.fields["redis.statement"] != "set k1 ?" {
		t.Errorf("got %+v", e)
	}
}
