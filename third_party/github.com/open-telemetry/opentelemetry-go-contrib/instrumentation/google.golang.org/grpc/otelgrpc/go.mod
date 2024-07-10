// Deprecated: Use the "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc" module instead.
// Deprecated: Use the "google.golang.org/grpc/stats/opentelemetry" module instead.
module github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc

go 1.21

require (
	github.com/searKing/golang/go v1.2.118
	go.opentelemetry.io/contrib v1.26.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.51.0
	go.opentelemetry.io/otel v1.26.0
	go.opentelemetry.io/otel/metric v1.26.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240509183442-62759503f434 // indirect
)

replace github.com/searKing/golang/go => ../../../../../../../../go
