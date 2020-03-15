// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiple_prefix

import (
	"strings"
	"unicode"

	strings_ "github.com/searKing/golang/go/strings"
)

func DecimalFormatInt(number int, precision int) string {
	return DecimalFormatInt64(int64(number), precision)
}

func DecimalFormatUint(number uint, precision int) string {
	return DecimalFormatUint64(uint64(number), precision)
}

func DecimalFormatInt64(number int64, precision int) string {
	return DecimalMultiplePrefixTODO.Copy().SetInt64(number).FormatInt64(number, precision)
}

func DecimalFormatUint64(number uint64, precision int) string {
	return DecimalMultiplePrefixTODO.Copy().SetUint64(number).FormatUint64(number, precision)
}

func DecimalFormatFloat(number float64, precision int) string {
	return DecimalMultiplePrefixTODO.Copy().SetFloat64(number).FormatFloat(number, precision)
}

// SplitDecimal splits s into number, multiple_prefix and unparsed strings
func SplitDecimal(s string) (number string, prefix *DecimalMultiplePrefix, unparsed string) {
	splits := strings_.SplitPrefixNumber(s)
	if len(splits) < 2 {
		return "", nil, unparsed
	}
	number = splits[0]
	// trim any space between numbers and symbols
	unparsed = strings.TrimLeftFunc(splits[1], unicode.IsSpace)

	for _, prefix := range decimalPositiveMultiplePrefixes {
		if strings.HasPrefix(unparsed, prefix.Symbol()) {
			return number, prefix.Copy(), strings.TrimPrefix(unparsed, prefix.Symbol())
		}
	}
	for _, prefix := range decimalNegativeMultiplePrefixes {
		if strings.HasPrefix(unparsed, prefix.Symbol()) {
			return number, prefix.Copy(), strings.TrimPrefix(unparsed, prefix.Symbol())
		}
	}
	for _, prefix := range decimalZeroMultiplePrefixes {
		if strings.HasPrefix(unparsed, prefix.Symbol()) {
			return number, prefix.Copy(), strings.TrimPrefix(unparsed, prefix.Symbol())
		}
	}
	return number, nil, unparsed
}
