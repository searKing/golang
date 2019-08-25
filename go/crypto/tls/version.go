package tls

import "crypto/tls"

var Versions = []int{
	tls.VersionSSL30,
	tls.VersionTLS10,
	tls.VersionTLS11,
	tls.VersionTLS12,
	tls.VersionTLS13,
}
