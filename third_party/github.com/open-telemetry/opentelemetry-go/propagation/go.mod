module github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go/propagation

go 1.26.0

require (
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel v1.46.0
	google.golang.org/grpc v1.84.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.opentelemetry.io/otel/trace v1.46.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.48.0 // indirect
)

// v1.2.148 was published with incorrect commit author metadata.
retract v1.2.148
