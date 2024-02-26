module github.com/searKing/golang/third_party/google.golang.org/grpc

go 1.21

require (
	github.com/searKing/golang/go v1.2.115
	golang.org/x/net v0.21.0
	golang.org/x/time v0.5.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240221002015-b0ce06bbee7c
	google.golang.org/grpc v1.62.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace github.com/searKing/golang/go => ../../../go
