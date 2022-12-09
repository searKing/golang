[![Build Status](https:travis-ci.org/searKing/travis-ci.svg?branch=go-enum)](https:travis-ci.org/searKing/travis-ci)
[![GoDoc](https:godoc.org/github.com/searKing/golang/tools/go-enum?status.svg)](https:godoc.org/github.com/searKing/golang/tools/go-enum)
[![Report card](https:goreportcard.com/badge/github.com/searKing/golang/tools/go-enum)](https:goreportcard.com/report/github.com/searKing/golang/tools/go-enum)
[![Sourcegraph](https:sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https:sourcegraph.com/github.com/searKing/travis-ci@go-enum?badge)

# go-enum

Generates Go code using a package as a generic template which implements interface fmt.Stringer | binary | json | text |
sql | yaml for enums.

go-enum is a tool to automate the creation of methods that satisfy such interfaces:

```text
	fmt         ==>  fmt.Stringer
	binary      ==>  encoding.BinaryMarshaler and encoding.BinaryUnmarshaler
	json        ==>  encoding/json.MarshalJSON and encoding/json.UnmarshalJSON
	text        ==>  encoding.TextMarshaler and encoding.TextUnmarshaler
	sql         ==>  database/sql.Scanner and database/sql/driver.Valuer
	yaml        ==>  gopkg.in/yaml.v2:yaml.Marshaler and gopkg.in/yaml.v2:yaml.Unmarshaler
```

Given the name of a (signed or unsigned) integer type T that has constants defined, stringer will create a new
self-contained Go source file implementing

```text
	fmt         ==>  fmt.Stringer
		func (t T) String() string
	binary      ==>  encoding.BinaryMarshaler and encoding.BinaryUnmarshaler
		func (t T) MarshalBinary() (data []byte, err error)
		func (t *T) UnmarshalBinary(data []byte) error
	json        ==>  encoding/json.MarshalJSON and encoding/json.UnmarshalJSON
		func (t T) MarshalJSON() ([]byte, error)
		func (t *T) UnmarshalJSON(data []byte) error
	text        ==>  encoding.TextMarshaler and encoding.TextUnmarshaler
		func (t T) MarshalText() ([]byte, error)
		func (t *T) UnmarshalText(text []byte) error
	sql         ==>  database/sql.Scanner and database/sql/driver.Valuer
		func (t T) Value() (driver.Value, error)
		func (t *T) Scan(value interface{}) error
	yaml        ==>  gopkg.in/yaml.v2:yaml.Marshaler and gopkg.in/yaml.v2:yaml.Unmarshaler
		func (t T) MarshalYAML() (interface{}, error)
		func (t *T) UnmarshalYAML(unmarshal func(interface{}) error) error
```

The file is created in the same package and directory as the package that defines T. It has helpful defaults designed
for use with go generate.

go-enum works best with constants that are consecutive values such as created using iota, but creates good code
regardless. In the future it might also provide custom support for constant sets that are bit patterns.

For example, given this snippet,

```go
    package painkiller

type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)
```

running this command

```bash
	go-enum -type=Pill
```

in the same directory will create the file pill_string.go, in package painkiller, containing a definition of interfaces
mentioned.

That method will translate the value of a Pill constant to the string representation of the respective constant name, so
that the call fmt.Print(painkiller.Aspirin) will print the string "Aspirin".

Typically this process would be run using go generate, like this:

```bash
	go:generate go-enum -type=Pill
```

If multiple constants have the same value, the lexically first matching name will be used (in the example, Acetaminophen
will print as "Paracetamol").

With no arguments, it processes the package in the current directory. Otherwise, the arguments must name a single
directory holding a Go package or a set of Go source files that represent a single Go package.

The -type flag accepts a comma-separated list of types so a single run can generate methods for multiple types. The
default output file is t_string.go, where t is the lower-cased name of the first type listed. It can be overridden with
the -output flag.

The -linecomment flag tells stringer to generate the text of any line comment, trimmed of leading spaces, instead of the
constant name. For instance, if the constants above had a Pill prefix, one could write PillAspirin Aspirin to suppress
it in the output.

## Download/Install

The easiest way to install is to run `go get install github.com/searKing/golang/tools/go-enum`
. You can also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/go-enum`.

## Inspiring projects

* [stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
* [jsonenums](https://github.com/campoy/jsonenums)
* [enumer](https://github.com/alvaroloes/enumer)