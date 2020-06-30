// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

var _ http.Handler = &rejectInsecure{}

//go:generate go-option -type "rejectInsecure"
type rejectInsecure struct {
	// ErrorLog specifies an optional logger for errors accepting
	// connections, unexpected behavior from handlers, and
	// underlying FileSystem errors.
	// If nil, logging will be done via the log package's standard logger.
	ErrorLog *log.Logger
	// ForceHttp allows any request, as a shortcut circuit
	ForceHttp bool
	// AllowedTlsCidrs allows any request which client or proxy's ip included
	// a cidr is a CIDR notation IP address and prefix length,
	// like "192.0.2.0/24" or "2001:db8::/32", as defined in
	// RFC 4632 and RFC 4291.
	AllowedTlsCidrs []string

	// WhitelistedPaths allows any request which http path matches
	WhitelistedPaths []string

	next http.Handler
}

// RejectInsecureServerInterceptor returns a new server interceptor with tls check.
// reject the request fulfills tls's constraints,
func RejectInsecureServerInterceptor(next http.Handler, opts ...RejectInsecureOption) *rejectInsecure {
	r := &rejectInsecure{
		next: next,
	}
	r.ApplyOptions(opts...)
	return r
}

func (m *rejectInsecure) logf(format string, args ...interface{}) {
	if m.ErrorLog != nil {
		m.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (m *rejectInsecure) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m == nil {
		m.next.ServeHTTP(w, r)
		return
	}
	m.rejectInsecureRequests(w, r)
	return
}

// rejectInsecureRequests refused if tls's constraints not passed
func (m *rejectInsecure) rejectInsecureRequests(w http.ResponseWriter, r *http.Request) {
	if m == nil || m.ForceHttp {
		m.next.ServeHTTP(w, r)
		return
	}

	err := DoesRequestSatisfyTlsTermination(r, m.WhitelistedPaths, m.AllowedTlsCidrs)
	if err != nil {
		m.logf("http: could not serve http connection %v: %v", r.RemoteAddr, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)

		if err := json.NewEncoder(w).Encode(fmt.Errorf("cannot serve request over insecure http: %w", err)); err != nil {
			// There was an error, but there's actually not a lot we can do except log that this happened.
			m.logf("http: could not write jsonError to response writer %v: %v", http.StatusBadGateway, err)
		}
		return
	}
	m.next.ServeHTTP(w, r)
	return
}

// DoesRequestSatisfyTlsTermination returns whether the request fulfills tls's constraints,
// https, path matches any whitelisted paths or ip inclued by any cidr
// whitelistedPath is http path that does not need to be checked
// allowedTLSCIDR is the network includes ip.
func DoesRequestSatisfyTlsTermination(r *http.Request, whitelistedPaths []string, allowedTLSCIDRs []string) error {
	// pass if the request is with tls, that is https
	if r.TLS != nil {
		return nil
	}

	// check if the http request can be passed

	// pass if the request belongs to whitelist
	for _, p := range whitelistedPaths {
		if r.URL.Path == p {
			return nil
		}
	}

	if len(allowedTLSCIDRs) == 0 {
		return errors.New("TLS termination is not enabled")
	}

	if err := matchesAnyCidr(r, allowedTLSCIDRs); err != nil {
		return err
	}

	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		return errors.New("X-Forwarded-Proto header is missing")
	}
	if proto != "https" {
		return fmt.Errorf("expected X-Forwarded-Proto header to be https, got %s", proto)
	}

	return nil
}

// matchesAnyCidr returns true if any of client and proxy's ip matches any cidr
// a cidr is a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in
// RFC 4632 and RFC 4291.
func matchesAnyCidr(r *http.Request, cidrs []string) error {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}

	check := []string{remoteIP}
	// X-Forwarded-For: client1, proxy1, proxy2
	for _, fwd := range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		check = append(check, strings.TrimSpace(fwd))
	}

	for _, rn := range cidrs {
		_, cidr, err := net.ParseCIDR(rn)
		if err != nil {
			return err
		}

		for _, ip := range check {
			addr := net.ParseIP(ip)
			if cidr.Contains(addr) {
				return nil
			}
		}
	}
	return fmt.Errorf("neither remote address nor any x-forwarded-for values match CIDR cidrs %v: %v, cidrs, check)", cidrs, check)
}
