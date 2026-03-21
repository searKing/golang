// Deprecated: Use the "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc" module instead.
// Deprecated: Use the "google.golang.org/grpc/stats/opentelemetry" module instead.
module github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc

go 1.25.0

require (
	github.com/searKing/golang/go v1.2.120
	go.opentelemetry.io/contrib v1.26.0
	go.opentelemetry.io/otel v1.42.0
	go.opentelemetry.io/otel/metric v1.42.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/trace v1.42.0 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240509183442-62759503f434 // indirect
)

replace github.com/searKing/golang/go => ../../../../../../../../go
