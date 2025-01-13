module github.com/searKing/golang/third_party/github.com/golang/go

go 1.23

toolchain go1.23.3

require (
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/searKing/golang/go v1.2.120
)

require (
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
)

replace github.com/searKing/golang/go => ../../../../go
