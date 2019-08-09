package tls_test

import (
	"github.com/searKing/golang/go/crypto/tls"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestHTTPSCertificatePool(t *testing.T) {
	certFixture := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVFRENDQXZpZ0F3SUJBZ0lKQU5mK0lUMU1HaHhCTUEwR0NTcUdTSWI` +
		`zRFFFQkN3VUFNSUdaTVFzd0NRWUQKVlFRR0V3SlZVekVMTUFrR0ExVUVDQXdDUTBFeEVqQVFCZ05WQkFjTUNWQmhiRzhnUVd4MGJ6RWlNQ0FHQ` +
		`TFVRQpDZ3daVDI1bFEyOXVZMlZ5YmlCYmRHVnpkQ0J3ZFhKd2IzTmxYVEVjTUJvR0ExVUVBd3dUYjI1bFkyOXVZMlZ5CmJpMTBaWE4wTG1OdmJ` +
		`URW5NQ1VHQ1NxR1NJYjNEUUVKQVJZWVpuSmxaR1Z5YVdOQVkyOXVaV052Ym1ObGNtNHUKWTI5dE1CNFhEVEU0TURnd016RTJNakUwT0ZvWERUR` +
		`TVNVEl4TmpFMk1qRTBPRm93Z1lReEN6QUpCZ05WQkFZVApBbFZUTVFzd0NRWURWUVFJREFKRFFURVNNQkFHQTFVRUJ3d0pVR0ZzYnlCQmJIUnZ` +
		`NU0l3SUFZRFZRUUxEQmxQCmJtVkRiMjVqWlhKdUlGdDBaWE4wSUhCMWNuQnZjMlZkTVRBd0xnWURWUVFERENkaGNHa3RjMlZ5ZG1salpTMXcKY` +
		`205NGFXVmtMbTl1WldOdmJtTmxjbTR0ZEdWemRDNWpiMjB3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQgpEd0F3Z2dFS0FvSUJBUURXVzF` +
		`KQnZweC9vZkYwei80QnkrYmdBcCtoYnlxblVsQ2FnYmlneE9QTHY3aUg4TSt1CjNENkRlSVkzQzdkV0thTjRnYXZHd1MvN3I0UWxXSWdvK09NR` +
		`HQ1M25OZDVvakwvNWY5R1E0ZGRObW53b25EeEYKVThrd1lMWURMTkJIQzJqMzFBNVNueHo0S1NkVE03Rmc0OFBJeTNBaWFGMkhEcURZVlJpWkV` +
		`ackl4U3JTSmFKZgp1WGVCSUVBcFBpUG1IOURObGw2VVo3ODZvZitJWWVLV2VuY0MvbGpPaGlJSnJWL3NEZTc2QVFjdXY5T29XaUdiCklGVFMyW` +
		`ExSRGF0YzByQXhWdlFiTnMzeWlFYjh3UzBaR0F4cTBuZk9pMGZkYVBIODdFc25MdkpqWk5PcXIvTVMKSW5BYmN2ZmlwckxxaEdLQTVIN2hKVGZ` +
		`EcFJ6WWxBcm5maTJMQWdNQkFBR2piakJzTUFrR0ExVWRFd1FDTUFBdwpDd1lEVlIwUEJBUURBZ1hnTUZJR0ExVWRFUVJMTUVtQ0htOWhkR2hyW` +
		`ldWd1pYSXViMjVsWTI5dVkyVnliaTEwClpYTjBMbU52YllJbllYQnBMWE5sY25acFkyVXRjSEp2ZUdsbFpDNXZibVZqYjI1alpYSnVMWFJsYzN` +
		`RdVkyOXQKTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFCMVBibCtSbW50RW9jbHlqWXpzeWtLb2lYczNwYTgzQ2dEWjZwQwpncnY0TFF4U29FZ` +
		`kowNGY4YkQ0SUlZRkdDWmZWTkcwVnBFWHJObGs2VWJzVmRUQUJ0cUNndUpUV3dER1VBaDZYCjNiRmhyWm5QZXhzLy9Rd2dEQWRxSWYwRWd3Y0R` +
		`VRzc2R0lkZms3MGUxWnV4Y2h4ZDhVQkNwQUlkZVUwOHZWa3kKNFBXdjJLNGFENEZqQ2hLeENONWtoTjUwRk1QY2FJK3hWZ2Q0N3RQaFZOOWxRa` +
		`W9HRENoc1Q1dkFSazdiYS9jZQowUTlOV2RpTWZMRWdMZGNCb2JaS0Z0RnJsS3R5ek9nRGpMdlh2TFFzL3MybWVyU0k5Zmt3b09CRVArN2o3Wm5` +
		`zCkFqeTlNZmh3cWJUcFc3S3BDU0ZhMFZULzJ1OTVaUmNQdnJYbGRLUnlnQjRXdUFScgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==`
	certFileContent := `-----BEGIN CERTIFICATE-----
MIIEEDCCAvigAwIBAgIJANf+IT1MGhxBMA0GCSqGSIb3DQEBCwUAMIGZMQswCQYD
VQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCVBhbG8gQWx0bzEiMCAGA1UE
CgwZT25lQ29uY2VybiBbdGVzdCBwdXJwb3NlXTEcMBoGA1UEAwwTb25lY29uY2Vy
bi10ZXN0LmNvbTEnMCUGCSqGSIb3DQEJARYYZnJlZGVyaWNAY29uZWNvbmNlcm4u
Y29tMB4XDTE4MDgwMzE2MjE0OFoXDTE5MTIxNjE2MjE0OFowgYQxCzAJBgNVBAYT
AlVTMQswCQYDVQQIDAJDQTESMBAGA1UEBwwJUGFsbyBBbHRvMSIwIAYDVQQLDBlP
bmVDb25jZXJuIFt0ZXN0IHB1cnBvc2VdMTAwLgYDVQQDDCdhcGktc2VydmljZS1w
cm94aWVkLm9uZWNvbmNlcm4tdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDWW1JBvpx/ofF0z/4By+bgAp+hbyqnUlCagbigxOPLv7iH8M+u
3D6DeIY3C7dWKaN4gavGwS/7r4QlWIgo+OMDt53nNd5ojL/5f9GQ4ddNmnwonDxF
U8kwYLYDLNBHC2j31A5Snxz4KSdTM7Fg48PIy3AiaF2HDqDYVRiZEZrIxSrSJaJf
uXeBIEApPiPmH9DNll6UZ786of+IYeKWencC/ljOhiIJrV/sDe76AQcuv9OoWiGb
IFTS2XLRDatc0rAxVvQbNs3yiEb8wS0ZGAxq0nfOi0fdaPH87EsnLvJjZNOqr/MS
InAbcvfiprLqhGKA5H7hJTfDpRzYlArnfi2LAgMBAAGjbjBsMAkGA1UdEwQCMAAw
CwYDVR0PBAQDAgXgMFIGA1UdEQRLMEmCHm9hdGhrZWVwZXIub25lY29uY2Vybi10
ZXN0LmNvbYInYXBpLXNlcnZpY2UtcHJveGllZC5vbmVjb25jZXJuLXRlc3QuY29t
MA0GCSqGSIb3DQEBCwUAA4IBAQB1Pbl+RmntEoclyjYzsykKoiXs3pa83CgDZ6pC
grv4LQxSoEfJ04f8bD4IIYFGCZfVNG0VpEXrNlk6UbsVdTABtqCguJTWwDGUAh6X
3bFhrZnPexs//QwgDAdqIf0EgwcDUG76GIdfk70e1Zuxchxd8UBCpAIdeU08vVky
4PWv2K4aD4FjChKxCN5khN50FMPcaI+xVgd47tPhVN9lQioGDChsT5vARk7ba/ce
0Q9NWdiMfLEgLdcBobZKFtFrlKtyzOgDjLvXvLQs/s2merSI9fkwoOBEP+7j7Zns
Ajy9MfhwqbTpW7KpCSFa0VT/2u95ZRcPvrXldKRygB4WuARr
-----END CERTIFICATE-----`
	tmpCertFile, _ := ioutil.TempFile("", "test-cert")
	tmpCert := tmpCertFile.Name()
	defer func() {
		_ = os.Remove(tmpCert)
	}()
	_ = ioutil.WriteFile(tmpCert, []byte(certFileContent), 0600)
	viper.AutomaticEnv() // read in environment variables that match

	// 1. no TLS
	certPool, err := tls.LoadCertificatePool(nil, "", "")
	assert.Nil(t, certPool)
	assert.EqualError(t, err, tls.ErrNoCertificatesConfigured.Error())

	// 2. inconsistent TLS (ii): warning only
	certPool, err = tls.LoadCertificatePool(nil, "x", "")
	assert.Nil(t, certPool)
	assert.Error(t, err)

	// 3. invalid TLS string (ii)
	certPool, err = tls.LoadCertificatePool(nil, "{}", "")
	assert.Nil(t, certPool)
	assert.Error(t, err)

	// 4. valid TLS files
	certPool, err = tls.LoadCertificatePool(nil, "", tmpCert)
	assert.NotNil(t, certPool)
	assert.NoError(t, err)

	// 5. valid TLS strings
	certPool, err = tls.LoadCertificatePool(nil, certFixture, "")
	assert.NotNil(t, certPool)
	assert.NoError(t, err)

	// 7. invalid TLS file content
	certPool, err = tls.LoadCertificatePool(nil, "", certFixture)
	assert.Nil(t, certPool)
	assert.Error(t, err)

	// 8. invalid TLS string content
	certPool, err = tls.LoadCertificatePool(nil, certFileContent, "")
	assert.Nil(t, certPool)
	assert.Error(t, err)
}
