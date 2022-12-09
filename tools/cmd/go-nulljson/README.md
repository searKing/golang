[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=go-nulljson)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/tools/go-nulljson?status.svg)](https://godoc.org/github.com/searKing/golang/tools/go-nulljson)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/go-nulljson)](https://goreportcard.com/report/github.com/searKing/golang/tools/go-nulljson)
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@go-nulljson?badge)

# go-nulljson

Generates Go code using a package as a generic template that implements database/sql.Scanner and
database/sql/driver.Valuer.

go-nulljson Generates Go code using a package as a generic template that implements database/sql.Scanner and
database/sql/driver.Valuer. Given the name of a NullJson type T , and the name of a type Value go-nulljson will create a
new self-contained Go source file implementing

```
func (m *T) Scan(src interface{}) error
func (m *T) Value() (driver.Value, error)
```

The file is created in the same package and directory as the package that defines T, Key. It has helpful defaults
designed for use with go generate.

For example, given this snippet,

```go
package painkiller


```

running this command

```
go-nulljson -type="Pill<time.Time>"
```

in the same directory will create the file pill_nulljson.go, in package painkiller, containing a definition of

```
func (m *Pill) Scan(src interface{}) error
func (m *Pill) Value() (driver.Value, error)
```

Typically this process would be run using go generate, like this:

```
//go:generate go-nulljson -type "Pill<int>"
//go:generate go-nulljson -type "Pill<*string>"
//go:generate go-nulljson -type "Pill<time.Time>"
//go:generate go-nulljson -type "Pill<*encoding/json.Token>"
```

If multiple constants have the same value, the lexically first matching name will be used (in the example, Acetaminophen
will print as "Paracetamol").

With no arguments, it processes the package in the current directory. Otherwise, the arguments must name a single
directory holding a Go package or a set of Go source files that represent a single Go package.

The -type flag accepts a comma-separated list of types so a single run can generate methods for multiple types. The
default output file is t_nulljson.go, where t is the lower-cased name of the first type listed. It can be overridden
with the -output flag.

## Download/Install

The easiest way to install is to run `go get install github.com/searKing/golang/tools/go-nulljson`
. You can also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/go-nulljson`.

## Inspiring projects

* [stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
