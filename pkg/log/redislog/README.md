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
- Internal logs of go-redis are logged by `redis.SetLogger(redislog.InternalLogger(logger))`, or by `WithSlogLogger`.

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
