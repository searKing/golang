module github.com/searKing/golang/third_party/google.golang.org/grpc

go 1.21

require (
	github.com/searKing/golang/go v1.2.116
	golang.org/x/net v0.25.0
	golang.org/x/time v0.5.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240509183442-62759503f434
	google.golang.org/grpc v1.63.2
)

require (
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace github.com/searKing/golang/go => ../../../go
