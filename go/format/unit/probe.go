package unit

import (
	"math"
	"math/big"
)

// number 1000000 => power 6 => unit M

func ProbeUnitFromUint(num uint, baseFormat BaseFormat) Unit {
	return ProbeUnitFromUint64(uint64(num), baseFormat)
}

func ProbeUnitFromInt(num int, baseFormat BaseFormat) Unit {
	return ProbeUnitFromInt64(int64(num), baseFormat)
}

func ProbeUnitFromUint64(num uint64, baseFormat BaseFormat) Unit {
	if baseFormat == 1 {
		return One
	}
	kiloBase := Kilo.Base(baseFormat).Uint64() // 1000

	var power uint
	// 幂
	for {
		num = num / kiloBase
		if num == 0 {
			break
		}
		if power >= uint(Max) {
			break
		}
		power = power + 1

	}
	return ParsePower(power)
}

func ProbeUnitFromInt64(num int64, baseFormat BaseFormat) Unit {
	if num >= 0 {
		return ProbeUnitFromUint64(uint64(num), baseFormat)
	}
	return ProbeUnitFromUint64(uint64(-num), baseFormat)
}

func ProbeUnitFromFloat32(num float32, baseFormat BaseFormat) Unit {
	return ProbeUnitFromFloat64(float64(num), baseFormat)
}

func ProbeUnitFromFloat64(num float64, baseFormat BaseFormat) Unit {
	num = math.Abs(num)
	if num > math.MaxUint64 {
		return Geop
	}
	return ProbeUnitFromUint64(uint64(num), baseFormat)
}

func ProbeUnitFromBigInt(num *big.Int, baseFormat BaseFormat) Unit {
	if baseFormat == 1 {
		return One
	}
	kiloBase := Kilo.Base(baseFormat) // 1000

	num.Abs(num)
	var power uint
	// 幂
	for {
		num.Div(num, kiloBase)
		if num.Cmp(big.NewInt(0)) == 0 {
			break
		}
		if power >= uint(Max) {
			break
		}
		power = power + 1
	}
	return ParsePower(power)
}

func ProbeUnitFromBigFloat(num *big.Float, baseFormat BaseFormat) Unit {
	var numInt big.Int
	num.Int(&numInt)
	return ProbeUnitFromBigInt(&numInt, baseFormat)
}

func ProbeUnitFromBigRat(num *big.Rat, baseFormat BaseFormat) Unit {
	var numFloat big.Float
	numFloat.SetRat(num)
	return ProbeUnitFromBigFloat(&numFloat, baseFormat)
}
