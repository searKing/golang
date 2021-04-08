// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aes_test

// Code borrowed from https://github.com/golang/go/blob/master/src/crypto/cipher/example_test.go
import (
	"encoding/hex"
	"fmt"

	aes_ "github.com/searKing/golang/go/crypto/aes"
)

func ExampleGCMEncrypt() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	ciphertext, err := aes_.GCMEncrypt(key, plaintext, nil)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	fmt.Printf("%x\n", ciphertext)

	plaintext, err = aes_.GCMDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 0000000000000000000000000b2591cb60a33bdfeb61a0d35207a4196f1e5d1f6fe2b32198e09765ee28bd17d08567996caca50cd7049c07d1db15b5
	// decrypted:
	// exampleplaintext
}

func ExampleCBCEncrypt() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	ciphertext, err := aes_.CBCEncrypt(key, plaintext, nil)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	fmt.Printf("%x\n", ciphertext)

	plaintext, err = aes_.CBCDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 00000000000000000000000000000000f42512e1e4039213bd449ba47faa1b7408eac45dbf536e5016511f86035707c6
	// decrypted:
	// exampleplaintext
}

func ExampleCFBEncryptIV() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	ciphertext, err := aes_.CFBEncrypt(key, plaintext, nil)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	fmt.Printf("%x\n", ciphertext)

	plaintext, err = aes_.CFBDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 00000000000000000000000000000000d91399c43f6adaef3d909876e79904a91d84244a42a2a96b67bd9ae936651a2e
	// decrypted:
	// exampleplaintext
}

func ExampleCTREncrypt() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	ciphertext, err := aes_.CTREncrypt(key, plaintext, nil)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	fmt.Printf("%x\n", ciphertext)

	plaintext, err = aes_.CTRDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 00000000000000000000000000000000d91399c43f6adaef3d909876e79904a93b1144928bc98a4e78f38c23b706610a
	// decrypted:
	// exampleplaintext
}

func ExampleOFBEncrypt() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	ciphertext, err := aes_.OFBEncrypt(key, plaintext, nil)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	fmt.Printf("%x\n", ciphertext)

	plaintext, err = aes_.OFBDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 00000000000000000000000000000000d91399c43f6adaef3d909876e79904a9729f985422616b9f45b757ac9e7a879e
	// decrypted:
	// exampleplaintext
}
