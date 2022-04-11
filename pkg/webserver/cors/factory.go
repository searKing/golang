// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cors

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
	gincors "github.com/rs/cors/wrapper/gin"
)

// FactoryConfigFunc is an alias for a function that will take in a pointer to an FactoryConfig and modify it
type FactoryConfigFunc func(os *FactoryConfig) error

// FactoryConfig 人脸配准SDK基础设置工厂函数配置
type FactoryConfig struct {
	Validator *validator.Validate

	Enable bool
	// returns Access-Control-Allow-Origin: * if false
	UseConditional bool
	// allowed_origins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	// return Access-Control-Allow-Origin
	AllowedOrigins []string
	// allowed_methods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string
	// allowed_headers is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string
	// exposed_headers indicates which headers are safe to expose to the API of a CORS
	// API specification
	// return Access-Control-Expose-Headers
	ExposedHeaders []string
	// allow_credentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	// return Access-Control-Allow-Credentials
	AllowCredentials bool
	// options_passthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
	// max_age indicates how long the results of a preflight request
	// can be cached
	MaxAge time.Duration
	// debug flag adds additional output to debug server side CORS issues
	Debug bool
}

type WebNet struct {
	Host    string
	Domains []string // service name to register for dns
	Port    int32
}

type WebTls struct {
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
	fc.AllowedOrigins = []string{"*"}
	fc.AllowedMethods = []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete}
	fc.AllowedHeaders = []string{"*"}
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

func (c Factory) options() cors.Options {
	return cors.Options{
		AllowedOrigins:     c.fc.AllowedOrigins,
		AllowedMethods:     c.fc.AllowedMethods,
		AllowedHeaders:     c.fc.AllowedHeaders,
		ExposedHeaders:     c.fc.ExposedHeaders,
		AllowCredentials:   c.fc.AllowCredentials,
		MaxAge:             int(c.fc.MaxAge.Truncate(time.Second).Seconds()),
		OptionsPassthrough: c.fc.OptionsPassthrough,
	}
}

func (c Factory) New() *cors.Cors {
	return cors.New(c.options())
}

func (c Factory) NewWrapper() func(http.Handler) http.Handler {
	return c.New().Handler
}

func (c Factory) NewGinHandler() gin.HandlerFunc {
	return gincors.New(c.options())
}
