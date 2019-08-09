package tls_test

import (
	"github.com/searKing/golang/go/crypto/tls"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadX509CertificatePool(t *testing.T) {
	tmpCertFile, _ := ioutil.TempFile("", "test-cert")
	tmpCertPath := tmpCertFile.Name()
	defer func() {
		_ = os.Remove(tmpCertPath)
	}()
	_ = ioutil.WriteFile(tmpCertPath, []byte(certFileContent), 0600)
	tmpCert, err := tls.LoadCertificates(certFixture, keyFixture, "", "")
	assert.NotNil(t, tmpCert)
	assert.NoError(t, err)

	viper.AutomaticEnv() // read in environment variables that match

	// 1. no TLS
	certPool, err := tls.LoadX509CertificatePool(nil, "", "")
	assert.Nil(t, certPool)
	assert.EqualError(t, err, tls.ErrNoCertificatesConfigured.Error())

	// 2. inconsistent TLS (ii): warning only
	certPool, err = tls.LoadX509CertificatePool(nil, "x", "")
	assert.Nil(t, certPool)
	assert.Error(t, err)

	// 3. invalid TLS string (ii)
	certPool, err = tls.LoadX509CertificatePool(nil, "{}", "")
	assert.Nil(t, certPool)
	assert.Error(t, err)

	// 4. valid TLS files
	certPool, err = tls.LoadX509CertificatePool(nil, "", tmpCertPath)
	assert.NotNil(t, certPool)
	assert.NoError(t, err)

	// 5. valid TLS strings
	certPool, err = tls.LoadX509CertificatePool(nil, certFixture, "")
	assert.NotNil(t, certPool)
	assert.NoError(t, err)

	// 6. valid TLS cert
	certPool, err = tls.LoadX509CertificatePool(nil, "", "", tmpCert)
	assert.NotNil(t, certPool)
	assert.NoError(t, err)

	// 7. invalid TLS file content
	certPool, err = tls.LoadX509CertificatePool(nil, "", certFixture)
	assert.Nil(t, certPool)
	assert.Error(t, err)

	// 8. invalid TLS string content
	certPool, err = tls.LoadX509CertificatePool(nil, certFileContent, "")
	assert.Nil(t, certPool)
	assert.Error(t, err)
}
