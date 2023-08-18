// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"io"
	"log/slog"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/searKing/golang/go/log/slog/internal/buffer"
	"golang.org/x/term"
)

var (
	timeNow       = time.Now // Stubbed out for testing.
	baseTimestamp time.Time
	getPid        = os.Getpid // Stubbed out for testing.
)

func init() {
	baseTimestamp = timeNow()
}

// Keys for "built-in" attributes.
const (
	// ErrorKey is the key used by the handlers for the error
	// when the log method is called. The associated Value is an [error].
	ErrorKey = "error"
)

type TimestampMode int

const (
	_ TimestampMode = iota

	// DisableTimestamp disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp

	// SinceStartTimestamp enable the time passed since beginning of execution instead of
	// logging the full timestamp when a TTY is attached.
	SinceStartTimestamp
)

// sharedVar shared const expvar among handler and children handler...
type sharedVar struct {
	once *sync.Once

	// Whether the logger's out is to a terminal
	isTerminal bool
	// The max length of the level text, generated dynamically on init
	maxLevelText int
	// The process id of the caller.
	pid int
}

func (h *sharedVar) init(w io.Writer) {
	h.once.Do(func() {
		if f, ok := w.(*os.File); ok {
			h.isTerminal = term.IsTerminal(int(f.Fd()))
		}
		// Get the max length of the level text
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			levelTextLength := utf8.RuneCount([]byte(level.String()))
			if levelTextLength > h.maxLevelText {
				h.maxLevelText = levelTextLength
			}
		}
		h.pid = getPid()
	})
}

type commonHandler struct {
	// replace level.String()
	ReplaceLevelString func(l slog.Level) string

	// the separator between attributes.
	AttrSep string

	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Force quoting of all values
	ForceQuote bool

	// DisableQuote disables quoting for all values.
	// DisableQuote will have a lower priority than ForceQuote.
	// If both of them are set to true, quote will be forced on all values.
	DisableQuote bool

	// ForceGoroutineId enables goroutine id instead of pid.
	ForceGoroutineId bool

	TimestampMode TimestampMode

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// Disables the glog style ：[IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg msg...
	// replace with ：[IWEF] [yyyymmdd] [hh:mm:ss.uuuuuu] [threadid] [file:line] msg msg...
	HumanReadable bool

	// PadLevelText Adds padding the level text so that all the levels output at the same length
	// PadLevelText is a superset of the DisableLevelTruncation option
	PadLevelText bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// SourcePrettier can be set by the user to modify the content
	// of the file, function and file keys when AddSource is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from slog attrs.
	SourcePrettier func(r slog.Record) *slog.Source

	// WithFuncName append Caller's func name
	WithFuncName bool

	opts slog.HandlerOptions

	sharedVar *sharedVar

	preformattedAttrs []byte
	// groupPrefix is for the text handler only.
	// It holds the prefix for groups that were already pre-formatted.
	// A group will appear here when a call to WithGroup is followed by
	// a call to WithAttrs.
	groupPrefix string
	groups      []string    // all groups started from WithGroup
	nOpenGroups int         // the number of groups opened in preformattedAttrs
	mu          *sync.Mutex // mutex shared among all clones of this handler
	w           io.Writer
}

// NewCommonHandler creates a CommonHandler that writes to w,
// using the given options.
// If opts is nil, the default options are used.
// A [CommonHandler] is a low-level primitive for making structured log.
// [NewGlogHandler] or [NewGlogHumanHandler] recommended.
func NewCommonHandler(w io.Writer, opts *slog.HandlerOptions) *commonHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &commonHandler{
		opts:      *opts,
		sharedVar: &sharedVar{once: &sync.Once{}},
		mu:        &sync.Mutex{},
		w:         w,
	}
}

func (h *commonHandler) clone() *commonHandler {
	// We can't use assignment because we can't copy the mutex.
	h2 := *h
	h2.preformattedAttrs = slices.Clip(h.preformattedAttrs)
	h2.groups = slices.Clip(h.groups)
	return &h2
}

// enabled reports whether l is greater than or equal to the
// minimum level.
func (h *commonHandler) enabled(l slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return l >= minLevel
}

func (h *commonHandler) withAttrs(as []slog.Attr) *commonHandler {
	// We are going to ignore empty groups, so if the entire slice consists of
	// them, there is nothing to do.
	if countEmptyGroups(as) == len(as) {
		return h
	}
	h2 := h.clone()
	// Pre-format the attributes as an optimization.
	buf := buffer.New()
	buf.Write(h2.preformattedAttrs)
	state := h2.newHandleState(buf, false, "")
	defer state.free()
	defer func() {
		h2.preformattedAttrs = buf.Bytes()
	}()
	state.prefix.WriteString(h.groupPrefix)
	if len(h2.preformattedAttrs) > 0 {
		state.sep = h.attrSep()
	}
	state.openGroups()
	for _, a := range as {
		state.appendAttr(a)
	}
	// Remember the new prefix for later keys.
	h2.groupPrefix = state.prefix.String()
	// Remember how many opened groups are in preformattedAttrs,
	// so we don't open them again when we handle a Record.
	h2.nOpenGroups = len(h2.groups)
	return h2
}

func (h *commonHandler) withGroup(name string) *commonHandler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

// handle is the internal implementation of Handler.Handle
// used by GlogHandler and HumanGlogHandler.
// header formats a log header as defined by the C++ implementation.
// It returns a buffer containing the formatted header and the user's file and line number.
// The depth specifies how many stack frames above lives the source line to be identified in the log message.
//
// # LOG LINE PREFIX FORMAT
//
// Log lines have this form:
//
//	Lyyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg...
//
// where the fields are defined as follows:
//
//	L                A single character, representing the log level
//	                 (eg 'I' for INFO)
//	yyyy             The year
//	mm               The month (zero padded; ie May is '05')
//	dd               The day (zero padded)
//	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
//	threadid         The space-padded thread ID as returned by GetTID()
//	                 (this matches the PID on Linux)
//	file             The file name
//	line             The line number
//	msg              The user-supplied message
//
// Example:
//
//	I1103 11:57:31.739339 24395 google.cc:2341] Command line: ./some_prog
//	I1103 11:57:31.739403 24395 google.cc:2342] Process id 24395
//
// NOTE: although the microseconds are useful for comparing events on
// a single machine, clocks on different machines may not be well
// synchronized.  Hence, use caution when comparing the low bits of
// timestamps from different machines.
func (h *commonHandler) handle(r slog.Record) error {
	h.sharedVar.init(h.w)

	state := h.newHandleState(buffer.New(), true, "")
	defer state.free()
	// Built-in attributes. They are not in a group.
	stateGroups := state.groups
	state.groups = nil // So ReplaceAttrs sees no groups instead of the pre groups.
	rep := h.opts.ReplaceAttr
	if h.isColored() {
		state.color = levelColor(r.Level)
	}
	// level
	state.appendLevel(r.Level, h.PadLevelText, h.sharedVar.maxLevelText, h.HumanReadable)
	// time
	t := r.Time // strip monotonic to match Attr behavior
	mode := h.TimestampMode
	if rep != nil {
		a := rep(nil, slog.Time(slog.TimeKey, r.Time))
		if a.Equal(slog.Attr{}) {
			// disable timestamp logging if time is removed.
			t = time.Time{}
			mode = DisableTimestamp
		} else if a.Value.Kind() == slog.KindTime {
			t = a.Value.Time()
		}
	}
	state.appendGlogTime(t, h.TimestampFormat, mode, h.HumanReadable)
	state.appendPid(h.ForceGoroutineId, h.HumanReadable)
	// source
	if h.opts.AddSource {
		if h.SourcePrettier != nil {
			state.appendSource(h.SourcePrettier(r), h.WithFuncName, h.HumanReadable)
		} else {
			state.appendSource(source(r), h.WithFuncName, h.HumanReadable)
		}
	} else {
		if !h.HumanReadable {
			state.buf.WriteString("]")
		}
	}

	var hasMessage bool
	if rep != nil {
		a := rep(nil, slog.String(slog.MessageKey, r.Message))
		if !isEmptyAttr(a) {
			state.buf.WriteString(" ")
			state.appendValue(a.Value)
			hasMessage = true
		}
	} else if r.Message != "" {
		state.buf.WriteString(" ")
		// take message as well formatted raw string, may be with color and so on, disable quote
		state.appendString(r.Message)
		hasMessage = true
	}
	if !hasMessage {
		state.sep = " "
	} else {
		state.sep = h.attrSep()
	}

	state.groups = stateGroups // Restore groups passed to ReplaceAttrs.
	state.appendNonBuiltIns(r)
	state.buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(state.buf.Bytes())
	return err
}

func (s *handleState) appendNonBuiltIns(r slog.Record) {
	// preformatted Attrs
	if len(s.h.preformattedAttrs) > 0 {
		s.buf.WriteString(s.sep)
		s.buf.Write(s.h.preformattedAttrs)
		s.sep = s.h.attrSep()
	}
	// Attrs in Record -- unlike the built-in ones, they are in groups started
	// from WithGroup.
	// If the record has no Attrs, don't output any groups.
	if r.NumAttrs() > 0 {
		s.prefix.WriteString(s.h.groupPrefix)
		s.openGroups()
		r.Attrs(func(a slog.Attr) bool {
			s.appendAttr(a)
			return true
		})
	}
}

// attrSep returns the separator between attributes.
func (h *commonHandler) attrSep() string {
	if h.AttrSep != "" {
		return h.AttrSep
	}
	return " "
}

var groupPool = sync.Pool{New: func() any {
	s := make([]string, 0, 10)
	return &s
}}

func (h *commonHandler) newHandleState(buf *buffer.Buffer, freeBuf bool, sep string) handleState {
	s := handleState{
		h:       h,
		buf:     buf,
		freeBuf: freeBuf,
		sep:     sep,
		prefix:  buffer.New(),
	}
	if h.opts.ReplaceAttr != nil {
		s.groups = groupPool.Get().(*[]string)
		*s.groups = append(*s.groups, h.groups[:h.nOpenGroups]...)
	}
	return s
}

func (h *commonHandler) isColored() bool {
	isColored := h.ForceColors || (h.sharedVar.isTerminal && (runtime.GOOS != "windows"))

	if h.EnvironmentOverrideColors {
		switch force, ok := os.LookupEnv("CLICOLOR_FORCE"); {
		case ok && force != "0":
			isColored = true
		case ok && force == "0", os.Getenv("CLICOLOR") == "0":
			isColored = false
		}
	}

	return isColored && !h.DisableColors
}
