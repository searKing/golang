package unit

import (
	"fmt"

	"github.com/searKing/golang/go/math"
)

func FormatInt(number int, baseFormat BaseFormat, precision int) string {
	return FormatInt64(int64(number), baseFormat, precision)
}

func FormatUint(number uint, baseFormat BaseFormat, precision int) string {
	return FormatUint64(uint64(number), baseFormat, precision)
}

func FormatInt64(number int64, baseFormat BaseFormat, precision int) string {
	unit := ProbeUnitFromInt64(number, baseFormat)
	return FormatFloatWithUnit(float64(number), baseFormat, unit, precision)
}

func FormatUint64(number uint64, baseFormat BaseFormat, precision int) string {
	unit := ProbeUnitFromUint64(number, baseFormat)
	return FormatFloatWithUnit(float64(number), baseFormat, unit, precision)

}

func FormatFloat(number float64, baseFormat BaseFormat, precision int) string {
	unit := ProbeUnitFromFloat64(number, baseFormat)
	return FormatFloatWithUnit(number, baseFormat, unit, precision)
}

func FormatFloatWithUnit(number float64, baseFormat BaseFormat, unit Unit, precision int) string {
	humanBase := unit.Base(baseFormat).Uint64()
	humanNumber := number / float64(humanBase)
	if precision >= 0 {
		humanNumber = math.TruncPrecision(humanNumber, precision)
	}
	return fmt.Sprintf("%g%s", humanNumber, unit)
}
