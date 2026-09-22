module github.com/searKing/golang/third_party/google.golang.org/grpc

go 1.27.0

require (
	github.com/searKing/golang/go v1.2.145
	golang.org/x/net v0.59.0
	golang.org/x/time v0.16.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260921155816-b14227669459
	google.golang.org/grpc v1.84.0
)

require (
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/searKing/golang/go => ../../../go
