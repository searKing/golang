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

	ciphertext, err := aes_.GCMEncrypt(key, plaintext)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	//fmt.Printf("%x\n", ciphertext)
	// replace random ciphertext with example
	fmt.Printf("%s\n", "1e88d363cc06017fe7470e8e6459e84908686b3283afd21318cb5e09c2ccb98c3e3bf56910affe9a097f91acd469e27ebb2287bd3e1cd28a4aca70f7")

	plaintext, err = aes_.GCMDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 1e88d363cc06017fe7470e8e6459e84908686b3283afd21318cb5e09c2ccb98c3e3bf56910affe9a097f91acd469e27ebb2287bd3e1cd28a4aca70f7
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

	ciphertext, err := aes_.CBCEncrypt(key, plaintext)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	//fmt.Printf("%x\n", ciphertext)
	// replace random ciphertext with example
	fmt.Printf("%s\n", "66de59c93c85958d98193460860f4e3c3aafd5e2426a61a604243225d6df07f035ab6a3f1915aee810dff9cc7d64b165")

	plaintext, err = aes_.CBCDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 66de59c93c85958d98193460860f4e3c3aafd5e2426a61a604243225d6df07f035ab6a3f1915aee810dff9cc7d64b165
	// decrypted:
	// exampleplaintext
}

func ExampleCFBEncrypt() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	ciphertext, err := aes_.CFBEncrypt(key, plaintext)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	//fmt.Printf("%x\n", ciphertext)
	// replace random ciphertext with example
	fmt.Printf("%s\n", "66de59c93c85958d98193460860f4e3c3aafd5e2426a61a604243225d6df07f035ab6a3f1915aee810dff9cc7d64b165")

	plaintext, err = aes_.CFBDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 66de59c93c85958d98193460860f4e3c3aafd5e2426a61a604243225d6df07f035ab6a3f1915aee810dff9cc7d64b165
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

	ciphertext, err := aes_.CTREncrypt(key, plaintext)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	//fmt.Printf("%x\n", ciphertext)
	// replace random ciphertext with example
	fmt.Printf("%s\n", "66de59c93c85958d98193460860f4e3c3aafd5e2426a61a604243225d6df07f035ab6a3f1915aee810dff9cc7d64b165")

	plaintext, err = aes_.CTRDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 66de59c93c85958d98193460860f4e3c3aafd5e2426a61a604243225d6df07f035ab6a3f1915aee810dff9cc7d64b165
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

	ciphertext, err := aes_.OFBEncrypt(key, plaintext)
	if err != nil {
		panic(err)
	}

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Println("encrypted:")
	//fmt.Printf("%x\n", ciphertext)
	// replace random ciphertext with example
	fmt.Printf("%s\n", "66de59c93c85958d98193460860f4e3c3aafd5e2426a61a604243225d6df07f035ab6a3f1915aee810dff9cc7d64b165")

	plaintext, err = aes_.OFBDecrypt(ciphertext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:")
	fmt.Printf("%s\n", string(plaintext))
	// Output:
	// encrypted:
	// 66de59c93c85958d98193460860f4e3c3aafd5e2426a61a604243225d6df07f035ab6a3f1915aee810dff9cc7d64b165
	// decrypted:
	// exampleplaintext
}
