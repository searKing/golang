// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/searKing/golang/go/log/slog/internal/buffer"
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

// WalkAttrFunc is a function type for processing slog.Attr.
// It receives an attribute and returns the processed attribute.
// Can be used to modify, filter, or transform attribute values.
type WalkAttrFunc func(slog.Attr) slog.Attr

// WalkAttr recursively walks through the attribute tree and applies walkFn to each attribute.
//
// Processing steps:
//  1. Resolves the attribute value first (handles lazy LogValuer evaluation)
//  2. If it's a Group type, recursively processes all attributes within the group
//  3. Applies walkFn to each attribute (including the group itself)
//
// Example:
//
//	// Redact sensitive fields
//	sanitized := WalkAttr(attr, func(a slog.Attr) slog.Attr {
//	    if a.Key == "password" {
//	        return slog.String(a.Key, "[REDACTED]")
//	    }
//	    return a
//	})
func WalkAttr(attr slog.Attr, walkFn WalkAttrFunc) (tattr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	switch v := attr.Value; v.Kind() {
	case slog.KindGroup:
		attrs := v.Group()
		as := make([]any, 0, len(attrs))
		for _, a := range attrs {
			as = append(as, WalkAttr(a, walkFn))
		}
		return walkFn(slog.Group(attr.Key, as...))
	default:
		return walkFn(attr)
	}
}

// ReplaceAttrTruncate returns [ReplaceAttr] which shrinks attr's key and value[string]'s len to n at most.
func ReplaceAttrTruncate(n int) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		return truncateAttr(a, n, nil)
	}
}

func ReplaceAttrJsonTruncate(n int) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		return truncateAttr(a, n, truncateAttrAny(true))
	}
}

func ReplaceAttrTextTruncate(n int) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		return truncateAttr(a, n, truncateAttrAny(false))
	}
}

func truncateAttrAny(isJson bool) func(attr slog.Attr, n int) slog.Attr {
	return func(attr slog.Attr, n int) slog.Attr {
		if n <= 0 {
			return attr
		}
		k, v := attr.Key, attr.Value
		a := v.Any()
		_, jm := a.(json.Marshaler)
		if err, ok := a.(error); ok && !jm {
			return slog.String(k, err.Error())
		} else {
			var b buffer.Buffer
			var err error
			if isJson {
				err = b.AppendJSONMarshal(a)
			} else {
				err = b.AppendTextMarshal(a)
			}
			if err != nil {
				return slog.String(k, fmt.Sprintf("!ERROR:%v", err))
			}
			if len(b.Bytes()) > n {
				return slog.String(k, truncate(b.String(), n))
			}
			return attr
		}
	}
}

func truncateAttr(attr slog.Attr, n int, truncateAttrAny func(slog.Attr, int) slog.Attr) (tattr slog.Attr) {
	if n <= 0 {
		return attr
	}
	defer func() {
		if r := recover(); r != nil {
			// If it panics with a nil pointer, the most likely cases are
			// an encoding.TextMarshaler or error fails to guard against nil,
			// in which case "<nil>" seems to be the feasible choice.
			//
			// Adapted from the code in fmt/print.go.
			v := attr.Value
			if v := reflect.ValueOf(v.Any()); v.Kind() == reflect.Pointer && v.IsNil() {
				tattr = slog.String(attr.Key, "<nil>")
				return
			}

			// Otherwise just print the original panic message.
			tattr = slog.String(attr.Key, fmt.Sprintf("!PANIC: %v", r))
		}
	}()
	return WalkAttr(attr, func(a slog.Attr) slog.Attr {
		attr.Key = truncate(attr.Key, n)
		switch v := attr.Value; v.Kind() {
		case slog.KindString:
			return slog.String(attr.Key, truncate(v.String(), n))
		case slog.KindAny:
			if truncateAttrAny != nil {
				return truncateAttrAny(attr, n)
			}
			return attr
		default:
			return attr
		}
	})
}

// ReplaceAttrShortSource returns [ReplaceAttr] which shortens source's function and file.
func ReplaceAttrShortSource() func(groups []string, s slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			if src, ok := a.Value.Any().(*slog.Source); ok {
				if isSourceEmpty(src) {
					return a
				}
				src.Function = shortFunction(src.Function)
				src.File = shortFile(src.File)
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

// isSourceEmpty returns whether the Source struct is nil or only contains zero fields.
func isSourceEmpty(s *slog.Source) bool { return s == nil || *s == slog.Source{} }
