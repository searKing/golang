// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"net"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	net_ "github.com/searKing/golang/go/net"
)

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
