// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"log/slog"
	"net"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ redis.Hook = (*hook)(nil)

var (
	// SystemTag is tag representing an event inside redis call.
	SystemTag = []string{"protocol", "redis"}
	// ComponentFieldKey is a tag representing the client that is calling, valued as KindClientFieldValue.
	ComponentFieldKey    = "redis.component"
	KindClientFieldValue = "client"
	// CallTypeFieldKey is a tag representing the type of the call, valued as CallType.
	CallTypeFieldKey = "redis.call_type"
	// NetworkFieldKey and AddrFieldKey are tags representing the network and address dialed, of dials only.
	NetworkFieldKey = "redis.network"
	AddrFieldKey    = "redis.addr"
	// CommandFieldKey and KeyFieldKey are tags representing the full name and the first key of the command,
	// of commands only, where KeyFieldKey is left out if the command has no key.
	CommandFieldKey = "redis.command"
	KeyFieldKey     = "redis.key"
	// CommandsFieldKey and KeysFieldKey are tags representing full names and first keys of commands,
	// of pipelines only, where keys are "" for commands without keys, truncated to 32 commands as SLOWLOG does.
	// PipelineLengthFieldKey is a tag representing the number of commands, of pipelines only,
	// including MULTI and EXEC of transactions.
	CommandsFieldKey       = "redis.commands"
	KeysFieldKey           = "redis.keys"
	PipelineLengthFieldKey = "redis.pipeline_length"
	// DBFieldKey is a tag representing the database index configured of the client, by Instrument only,
	// which may differ from the database selected by SELECT.
	DBFieldKey = "redis.db"
)

// CallType is the type of the call hooked, as interceptors.GRPCType of go-grpc-middleware.
type CallType string

const (
	Dial     CallType = "dial"
	Command  CallType = "command"
	Pipeline CallType = "pipeline"
)

// CallMeta is the meta of the call hooked, as interceptors.CallMeta of go-grpc-middleware.
type CallMeta struct {
	Typ CallType
	// Network and Addr are the network and address dialed, of dials only.
	Network string
	Addr    string
	// Cmds are the commands of the call, a single command for commands, nil for dials.
	Cmds []redis.Cmder
}

// Hook returns a new redis.Hook that optionally logs dials, commands and pipelines of go-redis clients to l,
// slog.Default() if nil.
// Logs and fields are consistent with logging.UnaryClientInterceptor of go-grpc-middleware for gRPC,
// except that results are logged as `redis.code` by codes of redis rather than of gRPC:
//   - the prefix of errors replied by the redis server, such as WRONGTYPE, NOAUTH or MOVED,
//     or ERR if missing;
//   - otherwise the outcome of the redis client: OK, NIL for redis.Nil, CANCELED, CLOSED for redis.ErrClosed,
//     TIMEOUT, POOL_TIMEOUT for redis.ErrPoolTimeout, IO for failures of connections, or UNKNOWN.
//
// Logs are leveled by DefaultErrorToLevel, customized by WithLevels.
// Each call is logged in a single line by default, on FinishCall, customized by WithLogOnEvents.
// Commands are logged with names and keys only, leaving other arguments and replies out of logs,
// as they may carry credentials or bulky data, unless enabled by WithPayloads,
// or by REDIS_GO_LOG_PAYLOAD if true, with credentials redacted.
func Hook(l *slog.Logger, opts ...LoggingOption) redis.Hook {
	return &hook{logger: l, opts: evaluateLoggingOptions(opts)}
}

// Instrument adds Hook to rdb, logging to l, slog.Default() if nil, with fields of the configuration of rdb:
// DBFieldKey and AddrFieldKey of *redis.Client, as the database index and the address configured,
// where commands are not logged with the database selected by SELECT, if differs.
// Other clients, such as *redis.ClusterClient, are added Hook only.
func Instrument(rdb redis.UniversalClient, l *slog.Logger, opts ...LoggingOption) {
	if c, ok := rdb.(*redis.Client); ok {
		opt := c.Options()
		opts = append([]LoggingOption{withClientFields(Fields{AddrFieldKey, opt.Addr, DBFieldKey, opt.DB})}, opts...)
	}
	rdb.AddHook(Hook(l, opts...))
}

type hook struct {
	logger *slog.Logger
	opts   *loggingOptions
}

func (h *hook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		r := h.reportable(ctx, CallMeta{Typ: Dial, Network: network, Addr: addr})
		conn, err := next(ctx, network, addr)
		r.PostCall(err)
		return conn, err
	}
}

func (h *hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		r := h.reportable(ctx, CallMeta{Typ: Command, Cmds: []redis.Cmder{cmd}})
		err := next(ctx, cmd)
		r.PostCall(err)
		return err
	}
}

func (h *hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		r := h.reportable(ctx, CallMeta{Typ: Pipeline, Cmds: cmds})
		// The pipeline returns the first error of its commands, such as redis.Nil if any command misses its key,
		// while the rest are executed as usual.
		err := next(ctx, cmds)
		r.PostCall(err)
		return err
	}
}

type reporter struct {
	CallMeta

	ctx    context.Context
	start  time.Time
	opts   *loggingOptions
	logger *slog.Logger
	// fields are built on the first log of the call, so that calls of disabled levels, mostly succeeded ones
	// at debug level, cost no fields.
	fields Fields
}

// reportable returns the reporter of the call c, logging StartCall if enabled,
// as reportable of logging interceptors of go-grpc-middleware does.
func (h *hook) reportable(ctx context.Context, c CallMeta) *reporter {
	logger := h.logger
	if logger == nil {
		logger = slog.Default()
	}
	r := &reporter{CallMeta: c, ctx: ctx, start: time.Now(), opts: h.opts, logger: logger}
	if h.opts.has(StartCall) {
		r.log(r.opts.levelFunc(nil), "started call", nil)
	}
	return r
}

// PostCall logs FinishCall with the error of the call, if enabled.
func (r *reporter) PostCall(err error) {
	if !r.opts.has(FinishCall) {
		return
	}
	duration := time.Since(r.start)
	level := r.opts.levelFunc(err)
	slow := r.isSlow(duration)
	if slow {
		level = max(level, slog.LevelWarn)
	}
	r.log(level, "finished call", func() Fields {
		fields := Fields{"redis.code", string(errorToCode(err))}
		if slow {
			fields = append(fields, "redis.slow", true)
		}
		if err != nil {
			fields = fields.appendUnique(Fields{"redis.error", err.Error()})
			if r.opts.errorToFieldsFunc != nil {
				fields = fields.appendUnique(r.opts.errorToFieldsFunc(err))
			}
		}
		if content, ok := r.responseContent(err); ok {
			fields = fields.appendUnique(Fields{"redis.response.content", content})
		}
		return fields.appendUnique(r.opts.durationFieldFunc(duration))
	})
}

// isSlow reports whether the call taking duration is slow, by WithSlowThreshold.
func (r *reporter) isSlow(duration time.Duration) bool {
	if r.opts.slowThreshold <= 0 || duration <= r.opts.slowThreshold {
		return false
	}
	return !slices.ContainsFunc(r.Cmds, isBlocking)
}

// statement returns the commands with keys and subcommands only, and other arguments replaced by "?",
// if StatementPayload is enabled.
func (r *reporter) statement() (any, bool) {
	if r.Typ == Dial || !r.opts.hasPayload(StatementPayload) {
		return nil, false
	}
	if r.Typ == Command {
		return commandStatement(r.Cmds[0]), true
	}
	statements := make([]string, 0, len(r.Cmds))
	for _, cmd := range r.Cmds {
		statements = append(statements, commandStatement(cmd))
	}
	return truncateElems(statements, "commands"), true
}

// requestContent returns arguments of the commands, with credentials redacted, if RequestPayload is enabled.
func (r *reporter) requestContent() (any, bool) {
	if r.Typ == Dial || !r.opts.hasPayload(RequestPayload) {
		return nil, false
	}
	if r.Typ == Command {
		return commandArgs(r.Cmds[0]), true
	}
	args := make([][]any, 0, len(r.Cmds))
	for _, cmd := range r.Cmds {
		args = append(args, commandArgs(cmd))
	}
	return truncateElems(args, "commands"), true
}

// responseContent returns replies of the commands, with credentials redacted, if ResponsePayload is enabled.
// Replies of a failed command are left out, as its error is logged,
// whereas replies of a pipeline are returned even if some commands failed, with their errors in place of replies.
func (r *reporter) responseContent(err error) (any, bool) {
	if r.Typ == Dial || !r.opts.hasPayload(ResponsePayload) || (err != nil && r.Typ == Command) {
		return nil, false
	}
	if r.Typ == Command {
		return commandReply(r.Cmds[0]), true
	}
	replies := make([]any, 0, len(r.Cmds))
	for _, cmd := range r.Cmds {
		replies = append(replies, commandReply(cmd))
	}
	return truncateElems(replies, "commands"), true
}

// log logs msg at level with fields of the call, followed by fields of extra, nil if none,
// where neither is built if level is disabled.
func (r *reporter) log(level slog.Level, msg string, extra func() Fields) {
	if !r.logger.Enabled(r.ctx, level) {
		return
	}
	if r.fields == nil {
		r.fields = r.callFields()
	}
	fields := r.fields
	if extra != nil {
		fields = fields.appendUnique(extra())
	}
	record := slog.NewRecord(time.Now(), level, msg, callerPC(isOfRedis))
	record.Add(fields...)
	_ = r.logger.Handler().Handle(r.ctx, record)
}

// callerPC returns the program counter of the first caller outside frames skipped by skip,
// for the source of logs, as Logger.Log does for its caller,
// or 0 if none, such as dials of connection pools of go-redis in goroutines of their own.
func callerPC(skip func(function string) bool) uintptr {
	var pcs [64]uintptr
	// Skip runtime.Callers, callerPC and the caller of callerPC.
	frames := runtime.CallersFrames(pcs[:runtime.Callers(3, pcs[:])])
	for {
		f, more := frames.Next()
		if !skip(f.Function) && !strings.HasPrefix(f.Function, "runtime.") {
			// The PC of the frame is of the call, while the record takes the return address as runtime.Callers,
			// which works for frames of inlined calls too, of which the PC is of the inlined call.
			return f.PC + 1
		}
		if !more {
			return 0
		}
	}
}

// isOfRedis reports whether the function is of go-redis, including its extra packages, such as redisotel,
// or of this package, skipped to locate the caller of redis, such as the caller of rdb.Get,
// as redisotel does for code.function of spans.
func isOfRedis(function string) bool {
	return strings.HasPrefix(function, "github.com/redis/go-redis/") || isOfThisPackage(function)
}

// isOfThisPackage reports whether the function is of this package, excluding its subpackages.
func isOfThisPackage(function string) bool {
	rest, ok := strings.CutPrefix(function, thisPackage)
	return ok && strings.HasPrefix(rest, ".")
}

// thisPackage is the import path of this package.
var thisPackage = reflect.TypeFor[hook]().PkgPath()

// callFields returns fields shared by all logs of the call.
func (r *reporter) callFields() Fields {
	fields := newCommonFields(r.CallMeta).appendUnique(r.opts.clientFields)
	for _, key := range r.opts.disableLogFields {
		fields.delete(key)
	}
	if r.opts.fieldsFromCtxCallMetaFn != nil {
		// fieldsFromCtxCallMetaFn dups override the existing fields.
		fields = r.opts.fieldsFromCtxCallMetaFn(r.ctx, r.CallMeta).appendUnique(fields)
	}
	fields = fields.appendUnique(Fields{"redis.start_time", r.start.Format(r.opts.timestampFormat)})
	if d, ok := r.ctx.Deadline(); ok {
		fields = fields.appendUnique(Fields{"redis.request.deadline", d.Format(r.opts.timestampFormat)})
	}
	if statement, ok := r.statement(); ok {
		fields = fields.appendUnique(Fields{"redis.statement", statement})
	}
	if content, ok := r.requestContent(); ok {
		fields = fields.appendUnique(Fields{"redis.request.content", content})
	}
	// Clipped, as fields are appended by every log of the call, which may reuse the array otherwise.
	return slices.Clip(fields)
}

// newCommonFields returns the common fields of the call c,
// as newCommonFields of logging interceptors of go-grpc-middleware does.
func newCommonFields(c CallMeta) Fields {
	fields := Fields{
		SystemTag[0], SystemTag[1],
		ComponentFieldKey, KindClientFieldValue,
		CallTypeFieldKey, string(c.Typ),
	}
	switch c.Typ {
	case Dial:
		fields = append(fields, NetworkFieldKey, c.Network, AddrFieldKey, c.Addr)
	case Command:
		fields = append(fields, CommandFieldKey, c.Cmds[0].FullName())
		if key := commandKey(c.Cmds[0]); key != "" {
			fields = append(fields, KeyFieldKey, key)
		}
	case Pipeline:
		names := make([]string, 0, len(c.Cmds))
		keys := make([]string, 0, len(c.Cmds))
		for _, cmd := range c.Cmds {
			names = append(names, cmd.FullName())
			keys = append(keys, commandKey(cmd))
		}
		fields = append(fields, CommandsFieldKey, truncateElems(names, "commands"),
			KeysFieldKey, truncateElems(keys, "commands"), PipelineLengthFieldKey, len(c.Cmds))
	}
	return fields
}
