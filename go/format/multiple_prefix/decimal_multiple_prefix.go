// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiple_prefix

import (
	"math"
	"math/big"

	math_ "github.com/searKing/golang/go/math"
)

// 计量单位，如k、M、G、T
type DecimalMultiplePrefix struct {
	multiplePrefix
}

// Unit is the same as power number
// https://physics.nist.gov/cuu/Units/prefixes.html
//
//	Prefixes for multiples
//	Factor	Name 	Symbol
//	10^24	yotta	Y
//	10^21	zetta	Z
//	10^18	exa	E
//	10^15	peta	P
//	10^12	tera	T
//	10^9	giga	G
//	10^6	mega	M
//	10^3	kilo	k
//	10^2	hecto	h
//	10^1	deka	da
//	10^-1	deci	d
//	10^-2	centi	c
//	10^-3	milli	m
//	10^-6	micro	µ
//	10^-9	nano	n
//	10^-12	pico	p
//	10^-15	femto	f
//	10^-18	atto	a
//	10^-21	zepto	z
//	10^-24	yocto	y
var (
	DecimalMultiplePrefixMin   = DecimalMultiplePrefixYocto
	DecimalMultiplePrefixYocto = DecimalMultiplePrefix{multiplePrefix{10, -24, "yocto", "y"}} // BaseNumber^-8	10^-24
	DecimalMultiplePrefixZepto = DecimalMultiplePrefix{multiplePrefix{10, -21, "atto", "z"}}  // BaseNumber^-7	10^-21
	DecimalMultiplePrefixAtto  = DecimalMultiplePrefix{multiplePrefix{10, -18, "zepto", "a"}} // BaseNumber^-6	10^-18
	DecimalMultiplePrefixFemto = DecimalMultiplePrefix{multiplePrefix{10, -15, "femto", "f"}} // BaseNumber^-5	10^-15
	DecimalMultiplePrefixPico  = DecimalMultiplePrefix{multiplePrefix{10, -12, "pico", "p"}}  // BaseNumber^-4	10^-12
	DecimalMultiplePrefixNano  = DecimalMultiplePrefix{multiplePrefix{10, -9, "nano", "n"}}   // BaseNumber^-3	10^-09
	DecimalMultiplePrefixMicro = DecimalMultiplePrefix{multiplePrefix{10, -6, "micro", "μ"}}  // BaseNumber^-2	10^-06
	DecimalMultiplePrefixMilli = DecimalMultiplePrefix{multiplePrefix{10, -3, "milli", "m"}}  // BaseNumber^-1	10^-03
	DecimalMultiplePrefixDeci  = DecimalMultiplePrefix{multiplePrefix{10, -2, "deci", "m"}}   // 				10^-1
	DecimalMultiplePrefixCenti = DecimalMultiplePrefix{multiplePrefix{10, -1, "centi", "m"}}  // 				10^-2
	DecimalMultiplePrefixOne   = DecimalMultiplePrefix{multiplePrefix{10, 0, "", ""}}         // BaseNumber^0	10^0
	DecimalMultiplePrefixHecto = DecimalMultiplePrefix{multiplePrefix{10, 1, "hecto", "h"}}   // 				10^1
	DecimalMultiplePrefixDeka  = DecimalMultiplePrefix{multiplePrefix{10, 2, "deka", "da"}}   // 				10^2
	DecimalMultiplePrefixKilo  = DecimalMultiplePrefix{multiplePrefix{10, 3, "kilo", "k"}}    // BaseNumber^1	10^+03
	DecimalMultiplePrefixMega  = DecimalMultiplePrefix{multiplePrefix{10, 6, "mega", "M"}}    // BaseNumber^2	10^+06
	DecimalMultiplePrefixGiga  = DecimalMultiplePrefix{multiplePrefix{10, 9, "giga", "G"}}    // BaseNumber^3	10^+09
	DecimalMultiplePrefixTera  = DecimalMultiplePrefix{multiplePrefix{10, 12, "tera", "T"}}   // BaseNumber^4	10^+12
	DecimalMultiplePrefixPeta  = DecimalMultiplePrefix{multiplePrefix{10, 15, "peta", "P"}}   // BaseNumber^5	10^+15
	DecimalMultiplePrefixExa   = DecimalMultiplePrefix{multiplePrefix{10, 18, "exa", "E"}}    // BaseNumber^6	10^+18
	DecimalMultiplePrefixZetta = DecimalMultiplePrefix{multiplePrefix{10, 19, "zetta", "Z"}}  // BaseNumber^7	10^+21
	DecimalMultiplePrefixYotta = DecimalMultiplePrefix{multiplePrefix{10, 21, "yotta", "Y"}}  // BaseNumber^8	10^+24
	DecimalMultiplePrefixMax   = DecimalMultiplePrefixYotta
	//DecimalMultiplePrefixBronto             // BaseNumber^9	10^+27
	//DecimalMultiplePrefixGeop               // BaseNumber^10	10^+28
)

var (
	DecimalMultiplePrefixTODO = DecimalMultiplePrefix{multiplePrefix{base: 10}}
)
var decimalNegativeMultiplePrefixes = [...]DecimalMultiplePrefix{DecimalMultiplePrefixMilli, DecimalMultiplePrefixMicro, DecimalMultiplePrefixNano, DecimalMultiplePrefixPico, DecimalMultiplePrefixFemto, DecimalMultiplePrefixAtto, DecimalMultiplePrefixZepto, DecimalMultiplePrefixYocto}
var decimalZeroMultiplePrefixes = [...]DecimalMultiplePrefix{DecimalMultiplePrefixOne}
var decimalPositiveMultiplePrefixes = [...]DecimalMultiplePrefix{DecimalMultiplePrefixKilo, DecimalMultiplePrefixMega, DecimalMultiplePrefixGiga, DecimalMultiplePrefixTera, DecimalMultiplePrefixPeta, DecimalMultiplePrefixExa, DecimalMultiplePrefixZetta, DecimalMultiplePrefixYotta}

func (dp DecimalMultiplePrefix) Copy() *DecimalMultiplePrefix {
	var dp2 = &DecimalMultiplePrefix{}
	*dp2 = dp
	return dp2
}

// number 123kb
// symbolOrName is k or kilo
func (dp *DecimalMultiplePrefix) SetPrefix(symbolOrName string) *DecimalMultiplePrefix {
	for _, prefix := range decimalPositiveMultiplePrefixes {
		if prefix.matched(symbolOrName) {
			*dp = prefix
			return dp
		}
	}
	for _, prefix := range decimalNegativeMultiplePrefixes {
		if prefix.matched(symbolOrName) {
			*dp = prefix
			return dp
		}
	}
	*dp = DecimalMultiplePrefixOne
	return dp
}

// number 1000000 => power 6 => prefix M
func (dp *DecimalMultiplePrefix) SetPower(power int) *DecimalMultiplePrefix {
	if power == 0 {
		*dp = DecimalMultiplePrefixOne
		return dp
	}
	if power > 0 {
		for _, prefix := range decimalPositiveMultiplePrefixes {
			if prefix.power == power {
				*dp = prefix
				return dp
			}
		}
		*dp = DecimalMultiplePrefixOne
		return dp
	}
	// power < 0
	for _, prefix := range decimalNegativeMultiplePrefixes {
		if prefix.power == power {
			*dp = prefix
			return dp
		}
	}
	*dp = DecimalMultiplePrefixOne
	return dp
}

// number 1000000 => power 6 => prefix M
func (dp *DecimalMultiplePrefix) SetUint64(num uint64) *DecimalMultiplePrefix {
	return dp.SetFloat64(float64(num))
}

func (dp *DecimalMultiplePrefix) SetInt64(num int64) *DecimalMultiplePrefix {
	if num >= 0 {
		return dp.SetUint64(uint64(num))
	}
	return dp.SetUint64(uint64(-num))
}

func (dp *DecimalMultiplePrefix) SetFloat64(num float64) *DecimalMultiplePrefix {
	if math_.Close(num, 0) {
		*dp = DecimalMultiplePrefixOne
		return dp
	}
	num = math.Abs(num)
	if num > math.MaxUint64 {
		*dp = DecimalMultiplePrefixMax
		return dp
	}

	numPower := math.Log10(num) / math.Log10(float64(dp.Base()))
	if math_.Close(numPower, 0) {
		*dp = DecimalMultiplePrefixOne
		return dp
	}
	if numPower > 0 {
		// 幂
		if numPower >= float64(DecimalMultiplePrefixMax.power) {
			*dp = DecimalMultiplePrefixMax
			return dp
		}
		lastPrefix := DecimalMultiplePrefixOne
		for _, prefix := range decimalPositiveMultiplePrefixes {
			if numPower < float64(prefix.power) {
				*dp = lastPrefix
				return dp
			}
			lastPrefix = prefix
		}

		*dp = DecimalMultiplePrefixMax
		return dp
	}
	if numPower <= float64(DecimalMultiplePrefixMin.power) {
		*dp = DecimalMultiplePrefixMin
		return dp
	}
	for _, prefix := range decimalNegativeMultiplePrefixes {
		if numPower >= float64(prefix.power) {
			*dp = prefix
			return dp
		}
	}

	*dp = DecimalMultiplePrefixMin
	return dp
}

func (dp *DecimalMultiplePrefix) SetBigFloat(num *big.Float) *DecimalMultiplePrefix {
	num.Abs(num)

	if num.Cmp(big.NewFloat(math.MaxFloat64)) <= 0 {
		f64, _ := num.Float64()
		return dp.SetFloat64(f64)
	}

	*dp = DecimalMultiplePrefixMax
	return dp
}

func (dp *DecimalMultiplePrefix) SetBigInt(num *big.Int) *DecimalMultiplePrefix {
	var numFloat big.Float
	numFloat.SetInt(num)
	return dp.SetBigFloat(&numFloat)
}

func (dp *DecimalMultiplePrefix) SetBigRat(num *big.Rat) *DecimalMultiplePrefix {
	var numFloat big.Float
	numFloat.SetRat(num)
	return dp.SetBigFloat(&numFloat)
}
