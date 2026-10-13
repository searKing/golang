// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"errors"
	"log/slog"

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
