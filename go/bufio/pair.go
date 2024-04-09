// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bufio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/searKing/golang/go/container/stack"
	slices_ "github.com/searKing/golang/go/exp/slices"
)

var (
	ErrMismatchTokenPair = errors.New("mismatch token pair")
	ErrInvalidStartToken = errors.New("invalid start token")
)

type DelimiterPair struct {
	start byte
	end   byte
}

// A PairScanner reads and decodes Pair wrapped values from an input stream, like a json、xml.
type PairScanner struct {
	r              io.Reader
	discardLeading bool // discard any char until we meet a start delimiter at list
	buf            []byte
	scanp          int   // start of unread data in buf
	scanned        int64 // amount of data already scanned
	err            error
}

// NewPairScanner returns a new scanner that reads from r.
//
// The scanner introduces its own buffering and may
// read data from r beyond the paired values requested.
func NewPairScanner(r io.Reader) *PairScanner {
	return &PairScanner{r: r}
}

func (s *PairScanner) SetDiscardLeading(discard bool) *PairScanner {
	s.discardLeading = discard
	return s
}

func (s *PairScanner) ScanDelimiters(delimiters string) (line []byte, err error) {
	var pairs []DelimiterPair
	var isPair bool
	var lastDelimiter byte

	for _, delimiter := range []byte(delimiters) {
		if !isPair {
			lastDelimiter = delimiter
			isPair = true
			continue
		}
		pairs = append(pairs, DelimiterPair{
			lastDelimiter, delimiter,
		})
		isPair = false
	}

	return s.Scan(pairs)

}

// Scan reads the next value complete wrapped by pair delimiters from its
// input and stores it in the value pointed to by v.
func (s *PairScanner) Scan(pairs []DelimiterPair) (line []byte, err error) {
	if s.err != nil {
		return nil, s.err
	}

	// Read whole value into buffer.
	n, err := s.readValue(pairs)
	if err != nil {
		return nil, err
	}
	line = s.buf[s.scanp : s.scanp+n]
	s.scanp += n

	return line, nil
}

// Buffered returns a reader of the data remaining in the PairScanner's
// buffer. The reader is valid until the next call to Decode.
func (s *PairScanner) Buffered() io.Reader {
	return bytes.NewReader(s.buf[s.scanp:])
}

// readValue reads a JSON value into dec.buf.
// It returns the length of the encoding.
func (s *PairScanner) readValue(pairs []DelimiterPair) (int, error) {
	var delimiters stack.Stack
	scanp := s.scanp
	var err error
Input:
	// help the compiler see that scanp is never negative, so it can remove
	// some bounds checks below.
	for scanp >= 0 {

		// Look in the buffer for a new value.
		for i, c := range s.buf[scanp:] {
			delimiterPair, ok := findMatchedTokenPair(c, pairs)
			if !ok {
				if delimiters.Len() == 0 {
					// no delimiter have been seen yet
					// discard any char until we meet a start delimiter at list
					if !s.discardLeading {
						return 0, s.tokenError(c, ErrInvalidStartToken)
					}
					s.scanp++
				}
				// read next char
				continue
			}
			if c == delimiterPair.start {
				delimiters.Push(c)
				continue
			} // c == delimiterPair.end
			// no delimiter have been seen yet
			if delimiters.Len() == 0 {
				// discard any char until we meet a start delimiter at list
				if !s.discardLeading {
					return 0, s.tokenError(c, ErrInvalidStartToken)
				}
				s.scanp++
				continue
			}

			lastDelimiter := delimiters.Peek().Value.(byte)
			if lastDelimiter != delimiterPair.start {
				return 0, s.tokenError(c, ErrMismatchTokenPair)
			}
			delimiters.Pop()
			// a perfect object is get, just return
			if delimiters.Len() == 0 {
				scanp += i + 1
				break Input
			}

		}
		scanp = len(s.buf)

		// Did the last read have an error?
		// Delayed until now to allow buffer scan.
		if err != nil {
			if err == io.EOF {
				if len(s.buf) > 0 {
					err = io.ErrUnexpectedEOF
				}
			}
			s.err = err
			return 0, err
		}

		n := scanp - s.scanp
		err = s.refill()
		scanp = s.scanp + n
	}
	return scanp - s.scanp, nil
}

func (s *PairScanner) refill() error {
	// Make room to read more into the buffer.
	// First slide down data already consumed.
	if s.scanp > 0 {
		s.scanned += int64(s.scanp)
		n := copy(s.buf, s.buf[s.scanp:])
		s.buf = s.buf[:n]
		s.scanp = 0
	}

	// Grow buffer if not large enough.
	const minRead = 512
	if cap(s.buf)-len(s.buf) < minRead {
		newBuf := make([]byte, len(s.buf), 2*cap(s.buf)+minRead)
		copy(newBuf, s.buf)
		s.buf = newBuf
	}

	// Read. Delay error for next iteration (after scan).
	n, err := s.r.Read(s.buf[len(s.buf):cap(s.buf)])
	s.buf = s.buf[0 : len(s.buf)+n]

	return err
}

func (s *PairScanner) tokenError(c byte, err error) error {
	return fmt.Errorf("invalid character %s at %d: %w", quoteChar(c), s.InputOffset(), err)
}

// More reports whether there is another element in the
// current array or object being parsed.
func (s *PairScanner) More() bool {
	c, err := s.peek()
	return err == nil && c != ']' && c != '}'
}

func (s *PairScanner) peek() (byte, error) {
	var err error
	for {
		if s.scanp < len(s.buf) {
			c := s.buf[s.scanp]
			return c, nil
		}
		// buffer has been scanned, now report any error
		if err != nil {
			return 0, err
		}
		err = s.refill()
	}
}

// InputOffset returns the input stream byte offset of the current scanner position.
// The offset gives the location of the end of the most recently returned token
// and the beginning of the next token.
func (s *PairScanner) InputOffset() int64 {
	return s.scanned + int64(s.scanp)
}

func findMatchedTokenPair(c byte, pairs []DelimiterPair) (tokenPair DelimiterPair, has bool) {
	return slices_.FirstFunc(pairs, func(pair DelimiterPair) bool {
		return c == pair.start || c == pair.end
	})
}

// quoteChar formats c as a quoted character literal.
func quoteChar(c byte) string {
	// special cases - different from quoted strings
	if c == '\'' {
		return `'\''`
	}
	if c == '"' {
		return `'"'`
	}

	// use quoted string with different quotation marks
	s := strconv.Quote(string(c))
	return "'" + s[1:len(s)-1] + "'"
}
