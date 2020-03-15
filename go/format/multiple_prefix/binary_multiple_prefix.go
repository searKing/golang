// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiple_prefix

import (
	"math"
	"math/big"

	math_ "github.com/searKing/golang/go/math"
)

// 计量单位，如Ki、Mi、Gi、Ti
type BinaryMultiplePrefix struct {
	multiplePrefix
}

// https://physics.nist.gov/cuu/Units/binary.html
//
// 				Prefixes for binary multiples
// 	Factor 	Name 	Symbol 	Origin					Derivation
// 	210		kibi	Ki		kilobinary: (2^10)^1	kilo: (10^3)^1
// 	220		mebi	Mi		megabinary: (2^10)^2 	mega: (10^3)^2
// 	230		gibi	Gi		gigabinary: (2^10)^3	giga: (10^3)^3
// 	240		tebi	Ti		terabinary: (2^10)^4	tera: (10^3)^4
// 	250		pebi	Pi		petabinary: (2^10)^5	peta: (10^3)^5
// 	260		exbi	Ei		exabinary:  (2^10)^6	exa:  (10^3)^6
var (
	BinaryMultiplePrefixMin  = BinaryMultiplePrefixOne
	BinaryMultiplePrefixOne  = BinaryMultiplePrefix{multiplePrefix{2, 0, "", ""}}        // BaseNumber^0	10^0
	BinaryMultiplePrefixKibi = BinaryMultiplePrefix{multiplePrefix{2, 10, "kibi", "Ki"}} // BaseNumber^1	10^+03
	BinaryMultiplePrefixMebi = BinaryMultiplePrefix{multiplePrefix{2, 20, "mebi", "Mi"}} // BaseNumber^2	10^+06
	BinaryMultiplePrefixGibi = BinaryMultiplePrefix{multiplePrefix{2, 30, "gibi", "Gi"}} // BaseNumber^3	10^+09
	BinaryMultiplePrefixTebi = BinaryMultiplePrefix{multiplePrefix{2, 40, "tebi", "Ti"}} // BaseNumber^4	10^+12
	BinaryMultiplePrefixPebi = BinaryMultiplePrefix{multiplePrefix{2, 50, "pebi", "Pi"}} // BaseNumber^5	10^+15
	BinaryMultiplePrefixExbi = BinaryMultiplePrefix{multiplePrefix{2, 60, "exbi", "Ei"}} // BaseNumber^6	10^+18
	BinaryMultiplePrefixMax  = BinaryMultiplePrefixExbi
)

var (
	BinaryMultiplePrefixTODO = BinaryMultiplePrefix{multiplePrefix{base: 2}}
)
var binaryNegativeMultiplePrefixes = [...]BinaryMultiplePrefix{}
var binaryZeroMultiplePrefixes = [...]BinaryMultiplePrefix{BinaryMultiplePrefixOne}
var binaryPositiveMultiplePrefixes = [...]BinaryMultiplePrefix{BinaryMultiplePrefixKibi, BinaryMultiplePrefixMebi, BinaryMultiplePrefixGibi, BinaryMultiplePrefixTebi, BinaryMultiplePrefixPebi, BinaryMultiplePrefixExbi}

func (dp BinaryMultiplePrefix) Copy() *BinaryMultiplePrefix {
	var dp2 = &BinaryMultiplePrefix{}
	*dp2 = dp
	return dp2
}

// number 123kb
// symbolOrName is k or kilo
func (dp *BinaryMultiplePrefix) SetPrefix(symbolOrName string) *BinaryMultiplePrefix {
	for _, prefix := range binaryPositiveMultiplePrefixes {
		if prefix.matched(symbolOrName) {
			*dp = prefix
			return dp
		}
	}
	*dp = BinaryMultiplePrefixOne
	return dp
}

// number 1000000 => power 6 => prefix M
func (dp *BinaryMultiplePrefix) SetPower(power int) *BinaryMultiplePrefix {
	if power == 0 {
		*dp = BinaryMultiplePrefixOne
		return dp
	}
	if power > 0 {
		for _, prefix := range binaryPositiveMultiplePrefixes {
			if prefix.power == power {
				*dp = prefix
				return dp
			}
		}
		*dp = BinaryMultiplePrefixOne
		return dp
	}
	*dp = BinaryMultiplePrefixOne
	return dp
}

// number 1000000 => power 6 => prefix M
func (dp *BinaryMultiplePrefix) SetUint64(num uint64) *BinaryMultiplePrefix {
	return dp.SetFloat64(float64(num))
}

func (dp *BinaryMultiplePrefix) SetInt64(num int64) *BinaryMultiplePrefix {
	if num >= 0 {
		return dp.SetUint64(uint64(num))
	}
	return dp.SetUint64(uint64(-num))
}

func (dp *BinaryMultiplePrefix) SetFloat64(num float64) *BinaryMultiplePrefix {
	if math_.Close(num, 0) {
		*dp = BinaryMultiplePrefixOne
		return dp
	}
	num = math.Abs(num)
	if num > math.MaxUint64 {
		*dp = BinaryMultiplePrefixMax
		return dp
	}

	numPower := math.Log10(num) / math.Log10(float64(dp.Base()))
	if math_.Close(numPower, 0) {
		*dp = BinaryMultiplePrefixOne
		return dp
	}
	if numPower > 0 {
		// 幂
		if numPower >= float64(BinaryMultiplePrefixMax.power) {
			*dp = BinaryMultiplePrefixMax
			return dp
		}
		lastPrefix := BinaryMultiplePrefixOne
		for _, prefix := range binaryPositiveMultiplePrefixes {
			if numPower < float64(prefix.power) {
				*dp = lastPrefix
				return dp
			}
			lastPrefix = prefix
		}

		*dp = BinaryMultiplePrefixMax
		return dp
	}
	if numPower <= float64(BinaryMultiplePrefixMin.power) {
		*dp = BinaryMultiplePrefixMin
		return dp
	}

	*dp = BinaryMultiplePrefixMin
	return dp
}

func (dp *BinaryMultiplePrefix) SetBigFloat(num *big.Float) *BinaryMultiplePrefix {
	num.Abs(num)

	if num.Cmp(big.NewFloat(math.MaxFloat64)) <= 0 {
		f64, _ := num.Float64()
		return dp.SetFloat64(f64)
	}

	*dp = BinaryMultiplePrefixMax
	return dp
}

func (dp *BinaryMultiplePrefix) SetBigInt(num *big.Int) *BinaryMultiplePrefix {
	var numFloat big.Float
	numFloat.SetInt(num)
	return dp.SetBigFloat(&numFloat)
}

func (dp *BinaryMultiplePrefix) SetBigRat(num *big.Rat) *BinaryMultiplePrefix {
	var numFloat big.Float
	numFloat.SetRat(num)
	return dp.SetBigFloat(&numFloat)
}
