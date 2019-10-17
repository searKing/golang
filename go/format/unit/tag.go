package unit

import "strings"

type tagPair struct {
	long  string // full unit tag
	short string // short unit tag
}

func (p *tagPair) Matched(tag string) bool {
	return strings.EqualFold(p.short, tag) || strings.EqualFold(p.long, tag)
}

var unitTagPairs = map[Unit]tagPair{
	One:    {"One", ""},
	Kilo:   {"Kilo", "K"},
	Mega:   {"Mega", "M"},
	Giga:   {"Giga", "G"},
	Tera:   {"Tera", "T"},
	Peta:   {"Peta", "P"},
	Exa:    {"Exa", "E"},
	Zetta:  {"Zetta", "Z"},
	Yotta:  {"Yotta", "Y"},
	Bronto: {"Bronto", "B"},
	Geop:   {"Geop", "Geop"},
}
