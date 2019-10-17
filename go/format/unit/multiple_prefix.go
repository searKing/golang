package unit

import (
	"fmt"
	"math"
	"strings"

	math_ "github.com/searKing/golang/go/math"
)

// 计量单位，如k、M、G、T
type multiplePrefix struct {
	base   int
	power  int
	name   string
	symbol string
}

func (dp multiplePrefix) FormatInt64(number int64, precision int) string {
	return dp.FormatFloat(float64(number), precision)
}

func (dp multiplePrefix) FormatUint64(number uint64, precision int) string {
	return dp.FormatFloat(float64(number), precision)
}

func (dp multiplePrefix) FormatFloat(number float64, precision int) string {
	humanBase := dp.Factor()
	humanNumber := number / humanBase
	if precision >= 0 {
		humanNumber = math_.TruncPrecision(humanNumber, precision)
	}
	return fmt.Sprintf("%g%s", humanNumber, dp)
}

func (dp multiplePrefix) Factor() float64 {
	if dp.Base() == 10 {
		return math.Pow10(dp.Power())
	}
	return math.Pow(float64(dp.Base()), float64(dp.Power()))
}

func (dp multiplePrefix) String() string {
	return dp.Symbol()
}

func (dp multiplePrefix) Base() int {
	return dp.base
}

func (dp multiplePrefix) Power() int {
	return dp.power
}

func (dp multiplePrefix) Symbol() string {
	return dp.symbol
}

func (dp multiplePrefix) Name() string {
	return dp.name
}

func (dp multiplePrefix) matched(prefix string) bool {
	return strings.Compare(dp.symbol, prefix) == 0 || strings.Compare(dp.name, prefix) == 0
}
