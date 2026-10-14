// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config configures NewHook, as Config of the logger of gorm.
// The zero Config logs full commands with slow calls disabled; DefaultConfig is recommended instead.
type Config struct {
	// SlowThreshold logs calls slower than it at warn level at least, with the field slow of true,
	// as SlowThreshold of gorm, disabled if not positive.
	// Commands blocking by design, such as BLPOP or XREAD with BLOCK, and pipelines of them,
	// are never slow, as SLOWLOG of redis does not count blocking.
	SlowThreshold time.Duration
	// Parameterized logs commands with keys and subcommands only, and other arguments replaced by "?",
	// such as "set k1 ?", as ParameterizedQueries of gorm, and db.query.text of OpenTelemetry for redis,
	// which is safe to log in production, where values may carry tokens or bulky data.
	Parameterized bool
	// Levels maps errors of calls, nil if succeeded, to log levels, DefaultErrorToLevel if nil,
	// which logs succeeded calls and redis.Nil at debug level, as IgnoreRecordNotFoundError of gorm does.
	Levels ErrorToLevel
	// FieldsFromContext returns fields logged with calls, out of the group "redis", such as the tenant of ctx,
	// as arguments of slog.Logger.Log, called only if the log is enabled.
	// Dials are called with contexts of go-redis, rather than of the callers of commands.
	FieldsFromContext func(ctx context.Context) []any
	// Skip reports whether to skip logging the call of cmds, a single command for commands and nil for dials,
	// such as PING of health checks, as filters of redisotel, logging all calls if nil.
	Skip func(cmds []redis.Cmder) bool
}

// DefaultSlowThreshold is the threshold of slow calls of DefaultConfig.
// It is looser than 10ms of slowlog-log-slower-than of redis, which counts the execution of the server only,
// as calls are timed by the client, including round trips, waits for connections of the pool and dials.
const DefaultSlowThreshold = 100 * time.Millisecond

// DefaultConfig is the recommended Config, as Default of the logger of gorm.
var DefaultConfig = Config{SlowThreshold: DefaultSlowThreshold}
