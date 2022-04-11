// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/searKing/golang/pkg/webserver/cors"
	"github.com/searKing/golang/pkg/webserver/healthz"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/burstlimit"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/timeoutlimit"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// ClientMaxReceiveMessageSize use 4GB as the default message size limit.
// grpc library default is 4MB
var defaultMaxReceiveMessageSize = math.MaxInt32 // 1024 * 1024 * 4
var defaultMaxSendMessageSize = math.MaxInt32

// FactoryConfigFunc is an alias for a function that will take in a pointer to an FactoryConfig and modify it
type FactoryConfigFunc func(os *FactoryConfig) error

// FactoryConfig for config
type FactoryConfig struct {
	Validator *validator.Validate

	GatewayOptions []grpc_.GatewayOption
	GinMiddlewares []gin.HandlerFunc

	TlsConfig *tls.Config

	// Name is the human-readable server name, optional
	Name string
	// BindAddress is the host name to use for bind (local internet) facing URLs (e.g. Loopback)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	BindAddress string
	// ExternalAddress is the host name to use for external (public internet) facing URLs (e.g. Swagger)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	ExternalAddress string
	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration

	BindAddr                     Address       // for listen
	AdvertiseAddr                Address       // for expose
	Tls                          *Tls          // for tls such as https
	Cors                         *cors.Factory // for cors
	MaxConcurrencyUnary          int           // for concurrent parallel requests of unary server, The default is 0 (no limit is given)
	BurstLimitTimeoutUnary       time.Duration // for concurrent parallel requests of unary server, The default is 0 (no limit is given)
	MaxConcurrencyStream         int           // for concurrent parallel requests of stream server, The default is 0 (no limit is given)
	BurstLimitTimeoutStream      time.Duration // for concurrent parallel requests of stream server, The default is 0 (no limit is given)
	HandledTimeoutUnary          time.Duration // for max handing time of unary server, The default is 0 (no limit is given)
	HandledTimeoutStream         time.Duration // for max handing time of unary server, The default is 0 (no limit is given)
	MaxReceiveMessageSizeInBytes int           // sets the maximum message size in bytes the grpc can receive, The default is 0 (no limit is given).
	MaxSendMessageSizeInBytes    int           // sets the maximum message size in bytes the grpc can receive, The default is 0 (no limit is given).
	// for debug
	ForceDisableTls bool                // disable tls
	LocalIpResolver *WebLocalIpResolver // for resolve local ip to expose, used if advertise_addr is empty
	NoGrpcProxy     bool                // disable http proxy for grpc client to connect grpc server
}

type Address struct {
	Host    string
	Domains []string // service name to register for dns
	Port    int32
}

type Tls struct {
	Enable        bool
	KeyPairBase64 *WebTlsKeyPair // key pair in base64 format encoded from pem
	KeyPairPath   *WebTlsKeyPair // key pair stored in file from pem
	// service_name is used to verify the hostname on the returned
	// certificates unless InsecureSkipVerify is given. It is also included
	// in the client's handshake to support virtual hosting unless it is
	// an IP address.
	ServiceName      string
	AllowedTlsCidrs  []string //"127.0.0.1/24"
	WhitelistedPaths []string
}

// WebTlsKeyPair represents a public/private key pair
type WebTlsKeyPair struct {
	Cert string // public key
	Key  string // private key
}
type WebLocalIpResolver struct {
	Networks  []string
	Addresses []string
	Timeout   time.Duration
}

func (fc *FactoryConfig) ApplyOptions(configs ...FactoryConfigFunc) error {
	// Apply all Configurations passed in
	for _, config := range configs {
		// Pass the FactoryConfig into the configuration function
		err := config(fc)
		if err != nil {
			return fmt.Errorf("failed to apply configuration function: %w", err)
		}
	}
	return nil
}

// SetDefaults sets sensible values for unset fields in config. This is
// exported for testing: Configs passed to repository functions are copied and have
// default values set automatically.
func (fc *FactoryConfig) SetDefaults() {
}

// Validate inspects the fields of the type to determine if they are valid.
func (fc FactoryConfig) Validate() error {
	valid := fc.Validator
	if valid == nil {
		valid = validator.New()
	}
	return valid.Struct(fc)
}

type Factory struct {
	// it's better to keep FactoryConfig as a private attribute,
	// thanks to that we are always sure that our configuration is not changed in the not allowed way
	fc FactoryConfig
}

func NewFactory(fc FactoryConfig) (Factory, error) {
	if err := fc.Validate(); err != nil {
		return Factory{}, fmt.Errorf("invalid config passed to factory: %w", err)
	}

	return Factory{fc: fc}, nil
}

func (f Factory) Config() FactoryConfig {
	return f.fc
}

func (f Factory) NewWebServer() (*WebServer, error) {

	// if there is no port, and we listen on one securely, use that one
	if _, _, err := net.SplitHostPort(f.fc.ExternalAddress); err != nil {
		if f.fc.BindAddress == "" {
			logrus.WithError(err).Errorf("cannot derive external address port without listening on a secure port.")
			return nil, err
		}

		_, port, err := net.SplitHostPort(f.fc.BindAddress)
		if err != nil {
			logrus.WithError(err).Errorf("cannot derive external address from the secure port: %v", err)
			return nil, err
		}
		f.fc.ExternalAddress = net.JoinHostPort(f.fc.ExternalAddress, port)
	}

	grpclog.SetLoggerV2(grpclog.NewLoggerV2(
		logrus.StandardLogger().WriterLevel(logrus.DebugLevel),
		logrus.StandardLogger().WriterLevel(logrus.WarnLevel),
		logrus.StandardLogger().WriterLevel(logrus.ErrorLevel)))
	opts := grpc_.WithDefault()
	if f.fc.NoGrpcProxy {
		opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithNoProxy()))
	}
	opts = append(opts, grpc_.WithLogrusLogger(logrus.StandardLogger()))
	{
		// 设置GRPC最大消息大小
		opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithNoProxy()))
		// http -> grpc client -> grpc server
		if f.fc.MaxReceiveMessageSizeInBytes > 0 {
			opts = append(opts, grpc_.WithGrpcServerOption(grpc.MaxRecvMsgSize(f.fc.MaxReceiveMessageSizeInBytes)))
			opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(f.fc.MaxReceiveMessageSizeInBytes))))
		} else {
			opts = append(opts, grpc_.WithGrpcServerOption(grpc.MaxRecvMsgSize(defaultMaxReceiveMessageSize)))
			opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaultMaxReceiveMessageSize))))
		}
		// http <- grpc client <- grpc server
		if f.fc.MaxSendMessageSizeInBytes > 0 {
			opts = append(opts, grpc_.WithGrpcServerOption(grpc.MaxSendMsgSize(f.fc.MaxSendMessageSizeInBytes)))
			opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(f.fc.MaxSendMessageSizeInBytes))))
		} else {
			opts = append(opts, grpc_.WithGrpcServerOption(grpc.MaxSendMsgSize(defaultMaxSendMessageSize)))
			opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaultMaxSendMessageSize))))
		}
	}
	{
		// recover
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logrus.WithError(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack())).Errorf("recovered in grpc")
			{
				_, _ = os.Stderr.Write([]byte(fmt.Sprintf("panic: %s", p)))
				debug.PrintStack()
				_, _ = os.Stderr.Write([]byte(" [recovered]"))
				_, _ = os.Stderr.Write([]byte("\n"))
			}
			return status.Errorf(codes.Internal, "%s", p)
		}))))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logrus.WithError(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack())).Errorf("recovered in grpc")
			{
				_, _ = os.Stderr.Write([]byte(fmt.Sprintf("panic: %s", p)))
				debug.PrintStack()
				_, _ = os.Stderr.Write([]byte(" [recovered]"))
				_, _ = os.Stderr.Write([]byte("\n"))
			}
			return status.Errorf(codes.Internal, "%s", p)
		}))))
	}
	{
		// handle request timeout
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(
			timeoutlimit.UnaryServerInterceptor(f.fc.HandledTimeoutUnary)))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(
			timeoutlimit.StreamServerInterceptor(f.fc.HandledTimeoutStream)))
	}
	{
		// burst limit
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(
			burstlimit.UnaryServerInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary)))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(
			burstlimit.StreamServerInterceptor(f.fc.MaxConcurrencyStream, f.fc.BurstLimitTimeoutStream)))
	}

	opts = append(opts, f.fc.GatewayOptions...)
	grpcBackend := grpc_.NewGatewayTLS(f.fc.BindAddress, f.fc.TlsConfig, opts...)
	grpcBackend.ApplyOptions()
	grpcBackend.ErrorLog = logrus_.AsStdLogger(logrus.StandardLogger(), logrus.ErrorLevel, "", 0)
	ginBackend := gin.New()
	ginBackend.Use(gin.LoggerWithWriter(logrus.StandardLogger().Writer()))
	ginBackend.Use(gin_.RecoveryWithWriter(grpcBackend.ErrorLog.Writer()))
	ginBackend.Use(gin_.UseHTTPPreflight())
	ginBackend.Use(f.fc.GinMiddlewares...)

	defaultHealthChecks := []healthz.HealthCheck{healthz.PingHealthCheck, healthz.LogHealthCheck}

	s := &WebServer{
		Name:                  f.fc.Name,
		BindAddress:           f.fc.BindAddress,
		ExternalAddress:       f.fc.ExternalAddress,
		ShutdownDelayDuration: f.fc.ShutdownDelayDuration,
		GrpcBackend:           grpcBackend,
		GinBackend:            ginBackend,

		postStartHooks:   map[string]postStartHookEntry{},
		preShutdownHooks: map[string]preShutdownHookEntry{},
		healthzChecks:    defaultHealthChecks,
		livezChecks:      defaultHealthChecks,
		readyzChecks:     defaultHealthChecks,
		readinessStopCh:  make(chan struct{}),
	}

	return s, nil
}
