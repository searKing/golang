module github.com/searKing/golang/tools/cmd/go-sqlx

go 1.21

toolchain go1.21.5

require (
	github.com/searKing/golang/go v1.2.112
	github.com/searKing/golang/tools v1.2.112
	golang.org/x/tools v0.16.0
)

require golang.org/x/mod v0.14.0 // indirect

replace github.com/searKing/golang/go => ../../../go

replace github.com/searKing/golang/tools => ../../../tools
