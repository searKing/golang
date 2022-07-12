package golang

import _ "github.com/searKing/golang/go"
import _ "github.com/searKing/golang/tools"
import _ "github.com/searKing/golang/third_party/google.golang.org/grpc"
import _ "github.com/searKing/golang/third_party/google.golang.org/grpc/grpclog/logruslogger"

import _ "github.com/searKing/golang/third_party/google.golang.org/protobuf"
import _ "github.com/searKing/golang/third_party/github.com/gorilla/websocket"
import _ "github.com/searKing/golang/third_party/github.com/jmoiron/sqlx"
import _ "github.com/searKing/golang/third_party/github.com/urfave/negroni"
import _ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"

import _ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2"
import _ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/go-grpc-middleware"

import _ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway"

import _ "github.com/searKing/golang/third_party/github.com/go-sql-driver/mysql"

import _ "github.com/searKing/golang/third_party/github.com/golang/go"

import _ "github.com/searKing/golang/third_party/github.com/golang/protobuf"
import _ "github.com/searKing/golang/third_party/github.com/google/uuid"
import _ "github.com/searKing/golang/third_party/github.com/julienschmidt/httprouter"
import _ "github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go/metric"
import _ "github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc"
import _ "github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/otlpsql"
import _ "github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/github.com/searKing/otelhttp"
import _ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
import _ "github.com/searKing/golang/third_party/github.com/spf13/viper"
import _ "github.com/searKing/golang/third_party/github.com/spf13/pflag"
import _ "github.com/searKing/golang/third_party/github.com/syndtr/goleveldb"
import _ "github.com/searKing/golang/third_party/github.com/gtank/cryptopasta"
import _ "github.com/searKing/golang/third_party/gocloud.dev"
import _ "github.com/searKing/golang/pkg/webserver"

// PlaceHolder file, so this can be seen as a module.
// https://github.com/ugorji/go/blob/master/FAQ.md#resolving-module-issues
// https://github.com/spf13/cobra/pull/1233

// for f in $(find . -name go.mod)
//	do (cd $(dirname $f); go mod tidy)
// done
// go clean -modcache
// GOPROXY=direct go mod tidy -v
