[![GoDoc](https://godoc.org/github.com/searKing/golang?status.svg)](https://godoc.org/github.com/searKing/golang)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang)](https://goreportcard.com/report/github.com/searKing/golang) 
# golang
Useful libs or tools for Golang

# GoLibs
* [generator](https://godoc.org/github.com/searKing/golang/go/go/generator) - Generator behaves like Generator in python or ES6, with yield and next statements.
* [ternary_search_tree](https://godoc.org/github.com/searKing/golang/go/container/trie_tree/ternary_search_tree) - A type of trie (sometimes called a prefix tree) where nodes are arranged in a manner similar to a binary search tree, but with up to three children rather than the binary tree's limit of two.
* [connection mux](https://godoc.org/github.com/searKing/golang/go/net/cmux) - Connection Mux is a generic Go library to multiplex connections based on their payload. Using cmux, you can serve gRPC, SSH, HTTPS, HTTP, Go RPC, and pretty much any other protocol on the same TCP listener.
* [SniffReader](https://godoc.org/github.com/searKing/golang/go/io) - A Reader that allows sniff and read from the provided input reader. data is buffered if Sniff(true) is called. buffered data is taken first, if Sniff(false) is called.
* [multiple_prefix](https://godoc.org/github.com/searKing/golang/go/format/multiple_prefix) - Prefixes for decimal and binary multiples, [Prefixes for decimal multiples](https://physics.nist.gov/cuu/Units/prefixes.html), [Prefixes for binary multiples](https://physics.nist.gov/cuu/Units/binary.html).

# GoGenerateTools
[`go generate`](https://blog.golang.org/generate) is only useful if you have tools to use it with! Here is an incomplete list of useful tools that generate code.

* [go-syncmap](https://godoc.org/github.com/searKing/golang/tools/cmd/go-syncmap) - Generates Go code using a package as a generic template for sync.Map.
* [go-syncpool](https://godoc.org/github.com/searKing/golang/tools/cmd/go-syncpool) - Generates Go code using a package as a generic template for sync.Pool.
* [go-atomicvalue](https://godoc.org/github.com/searKing/golang/tools/cmd/go-atomicvalue) - Generates Go code using a package as a generic template for atomic.Value.
* [go-option](https://godoc.org/github.com/searKing/golang/tools/cmd/go-option) - Generates Go code using a package as a graceful option.
* [go-nulljson](https://godoc.org/github.com/searKing/golang/tools/cmd/go-nulljson) - Generates Go code using a package as a generic template that implements sql.Scanner and sql.Valuer.
* [go-enum](https://godoc.org/github.com/searKing/golang/tools/cmd/go-enum) - Generates Go code using a package as a generic template which implements interface fmt.Stringer | binary | json | text | sql | yaml for enums.