// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

// https://jwt.io/
const (
	SigningMethodNone  = "none"
	SigningMethodHS256 = "HS256" // HS256: HMAC using SHA-256
	SigningMethodHS384 = "HS384" // HS384: HMAC using SHA-384
	SigningMethodHS512 = "HS512" // HS512: HMAC using SHA-512
	SigningMethodRS256 = "RS256" // RS256: RSASSA-PKCS-v1_5 using SHA-256
	SigningMethodRS384 = "RS384" // RS384: RSASSA-PKCS-v1_5 using SHA-384
	SigningMethodRS512 = "RS512" // RS512: RSASSA-PKCS-v1_5 using SHA-512
	SigningMethodPS256 = "PS256" // PS256: RSASSA-PSS using SHA-256 and MGF1 with SHA-256
	SigningMethodPS384 = "PS384" // PS384: RSASSA-PSS using SHA-384 and MGF1 with SHA-384
	SigningMethodPS512 = "PS512" // PS512: RSASSA-PSS using SHA-512 and MGF1 with SHA-512
	SigningMethodES256 = "ES256" // ES256: ECDSA using P-256 and SHA-256
	SigningMethodES384 = "ES384" // ES384: ECDSA using P-384 and SHA-384
	SigningMethodES512 = "ES512" // ES512: ECDSA using P-521 and SHA-512
)
