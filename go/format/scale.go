package format

import (
	"fmt"
	"github.com/searKing/golang/go/format/radix"
	"github.com/searKing/golang/go/format/unit"
	"strconv"
	"strings"
	"unicode"
)

type kind int

const (
	kindUint = iota
	kindInt
	kindFloat
	kindUnset
)

// Scale num's coefficient in string
type Scale struct {
	scale string
	kind  kind
}

// num = 123kb
// scale is 123k
func New(s string) Scale {
	return Scale{scale: s, kind: kindUnset}
}

func NewFromNum(scaleNum interface{}, r radix.Radix) Scale {
	var k kind
	switch scaleNum.(type) {
	case uint64:
		k = kindUint
	case int64:
		k = kindInt
	case float64:
		k = kindFloat
	default:
		k = kindUnset
	}
	rs, _ := r.String()
	return Scale{scale: fmt.Sprintf("%v%v", scaleNum, rs), kind: k}
}

// parse -123k  and get num's coeffieient
func (s Scale) ParseScale(r radix.Radix) (scale int64, err error) {
	var u unit.Unit
	switch s.scale[len(s.scale)-1] {
	case 'k':
		u = unit.Kilo
	case 'm':
		u = unit.Mega
	case 'g':
		u = unit.Giga
	case 't':
		u = unit.Tera
	case 'p':
		u = unit.Peta
	case 'b':
		u = unit.Bronto
	default:
		return 0, fmt.Errorf("invalid number unit format %s", s)
	}
	return u.BaseFromRadix(r)
}

// parse -123k  and get number
func (s Scale) ParseNum() (num interface{}, err error) {
	switch s.kind {
	case kindUint:
		return s.ParseUint(10, 0)
	case kindInt:
		return s.ParseInt(10, 0)
	case kindFloat:
		return s.ParseFloat(0)
	default:
		return 0, fmt.Errorf("invalid number kind format %s", s)
	}
}

func (s Scale) ParseInt(base int, bitSize int) (num int64, err error) {
	return strconv.ParseInt(s.trimUnit(), 10, 64)
}
func (s Scale) ParseUint(base int, bitSize int) (num uint64, err error) {
	return strconv.ParseUint(s.trimUnit(), 10, 64)
}
func (s Scale) ParseFloat(bitSize int) (num float64, err error) {
	return strconv.ParseFloat(s.trimUnit(), bitSize)
}

func (s Scale) ParseNumStr() (num string) {
	return s.trimUnit()
}

// parse -123k and get sign
func (s Scale) ParseSign() (neg bool) {
	if s.scale[0] == '-' {
		return true
	}
	return false
}

// format and get string like -123k
func (s Scale) FormatScale(r radix.Radix) (string, error) {
	switch s.kind {
	case kindUint:
		return s.FormatScaleUint(r, 10, 0)
	case kindInt:
		return s.FormatScaleInt(r, 10, 0)
	case kindFloat:
		return s.FormatScaleFloat(r, 0)
	default:
		return s.scale, nil
	}
}

func (s Scale) FormatScaleUint(r radix.Radix, base int, bitSize int) (string, error) {
	num, err := s.ParseUint(base, bitSize)
	if err != nil {
		return "", err
	}
	return s.formatScale(float64(num), r)
}

func (s Scale) FormatScaleInt(r radix.Radix, base int, bitSize int) (string, error) {
	num, err := s.ParseInt(base, bitSize)
	if err != nil {
		return "", err
	}
	return s.formatScale(float64(num), r)
}

func (s Scale) FormatScaleFloat(r radix.Radix, bitSize int) (string, error) {
	num, err := s.ParseFloat(bitSize)
	if err != nil {
		return "", err
	}
	return s.formatScale(num, r)
}

func (s Scale) formatScale(num float64, r radix.Radix) (string, error) {
	scale, err := s.ParseScale(r)
	if err != nil {
		return "", err
	}

	humanNumber, humanBase, err := unit.ParseBase(float64(num), float64(scale), r)
	if err != nil {
		return "", err
	}

	humanUnit, err := unit.NewFromBase(humanBase, r)
	if err != nil {
		return "", err
	}
	humanUnitStr, err := humanUnit.String()
	if err != nil {
		return "", err
	}
	switch s.kind {
	case kindUint:
		return fmt.Sprintf("%v%v", uint64(humanNumber), humanUnitStr), nil
	case kindInt:
		return fmt.Sprintf("%v%v", int64(humanNumber), humanUnitStr), nil
	case kindFloat:
		return fmt.Sprintf("%v%v", float64(humanNumber), humanUnitStr), nil
	default:
		return fmt.Sprintf("%v%v", humanNumber, humanUnitStr), nil
	}
}

func (s Scale) trimUnit() string {
	return strings.TrimRightFunc(s.scale, func(tag rune) bool {
		return !unicode.IsDigit(tag)
	})
}
func (s Scale) trimSign() string {
	return strings.TrimRightFunc(s.scale, func(tag rune) bool {
		return tag == '-'
	})
}
