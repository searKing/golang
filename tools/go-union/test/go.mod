module github.com/searKing/golang/tools/go-union/test

go 1.27.0

replace github.com/searKing/golang/tools/go-union => ../

replace github.com/searKing/golang/tools => ../..

replace github.com/searKing/golang/go => ../../../go

tool github.com/searKing/golang/tools/go-union

require (
	github.com/searKing/golang/go v1.2.149 // indirect
	github.com/searKing/golang/tools v1.2.149 // indirect
	github.com/searKing/golang/tools/go-union v1.2.149 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/tools v0.48.0 // indirect
)
