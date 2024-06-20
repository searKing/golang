module github.com/searKing/golang/third_party/google.golang.org/grpc

go 1.21

require (
	github.com/searKing/golang/go v1.2.118
	golang.org/x/net v0.26.0
	golang.org/x/time v0.5.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240617180043-68d350f18fd4
	google.golang.org/grpc v1.64.0
)

require (
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/searKing/golang/go => ../../../go
