// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiple_prefix

import (
	"strings"
	"unicode"

	strings_ "github.com/searKing/golang/go/strings"
)

func BinaryFormatInt(number int, precision int) string {
	return BinaryFormatInt64(int64(number), precision)
}

func BinaryFormatUint(number uint, precision int) string {
	return BinaryFormatUint64(uint64(number), precision)
}

func BinaryFormatInt64(number int64, precision int) string {
	return BinaryMultiplePrefixTODO.Copy().SetInt64(number).FormatInt64(number, precision)
}

func BinaryFormatUint64(number uint64, precision int) string {
	return BinaryMultiplePrefixTODO.Copy().SetUint64(number).FormatUint64(number, precision)
}

func BinaryFormatFloat(number float64, precision int) string {
	return BinaryMultiplePrefixTODO.Copy().SetFloat64(number).FormatFloat(number, precision)
}

// SplitBinary splits s into number, multiple_prefix and unparsed strings
func SplitBinary(s string) (number string, prefix *BinaryMultiplePrefix, unparsed string) {
	splits := strings_.SplitPrefixNumber(s)
	if len(splits) < 2 {
		return "", nil, unparsed
	}
	number = splits[0]
	unparsed = splits[1]
	// trim any space between numbers and symbols
	unparsed = strings.TrimLeftFunc(splits[1], unicode.IsSpace)

	for _, prefix := range binaryPositiveMultiplePrefixes {
		if strings.HasPrefix(unparsed, prefix.Symbol()) {
			return number, prefix.Copy(), strings.TrimPrefix(unparsed, prefix.Symbol())
		}
	}
	for _, prefix := range binaryNegativeMultiplePrefixes {
		if strings.HasPrefix(unparsed, prefix.Symbol()) {
			return number, prefix.Copy(), strings.TrimPrefix(unparsed, prefix.Symbol())
		}
	}
	for _, prefix := range binaryZeroMultiplePrefixes {
		if strings.HasPrefix(unparsed, prefix.Symbol()) {
			return number, prefix.Copy(), strings.TrimPrefix(unparsed, prefix.Symbol())
		}
	}
	return number, nil, unparsed
}
