// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/searKing/golang/go/runtime/goroutine"
	strings_ "github.com/searKing/golang/go/strings"
	time_ "github.com/searKing/golang/go/time"
	"github.com/sirupsen/logrus"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

var (
	timeNow       = time.Now // Stubbed out for testing.
	baseTimestamp time.Time
	getPid        = os.Getpid // Stubbed out for testing.
)

func init() {
	baseTimestamp = timeNow()
}

// GlogFormatter formats logs into text
// https://medium.com/technical-tips/google-log-glog-output-format-7eb31b3f0ce5
// [IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg
// IWEF — Log Levels, I for INFO, W for WARNING, E for ERROR and `F` for FATAL.
// yyyymmdd — Year, Month and Date.
// hh:mm:ss.uuuuuu — Hours, Minutes, Seconds and Microseconds.
// threadid — PID/TID of the process/thread.
// file:line — File name and line number.
// msg — Actual user-specified log message.
type GlogFormatter struct {
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

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable the time passed since beginning of execution instead of
	// logging the full timestamp when a TTY is attached.
	SinceStartTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// replace level.String()
	LevelStringFunc func(level logrus.Level) string

	// replace message.String()
	MessageStringFunc func(value interface{}) string

	// replace key
	KeyStringFunc func(key string) string

	// replace value.String()
	ValueStringFunc func(value interface{}) string

	// Set the truncation of the level text to n characters.
	// >0, truncate the level text to n characters at most.
	// =0, truncate the level text to 1 characters at most.
	// <0, don't truncate
	LevelTruncationLimit int

	// Disables the glog style ：[IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg msg...
	// replace with ：[IWEF] [yyyymmdd] [hh:mm:ss.uuuuuu] [threadid] [file:line] msg msg...
	HumanReadable bool

	// PadLevelText Adds padding the level text so that all the levels output at the same length
	// PadLevelText is a superset of the DisableLevelTruncation option
	PadLevelText bool

	// WithFuncName append Caller's func name
	WithFuncName bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &GlogFormatter{
	//     FieldMap: FieldMap{
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"}}
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	terminalInitOnce sync.Once

	// The max length of the level text, generated dynamically on init
	levelTextMaxLength int

	pid int
}

func NewGlogFormatter() *GlogFormatter {
	return &GlogFormatter{}
}

func NewGlogEnhancedFormatter() *GlogFormatter {
	return &GlogFormatter{
		DisableQuote:         true,
		LevelTruncationLimit: 5,
		PadLevelText:         true,
		HumanReadable:        true,
		WithFuncName:         true,
		QuoteEmptyFields:     true,
		DisableSorting:       false,
		LevelStringFunc: func(level logrus.Level) string {
			if level == logrus.WarnLevel {
				return "WARN"
			}
			return strings.ToUpper(level.String())
		},
	}
}

func (f *GlogFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
	// Get the max length of the level text
	for _, level := range logrus.AllLevels {
		levelTextLength := utf8.RuneCount([]byte(f.levelString(level)))
		if levelTextLength > f.levelTextMaxLength {
			f.levelTextMaxLength = levelTextLength
		}
	}
	f.pid = getPid()
}

func (f *GlogFormatter) levelString(level logrus.Level) string {
	if f.LevelStringFunc != nil {
		return f.LevelStringFunc(level)
	}
	return strings.ToUpper(level.String())
}

func (f *GlogFormatter) messageString(value interface{}) string {
	if f.MessageStringFunc != nil {
		return f.MessageStringFunc(value)
	}
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	return stringVal
}

func (f *GlogFormatter) keyString(key string) string {
	if f.KeyStringFunc != nil {
		return f.KeyStringFunc(key)
	}
	return key
}

func (f *GlogFormatter) valueString(value interface{}) string {
	if f.ValueStringFunc != nil {
		return f.ValueStringFunc(value)
	}
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	return stringVal
}

func (f *GlogFormatter) isColored() bool {
	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))

	if f.EnvironmentOverrideColors {
		switch force, ok := os.LookupEnv("CLICOLOR_FORCE"); {
		case ok && force != "0":
			isColored = true
		case ok && force == "0", os.Getenv("CLICOLOR") == "0":
			isColored = false
		}
	}

	return isColored && !f.DisableColors
}

// Format renders a single log entry
func (f *GlogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data {
		if k == logrus.ErrorKey {
			continue
		}
		keys = append(keys, k)
	}

	fixedKeys := make([]string, 0, 4+len(data))
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyMsg))
	}
	if _, has := data[logrus.ErrorKey]; has {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(fieldKey(logrus.ErrorKey)))
	}

	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
		} else {
			f.SortingFunc(keys)
		}
		fixedKeys = append(fixedKeys, keys...)
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.terminalInitOnce.Do(func() { f.init(entry) })

	levelColor, levelText := f.level(entry.Level)

	b.Write(f.header(entry, 0, levelColor, levelText))

	for _, key := range fixedKeys {
		var value interface{}
		switch {
		case key == f.FieldMap.resolve(FieldKeyMsg):
			// Remove a single newline if it already exists in the message to keep
			// the behavior of logrus glog_formatter the same as the stdlib log package
			if levelColor > 0 {
				value = strings.TrimSuffix(entry.Message, "\n")
			} else {
				value = entry.Message
			}
			if levelColor > 0 {
				fmt.Fprintf(b, "\x1b[0m")
				f.appendMessage(b, value)
			} else {
				f.appendMessage(b, value)
			}
			continue
		case key == f.FieldMap.resolve(fieldKey(logrus.ErrorKey)):
			value = data[logrus.ErrorKey]
		default:
			value = data[key]
		}
		if levelColor > 0 {
			_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m=", levelColor, f.keyString(key))
			f.appendValue(b, value)
		} else {
			f.appendKeyValue(b, key, value)
		}
	}

	if levelColor <= 0 {
		b.WriteByte('\n')
	}
	return b.Bytes(), nil
}

func (f *GlogFormatter) needsQuoting(text string, message bool) bool {
	if f.ForceQuote {
		return true
	}
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.DisableQuote {
		return false
	}

	if message {
		return false
	}

	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *GlogFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteString(", ")
	}
	b.WriteString(f.keyString(key))
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *GlogFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal := f.valueString(value)

	if !f.needsQuoting(stringVal, false) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func (f *GlogFormatter) appendMessage(b *bytes.Buffer, value interface{}) {
	stringVal := f.messageString(value)

	if !f.needsQuoting(stringVal, true) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

/*
header formats a log header as defined by the C++ implementation.
It returns a buffer containing the formatted header and the user's file and line number.
The depth specifies how many stack frames above lives the source line to be identified in the log message.

Log lines have this form:
	[IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line(func)] ms
where the fields are defined as follows:
	L                A single character, representing the log level (eg 'I' for INFO)
	mm               The month (zero padded; ie May is '05')
	dd               The day (zero padded)
	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
	threadid         The space-padded thread ID as returned by GetTID()
	file             The file name
	line             The line number
	msg              The user-supplied message
*/
func (f *GlogFormatter) header(entry *logrus.Entry, depth int, levelColor int, levelText string) []byte {
	var function string
	var fileline string
	if !entry.HasCaller() {
		_, file, line, ok := runtime.Caller(3 + depth)
		if !ok {
			file = "???"
			line = 1
		} else {
			slash := strings.LastIndex(file, "/")
			if slash >= 0 {
				file = file[slash+1:]
			}
		}
		fileline = fmt.Sprintf("%s:%d", file, line)
	} else {
		var file = "???"
		if f.CallerPrettyfier != nil {
			function, file = f.CallerPrettyfier(entry.Caller)
		} else {
			function = entry.Caller.Function
			file = entry.Caller.File
			line := entry.Caller.Line
			if line < 0 {
				line = 0 // not a real line number, but acceptable to someDigits
			}
			slash := strings.LastIndex(function, ".")
			if slash >= 0 {
				function = function[slash+1:]
			}
			slash = strings.LastIndex(file, "/")
			if slash >= 0 {
				file = file[slash+1:]
			}
			fileline = fmt.Sprintf("%s:%d", file, line)
		}
	}
	return f.formatHeader(entry, levelColor, levelText, fileline, function)
}

func (f *GlogFormatter) level(level logrus.Level) (levelColor int, levelText string) {
	if level > logrus.TraceLevel {
		level = logrus.InfoLevel // for safety.
	}
	if f.isColored() {
		switch level {
		case logrus.DebugLevel, logrus.TraceLevel:
			levelColor = gray
		case logrus.WarnLevel:
			levelColor = yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = red
		default:
			levelColor = blue
		}
	}

	levelText = f.levelString(level)
	{
		limit := f.LevelTruncationLimit
		if limit > f.levelTextMaxLength {
			limit = f.levelTextMaxLength
		}
		if limit == 0 {
			limit = 1
		}
		if limit < 0 {
			limit = f.levelTextMaxLength
		}
		if limit > 0 && limit < len(levelText) {
			levelText = levelText[0:limit]
		}
		if f.PadLevelText {
			// Generates the format string used in the next line, for example "%-6s" or "%-7s".
			// Based on the max level text length.
			formatString := "%-" + strconv.Itoa(limit) + "s"
			// Formats the level text by appending spaces up to the max length, for example:
			// 	- "INFO   "
			//	- "WARNING"
			levelText = fmt.Sprintf(formatString, levelText)
		}
	}
	return levelColor, levelText
}

// formatHeader formats a log header using the provided file name and line number.
func (f *GlogFormatter) formatHeader(entry *logrus.Entry, levelColor int, levelText, fileline string, function string) []byte {
	var buf bytes.Buffer
	// Avoid Fprintf, for speed. The format is so simple that we can do it quickly by hand.
	// It's worth about 3X. Fprintf is hard.
	if levelColor > 0 {
		switch {
		case f.DisableTimestamp:
			buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, levelText))
		case f.SinceStartTimestamp:
			buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m[%04d]", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second)))
		default:
			buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m[%s]", levelColor, levelText,
				entry.Time.Format(strings_.ValueOrDefault(f.TimestampFormat, time_.GLogDate))))
		}
	} else {
		// Log line format: [IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg
		// I20200308 23:47:32.089828 400441 config.cc:27] Loading user configuration: /home/aesophor/.config/wmderland/config
		switch {
		case f.DisableTimestamp:
			if f.HumanReadable {
				buf.WriteString(fmt.Sprintf("[%s]", levelText))
			} else {
				buf.WriteString(fmt.Sprintf("%s", levelText))
			}
		case f.SinceStartTimestamp:
			if f.HumanReadable {
				buf.WriteString(fmt.Sprintf("[%s] [%04d]", levelText, int(entry.Time.Sub(baseTimestamp)/time.Second)))
			} else {
				buf.WriteString(fmt.Sprintf("%s%04d", levelText, int(entry.Time.Sub(baseTimestamp)/time.Second)))
			}
		default:
			layout := strings_.ValueOrDefault(f.TimestampFormat, time_.GLogDate)
			var formatString string
			if f.HumanReadable {
				formatString = "[%s] [%s]"
			} else {
				formatString = "%s%s"
			}
			buf.WriteString(fmt.Sprintf(formatString, levelText, entry.Time.Format(layout)))

		}
	}

	if f.ForceGoroutineId {
		if f.HumanReadable {
			buf.WriteString(fmt.Sprintf(" [%-3d]", goroutine.ID()))
		} else {
			buf.WriteString(fmt.Sprintf(" %-3d", goroutine.ID()))
		}
	} else {
		if f.HumanReadable {
			buf.WriteString(fmt.Sprintf(" [%d]", f.pid))
		} else {
			buf.WriteString(fmt.Sprintf(" %d", f.pid))
		}
	}
	if f.WithFuncName && function != "" {
		if f.HumanReadable {
			buf.WriteString(fmt.Sprintf(" [%s](%s)", fileline, function))
		} else {
			buf.WriteString(fmt.Sprintf(" %s(%s)]", fileline, function))
		}
	} else {
		if f.HumanReadable {
			buf.WriteString(fmt.Sprintf(" [%s]", fileline))
		} else {
			buf.WriteString(fmt.Sprintf(" %s]", fileline))
		}
	}
	buf.WriteString(" ")
	return buf.Bytes()
}
