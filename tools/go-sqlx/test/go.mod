module github.com/searKing/golang/tools/go-sqlx/test

go 1.27.0

replace github.com/searKing/golang/tools/go-sqlx => ../

replace github.com/searKing/golang/tools => ../..

replace github.com/searKing/golang/go => ../../../go

tool github.com/searKing/golang/tools/go-sqlx

require (
	github.com/jmoiron/sqlx v1.4.0
	github.com/searKing/golang/go v1.2.149
	github.com/searKing/golang/third_party/github.com/jmoiron/sqlx v1.2.149
)

require (
	github.com/searKing/golang/tools v1.2.149 // indirect
	github.com/searKing/golang/tools/go-sqlx v1.2.149 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/tools v0.48.0 // indirect
)

replace github.com/searKing/golang/third_party/github.com/jmoiron/sqlx => ../../../third_party/github.com/jmoiron/sqlx
