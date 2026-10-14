module github.com/searKing/golang/tools/go-validator

go 1.27.0

require (
	github.com/searKing/golang/go v1.2.122
	github.com/searKing/golang/tools v1.2.122
	golang.org/x/tools v0.48.0
)

require (
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
)

replace github.com/searKing/golang/tools => ../

replace github.com/searKing/golang/go => ../../go
