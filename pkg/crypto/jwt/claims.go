// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

import (
	"time"

	"github.com/google/uuid"
)

// See: https://tools.ietf.org/html/rfc7519
const (
	// Registered Claim Names
	ClaimNameIssuer         = "iss" // RFC 7519, 4.1.1
	ClaimNameSubject        = "sub" // RFC 7519, 4.1.2
	ClaimNameAudience       = "aud" // RFC 7519, 4.1.3
	ClaimNameExpirationTime = "exp" // RFC 7519, 4.1.4
	ClaimNameNotBefore      = "nbf" // RFC 7519, 4.1.5
	ClaimNameIssuedAt       = "iat" // RFC 7519, 4.1.6
	ClaimNameJWTID          = "jti" // RFC 7519, 4.1.7

	claimNameScope = "scope"
)

var claimNamesText = map[string]string{
	ClaimNameIssuer:         "Issuer",
	ClaimNameSubject:        "Subject",
	ClaimNameAudience:       "Audience",
	ClaimNameExpirationTime: "Expiration Time",
	ClaimNameNotBefore:      "Not Before",
	ClaimNameIssuedAt:       "Issued At",
	ClaimNameJWTID:          "JWT ID",
}

// ClaimNamesText returns a text for the Claim Names. It returns the empty
// string if the code is "".
func ClaimNamesText(name string) string {
	if name, ok := claimNamesText[name]; ok {
		return name
	}
	return name
}

type RegisteredClaims struct {
	Subject   string    `json:"sub"`
	Issuer    string    `json:"iss"`
	Audience  []string  `json:"aud"`
	JWTID     string    `json:"jti"`
	IssuedAt  time.Time `json:"iat"`
	NotBefore time.Time `json:"nbf"`
	ExpiresAt time.Time `json:"exp"`
}

func (c *RegisteredClaims) SetDefaults() {
	if c.JWTID == "" {
		c.JWTID = uuid.New().String()
	}
}

type ClaimsContainer interface {
	// With returns a copy of itself with expiresAt, scope, audience set to the given values.
	With(expiry time.Time, scope, audience []string) ClaimsContainer

	// WithDefaults returns a copy of itself with issuedAt and issuer set to the given default values. If those
	// values are already set in the claims, they will not be updated.
	WithDefaults(iat time.Time, issuer string) ClaimsContainer

	// ToMapClaims returns the claims as a github.com/dgrijalva/jwt-go.MapClaims type.
	ToMapClaims() jwt.MapClaims
}

// Claims represent a token's claims.
type Claims struct {
	RegisteredClaims
	Scope []string `json:"scope"`
	Extra map[string]interface{}
}

func (c *Claims) With(expiry time.Time, scope, audience []string) ClaimsContainer {
	c.ExpiresAt = expiry
	c.Scope = scope
	c.Audience = audience
	return c
}

func (c *Claims) WithDefaults(iat time.Time, issuer string) ClaimsContainer {
	if c.IssuedAt.IsZero() {
		c.IssuedAt = iat
	}

	if c.Issuer == "" {
		c.Issuer = issuer
	}
	return c
}

// ToMap will transform the headers to a map structure
func (c *Claims) ToMap() map[string]interface{} {
	var ret = Copy(c.Extra)

	ret[ClaimNameJWTID] = c.JWTID
	if c.JWTID == "" {
		ret[ClaimNameJWTID] = uuid.New()
	}

	ret[ClaimNameSubject] = c.Subject
	ret[ClaimNameIssuer] = c.Issuer
	ret[ClaimNameAudience] = c.Audience

	if !c.IssuedAt.IsZero() {
		ret[ClaimNameIssuedAt] = float64(c.IssuedAt.Unix()) // jwt-go does not support int64 as datatype
	}

	if !c.NotBefore.IsZero() {
		ret[ClaimNameNotBefore] = float64(c.NotBefore.Unix()) // jwt-go does not support int64 as datatype
	}

	ret[ClaimNameExpirationTime] = float64(c.ExpiresAt.Unix()) // jwt-go does not support int64 as datatype

	if c.Scope != nil {
		ret[claimNameScope] = c.Scope
	}

	return ret
}

// FromMap will set the claims based on a mapping
func (c *Claims) FromMap(m map[string]interface{}) {
	c.Extra = make(map[string]interface{})
	for k, v := range m {
		switch k {
		case ClaimNameJWTID:
			if s, ok := v.(string); ok {
				c.JWTID = s
			}
		case ClaimNameSubject:
			if s, ok := v.(string); ok {
				c.Subject = s
			}
		case ClaimNameIssuer:
			if s, ok := v.(string); ok {
				c.Issuer = s
			}
		case ClaimNameAudience:
			if s, ok := v.(string); ok {
				c.Audience = []string{s}
			} else if s, ok := v.([]string); ok {
				c.Audience = s
			}
		case ClaimNameIssuedAt:
			switch v.(type) {
			case float64:
				c.IssuedAt = time.Unix(int64(v.(float64)), 0).UTC()
			case int64:
				c.IssuedAt = time.Unix(v.(int64), 0).UTC()
			}
		case ClaimNameNotBefore:
			switch v.(type) {
			case float64:
				c.NotBefore = time.Unix(int64(v.(float64)), 0).UTC()
			case int64:
				c.NotBefore = time.Unix(v.(int64), 0).UTC()
			}
		case ClaimNameExpirationTime:
			switch v.(type) {
			case float64:
				c.ExpiresAt = time.Unix(int64(v.(float64)), 0).UTC()
			case int64:
				c.ExpiresAt = time.Unix(v.(int64), 0).UTC()
			}
		case claimNameScope:
			switch v.(type) {
			case []string:
				c.Scope = v.([]string)
			case []interface{}:
				c.Scope = make([]string, len(v.([]interface{})))
				for i, vi := range v.([]interface{}) {
					if s, ok := vi.(string); ok {
						c.Scope[i] = s
					}
				}
			}
		default:
			c.Extra[k] = v
		}
	}
}

// Add will add a key-value pair to the extra field
func (c *Claims) Add(key string, value interface{}) {
	if c.Extra == nil {
		c.Extra = make(map[string]interface{})
	}
	c.Extra[key] = value
}

// Get will get a value from the extra field based on a given key
func (c Claims) Get(key string) interface{} {
	return c.ToMap()[key]
}

// ToMapClaims will return a jwt-go MapClaims representation
func (c Claims) ToMapClaims() jwt.MapClaims {
	return c.ToMap()
}

// FromMapClaims will populate claims from a jwt-go MapClaims representation
func (c *Claims) FromMapClaims(mc jwt.MapClaims) {
	c.FromMap(mc)
}

// Copy will copy all elements in a map and return a new representational map
func Copy(elements map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{}, len(elements))
	for k, v := range elements {
		result[k] = v
	}

	return result
}
