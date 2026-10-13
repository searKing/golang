// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

// Fields loosely represents key value pairs that adds context to log lines, as logging.Fields of go-grpc-middleware.
// The key has to be type of string, whereas value can be an arbitrary object.
type Fields []any

func (f Fields) has(key any) bool {
	for i := 0; i+1 < len(f); i += 2 {
		if f[i] == key {
			return true
		}
	}
	return false
}

// appendUnique appends fields of add which do not occur in f, where the first occurrence takes precedence.
func (f Fields) appendUnique(add Fields) Fields {
	for i := 0; i+1 < len(add); i += 2 {
		if !f.has(add[i]) {
			f = append(f, add[i], add[i+1])
		}
	}
	return f
}

func (f *Fields) delete(key string) {
	for i := 0; i+1 < len(*f); i += 2 {
		if (*f)[i] == key {
			*f = slices.Delete(*f, i, i+2)
			return
		}
	}
}

// LoggableEvent defines the events a log line can be added on, as logging.LoggableEvent of go-grpc-middleware,
// without events of payloads, as a redis call sends its commands and receives their replies in a single round trip,
// of which payloads are logged as fields of the call, enabled by WithPayloads.
type LoggableEvent uint

const (
	// StartCall is a loggable event representing start of the redis call,
	// for calls which may block, such as BLPOP, or hang.
	StartCall LoggableEvent = iota
	// FinishCall is a loggable event representing finish of the redis call.
	FinishCall
)

// Payload defines the payloads of commands logged as fields of calls, of commands and pipelines only.
// NOTE: Payloads can get quite verbose and sensitive, use with caution (e.g. debug only purposes).
//
// Payloads are truncated as SLOWLOG of redis does, up to 32 arguments, elements or commands,
// with the rest replaced by "... (N more arguments)", and up to 128 bytes of strings and bytes,
// followed by "... (N more bytes)".
type Payload uint

const (
	// RequestPayload logs commands in "redis.request.content" field, of all logs of the call,
	// as lines of redis-cli, quoted and escaped as redis-cli parses, such as `set k1 "hello world"`,
	// which can be run by redis-cli as is, and lines of commands of pipelines,
	// with arguments which may carry credentials redacted, such as passwords of AUTH.
	RequestPayload Payload = iota
	// ResponsePayload logs replies of commands in "redis.response.content" field, of FinishCall,
	// with those which may carry credentials redacted, such as passwords of CONFIG GET requirepass.
	ResponsePayload
	// StatementPayload logs commands in "redis.statement" field, of all logs of the call,
	// with keys and subcommands only, and other arguments replaced by "?", such as "set k1 ?",
	// as db.query.text of OpenTelemetry for redis, which is sanitized and safe to log in production.
	StatementPayload
)

// LoggingOption customizes Hook, as logging.Option of go-grpc-middleware does for logging interceptors of gRPC.
type LoggingOption func(*loggingOptions)

type loggingOptions struct {
	levelFunc               ErrorToLevel
	slowThreshold           time.Duration
	loggableEvents          []LoggableEvent
	payloads                []Payload
	errorToFieldsFunc       ErrorToFields
	durationFieldFunc       DurationToFields
	timestampFormat         string
	fieldsFromCtxCallMetaFn func(ctx context.Context, c CallMeta) Fields
	disableLogFields        []string
	// clientFields are fields of the configuration of the client, by Instrument.
	clientFields Fields
}

func evaluateLoggingOptions(opts []LoggingOption) *loggingOptions {
	o := &loggingOptions{
		levelFunc:         DefaultErrorToLevel,
		slowThreshold:     DefaultSlowThreshold,
		loggableEvents:    []LoggableEvent{FinishCall},
		payloads:          payloadsFromEnv(),
		durationFieldFunc: DefaultDurationToFields,
		timestampFormat:   time.RFC3339,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *loggingOptions) has(event LoggableEvent) bool {
	return slices.Contains(o.loggableEvents, event)
}

func (o *loggingOptions) hasPayload(payload Payload) bool {
	return slices.Contains(o.payloads, payload)
}

// DurationToFields function defines how to produce duration fields for logging.
type DurationToFields func(duration time.Duration) Fields

// ErrorToFields function extract fields from error.
type ErrorToFields func(err error) Fields

// DefaultDurationToFields is the default implementation of converting request duration to a field.
var DefaultDurationToFields = DurationToTimeMillisFields

// DurationToTimeMillisFields converts the duration to milliseconds and uses the key `redis.time_ms`,
// as logging.DurationToTimeMillisFields of go-grpc-middleware does with the key `grpc.time_ms`.
func DurationToTimeMillisFields(duration time.Duration) Fields {
	return Fields{"redis.time_ms", fmt.Sprintf("%v", float32(duration.Nanoseconds()/1000)/1000)}
}

// DurationToDurationField uses a Duration field to log the request duration with the key `redis.duration`,
// as logging.DurationToDurationField of go-grpc-middleware does with the key `grpc.duration`.
func DurationToDurationField(duration time.Duration) Fields {
	return Fields{"redis.duration", duration.String()}
}

// WithFieldsFromContext allows overriding existing or adding extra fields to all log messages per given context.
// f is called once per call, on its first log, and not at all if no log of the call is enabled.
// Only one of WithFieldsFromContext or WithFieldsFromContextAndCallMeta should be used,
// using both will result in the last one overwriting the previous.
func WithFieldsFromContext(f func(ctx context.Context) Fields) LoggingOption {
	return func(o *loggingOptions) {
		o.fieldsFromCtxCallMetaFn = func(ctx context.Context, _ CallMeta) Fields {
			return f(ctx)
		}
	}
}

// WithFieldsFromContextAndCallMeta allows overriding existing or adding extra fields to all log messages
// per given context and CallMeta.
// f is called once per call, on its first log, and not at all if no log of the call is enabled.
// Only one of WithFieldsFromContext or WithFieldsFromContextAndCallMeta should be used,
// using both will result in the last one overwriting the previous.
func WithFieldsFromContextAndCallMeta(f func(ctx context.Context, c CallMeta) Fields) LoggingOption {
	return func(o *loggingOptions) {
		o.fieldsFromCtxCallMetaFn = f
	}
}

// WithLogOnEvents customizes on what events the redis hook should log on, FinishCall only by default.
func WithLogOnEvents(events ...LoggableEvent) LoggingOption {
	return func(o *loggingOptions) {
		o.loggableEvents = events
	}
}

// payloadsFromEnv returns all payloads if REDIS_GO_LOG_PAYLOAD is true, as HTTP_GO_LOG_HTTP_HEADER does
// for headers of logging.http.go of grpc-gateway-v2, or none otherwise.
func payloadsFromEnv() []Payload {
	if v, _ := strconv.ParseBool(os.Getenv("REDIS_GO_LOG_PAYLOAD")); v {
		return []Payload{RequestPayload, ResponsePayload}
	}
	return nil
}

// WithPayloads enables payloads of commands logged as fields of calls, overriding REDIS_GO_LOG_PAYLOAD,
// which enables all payloads if true, read by Hook. None are logged by default.
// WithPayloads() without payloads disables payloads, even if REDIS_GO_LOG_PAYLOAD is true.
func WithPayloads(payloads ...Payload) LoggingOption {
	return func(o *loggingOptions) {
		o.payloads = payloads
	}
}

// withClientFields adds fields of the configuration of the client to all logs.
func withClientFields(fields Fields) LoggingOption {
	return func(o *loggingOptions) {
		o.clientFields = fields
	}
}

// WithErrorFields allows to extract logging fields from an error of redis.
func WithErrorFields(f ErrorToFields) LoggingOption {
	return func(o *loggingOptions) {
		o.errorToFieldsFunc = f
	}
}

// WithLevels customizes the function for mapping errors of calls to log levels, DefaultErrorToLevel by default,
// where events other than FinishCall are logged at the level of nil.
func WithLevels(f ErrorToLevel) LoggingOption {
	return func(o *loggingOptions) {
		o.levelFunc = f
	}
}

// DefaultSlowThreshold is the default threshold of WithSlowThreshold.
// It is looser than 10ms of slowlog-log-slower-than of redis, which counts the execution of the server only,
// as calls are timed by the client, including round trips, waits for connections of the pool and dials.
const DefaultSlowThreshold = 100 * time.Millisecond

// WithSlowThreshold logs calls taking longer than threshold on FinishCall at warn level at least,
// with "redis.slow" field of true, as SlowThreshold of the logger of gorm does,
// DefaultSlowThreshold by default, or disabled if not positive.
// Commands blocking by design, such as BLPOP or XREAD with BLOCK, and pipelines of them,
// are never slow, as SLOWLOG of redis does.
func WithSlowThreshold(threshold time.Duration) LoggingOption {
	return func(o *loggingOptions) {
		o.slowThreshold = threshold
	}
}

// WithDurationField customizes the function for mapping request durations to log fields.
func WithDurationField(f DurationToFields) LoggingOption {
	return func(o *loggingOptions) {
		o.durationFieldFunc = f
	}
}

// WithTimestampFormat customizes the timestamps emitted in the log fields.
func WithTimestampFormat(format string) LoggingOption {
	return func(o *loggingOptions) {
		o.timestampFormat = format
	}
}

// WithDisableLoggingFields disables logging of redis fields provided.
// The following are the default logging fields:
//   - SystemTag[0]
//   - ComponentFieldKey
//   - CallTypeFieldKey
//   - NetworkFieldKey and AddrFieldKey, of dials only
//   - CommandFieldKey and KeyFieldKey, of commands only
//   - CommandsFieldKey, KeysFieldKey and PipelineLengthFieldKey, of pipelines only
//   - AddrFieldKey and DBFieldKey, of clients instrumented by Instrument
//
// Usage example - WithDisableLoggingFields(redis.KeyFieldKey, redis.KeysFieldKey)
func WithDisableLoggingFields(disableLogFields ...string) LoggingOption {
	return func(o *loggingOptions) {
		o.disableLogFields = disableLogFields
	}
}
