// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strings"

	net_ "github.com/searKing/golang/go/net"
	strings_ "github.com/searKing/golang/go/strings"
)

func (fc *FactoryConfig) HTTPScheme() string {
	if fc.ForceDisableTls {
		return "http"
	}
	return "https"
}

func (fc *FactoryConfig) ResolveLocalIp() string {
	resolvers := fc.LocalIpResolver
	ip, err := net_.ServeIP(resolvers.Networks, resolvers.Addresses, resolvers.Timeout)
	if err != nil {
		return "localhost"
	}
	return ip.String()
}

// GetBackendBindHostPort returns a address to listen.
func (fc *FactoryConfig) GetBackendBindHostPort() string {
	local := fc.BindAddr
	return getHostPort(local.Host, local.Port)
}

// GetBackendAdvertiseHostPort returns a address to expose with domain, if not set, use host instead.
func (fc *FactoryConfig) GetBackendAdvertiseHostPort() string {
	extern := fc.AdvertiseAddr
	host := strings_.ValueOrDefault(extern.Domains...)
	if host == "" {
		host = fc.AdvertiseAddr.Host
	}
	if host == "" {
		return fc.GetBackendBindHostPort()
	}
	return getHostPort(host, extern.Port)
}

// GetBackendServeHostPort returns a address to expose without domain, if not set, use resolver to resolve a ip
func (fc *FactoryConfig) GetBackendServeHostPort(advertise bool) string {
	if advertise {
		host := fc.AdvertiseAddr.Host
		if host != "" {
			return getHostPort(host, fc.AdvertiseAddr.Port)
		}
	}

	host := fc.BindAddr.Host
	if host != "" && host != "0.0.0.0" {
		return getHostPort(host, fc.BindAddr.Port)
	}
	return getHostPort(fc.ResolveLocalIp(), fc.BindAddr.Port)
}

func (fc *FactoryConfig) ResolveBackendLocalUrl(relativePaths ...string) string {
	return resolveLocalUrl(
		fc.HTTPScheme(),
		fc.GetBackendServeHostPort(true),
		filepath.Join(relativePaths...)).String()
}

func getHostPort(hostname string, port int32) string {
	if strings.HasPrefix(hostname, "unix:") {
		return hostname
	}
	return fmt.Sprintf("%s:%d", hostname, port)
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
