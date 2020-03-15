// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	jwt_ "github.com/ory/fosite/token/jwt"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

// HS256JWTStrategy is responsible for generating and validating JWT challenges
type HS256JWTStrategy struct {
	Key []byte
}

// Generate generates a new authorize code or returns an error. set secret
func (j *HS256JWTStrategy) Generate(ctx context.Context, claims jwt.Claims, header jwt_.Mapper) (string, string, error) {
	if header == nil || claims == nil {
		return "", "", errors.New("Either claims or header is nil.")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header = assign(token.Header, header.ToMap())

	var sig, sstr string
	var err error
	if sstr, err = token.SigningString(); err != nil {
		return "", "", errors.WithStack(err)
	}

	if sig, err = token.Method.Sign(sstr, j.Key); err != nil {
		return "", "", errors.WithStack(err)
	}

	return fmt.Sprintf("%s.%s", sstr, sig), sig, nil
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (j *HS256JWTStrategy) Validate(ctx context.Context, token string) (string, error) {
	if _, err := j.Decode(ctx, token); err != nil {
		return "", errors.WithStack(err)
	}

	return j.GetSignature(ctx, token)
}

// Decode will decode a JWT token
func (j *HS256JWTStrategy) Decode(ctx context.Context, token string) (*jwt.Token, error) {
	// Parse the token.
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return j.Key, nil
	})

	if err != nil {
		return parsedToken, errors.WithStack(err)
	} else if !parsedToken.Valid {
		return parsedToken, errors.WithStack(fosite.ErrInactiveToken)
	}

	return parsedToken, err
}

// GetSignature will return the signature of a token
func (j *HS256JWTStrategy) GetSignature(ctx context.Context, token string) (string, error) {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return "", errors.New("Header, body and signature must all be set")
	}
	return split[2], nil
}

// Hash will return a given hash based on the byte input or an error upon fail
func (j *HS256JWTStrategy) Hash(ctx context.Context, in []byte) ([]byte, error) {
	// SigningMethodRS256
	hash := sha256.New()
	_, err := hash.Write(in)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return hash.Sum([]byte{}), nil
}

// GetSigningMethodLength will return the length of the signing method
func (j *HS256JWTStrategy) GetSigningMethodLength() int {
	return jwt.SigningMethodHS256.Hash.Size()
}

func assign(a, b map[string]interface{}) map[string]interface{} {
	for k, w := range b {
		if _, ok := a[k]; ok {
			continue
		}
		a[k] = w
	}
	return a
}
