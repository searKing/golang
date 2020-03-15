// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Split functions
var (
	// ScanBytes is a split function for a Scanner that returns each byte as a token.
	ScanBytes = bufio.ScanBytes

	// ScanRunes is a split function for a Scanner that returns each
	// UTF-8-encoded rune as a token. The sequence of runes returned is
	// equivalent to that from a range loop over the input as a string, which
	// means that erroneous UTF-8 encodings translate to U+FFFD = "\xef\xbf\xbd".
	// Because of the Scan interface, this makes it impossible for the client to
	// distinguish correctly encoded replacement runes from encoding errors.
	ScanRunes = bufio.ScanRunes

	// ScanWords is a split function for a Scanner that returns each
	// space-separated word of text, with surrounding spaces deleted. It will
	// never return an empty string. The definition of space is set by
	// unicode.IsSpace.
	ScanWords = bufio.ScanWords

	// ScanLines is a split function for a Scanner that returns each line of
	// text, stripped of any trailing end-of-line marker. The returned line may
	// be empty. The end-of-line marker is one optional carriage return followed
	// by one mandatory newline. In regular expression notation, it is `\r?\n`.
	// The last non-empty line of input will be returned even if it has no
	// newline.
	ScanLines = bufio.ScanLines
)

// ScanRawStrings is a split function for a Scanner that returns each string quoted by ` of
// text. The returned line may be empty. Escape is disallowed
// Raw string literals are character sequences between back quotes, as in `foo`.
// Within the quotes, any character may appear except back quote.
// The value of a raw string literal is the string composed of the uninterpreted (implicitly UTF-8-encoded) characters
// between the quotes; in particular, backslashes have no special meaning and the string may contain newlines.
// Carriage return characters ('\r') inside raw string literals are discarded from the raw string value.
// https://golang.org/ref/spec#String_literals
// raw_string_lit         = "`" { unicode_char | newline } "`" .
func ScanRawStrings(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return scanStrings(data, atEOF, '`')
}

// ScanInterpretedStrings is a split function for a Scanner that returns each string quoted by " of
// text. The returned line may be empty.
// Interpreted string literals are character sequences between double quotes, as in "bar".
// Within the quotes, any character may appear except newline and unescaped double quote.
// The text between the quotes forms the value of the literal,
// with backslash escapes interpreted as they are in rune literals
// (except that \' is illegal and \" is legal), with the same restrictions.
// The three-digit octal (\nnn) and two-digit hexadecimal (\xnn)
// escapes represent individual bytes of the resulting string;
// all other escapes represent the (possibly multi-byte) UTF-8 encoding of individual characters. Thus inside a string
// literal \377 and \xFF represent a single byte of value 0xFF=255, while ÿ, \u00FF, \U000000FF and \xc3\xbf represent
// the two bytes 0xc3 0xbf of the UTF-8 encoding of character U+00FF.
// https://golang.org/ref/spec#String_literals
// interpreted_string_lit = `"` { unicode_value | byte_value } `"` .
func ScanInterpretedStrings(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return scanStrings(data, atEOF, '"')
}

// ScanEscapes is a split function wrapper for a Scanner that returns each string which is an escape format of
// text. The returned line may be empty.
func ScanEscapes(quote rune) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return scanEscapes(data, atEOF, quote)
	}
}

// ScanMantissas is a split function wrapper for a Scanner that returns each string which is an n-base number format of
// text. The returned line may be empty.
func ScanMantissas(base int) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return ScanWhile(func(r rune) bool {
		return digitVal(r) < base
	})
}

// https://golang.org/ref/spec#Integer_literals
// https://golang.org/ref/spec#Floating-point_literals
// https://golang.org/ref/spec#Imaginary_literals
// ScanNumbers is a split function wrapper for a Scanner that returns each string which is an integer, floating-point
// or imaginary format of text. The returned line may be empty.
func ScanNumbers(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return needMoreData()
	}
	var off int
	var seenSign bool
	var seenDecimalPoint bool
	var seenDecimalNumber bool

	var lookforFraction bool
	var lookforExponent bool
	// First character 1: digitVal(ch) < 10.
	// Handle .989 or 0x888
	for {
		// read a rune
		advance, token, err := handleSplitError(ScanRunes(data[off:], atEOF))
		off = off + advance
		if err != nil || len(token) == 0 {
			return advance, token, err
		}
		ch := bytes.Runes(token)[0]
		if ch == '.' {
			// . can be seen once only
			if seenDecimalPoint {
				off--
				return off, data[:off], nil
			}
			seenDecimalPoint = true
			continue
		}

		// sign can be seen leading or after E or e
		if ch == '+' || ch == '-' {
			// sign can be seen once only, and can never be after "."
			if seenSign || seenDecimalPoint {
				off--
				return off, data[:off], nil
			}
			seenSign = true
			continue
		}

		// number must be leading with "." "+" "-" or "0-9"
		if !seenDecimalNumber && digitVal(ch) > 10 {
			msg := fmt.Sprintf("illegal character %#U leading escape sequence, expect \\", token)
			return 0, nil, errors.New(msg)
		}
		seenDecimalNumber = true

		// .989777
		if seenDecimalPoint {
			advance, token, err := handleSplitError(ScanMantissas(10)(data[off:], atEOF))
			off = off + advance
			if err != nil || len(token) == 0 {
				return advance, token, err
			}
			// look for "E" or "e"
			lookforExponent = true
			break
		}

		// 0x12
		if ch == '0' {
			// int or float
			advance, token, err := handleSplitError(ScanRunes(data[off:], atEOF))
			off = off + advance
			if err != nil {
				return advance, token, err
			}
			if len(token) == 0 {
				return off, data[:off], nil
			}
			ch = bytes.Runes(token)[0]

			if ch == 'x' || ch == 'X' {
				// hexadecimal int
				advance, token, err := handleSplitError(ScanMantissas(16)(data[off:], atEOF))
				off = off + advance
				if err != nil || len(token) == 0 {
					return advance, token, err
				}
				if len(token) <= 0 {
					// only scanned "0x" or "0X"
					return 0, nil, errors.New("illegal hexadecimal number")
				}
				return off, data[:off], nil
			} else {
				// octal int or float
				seenDecimalDigit := false
				advance, token, err := handleSplitError(ScanMantissas(8)(data[off:], atEOF))
				off = off + advance
				if err != nil {
					return advance, token, err
				}

				// read new rune
				advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
				off = off + advance
				if err != nil {
					return advance, token, err
				}
				if len(token) == 0 {
					return off, data[:off], nil
				}
				ch = bytes.Runes(token)[0]

				if ch == '8' || ch == '9' {
					// illegal octal int or float
					seenDecimalDigit = true
					advance, token, err := handleSplitError(ScanMantissas(10)(data[off:], atEOF))
					off = off + advance
					if err != nil || len(token) == 0 {
						return advance, token, err
					}
					advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
					off = off + advance
					if err != nil || len(token) == 0 {
						return advance, token, err
					}
					ch = bytes.Runes(token)[0]
				}
				if ch == '.' || ch == 'e' || ch == 'E' || ch == 'i' {
					off-- //backward for fraction "." "e" "E" or "i"
					lookforFraction = true
					break
				}
				// octal int
				if seenDecimalDigit {
					return 0, nil, errors.New("illegal octal number")
				}

				off-- //backward for exit

			}
			return off, data[:off], nil
		}

		// decimal int or float
		advance, token, err = handleSplitError(ScanMantissas(10)(data[off:], atEOF))
		off = off + advance
		if err != nil || len(token) == 0 {
			return advance, token, err
		}
		lookforFraction = true
		break
	}

	// read a rune
	advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
	off = off + advance
	if err != nil {
		return advance, token, err
	}
	if len(token) == 0 {
		return off, data[:off], nil
	}
	ch := bytes.Runes(token)[0]

	if lookforFraction && ch == '.' {
		advance, token, err := handleSplitError(ScanMantissas(10)(data[off:], atEOF))
		off = off + advance
		if err != nil {
			return advance, token, err
		}
		if len(token) == 0 {
			return off, data[:off], nil
		}
		lookforExponent = true

		// read new rune
		advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
		off = off + advance
		if err != nil {
			return advance, token, err
		}
		if len(token) == 0 {
			return off, data[:off], nil
		}
		ch = bytes.Runes(token)[0]
	}

	if lookforExponent && (ch == 'e' || ch == 'E') {
		advance, token, err := handleSplitError(ScanRunes(data[off:], atEOF))
		off = off + advance
		if err != nil {
			return advance, token, err
		}
		if len(token) == 0 {
			return off, data[:off], nil
		}
		ch = bytes.Runes(token)[0]

		if ch == '-' || ch == '+' {
			advance, token, err := handleSplitError(ScanRunes(data[off:], atEOF))
			off = off + advance
			if err != nil {
				return advance, token, err
			}
			if len(token) == 0 {
				return off, data[:off], nil
			}
			ch = bytes.Runes(token)[0]
		}
		if digitVal(ch) < 10 {
			advance, token, err := handleSplitError(ScanMantissas(10)(data[off:], atEOF))
			off = off + advance
			if err != nil {
				return advance, token, err
			}
			if len(token) == 0 {
				return off, data[:off], nil
			}

			// read new rune
			advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
			off = off + advance
			if err != nil {
				return advance, token, err
			}
			if len(token) == 0 {
				return off, data[:off], nil
			}
		} else {
			return 0, nil, errors.New("illegal floating-point exponent")
		}
	}

	if ch != 'i' {
		// backward
		off = off - utf8.RuneLen(ch)
	}
	return off, data[:off], nil
}

// https://golang.org/ref/spec#Identifiers
// ScanIdentifier is a split function wrapper for a Scanner that returns each string which is an identifier format of text.
// The returned line may be empty.
// identifier = letter { letter | unicode_digit } .
func ScanIdentifier(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return needMoreData()
	}
	var off int

	// First character 1: \.
	advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
	off = off + advance
	if err != nil || len(token) == 0 {
		return advance, token, err
	}
	ch := bytes.Runes(token)[0]

	if isLetter(ch) {
		for isLetter(ch) || isDigit(ch) {
			advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
			off = off + advance
			if err != nil {
				return advance, token, err
			}
			if token == nil {
				return off, data[:off], nil
			}
			ch = bytes.Runes(token)[0]
		}
	}
	off -= utf8.RuneLen(ch) // backward
	return off, data[:off], nil
}

// ScanUntil is a split function wrapper for a Scanner that returns each string until filter case is meet.
// The returned line may be empty.
func ScanUntil(filter func(r rune) bool) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return ScanWhile(func(r rune) bool {
		if filter == nil {
			return false
		}
		return !filter(r)
	})
}

// ScanUntil is a split function wrapper for a Scanner that returns each string until filter case is not meet.
// The returned line may be empty.
func ScanWhile(filter func(r rune) bool) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if filter == nil || atEOF && len(data) == 0 {
			return needMoreData()
		}
		var off int

		// First character 1: \.
		advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
		off = off + advance
		if err != nil || len(token) == 0 {
			return advance, token, err
		}
		ch := bytes.Runes(token)[0]

		for filter(ch) {
			advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
			off = off + advance
			if err != nil {
				return advance, token, err
			}
			if token == nil {
				return off, data[:off], nil
			}
			ch = bytes.Runes(token)[0]
		}
		off -= utf8.RuneLen(ch) // backward

		return off, data[:off], nil
	}
}

// ScanRegexp is a split function wrapper for a Scanner that returns each string until regexp case is not meet.
// The returned line may be empty.
func ScanRegexp(regs ...*regexp.Regexp) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return needMoreData()
		}
		var off int

		// First character 1: \.
		// regex mode
		for _, reg := range regs {
			if reg == nil {
				continue
			}

			locs := reg.FindStringIndex(string(data[off:]))
			if len(locs) == 0 {
				continue
			}
			off = locs[1]
			return off, data[locs[0]:off], nil
		}

		return off, data[:off], nil
	}
}

// ScanRegexpPerl is a split function wrapper for a Scanner that returns each string until regexp case is not meet.
// The returned line may be empty.
// This so-called leftmost-first matching is the same semantics
// that Perl, Python, and other implementations use, although this
// package implements it without the expense of backtracking.
// For POSIX leftmost-longest matching, see ScanRegexpPosix.
func ScanRegexpPerl(expectStrs ...string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var regs []*regexp.Regexp
	for _, expect := range expectStrs {
		expect = "^" + strings.TrimPrefix(expect, "^")

		regs = append(regs, regexp.MustCompile(expect))
	}
	return ScanRegexp(regs...)
}

// ScanRegexpPosix is a split function wrapper for a Scanner that returns each string until regexp case is not meet.
// The returned line may be empty.
// ScanRegexpPosix is like ScanRegexpPerl but restricts the regular expression
// to POSIX ERE (egrep) syntax and changes the match semantics to
// leftmost-longest.
func ScanRegexpPosix(expectStrs ...string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var regs []*regexp.Regexp
	for _, expect := range expectStrs {
		expect = "^" + strings.TrimPrefix(expect, "^")

		regs = append(regs, regexp.MustCompilePOSIX(expect))
	}
	return ScanRegexp(regs...)
}

// https://golang.org/ref/spec#String_literals
// string_lit             = raw_string_lit | interpreted_string_lit .
// raw_string_lit         = "`" { unicode_char | newline } "`" .
// interpreted_string_lit = `"` { unicode_value | byte_value } `"` .
func scanStrings(data []byte, atEOF bool, quote rune) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return needMoreData()
	}
	var off int

	// First character 1: ".
	advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
	off = off + advance
	if err != nil || len(token) == 0 {
		return advance, token, err
	}
	if !bytes.ContainsRune(token, quote) {
		msg := fmt.Sprintf("illegal character %#U leading escape sequence, expect \\", token)
		return 0, nil, errors.New(msg)
	}

	var allowEscape bool
	if quote == '"' {
		allowEscape = true
	}
	// '"' opening already consumed
	for _, ch := range data[off:] {
		off++
		if ch == '\n' || ch < 0 {
			return 0, nil, errors.New("string literal not terminated")
		}

		if rune(ch) == quote {
			break
		}

		if allowEscape && ch == '\\' {
			// backward
			off--
			advance, token, err = handleSplitError(ScanEscapes(quote)(data[off:], atEOF))
			off = off + advance
			if err != nil || len(token) == 0 {
				return advance, token, err
			}

		}
	}
	return off, data[:off], nil
}

func scanEscapes(data []byte, atEOF bool, quote rune) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return needMoreData()
	}
	var off int

	// First character 1: \.
	advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
	off = off + advance
	if err != nil || len(token) == 0 {
		return advance, token, err
	}

	if !bytes.ContainsRune(token, '\\') {
		msg := fmt.Sprintf("illegal character %#U leading escape sequence, expect \\", token)
		return 0, nil, errors.New(msg)
	}

	// Second character 2: char.
	advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
	off = off + advance
	if err != nil || len(token) == 0 {
		return advance, token, err
	}

	ch := bytes.Runes(token)[0]

	var n int
	var base, max uint32
	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		return off, data[0:off], nil
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		n, base, max = 2, 16, 255
	case 'u':
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if ch < 0 {
			msg = "escape sequence not terminated"
		}
		return 0, nil, errors.New(msg)
	}

	switch ch {
	case 'x', 'u', 'U':
		advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
		off = off + advance
		if err != nil || len(token) == 0 {
			return advance, token, err
		}

		ch = bytes.Runes(token)[0]
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", ch)
			if ch < 0 {
				msg = "escape sequence not terminated"
			}
			return 0, nil, errors.New(msg)
		}
		x = x*base + d

		advance, token, err = handleSplitError(ScanRunes(data[off:], atEOF))
		off = off + advance
		if err != nil || len(token) == 0 {
			return advance, token, err
		}
		ch = bytes.Runes(token)[0]

		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		return 0, nil, errors.New("escape sequence is invalid Unicode code point")
	}
	return off, data[:off], nil
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func needMoreData() (advance int, token []byte, err error) {
	return 0, nil, nil
}

func handleSplitError(advance int, token []byte, err error) (int, []byte, error) {
	if err != nil {
		if err == bufio.ErrFinalToken {
			return 0, nil, nil
		}
		return 0, nil, err
	}

	if len(token) == 0 {
		// needMoreData
		return 0, nil, nil
	}

	return advance, token, nil
}
