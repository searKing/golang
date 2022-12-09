[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=go-syncpool)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/tools/go-syncpool?status.svg)](https://godoc.org/github.com/searKing/golang/tools/go-syncpool)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/go-syncpool)](https://goreportcard.com/report/github.com/searKing/golang/tools/go-syncpool)
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@go-syncpool?badge)

# go-syncpool

Generates Go code using a package as a generic template for sync.Pool.

go-syncpool Generates Go code using a package as a generic template for sync.Pool. Given the name of a sync.Pool type T
, and the name of a type Value go-syncpool will create a new self-contained Go source file implementing

```
// type T sync.Pool
// T<Value>
func (m *T) Put(value Value)
func (m *T) Get() Value
```

The file is created in the same package and directory as the package that defines T, Key. It has helpful defaults
designed for use with go generate.

For example, given this snippet,

```go
package painkiller

import "sync"

type Pill sync.Pool
```

running this command

```
go-syncpool -type "Pill<time.Time>"
```

in the same directory will create the file pill_syncpool.go, in package painkiller, containing a definition of

```
func (m *Pill) Put(value time.Time)
func (m *Pill) Get() time.Time
```

Typically this process would be run using go generate, like this:

```
//go:generate go-syncpool -type "Pill<int>"
//go:generate go-syncpool -type "Pill<string>"
//go:generate go-syncpool -type "Pill<time.Time>"
//go:generate go-syncpool -type "Pill<encoding/json.Token>"
```

If multiple constants have the same value, the lexically first matching name will be used (in the example, Acetaminophen
will print as "Paracetamol").

With no arguments, it processes the package in the current directory. Otherwise, the arguments must name a single
directory holding a Go package or a set of Go source files that represent a single Go package.

The -type flag accepts a comma-separated list of types so a single run can generate methods for multiple types. The
default output file is t_syncpool.go, where t is the lower-cased name of the first type listed. It can be overridden
with the -output flag.

## Download/Install

The easiest way to install is to run `go get install github.com/searKing/golang/tools/go-syncpool`
. You can also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/go-syncpool`.

## Inspiring projects

* [stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)

