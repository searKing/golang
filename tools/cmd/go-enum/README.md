[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=go-jsonenum)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/tools/cmd/go-jsonenum?status.svg)](https://godoc.org/github.com/searKing/golang/tools/cmd/go-jsonenum)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/cmd/go-jsonenum)](https://goreportcard.com/report/github.com/searKing/golang/tools/cmd/go-jsonenum) 
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@go-jsonenum?badge)
# go-jsonenum
Generates Go code using a package as a generic template for atomic.Value.

go-jsonenum Generates Go code using a package as a generic template for atomic.Value.
Given the name of a atomic.Value type T , and the name of a type Value
go-jsonenum will create a new self-contained Go source file implementing
```
// type T atomic.Value
// T<Value>
func (m *T) Store(value Value)
func (m *T) Load() Value
```

The file is created in the same package and directory as the package that defines T, Key.
It has helpful defaults designed for use with go generate.

For example, given this snippet,

```go
package painkiller

import "sync/atomic"

type Pill atomic.Value
```

running this command
```
go-jsonenum -type="Pill<time.Time>"
```

in the same directory will create the file pill_jsonenum.go, in package painkiller,
containing a definition of

```
func (m *Pill) Store(value time.Time)
func (m *Pill) Load() time.Time
```

Typically this process would be run using go generate, like this:
```
//go:generate go-jsonenum -type "Pill<int>"
//go:generate go-jsonenum -type "Pill<*string>"
//go:generate go-jsonenum -type "Pill<time.Time>"
//go:generate go-jsonenum -type "Pill<*encoding/json.Token>"
```

If multiple constants have the same value, the lexically first matching name will
be used (in the example, Acetaminophen will print as "Paracetamol").

With no arguments, it processes the package in the current directory.
Otherwise, the arguments must name a single directory holding a Go package
or a set of Go source files that represent a single Go package.

The -type flag accepts a comma-separated list of types so a single run can
generate methods for multiple types. The default output file is t_jsonenum.go,
where t is the lower-cased name of the first type listed. It can be overridden
with the -output flag.

## Download/Install

The easiest way to install is to run `go get -u github.com/searKing/golang/tools/cmd/go-jsonenum`. You can
also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/cmd/go-jsonenum`.

