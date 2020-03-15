// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import (
	"encoding/base64"
)

// Numerical elements
const (
	CharsetBinaryDigits      = "01"
	CharsetOctalDigits       = "01234567"
	CharsetDecimalDigits     = "0123456789"
	CharsetHexadecimalDigits = "0123456789aAbBcCdDeEfF"
	CharsetSign              = "+-"
	CharsetPeriod            = "."
	CharsetExponent          = "eEp"
	CharsetAlpha             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetAlphaNum          = CharsetDecimalDigits + CharsetAlpha
	CharsetBase64            = CharsetAlphaNum + "+/"

	// StdEncoding is the standard base64 encoding, as defined in
	// RFC 4648.
	CharsetBase64StdEncoding = CharsetBase64RawStdEncoding + string(base64.StdPadding)

	// URLEncoding is the alternate base64 encoding defined in RFC 4648.
	// It is typically used in URLs and file names.
	CharsetBase64URLEncoding = CharsetBase64RawURLEncoding + string(base64.StdPadding)

	// RawStdEncoding is the standard raw, unpadded base64 encoding,
	// as defined in RFC 4648 section 3.2.
	// This is the same as StdEncoding but omits padding characters.
	CharsetBase64RawStdEncoding = CharsetAlphaNum + "+/"

	// RawURLEncoding is the unpadded alternate base64 encoding defined in RFC 4648.
	// It is typically used in URLs and file names.
	// This is the same as URLEncoding but omits padding characters.
	CharsetBase64RawURLEncoding = CharsetAlphaNum + "-_"
)

// take in a character set and a length and will generate a random string using that character set.
func StringWithCharset(len int, charset string) string {
	s, err := StringCryptoWithCharset(int64(len), charset)
	if err == nil {
		return s
	}
	return StringMathWithCharset(len, charset)
}

// only take in a length, and will use a default characters set to generate a random string
func String(len int) string {
	return StringWithCharset(len, CharsetAlphaNum)
}
