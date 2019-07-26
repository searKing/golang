package scanner_test

import (
	"bufio"
	"bytes"
	"github.com/searKing/golang/go/go/scanner"
	"go/token"
	"testing"
)

type ScanTests struct {
	split  bufio.SplitFunc
	input  []byte // input is the input that we want to tokenize.
	output [][]byte
}

var (
	scanTests = []ScanTests{
		{
			split:  bufio.ScanWords,
			input:  []byte("cos(x) + 1i*sin(x) // Euler"),
			output: [][]byte{[]byte("cos(x)"), []byte("+"), []byte("1i*sin(x)"), []byte("//"), []byte("Euler")},
		},
		{
			split: bufio.ScanBytes,
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
			split:  bufio.ScanLines,
			input:  []byte("cos(x) + 1i*sin(x) \n// Euler"),
			output: [][]byte{[]byte("cos(x) + 1i*sin(x) "), []byte("// Euler")},
		},
		{
			split: bufio.ScanRunes,
			input: []byte("cos(x) + 1i*sin(x) // Euler"),
			output: [][]byte{[]byte("c"), []byte("o"), []byte("s"), []byte("("), []byte("x"), []byte(")"),
				[]byte(" "), []byte("+"), []byte(" "),
				[]byte("1"), []byte("i"), []byte("*"), []byte("s"), []byte("i"), []byte("n"), []byte("("), []byte("x"), []byte(")"),
				[]byte(" "),
				[]byte("/"), []byte("/"),
				[]byte(" "),
				[]byte("E"), []byte("u"), []byte("l"), []byte("e"), []byte("r")},
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
		s.Split(test.split)
		var i int
		for ; ; i++ {
			token, has := s.Scan()
			if !has {
				break
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
