// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cryptopasta

import (
	"encoding/base64"
	"fmt"

	"github.com/gtank/cryptopasta"
)

// copy & paste-friendly golang crypto
func fixedKey(key []byte) *[32]byte {
	var result [32]byte
	copy(result[:], key[:32])
	return &result
}

func Encrypt(plaintext []byte, keys ...[]byte) (ciphertext string, err error) {
	if len(keys) == 0 {
		return "", fmt.Errorf("at least one encryption key must be defined but none were")
	}
	for _, key := range keys {
		if ciphertext, err = encrypt(plaintext, key); err == nil {
			return ciphertext, nil
		}
	}

	return "", err
}

func Decrypt(ciphertext string, keys ...[]byte) (p []byte, err error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("at least one decryption key must be defined but none were")
	}

	for _, key := range keys {
		if p, err = decrypt(ciphertext, key); err == nil {
			return p, nil
		}
	}

	return nil, err
}

func encrypt(plaintext []byte, key []byte) (string, error) {
	if len(key) != 32 {
		return "", fmt.Errorf("key must be exactly 32 long bytes, got %d bytes", len(key))
	}

	ciphertext, err := cryptopasta.Encrypt(plaintext, fixedKey(key))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(ciphertext string, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be exactly 32 long bytes, got %d bytes", len(key))
	}

	raw, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	plaintext, err := cryptopasta.Decrypt(raw, fixedKey(key))
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
