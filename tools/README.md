# Go Tools

This subrepository holds the source for various packages and tools that support
the Go programming language.

All of the tools can be fetched with `go get`.

Packages include an implementation of the
Static Single Assignment form (SSA) representation for Go programs.

## Download/Install

The easiest way to install is to run `go get install github.com/searKing/golang/tools/...`. You can
also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/`.

### Tips
```bash
go get install github.com/searKing/golang/tools/go-atomicvalue
go get install github.com/searKing/golang/tools/go-enum
go get install github.com/searKing/golang/tools/go-import
go get install github.com/searKing/golang/tools/go-nulljson
go get install github.com/searKing/golang/tools/go-option
go get install github.com/searKing/golang/tools/go-sqlx
go get install github.com/searKing/golang/tools/go-syncmap
go get install github.com/searKing/golang/tools/go-syncpool
go get install github.com/searKing/golang/tools/go-validator
go get install github.com/searKing/golang/tools/protoc-gen-go-tag
```

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The main issue tracker for the tools repository is located at
https://github.com/searKing/golang/issues. Prefix your issue with "golang/tools/(your
subdir):" in the subject line, so it is easy to find.