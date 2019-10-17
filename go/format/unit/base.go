package unit

import "math/big"

//go:generate go-enum -type BaseFormat -trimprefix=BaseFormat
type BaseFormat uint

const (
	BaseFormatDecimal BaseFormat = 1000
	BaseFormatBinary  BaseFormat = 1024
)

// get base number
// decimal	K -> 1000
// binary	K -> 1024
func (u Unit) Base(baseFormat BaseFormat) *big.Int {
	res := big.NewInt(1)
	for i := uint64(0); i < uint64(u); i++ {
		res = res.Mul(res, big.NewInt(int64(baseFormat)))
	}

	return res
}
