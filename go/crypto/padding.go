// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crypto

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	rand_ "github.com/searKing/golang/go/crypto/rand"
)

// See https://en.wikipedia.org/wiki/Data_source_name
var ZeroPadding = ISO9791Method1Padding
var ZeroUnPadding = ISO9791Method1UnPadding
var BitPadding = ISO9791Method2Padding
var BitUnPadding = ISO9791Method2UnPadding

// ISO/IEC 7816-4:2005[9] is identical to the bit padding scheme, applied to a plain text of N bytes.
// This means in practice that the first byte is a mandatory byte valued '80' (Hexadecimal) followed, if needed,
// by 0 to N − 1 bytes set to '00', until the end of the block is reached. ISO/IEC 7816-4 itself
// is a communication standard for smart cards containing a file system,
// and in itself does not contain any cryptographic specifications.
var ISO78164Padding = ISO9791Method2Padding
var ISO78164UnPadding = ISO9791Method2UnPadding

// the data string D to be input to the MAC algorithm shall be right-padded with as few (possible none) '0' bits
// as necessary to obtain a data string whose length(in bits) is a positive integer multiple of n.
// NOTE
// 1 MAC algorithm using Padding Method 1 may be subject to trivial forgery attacks.
// 2 If the data string is empty, Padding Method1 specifies that it is right-padded with n '0' bits.
// See https://en.wikipedia.org/wiki/ISO/IEC_9797-1#Padding_method_1
func ISO9791Method1Padding(s []byte, block int) []byte {
	// If the data string is empty, Padding Method1 specifies that it is right-padded with n '0' bits.
	if block <= 0 || len(s) == 0 {
		return s
	}
	padding := block - (len(s) % block)
	if padding == 0 {
		return s
	}
	tmp := make([]byte, padding)
	return append(s, tmp...)
}

func ISO9791Method1UnPadding(s []byte, block int) []byte {
	var offset int
	for i := 0; i < block; i++ {
		j := len(s) - 1 - i
		if j < 0 || j >= len(s) {
			break
		}
		if s[j] != 0x00 {
			offset = j
			break
		}
	}
	return s[:offset+1]
}

// the data string D to be input to the MAC algorithm shall be right-padded with a single '1' bit.
// The resulting string shall then be right-padded with as few (possible none) '0' bits
// as necessary to obtain a data string whose length(in bits) is a positive integer multiple of n.
// See https://en.wikipedia.org/wiki/ISO/IEC_9797-1#Padding_method_2
func ISO9791Method2Padding(s []byte, block int) []byte {
	padding := block - (len(s) % block)
	if padding == 0 {
		padding = block
	}
	tmp := make([]byte, padding)
	tmp[0] = 0x80
	return append(s, tmp...)
}

func ISO9791Method2UnPadding(s []byte, block int) ([]byte, error) {
	var offset int
	var found bool
	for i := 0; i < block; i++ {
		j := len(s) - 1 - i
		if j < 0 || j >= len(s) {
			break
		}
		if s[j] == 0x00 {
			continue
		}
		if s[j] == 0x80 {
			offset = j
			found = true
			break
		}
		break
	}
	if !found {
		return nil, fmt.Errorf("missing %#x, unpad failed", 0x80)
	}
	return s[:offset], nil
}

// the data string D to be input to the MAC algorithm shall be right-padded with as few (possible none) '0' bits
// as necessary to obtain a data string whose length(in bits) is a positive integer multiple of n.
// The resulting string shall then be left-padded with a block L.
// The block L consists of the binary representation of the length(in bits) Ld of the unpadded data string D,
// left-padded with as few (possible none) '0' bits as necessary to obtain an n-bit block.
// The right-most bit of the block L corresponds to the least significant bit of the binary representation of Ld.
// NOTE
// 1 Padding Method 3 is not suitable for use in situations where the length of the data string is not available prior
// to the start of the MAC calculation.
// See https://en.wikipedia.org/wiki/ISO/IEC_9797-1#Padding_method_2
func ISO9791Method3Padding(s []byte, block int) []byte {
	// If the data string is empty, specifies that it is right-padded with n '0' bits.
	if block <= 0 {
		return s
	}
	unpaddedSizeInByte := len(s)

	s = ISO9791Method1Padding(s, block)
	padding := block
	tmp := make([]byte, padding)
	switch block {
	case 1:
		tmp[len(tmp)-1] = byte(unpaddedSizeInByte << 3) // append length in bits
	case 2, 3:
		binary.BigEndian.PutUint16(tmp[len(tmp)-2:], uint16(unpaddedSizeInByte<<3)) // append length in bits
	case 4, 5, 6, 7:
		binary.BigEndian.PutUint32(tmp[len(tmp)-2:], uint32(unpaddedSizeInByte<<3)) // append length in bits
	default:
		binary.BigEndian.PutUint64(tmp[len(tmp)-8:], uint64(unpaddedSizeInByte<<3)) // append length in bits
	}
	return append(tmp, s...)
}

func ISO9791Method3UnPadding(s []byte, block int) ([]byte, error) {
	// If the data string is empty, specifies that it is right-padded with n '0' bits.
	if block <= 0 {
		return s, nil
	}
	var unpadded int
	switch block {
	case 1:
		unpadded = int(s[block-1] >> 3)
	case 2, 3:
		unpadded = int(binary.BigEndian.Uint16(s[block-2:]) >> 3)
	case 4, 5, 6, 7:
		unpadded = int(binary.BigEndian.Uint32(s[block-4:]) >> 3)
	default:
		unpadded = int(binary.BigEndian.Uint64(s[block-8:]) >> 3)
	}
	if block+unpadded > len(s) {
		return nil, fmt.Errorf("malformed unpadded length %d, unpad failed", unpadded)
	}
	return s[block : block+unpadded], nil
}

// In ANSI X9.23, between 1 and 8 bytes are always added as padding.
// The block is padded with random bytes (although many implementations use 00) and
// the last byte of the block is set to the number of bytes added.[6]
// Example: In the following example the block size is 8 bytes,
// and padding is required for 4 bytes (in hexadecimal format)
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 00 00 00 04 |
// See https://en.wikipedia.org/wiki/Padding_(cryptography)#ANSI_X9.23
func ANSIX923Padding(s []byte, block int) []byte {
	// If the data string is empty, specifies that it is right-padded with n '0' bits.
	if block <= 0 {
		return s
	}

	padding := block - (len(s) % block)
	if padding == 0 {
		padding = block
	}
	tmp := make([]byte, padding)
	tmp[len(tmp)-1] = byte(padding) // append length in bytes
	return append(s, tmp...)
}

func ANSIX923UnPadding(s []byte, block int) ([]byte, error) {
	// If the data string is empty, specifies that it is right-padded with n '0' bits.
	if block <= 0 || len(s) <= 0 {
		return s, nil
	}
	offset := len(s) - int(s[len(s)-1])
	if offset > len(s) || offset < 0 {
		return nil, fmt.Errorf("malformed padding length %d, unpad failed", int(s[len(s)-1]))
	}
	return s[:offset], nil
}

// ISO 10126 (withdrawn, 2007[7][8]) specifies that the padding should be done
// at the end of that last block with random bytes, and the padding boundary should be specified by the last byte.
// Example: In the following example the block size is 8 bytes and padding is required for 4 bytes
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 81 A6 23 04 |
// See https://en.wikipedia.org/wiki/Padding_(cryptography)#ISO_10126
func IOS10126Padding(s []byte, block int) []byte {
	// If the data string is empty, specifies that it is right-padded with n '0' bits.
	if block <= 0 {
		return s
	}

	padding := block - (len(s) % block)
	if padding == 0 {
		padding = block
	}
	tmp, err := rand_.BytesCrypto(padding)
	if err != nil {
		tmp = make([]byte, padding)
		for i := 0; i < padding; i++ {
			tmp[i] = byte(rand.Uint32() & 0xFF)
		}
	}

	tmp[len(tmp)-1] = byte(padding) // append length in bytes

	return append(s, tmp...)
}

func IOS10126UnPadding(s []byte, block int) ([]byte, error) {
	// If the data string is empty, specifies that it is right-padded with n '0' bits.
	if block <= 0 || len(s) == 0 {
		return s, nil
	}
	offset := len(s) - int(s[len(s)-1])
	if offset > len(s) || offset < 0 {
		return nil, fmt.Errorf("malformed padding length %d, unpad failed", int(s[len(s)-1]))
	}
	return s[:offset], nil
}

// PKCS#5 padding is identical to PKCS#7 padding,
// except that it has only been defined for block ciphers that use a 64-bit (8-byte) block size.
// In practice the two can be used interchangeably.
func PKCS5Padding(s []byte) []byte {
	return PKCS7Padding(s, 8)
}

func PKCS5UnPadding(s []byte) ([]byte, error) {
	return PKCS7UnPadding(s, 8)
}

// Padding is in whole bytes.
// The value of each added byte is the number of bytes that are added, i.e. N bytes, each of value N are added.
// The number of bytes added will depend on the block boundary to which the message needs to be extended.
// The padding will be one of:
// 01
// 02 02
// 03 03 03
// 04 04 04 04
// 05 05 05 05 05
// 06 06 06 06 06 06
// etc.
// This padding method (as well as the previous two) is well-defined if and only if N is less than 256.
// Example: In the following example the block size is 8 bytes and padding is required for 4 bytes
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 04 04 04 04 |
// If the length of the original data is an integer multiple of the block size B,
// then an extra block of bytes with value B is added.
// This is necessary so the deciphering algorithm can determine with certainty
// whether the last byte of the last block is a pad byte indicating the number of padding bytes added
// or part of the plaintext message. Consider a plaintext message
// that is an integer multiple of B bytes with the last byte of plaintext being 01.
// With no additional information, the deciphering algorithm will not be able to
// determine whether the last byte is a plaintext byte or a pad byte.
// However, by adding B bytes each of value B after the 01 plaintext byte,
// the deciphering algorithm can always treat the last byte as a pad byte and strip
// the appropriate number of pad bytes off the end of the ciphertext;
// said number of bytes to be stripped based on the value of the last byte.
// See https://tools.ietf.org/html/rfc5652#section-6.3
func PKCS7Padding(s []byte, block int) []byte {
	// If the data string is empty, specifies that it is right-padded with n '0' bits.
	if block <= 0 {
		return s
	}

	padding := block - (len(s) % block)
	if padding == 0 {
		padding = block
	}
	tmp := make([]byte, padding)
	for i := 0; i < padding; i++ {
		tmp[i] = byte(padding) // append length in bytes
	}
	return append(s, tmp...)
}

func PKCS7UnPadding(s []byte, block int) ([]byte, error) {
	// If the data string is empty, specifies that it is right-padded with n '0' bits.
	if block <= 0 || len(s) == 0 {
		return s, nil
	}
	offset := len(s) - int(s[len(s)-1])
	if offset > len(s) || offset < 0 {
		return nil, fmt.Errorf("malformed padding length %d, unpad failed", int(s[len(s)-1]))
	}

	return s[:offset], nil
}
