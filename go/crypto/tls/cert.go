// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
)

// This code is borrowed from https://github.com/ory/x/blob/master/tlsx/cert.go

// ErrNoCertificatesConfigured is returned when no TLS configuration was found.
var ErrNoCertificatesConfigured = errors.New("no tls configuration was found")

// ErrInvalidCertificateConfiguration is returned when an invalid TLS configuration was found.
var ErrInvalidCertificateConfiguration = errors.New("tls configuration is invalid")

// LoadCertificates returns loads a TLS LoadCertificates.
// certString: Base64 encoded (without padding) string of the TLS certificate (PEM encoded) to be used for HTTP over TLS (HTTPS).
// Example: certString="-----BEGIN CERTIFICATE-----\nMIIDZTCCAk2gAwIBAgIEV5xOtDANBgkqhkiG9w0BAQ0FADA0MTIwMAYDVQQDDClP..."
// keyString: Base64 encoded (without padding) string of the private key (PEM encoded) to be used for HTTP over TLS (HTTPS).
// Example: keyString="-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFDjBABgkqhkiG9w0BBQ0wMzAbBgkqhkiG9w0BBQwwDg..."
// certPath: The path to the TLS certificate (pem encoded).
// Example: certPath=~/cert.pem
// keyPath: The path to the TLS private key (pem encoded).
// Example: keyPath=~/key.pem
// certs: certs of tls.Certificate, *tls.Certificate
func LoadCertificates(
	certString, keyString string,
	certFile, keyFile string,
	certs ...any,
) ([]tls.Certificate, error) {
	if certString == "" && keyString == "" && certFile == "" && keyFile == "" && len(certs) == 0 {
		return nil, ErrNoCertificatesConfigured
	}
	if certString != "" && keyString != "" {
		tlsCertBytes, err := base64.StdEncoding.DecodeString(certString)
		if err != nil {
			return nil, fmt.Errorf("unable to base64 decode the TLS certificate: %v", err)
		}
		tlsKeyBytes, err := base64.StdEncoding.DecodeString(keyString)
		if err != nil {
			return nil, fmt.Errorf("unable to base64 decode the TLS private key: %v", err)
		}

		cert, err := tls.X509KeyPair(tlsCertBytes, tlsKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to load X509 key pair: %v", err)
		}
		return []tls.Certificate{cert}, nil
	}

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("unable to load X509 key pair from files: %v", err)
		}
		return []tls.Certificate{cert}, nil
	}
	var uniformedCerts []tls.Certificate
	for _, cert := range of(certs...) {
		switch cert.(type) {
		case *tls.Certificate:
			tlsCert := cert.(*tls.Certificate)
			uniformedCerts = append(uniformedCerts, *tlsCert)
		case tls.Certificate:
			tlsCert := cert.(tls.Certificate)
			uniformedCerts = append(uniformedCerts, tlsCert)
		default:
			return nil, fmt.Errorf("unable to load X509 key pair from cert: %v", cert)
		}
	}

	return nil, ErrInvalidCertificateConfiguration
}

// LoadX509Certificates returns loads a TLS LoadCertificates of x509.
func LoadX509Certificates(
	certString, keyString string,
	certFile, keyFile string,
) ([]*x509.Certificate, error) {
	certs, err := LoadCertificates(certString, keyString, certFile, keyFile)
	if err != nil {
		return nil, err
	}
	var x509Certs []*x509.Certificate
	for _, cert := range certs {
		for _, certBytes := range cert.Certificate {
			x509Cert, err := x509.ParseCertificate(certBytes)
			if err != nil {
				return nil, err
			}
			x509Certs = append(x509Certs, x509Cert)
		}
	}
	if len(x509Certs) == 0 {
		return nil, ErrNoCertificatesConfigured
	}
	return x509Certs, nil
}

func LoadCertificateAndPool(
	certPool *x509.CertPool,
	certString, keyString string,
	certFile, keyFile string,
) ([]tls.Certificate, *x509.CertPool, error) {
	certs, err := LoadCertificates(certString, keyString, certFile, keyFile)
	if err != nil {
		return nil, nil, err
	}
	certPool, err = LoadX509CertificatePool(certPool, "", "", certs)
	if err != nil {
		return nil, nil, err
	}
	return certs, certPool, nil

}
