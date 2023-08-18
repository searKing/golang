// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConcurrentWrites(t *testing.T) {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	ctx := context.Background()
	count := 1000
	for _, handlerType := range []string{"text", "json", "glog", "glog_human"} {
		t.Run(handlerType, func(t *testing.T) {
			var buf bytes.Buffer
			var h slog.Handler
			switch handlerType {
			case "text":
				h = slog.NewTextHandler(&buf, nil)
			case "json":
				h = slog.NewJSONHandler(&buf, nil)
			case "glog":
				h = NewGlogHandler(&buf, nil)
			case "glog_human":
				h = NewGlogHumanHandler(&buf, nil)
			default:
				t.Fatalf("unexpected handlerType %q", handlerType)
			}
			sub1 := h.WithAttrs([]slog.Attr{slog.Bool("sub1", true)})
			sub2 := h.WithAttrs([]slog.Attr{slog.Bool("sub2", true)})
			var wg sync.WaitGroup
			for i := 0; i < count; i++ {
				sub1Record := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)
				sub1Record.AddAttrs(slog.Int("i", i))
				sub2Record := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub2", 0)
				sub2Record.AddAttrs(slog.Int("i", i))
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := sub1.Handle(ctx, sub1Record); err != nil {
						t.Error(err)
					}
					if err := sub2.Handle(ctx, sub2Record); err != nil {
						t.Error(err)
					}
				}()
			}
			wg.Wait()
			for i := 1; i <= 2; i++ {
				want := "hello from sub" + strconv.Itoa(i)
				n := strings.Count(buf.String(), want)
				if n != count {
					t.Fatalf("want %d occurrences of %q, got %d", count, want, n)
				}
			}
		})
	}
}

type replace struct {
	v slog.Value
}

func (r *replace) LogValue() slog.Value { return r.v }

// Verify the common parts of TextHandler and JSONHandler.
func TestHandlers(t *testing.T) {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	// remove all Attrs
	removeAll := func(_ []string, a slog.Attr) slog.Attr { return slog.Attr{} }

	attrs := []slog.Attr{slog.String("a", "one"), slog.Int("b", 2), slog.Any("", nil)}
	preAttrs := []slog.Attr{slog.Int("pre", 3), slog.String("x", "y")}

	for _, test := range []struct {
		name          string
		replace       func([]string, slog.Attr) slog.Attr
		addSource     bool
		with          func(slog.Handler) slog.Handler
		preAttrs      []slog.Attr
		attrs         []slog.Attr
		wantText      string
		wantJSON      string
		wantGlog      string
		wantGlogHuman string
	}{
		{
			name:          "basic",
			attrs:         attrs,
			wantText:      "time=2000-01-02T03:04:05.000Z level=INFO msg=message a=one b=2",
			wantJSON:      `{"time":"2000-01-02T03:04:05Z","level":"INFO","msg":"message","a":"one","b":2}`,
			wantGlog:      `I20000102 03:04:05.000000 0] message, a=one, b=2`,
			wantGlogHuman: `[INFO ][20000102 03:04:05.000000] [0] message, a=one, b=2`,
		},
		{
			name:          "empty key",
			attrs:         append(slices.Clip(attrs), slog.Any("", "v")),
			wantText:      `time=2000-01-02T03:04:05.000Z level=INFO msg=message a=one b=2 ""=v`,
			wantJSON:      `{"time":"2000-01-02T03:04:05Z","level":"INFO","msg":"message","a":"one","b":2,"":"v"}`,
			wantGlog:      `I20000102 03:04:05.000000 0] message, a=one, b=2`,
			wantGlogHuman: `[INFO ][20000102 03:04:05.000000] [0] message, a=one, b=2`,
		},
		{
			name:          "cap keys",
			replace:       upperCaseKey,
			attrs:         attrs,
			wantText:      "TIME=2000-01-02T03:04:05.000Z LEVEL=INFO MSG=message A=one B=2",
			wantJSON:      `{"TIME":"2000-01-02T03:04:05Z","LEVEL":"INFO","MSG":"message","A":"one","B":2}`,
			wantGlog:      `I20000102 03:04:05.000000 0] message, A=one, B=2`,
			wantGlogHuman: `[INFO ][20000102 03:04:05.000000] [0] message, A=one, B=2`,
		},
		{
			name:          "remove all",
			replace:       removeAll,
			attrs:         attrs,
			wantText:      "",
			wantJSON:      `{}`,
			wantGlog:      `0]`,
			wantGlogHuman: `[] [0]`,
		},
		{
			name:          "preformatted",
			with:          func(h slog.Handler) slog.Handler { return h.WithAttrs(preAttrs) },
			preAttrs:      preAttrs,
			attrs:         attrs,
			wantText:      "time=2000-01-02T03:04:05.000Z level=INFO msg=message pre=3 x=y a=one b=2",
			wantJSON:      `{"time":"2000-01-02T03:04:05Z","level":"INFO","msg":"message","pre":3,"x":"y","a":"one","b":2}`,
			wantGlog:      `I20000102 03:04:05.000000 0] message, pre=3, x=y, a=one, b=2`,
			wantGlogHuman: `[INFO ][20000102 03:04:05.000000] [0] message, pre=3, x=y, a=one, b=2`,
		},
		{
			name:          "preformatted cap keys",
			replace:       upperCaseKey,
			with:          func(h slog.Handler) slog.Handler { return h.WithAttrs(preAttrs) },
			preAttrs:      preAttrs,
			attrs:         attrs,
			wantText:      "TIME=2000-01-02T03:04:05.000Z LEVEL=INFO MSG=message PRE=3 X=y A=one B=2",
			wantJSON:      `{"TIME":"2000-01-02T03:04:05Z","LEVEL":"INFO","MSG":"message","PRE":3,"X":"y","A":"one","B":2}`,
			wantGlog:      `I20000102 03:04:05.000000 0] message, PRE=3, X=y, A=one, B=2`,
			wantGlogHuman: `[INFO ][20000102 03:04:05.000000] [0] message, PRE=3, X=y, A=one, B=2`,
		},
		{
			name:          "preformatted remove all",
			replace:       removeAll,
			with:          func(h slog.Handler) slog.Handler { return h.WithAttrs(preAttrs) },
			preAttrs:      preAttrs,
			attrs:         attrs,
			wantText:      "",
			wantJSON:      "{}",
			wantGlog:      `0]`,
			wantGlogHuman: `[] [0]`,
		},
		{
			name:          "remove built-in",
			replace:       removeKeys(slog.TimeKey, slog.LevelKey, slog.MessageKey),
			attrs:         attrs,
			wantText:      "a=one b=2",
			wantJSON:      `{"a":"one","b":2}`,
			wantGlog:      `0] a=one, b=2`,
			wantGlogHuman: `[] [0] a=one, b=2`,
		},
		{
			name:          "preformatted remove built-in",
			replace:       removeKeys(slog.TimeKey, slog.LevelKey, slog.MessageKey),
			with:          func(h slog.Handler) slog.Handler { return h.WithAttrs(preAttrs) },
			attrs:         attrs,
			wantText:      "pre=3 x=y a=one b=2",
			wantJSON:      `{"pre":3,"x":"y","a":"one","b":2}`,
			wantGlog:      `0] pre=3, x=y, a=one, b=2`,
			wantGlogHuman: `[] [0] pre=3, x=y, a=one, b=2`,
		},
		{
			name:    "groups",
			replace: removeKeys(slog.TimeKey, slog.LevelKey), // to simplify the result
			attrs: []slog.Attr{
				slog.Int("a", 1),
				slog.Group("g",
					slog.Int("b", 2),
					slog.Group("h", slog.Int("c", 3)),
					slog.Int("d", 4)),
				slog.Int("e", 5),
			},
			wantText:      "msg=message a=1 g.b=2 g.h.c=3 g.d=4 e=5",
			wantJSON:      `{"msg":"message","a":1,"g":{"b":2,"h":{"c":3},"d":4},"e":5}`,
			wantGlog:      `0] message, a=1, g.b=2, g.h.c=3, g.d=4, e=5`,
			wantGlogHuman: `[] [0] message, a=1, g.b=2, g.h.c=3, g.d=4, e=5`,
		},
		{
			name:          "empty group",
			replace:       removeKeys(slog.TimeKey, slog.LevelKey),
			attrs:         []slog.Attr{slog.Group("g"), slog.Group("h", slog.Int("a", 1))},
			wantText:      "msg=message h.a=1",
			wantJSON:      `{"msg":"message","h":{"a":1}}`,
			wantGlog:      `0] message, h.a=1`,
			wantGlogHuman: `[] [0] message, h.a=1`,
		},
		{
			name:    "nested empty group",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			attrs: []slog.Attr{
				slog.Group("g",
					slog.Group("h",
						slog.Group("i"), slog.Group("j"))),
			},
			wantText:      `msg=message`,
			wantJSON:      `{"msg":"message"}`,
			wantGlog:      `0] message`,
			wantGlogHuman: `[] [0] message`,
		},
		{
			name:    "nested non-empty group",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			attrs: []slog.Attr{
				slog.Group("g",
					slog.Group("h",
						slog.Group("i"), slog.Group("j", slog.Int("a", 1)))),
			},
			wantText:      `msg=message g.h.j.a=1`,
			wantJSON:      `{"msg":"message","g":{"h":{"j":{"a":1}}}}`,
			wantGlog:      `0] message, g.h.j.a=1`,
			wantGlogHuman: `[] [0] message, g.h.j.a=1`,
		},
		{
			name:    "escapes",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			attrs: []slog.Attr{
				slog.String("a b", "x\t\n\000y"),
				slog.Group(" b.c=\"\\x2E\t",
					slog.String("d=e", "f.g\""),
					slog.Int("m.d", 1)), // dot is not escaped
			},
			wantText:      `msg=message "a b"="x\t\n\x00y" " b.c=\"\\x2E\t.d=e"="f.g\"" " b.c=\"\\x2E\t.m.d"=1`,
			wantJSON:      `{"msg":"message","a b":"x\t\n\u0000y"," b.c=\"\\x2E\t":{"d=e":"f.g\"","m.d":1}}`,
			wantGlog:      "0] message, a b=\"x\\t\\n\\x00y\", ` b.c=\"\\x2E\t.d=e`=`f.g\"`, ` b.c=\"\\x2E\t.m.d`=1",
			wantGlogHuman: "[] [0] message, a b=\"x\\t\\n\\x00y\", ` b.c=\"\\x2E\t.d=e`=`f.g\"`, ` b.c=\"\\x2E\t.m.d`=1",
		},
		{
			name:    "LogValuer",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			attrs: []slog.Attr{
				slog.Int("a", 1),
				slog.Any("name", logValueName{"Ren", "Hoek"}),
				slog.Int("b", 2),
			},
			wantText:      "msg=message a=1 name.first=Ren name.last=Hoek b=2",
			wantJSON:      `{"msg":"message","a":1,"name":{"first":"Ren","last":"Hoek"},"b":2}`,
			wantGlog:      `0] message, a=1, name.first=Ren, name.last=Hoek, b=2`,
			wantGlogHuman: `[] [0] message, a=1, name.first=Ren, name.last=Hoek, b=2`,
		},
		{
			// Test resolution when there is no ReplaceAttr function.
			name: "resolve",
			attrs: []slog.Attr{
				slog.Any("", &replace{slog.Value{}}), // should be elided
				slog.Any("name", logValueName{"Ren", "Hoek"}),
			},
			wantText:      "time=2000-01-02T03:04:05.000Z level=INFO msg=message name.first=Ren name.last=Hoek",
			wantJSON:      `{"time":"2000-01-02T03:04:05Z","level":"INFO","msg":"message","name":{"first":"Ren","last":"Hoek"}}`,
			wantGlog:      `I20000102 03:04:05.000000 0] message, name.first=Ren, name.last=Hoek`,
			wantGlogHuman: `[INFO ][20000102 03:04:05.000000] [0] message, name.first=Ren, name.last=Hoek`,
		},
		{
			name:          "with-group",
			replace:       removeKeys(slog.TimeKey, slog.LevelKey),
			with:          func(h slog.Handler) slog.Handler { return h.WithAttrs(preAttrs).WithGroup("s") },
			attrs:         attrs,
			wantText:      "msg=message pre=3 x=y s.a=one s.b=2",
			wantJSON:      `{"msg":"message","pre":3,"x":"y","s":{"a":"one","b":2}}`,
			wantGlog:      `0] message, pre=3, x=y, s.a=one, s.b=2`,
			wantGlogHuman: `[] [0] message, pre=3, x=y, s.a=one, s.b=2`,
		},
		{
			name:    "preformatted with-groups",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			with: func(h slog.Handler) slog.Handler {
				return h.WithAttrs([]slog.Attr{slog.Int("p1", 1)}).
					WithGroup("s1").
					WithAttrs([]slog.Attr{slog.Int("p2", 2)}).
					WithGroup("s2").
					WithAttrs([]slog.Attr{slog.Int("p3", 3)})
			},
			attrs:         attrs,
			wantText:      "msg=message p1=1 s1.p2=2 s1.s2.p3=3 s1.s2.a=one s1.s2.b=2",
			wantJSON:      `{"msg":"message","p1":1,"s1":{"p2":2,"s2":{"p3":3,"a":"one","b":2}}}`,
			wantGlog:      `0] message, p1=1, s1.p2=2, s1.s2.p3=3, s1.s2.a=one, s1.s2.b=2`,
			wantGlogHuman: `[] [0] message, p1=1, s1.p2=2, s1.s2.p3=3, s1.s2.a=one, s1.s2.b=2`,
		},
		{
			name:    "two with-groups",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			with: func(h slog.Handler) slog.Handler {
				return h.WithAttrs([]slog.Attr{slog.Int("p1", 1)}).
					WithGroup("s1").
					WithGroup("s2")
			},
			attrs:         attrs,
			wantText:      "msg=message p1=1 s1.s2.a=one s1.s2.b=2",
			wantJSON:      `{"msg":"message","p1":1,"s1":{"s2":{"a":"one","b":2}}}`,
			wantGlog:      `0] message, p1=1, s1.s2.a=one, s1.s2.b=2`,
			wantGlogHuman: `[] [0] message, p1=1, s1.s2.a=one, s1.s2.b=2`,
		},
		{
			name:    "empty with-groups",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			with: func(h slog.Handler) slog.Handler {
				return h.WithGroup("x").WithGroup("y")
			},
			wantText:      "msg=message",
			wantJSON:      `{"msg":"message"}`,
			wantGlog:      `0] message`,
			wantGlogHuman: `[] [0] message`,
		},
		{
			name:    "empty with-groups, no non-empty attrs",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			with: func(h slog.Handler) slog.Handler {
				return h.WithGroup("x").WithAttrs([]slog.Attr{slog.Group("g")}).WithGroup("y")
			},
			wantText:      "msg=message",
			wantJSON:      `{"msg":"message"}`,
			wantGlog:      `0] message`,
			wantGlogHuman: `[] [0] message`,
		},
		{
			name:    "one empty with-group",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			with: func(h slog.Handler) slog.Handler {
				return h.WithGroup("x").WithAttrs([]slog.Attr{slog.Int("a", 1)}).WithGroup("y")
			},
			attrs:         []slog.Attr{slog.Group("g", slog.Group("h"))},
			wantText:      "msg=message x.a=1",
			wantJSON:      `{"msg":"message","x":{"a":1}}`,
			wantGlog:      `0] message, x.a=1`,
			wantGlogHuman: `[] [0] message, x.a=1`,
		},
		{
			name:          "GroupValue as Attr value",
			replace:       removeKeys(slog.TimeKey, slog.LevelKey),
			attrs:         []slog.Attr{{"v", slog.AnyValue(slog.IntValue(3))}},
			wantText:      "msg=message v=3",
			wantJSON:      `{"msg":"message","v":3}`,
			wantGlog:      `0] message, v=3`,
			wantGlogHuman: `[] [0] message, v=3`,
		},
		{
			name:          "byte slice",
			replace:       removeKeys(slog.TimeKey, slog.LevelKey),
			attrs:         []slog.Attr{slog.Any("bs", []byte{1, 2, 3, 4})},
			wantText:      `msg=message bs="\x01\x02\x03\x04"`,
			wantJSON:      `{"msg":"message","bs":"AQIDBA=="}`,
			wantGlog:      `0] message, bs="\x01\x02\x03\x04"`,
			wantGlogHuman: `[] [0] message, bs="\x01\x02\x03\x04"`,
		},
		{
			name:          "json.RawMessage",
			replace:       removeKeys(slog.TimeKey, slog.LevelKey),
			attrs:         []slog.Attr{slog.Any("bs", json.RawMessage("1234"))},
			wantText:      `msg=message bs="1234"`,
			wantJSON:      `{"msg":"message","bs":1234}`,
			wantGlog:      `0] message, bs=1234`,
			wantGlogHuman: `[] [0] message, bs=1234`,
		},
		{
			name:    "inline group",
			replace: removeKeys(slog.TimeKey, slog.LevelKey),
			attrs: []slog.Attr{
				slog.Int("a", 1),
				slog.Group("", slog.Int("b", 2), slog.Int("c", 3)),
				slog.Int("d", 4),
			},
			wantText:      `msg=message a=1 b=2 c=3 d=4`,
			wantJSON:      `{"msg":"message","a":1,"b":2,"c":3,"d":4}`,
			wantGlog:      `0] message, a=1, d=4`,
			wantGlogHuman: `[] [0] message, a=1, d=4`,
		},
		{
			name: "Source",
			replace: func(gs []string, a slog.Attr) slog.Attr {
				if a.Key == slog.SourceKey {
					s := a.Value.Any().(*slog.Source)
					s.Function = filepath.Join(filepath.Base(filepath.Dir(filepath.Base(s.Function))), filepath.Base(s.Function))
					s.File = filepath.Base(s.File)
					return slog.Any(a.Key, s)
				}
				return removeKeys(slog.TimeKey, slog.LevelKey)(gs, a)
			},
			addSource:     true,
			wantText:      `source=handler_test.go:$LINE msg=message`,
			wantJSON:      `{"source":{"function":"slog.TestHandlers","file":"handler_test.go","line":$LINE},"msg":"message"}`,
			wantGlog:      `0 handler_test.go:$LINE] message`,
			wantGlogHuman: `[] [0] [handler_test.go:$LINE](TestHandlers) message`,
		},
		{
			name: "replace built-in with group",
			replace: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Group(slog.TimeKey, "mins", 3, "secs", 2)
				}
				if a.Key == slog.LevelKey {
					return slog.Attr{}
				}
				return a
			},
			wantText:      `time.mins=3 time.secs=2 msg=message`,
			wantJSON:      `{"time":{"mins":3,"secs":2},"msg":"message"}`,
			wantGlog:      `20000102 03:04:05.000000 0] message`,
			wantGlogHuman: `[][20000102 03:04:05.000000] [0] message`,
		},
	} {
		r := slog.NewRecord(testTime, slog.LevelInfo, "message", callerPC(2))
		line := strconv.Itoa(source(r).Line)
		r.AddAttrs(test.attrs...)
		var buf bytes.Buffer
		opts := slog.HandlerOptions{ReplaceAttr: test.replace, AddSource: test.addSource}
		t.Run(test.name, func(t *testing.T) {
			for _, handler := range []struct {
				name string
				h    slog.Handler
				want string
			}{
				{"text", slog.NewTextHandler(&buf, &opts), test.wantText},
				{"json", slog.NewJSONHandler(&buf, &opts), test.wantJSON},
				{"glog", NewGlogHandler(&buf, &opts), test.wantGlog},
				{"glog_human", NewGlogHumanHandler(&buf, &opts), test.wantGlogHuman},
			} {
				t.Run(handler.name, func(t *testing.T) {
					h := handler.h
					if test.with != nil {
						h = test.with(h)
					}
					buf.Reset()
					if err := h.Handle(nil, r); err != nil {
						t.Fatal(err)
					}
					want := strings.ReplaceAll(handler.want, "$LINE", line)
					got := strings.TrimSuffix(buf.String(), "\n")
					if got != want {
						t.Errorf("\ngot  %s\nwant %s\n", got, want)
					}
				})
			}
		})
	}
}

// removeKeys returns a function suitable for HandlerOptions.ReplaceAttr
// that removes all Attrs with the given keys.
func removeKeys(keys ...string) func([]string, slog.Attr) slog.Attr {
	return func(_ []string, a slog.Attr) slog.Attr {
		for _, k := range keys {
			if a.Key == k {
				return slog.Attr{}
			}
		}
		return a
	}
}

func upperCaseKey(_ []string, a slog.Attr) slog.Attr {
	a.Key = strings.ToUpper(a.Key)
	return a
}

type logValueName struct {
	first, last string
}

func (n logValueName) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("first", n.first),
		slog.String("last", n.last))
}

func TestHandlerEnabled(t *testing.T) {
	levelVar := func(l slog.Level) *slog.LevelVar {
		var al slog.LevelVar
		al.Set(l)
		return &al
	}

	for _, test := range []struct {
		leveler slog.Leveler
		want    bool
	}{
		{nil, true},
		{slog.LevelWarn, false},
		{&slog.LevelVar{}, true}, // defaults to Info
		{levelVar(slog.LevelWarn), false},
		{slog.LevelDebug, true},
		{levelVar(slog.LevelDebug), true},
	} {
		h := &commonHandler{opts: slog.HandlerOptions{Level: test.leveler}}
		got := h.enabled(slog.LevelInfo)
		if got != test.want {
			t.Errorf("%v: got %t, want %t", test.leveler, got, test.want)
		}
	}
}

func TestSecondWith(t *testing.T) {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()
	// Verify that a second call to Logger.With does not corrupt
	// the original.
	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, &slog.HandlerOptions{ReplaceAttr: removeKeys(slog.TimeKey)})
	logger := slog.New(h).With(
		slog.String("app", "playground"),
		slog.String("role", "tester"),
		slog.Int("data_version", 2),
	)
	appLogger := logger.With("type", "log") // this becomes type=met
	_ = logger.With("type", "metric")
	appLogger.Info("foo")
	got := strings.TrimSpace(buf.String())
	want := `level=INFO msg=foo app=playground role=tester data_version=2 type=log`
	if got != want {
		t.Errorf("\ngot  %s\nwant %s", got, want)
	}
}

func TestReplaceAttrGroups(t *testing.T) {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()
	// Verify that ReplaceAttr is called with the correct groups.
	type ga struct {
		groups string
		key    string
		val    string
	}

	var got []ga

	h := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{ReplaceAttr: func(gs []string, a slog.Attr) slog.Attr {
		v := a.Value.String()
		if a.Key == slog.TimeKey {
			v = "<now>"
		}
		got = append(got, ga{strings.Join(gs, ","), a.Key, v})
		return a
	}})
	slog.New(h).
		With(slog.Int("a", 1)).
		WithGroup("g1").
		With(slog.Int("b", 2)).
		WithGroup("g2").
		With(
			slog.Int("c", 3),
			slog.Group("g3", slog.Int("d", 4)),
			slog.Int("e", 5)).
		Info("m",
			slog.Int("f", 6),
			slog.Group("g4", slog.Int("h", 7)),
			slog.Int("i", 8))

	want := []ga{
		{"", "a", "1"},
		{"g1", "b", "2"},
		{"g1,g2", "c", "3"},
		{"g1,g2,g3", "d", "4"},
		{"g1,g2", "e", "5"},
		{"", "time", "<now>"},
		{"", "level", "INFO"},
		{"", "msg", "m"},
		{"g1,g2", "f", "6"},
		{"g1,g2,g4", "h", "7"},
		{"g1,g2", "i", "8"},
	}
	if !slices.Equal(got, want) {
		t.Errorf("\ngot  %v\nwant %v", got, want)
	}
}
