// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/redis/go-redis/v9"
)

// ErrorToLevel maps the error of a redis dial, command or pipeline, nil if succeeded, to the log level.
type ErrorToLevel func(err error) slog.Level

// DefaultErrorToLevel maps the error to the log level, by helpers of go-redis:
//   - debug for results expected by callers: nil, redis.Nil, context.Canceled,
//     and NOSCRIPT, which redis.Script recovers from by EVAL;
//   - error for failures to be fixed by operators: errors of authentication, such as NOAUTH and WRONGPASS,
//     of permissions, NOPERM, and of memory, OOM;
//   - warn for others.
func DefaultErrorToLevel(err error) slog.Level {
	switch {
	case err == nil, errors.Is(err, redis.Nil), errors.Is(err, context.Canceled), redis.HasErrorPrefix(err, "NOSCRIPT"):
		return slog.LevelDebug
	case redis.IsAuthError(err), redis.IsPermissionError(err), redis.IsOOMError(err),
		// Is*Error of go-redis misses prefixes following ERR added by KVRocks.
		redis.HasErrorPrefix(err, "NOAUTH"), redis.HasErrorPrefix(err, "WRONGPASS"),
		redis.HasErrorPrefix(err, "NOPERM"), redis.HasErrorPrefix(err, "OOM"):
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}

// code classifies the result of a redis call, logged as `redis.code`, by the prefix of the error replied by
// the redis server, such as WRONGTYPE, or by the outcome of the redis client otherwise, such as NIL for redis.Nil.
type code string

// Outcomes of the redis client.
const (
	codeOK          code = "OK"           // succeeded
	codeNil         code = "NIL"          // redis.Nil, such as GET of a missing key
	codeCanceled    code = "CANCELED"     // the context is canceled
	codeClosed      code = "CLOSED"       // redis.ErrClosed, the client is closed
	codeTimeout     code = "TIMEOUT"      // the deadline of the context, or the read or write timeout, is exceeded
	codePoolTimeout code = "POOL_TIMEOUT" // redis.ErrPoolTimeout, no connection of the pool is available in time
	codeIO          code = "IO"           // the connection fails to dial, read or write
	codeUnknown     code = "UNKNOWN"      // other errors of the client
)

// codeErr is the generic prefix of errors replied by the redis server, also reported for errors without a prefix.
// Errors replied with other prefixes are reported as their prefixes, such as WRONGTYPE.
//
// See https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-errors
const codeErr code = "ERR"

// errorToCode maps errors replied by the redis server to their prefixes, such as WRONGTYPE,
// and errors of the redis client to its outcomes, such as NIL for redis.Nil.
func errorToCode(err error) code {
	switch {
	case err == nil:
		return codeOK
	case errors.Is(err, redis.Nil): // redis.Nil is a redis.Error too.
		return codeNil
	case errors.Is(err, context.Canceled):
		return codeCanceled
	case errors.Is(err, redis.ErrClosed):
		return codeClosed
	case errors.Is(err, context.DeadlineExceeded):
		return codeTimeout
	case errors.Is(err, redis.ErrPoolTimeout):
		return codePoolTimeout
	}
	if redisErr, ok := errors.AsType[redis.Error](err); ok {
		return errorPrefix(redisErr.Error())
	}
	if netErr, ok := errors.AsType[net.Error](err); ok {
		if netErr.Timeout() {
			return codeTimeout
		}
		return codeIO
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return codeIO
	}
	return codeUnknown
}

// standardErrorPrefixes are prefixes of errors defined by redis, other than ERR.
var standardErrorPrefixes = map[string]struct{}{
	"ASK": {}, "BUSY": {}, "BUSYGROUP": {}, "BUSYKEY": {}, "CLUSTERDOWN": {}, "CROSSSLOT": {}, "EXECABORT": {},
	"LOADING": {}, "MASTERDOWN": {}, "MISCONF": {}, "MOVED": {}, "NOAUTH": {}, "NOGROUP": {}, "NOPERM": {},
	"NOPROTO": {}, "NOREPLICAS": {}, "NOSCRIPT": {}, "NOTBUSY": {}, "OOM": {}, "READONLY": {}, "TRYAGAIN": {},
	"UNBLOCKED": {}, "UNKILLABLE": {}, "WRONGPASS": {}, "WRONGTYPE": {},
}

// errorPrefix returns the prefix of the error replied by the redis server, the first word in upper case,
// or ERR if missing.
// A standard prefix following ERR is returned instead, such as WRONGTYPE of "ERR WRONGTYPE ...",
// as KVRocks adds ERR to prefixes, which redis.HasErrorPrefix trims too,
// whereas other words, such as DB of "ERR DB index is out of range", are messages of ERR.
func errorPrefix(msg string) code {
	prefix, rest, _ := strings.Cut(msg, " ")
	if !isErrorPrefix(prefix) {
		return codeErr
	}
	if prefix == string(codeErr) {
		next, _, _ := strings.Cut(rest, " ")
		if _, ok := standardErrorPrefixes[next]; ok {
			return code(next)
		}
	}
	return code(prefix)
}

// isErrorPrefix reports whether s is a prefix of errors of redis, a word in upper case.
func isErrorPrefix(s string) bool {
	return s != "" && !strings.ContainsFunc(s, func(r rune) bool {
		return (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' && r != '-'
	})
}
