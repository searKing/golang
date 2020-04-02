// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package md5

import (
	"crypto/md5"
	"io"
	"log"
	"os"
)

func MySelf() ([]byte, error) {
	f, err := os.Open(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func SumBytes(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func SumString(b string) string {
	return string(SumBytes([]byte(b)))
}

func SumReader(r io.Reader) ([]byte, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func SumFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	return SumReader(f)
}
