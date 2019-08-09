package tls

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
)

// LoadCertificatePool returns loads a TLS x509.CertPool or update a TLS x509.CertPool if nil.
// certString: Base64 encoded (without padding) string of the TLS certificate (PEM encoded) to be used for HTTP over TLS (HTTPS).
// Example: certString="-----BEGIN CERTIFICATE-----\nMIIDZTCCAk2gAwIBAgIEV5xOtDANBgkqhkiG9w0BAQ0FADA0MTIwMAYDVQQDDClP..."
// certPath: The path to the TLS certificate (pem encoded).
// Example: certPath=~/cert.pem
func LoadCertificatePool(
	certPool *x509.CertPool,
	certString string,
	certFile string,
) (*x509.CertPool, error) {
	var tlsCertBytes []byte
	var err error
	if certString == "" && certFile == "" {
		return nil, errors.WithStack(ErrNoCertificatesConfigured)
	} else if certString != "" {
		tlsCertBytes, err = base64.StdEncoding.DecodeString(certString)
		if err != nil {
			return nil, fmt.Errorf("unable to base64 decode the TLS certificate: %v", err)
		}
	} else if certFile != "" {
		tlsCertBytes, err = ioutil.ReadFile(certFile)
		if err != nil {
			return nil, err
		}
	}
	if len(tlsCertBytes) == 0 {
		return nil, errors.WithStack(ErrInvalidCertificateConfiguration)
	}
	if certPool == nil {
		certPool = x509.NewCertPool()
	}
	if !certPool.AppendCertsFromPEM(tlsCertBytes) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return certPool, nil

}
