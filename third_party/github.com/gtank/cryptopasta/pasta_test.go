// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cryptopasta_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/searKing/golang/third_party/github.com/gtank/cryptopasta"
)

func TestEncryptDecryptGCM(t *testing.T) {
	randomKey := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, randomKey[:])
	if err != nil {
		t.Fatal(err)
	}

	gcmTests := []struct {
		plaintext []byte
		key       []byte
	}{
		{
			plaintext: []byte("Hello, world!"),
			key:       randomKey[:],
		},
	}

	for _, tt := range gcmTests {
		ciphertext, err := cryptopasta.Encrypt(tt.plaintext, tt.key)
		if err != nil {
			t.Fatal(err)
		}

		plaintext, err := cryptopasta.Decrypt(ciphertext, tt.key)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(plaintext, tt.plaintext) {
			t.Errorf("plaintexts don't match")
		}
		malformedCiphertext := string(ciphertext[0]^0xff) + ciphertext[1:]
		plaintext, err = cryptopasta.Decrypt(malformedCiphertext)
		if err == nil {
			t.Errorf("gcmOpen should not have worked, but did")
		}
	}
}

func BenchmarkAESGCM(b *testing.B) {
	randomKey := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, randomKey[:])
	if err != nil {
		b.Fatal(err)
	}

	data, err := ioutil.ReadFile("testdata/big")
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		cryptopasta.Encrypt(data, randomKey[:])
	}
}
