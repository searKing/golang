package jwt

// See: https://tools.ietf.org/html/rfc7519
const (
	// Registered Claim Names
	ClaimsIssuer         = "iss" // RFC 7519, 4.1.1
	ClaimsSubject        = "sub" // RFC 7519, 4.1.2
	ClaimsAudience       = "aud" // RFC 7519, 4.1.3
	ClaimsExpirationTime = "exp" // RFC 7519, 4.1.4
	ClaimsNotBefore      = "nbf" // RFC 7519, 4.1.5
	ClaimsIssuedAt       = "iat" // RFC 7519, 4.1.6
	ClaimsJWTID          = "jti" // RFC 7519, 4.1.7
)

var claimsText = map[string]string{
	ClaimsIssuer:         "Issuer",
	ClaimsSubject:        "Subject",
	ClaimsAudience:       "Audience",
	ClaimsExpirationTime: "Expiration Time",
	ClaimsNotBefore:      "Not Before",
	ClaimsIssuedAt:       "Issued At",
	ClaimsJWTID:          "JWT ID",
}

// ClaimsText returns a text for the Claim Names. It returns the empty
// string if the code is "".
func ClaimsText(name string) string {
	return claimsText[name]
}
