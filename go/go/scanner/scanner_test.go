// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner_test

import (
	"bufio"
	"bytes"
	"go/token"
	"testing"
	"unicode"

	"github.com/searKing/golang/go/go/scanner"
)

type ScanTests struct {
	split  bufio.SplitFunc
	input  []byte // input is the input that we want to tokenize.
	output [][]byte
}

var (
	scanTests = []ScanTests{
		{
			split:  scanner.ScanWords,
			input:  []byte("cos(x) + 1i*sin(x) // Euler"),
			output: [][]byte{[]byte("cos(x)"), []byte("+"), []byte("1i*sin(x)"), []byte("//"), []byte("Euler")},
		},
		{
			split: scanner.ScanBytes,
			input: []byte("cos(x) + 1i*sin(x) // Euler"),
			output: [][]byte{[]byte("c"), []byte("o"), []byte("s"), []byte("("), []byte("x"), []byte(")"),
				[]byte(" "), []byte("+"), []byte(" "),
				[]byte("1"), []byte("i"), []byte("*"), []byte("s"), []byte("i"), []byte("n"), []byte("("), []byte("x"), []byte(")"),
				[]byte(" "),
				[]byte("/"), []byte("/"),
				[]byte(" "),
				[]byte("E"), []byte("u"), []byte("l"), []byte("e"), []byte("r")},
		},
		{
			split:  scanner.ScanLines,
			input:  []byte("cos(x) + 1i*sin(x) \n// Euler"),
			output: [][]byte{[]byte("cos(x) + 1i*sin(x) "), []byte("// Euler")},
		},
		{
			split: scanner.ScanRunes,
			input: []byte("cos(x) + 1i*sin(x) // Euler"),
			output: [][]byte{[]byte("c"), []byte("o"), []byte("s"), []byte("("), []byte("x"), []byte(")"),
				[]byte(" "), []byte("+"), []byte(" "),
				[]byte("1"), []byte("i"), []byte("*"), []byte("s"), []byte("i"), []byte("n"), []byte("("), []byte("x"), []byte(")"),
				[]byte(" "),
				[]byte("/"), []byte("/"),
				[]byte(" "),
				[]byte("E"), []byte("u"), []byte("l"), []byte("e"), []byte("r")},
		},
		{
			split:  scanner.ScanEscapes('"'),
			input:  []byte(`\\\r\n\"cos(x) + 1i*sin(x) // Euler`),
			output: [][]byte{[]byte(`\\`), []byte(`\r`), []byte(`\n`), []byte(`\"`)},
		},
		{
			split:  scanner.ScanEscapes('"'),
			input:  []byte(`\y`),
			output: [][]byte{},
		},
		{
			split:  scanner.ScanEscapes('"'),
			input:  []byte(`\xGGGG`),
			output: [][]byte{},
		},
		{
			split:  scanner.ScanInterpretedStrings,
			input:  []byte(`\xGGGG`),
			output: [][]byte{},
		},
		{
			split:  scanner.ScanInterpretedStrings,
			input:  []byte(`"""Hello World\""\xGGGG`),
			output: [][]byte{[]byte(`""`), []byte(`"Hello World\""`)},
		},
		{
			split:  scanner.ScanRawStrings,
			input:  []byte("```Hello World\r`"),
			output: [][]byte{[]byte("``"), []byte("`Hello World\r`")},
		},
		{
			split:  scanner.ScanMantissas(10),
			input:  []byte("123"),
			output: [][]byte{[]byte("123")},
		},
		{
			split:  scanner.ScanMantissas(8),
			input:  []byte("1238"),
			output: [][]byte{[]byte("123")},
		},
		{
			split:  scanner.ScanNumbers,
			input:  []byte("123"),
			output: [][]byte{[]byte("123")},
		},
		{
			split:  scanner.ScanNumbers,
			input:  []byte("078"),
			output: [][]byte{},
		},
		{
			split:  scanner.ScanNumbers,
			input:  []byte("123.56"),
			output: [][]byte{[]byte("123.56")},
		},
		{
			split:  scanner.ScanNumbers,
			input:  []byte("123.56i78.96i"),
			output: [][]byte{[]byte("123.56i"), []byte("78.96i")},
		},
		{
			split:  scanner.ScanNumbers,
			input:  []byte("0x55.FF"),
			output: [][]byte{[]byte("0x55")},
		},
		{
			split:  scanner.ScanNumbers,
			input:  []byte("077x"),
			output: [][]byte{[]byte("077")},
		},
		{
			split:  scanner.ScanNumbers,
			input:  []byte("0xFG"),
			output: [][]byte{[]byte("0xF")},
		},
		{
			split:  scanner.ScanIdentifier,
			input:  []byte("_HelloWorld_&"),
			output: [][]byte{[]byte("_HelloWorld_")},
		},
		{
			split:  scanner.ScanIdentifier,
			input:  []byte("_Hello World_&"),
			output: [][]byte{[]byte("_Hello")},
		},
		{
			split:  scanner.ScanIdentifier,
			input:  []byte("_你好_&"),
			output: [][]byte{[]byte("_你好_")},
		},
		{
			split:  scanner.ScanIdentifier,
			input:  []byte("0你好_&"),
			output: [][]byte{},
		},
		{
			split:  scanner.ScanIdentifier,
			input:  []byte("H"),
			output: [][]byte{[]byte("H")},
		},
		{
			split:  scanner.ScanUntil(unicode.IsSpace),
			input:  []byte("Hello World"),
			output: [][]byte{[]byte("Hello")},
		},
		{
			split:  scanner.ScanWhile(unicode.IsSpace),
			input:  []byte("  \t Hello World"),
			output: [][]byte{[]byte("  \t ")},
		},
		{
			split:  scanner.ScanRegexpPosix(`^Hello.*[ W]`),
			input:  []byte("Hello World"),
			output: [][]byte{[]byte("Hello W")},
		},
		{
			split:  scanner.ScanRegexpPerl(`^Hello.*[ W]`),
			input:  []byte("Hello World"),
			output: [][]byte{[]byte("Hello W")},
		},
	}
)

func TestScan(t *testing.T) {
Outer:
	for n, test := range scanTests {
		// Initialize the scanner.
		var s scanner.Scanner
		fset := token.NewFileSet()                             // positions are relative to fset
		file := fset.AddFile("", fset.Base(), len(test.input)) // register input "file"
		s.Init(file, test.input, nil /* no error handler */, scanner.ModeCaseSensitive)
		var i int
		for ; ; i++ {
			token, has := s.ScanSplits(test.split)
			if !has {
				break
			}
			if i >= len(test.output) {
				t.Errorf("#%d: scan[%d] got %v; nothing expected", n, i, string(token))
				continue Outer
			}
			if bytes.Compare(token, test.output[i]) != 0 {
				t.Errorf("#%d: scan[%d] got %v; expected %v", n, i, string(token), string(test.output[i]))
				continue Outer
			}
		}
		if i != len(test.output) {
			t.Errorf("#%d: %v: got %d runs; expected %d", n, string(test.input), i, len(test.output))
			continue
		}
	}
}
