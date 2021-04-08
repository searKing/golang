// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crypto_test

import (
	"encoding/hex"
	"fmt"

	"github.com/searKing/golang/go/crypto"
)

var unpad = []byte{0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD}

const (
	blockSizeInByte = 8
)

func ExampleZeroPadding() {
	padded := crypto.ZeroPadding(unpad, blockSizeInByte)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded := crypto.ZeroUnPadding(padded, blockSizeInByte)
	fmt.Println("unpadding:")
	fmt.Println(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 00 00 00 00  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExampleBitPadding() {
	padded := crypto.BitPadding(unpad, blockSizeInByte)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded, err := crypto.BitUnPadding(padded, blockSizeInByte)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 80 00 00 00  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExamplePKCS5Padding() {
	padded := crypto.PKCS5Padding(unpad)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded, err := crypto.PKCS5UnPadding(padded)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 04 04 04 04  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExampleISO78164Padding() {
	padded := crypto.ISO78164Padding(unpad, blockSizeInByte)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded, err := crypto.ISO78164UnPadding(padded, blockSizeInByte)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 80 00 00 00  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExampleISO9791Method1Padding() {
	padded := crypto.ISO9791Method1Padding(unpad, blockSizeInByte)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded := crypto.ISO9791Method1UnPadding(padded, blockSizeInByte)
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 00 00 00 00  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExampleISO9791Method2Padding() {
	padded := crypto.ISO9791Method2Padding(unpad, blockSizeInByte)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded, err := crypto.ISO9791Method2UnPadding(padded, blockSizeInByte)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 80 00 00 00  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExampleISO9791Method3Padding() {
	padded := crypto.ISO9791Method3Padding(unpad, blockSizeInByte)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded, err := crypto.ISO9791Method3UnPadding(padded, blockSizeInByte)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  00 00 00 00 00 00 00 60  dd dd dd dd dd dd dd dd  |.......`........|
	// 00000010  dd dd dd dd 00 00 00 00                           |........|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExampleANSIX923Padding() {
	padded := crypto.ANSIX923Padding(unpad, blockSizeInByte)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded, err := crypto.ANSIX923UnPadding(padded, blockSizeInByte)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 00 00 00 04  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExampleIOS10126Padding() {
	padded := crypto.IOS10126Padding(unpad, blockSizeInByte)
	for i := len(unpad); i < len(padded)-1; i++ {
		padded[i] = 0
	}
	fmt.Println("padding, replace random with 0x00:")
	fmt.Print(hex.Dump(padded))
	decoded, err := crypto.IOS10126UnPadding(padded, blockSizeInByte)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding, replace random with 0x00:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 00 00 00 04  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}

func ExamplePKCS7Padding() {
	padded := crypto.PKCS7Padding(unpad, blockSizeInByte)
	fmt.Println("padding:")
	fmt.Print(hex.Dump(padded))
	decoded, err := crypto.PKCS7UnPadding(padded, blockSizeInByte)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("unpadding:")
	fmt.Print(hex.Dump(decoded))

	// Output:
	// padding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd 04 04 04 04  |................|
	// unpadding:
	// 00000000  dd dd dd dd dd dd dd dd  dd dd dd dd              |............|
}
