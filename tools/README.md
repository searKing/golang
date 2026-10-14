# Go Tools

This subrepository holds the source for various packages and tools that support
the Go programming language.

All of the tools can be fetched with `go get`.

Packages include an implementation of the
Static Single Assignment form (SSA) representation for Go programs.

## Download/Install

Each tool is a module of its own. Since Go 1.24, track tools by tool directives of go.mod, and run them by `go tool`:

```bash
go get -tool github.com/searKing/golang/tools/go-option@latest
go tool go-option -type=Pill
```

or install them into GOBIN:

```bash
go install github.com/searKing/golang/tools/go-atomicvalue@latest
go install github.com/searKing/golang/tools/go-enum@latest
go install github.com/searKing/golang/tools/go-import@latest
go install github.com/searKing/golang/tools/go-nulljson@latest
go install github.com/searKing/golang/tools/go-option@latest
go install github.com/searKing/golang/tools/go-sqlx@latest
go install github.com/searKing/golang/tools/go-syncmap@latest
go install github.com/searKing/golang/tools/go-syncpool@latest
go install github.com/searKing/golang/tools/go-union@latest
go install github.com/searKing/golang/tools/go-validator@latest
go install github.com/searKing/golang/tools/protoc-gen-go-tag@latest
```

In this repository, tools are tracked by internal/tools/go.mod, with tools of this repository replaced by local paths,
and `make generate` runs `go generate` in every module with them.

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The main issue tracker for the tools repository is located at
https://github.com/searKing/golang/issues. Prefix your issue with "golang/tools/(your
subdir):" in the subject line, so it is easy to find.