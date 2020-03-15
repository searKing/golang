// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

const (
	// JOSE Header
	HeaderParameterType        = "type" // RFC 7231, 5.1
	HeaderParameterContentType = "cty"  // RFC 7231, 5.2

	// Replicating Claims as Header Parameters
	HeaderParameterIssuer         = ClaimNameIssuer         // RFC 7231, 5.3 RFC 7519, 10.4.1
	HeaderParameterSubject        = ClaimNameSubject        // RFC 7231, 5.3 RFC 7519, 10.4.1
	HeaderParameterAudience       = ClaimNameAudience       // RFC 7231, 5.3 RFC 7519, 10.4.1
	HeaderParameterExpirationTime = ClaimNameExpirationTime // RFC 7231, 5.3
	HeaderParameterNotBefore      = ClaimNameNotBefore      // RFC 7231, 5.3
	HeaderParameterIssuedAt       = ClaimNameIssuedAt       // RFC 7231, 5.3
	HeaderParameterJWTID          = ClaimNameJWTID          // RFC 7231, 5.3
)

var headerParameterText = map[string]string{
	HeaderParameterType:        "type",
	HeaderParameterContentType: "Content Type",

	HeaderParameterIssuer:         ClaimNamesText(ClaimNameIssuer),
	HeaderParameterSubject:        ClaimNamesText(ClaimNameSubject),
	HeaderParameterAudience:       ClaimNamesText(ClaimNameAudience),
	HeaderParameterExpirationTime: ClaimNamesText(ClaimNameExpirationTime),
	HeaderParameterNotBefore:      ClaimNamesText(ClaimNameNotBefore),
	HeaderParameterIssuedAt:       ClaimNamesText(ClaimNameIssuedAt),
	HeaderParameterJWTID:          ClaimNamesText(ClaimNameJWTID),
}

// HeaderParameterText returns a text for the Claim Names. It returns the empty
// string if the code is "".
func HeaderParameterText(param string) string {
	return headerParameterText[param]
}
