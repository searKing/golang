module github.com/searKing/golang/pkg/webserver

go 1.27.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/go-playground/validator/v10 v10.30.5
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.4
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.30.0
	github.com/rs/cors v1.11.1
	github.com/searKing/golang/go v1.2.145
	github.com/searKing/golang/third_party/github.com/gin-gonic/gin v1.2.129
	github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2 v1.2.136
	github.com/searKing/golang/third_party/google.golang.org/grpc v1.2.129
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.71.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.71.0
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/sdk v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
	golang.org/x/exp v0.0.0-20260908205506-85c1c2202aba
	golang.org/x/sync v0.23.0
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/bytedance/gopkg v0.1.4 // indirect
	github.com/bytedance/sonic v1.15.4 // indirect
	github.com/bytedance/sonic/loader v0.5.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.7 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.15 // indirect
	github.com/gin-contrib/sse v1.1.2 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.2 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.4.0 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.63.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.2 // indirect
	go.mongodb.org/mongo-driver/v2 v2.9.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.46.0 // indirect
	golang.org/x/arch v0.31.0 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	golang.org/x/time v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260921155816-b14227669459 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/searKing/golang/go => ../../go

replace github.com/searKing/golang/third_party/github.com/gin-gonic/gin => ../../third_party/github.com/gin-gonic/gin

replace github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2 => ../../third_party/github.com/grpc-ecosystem/grpc-gateway-v2

replace github.com/searKing/golang/third_party/google.golang.org/grpc => ../../third_party/google.golang.org/grpc

replace github.com/searKing/golang/pkg/instrumentation/otel => ../../pkg/instrumentation/otel
