package radix

import "fmt"

type Radix int

const (
	Binary      = iota // 二进制
	Octonary           // 八进制
	Decimal            // 十进制
	Hexadecimal        // 十六进制
)

var radixs = map[Radix]string{
	Binary:      "Binary",
	Octonary:    "Octonary",
	Decimal:     "Decimal",
	Hexadecimal: "Hexadecimal",
}

func (r Radix) String() (string, error) {
	if v, ok := radixs[r]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unimplement radix: %v", r)
}

func (r Radix) Valid() (bool, error) {
	if _, ok := radixs[r]; ok {
		return true, nil
	}
	return false, fmt.Errorf("unimplement radix: %v", r)
}

func IsDecimal(r Radix) (bool, error) {
	if _, err := r.Valid(); err != nil {
		return false, err
	}
	if r == Decimal {
		return true, nil
	}
	return false, nil
}
