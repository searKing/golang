package unit

import (
	"fmt"
	"github.com/searKing/golang/go/format/radix"
	"math"
)

type Unit int

const HumanBasePower = 3

//power
const (
	One Unit = 3 * iota
	Kilo
	Mega
	Giga
	Tera
	Peta
	Exa
	Zetta
	Yotta
	Bronto
	Geop
	Ten           = 1
	Hundred       = 2
	Unimplemented = -1 //Unimplemented Unit
)

type pair struct {
	long  string
	short string
}

func (p *pair) Matched(tag string) bool {
	return p.short == tag || p.long == tag
}

var unitTags = map[Unit]pair{
	One:     {"One", ""},
	Kilo:    {"Kilo", "K"},
	Mega:    {"Mega", "M"},
	Giga:    {"Giga", "G"},
	Tera:    {"Tera", "T"},
	Peta:    {"Peta", "P"},
	Exa:     {"Exa", "E"},
	Zetta:   {"Zetta", "Z"},
	Yotta:   {"Yotta", "Y"},
	Bronto:  {"Bronto", "B"},
	Geop:    {"Geop", "Geop"},
	Ten:     {"Ten", "Ten"},
	Hundred: {"Hundred", "Hundred"},
}

// num 123kb
// tag is k
func New(tag string) (Unit, error) {
	for u, pair := range unitTags {
		if pair.Matched(tag) {
			return u, nil
		}
	}
	return Unimplemented, fmt.Errorf("Unimplemented unit tag: %v", tag)
}

func NewFromPower(power int64) Unit {
	return Unit(HumanBasePower * (power / 3))
}
func NewFromBase(base float64, r radix.Radix) (Unit, error) {
	kiloBase, err := Unit(Kilo).BaseFromRadix(r)
	if err != nil {
		return Unimplemented, err
	}
	// 幂

	power := int64(math.Floor(math.Log(math.Log(base)) / math.Log(float64(kiloBase))))
	return NewFromPower(power), nil
}

func (u Unit) String() (string, error) {
	if k, v := unitTags[u]; v {
		return k.short, nil
	}
	return "", fmt.Errorf("Unimplemented unit: %v", u)
}

func (u Unit) decimalBase() (int64, error) {
	return int64(math.Pow(1000, float64(u/3))), nil
}

func (u Unit) binaryBase() (int64, error) {
	return int64(math.Pow(1024, float64(u/3))), nil
}

func (u Unit) BaseFromRadix(r radix.Radix) (int64, error) {
	d, err := radix.IsDecimal(r)
	if err != nil {
		return 0, err
	}
	return u.Base(d)
}

// get base number
func (u Unit) Base(decimal bool) (int64, error) {
	if _, err := u.Valid(); err != nil {
		return 0, err
	}
	if decimal {
		return u.decimalBase()
	} else {
		return u.binaryBase()
	}
}
func (u Unit) Valid() (bool, error) {
	if _, ok := unitTags[u]; ok {
		return true, nil
	}
	return false, fmt.Errorf("unimplement unit: %v", u)
}

func ParseBase(num, scale float64, r radix.Radix) (humanNumber, humanBase float64, err error) {
	var u Unit = Kilo
	base, err := u.BaseFromRadix(r)
	if err != nil {
		return 0, 0, err
	}
	// 幂
	power := math.Floor(math.Log(math.Abs(num*scale)) / math.Log(float64(base)))

	humanUnit := NewFromPower(int64(power))
	humanBaseInt, err := humanUnit.BaseFromRadix(r)
	if err != nil {
		return 0, 0, err
	}

	humanBase = float64(humanBaseInt)
	humanNumber = num * scale / float64(humanBase)
	return humanNumber, humanBase, nil
}
