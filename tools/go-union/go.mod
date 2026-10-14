module github.com/searKing/golang/tools/go-union

go 1.27.0

require (
	github.com/searKing/golang/go v1.2.147
	github.com/searKing/golang/tools v1.2.147
	golang.org/x/tools v0.48.0
)

require (
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
)

replace github.com/searKing/golang/go => ../../go

replace github.com/searKing/golang/tools => ../
