[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=github.com/golang/go/error)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/go/error?status.svg)](https://godoc.org/github.com/searKing/golang/go/error)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/go/error)](https://goreportcard.com/report/github.com/searKing/golang/go/error) 
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@github.com/golang/go/error?badge)
# cause 
Package errors provides simple error handling primitives.

`go get github.com/searKing/golang/go/error`

The traditional error handling idiom in Go is roughly akin to
```go
if err != nil {
        return err
}
```
which applied recursively up the call stack results in error reports without context or debugging information. 
The errors package allows programmers to use errors as Exception in Java in their code.


[Read the package documentation for more information](https://godoc.org/github.com/searKing/golang/go/error).

