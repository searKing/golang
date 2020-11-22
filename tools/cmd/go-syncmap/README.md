[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=go-syncmap)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/tools/cmd/go-syncmap?status.svg)](https://godoc.org/github.com/searKing/golang/tools/cmd/go-syncmap)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/cmd/go-syncmap)](https://goreportcard.com/report/github.com/searKing/golang/tools/cmd/go-syncmap) 
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@go-syncmap?badge)
# go-syncmap
Generates Go code using a package as a generic template for sync.Map.

go-syncmap Generates Go code using a package as a generic template for sync.Map.
Given the name of a sync.Map type T , and the name of a type Key and Value
go-syncmap will create a new self-contained Go source file implementing
```
// type T sync.Map
// T<Key,Value>

# from Go version 1.9 onward 
func (m *T) Load(key Key) (Value, bool)
func (m *T) Store(key Key, value Value)
func (m *T) LoadOrStore(key Key, value Value) (Value, bool)
func (m *T) Delete(key Key)
func (m *T) Range(f func(key Key, value Value) bool

# from Go version 1.15 onward 
func (m *T) LoadAndDelete(key Key) (Value, bool)
```

The file is created in the same package and directory as the package that defines T, Key and Value.
It has helpful defaults designed for use with go generate.

For example, given this snippet,

```go
package painkiller

import "sync"

type Pill sync.Map
```

running this command
```
go-syncmap -type="Pill<int,string>"
```

in the same directory will create the file pill_syncmap.go, in package painkiller,
containing a definition of

```

# from Go version 1.9 onward 
func (m *Pill) Store(key int, value string)
func (m *Pill) LoadOrStore(key int, value string) (string, bool)
func (m *Pill) Load(key int) (string, bool)
func (m *Pill) Delete(key int)
func (m *Pill) Range(f func(key int, value string) bool)

# from Go version 1.15 onward 
func (m *Pill) LoadAndDelete(key int) (string, bool)
```

Typically this process would be run using go generate, like this:
```
//go:generate go-syncmap -type "Pill<int, string>"
//go:generate go-syncmap -type "Pill<int, time.Time>"
//go:generate go-syncmap -type "Pill<int, encoding/json.Token>"
```

If multiple constants have the same value, the lexically first matching name will
be used (in the example, Acetaminophen will print as "Paracetamol").

With no arguments, it processes the package in the current directory.
Otherwise, the arguments must name a single directory holding a Go package
or a set of Go source files that represent a single Go package.

The -type flag accepts a comma-separated list of types so a single run can
generate methods for multiple types. The default output file is t_syncmap.go,
where t is the lower-cased name of the first type listed. It can be overridden
with the -output flag.

## Download/Install

The easiest way to install is to run `go get -u github.com/searKing/golang/tools/cmd/go-syncmap`. You can
also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/cmd/go-syncmap`.

## Inspiring projects
* [stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
