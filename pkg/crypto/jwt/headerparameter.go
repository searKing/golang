package jwt

const (
	// JOSE Header
	HeaderParameterType        = "type" // RFC 7231, 5.1
	HeaderParameterContentType = "cty"  // RFC 7231, 5.2

	// Replicating Claims as Header Parameters
	HeaderParameterIssuer         = ClaimsIssuer         // RFC 7231, 5.3 RFC 7519, 10.4.1
	HeaderParameterSubject        = ClaimsSubject        // RFC 7231, 5.3 RFC 7519, 10.4.1
	HeaderParameterAudience       = ClaimsAudience       // RFC 7231, 5.3 RFC 7519, 10.4.1
	HeaderParameterExpirationTime = ClaimsExpirationTime // RFC 7231, 5.3
	HeaderParameterNotBefore      = ClaimsNotBefore      // RFC 7231, 5.3
	HeaderParameterIssuedAt       = ClaimsIssuedAt       // RFC 7231, 5.3
	HeaderParameterJWTID          = ClaimsJWTID          // RFC 7231, 5.3
)

var headerParameterText = map[string]string{
	HeaderParameterType:        "type",
	HeaderParameterContentType: "Content Type",

	HeaderParameterIssuer:         ClaimsText(ClaimsIssuer),
	HeaderParameterSubject:        ClaimsText(ClaimsSubject),
	HeaderParameterAudience:       ClaimsText(ClaimsAudience),
	HeaderParameterExpirationTime: ClaimsText(ClaimsExpirationTime),
	HeaderParameterNotBefore:      ClaimsText(ClaimsNotBefore),
	HeaderParameterIssuedAt:       ClaimsText(ClaimsIssuedAt),
	HeaderParameterJWTID:          ClaimsText(ClaimsJWTID),
}

// HeaderParameterText returns a text for the Claim Names. It returns the empty
// string if the code is "".
func HeaderParameterText(param string) string {
	return headerParameterText[param]
}
