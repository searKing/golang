// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"context"
	"log/slog"
	"os"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"

	redis_ "github.com/searKing/golang/third_party/github.com/redis/go-redis-v9"
)

// exampleLogger logs to stdout without time and durations, to keep the output stable,
// at debug level, as calls succeeded of redis clients are logged at debug level.
func exampleLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey, "duration", "addr":
				return slog.Attr{}
			}
			return a
		},
	}))
}

type tenantKey struct{}

// This example verifies the logs of ExampleNewHook in the parent module against a miniredis server.
// Dials and handshakes, such as HELLO, are logged on the first command,
// where dials are logged without fields of the caller, as go-redis dials with a context of its own.
func Example_newHook() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr:                     mr.Addr(),
		DisableIdentity:          true,
		MaintNotificationsConfig: &maintnotifications.Config{Mode: maintnotifications.ModeDisabled},
	})
	defer rdb.Close()
	config := redis_.DefaultConfig
	config.FieldsFromContext = func(ctx context.Context) []any {
		if tenant, ok := ctx.Value(tenantKey{}).(string); ok {
			return []any{"tenant", tenant}
		}
		return nil
	}
	rdb.AddHook(redis_.NewHook(exampleLogger(), config))

	ctx := context.WithValue(context.Background(), tenantKey{}, "t1")
	_ = rdb.Set(ctx, "k1", "hello world", 0).Err()
	_ = rdb.Get(ctx, "missing").Err()

	// Output:
	// level=DEBUG msg="redis dialed" redis.network=tcp
	// level=DEBUG msg="redis executed" tenant=t1 redis.cmd="hello 3"
	// level=DEBUG msg="redis executed" tenant=t1 redis.cmd="set k1 \"hello world\""
	// level=DEBUG msg="redis executed" tenant=t1 redis.cmd="get missing" redis.error="redis: nil"
}
