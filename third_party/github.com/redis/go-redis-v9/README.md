# go-redis-v9

Structured logging for [go-redis v9](https://github.com/redis/go-redis) clients by `log/slog`,
modeled on the slog logger of [gorm](https://gorm.io/docs/logger.html), with a redis command in place of SQL.

```go
rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
rdb.AddHook(redis_.NewHook(slog.Default(), redis_.DefaultConfig))
```

Each dial, command and pipeline is logged in a single line, with fields grouped by `redis`:

```text
level=DEBUG msg="redis executed" redis.cmd="set k1 \"hello world\"" redis.duration=0.123ms
level=WARN msg="redis executed" redis.cmd="incr h" redis.duration=0.101ms redis.error="WRONGTYPE Operation against a key holding the wrong kind of value"
level=WARN msg="redis pipeline executed" redis.cmd="get k1\nget k2" redis.duration=120.000ms redis.slow=true
level=DEBUG msg="redis dialed" redis.network=tcp redis.addr=localhost:6379 redis.duration=0.512ms
```

- `cmd` is the command as a line of redis-cli, quoted and escaped as redis-cli parses, which can be run by redis-cli
  as is, or commands of pipelines as lines, logged as an array by JSON handlers, or joined by `\n` by text handlers.
  Arguments which may carry credentials, of `AUTH`, `HELLO`, `ACL`, `CONFIG` and `MIGRATE`, are redacted as `xxxxx`,
  and commands are truncated as `SLOWLOG` of redis does: up to 32 arguments of 128 bytes each, and 32 commands.
- Logs are sourced, if `AddSource` of handlers is enabled, from callers of redis, such as the caller of `rdb.Get`,
  rather than the hook, as redisotel locates callers for `code.function`; dials of connection pools have no source.
- Internal logs of go-redis are logged by `redis.SetLogger(redis_.InternalLogger(logger))`, or by `WithSlogLogger`.

`Config` corresponds to `Config` of gorm:

| gorm                        | Config              | Default                                                           |
|-----------------------------|---------------------|-------------------------------------------------------------------|
| `SlowThreshold`             | `SlowThreshold`     | 100ms of `DefaultConfig`, disabled if 0; blocking commands, such as `BLPOP`, are never slow |
| `ParameterizedQueries`      | `Parameterized`     | full commands; parameterized as `set k1 ?` if true, safe in production |
| `IgnoreRecordNotFoundError` | `Levels`            | `DefaultErrorToLevel`: debug for success, `redis.Nil`, `NOSCRIPT` and canceled; error for `NOAUTH`, `WRONGPASS`, `NOPERM` and `OOM`; warn for others |
| `LogLevel`                  | levels of handlers  |                                                                   |
|                             | `FieldsFromContext` | fields of contexts, such as tenants, logged out of the group      |

Slow calls are logged at warn level at least. Fields of clients, such as addresses, can be added by `slog.Logger.With`.

## Tests

The module depends on go-redis only. Tests against a redis server ([miniredis](https://github.com/alicebob/miniredis))
live in the nested module `test`, so that users of this module do not download test dependencies:

```sh
go test ./...                       # unit tests
(cd test && go test ./...)          # tests against miniredis and verified examples
```
