module github.com/searKing/golang/tools/protoc-gen-go-tag

go 1.27.0

require (
	github.com/searKing/golang/go v1.2.147
	google.golang.org/protobuf v1.36.12
)

// v1.2.148 was published with incorrect commit author metadata.
retract v1.2.148
