package unit

// 幂

// number 1000000 => power 6 => unit M
func ParsePower(power uint) Unit {
	return Unit(power)
}

func (u Unit) Power() uint {
	return uint(u)
}
