[![Go Reference](https://pkg.go.dev/badge/github.com/searKing/golang.svg)](https://pkg.go.dev/github.com/searKing/golang)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang)](https://goreportcard.com/report/github.com/searKing/golang)
[<img src="https://api.visitorbadge.io/api/visitors?path=https%3A%2F%2Fgithub.com%2FsearKing%2Fgolang&countColor=%23263759" height="20">](https://visitorbadge.io/status?path=https%3A%2F%2Fgithub.com%2FsearKing%2Fgolang)
[<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" height="20">](https://jb.gg/OpenSourceSupport)

# golang

various libs or tools for Golang

# GoLibs

* [exp](https://pkg.go.dev/github.com/searKing/golang/go/exp) — `exp` holds experimental packages and defines various
  functions useful with [generics](https://go.dev/doc/tutorial/generics) of any type.
    - [slices](https://pkg.go.dev/github.com/searKing/golang/go/exp/slices) defines various functions useful with slices
      of any type.
    - [maps](https://pkg.go.dev/github.com/searKing/golang/go/exp/maps) defines various functions useful with maps of
      any type.
    - [image](https://pkg.go.dev/github.com/searKing/golang/go/exp/image)
      defines [Point[T]](https://pkg.go.dev/github.com/searKing/golang/go/exp/image#Point)
      and [Rectangle[T]](https://pkg.go.dev/github.com/searKing/golang/go/exp/image#Rectangle) of any type.
    - [sync.LRU](https://pkg.go.dev/github.com/searKing/golang/go/exp/sync#LRU) implements a thread safe fixed size LRU
      cache, based on [not-thread safe lru](https://pkg.go.dev/github.com/searKing/golang/go@v1.2.82/exp/container/lru)
    - [sync.FixedPool](https://pkg.go.dev/github.com/searKing/golang/go/exp/sync#FixedPool) is a set of resident and
      temporary items that may be individually saved and retrieved.
    - [types.Optional](https://pkg.go.dev/github.com/searKing/golang/go/exp/types#Optional) represents a Value that may
      be null.
* [slog](https://pkg.go.dev/github.com/searKing/log/slog) - `slog`
  provides [GlogHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewGlogHandler),
  [GlogHumanHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewGlogHumanHandler),
  [NewRotateHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewRotateHandler) and
  [MultiHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#MultiHandler) handlers
  for [slog](https://pkg.go.dev/log/slog)
    - [GlogHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewGlogHandler) provides a Handler that
      writes Records to an io.Writer as line-delimited [glog](https://github.com/google/glog) formats, Log lines have
      this
      form: [Lyyyymmdd hh:mm:ss.uuuuuu threadid file:line\] msg...](https://github.com/google/glog/blob/v0.6.0/src/glog/logging.h.in#L346).
    - [GlogHumanHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewGlogHumanHandler) provides a
      Handler that writes Records to an io.Writer as line-delimited [glog](https://github.com/google/glog) formats, but
      human-readable, Log lines have this
      form: [\[LLLLL\] \[yyyymmdd hh:mm:ss.uuuuuu\] \[threadid\] \[file:line(func)\] msg...](https://github.com/searKing/golang/blob/go/v1.2.86/go/log/slog/glog_handler.go#L85).
    - [MultiHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#MultiHandler) that duplicates its writes
      to all the provided handlers, similar to the Unix tee(1) command. `MultiHandler` might be useful for write log to
      many Handlers, like writing log to both stdout and rotate file.
    - [RotateHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewRotateHandler) that provides a
      Handler that writes Records to
      a [RotateFile](https://pkg.go.dev/github.com/searKing/golang/go/os). `RotateHandler` might be useful for write log
      to ease administration of systems that generate large numbers of files. It allows automatic rotation,
      removal, and handler of files. Each file may be handled daily, weekly, monthly, strftimely, time_layoutly or when
      it grows too large. Here are some helper
      functions, [NewRotateGlogHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewRotateGlogHandler), [NewRotateGlogHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewRotateGlogHandler), [NewRotateGlogHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewRotateGlogHandler), [NewRotateGlogHandler](https://pkg.go.dev/github.com/searKing/golang/go/log/slog#NewRotateGlogHandler)
* [webhdfs](https://pkg.go.dev/github.com/searKing/webhdfs) - Hadoop WebHDFS REST API client library for Golang with fs
  module like (asynchronous) interface.
* [Thread](https://pkg.go.dev/github.com/searKing/golang/go/sync#Thread) — Thread should be used for such as
  calling OS services or non-Go library functions that depend on per-thread state, as runtime.LockOSThread().
* [BurstLimiter](https://pkg.go.dev/github.com/searKing/golang/go/time/rate#BurstLimiter) — BurstLimiter behaves
  like Limiter in `time`, BurstLimiter controls how frequently events are allowed to happen.
    - It implements a "token bucket" of size b, initially full、empty or any size, and refilled by `PutToken`
      or `PutTokenN`. The difference is
      that `time/rate.Limiter`initially full and refilled at rate r tokens per second.
    - It implements a [Reorder Buffer](https://en.wikipedia.org/wiki/Re-order_buffer) allocated by `Reserve`
      or `ReserveN` into account when allowing future events and `Wait` or `WaitN` blocks until lim permits n events to
      happen.
* [generator](https://pkg.go.dev/github.com/searKing/golang/go/go/generator#Generator) — Generator behaves
  like `Generator` in python or ES6, with yield and next statements.
* [signal](https://pkg.go.dev/github.com/searKing/golang/go/os/signal#Notify) — Signal enhances signal.Notify with the
  stacktrace of cgo.
* [sql](https://pkg.go.dev/github.com/searKing/golang/go/database/sql#NullDuration) — Enhance go std sql.
    - NullDuration
        - ```NullDuration represents an interface that may be null. NullDuration implements the Scanner interface so it can be used as a scan destination, similar to sql.NullString.```
    - NullJson
        - ```NullJson represents an interface that may be null. NullJson implements the Scanner interface so it can be used as a scan destination, similar to sql.NullString. Deprecate, use go-nulljson instead. For more information, see: https://pkg.go.dev/github.com/searKing/golang/tools/go-nulljson```
* [ternary_search_tree](https://pkg.go.dev/github.com/searKing/golang/go/container/trie_tree/ternary_search_tree#TernarySearchTree)
  — A type of trie (sometimes called a prefix tree) where nodes are arranged in a manner similar to a binary search
  tree, but with up to three children rather than the binary tree's limit of two.
* [mux](https://pkg.go.dev/github.com/searKing/golang/go/net/mux) — Mux is a generic Go library to multiplex
  connections based on their payload. Using mux, you can serve gRPC, SSH, HTTPS, HTTP, Go RPC, and pretty much any other
  protocol on the same TCP listener.
* [SniffReader](https://pkg.go.dev/github.com/searKing/golang/go/io#SniffReader) — A Reader that allows sniff
  and read from the provided input reader. data is buffered if Sniff(true) is called. buffered data is taken first, if
  Sniff(false) is called.
* [multiple_prefix](https://pkg.go.dev/github.com/searKing/golang/go/format/multiple_prefix) - Prefixes for
  decimal and binary multiples, [Prefixes for decimal multiples](https://physics.nist.gov/cuu/Units/prefixes.html)
  , [Prefixes for binary multiples](https://physics.nist.gov/cuu/Units/binary.html).
* [flag](https://pkg.go.dev/github.com/searKing/golang/go/flag) — Enhance go std flag, such as StringSlice that
  is a flag.Value that accumulates strings, e.g. --flag=one --flag=two would produce []string{"one", "two"}.
* [goroutine](https://pkg.go.dev/github.com/searKing/golang/go/runtime/goroutine) — goroutine is a collection of
  apis about goroutine.
    - ID() — returns goroutine id of the goroutine that calls it.
    - Lock — represents a goroutine ID, with goroutine ID checked, that is whether GoRoutines of lock newer and check
      caller differ.
* [hashring](https://pkg.go.dev/github.com/searKing/golang/go/container/hashring) — hashring provides a
  consistent hashing function, read more about consistent hashing on
  wikipedia:  [Consistent_hashing](http://en.wikipedia.org/wiki/Consistent_hashing).
* [RotateFile](https://pkg.go.dev/github.com/searKing/golang/go/os#RotateFile) — RotateFile derived from os.File, and is
  designed to ease administration of systems that generate large numbers of files. It allows automatic rotation,
  removal, and handler of files. Each file may be handled daily, weekly, monthly, strftimely, time_layoutly or when it
  grows too large.
    - [NewFactoryFromFile](https://pkg.go.dev/github.com/searKing/golang/third_party/github.com/sirupsen/logrus#NewFactoryFromFile) —
      NewFactoryFromFile is an example of os.RotateFile register for logrus.
* [CacheFile](https://pkg.go.dev/github.com/searKing/golang/go/os#CacheFile) - CacheFile is a package cache(Eventual
  consistency, behaves like sync.LRU[string, *os.File]), backed by a file system directory tree. It is safe for multiple
  processes on a single machine to use the
  same cache directory in a local file system simultaneously. They will coordinate using operating system file locks and
  may duplicate effort but will not corrupt the cache. It's usually used to support download cache, download if cache
  file not hit.
* [UnlinkOldestFiles](https://pkg.go.dev/github.com/searKing/golang/go/os#UnlinkOldestFilesFunc) - UnlinkOldestFiles
  unlink old files if [DiskQuota](https://pkg.go.dev/github.com/searKing/golang/go/os#DiskQuota) exceeds. It's usually
  used for disk clean.

# GoGenerateTools

* [go generate](https://blog.golang.org/generate) is only useful if you have tools to use it with! Here is an incomplete
  list of useful tools that generate code.

* [go-syncmap](https://pkg.go.dev/github.com/searKing/golang/tools/go-syncmap) — Generates Go code using a
  package as a generic template for sync.Map.
* [go-syncpool](https://pkg.go.dev/github.com/searKing/golang/tools/go-syncpool) — Generates Go code using a
  package as a generic template for sync.Pool.
* [go-atomicvalue](https://pkg.go.dev/github.com/searKing/golang/tools/go-atomicvalue) — Generates Go code using
  a package as a generic template for atomic.Value.
* [go-option](https://pkg.go.dev/github.com/searKing/golang/tools/go-option) — Generates Go code using a package
  as a graceful option.
* [go-nulljson](https://pkg.go.dev/github.com/searKing/golang/tools/go-nulljson) — Generates Go code using a
  package as a generic template that implements sql.Scanner and sql.Valuer.
* [go-enum](https://pkg.go.dev/github.com/searKing/golang/tools/go-enum) — Generates Go code using a package as
  a generic template, which implements interface fmt.Stringer | binary | json | text | sql | yaml for enums.
* [go-import](https://pkg.go.dev/github.com/searKing/golang/tools/go-import) — Performs auto import of non go
  files.
* [go-sqlx](https://pkg.go.dev/github.com/searKing/golang/tools/go-sqlx) — Generates Go code using a package as
  a generic template that implements sqlx.
                                                                               