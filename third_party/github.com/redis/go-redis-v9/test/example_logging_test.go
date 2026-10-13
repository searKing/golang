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

// exampleLogger logs to stdout only the message and fields of commands, to keep the output stable,
// at debug level, as calls succeeded of redis clients are logged at debug level.
func exampleLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.MessageKey, redis_.CallTypeFieldKey, redis_.CommandFieldKey, redis_.KeyFieldKey, "redis.code", "tenant":
				return a
			}
			return slog.Attr{}
		},
	}))
}

type tenantKey struct{}

// This example verifies the logs of ExampleHook in the parent module against a miniredis server.
// Dials and handshakes, such as HELLO, are logged on the first command,
// where dials are logged without fields of the caller, as go-redis dials with a context of its own.
func Example_hook() {
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
	rdb.AddHook(redis_.Hook(exampleLogger(),
		redis_.WithFieldsFromContext(func(ctx context.Context) redis_.Fields {
			if tenant, ok := ctx.Value(tenantKey{}).(string); ok {
				return redis_.Fields{"tenant", tenant}
			}
			return nil
		}),
	))

	ctx := context.WithValue(context.Background(), tenantKey{}, "t1")
	_ = rdb.Set(ctx, "k1", "v1", 0).Err()
	_ = rdb.Get(ctx, "missing").Err()

	// Output:
	// msg="finished call" redis.call_type=dial redis.code=OK
	// msg="finished call" tenant=t1 redis.call_type=command redis.command=hello redis.code=OK
	// msg="finished call" tenant=t1 redis.call_type=command redis.command=set redis.key=k1 redis.code=OK
	// msg="finished call" tenant=t1 redis.call_type=command redis.command=get redis.key=missing redis.code=NIL
}
