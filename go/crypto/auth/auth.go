// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://www.oauth.com/oauth2-servers/client-registration/client-id-secret/
package auth

import (
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/searKing/golang/go/crypto/rand"
)

const (
	DefaultSizeClientID      = 16
	DefaultSizeClientKey     = 16
	DefaultSizeClientSecret  = 32
	DefaultSizeAuthorizeCode = 16
	DefaultSizeUUID          = 64
)

// Here are some examples of client IDs from services that support OAuth 2.0:
// Foursquare: ZYDPLLBWSK3MVQJSIYHB1OR2JXCY0X2C5UJ2QAR2MAAIT5Q
// Github: 6779ef20e75817b79602
// Google: 292085223830.apps.googleusercontent.com
// Instagram: f2a1ed52710d4533bde25be6da03b6e3
// SoundCloud: 269d98e4922fb3895e9ae2108cbb5064
// Windows Live: 00000000400ECB04
func ClientID() string {
	return ClientIDWithSize(DefaultSizeClientID)
}
func ClientIDWithSize(len int) string {
	return rand.String(len)
}
func ClientSecret() string {
	return ClientSecretWithSize(DefaultSizeClientSecret)
}

func ClientSecretWithSize(len int) string {
	return rand.StringWithCharset(len, rand.CharsetHexadecimalDigits)
}

func ClientKey() string {
	return ClientKeyWithSize(DefaultSizeClientKey)
}

var ClientKeyWithSize = ClientSecretWithSize

func AuthorizeCode() string {
	var token []byte
	tokenN, err := uuid.NewRandom()
	if err != nil || len(tokenN) == 0 {
		token = rand.Bytes(DefaultSizeAuthorizeCode)
	} else {
		token = tokenN[:]
	}
	return base64.RawURLEncoding.EncodeToString(token)
}
func AuthorizeCodeWithSize(len int) string {
	return base64.RawURLEncoding.EncodeToString(rand.Bytes(len))
}

func UUID() string {
	return AuthorizeCode()
}
