[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=go-sqlx)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/tools/go-sqlx?status.svg)](https://godoc.org/github.com/searKing/golang/tools/go-sqlx)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/go-sqlx)](https://goreportcard.com/report/github.com/searKing/golang/tools/go-sqlx)
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@go-sqlx?badge)

# go-sqlx

Generates Go code using a package as a generic template that implements sqlx.

go-sqlx Generates Go code using a package as a generic template that implements sqlx. Given the StructName of a Struct
type T go-sqlx will create a new self-contained Go source file and rewrite T's "db" tag of struct field

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
go-sqlx -type=Pill --linecomment
```

in the same directory will create the file pill_sqlx.go, in package painkiller, containing a definition of

```
func (m *Pill) Scan(src interface{}) error
func (m *Pill) Value() (driver.Value, error)
```

Typically this process would be run using go generate, like this:

```
//go:generate go-sqlx -type "Pill"
//go:generate go-sqlx -type "Pill" --linecomment
//go:generate go-sqlx -type "Pill" --linecomment --with-dao
```

If multiple constants have the same value, the lexically first matching name will be used (in the example, Acetaminophen
will print as "Paracetamol").

With no arguments, it processes the package in the current directory. Otherwise, the arguments must name a single
directory holding a Go package or a set of Go source files that represent a single Go package.

The -type flag accepts a comma-separated list of types so a single run can generate methods for multiple types. The
default output file is t_sqlx.go, where t is the lower-cased name of the first type listed. It can be overridden with
the -output flag.

## Download/Install

The easiest way to install is to run `go get install github.com/searKing/golang/tools/go-sqlx`
. You can also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/go-sqlx`.

## Inspiring projects

* [stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
