// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bufio

import (
	"bytes"
	"io"
	"unicode"

	"github.com/pkg/errors"
	"github.com/searKing/golang/go/container/slice"
	"github.com/searKing/golang/go/container/stack"
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
	discardLeading bool // discard any char until we meet a start delimeter at list
	buf            []byte
	scanp          int   // start of unread data in buf
	scanned        int64 // amount of data already scanned
	err            error
}

// NewPairScanner returns a new scanner that reads from r.
//
// The scanner introduces its own buffering and may
// read data from r beyond the JSON values requested.
func NewPairScanner(r io.Reader) *PairScanner {
	return &PairScanner{r: r}
}

func (pairScanner *PairScanner) SetDiscardLeading(discard bool) *PairScanner {
	pairScanner.discardLeading = discard
	return pairScanner
}

func (pairScanner *PairScanner) ScanDelimiters(delimiters string) (line []byte, err error) {
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

	return pairScanner.Scan(pairs)

}

// Scan reads the next value complete wrapped by pair delimiters from its
// input and stores it in the value pointed to by v.
func (pairScanner *PairScanner) Scan(pairs []DelimiterPair) (line []byte, err error) {
	if pairScanner.err != nil {
		return nil, pairScanner.err
	}

	// Read whole value into buffer.
	n, err := pairScanner.readValue(pairs)
	if err != nil {
		return nil, err
	}
	line = pairScanner.buf[pairScanner.scanp : pairScanner.scanp+n]
	pairScanner.scanp += n

	return line, nil
}

// Buffered returns a reader of the data remaining in the PairScanner's
// buffer. The reader is valid until the next call to Decode.
func (pairScanner *PairScanner) Buffered() io.Reader {
	return bytes.NewReader(pairScanner.buf[pairScanner.scanp:])
}

// readValue reads a JSON value into dec.buf.
// It returns the length of the encoding.
func (pairScanner *PairScanner) readValue(pairs []DelimiterPair) (int, error) {
	var delimiters stack.Stack
	scanp := pairScanner.scanp
	var err error
Input:
	for {
		// Look in the buffer for a new value.
		for i, c := range pairScanner.buf[scanp:] {
			delimiterPair, ok := findMatchedTokenPair(c, pairs)
			if !ok && delimiters.Len() == 0 {
				// no delimiter have been seen yet
				// discard any char until we meet a start delimeter at list
				if pairScanner.discardLeading {
					pairScanner.scanp += 1
					continue
				}
				continue
			}
			if !ok {
				// read next char
				continue
			}
			if c == delimiterPair.start {
				delimiters.Push(c)
			} else { //c == delimiterPair.end
				// no delimiter have been seen yet
				if delimiters.Len() == 0 {
					// discard any char until we meet a start delimeter at list
					pairScanner.scanp += 1
					if pairScanner.discardLeading {
						continue
					}
					return 0, ErrInvalidStartToken
				}

				lastDelimiter := delimiters.Peek().Value.(byte)
				if lastDelimiter != delimiterPair.start {
					return 0, ErrMismatchTokenPair
				}
				delimiters.Pop()
				// a perfect object is get, just return
				if delimiters.Len() == 0 {
					scanp += i + 1
					break Input
				}
			}
		}
		scanp = len(pairScanner.buf)

		// Did the last read have an error?
		// Delayed until now to allow buffer scan.
		if err != nil {
			if err == io.EOF {
				if nonSpace(pairScanner.buf) {
					err = io.ErrUnexpectedEOF
				}
			}
			pairScanner.err = err
			return 0, err
		}

		n := scanp - pairScanner.scanp
		err = pairScanner.refill()
		scanp = pairScanner.scanp + n
	}
	return scanp - pairScanner.scanp, nil
}

func (pairScanner *PairScanner) refill() error {
	// Make room to read more into the buffer.
	// First slide down data already consumed.
	if pairScanner.scanp > 0 {
		pairScanner.scanned += int64(pairScanner.scanp)
		n := copy(pairScanner.buf, pairScanner.buf[pairScanner.scanp:])
		pairScanner.buf = pairScanner.buf[:n]
		pairScanner.scanp = 0
	}

	// Grow buffer if not large enough.
	const minRead = 512
	if cap(pairScanner.buf)-len(pairScanner.buf) < minRead {
		newBuf := make([]byte, len(pairScanner.buf), 2*cap(pairScanner.buf)+minRead)
		copy(newBuf, pairScanner.buf)
		pairScanner.buf = newBuf
	}

	// Read. Delay error for next iteration (after scan).
	n, err := pairScanner.r.Read(pairScanner.buf[len(pairScanner.buf):cap(pairScanner.buf)])
	pairScanner.buf = pairScanner.buf[0 : len(pairScanner.buf)+n]

	return err
}

func findMatchedTokenPair(c byte, pairs []DelimiterPair) (tokenPair DelimiterPair, has bool) {
	opt := slice.NewStream().WithSlice(pairs).FindFirst(func(e interface{}) bool {
		pair := e.(DelimiterPair)
		if pair.start == c || pair.end == c {
			return true
		}
		return false
	})
	if !opt.IsPresent() {
		return tokenPair, false
	}
	return opt.Get().(DelimiterPair), true
}

func nonSpace(b []byte) bool {
	for _, c := range b {
		if !unicode.IsSpace(rune(c)) {
			return true
		}
	}
	return false
}
