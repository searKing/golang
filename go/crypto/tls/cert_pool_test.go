// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/searKing/golang/go/crypto/tls"
	testing_ "github.com/searKing/golang/go/testing"
)

func TestLoadX509CertificatePool(t *testing.T) {
	tmpCertFile, _ := ioutil.TempFile("", "test-cert")
	tmpCertPath := tmpCertFile.Name()
	defer func() {
		_ = os.Remove(tmpCertPath)
	}()
	_ = ioutil.WriteFile(tmpCertPath, []byte(certFileContent), 0600)
	tmpCert, err := tls.LoadCertificates(certFixture, keyFixture, "", "")
	if ok, msg := testing_.NonNil(tmpCert); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.NonError(err); !ok {
		t.Error(msg)
	}

	// 1. no TLS
	certPool, err := tls.LoadX509CertificatePool(nil, "", "")
	if ok, msg := testing_.Nil(certPool); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.EqualError(err, tls.ErrNoCertificatesConfigured); !ok {
		t.Error(msg)
	}

	// 2. inconsistent TLS (ii): warning only
	certPool, err = tls.LoadX509CertificatePool(nil, "x", "")
	if ok, msg := testing_.Nil(certPool); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}
	// 3. invalid TLS string (ii)
	certPool, err = tls.LoadX509CertificatePool(nil, "{}", "")
	if ok, msg := testing_.Nil(certPool); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

	// 4. valid TLS files
	certPool, err = tls.LoadX509CertificatePool(nil, "", tmpCertPath)
	if ok, msg := testing_.NonNil(certPool); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.NonError(err); !ok {
		t.Error(msg)
	}

	// 5. valid TLS strings
	certPool, err = tls.LoadX509CertificatePool(nil, certFixture, "")
	if ok, msg := testing_.NonNil(certPool); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.NonError(err); !ok {
		t.Error(msg)
	}

	// 6. valid TLS cert
	certPool, err = tls.LoadX509CertificatePool(nil, "", "", tmpCert)
	if ok, msg := testing_.NonNil(certPool); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.NonError(err); !ok {
		t.Error(msg)
	}

	// 7. invalid TLS file content
	certPool, err = tls.LoadX509CertificatePool(nil, "", certFixture)
	if ok, msg := testing_.Nil(certPool); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

	// 8. invalid TLS string content
	certPool, err = tls.LoadX509CertificatePool(nil, certFileContent, "")
	if ok, msg := testing_.Nil(certPool); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}
}
