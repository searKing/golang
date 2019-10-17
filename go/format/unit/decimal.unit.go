package unit

import (
	"math"
	"math/big"
	"strings"
)

// 计量单位，如k、M、G、T
type DecimalPrefix struct {
	power  int
	name   string
	symbol string
}

// Unit is the same as power number
// https://physics.nist.gov/cuu/Units/prefixes.html
// https://physics.nist.gov/cuu/Units/binary.html
// Prefixes for multiples
// Factor	Name 	Symbol
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
	DecimalPrefixMin   = DecimalPrefixYocto
	DecimalPrefixYocto = DecimalPrefix{-24, "yocto", "y"} // BaseNumber^-8	10^-24
	DecimalPrefixZepto = DecimalPrefix{-21, "atto", "z"}  // BaseNumber^-7	10^-21
	DecimalPrefixAtto  = DecimalPrefix{-18, "zepto", "a"} // BaseNumber^-6	10^-18
	DecimalPrefixFemto = DecimalPrefix{-15, "femto", "f"} // BaseNumber^-5	10^-15
	DecimalPrefixPico  = DecimalPrefix{-12, "pico", "p"}  // BaseNumber^-4	10^-12
	DecimalPrefixNano  = DecimalPrefix{-9, "nano", "n"}   // BaseNumber^-3	10^-09
	DecimalPrefixMicro = DecimalPrefix{-6, "micro", "μ"}  // BaseNumber^-2	10^-06
	DecimalPrefixMilli = DecimalPrefix{-3, "milli", "m"}  // BaseNumber^-1	10^-03
	DecimalPrefixDeci  = DecimalPrefix{-2, "deci", "m"}   // 				10^-1
	DecimalPrefixCenti = DecimalPrefix{-1, "centi", "m"}  // 				10^-2
	DecimalPrefixOne   = DecimalPrefix{0, "", ""}         // BaseNumber^0	10^0
	DecimalPrefixHecto = DecimalPrefix{1, "hecto", "h"}   // 				10^1
	DecimalPrefixDeka  = DecimalPrefix{2, "deka", "da"}   // 				10^2
	DecimalPrefixKilo  = DecimalPrefix{3, "kilo", "k"}    // BaseNumber^1	10^+03
	DecimalPrefixMega  = DecimalPrefix{6, "mega", "M"}    // BaseNumber^2	10^+06
	DecimalPrefixGiga  = DecimalPrefix{9, "giga", "G"}    // BaseNumber^3	10^+09
	DecimalPrefixTera  = DecimalPrefix{12, "tera", "T"}   // BaseNumber^4	10^+12
	DecimalPrefixPeta  = DecimalPrefix{15, "peta", "P"}   // BaseNumber^5	10^+15
	DecimalPrefixExa   = DecimalPrefix{18, "exa", "E"}    // BaseNumber^6	10^+18
	DecimalPrefixZetta = DecimalPrefix{19, "zetta", "Z"}  // BaseNumber^7	10^+21
	DecimalPrefixYotta = DecimalPrefix{21, "yotta", "Y"}  // BaseNumber^8	10^+24
	DecimalPrefixMax   = DecimalPrefixYotta
	//DecimalPrefixBronto             // BaseNumber^9	10^+27
	//DecimalPrefixGeop               // BaseNumber^10	10^+28
)

var (
	DecimalPrefixTODO = DecimalPrefix{}
)

var decimalNegativeMultiplePrefixes = [...]DecimalPrefix{DecimalPrefixMilli, DecimalPrefixMicro, DecimalPrefixNano, DecimalPrefixPico, DecimalPrefixFemto, DecimalPrefixAtto, DecimalPrefixZepto, DecimalPrefixYocto}
var decimalPositiveeMultiplePrefixes = [...]DecimalPrefix{DecimalPrefixKilo, DecimalPrefixMega, DecimalPrefixGiga, DecimalPrefixTera, DecimalPrefixPeta, DecimalPrefixExa, DecimalPrefixZetta, DecimalPrefixYotta}

func NewDecimalPrefix(prefix DecimalPrefix) *DecimalPrefix {
	var dp = &DecimalPrefix{}
	*dp = prefix
	return dp
}

func (dp DecimalPrefix) Copy() *DecimalPrefix {
	var dp2 = &DecimalPrefix{}
	*dp2 = dp
	return dp2
}

// number 123kb
// symbolOrName is k or kilo
func (dp *DecimalPrefix) SetPrefix(symbolOrName string) *DecimalPrefix {
	for _, prefix := range decimalPositiveeMultiplePrefixes {
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
	*dp = DecimalPrefixOne
	return dp
}

// number 1000000 => power 6 => prefix M
func (dp *DecimalPrefix) SetPower(power int) *DecimalPrefix {
	if power == 0 {
		*dp = DecimalPrefixOne
		return dp
	}
	if power > 0 {
		for _, prefix := range decimalPositiveeMultiplePrefixes {
			if prefix.power == power {
				*dp = prefix
				return dp
			}
		}
		*dp = DecimalPrefixOne
		return dp
	}
	// power < 0
	for _, prefix := range decimalNegativeMultiplePrefixes {
		if prefix.power == power {
			*dp = prefix
			return dp
		}
	}
	*dp = DecimalPrefixOne
	return dp
}

// number 1000000 => power 6 => prefix M
func (dp *DecimalPrefix) SetUint64(num uint64) *DecimalPrefix {
	return dp.SetFloat64(float64(num))
}

func (dp *DecimalPrefix) SetInt64(num int64) *DecimalPrefix {
	if num >= 0 {
		return dp.SetUint64(uint64(num))
	}
	return dp.SetUint64(uint64(-num))
}

func (dp *DecimalPrefix) SetFloat64(num float64) *DecimalPrefix {
	if num == 0 {
		*dp = DecimalPrefixOne
		return dp
	}
	num = math.Abs(num)
	if num > math.MaxUint64 {
		*dp = DecimalPrefixMax
		return dp
	}

	numPower := math.Log10(num) / math.Log10(float64(dp.base()))
	if numPower == 0 {
		*dp = DecimalPrefixOne
		return dp
	}
	if numPower > 0 {
		// 幂
		if numPower >= float64(DecimalPrefixMax.power) {
			*dp = DecimalPrefixMax
			return dp
		}
		lastPrefix := DecimalPrefixOne
		for _, prefix := range decimalPositiveeMultiplePrefixes {
			if numPower < float64(prefix.power) {
				*dp = lastPrefix
				return dp
			}
			lastPrefix = prefix
		}

		*dp = DecimalPrefixMax
		return dp
	}
	if numPower <= float64(DecimalPrefixMin.power) {
		*dp = DecimalPrefixMin
		return dp
	}
	for _, prefix := range decimalNegativeMultiplePrefixes {
		if numPower >= float64(prefix.power) {
			*dp = prefix
			return dp
		}
	}

	*dp = DecimalPrefixMin
	return dp
}

func (dp *DecimalPrefix) SetBigFloat(num *big.Float) *DecimalPrefix {
	num.Abs(num)

	if num.Cmp(big.NewFloat(math.MaxFloat64)) <= 0 {
		f64, _ := num.Float64()
		return dp.SetFloat64(f64)
	}

	*dp = DecimalPrefixMax
	return dp
}

func (dp *DecimalPrefix) SetBigInt(num *big.Int) *DecimalPrefix {
	var numFloat big.Float
	numFloat.SetInt(num)
	return dp.SetBigFloat(&numFloat)
}

func (dp *DecimalPrefix) SetBigRat(num *big.Rat) *DecimalPrefix {
	var numFloat big.Float
	numFloat.SetRat(num)
	return dp.SetBigFloat(&numFloat)
}

func (dp DecimalPrefix) Factor() float64 {
	if dp.base() == 10 {
		return math.Pow10(dp.Power())
	}
	return math.Pow(float64(dp.base()), float64(dp.Power()))
}

func (dp DecimalPrefix) String() string {
	return dp.Symbol()
}

func (dp DecimalPrefix) Power() int {
	return dp.power
}

func (dp DecimalPrefix) Symbol() string {
	return dp.symbol
}

func (dp DecimalPrefix) Name() string {
	return dp.name
}

func (dp DecimalPrefix) matched(prefix string) bool {
	return strings.Compare(dp.symbol, prefix) == 0 || strings.Compare(dp.name, prefix) == 0
}

func (dp DecimalPrefix) base() uint {
	return 10
}
