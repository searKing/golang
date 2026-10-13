# go-redis-v9

Structured logging for [go-redis v9](https://github.com/redis/go-redis) clients by `log/slog`,
modeled on the logging interceptors of go-grpc-middleware, without depending on grpc.

```go
rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
rdb.AddHook(redis_.Hook(slog.Default()))
// or, logging the address and the database index configured of rdb, as redisotel does:
redis_.Instrument(rdb, slog.Default())
```

- Dials, commands and pipelines are logged in a single line each, as `finished call`;
  `started call` is logged too if enabled by `WithLogOnEvents(StartCall, FinishCall)`, such as for blocking commands.
- Commands are logged with names and first keys by default, `redis.command` and `redis.key`,
  or `redis.commands`, `redis.keys` and `redis.pipeline_length` of pipelines.
- Commands sanitized, with keys and subcommands only and other arguments replaced by `?`, such as `set k1 ?`,
  as `db.query.text` of OpenTelemetry, are logged as `redis.statement` if enabled by `WithPayloads(StatementPayload)`.
- Arguments and replies of commands are logged as fields `redis.request.content` and `redis.response.content`
  of the call only if enabled by `WithPayloads`, or by the environment variable `REDIS_GO_LOG_PAYLOAD=true`
  read on `Hook`, as `HTTP_GO_LOG_HTTP_HEADER` does for headers of HTTP logs; `WithPayloads` takes precedence.
- Results are logged as `redis.code`: the prefix of errors replied by the server,
  such as `WRONGTYPE` or `NOAUTH`, or the outcome of the client, such as `OK`, `NIL`, `TIMEOUT` or `IO`.
- Logs are leveled from errors by `DefaultErrorToLevel`, customized by `WithLevels`:
  `OK`, `NIL`, `CANCELED` and `NOSCRIPT` at debug level,
  `NOAUTH`, `WRONGPASS`, `NOPERM` and `OOM` at error level, and others at warn level.
- Calls slower than `WithSlowThreshold`, 100ms by default and disabled by `WithSlowThreshold(0)`, are logged at warn level at least with `redis.slow=true`,
  as `SlowThreshold` of gorm does, except blocking commands, such as `BLPOP` or `XREAD BLOCK`, as `SLOWLOG` does.
- Credentials of `AUTH`, `HELLO`, `ACL`, `CONFIG` and `MIGRATE` are redacted.
- Logs are sourced, if `AddSource` of handlers is enabled, from callers of redis, such as the caller of `rdb.Get`,
  rather than the hook, as redisotel locates callers for `code.function`; dials of connection pools have no source.
- Payloads are truncated as `SLOWLOG` of redis does: up to 32 arguments, elements or commands,
  and up to 128 bytes of each.

## Tests

The module depends on go-redis only. Tests against a redis server ([miniredis](https://github.com/alicebob/miniredis))
live in the nested module `test`, so that users of this module do not download test dependencies:

```sh
go test ./...                       # unit tests
(cd test && go test ./...)          # tests against miniredis and verified examples
```
