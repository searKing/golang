// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"

	redis_ "github.com/searKing/golang/third_party/github.com/redis/go-redis-v9"
)

type tenantKey struct{}

// This example logs commands of a redis client, with fields of the caller from the context by WithFieldsFromContext.
// Dials and handshakes, such as HELLO, are logged on the first command,
// where dials are logged without fields of the caller, as go-redis dials with a context of its own.
//
// It needs a redis server, so it is run with a miniredis server in the module test.
func ExampleHook() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	rdb.AddHook(redis_.Hook(logger,
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
}
