// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/searKing/golang/go/log/slog/internal/buffer"
	"github.com/searKing/golang/go/runtime/goroutine"
	strings_ "github.com/searKing/golang/go/strings"
	time_ "github.com/searKing/golang/go/time"
)

// handleState holds state for a single call to commonHandler.handle.
// The initial value of sep determines whether to emit a separator
// before the next key, after which it stays true.
type handleState struct {
	h       *commonHandler
	buf     *buffer.Buffer
	freeBuf bool           // should buf be freed?
	sep     string         // separator to write before next key
	prefix  *buffer.Buffer // for text: key prefix
	groups  *[]string      // pool-allocated slice of active groups, for ReplaceAttr
}

func (s *handleState) free() {
	if s.freeBuf {
		s.buf.Free()
	}
	if gs := s.groups; gs != nil {
		*gs = (*gs)[:0]
		groupPool.Put(gs)
	}
	s.prefix.Free()
}

func (s *handleState) openGroups() {
	for _, n := range s.h.groups[s.h.nOpenGroups:] {
		s.openGroup(n)
	}
}

// Separator for group names and keys.
const keyComponentSep = '.'

// openGroup starts a new group of attributes
// with the given name.
func (s *handleState) openGroup(name string) {
	s.prefix.WriteString(name)
	s.prefix.WriteByte(keyComponentSep)
	// Collect group names for ReplaceAttr.
	if s.groups != nil {
		*s.groups = append(*s.groups, name)
	}
}

// closeGroup ends the group with the given name.
func (s *handleState) closeGroup(name string) {
	s.prefix.Truncate(s.prefix.Len() - len(name) - 1) /* -1 for keyComponentSep */
	s.sep = s.h.attrSep()
	if s.groups != nil {
		*s.groups = (*s.groups)[:len(*s.groups)-1]
	}
}

// replaceAttr handles replacement and checking for an empty key after replacement.
func (s *handleState) replaceAttr(a slog.Attr) slog.Attr {
	if rep := s.h.opts.ReplaceAttr; rep != nil && a.Value.Kind() != slog.KindGroup {
		var gs []string
		if s.groups != nil {
			gs = *s.groups
		}
		// Resolve before calling ReplaceAttr, so the user doesn't have to.
		a.Value = a.Value.Resolve()
		a = rep(gs, a)
	}
	a.Value = a.Value.Resolve()
	// Elide empty Attrs.
	if isEmptyAttr(a) {
		return a
	}
	// Special case: Source.
	if v := a.Value; v.Kind() == slog.KindAny {
		if src, ok := v.Any().(*slog.Source); ok {
			a.Value = sourceAsGroup(src)
		}
	}
	return a
}

// replaceKey handles replacement and checking for an empty key after replacement.
func (s *handleState) replaceKey(key string) string {
	if s.prefix != nil && s.prefix.Len() > 0 {
		return s.prefix.String() + key
	}
	return key
}

// appendAttr appends the Attr's key and value using app.
// It handles replacement and checking for an empty key.
// after replacement.
func (s *handleState) appendAttr(a slog.Attr) {
	a = s.replaceAttr(a)
	// Elide empty Attrs.
	if isEmptyAttr(a) {
		return
	}
	if a.Value.Kind() == slog.KindGroup {
		attrs := a.Value.Group()
		// Output only non-empty groups.
		if len(attrs) > 0 {
			// Inline a group with an empty key.
			if a.Key != "" {
				s.openGroup(a.Key)
			}
			for _, aa := range attrs {
				s.appendAttr(aa)
			}
			if a.Key != "" {
				s.closeGroup(a.Key)
			}
		}
	} else {
		s.appendKey(a.Key)
		s.appendValue(a.Value)
	}
}

func (s *handleState) appendError(err error) {
	s.appendString(fmt.Sprintf("!ERROR:%v", err))
}

func (s *handleState) appendKey(key string) {
	s.buf.WriteString(s.sep)
	s.appendString(s.replaceKey(key))
	s.buf.WriteByte('=')
	s.sep = s.h.attrSep()
}

func (s *handleState) appendString(str string) {
	// text
	if s.h.ForceQuote ||
		(!s.h.DisableQuote && needsQuoting(str, false)) ||
		(s.h.DisableQuote && needsQuoting(str, true)) {
		s.buf.WriteString(strconv.Quote(str))
	} else {
		s.buf.WriteString(str)
	}
}

func (s *handleState) appendValue(v slog.Value) {
	var err error
	err = appendTextValue(s, v)
	if err != nil {
		s.appendError(err)
	}
}

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

func (s *handleState) appendLevel(level slog.Level, preferColor bool, padLevelText bool, maxLevelText int, humanReadable bool) (colored bool) {
	var c int
	if preferColor {
		switch level {
		case slog.LevelDebug:
			c = gray
		case slog.LevelWarn:
			c = yellow
		case slog.LevelError:
			c = red
		default:
			c = blue
		}
	}
	colored = c > 0

	// level
	key := slog.LevelKey
	var val string
	a := s.replaceAttr(slog.Any(key, level))
	if !a.Equal(slog.Attr{}) {
		var limit int
		// Handle custom level values.
		level, ok := a.Value.Any().(slog.Level)
		if ok {
			if f := s.h.ReplaceLevelString; f != nil {
				val = s.h.ReplaceLevelString(level)
			} else {
				val = level.String()
			}
			limit = maxLevelText
			if limit > 0 && limit < len(val) {
				val = val[0:limit]
			}
		} else {
			val = a.Value.String()
		}
		if padLevelText && limit > 0 {
			// Generates the format string used in the next line, for example "%-6s" or "%-7s".
			// Based on the max level text length.
			var pad strings.Builder
			pad.WriteString("%-")
			pad.WriteString(strconv.Itoa(limit))
			pad.WriteString("s")

			// Formats the level text by appending spaces up to the max length, for example:
			// 	- "INFO   "
			//	- "WARNING"
			val = fmt.Sprintf(pad.String(), val)
		}
	}

	// Avoid Fprintf, for speed. The format is so simple that we can do it quickly by hand.
	// It's worth about 3X. Fprintf is hard.
	if colored {
		s.buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, val))
	} else if humanReadable {
		s.buf.WriteString("[")
		s.buf.WriteString(val)
		s.buf.WriteString("]")
	} else {
		s.buf.WriteString(val)
	}
	return colored
}

func (s *handleState) appendGlogTime(t time.Time, layout string, mode TimestampMode, humanReadable bool) {
	if mode == DisableTimestamp {
		return
	}
	val := t.Round(0) // strip monotonic to match Attr behavior
	switch mode {
	case SinceStartTimestamp:
		if humanReadable {
			s.buf.WriteString(fmt.Sprintf(" [%04d]", int(val.Sub(baseTimestamp)/time.Second)))
			return
		}
		s.buf.WriteString(fmt.Sprintf("%04d", int(val.Sub(baseTimestamp)/time.Second)))
	default:
		if humanReadable {
			s.buf.WriteString(fmt.Sprintf("[%s]", val.Format(strings_.ValueOrDefault(layout, time_.GLogDate))))
			return
		}
		s.buf.WriteString(val.Format(strings_.ValueOrDefault(layout, time_.GLogDate)))
	}
}

func (s *handleState) appendTime(t time.Time) {
	writeTimeRFC3339Millis(s.buf, t)
}

func (s *handleState) appendPid(forceGoroutineId bool, humanReadable bool) {
	if s.buf.Len() > 0 {
		s.buf.WriteString(" ")
	}
	if forceGoroutineId {
		if humanReadable {
			s.buf.WriteString(fmt.Sprintf("[%-3d]", goroutine.ID()))
		} else {
			s.buf.WriteString(fmt.Sprintf("%-3d", goroutine.ID()))
		}
	} else {
		if humanReadable {
			// " [{pid}]"
			s.buf.WriteString("[")
			s.buf.WriteString(strconv.Itoa(s.h.pid))
			s.buf.WriteString("]")
		} else {
			// " {pid}"
			s.buf.WriteString(strconv.Itoa(s.h.pid))
		}
	}
}

func (s *handleState) appendSource(src *slog.Source, withFuncName bool, humanReadable bool) {
	if withFuncName && src.Function != "" {
		if humanReadable {
			s.buf.WriteString(fmt.Sprintf(" [%s:%d](%s)", src.File, src.Line, src.Function))
		} else {
			s.buf.WriteString(fmt.Sprintf(" %s:%d(%s)]", src.File, src.Line, src.Function))
		}
	} else {
		if humanReadable {
			s.buf.WriteString(fmt.Sprintf(" [%s:%d]", src.File, src.Line))
		} else {
			s.buf.WriteString(fmt.Sprintf(" %s:%d]", src.File, src.Line))
		}
	}
}

// This takes half the time of Time.AppendFormat.
func writeTimeRFC3339Millis(buf *buffer.Buffer, t time.Time) {
	year, month, day := t.Date()
	buf.WritePosIntWidth(year, 4)
	buf.WriteByte('-')
	buf.WritePosIntWidth(int(month), 2)
	buf.WriteByte('-')
	buf.WritePosIntWidth(day, 2)
	buf.WriteByte('T')
	hour, min, sec := t.Clock()
	buf.WritePosIntWidth(hour, 2)
	buf.WriteByte(':')
	buf.WritePosIntWidth(min, 2)
	buf.WriteByte(':')
	buf.WritePosIntWidth(sec, 2)
	ns := t.Nanosecond()
	buf.WriteByte('.')
	buf.WritePosIntWidth(ns/1e6, 3)
	_, offsetSeconds := t.Zone()
	if offsetSeconds == 0 {
		buf.WriteByte('Z')
	} else {
		offsetMinutes := offsetSeconds / 60
		if offsetMinutes < 0 {
			buf.WriteByte('-')
			offsetMinutes = -offsetMinutes
		} else {
			buf.WriteByte('+')
		}
		buf.WritePosIntWidth(offsetMinutes/60, 2)
		buf.WriteByte(':')
		buf.WritePosIntWidth(offsetMinutes%60, 2)
	}
}
