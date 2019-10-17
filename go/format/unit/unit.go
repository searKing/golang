package unit

import (
	"fmt"
	"strconv"
)

type Unit int // 单位，如K、M、G、T

// Unit is the same as power number
// https://physics.nist.gov/cuu/Units/prefixes.html
// https://physics.nist.gov/cuu/Units/binary.html
const (
	Yocto  Unit = iota - 8 // BaseFormat^-8	10^-24
	Zepto                  // BaseFormat^-7	10^-21
	Atto                   // BaseFormat^-6	10^-18
	Femto                  // BaseFormat^-5	10^-15
	Pico                   // BaseFormat^-4	10^-12
	Nano                   // BaseFormat^-3	10^-09
	Micro                  // BaseFormat^-2	10^-06
	Milli                  // BaseFormat^-1	10^-03
	One                    // BaseFormat^0	10^0
	//Hecto                  // BaseFormat^0	10^1
	//Deka                   // BaseFormat^0	10^2
	Kilo                   // BaseFormat^1	10^+03
	Mega                   // BaseFormat^2	10^+06
	Giga                   // BaseFormat^3	10^+09
	Tera                   // BaseFormat^4	10^+12
	Peta                   // BaseFormat^5	10^+15
	Exa                    // BaseFormat^6	10^+18
	Zetta                  // BaseFormat^7	10^+21
	Yotta                  // BaseFormat^8	10^+24
	Bronto                 // BaseFormat^9	10^+27
	Geop                   // BaseFormat^10	10^+28
	Max    = Geop
)

// number 123kb
// tag is k
func ParseUnit(tag string) (Unit, error) {
	for u, pair := range unitTagPairs {
		if pair.Matched(tag) {
			return u, nil
		}
	}
	return 0, fmt.Errorf("unimplemented unit tag: %v", tag)
}

func (u Unit) String() string {
	if k, v := unitTagPairs[u]; v {
		return k.short
	}
	return "Unit(" + strconv.FormatInt(int64(u), 10) + ")"
}

func (u Unit) Registered() bool {
	if _, ok := unitTagPairs[u]; ok {
		return true
	}
	return false
}
