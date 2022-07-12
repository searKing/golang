[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=go-validator)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/tools/go-validator?status.svg)](https://godoc.org/github.com/searKing/golang/tools/go-validator)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/go-validator)](https://goreportcard.com/report/github.com/searKing/golang/tools/go-validator)
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@go-validator?badge)

# go-validator

Generates Go code using a package as a generic template that implements validator.

go-validator Generates Go code using a package as a generic template that implements validator. Given the StructName of
a Struct type T go-validator will create a new self-contained Go source file and rewrite T's "db" tag of struct field

The file is created in the same package and directory as the package that defines T. It has helpful defaults designed
for use with go generate.

For example, given this snippet,

running this command

```go
package painkiller

import (
	"database/sql"
	"time"
)

type Pill struct {
	Id        uint      `db:"id" json:"sql_data_id,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"sql_data_created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"sql_data_updated_at,omitempty"`

	IsDeleted bool         `json:"sql_data_is_deleted,omitempty" db:"is_deleted"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"sql_data_deleted_at,omitempty"`

	Version uint `db:"version" json:"sql_data_version,omitempty"`
} // sql_data

```

running this command

```
go-validator -type=Pill
```

in the same directory will create the file pill_validator.go, in package painkiller, containing a definition of

```
func (m *Pill) Validate(validate *validator.Validate) error
```

Typically this process would be run using go generate, like this:

```
//go:generate go-validator --all-type
//go:generate go-validator -type "Pill"
//go:generate go-validator -type "Pill" --linecomment
//go:generate go-validator -type "Pill" --linecomment --with-dao
```

If multiple constants have the same value, the lexically first matching name will be used (in the example, Acetaminophen
will print as "Paracetamol").

With no arguments, it processes the package in the current directory. Otherwise, the arguments must name a single
directory holding a Go package or a set of Go source files that represent a single Go package.

The -type flag accepts a comma-separated list of types so a single run can generate methods for multiple types. The
default output file is t_validator.go, where t is the lower-cased name of the first type listed. It can be overridden
with the -output flag.

## Download/Install

The easiest way to install is to run `go get -u github.com/searKing/golang/tools/go-validator`
. You can also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/go-validator`.

## Inspiring projects

* [stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
