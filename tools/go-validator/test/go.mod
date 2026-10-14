module github.com/searKing/golang/tools/go-validator/test

go 1.27.0

replace github.com/searKing/golang/tools/go-validator => ../

replace github.com/searKing/golang/go => ../../../go

replace github.com/searKing/golang/tools => ../..

tool github.com/searKing/golang/tools/go-validator

require github.com/go-playground/validator/v10 v10.30.5

require (
	github.com/gabriel-vasile/mimetype v1.4.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/searKing/golang/go v1.2.149 // indirect
	github.com/searKing/golang/tools v1.2.149 // indirect
	github.com/searKing/golang/tools/go-validator v1.2.149 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
)
