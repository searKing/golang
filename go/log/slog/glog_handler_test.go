// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

var testTime = time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

func TestGlogHandler(t *testing.T) {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	for _, test := range []struct {
		name             string
		attr             slog.Attr
		wantKey, wantVal string
	}{
		{
			"unquoted",
			slog.Int("a", 1),
			"a", "1",
		},
		{
			"quoted",
			slog.String("x = y", `qu"o`),
			`x = y`, "`qu\"o`",
		},
		{
			"String method",
			slog.Any("name", name{"Ren", "Hoek"}),
			`name`, `Hoek, Ren`,
		},
		{
			"struct",
			slog.Any("x", &struct{ A, b int }{A: 1, b: 2}),
			`x`, `&{A:1 b:2}`,
		},
		{
			"TextMarshaler",
			slog.Any("t", text{"abc"}),
			`t`, "`text{\"abc\"}`",
		},
		{
			"TextMarshaler error",
			slog.Any("t", text{""}),
			`t`, `!ERROR:text: empty string`,
		},
		{
			"nil value",
			slog.Any("a", nil),
			`a`, `<nil>`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, opts := range []struct {
				name       string
				opts       slog.HandlerOptions
				wantPrefix string
				modKey     func(string) string
			}{
				{
					"none",
					slog.HandlerOptions{},
					`I20000102 03:04:05.000000 0] a message`,
					func(s string) string { return s },
				},
				{
					"replace",
					slog.HandlerOptions{ReplaceAttr: upperCaseKey},
					`I20000102 03:04:05.000000 0] a message`,
					strings.ToUpper,
				},
			} {
				t.Run(opts.name, func(t *testing.T) {
					var buf bytes.Buffer
					h := NewGlogHandler(&buf, &opts.opts)
					r := slog.NewRecord(testTime, slog.LevelInfo, "a message", 0)
					r.AddAttrs(test.attr)
					if err := h.Handle(context.Background(), r); err != nil {
						t.Fatal(err)
					}
					got := buf.String()
					// Remove final newline.
					got = got[:len(got)-1]
					want := opts.wantPrefix + ", " + opts.modKey(test.wantKey) + "=" + test.wantVal
					if got != want {
						t.Errorf("\ngot  %s\nwant %s", got, want)
					}
				})
			}
		})
	}
}

func TestGlogHumanHandler(t *testing.T) {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	for _, test := range []struct {
		name             string
		attr             slog.Attr
		wantKey, wantVal string
	}{
		{
			"unquoted",
			slog.Int("a", 1),
			"a", "1",
		},
		{
			"quoted",
			slog.String("x = y", `qu"o`),
			`x = y`, "`qu\"o`",
		},
		{
			"String method",
			slog.Any("name", name{"Ren", "Hoek"}),
			`name`, `Hoek, Ren`,
		},
		{
			"struct",
			slog.Any("x", &struct{ A, b int }{A: 1, b: 2}),
			`x`, `&{A:1 b:2}`,
		},
		{
			"TextMarshaler",
			slog.Any("t", text{"abc"}),
			`t`, "`text{\"abc\"}`",
		},
		{
			"TextMarshaler error",
			slog.Any("t", text{""}),
			`t`, `!ERROR:text: empty string`,
		},
		{
			"nil value",
			slog.Any("a", nil),
			`a`, `<nil>`,
		},
		{
			"typed nil value",
			slog.Any("a", (*text)(nil)),
			`a`, `<nil>`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, opts := range []struct {
				name       string
				opts       slog.HandlerOptions
				wantPrefix string
				modKey     func(string) string
			}{
				{
					"none",
					slog.HandlerOptions{},
					`[INFO ][20000102 03:04:05.000000] [0] a message`,
					func(s string) string { return s },
				},
				{
					"replace",
					slog.HandlerOptions{ReplaceAttr: upperCaseKey},
					`[INFO ][20000102 03:04:05.000000] [0] a message`,
					strings.ToUpper,
				},
			} {
				t.Run(opts.name, func(t *testing.T) {
					var buf bytes.Buffer
					h := NewGlogHumanHandler(&buf, &opts.opts)
					r := slog.NewRecord(testTime, slog.LevelInfo, "a message", 0)
					r.AddAttrs(test.attr)
					if err := h.Handle(context.Background(), r); err != nil {
						t.Fatal(err)
					}
					got := buf.String()
					// Remove final newline.
					got = got[:len(got)-1]
					want := opts.wantPrefix + ", " + opts.modKey(test.wantKey) + "=" + test.wantVal
					if got != want {
						t.Errorf("\ngot  %s\nwant %s", got, want)
					}
				})
			}
		})
	}
}

// for testing fmt.Sprint
type name struct {
	First, Last string
}

func (n name) String() string { return n.Last + ", " + n.First }

// for testing TextMarshaler
type text struct {
	s string
}

func (t text) Error() string { return t.s } // should be ignored

func (t text) MarshalText() ([]byte, error) {
	if t.s == "" {
		return nil, errors.New("text: empty string")
	}
	return []byte(fmt.Sprintf("text{%q}", t.s)), nil
}

func TestTextHandlerPreformatted(t *testing.T) {
	var buf bytes.Buffer
	var h slog.Handler = slog.NewTextHandler(&buf, nil)
	h = h.WithAttrs([]slog.Attr{slog.Duration("dur", time.Minute), slog.Bool("b", true)})
	// Also test omitting time.
	r := slog.NewRecord(time.Time{}, 0 /* 0 Level is INFO */, "m", 0)
	r.AddAttrs(slog.Int("a", 1))
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSuffix(buf.String(), "\n")
	want := `level=INFO msg=m dur=1m0s b=true a=1`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestTextHandlerAlloc(t *testing.T) {
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	for i := 0; i < 10; i++ {
		r.AddAttrs(slog.Int("x = y", i))
	}
	var h slog.Handler = slog.NewTextHandler(io.Discard, nil)
	wantAllocs(t, 0, func() { h.Handle(context.Background(), r) })

	h = h.WithGroup("s")
	r.AddAttrs(slog.Group("g", slog.Int("a", 1)))
	wantAllocs(t, 0, func() { h.Handle(context.Background(), r) })
}

func TestNeedsQuoting(t *testing.T) {
	for _, test := range []struct {
		in   string
		want bool
	}{
		{"", true},
		{"ab", false},
		{"a=b", true},
		{`"ab"`, true},
		{"\a\b", true},
		{"a\tb", true},
		{"µåπ", false},
		{"a b", true},
		{"badutf8\xF6", true},
	} {
		got := needsQuoting(test.in, false)
		if got != test.want {
			t.Errorf("%q: got %t, want %t", test.in, got, test.want)
		}
	}
}
