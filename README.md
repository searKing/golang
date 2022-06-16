[![GoDoc](https://pkg.go.dev/github.com/searKing/golang?status.svg)](https://pkg.go.dev/github.com/searKing/golang)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang)](https://goreportcard.com/report/github.com/searKing/golang)

# golang

Useful libs or tools for Golang

# GoLibs

* [webhdfs](https://github.com/searKing/webhdfs) - Hadoop WebHDFS REST API client library for Golang with fs module
  like (asynchronous) interface.
* [Thread](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/sync#Thread) — Thread should be used for such as
  calling OS services or non-Go library functions that depend on per-thread state, as runtime.LockOSThread().
* [BurstLimiter](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/time/rate#BurstLimiter) — BurstLimiter behaves
  like Limiter in `time`, BurstLimiter controls how frequently events are allowed to happen. It implements a "token
  bucket" of size b, initially full、empty or any size, and refilled by `PutToken` or `PutTokenN`. The difference is
  that `time/rate.Limiter`initially full and refilled at rate r tokens per second.
* [generator](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/go/generator#Generator) — Generator behaves like
  Generator in python or ES6, with yield and next statements.
* [signal](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/os/signal#Notify) — Signal enhances signal.Notify
  with the stacktrace of cgo.
* [sql](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/database/sql#NullDuration) — Enhance go std sql.
    - NullDuration
        - ```NullDuration represents an interface that may be null. NullDuration implements the Scanner interface so it can be used as a scan destination, similar to sql.NullString.```
    - NullJson
        - ```NullJson represents an interface that may be null. NullJson implements the Scanner interface so it can be used as a scan destination, similar to sql.NullString. Deprecate, use go-nulljson instead. For more information, see: https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-nulljson```
* [ternary_search_tree](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/container/trie_tree/ternary_search_tree#TernarySearchTree)
  — A type of trie (sometimes called a prefix tree) where nodes are arranged in a manner similar to a binary search
  tree, but with up to three children rather than the binary tree's limit of two.
* [mux](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/net/mux) — Mux is a generic Go library to multiplex
  connections based on their payload. Using mux, you can serve gRPC, SSH, HTTPS, HTTP, Go RPC, and pretty much any other
  protocol on the same TCP listener.
* [SniffReader](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/io#SniffReader) — A Reader that allows sniff
  and read from the provided input reader. data is buffered if Sniff(true) is called. buffered data is taken first, if
  Sniff(false) is called.
* [multiple_prefix](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/format/multiple_prefix) - Prefixes for
  decimal and binary multiples, [Prefixes for decimal multiples](https://physics.nist.gov/cuu/Units/prefixes.html)
  , [Prefixes for binary multiples](https://physics.nist.gov/cuu/Units/binary.html).
* [flag](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/flag) — Enhance go std flag, such as StringSlice that
  is a flag.Value that accumulates strings, e.g. --flag=one --flag=two would produce []string{"one", "two"}.
* [goroutine](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/runtime/goroutine) — goroutine is a collection of
  apis about goroutine.
    - ID() — returns goroutine id of the goroutine that calls it.
    - Lock — represents a goroutine ID, with goroutine ID checked, that is whether GoRoutines of lock newer and check
      caller differ.
* [hashring](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/container/hashring) — hashring provides a
  consistent hashing function, read more about consistent hashing on
  wikipedia:  [Consistent_hashing](http://en.wikipedia.org/wiki/Consistent_hashing).
* [RotateFile](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/go/os) — RotateFile derived from os.File, and is
  designed to ease administration of systems that generate large numbers of files. It allows automatic rotation,
  removal, and handler of files. Each file may be handled daily, weekly, monthly, strftimely, time_layoutly or when it
  grows too large.
    - [WithRotate]https://pkg.go.dev/github.com/searKing/golang@v1.1.18/third_party/github.com/sirupsen/logrus) -
      WithRotate is an example of os.RotateFile register for logrus.

# GoGenerateTools

* [go generate](https://blog.golang.org/generate) is only useful if you have tools to use it with! Here is an incomplete
  list of useful tools that generate code.

* [go-syncmap](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-syncmap) — Generates Go code using a
  package as a generic template for sync.Map.
* [go-syncpool](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-syncpool) — Generates Go code using a
  package as a generic template for sync.Pool.
* [go-atomicvalue](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-atomicvalue) — Generates Go code using
  a package as a generic template for atomic.Value.
* [go-option](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-option) — Generates Go code using a package
  as a graceful option.
* [go-nulljson](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-nulljson) — Generates Go code using a
  package as a generic template that implements sql.Scanner and sql.Valuer.
* [go-enum](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-enum) — Generates Go code using a package as
  a generic template, which implements interface fmt.Stringer | binary | json | text | sql | yaml for enums.
* [go-import](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-import) — Performs auto import of non go
  files.
* [go-sqlx](https://pkg.go.dev/github.com/searKing/golang@v1.1.18/tools/go-sqlx) — Generates Go code using a package as
  a generic template that implements sqlx.
                                                                               