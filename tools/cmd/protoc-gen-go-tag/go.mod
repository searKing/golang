module github.com/searKing/golang/tools/cmd/protoc-gen-go-tag

go 1.21

toolchain go1.21.5

require (
	github.com/searKing/golang/tools/protoc-gen-go-tag v1.2.112
	google.golang.org/protobuf v1.32.0
)

require github.com/searKing/golang/go v1.2.112 // indirect

replace github.com/searKing/golang/tools/protoc-gen-go-tag => ../../../tools/protoc-gen-go-tag
