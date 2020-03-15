// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/searKing/golang/go/crypto/auth"
	"github.com/searKing/golang/go/error/exception"
)

type AuthKey struct {
	alg string

	// Private key
	privKey crypto.PrivateKey

	// Public key
	pubKey crypto.PublicKey

	// Secret key used for signing. Required.
	symmetricKey []byte
}

func NewAuthKeyFromRandom(alg string) (*AuthKey, error) {
	if alg == "" {
		alg = SigningMethodNone
	}
	var privKey crypto.PrivateKey
	var pubKey crypto.PublicKey
	switch alg {
	case SigningMethodNone:
		return &AuthKey{alg: alg}, nil
	case SigningMethodHS256, SigningMethodHS384, SigningMethodHS512:
		return &AuthKey{alg: alg, symmetricKey: []byte(auth.ClientKeyWithSize(2048))}, nil
	case SigningMethodRS256, SigningMethodPS256, SigningMethodRS384, SigningMethodPS384, SigningMethodRS512, SigningMethodPS512:
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		privKey = priv
		pubKey = priv.Public()
	case SigningMethodES256:
		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		privKey = priv
		pubKey = priv.Public()
	case SigningMethodES384:
		priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, err
		}
		privKey = priv
		pubKey = priv.Public()
	case SigningMethodES512:
		priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, err
		}
		privKey = priv
		pubKey = priv.Public()
	default:
		return nil, fmt.Errorf("unsupport jwt.alg [%s]", alg)
	}
	authKey := &AuthKey{
		alg: alg,
	}
	authKey.privKey = privKey
	authKey.pubKey = pubKey
	return authKey, nil
}

// SymmetricKey : privateKey
// else: privKey publicKey
func NewAuthKey(alg string, privateKey []byte, publicKey []byte, password ...string) (*AuthKey, error) {
	if alg == "" {
		alg = SigningMethodNone
	}
	authKey := &AuthKey{
		alg: alg,
	}
	switch alg {
	case SigningMethodNone:
		return &AuthKey{alg: alg}, nil
	case SigningMethodHS256, SigningMethodHS384, SigningMethodHS512:
		if len(privateKey) == 0 {
			return authKey, exception.NewIllegalArgumentException1("privateKey is missing")
		}
		if err := authKey.setPrivateKey(privateKey, password...); err != nil {
			return nil, err
		}
	case SigningMethodRS256, SigningMethodPS256,
		SigningMethodRS384, SigningMethodPS384,
		SigningMethodRS512, SigningMethodPS512,
		SigningMethodES256, SigningMethodES384, SigningMethodES512:
		if len(privateKey) == 0 {
			return authKey, exception.NewIllegalArgumentException1("privateKey is missing")
		}
		if err := authKey.setPrivateKey(privateKey, password...); err != nil {
			return nil, err
		}
		if len(publicKey) == 0 {
			return authKey, nil
		}
		if err := authKey.setPublicKey(publicKey); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupport jwt.alg [%s]", alg)
	}
	return authKey, nil
}

func NewAuthKeyFromFile(alg string, privateKeyFile string, publicKeyFile string, password ...string) (*AuthKey, error) {
	authKey := &AuthKey{
		alg: alg,
	}
	if privateKeyFile == "" {
		return nil, exception.NewIllegalArgumentException1("key file path is empty")
	}
	if err := authKey.setPrivateKeyFromFile(privateKeyFile, password...); err != nil {
		return nil, err
	}
	if publicKeyFile == "" {
		return authKey, nil
	}
	if err := authKey.setPublicKeyFromFile(publicKeyFile); err != nil {
		return nil, err
	}
	return authKey, nil
}

func (a *AuthKey) setPrivateKeyFromFile(keyFile string, passwords ...string) error {
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}
	return a.setPrivateKey(keyData, passwords...)
}

func (a *AuthKey) setPublicKeyFromFile(keyFile string) error {
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}
	return a.setPublicKey(keyData)
}

func (a *AuthKey) setPrivateKey(keyData []byte, passwords ...string) error {
	switch a.alg {
	case SigningMethodHS256, SigningMethodHS384, SigningMethodHS512:
		a.symmetricKey = keyData
		return nil
	case SigningMethodRS256, SigningMethodRS384, SigningMethodRS512,
		SigningMethodPS256, SigningMethodPS384, SigningMethodPS512:
		passwordsLen := len(passwords)
		var privKey *rsa.PrivateKey
		if passwordsLen == 0 {
			priv, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
			if err != nil {
				return err
			}
			privKey = priv
		} else {
			priv, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(keyData, passwords[0])
			if err != nil {
				return err
			}
			privKey = priv
		}
		a.privKey = privKey
		a.pubKey = privKey.Public()
		return nil

	case SigningMethodES256, SigningMethodES384, SigningMethodES512:
		privKey, err := jwt.ParseECPrivateKeyFromPEM(keyData)
		if err != nil {
			return err
		}
		a.privKey = privKey
		a.pubKey = privKey.Public()
		return nil
	}
	return exception.NewIllegalArgumentException1(fmt.Sprintf("unsupport jwt.alg [%s]", a.alg))
}

func (a *AuthKey) setPublicKey(keyData []byte) error {
	switch a.alg {
	case SigningMethodHS256, SigningMethodHS384, SigningMethodHS512:
		a.symmetricKey = keyData
		return nil
	case SigningMethodRS256, SigningMethodRS384, SigningMethodRS512,
		SigningMethodPS256, SigningMethodPS384, SigningMethodPS512:
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
		if err != nil {
			return err
		}
		a.pubKey = pubKey
		return nil
	case SigningMethodES256, SigningMethodES384, SigningMethodES512:
		pubKey, err := jwt.ParseECPublicKeyFromPEM(keyData)
		if err != nil {
			return err
		}
		a.pubKey = pubKey
		return nil
	}
	return exception.NewIllegalArgumentException1(fmt.Sprintf("unsupport jwt.alg [%s]", a.alg))
}

func (a *AuthKey) GetSignedKey(token *jwt.Token) (interface{}, error) {
	if token != nil && jwt.GetSigningMethod(a.alg) != token.Method {
		return nil, errors.New("invalid signing method")
	}
	if a.IsSymmetricKey() {
		return a.symmetricKey, nil
	}
	return a.privKey, nil
}

func (a *AuthKey) GetVerifiedKey(token *jwt.Token) (interface{}, error) {
	if token != nil && jwt.GetSigningMethod(a.alg) != token.Method {
		return nil, errors.New("invalid signing method")
	}
	if a.IsSymmetricKey() {
		return a.symmetricKey, nil
	}
	return a.pubKey, nil
}

func (a *AuthKey) GetSignedMethod() jwt.SigningMethod {
	return jwt.GetSigningMethod(a.alg)
}

func (a *AuthKey) IsSymmetricKey() bool {
	switch a.alg {
	case SigningMethodRS256, SigningMethodRS384, SigningMethodRS512:
		return false
	}
	return true
}

func (a *AuthKey) IsRSAKey() bool {
	switch a.alg {
	case SigningMethodRS256, SigningMethodRS384, SigningMethodRS512,
		SigningMethodPS256, SigningMethodPS384, SigningMethodPS512:
		return true
	}
	return false
}

func (a *AuthKey) IsECMAKey() bool {
	switch a.alg {
	case SigningMethodES256, SigningMethodES384, SigningMethodES512:
		return true
	}
	return false
}
