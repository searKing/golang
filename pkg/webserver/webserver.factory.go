// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	slog_ "github.com/searKing/golang/go/log/slog"
	net_ "github.com/searKing/golang/go/net"
	"github.com/searKing/golang/pkg/webserver/healthz"
	"github.com/searKing/golang/pkg/webserver/pkg/recovery"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/burstlimit"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/timeoutlimit"
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
	// ExternalAddress is the address advertised, even if BindAddress is a loopback. By default, this
	// is set to BindAddress if the later no loopback, or to the first host interface address.
	ExternalAddress string
	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration

	TlsConfig                      *tls.Config
	Cors                           cors.Options     // for cors
	ForceDisableTls                bool             // disable tls
	LocalIpResolver                *LocalIpResolver // for resolve local ip to expose, used if advertise_addr is empty
	NoGrpcProxy                    bool             // disable http proxy for grpc client to connect grpc server
	PreferRegisterHTTPFromEndpoint bool             // prefer register http handler from endpoint

	// grpc middlewares
	MaxConcurrencyUnary          int           // for concurrent parallel requests of unary server, The default is 0 (no limit is given)
	MaxConcurrencyStream         int           // for concurrent parallel requests of stream server, The default is 0 (no limit is given)
	BurstLimitTimeoutUnary       time.Duration // for concurrent parallel requests of unary server, The default is 0 (no limit is given)
	BurstLimitTimeoutStream      time.Duration // for concurrent parallel requests of stream server, The default is 0 (no limit is given)
	HandledTimeoutUnary          time.Duration // for max handing time of unary server, The default is 0 (no limit is given)
	HandledTimeoutStream         time.Duration // for max handing time of unary server, The default is 0 (no limit is given)
	MaxReceiveMessageSizeInBytes int           // sets the maximum message size in bytes the grpc server can receive, The default is 0 (no limit is given).
	MaxSendMessageSizeInBytes    int           // sets the maximum message size in bytes the grpc server can send, The default is 0 (no limit is given).

	// Deprecated: takes no effect, use slog instead.
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
		Logger:         slog.NewLogLogger(slog.Default().Handler(), slog.LevelDebug),
	}
}

// Validate inspects the fields of the type to determine if they are valid.
func (fc *FactoryConfig) Validate() error {
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
			slog.Error("cannot derive external address port without listening on a secure port.", slog_.Error(err))
			os.Exit(1)
		}

		_, port, err := net.SplitHostPort(f.fc.BindAddress)
		if err != nil {
			slog.Error("cannot derive external address from the secure port", slog_.Error(err))
			os.Exit(1)
		}
		f.fc.ExternalAddress = net.JoinHostPort(f.fc.ExternalAddress, port)
	}

	opts := grpc_.WithDefault()
	if f.fc.NoGrpcProxy {
		opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithNoProxy()))
	}
	{
		// 设置GRPC最大消息大小
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
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(recovery.UnaryServerInterceptor()))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(recovery.StreamServerInterceptor()))
	}
	{
		// handle request timeout
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(timeoutlimit.UnaryServerInterceptor(f.fc.HandledTimeoutUnary)))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(timeoutlimit.StreamServerInterceptor(f.fc.HandledTimeoutStream)))
	}
	{
		// burst limit
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(burstlimit.UnaryServerInterceptor(f.fc.MaxConcurrencyUnary, f.fc.BurstLimitTimeoutUnary)))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(burstlimit.StreamServerInterceptor(f.fc.MaxConcurrencyStream, f.fc.BurstLimitTimeoutStream)))
	}

	// cors
	opts = append(opts, grpc_.WithHttpWrapper(cors.New(f.fc.Cors).Handler))
	opts = append(opts, f.fc.GatewayOptions...)
	opts = append(opts, grpc_.WithSlogLoggerConfig(slog.Default().Handler(), grpc_.ExtractLoggingOptions(opts...))...)
	grpcBackend := grpc_.NewGatewayTLS(f.fc.BindAddress, f.fc.TlsConfig, opts...)
	{
		l := slog.NewLogLogger(slog.Default().Handler(), slog.LevelError)
		l.SetFlags(log.Lshortfile)
		grpcBackend.ErrorLog = l
	}
	ginBackend := gin.New()

	{
		l := slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo)
		l.SetFlags(log.Lshortfile)
		ginBackend.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: GinLogFormatter("GIN over HTTP"),
			Output:    l.Writer(),
		}))
	}
	ginBackend.Use(gin_.RecoveryWithWriter(grpcBackend.ErrorLog.Writer()))
	ginBackend.Use(gin_.UseHTTPPreflight())
	ginBackend.Use(f.fc.GinMiddlewares...)

	defaultHealthChecks := []healthz.HealthChecker{healthz.PingHealthzCheck, healthz.LogHealthCheck}

	s := &WebServer{
		Name:                           f.fc.Name,
		BindAddress:                    f.fc.BindAddress,
		ExternalAddress:                f.fc.ExternalAddress,
		PreferRegisterHTTPFromEndpoint: f.fc.PreferRegisterHTTPFromEndpoint,
		ShutdownDelayDuration:          f.fc.ShutdownDelayDuration,
		grpcBackend:                    grpcBackend,
		ginBackend:                     ginBackend,

		postStartHooks:   map[string]postStartHookEntry{},
		preShutdownHooks: map[string]preShutdownHookEntry{},
		healthzChecks:    defaultHealthChecks,
		livezChecks:      defaultHealthChecks,
		readyzChecks:     defaultHealthChecks,
		readinessStopCh:  make(chan struct{}),
	}

	err := s.AddBootSequencePostStartHook("__bind_addr__", func(ctx context.Context) error {
		host, port, err := net.SplitHostPort(s.ExternalAddress)
		if err != nil {
			return fmt.Errorf("malformed external address: %w", err)
		}
		_, bindPort, err := net.SplitHostPort(s.grpcBackend.BindAddr())
		if err != nil {
			return fmt.Errorf("malformed bind address: %w", err)
		}
		if bindPort == "" || bindPort == "0" {
			return nil
		}

		if port == "0" {
			logger := slog.With("old_external_address", s.ExternalAddress).
				With("bind_address", s.BindAddress).
				With("bind_port", bindPort)
			s.ExternalAddress = net.JoinHostPort(host, bindPort)
			logger.With("new_external_address", s.ExternalAddress).
				Info("update external address with bind port")
		}
		return nil
	})
	if err != nil {
		return nil, err
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
		if err == nil && len(ip) > 0 {
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

// GetBackendExternalHostPort returns an address to expose with domain, if not set, use host instead.
func (f *Factory) GetBackendExternalHostPort() string {
	host, port, _ := net_.SplitHostPort(f.fc.ExternalAddress)
	if host == "" {
		return f.GetBackendBindHostPort()
	}
	return getHostPort(host, port)
}

// GetBackendServeHostPort returns an address to expose without domain, if not set, use resolver to resolve an ip
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
