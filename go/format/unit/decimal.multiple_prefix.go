package unit

//import "strings"
//
//// https://physics.nist.gov/cuu/Units/prefixes.html
//// Prefixes for multiples
//// Factor	Name 	Symbol
////	10^24	yotta	Y
////	10^21	zetta	Z
////	10^18	exa	E
////	10^15	peta	P
////	10^12	tera	T
////	10^9	giga	G
////	10^6	mega	M
////	10^3	kilo	k
////	10^2	hecto	h
////	10^1	deka	da
////	10^-1	deci	d
////	10^-2	centi	c
////	10^-3	milli	m
////	10^-6	micro	µ
////	10^-9	nano	n
////	10^-12	pico	p
////	10^-15	femto	f
////	10^-18	atto	a
////	10^-21	zepto	z
////	10^-24	yocto	y
//type decimalMultiplePrefix struct {
//	power  int
//	name   string
//	symbol string
//}
//
//var decimalMultiplePrefixes = map[DecimalPrefix]decimalMultiplePrefix{
//	DecimalPrefixYocto: {-24, "yocto", "y"},
//	DecimalPrefixAtto:  {-21, "atto", "z"},
//	DecimalPrefixZepto: {-18, "zepto", "a"},
//	DecimalPrefixFemto: {-15, "femto", "f"},
//	DecimalPrefixPico:  {-12, "pico", "p"},
//	DecimalPrefixNano:  {-9, "nano", "n"},
//	DecimalPrefixMicro: {-6, "micro", "μ"},
//	DecimalPrefixMilli: {-3, "milli", "m"},
//	DecimalPrefixDeci:  {-2, "deci", "m"},
//	DecimalPrefixCenti: {-1, "centi", "m"},
//	DecimalPrefixOne:   {0, "", ""},
//	DecimalPrefixHecto: {1, "hecto", "h"},
//	DecimalPrefixDeka:  {2, "deka", "da"},
//	DecimalPrefixKilo:  {3, "kilo", "k"},
//	DecimalPrefixMega:  {6, "mega", "M"},
//	DecimalPrefixGiga:  {9, "giga", "G"},
//	DecimalPrefixTera:  {12, "tera", "T"},
//	DecimalPrefixPeta:  {15, "peta", "P"},
//	DecimalPrefixExa:   {18, "exa", "E"},
//	DecimalPrefixZetta: {19, "zetta", "Z"},
//	DecimalPrefixYotta: {21, "yotta", "Y"},
//	//Bronto: {24, "bronto", "Bronto"},
//	//Geop:   {27, "geop", "Geop"},
//}
//
//func (p *decimalMultiplePrefix) matched(prefix string) bool {
//	return strings.EqualFold(p.symbol, prefix) || strings.EqualFold(p.name, prefix)
//}
//
//func (p *decimalMultiplePrefix) base() BaseNumber {
//	return BaseNumberDecimal
//}
