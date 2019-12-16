[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=go-import)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/tools/cmd/go-import?status.svg)](https://godoc.org/github.com/searKing/golang/tools/cmd/go-import)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/cmd/go-import)](https://goreportcard.com/report/github.com/searKing/golang/tools/cmd/go-import) 
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@go-import?badge)
# go-import
Performs auto import of non go files.

go-import Performs auto import of non go files.
Given the directory to be imported
go-import will create gokeep.go Go source files and a new self-contained goimport.go Go source file.
+ The gokeep.go file is created in the same package and directory as the cwd package.
+ The goimport.go file is created in the package and directory under directories to be imported,
It has helpful defaults designed for use with go generate.

For example, given this snippet,

```go
package painkiller

```

running this command
```bash
go-import /dirs_to_be_force_imported
```

in the same directory will create the file goimport.go,
and in /dirs_to_be_force_imported will create the file gokeep.go

Typically, this process would be run using go generate, like this:
```bash
//go:generate go-import
```

With no arguments, it processes the package in the current directory.
Otherwise, the arguments must name a single directory holding a Go package
or a set of Go source files that represent a single Go package.

 The -tag flag accepts a build tag string.

## Download/Install

The easiest way to install is to run `go get -u github.com/searKing/golang/tools/cmd/go-import`. You can
also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/cmd/go-import`.
