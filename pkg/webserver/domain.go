// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"time"
)

type Cors struct {
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
