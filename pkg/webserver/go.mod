module github.com/searKing/golang/pkg/webserver

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator/v10 v10.9.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/rs/cors v1.8.2
	github.com/rs/cors/wrapper/gin v0.0.0-20220223021805-a4a5ce87d5a2
	github.com/searKing/golang v1.2.0
	github.com/searKing/golang/third_party/github.com/gin-gonic/gin v1.1.72
	github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2 v1.1.74
	github.com/searKing/golang/third_party/github.com/sirupsen/logrus v1.1.70
	github.com/searKing/golang/third_party/google.golang.org/grpc v1.1.72
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.45.0
)

replace github.com/searKing/golang v1.2.0 => ../../../../../github.com/searKing/golang
