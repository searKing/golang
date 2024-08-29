// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

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
