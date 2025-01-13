module github.com/searKing/golang/third_party/google.golang.org/grpc/grpclog/logruslogger

go 1.23

toolchain go1.23.3

require (
	github.com/searKing/golang/go v1.2.120
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.59.0
)

require golang.org/x/sys v0.28.0 // indirect

replace github.com/searKing/golang/go => ../../../../../go
