module github.com/searKing/golang/pkg/webserver

go 1.21

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0
	github.com/rs/cors v1.11.0
	github.com/searKing/golang/go v1.2.118
	github.com/searKing/golang/third_party/github.com/gin-gonic/gin v1.2.118
	github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2 v1.2.118
	github.com/searKing/golang/third_party/google.golang.org/grpc v1.2.118
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842
	golang.org/x/sync v0.7.0
	google.golang.org/grpc v1.64.1
)

require (
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240513163218-0867130af1f8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/searKing/golang/go => ../../go

replace github.com/searKing/golang/third_party/github.com/gin-gonic/gin => ../../third_party/github.com/gin-gonic/gin

replace github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2 => ../../third_party/github.com/grpc-ecosystem/grpc-gateway-v2

replace github.com/searKing/golang/third_party/google.golang.org/grpc => ../../third_party/google.golang.org/grpc
