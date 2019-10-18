package multiple_prefix

import (
	"fmt"
	"io"
	"strings"
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

func TrimBinaryMultiplePrefix(s string) string {
	var value float64
	var unparsed string
	count, err := fmt.Sscanf(s, `%v%s`, &value, &unparsed)

	if (err != nil && err != io.EOF) || (count == 0) {
		var value int64
		count, err := fmt.Sscanf(s, `%v%s`, &value, &unparsed)
		if (err != nil && err != io.EOF) || (count == 0) {
			return s
		}
	}

	for _, prefix := range binaryPositiveeMultiplePrefixes {
		if strings.HasPrefix(unparsed, prefix.Symbol()) {
			return strings.TrimPrefix(unparsed, prefix.Symbol())
		}
	}
	return unparsed
}
