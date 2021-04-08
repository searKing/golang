// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/searKing/golang/go/crypto"
)

const (
	gcmBlockSize         = 16
	gcmTagSize           = 16
	gcmMinimumTagSize    = 12 // NIST SP 800-38D recommends tags with 12 or more bytes.
	gcmStandardNonceSize = 12
)

func GCMEncryptRandom(key, plaintext []byte) ([]byte, error) {
	var nonce = [gcmStandardNonceSize]byte{}
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}
	return CFBEncrypt(key, plaintext, nonce[:])
}

func GCMEncrypt(key, plaintext []byte, nonce []byte) ([]byte, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).

	paddingtext := crypto.PKCS7Padding(plaintext, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The nounce needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.

	ciphertext := make([]byte, gcmStandardNonceSize)
	if len(nonce) < gcmStandardNonceSize {
		copy(ciphertext[:gcmStandardNonceSize], nonce[:])
	} else {
		copy(ciphertext[:gcmStandardNonceSize], nonce[:gcmStandardNonceSize])
	}
	nonce = ciphertext[:gcmStandardNonceSize]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	sealedtext := aesgcm.Seal(nil, nonce, paddingtext, nil)

	ciphertext = append(nonce, sealedtext...)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	return ciphertext, nil
}

func GCMDecrypt(ciphertext, key []byte) ([]byte, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).

	if len(ciphertext) < gcmStandardNonceSize {
		panic("ciphertext too short")
	}
	nonce := ciphertext[:gcmStandardNonceSize]
	ciphertext = ciphertext[gcmStandardNonceSize:]

	// GCM mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	paddingtext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := crypto.PKCS7UnPadding(paddingtext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
