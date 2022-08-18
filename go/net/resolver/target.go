// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"fmt"
	"net/url"
	"strings"
)

// Target represents a target for gRPC, as specified in:
// https://github.com/grpc/grpc/blob/master/doc/naming.md.
// It is parsed from the target string that gets passed into Dial or DialContext
// by the user. And gRPC passes it to the resolver and the balancer.
//
// If the target follows the naming spec, and the parsed scheme is registered
// with gRPC, we will parse the target string according to the spec. If the
// target does not contain a scheme or if the parsed scheme is not registered
// (i.e. no corresponding resolver available to resolve the endpoint), we will
// apply the default scheme, and will attempt to reparse it.
//
// Examples:
//
//   - "dns://some_authority/foo.bar"
//     Target{Scheme: "dns", Authority: "some_authority", Endpoint: "foo.bar"}
//   - "foo.bar"
//     Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "foo.bar"}
//   - "unknown_scheme://authority/endpoint"
//     Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "unknown_scheme://authority/endpoint"}
//
// If the target does not contain a scheme, we will apply the default scheme, and set the Target to
// be the full target string. e.g. "foo.bar" will be parsed into
// &Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "foo.bar"}.
//
// If the parsed scheme is not registered (i.e. no corresponding resolver available to resolve the
// endpoint), we set the Scheme to be the default scheme, and set the Endpoint to be the full target
// string. e.g. target string "unknown_scheme://authority/endpoint" will be parsed into
// &Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "unknown_scheme://authority/endpoint"}.

type Target struct {
	// Deprecated: use URL.Scheme instead.
	Scheme string
	// Deprecated: use URL.Host instead.
	Authority string
	// Deprecated: use URL.Path or URL.Opaque instead. The latter is set when
	// the former is empty.
	Endpoint string
	// URL contains the parsed dial target with an optional default scheme added
	// to it if the original dial target contained no scheme or contained an
	// unregistered scheme. Any query params specified in the original dial
	// target can be accessed from here.
	URL url.URL
}

func (t *Target) key() targetKey {
	return targetKey{
		Scheme:    t.Scheme,
		Authority: t.Authority,
		Endpoint:  t.Endpoint,
	}
}

type targetKey Target

func (k targetKey) String() string {
	// Only used by tests.
	return fmt.Sprintf("%s|%s|%s", k.Scheme, k.Authority, k.Endpoint)
}

// ParseTarget uses RFC 3986 semantics to parse the given target into a
// resolver.Target struct containing scheme, authority and endpoint. Query
// params are stripped from the endpoint.
//
// If target is not a valid scheme://authority/endpoint as specified in
// https://github.com/grpc/grpc/blob/master/doc/naming.md,
// it returns {Endpoint: target}.
// Code borrowed from https://github.com/grpc/grpc-go/blob/v1.48.0/clientconn.go#L1619
// See https://github.com/grpc/grpc-go/pull/4817
func ParseTarget(target string) Target {
	parsedTarget, err := parseTarget(target)
	if err == nil && parsedTarget.Scheme != "" {
		return parsedTarget
	}

	// We are here because the user's dial target did not contain a scheme or
	// specified an unregistered scheme. We should fallback to the default
	// scheme, except when a custom dialer is specified in which case, we should
	// always use passthrough scheme.
	defScheme := GetDefaultScheme()
	canonicalTarget := defScheme + ":///" + target

	parsedTarget, err = parseTarget(canonicalTarget)
	if err != nil {
		return Target{
			Endpoint: target,
			URL: url.URL{
				Path: target,
			},
		}
	}
	parsedTarget.Scheme = "" // trim scheme
	parsedTarget.URL.Scheme = ""
	return parsedTarget
}

// parseTarget uses RFC 3986 semantics to parse the given target into a
// resolver.Target struct containing scheme, authority and endpoint. Query
// params are stripped from the endpoint.
func parseTarget(target string) (Target, error) {
	u, err := url.Parse(target)
	if err != nil {
		return Target{}, err
	}
	// For targets of the form "[scheme]://[authority]/endpoint, the endpoint
	// value returned from url.Parse() contains a leading "/". Although this is
	// in accordance with RFC 3986, we do not want to break existing resolver
	// implementations which expect the endpoint without the leading "/". So, we
	// end up stripping the leading "/" here. But this will result in an
	// incorrect parsing for something like "unix:///path/to/socket". Since we
	// own the "unix" resolver, we can workaround in the unix resolver by using
	// the `URL` field instead of the `Endpoint` field.
	endpoint := u.Path
	if endpoint == "" {
		endpoint = u.Opaque
	}
	endpoint = strings.TrimPrefix(endpoint, "/")
	return Target{
		Scheme:    u.Scheme,
		Authority: u.Host,
		Endpoint:  endpoint,
		URL:       *u,
	}, nil
}
