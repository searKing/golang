module github.com/searKing/golang/tools/cmd/go-import

go 1.21

toolchain go1.21.5

require (
	github.com/searKing/golang/go v1.2.112
	github.com/searKing/golang/tools/go-import v1.2.112
)

require (
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
)

replace github.com/searKing/golang/go => ../../../go

replace github.com/searKing/golang/tools/go-import => ../../../tools/go-import
