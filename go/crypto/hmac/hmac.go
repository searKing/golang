// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hmac

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"hash"
	"io"
	"log"
	"os"
)

func MySelf(hasher func() hash.Hash, key []byte) ([]byte, error) {
	f, err := os.Open(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := hmac.New(hasher, key)
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func SumBytes(hasher func() hash.Hash, key []byte, b []byte) []byte {
	h := hmac.New(hasher, key)
	h.Write(b)
	return h.Sum(nil)
}

func SumString(hasher func() hash.Hash, key []byte, b string) string {
	return string(SumBytes(hasher, key, []byte(b)))
}

func SumHex(hasher func() hash.Hash, key []byte, b string) string {
	return hex.EncodeToString(SumBytes(hasher, key, []byte(b)))
}

func SumReader(hasher func() hash.Hash, key []byte, r io.Reader) ([]byte, error) {
	h := hmac.New(hasher, key)
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func SumFile(hasher func() hash.Hash, key []byte, name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return SumReader(hasher, key, f)
}

func SumReaderAt(hasher func() hash.Hash, key []byte, r io.ReaderAt, offset, size int64, buf []byte) ([]byte, error) {
	h := hmac.New(hasher, key)
	if len(buf) == 0 {
		buf = make([]byte, 1024)
	}

	var total int64
	for total < size {
		n, err := r.ReadAt(buf, offset)
		if err != nil {
			return nil, err
		}
		offset += int64(n)
		total += int64(n)
		var copied int64
		for copied < int64(n) {
			copiedCount, err := io.CopyN(h, bytes.NewReader(buf), int64(n))
			if err != nil {
				return nil, err
			}
			copied += copiedCount
		}
	}

	return h.Sum(nil), nil
}

// SumFileAt return ms5sum of data by offset and len is size
// buf is a buffer to read from file every time.
func SumFileAt(hasher func() hash.Hash, key []byte, name string, offset, size int64, buf []byte) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return SumReaderAt(hasher, key, f, offset, size, buf)
}
