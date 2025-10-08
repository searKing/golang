module github.com/searKing/golang/third_party/google.golang.org/grpc

go 1.24.0

require (
	github.com/searKing/golang/go v1.2.120
	golang.org/x/net v0.42.0
	golang.org/x/time v0.11.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b
	google.golang.org/grpc v1.76.0
)

require (
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/searKing/golang/go => ../../../go
