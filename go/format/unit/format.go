package unit

import (
	"fmt"

	"github.com/searKing/golang/go/math"
)

func FormatInt(number int, precision int) string {
	return FormatInt64(int64(number), precision)
}

func FormatUint(number uint, precision int) string {
	return FormatUint64(uint64(number), precision)
}

func FormatInt64(number int64, precision int) string {
	return FormatFloatWithUnit(float64(number), *DecimalPrefixTODO.Copy().SetInt64(number), precision)
}

func FormatUint64(number uint64, precision int) string {
	return FormatFloatWithUnit(float64(number), *DecimalPrefixTODO.Copy().SetUint64(number), precision)
}

func FormatFloat(number float64, precision int) string {
	return FormatFloatWithUnit(number, *DecimalPrefixTODO.Copy().SetFloat64(number), precision)
}

func FormatFloatWithUnit(number float64, prefix DecimalPrefix, precision int) string {
	humanBase := prefix.Factor()
	humanNumber := number / humanBase
	if precision >= 0 {
		humanNumber = math.TruncPrecision(humanNumber, precision)
	}
	return fmt.Sprintf("%g%s", humanNumber, prefix)
}
