module github.com/searKing/golang/pkg/webserver

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.1
	github.com/rs/cors v1.10.1
	github.com/searKing/golang/go v1.2.115
	github.com/searKing/golang/third_party/github.com/gin-gonic/gin v1.2.115
	github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2 v1.2.115
	github.com/searKing/golang/third_party/google.golang.org/grpc v1.2.115
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225
	golang.org/x/sync v0.6.0
	google.golang.org/grpc v1.62.0
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240125205218-1f4bbc51befe // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240221002015-b0ce06bbee7c // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/searKing/golang/go => ../../go

replace github.com/searKing/golang/third_party/github.com/gin-gonic/gin => ../../third_party/github.com/gin-gonic/gin

replace github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2 => ../../third_party/github.com/grpc-ecosystem/grpc-gateway-v2

replace github.com/searKing/golang/third_party/google.golang.org/grpc => ../../third_party/google.golang.org/grpc
