// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	strings_ "github.com/searKing/golang/go/strings"
)

// Error returns an Attr for an error value.
func Error(err error) slog.Attr {
	return slog.Any(ErrorKey, err)
}

// isEmptyAttr reports whether a has an empty key and a nil value.
// That can be written as Attr{} or Any("", nil).
func isEmptyAttr(a slog.Attr) bool {
	if a.Key == "" {
		return true
	}
	switch a.Value.Kind() {
	case slog.KindAny:
		return a.Equal(slog.Any("", nil))
	case slog.KindBool:
		return a.Equal(slog.Bool("", false))
	case slog.KindDuration:
		return a.Equal(slog.Duration("", 0))
	case slog.KindFloat64:
		return a.Equal(slog.Float64("", 0))
	case slog.KindInt64:
		return a.Equal(slog.Int64("", 0))
	case slog.KindString:
		return a.Equal(slog.String("", ""))
	case slog.KindTime:
		return a.Equal(slog.Time("", time.Time{}))
	case slog.KindUint64:
		return a.Equal(slog.Uint64("", 0))
	case slog.KindGroup:
		return a.Equal(slog.Group(""))
	default:
		s := a.String()
		return s == "" || s == "<nil>"
	}
}

// ReplaceAttrTruncate returns [ReplaceAttr] which shrinks attr's key and value[string]'s len to n at most.
func ReplaceAttrTruncate(n int) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if n > 0 {
			k := truncate(a.Key, n)
			switch a.Value.Kind() {
			case slog.KindString:
				return slog.String(k, truncate(a.Value.String(), n))
			default:
				return a
			}
		}
		return a
	}
}

// ReplaceAttrKeys allows customization of the key names for default fields.
type ReplaceAttrKeys map[string]string

func (f ReplaceAttrKeys) resolve(key string) string {
	if k, ok := f[key]; ok {
		return k
	}
	return key
}

// This is to not silently overwrite `time`, `msg`, `func` and `level` fields when
// dumping it. If this code wasn't there doing:
//
//	slog.With("level", 1).Info("hello")
//
// Would just silently drop the user provided level. Instead with this code
// it'll logged as:
//
//	{"level": "info", "fields.level": 1, "msg": "hello", "time": "..."}
//
// It's not exported because it's still using Data in an opinionated way. It's to
// avoid code duplication between the two default formatters.
func prefixAttrClashes(attrs []slog.Attr, builtinAttrKeys ReplaceAttrKeys) []slog.Attr {
	if len(builtinAttrKeys) == 0 {
		return attrs
	}
	attrs = replaceAttrClash(attrs, builtinAttrKeys.resolve(slog.TimeKey))
	attrs = replaceAttrClash(attrs, builtinAttrKeys.resolve(slog.MessageKey))
	attrs = replaceAttrClash(attrs, builtinAttrKeys.resolve(slog.LevelKey))
	attrs = replaceAttrClash(attrs, builtinAttrKeys.resolve(slog.SourceKey))
	return attrs
}

func replaceAttrClash(attrs []slog.Attr, key string) []slog.Attr {
	var val slog.Value
	if !slices.ContainsFunc(attrs, func(attr slog.Attr) bool {
		return attr.Key == key
	}) {
		return attrs
	}
	attrs = append(attrs, slog.Attr{
		Key:   "fields." + key,
		Value: val,
	})
	return slices.DeleteFunc(attrs, func(attr slog.Attr) bool {
		return attr.Key == key
	})
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("size: %d, string: ", len(s)))
	buf.WriteString(strings_.Truncate(s, n))
	return buf.String()
}
