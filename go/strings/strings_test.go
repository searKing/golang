// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings_test

import (
	"testing"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/searKing/golang/go/strings"
)

type SliceContainsTest struct {
	inputSS []string
	inputTT []string
	output  bool
}

var (
	sliceContainsAnyTests = []SliceContainsTest{
		{
			[]string{"A", "B", "C", "D"},
			nil,
			true,
		},
		{
			[]string{},
			[]string{"E"},
			false,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"A"},
			true,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"A", "D"},
			true,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"A", "E"},
			true,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"E", "A"},
			true,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"E"},
			false,
		},
	}
)

func TestSliceContainsAny(t *testing.T) {
	for n, test := range sliceContainsAnyTests {
		out := strings.SliceContainsAny(test.inputSS, test.inputTT...)
		if out != test.output {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	sliceContainsTests = []SliceContainsTest{
		{
			[]string{"A", "B", "C", "D"},
			nil,
			true,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"A"},
			true,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"A", "D"},
			true,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"A", "E"},
			false,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"E", "A"},
			false,
		},
		{
			[]string{"A", "B", "C", "D"},
			[]string{"E"},
			false,
		},
	}
)

func TestSliceContains(t *testing.T) {
	for n, test := range sliceContainsTests {
		out := strings.SliceContains(test.inputSS, test.inputTT...)
		if out != test.output {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

func tenRunes(ch rune) string {
	r := make([]rune, 10)
	for i := range r {
		r[i] = ch
	}
	return string(r)
}

func leadingTenRunes(lead, ch rune) string {
	r := make([]rune, 10)
	for i := range r {
		if i == 0 {
			if lead < 0 {
				continue
			}
			r[i] = lead
			continue
		}
		if ch < 0 {
			continue
		}
		r[i] = ch
	}
	return string(r)
}

// User-defined self-inverse mapping function
func rot13(r rune) rune {
	step := rune(13)
	if r >= 'a' && r <= 'z' {
		return ((r - 'a' + step) % 26) + 'a'
	}
	if r >= 'A' && r <= 'Z' {
		return ((r - 'A' + step) % 26) + 'A'
	}
	return r
}

func TestMapLeading(t *testing.T) {
	// Run a couple of awful growth/shrinkage tests
	a := tenRunes('a')
	// 1.  Grow. This triggers two reallocations in Map.
	maxRune := func(rune) rune { return unicode.MaxRune }
	m := strings.MapLeading(maxRune, a)
	expect := leadingTenRunes(unicode.MaxRune, 'a')
	if m != expect {
		t.Errorf("growing: expected %q got %q", expect, m)
	}

	// 2. Shrink
	minRune := func(rune) rune { return 'a' }
	m = strings.MapLeading(minRune, leadingTenRunes(unicode.MaxRune, 'a'))
	expect = a
	if m != expect {
		t.Errorf("shrinking: expected %q got %q", expect, m)
	}

	// 3. Rot13
	m = strings.MapLeading(rot13, "a to zed")
	expect = "n to zed"
	if m != expect {
		t.Errorf("rot13: expected %q got %q", expect, m)
	}

	// 4. Rot13^2
	m = strings.MapLeading(rot13, strings.MapLeading(rot13, "a to zed"))
	expect = "a to zed"
	if m != expect {
		t.Errorf("rot13: expected %q got %q", expect, m)
	}

	// 5. Drop
	dropNotLatin := func(r rune) rune {
		if unicode.Is(unicode.Latin, r) {
			return r
		}
		return -1
	}
	m = strings.MapLeading(dropNotLatin, "세계, Hello")
	expect = "계, Hello"
	if m != expect {
		t.Errorf("drop: expected %q got %q", expect, m)
	}

	// 6. Identity
	identity := func(r rune) rune {
		return r
	}
	orig := "Input string that we expect not to be copied."
	m = strings.MapLeading(identity, orig)
	if unsafe.StringData(orig) != unsafe.StringData(m) {
		t.Error("unexpected copy during identity map")
	}

	// 7. Handle invalid UTF-8 sequence
	replaceNotLatin := func(r rune) rune {
		if unicode.Is(unicode.Latin, r) {
			return r
		}
		return utf8.RuneError
	}
	m = strings.MapLeading(replaceNotLatin, "中 Hello\255World")
	expect = "\uFFFD Hello\255World"
	if m != expect {
		t.Errorf("replace invalid sequence: expected %q got %q", expect, m)
	}

	// 8. Check utf8.RuneSelf and utf8.MaxRune encoding
	encode := func(r rune) rune {
		switch r {
		case utf8.RuneSelf:
			return unicode.MaxRune
		case unicode.MaxRune:
			return utf8.RuneSelf
		}
		return r
	}
	s := string(rune(utf8.RuneSelf)) + string(utf8.MaxRune)
	r := string(utf8.MaxRune) + string(utf8.MaxRune) // reverse of s
	m = strings.MapLeading(encode, s)
	if m != r {
		t.Errorf("encoding not handled correctly: expected %q got %q", r, m)
	}
	m = strings.MapLeading(encode, r)
	if m != s {
		t.Errorf("encoding not handled correctly: expected %q got %q", s, m)
	}

	// 9. Check mapping occurs in the front, middle and back
	trimSpaces := func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}
	m = strings.MapLeading(trimSpaces, "   abc    123   ")
	expect = "  abc    123   "
	if m != expect {
		t.Errorf("trimSpaces: expected %q got %q", expect, m)
	}
}

func TestContainsAsciiVisual(t *testing.T) {
	table := []struct {
		Q string
		R bool
	}{
		{
			Q: string(rune(0x00)),
			R: false,
		},
		{
			Q: " ",
			R: false,
		},
		{
			Q: "!",
			R: true,
		},
		{
			Q: `"`,
			R: true,
		},
		{
			Q: "0",
			R: true,
		},
		{
			Q: ":",
			R: true,
		},
		{
			Q: "A",
			R: true,
		},
		{
			Q: "{",
			R: true,
		},
		{
			Q: "~",
			R: true,
		},
		{
			Q: string(rune(0xFF)),
			R: false,
		},
	}

	for i, test := range table {
		qr := strings.ContainsAsciiVisual(test.Q)
		if qr != test.R {
			t.Errorf("#%d. got %t, want %t", i, qr, test.R)
		}
	}
}

func TestContainsOnlyAsciiVisual(t *testing.T) {
	table := []struct {
		Q string
		R bool
	}{
		//{
		//	Q: "123qwe<>?+_{",
		//	R: true,
		//},
		{
			Q: string(rune(0x00)) + "a",
			R: false,
		},
		{
			Q: string(rune(0xFF)) + "a",
			R: false,
		},
	}

	for i, test := range table {
		qr := strings.ContainsOnlyAsciiVisual(test.Q)
		if qr != test.R {
			t.Errorf("#%d. got %t, want %t", i, qr, test.R)
		}
	}
}

func TestJoinRepeat(t *testing.T) {
	table := []struct {
		Q   string
		sep string
		n   int
		R   string
	}{
		{
			Q:   "a",
			sep: ",",
			n:   -1,
			R:   "",
		},
		{
			Q:   "a",
			sep: ",",
			n:   0,
			R:   "",
		},
		{
			Q:   "a",
			sep: ",",
			n:   1,
			R:   "a",
		},
		{
			Q:   "a",
			sep: ",",
			n:   10,
			R:   "a,a,a,a,a,a,a,a,a,a",
		},
	}

	for i, test := range table {
		qr := strings.JoinRepeat(test.Q, test.sep, test.n)
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestPadLeft(t *testing.T) {
	table := []struct {
		Q   string
		pad string
		n   int
		R   string
	}{
		{
			Q:   "a",
			pad: "*",
			n:   -1,
			R:   "a",
		},
		{
			Q:   "a",
			pad: "*",
			n:   10,
			R:   "*********a",
		},
		{
			Q:   "a",
			pad: "*^",
			n:   5,
			R:   "*^*^a",
		},
		{
			Q:   "a",
			pad: "*^",
			n:   6,
			R:   "*^*^ a",
		},
	}

	for i, test := range table {
		qr := strings.PadLeft(test.Q, test.pad, test.n)
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestPadRight(t *testing.T) {
	table := []struct {
		Q   string
		pad string
		n   int
		R   string
	}{
		{
			Q:   "a",
			pad: "*",
			n:   -1,
			R:   "a",
		},
		{
			Q:   "a",
			pad: "*",
			n:   1,
			R:   "a",
		},
		{
			Q:   "a",
			pad: "*",
			n:   10,
			R:   "a*********",
		},
		{
			Q:   "a",
			pad: "*^",
			n:   5,
			R:   "a*^*^",
		},
		{
			Q:   "a",
			pad: "*^",
			n:   6,
			R:   "a *^*^",
		},
	}

	for i, test := range table {
		qr := strings.PadRight(test.Q, test.pad, test.n)
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestReverseByByte(t *testing.T) {
	table := []struct {
		Q string
		R string
	}{
		{
			Q: "abc123",
			R: "321cba",
		},
		{
			Q: "Hello, 世界",
			R: "\x8c\x95疸\xe4 ,olleH",
		},
	}

	for i, test := range table {
		qr := strings.ReverseByByte(test.Q)
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestReverseByRune(t *testing.T) {
	table := []struct {
		Q string
		R string
	}{
		{
			Q: "abc123",
			R: "321cba",
		},
		{
			Q: "Hello, 世界",
			R: "界世 ,olleH",
		},
	}

	for i, test := range table {
		qr := strings.ReverseByRune(test.Q)
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestCountPrefix(t *testing.T) {
	table := []struct {
		Q string
		p string
		R int
	}{
		{
			Q: "abc123",
			p: "a",
			R: 1,
		},
		{
			Q: "世界, Hello",
			p: "世界",
			R: 1,
		},
		{
			Q: "aaaabc123",
			p: "a",
			R: 4,
		},
		{
			Q: "aaaabc123",
			p: "aa",
			R: 2,
		},
		{
			Q: "aaaabc123",
			p: "aaa",
			R: 1,
		},
		{
			Q: "aaaabc123",
			p: "aaaa",
			R: 1,
		},
		{
			Q: "aaaabc123",
			p: "aaaaa",
			R: 0,
		},
	}

	for i, test := range table {
		qr := strings.CountPrefix(test.Q, test.p)
		if qr != test.R {
			t.Errorf("#%d. got %d, want %d", i, qr, test.R)
		}
	}
}

func TestCountSuffix(t *testing.T) {
	table := []struct {
		Q string
		p string
		R int
	}{
		{
			Q: "abc123",
			p: "3",
			R: 1,
		},
		{
			Q: "Hello, 世界",
			p: "世界",
			R: 1,
		},
		{
			Q: "abc123333",
			p: "3",
			R: 4,
		},
		{
			Q: "abc123333",
			p: "33",
			R: 2,
		},
		{
			Q: "abc123333",
			p: "333",
			R: 1,
		},
		{
			Q: "abc123333",
			p: "3333",
			R: 1,
		},
		{
			Q: "abc123333",
			p: "33333",
			R: 0,
		},
	}

	for i, test := range table {
		qr := strings.CountSuffix(test.Q, test.p)
		if qr != test.R {
			t.Errorf("#%d. got %d, want %d", i, qr, test.R)
		}
	}
}
