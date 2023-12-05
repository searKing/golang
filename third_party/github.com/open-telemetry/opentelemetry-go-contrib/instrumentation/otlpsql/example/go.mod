module github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/otlpsql/example

go 1.21

require (
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/otlpsql v1.2.108
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
)

require (
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/contrib v1.21.1 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
)

replace github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/otlpsql => ../../otlpsql
