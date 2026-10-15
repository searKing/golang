# redislog

Structured logging for [go-redis v9](https://github.com/redis/go-redis) clients by `log/slog`,
modeled on the slog logger of [gorm](https://gorm.io/docs/logger.html), with a redis command in place of SQL.

```go
rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
rdb.AddHook(redislog.NewHook(slog.Default(), redislog.DefaultConfig))
```

Each dial, command and pipeline is logged in a single line, with fields grouped by `redis`:

```text
level=DEBUG msg="redis executed" redis.name=set redis.cmd="set k1 \"hello world\"" redis.duration=0.123ms
level=WARN msg="redis executed" redis.name=incr redis.cmd="incr h" redis.duration=0.101ms redis.error="WRONGTYPE Operation against a key holding the wrong kind of value"
level=WARN msg="redis pipeline executed" redis.name=[get] redis.cmd="get k1\nget k2" redis.duration=120.000ms redis.slow=true
level=DEBUG msg="redis dialed" redis.network=tcp redis.addr=localhost:6379 redis.duration=0.512ms
```

- `name` is the name of the command, with its subcommand if any, such as `config get`, for aggregations,
  or unique names of commands of pipelines.
- `cmd` is the command as a line of redis-cli, quoted and escaped as redis-cli parses, which can be run by redis-cli
  as is, or commands of pipelines as lines, logged as an array by JSON handlers, or joined by `\n` by text handlers.
  Credentials are redacted as `xxxxx`, as `rediscmd` of redisotel does: passwords of `AUTH`, `HELLO` and `MIGRATE`,
  values of `requirepass`, `masterauth` and TLS passphrases of `CONFIG SET`, and password rules of `ACL SETUSER`.
  Commands are truncated as `SLOWLOG` of redis does: up to 32 arguments of 128 bytes each, and 32 commands.
- Logs are sourced, if `AddSource` of handlers is enabled, from callers of redis, such as the caller of `rdb.Get`,
  rather than the hook, as redisotel locates callers for `code.function`; dials of connection pools have no source.
- Dials and handshakes, such as `HELLO`, are logged on the first command of a connection. Pipelines run on a pool
  of their own since go-redis v9.23.0, so the first pipeline dials again, unless `PipelinePoolSize` is -1.

`Config` corresponds to `Config` of gorm:

| gorm                        | Config              | Default                                                           |
|-----------------------------|---------------------|-------------------------------------------------------------------|
| `SlowThreshold`             | `SlowThreshold`     | 100ms of `DefaultConfig`, disabled if 0; blocking commands, such as `BLPOP`, are never slow |
| `ParameterizedQueries`      | `Parameterized`     | full commands; parameterized as `set k1 ?` if true, safe in production |
| `IgnoreRecordNotFoundError` | `Levels`            | `DefaultErrorToLevel`: debug for success, `redis.Nil`, `NOSCRIPT` and canceled; error for `NOAUTH`, `WRONGPASS`, `NOPERM` and `OOM`; warn for others |
| `LogLevel`                  | levels of handlers  |                                                                   |
|                             | `FieldsFromContext` | fields of contexts, such as tenants, logged out of the group      |
|                             | `Skip`              | all commands logged; skips commands or pipelines, such as `PING`, as `WithCommandFilter` of redisotel |
|                             | `SkipDials`         | dials logged; skips dials if true, as `WithDialFilter` of redisotel |

Slow calls are logged at warn level at least. Fields of clients, such as addresses, can be added by `slog.Logger.With`.

### Internal logs of go-redis

go-redis logs internally, such as failures of the connection pool, reconnections and handoffs of maintenance
notifications, by `Printf` of the global logger set by `redis.SetLogger`, rather than by hooks.
`InternalLogger` logs them to slog, at warn level, for go-redis logs without levels:

```go
redis.SetLogger(redislog.InternalLogger(logger))
// go-redis filters internal logs by its own level, error by default, before they reach the logger.
logging.SetLogLevel(logging.LogLevelInfo)
rdb.AddHook(redislog.NewHook(logger, redislog.DefaultConfig))
```

`WithSlogLogger(h, config)` does both `redis.SetLogger` and `NewHook`, and logs to `slog.Default()` of the time
of each log if `h` is nil, following `slog.SetDefault` called later. The package
[logging](https://pkg.go.dev/github.com/redis/go-redis/v9/logging) of go-redis sets the level, disables internal
logs by `logging.Disable()`, or filters them by substrings by `logging.NewBlacklistLogger` and
`logging.NewWhitelistLogger`, which log by the standard `log` package rather than slog.

### Cluster and ring

`ClusterClient` and `Ring` do not dial by themselves, but by clients of their nodes, so either:

```go
// log each call once, with redirections, such as MOVED, retried inside, but without dials and addresses of nodes;
cluster.AddHook(redislog.NewHook(logger, redislog.DefaultConfig))
// or log calls of each node, as redisotel does, with dials, retries of redirections and addresses of nodes,
// and commands of the cluster itself, such as CLUSTER SLOTS.
cluster.OnNewNode(func(node *redis.Client) {
	node.AddHook(redislog.NewHook(logger.With("node", node.Options().Addr), redislog.DefaultConfig))
})
```

## Tests

The module depends on go-redis only. Tests against a redis server ([miniredis](https://github.com/alicebob/miniredis))
live in the nested module `test`, so that users of this module do not download test dependencies:

```sh
go test ./...                       # unit tests
(cd test && go test ./...)          # tests against miniredis and verified examples
```
