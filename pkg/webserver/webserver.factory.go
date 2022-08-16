// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/cors"
	net_ "github.com/searKing/golang/go/net"
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

// FactoryConfigFunc is an alias for a function that will take in a pointer to an FactoryConfig and modify it
type FactoryConfigFunc func(os *FactoryConfig) error

// FactoryConfig Config of Factory
type FactoryConfig struct {
	// Name is the human-readable server name, optional
	Name string
	// BindAddress is the host port to bind to (local internet)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	BindAddress string
	// ExternalAddress is the address advertised, even if BindAddress is a loopback. By default this
	// is set to BindAddress if the later no loopback, or to the first host interface address.
	ExternalAddress string
	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration

	TlsConfig       *tls.Config
	Cors            cors.Options     // for cors
	ForceDisableTls bool             // disable tls
	LocalIpResolver *LocalIpResolver // for resolve local ip to expose, used if advertise_addr is empty
	NoGrpcProxy     bool             // disable http proxy for grpc client to connect grpc server

	// grpc middlewares
	MaxConcurrencyUnary          int           // for concurrent parallel requests of unary server, The default is 0 (no limit is given)
	MaxConcurrencyStream         int           // for concurrent parallel requests of stream server, The default is 0 (no limit is given)
	BurstLimitTimeoutUnary       time.Duration // for concurrent parallel requests of unary server, The default is 0 (no limit is given)
	BurstLimitTimeoutStream      time.Duration // for concurrent parallel requests of stream server, The default is 0 (no limit is given)
	HandledTimeoutUnary          time.Duration // for max handing time of unary server, The default is 0 (no limit is given)
	HandledTimeoutStream         time.Duration // for max handing time of unary server, The default is 0 (no limit is given)
	MaxReceiveMessageSizeInBytes int           // sets the maximum message size in bytes the grpc server can receive, The default is 0 (no limit is given).
	MaxSendMessageSizeInBytes    int           // sets the maximum message size in bytes the grpc server can send, The default is 0 (no limit is given).

	EnableLogrusMiddleware bool // disable logrus middleware

	GatewayOptions []grpc_.GatewayOption
	GinMiddlewares []gin.HandlerFunc
}

// SetDefaults sets sensible values for unset fields in config. This is
// exported for testing: Configs passed to repository functions are copied and have
// default values set automatically.
func (fc *FactoryConfig) SetDefaults() {
	fc.Cors = cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete},
		AllowedHeaders: []string{"*"},
	}
}

// Validate inspects the fields of the type to determine if they are valid.
func (fc FactoryConfig) Validate() error {
	return nil
}

type Factory struct {
	// it's better to keep FactoryConfig as a private attribute,
	// thanks to that we are always sure that our configuration is not changed in the not allowed way
	fc FactoryConfig
}

func NewFactory(fc FactoryConfig, configs ...FactoryConfigFunc) (Factory, error) {
	// Apply all Configurations passed in
	for _, config := range configs {
		// Pass the FactoryConfig into the configuration function
		err := config(&fc)
		if err != nil {
			return Factory{}, fmt.Errorf("failed to apply configuration function: %w", err)
		}
	}

	if err := fc.Validate(); err != nil {
		return Factory{}, fmt.Errorf("invalid config passed to factory: %w", err)
	}

	f := Factory{fc: fc}

	return f, nil
}

func (f *Factory) Config() FactoryConfig {
	return f.fc
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (f *Factory) New() (*WebServer, error) {
	f.fc.BindAddress = f.GetBackendBindHostPort()
	f.fc.ExternalAddress = f.GetBackendServeHostPort(true)

	// if there is no port, and we listen on one securely, use that one
	if _, _, err := net.SplitHostPort(f.fc.ExternalAddress); err != nil {
		if f.fc.BindAddress == "" {
			logrus.WithError(err).Fatalf("cannot derive external address port without listening on a secure port.")
		}

		_, port, err := net.SplitHostPort(f.fc.BindAddress)
		if err != nil {
			logrus.WithError(err).Fatalf("cannot derive external address from the secure port: %v", err)
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

	// cors
	{
		opts = append(opts, grpc_.WithHttpWrapper(cors.New(f.fc.Cors).Handler))
	}

	opts = append(opts, f.fc.GatewayOptions...)
	if f.fc.EnableLogrusMiddleware {
		opts = append(opts, grpc_.WithLogrusLogger(logrus.StandardLogger()))
	}
	grpcBackend := grpc_.NewGatewayTLS(f.fc.BindAddress, f.fc.TlsConfig, opts...)
	grpcBackend.ApplyOptions()
	grpcBackend.ErrorLog = logrus_.AsStdLogger(logrus.StandardLogger(), logrus.ErrorLevel, "", 0)
	ginBackend := gin.New()
	if f.fc.EnableLogrusMiddleware {
		ginBackend.Use(gin.LoggerWithWriter(logrus.StandardLogger().Writer()))
	}
	ginBackend.Use(gin_.RecoveryWithWriter(grpcBackend.ErrorLog.Writer()))
	ginBackend.Use(gin_.UseHTTPPreflight())
	ginBackend.Use(f.fc.GinMiddlewares...)

	defaultHealthChecks := []healthz.HealthChecker{healthz.PingHealthzCheck, healthz.LogHealthCheck}

	s := &WebServer{
		Name:                  f.fc.Name,
		BindAddress:           f.fc.BindAddress,
		ExternalAddress:       f.fc.ExternalAddress,
		ShutdownDelayDuration: f.fc.ShutdownDelayDuration,
		grpcBackend:           grpcBackend,
		ginBackend:            ginBackend,

		postStartHooks:   map[string]postStartHookEntry{},
		preShutdownHooks: map[string]preShutdownHookEntry{},
		healthzChecks:    defaultHealthChecks,
		livezChecks:      defaultHealthChecks,
		readyzChecks:     defaultHealthChecks,
		readinessStopCh:  make(chan struct{}),
	}

	return s, nil
}

// ClientMaxReceiveMessageSize use 4GB as the default message size limit.
// grpc library default is 4MB
var defaultMaxReceiveMessageSize = math.MaxInt32 // 1024 * 1024 * 4
var defaultMaxSendMessageSize = math.MaxInt32

type Net struct {
	Host    string
	Domains []string // service name to register to consul for dns
	Port    int32
}

// CertKey a public/private key pair
type CertKey struct {
	Cert string // public key, containing a PEM-encoded certificate, and possibly the complete certificate chain
	Key  string // private key, containing a PEM-encoded private key for the certificate specified by CertFile
}

//type CertKey struct {
//	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
//	CertFile string
//	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
//	KeyFile string
//}

type TLS struct {
	Enable        bool
	KeyPairBase64 *CertKey // key pair in base64 format encoded from pem
	KeyPairPath   *CertKey // key pair stored in file from pem
	// service_name is used to verify the hostname on the returned
	// certificates unless InsecureSkipVerify is given. It is also included
	// in the client's handshake to support virtual hosting unless it is
	// an IP address.
	ServiceName      string
	AllowedTlsCidrs  []string //"127.0.0.1/24"
	WhitelistedPaths []string
}

type LocalIpResolver struct {
	Networks  []string
	Addresses []string
	Timeout   time.Duration
}

func (f *Factory) HTTPScheme() string {
	if f.fc.ForceDisableTls {
		return "http"
	}
	return "https"
}

func (f *Factory) ResolveLocalIp() string {
	resolver := f.fc.LocalIpResolver
	if resolver != nil {
		ip, err := net_.ServeIP(resolver.Networks, resolver.Addresses, resolver.Timeout)
		if err != nil && len(ip) > 0 {
			return ip.String()
		}
	}

	// use local ip
	localIP, err := net_.ListenIP()
	if err == nil && len(localIP) > 0 {
		return localIP.String()
	}
	return "localhost"
}

// GetBackendBindHostPort returns a address to listen.
func (f *Factory) GetBackendBindHostPort() string {
	host, port, _ := net_.SplitHostPort(f.fc.BindAddress)
	return getHostPort(host, port)
}

// GetBackendExternalHostPort returns a address to expose with domain, if not set, use host instead.
func (f *Factory) GetBackendExternalHostPort() string {
	host, port, _ := net_.SplitHostPort(f.fc.ExternalAddress)
	if host == "" {
		return f.GetBackendBindHostPort()
	}
	return getHostPort(host, port)
}

// GetBackendServeHostPort returns a address to expose without domain, if not set, use resolver to resolve a ip
func (f *Factory) GetBackendServeHostPort(external bool) string {
	if external {
		host, _, _ := net_.SplitHostPort(f.fc.ExternalAddress)
		if host != "" {
			return f.GetBackendExternalHostPort()
		}
	}

	host, port, _ := net_.SplitHostPort(f.fc.BindAddress)
	if host != "" && host != "0.0.0.0" {
		return f.GetBackendBindHostPort()
	}
	return getHostPort(f.ResolveLocalIp(), port)
}

func (f *Factory) ResolveBackendLocalUrl(relativePaths ...string) string {
	return resolveLocalUrl(
		f.HTTPScheme(),
		f.GetBackendServeHostPort(true),
		filepath.Join(relativePaths...)).String()
}

func getHostPort(hostname string, port string) string {
	if strings.HasPrefix(hostname, "unix:") {
		return hostname
	}

	return net.JoinHostPort(hostname, port)
}

func resolveLocalUrl(scheme, hostport, path string) *url.URL {
	u := &url.URL{
		Scheme: scheme,
		Host:   hostport,
		Path:   path,
	}
	if u.Hostname() == "" {
		// use local host
		localHost := "localhost"

		// use local ip
		localIP, err := net_.ListenIP()
		if err == nil && len(localIP) > 0 {
			localHost = localIP.String()
		}
		u.Host = net.JoinHostPort(localHost, u.Port())
	}
	return u
}
