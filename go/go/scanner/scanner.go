// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// A mode value is a set of flags (or 0).
// They control scanner behavior.
type Mode uint

const (
	ModeCaseSensitive Mode = 1 << iota
	ModeRegexpPerl
	ModeRegexpPosix
)

// An ErrorHandler may be provided to Scanner.Init. If a syntax error is
// encountered and a handler was installed, the handler is called with a
// position and an error message. The position points to the beginning of
// the offending token.
type ErrorHandler func(pos token.Position, msg string)

// A Scanner holds the scanner's internal state while processing
// a given text. It can be allocated as part of another data
// structure but must be initialized via Init before use.
type Scanner struct {
	// immutable state
	file *token.File  // source file handle
	dir  string       // directory portion of file.Name()
	src  []byte       // source
	err  ErrorHandler // error reporting; or nil
	mode Mode         // scanning mode

	// scanning state
	offset     int // character offset
	rdOffset   int // reading offset (position after current character)
	lineOffset int // current line offset

	// public state - ok to modify
	ErrorCount int // number of errors encountered
}

const bom = 0xFEFF // byte order mark, only permitted as very first character

func (s *Scanner) AtEOF() bool {
	return s.rdOffset >= len(s.src)
}

func (s *Scanner) CurrentBytes() []byte {
	return s.src[s.offset:s.rdOffset]
}

func (s *Scanner) CurrentString() string {
	return string(s.CurrentBytes())
}

func (s *Scanner) CurrentRunes() []rune {
	return []rune(s.CurrentString())
}

func (s *Scanner) CurrentRune() rune {
	runes := s.CurrentRunes()
	if len(runes) > 0 {
		return runes[0]
	}
	return -1
}

func (s *Scanner) CurrentLength() int {
	return s.rdOffset - s.offset
}

// walk until current is consumed
func (s *Scanner) Consume() {
	chars := s.CurrentBytes()
	if len(chars) == 0 {
		return
	}

	lines := bytes.Split(chars, []byte{'\n'})
	var hasCL bool
	if len(lines) > 1 {
		hasCL = true
	}

	for _, line := range lines {
		lineLen := len(line)
		if hasCL {
			lineLen++
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}

		s.offset = s.offset + lineLen
	}
	s.offset = s.rdOffset
}

func (s *Scanner) NextByte() {
	s.NextBytesN(1)
}

func (s *Scanner) NextBytesN(n int) {
	s.Consume()
	if s.rdOffset+n <= len(s.src) {
		s.rdOffset += n
	} else {
		s.offset = len(s.src)
	}
}

// Read the NextRune Unicode char into s.ch.
// s.AtEOF() == true means end-of-file.
func (s *Scanner) NextRune() {
	if s.rdOffset < len(s.src) {
		s.Consume()
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
			}
		}
		s.rdOffset += w
	} else {
		s.Consume()
		s.offset = len(s.src)
	}
}

func (s *Scanner) PeekRune() rune {
	if s.rdOffset < len(s.src) {
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
			}
		}
		return r
	}
	return -1
}

// PeekByte returns the byte following the most recently read character without
// advancing the scanner. If the scanner is at EOF, PeekByte returns 0.
func (s *Scanner) PeekByte() byte {
	if s.rdOffset < len(s.src) {
		return s.src[s.rdOffset]
	}
	return 0
}

// Read the NextRune Unicode chars into s.ch.
// s.ch < 0 means end-of-file.
func (s *Scanner) NextRunesN(n int) {
	offsetBegin := s.rdOffset

	for i := 0; i < n; i++ {
		s.NextRune()
	}
	s.offset = offsetBegin
}

// Read the NextRune Unicode chars into s.ch.
// s.ch < 0 means end-of-file.
func (s *Scanner) NextRegexp(expectStrs ...string) {
	match := s.PeekRegexpAny()
	if match == "" {
		return
	}
	offsetBegin := s.rdOffset

	for range match {
		s.NextRune()
	}
	s.offset = offsetBegin
}

// PeekRegexpAny returns the string following the most recently read character which matches the regexp case without
// advancing the scanner. If the scanner is at EOF or regexp unmatched, PeekRegexpAny returns nil.
func (s *Scanner) PeekRegexpAny(expectStrs ...string) string {
	if s.AtEOF() {
		return ""
	}
	if s.mode&ModeRegexpPosix != 0 {
		return s.peekRegexpPosix(expectStrs...)
	} else if s.mode&ModeRegexpPerl != 0 {
		return s.peekRegexpPerl(expectStrs...)
	}

	return s.PeekString(expectStrs...)
}

func (s *Scanner) PeekString(expectStrs ...string) string {
	if s.AtEOF() {
		return ""
	}

	// regex mode
	for _, expect := range expectStrs {
		endPos := s.rdOffset + len(expect)
		if endPos > len(s.src) {
			continue
		}
		selected := s.src[s.rdOffset:endPos]
		if string(selected) == expect {
			return string(selected)
		}

		if ((s.mode&ModeCaseSensitive != 0) && strings.EqualFold(string(selected), expect)) ||
			string(selected) == expect {
			return string(selected)
		}
	}
	return ""
}

func (s *Scanner) peekRegexpPosix(expectStrs ...string) string {
	if s.AtEOF() {
		return ""
	}

	// regex mode
	for _, expect := range expectStrs {
		expect = "^" + strings.TrimPrefix(expect, "^")

		reg := regexp.MustCompilePOSIX(expect)
		matches := reg.FindStringSubmatch(string(s.src[s.rdOffset:]))
		if len(matches) == 0 {
			continue
		}

		return matches[0]
	}
	return ""
}

func (s *Scanner) peekRegexpPerl(expectStrs ...string) string {
	if s.AtEOF() {
		return ""
	}

	// regex mode
	for _, expect := range expectStrs {
		expect = "^" + strings.TrimPrefix(expect, "^")

		reg := regexp.MustCompile(expect)
		matches := reg.FindStringSubmatch(string(s.src[s.rdOffset:]))
		if len(matches) == 0 {
			continue
		}

		return matches[0]
	}
	return ""
}

// Init prepares the scanner s to tokenize the text src by setting the
// scanner at the beginning of src. The scanner uses the file set file
// for position information and it adds line information for each line.
// It is ok to re-use the same file when re-scanning the same file as
// line information which is already present is ignored. Init causes a
// panic if the file size does not match the src size.
//
// Calls to Scan will invoke the error handler err if they encounter a
// syntax error and err is not nil. Also, for each error encountered,
// the Scanner field ErrorCount is incremented by one. The mode parameter
// determines how comments are handled.
//
// Note that Init may call err if there is an error in the first character
// of the file.
func (s *Scanner) Init(file *token.File, src []byte, err ErrorHandler, mode Mode) {
	// Explicitly initialize all fields since a scanner may be reused.
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}
	s.file = file
	s.dir, _ = filepath.Split(file.Name())
	s.src = src
	s.err = err
	s.mode = mode

	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.ErrorCount = 0

	if s.PeekRune() == bom {
		s.NextRune() // ignore BOM at file beginning
	}
}

func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(s.file.Position(s.file.Pos(offs)), msg)
	}
	s.ErrorCount++
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}

// ScanEscape parses an escape sequence where rune is the accepted
// escaped quote. In case of a syntax error, it stops at the offending
// character (without consuming it) and returns false. Otherwise
// it returns true.
func (s *Scanner) ScanEscape(quote rune) bool {
	offs := s.offset

	var ch = s.CurrentRune()

	var n int
	var base, max uint32
	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.NextRune()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.NextRune()
		n, base, max = 2, 16, 255
	case 'u':
		s.NextRune()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.NextRune()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if ch < 0 {
			msg = "escape sequence not terminated"
		}
		s.error(offs, msg)
		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", ch)
			if ch < 0 {
				msg = "escape sequence not terminated"
			}
			s.error(s.offset, msg)
			return false
		}
		x = x*base + d
		s.NextRune()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		s.error(offs, "escape sequence is invalid Unicode code point")
		return false
	}

	return true
}

func (s *Scanner) ScanRune() string {
	// '\'' opening already consumed
	offs := s.offset - 1

	valid := true
	n := 0
	for {
		var ch = s.CurrentRune()

		if ch == '\n' || ch < 0 {
			// only report error if we don't have one already
			if valid {
				s.error(offs, "rune literal not terminated")
				valid = false
			}
			break
		}
		s.NextRune()
		if ch == '\'' {
			break
		}
		n++
		if ch == '\\' {
			if !s.ScanEscape('\'') {
				valid = false
			}
			// continue to read to closing quote
		}
	}

	if valid && n != 1 {
		s.error(offs, "illegal rune literal")
	}

	return string(s.src[offs:s.offset])
}

func (s *Scanner) ScanString() string {
	// '"' opening already consumed
	offs := s.offset - 1

	for {
		var ch = s.CurrentRune()
		if ch == '\n' || ch < 0 {
			s.error(offs, "string literal not terminated")
			break
		}
		s.NextRune()
		if ch == '"' {
			break
		}
		if ch == '\\' {
			s.ScanEscape('"')
		}
	}

	return string(s.src[offs:s.offset])
}

func stripCR(b []byte, comment bool) []byte {
	c := make([]byte, len(b))
	i := 0
	for j, ch := range b {
		// In a /*-style comment, don't strip \r from *\r/ (incl.
		// sequences of \r from *\r\r...\r/) since the resulting
		// */ would terminate the comment too early unless the \r
		// is immediately following the opening /* in which case
		// it's ok because /*/ is not closed yet (issue #11151).
		if ch != '\r' || comment && i > len("/*") && c[i-1] == '*' && j+1 < len(b) && b[j+1] == '/' {
			c[i] = ch
			i++
		}
	}
	return c[:i]
}

func (s *Scanner) ScanRawString() string {
	// '`' opening already consumed
	offs := s.offset - 1

	hasCR := false
	for {
		var ch = s.CurrentRune()
		if ch < 0 {
			s.error(offs, "raw string literal not terminated")
			break
		}
		s.NextRune()
		if ch == '`' {
			break
		}
		if ch == '\r' {
			hasCR = true
		}
	}

	lit := s.src[offs:s.offset]
	if hasCR {
		lit = stripCR(lit, false)
	}

	return string(lit)
}

func (s *Scanner) ScanLine() string {
	// '"' opening already consumed
	offs := s.offset

	for {
		var ch = s.CurrentRune()
		if ch < 0 {
			s.error(offs, "string literal not terminated")
			break
		}
		s.NextRune()
		if ch == '\n' {
			break
		}
	}

	return string(s.src[offs:s.offset])
}

// ScanSplits advances the Scanner to the next token by splits when first meet, which will then be
// available through the Bytes or Text method. It returns false when the
// scan stops, either by reaching the end of the input or an error.
// After Scan returns false, the Err method will return any error that
// occurred during scanning, except that if it was io.EOF, Err
// will return nil.
func (s *Scanner) ScanSplits(splits ...bufio.SplitFunc) ([]byte, bool) {
	s.Consume()

	for _, split := range splits {
		if split == nil {
			continue
		}
		// See if we can get a token with what we already have.
		// If we've run out of data but have an error, give the split function
		// a chance to recover any remaining, possibly empty token.
		// atEOF is true always, for we consume by a byte slice
		advance, token, err := split(s.src[s.rdOffset:], true)
		if err != nil && err != bufio.ErrFinalToken {
			s.error(s.offset, err.Error())
			return nil, false
		}
		s.NextBytesN(advance)
		if len(token) != 0 {
			return token, true
		}
	}
	return nil, false
}
