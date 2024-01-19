// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson

import (
	"bytes"
	"encoding"
	"fmt"
	"log"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

type Optionals struct {
	Sr string `json:"sr"`
	So string `json:"so,omitempty"`
	Sw string `json:"-"`

	Ir int `json:"omitempty"` // actually named omitempty, not an option
	Io int `json:"io,omitempty"`

	Slr []string `json:"slr,random"`
	Slo []string `json:"slo,omitempty"`

	Mr map[string]any `json:"mr"`
	Mo map[string]any `json:",omitempty"`

	Fr float64 `json:"fr"`
	Fo float64 `json:"fo,omitempty"`

	Br bool `json:"br"`
	Bo bool `json:"bo,omitempty"`

	Ur uint `json:"ur"`
	Uo uint `json:"uo,omitempty"`

	Str struct{} `json:"str"`
	Sto struct{} `json:"sto,omitempty"`
}

var optionalsExpected = `{
 "sr": "",
 "omitempty": 0,
 "slr": null,
 "mr": {},
 "fr": 0,
 "br": false,
 "ur": 0,
 "str": {}
}`

func TestOmitEmpty(t *testing.T) {
	var o Optionals
	o.Sw = "something"
	o.Mr = map[string]any{}
	o.Mo = map[string]any{}

	got, err := MarshalIndent(&o, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	if got := string(got); got != optionalsExpected {
		t.Errorf(" got: %s\nwant: %s\n", got, optionalsExpected)
	}
}

type StringTag struct {
	BoolStr    bool    `json:",string"`
	IntStr     int64   `json:",string"`
	UintptrStr uintptr `json:",string"`
	StrStr     string  `json:",string"`
	NumberStr  Number  `json:",string"`
}

// byte slices are special even if they're renamed types.
type renamedByte byte
type renamedByteSlice []byte
type renamedRenamedByteSlice []renamedByte

func TestEncodeRenamedByteSlice(t *testing.T) {
	s := renamedByteSlice("abc")
	result, err := Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	expect := `"YWJj"`
	if string(result) != expect {
		t.Errorf(" got %s want %s", result, expect)
	}
	r := renamedRenamedByteSlice("abc")
	result, err = Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != expect {
		t.Errorf(" got %s want %s", result, expect)
	}
}

type SamePointerNoCycle struct {
	Ptr1, Ptr2 *SamePointerNoCycle
}

var samePointerNoCycle = &SamePointerNoCycle{}

type PointerCycle struct {
	Ptr *PointerCycle
}

var pointerCycle = &PointerCycle{}

type PointerCycleIndirect struct {
	Ptrs []any
}

type RecursiveSlice []RecursiveSlice

var (
	pointerCycleIndirect = &PointerCycleIndirect{}
	mapCycle             = make(map[string]any)
	sliceCycle           = []any{nil}
	sliceNoCycle         = []any{nil, nil}
	recursiveSliceCycle  = []RecursiveSlice{nil}
)

func init() {
	ptr := &SamePointerNoCycle{}
	samePointerNoCycle.Ptr1 = ptr
	samePointerNoCycle.Ptr2 = ptr

	pointerCycle.Ptr = pointerCycle
	pointerCycleIndirect.Ptrs = []any{pointerCycleIndirect}

	mapCycle["x"] = mapCycle
	sliceCycle[0] = sliceCycle
	sliceNoCycle[1] = sliceNoCycle[:1]
	for i := startDetectingCyclesAfter; i > 0; i-- {
		sliceNoCycle = []any{sliceNoCycle}
	}
	recursiveSliceCycle[0] = recursiveSliceCycle
}

func TestSamePointerNoCycle(t *testing.T) {
	if _, err := Marshal(samePointerNoCycle); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSliceNoCycle(t *testing.T) {
	if _, err := Marshal(sliceNoCycle); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

var unsupportedValues = []any{
	math.NaN(),
	math.Inf(-1),
	math.Inf(1),
	pointerCycle,
	pointerCycleIndirect,
	mapCycle,
	sliceCycle,
	recursiveSliceCycle,
}

func TestUnsupportedValues(t *testing.T) {
	for _, v := range unsupportedValues {
		if _, err := Marshal(v); err != nil {
			if _, ok := err.(*UnsupportedValueError); !ok {
				t.Errorf("for %v, got %T want UnsupportedValueError", v, err)
			}
		} else {
			t.Errorf("for %v, expected error", v)
		}
	}
}

// Issue 43207
func TestMarshalTextFloatMap(t *testing.T) {
	m := map[textfloat]string{
		textfloat(math.NaN()): "1",
		textfloat(math.NaN()): "1",
	}
	got, err := Marshal(m)
	if err != nil {
		t.Errorf("Marshal() error: %v", err)
	}
	want := `{"TF:NaN":"1","TF:NaN":"1"}`
	if string(got) != want {
		t.Errorf("Marshal() = %s, want %s", got, want)
	}
}

// Ref has Marshaler and Unmarshaler methods with pointer receiver.
type Ref int

func (*Ref) MarshalJSON() ([]byte, error) {
	return []byte(`"ref"`), nil
}

func (r *Ref) UnmarshalJSON([]byte) error {
	*r = 12
	return nil
}

// Val has Marshaler methods with value receiver.
type Val int

func (Val) MarshalJSON() ([]byte, error) {
	return []byte(`"val"`), nil
}

// RefText has Marshaler and Unmarshaler methods with pointer receiver.
type RefText int

func (*RefText) MarshalText() ([]byte, error) {
	return []byte(`"ref"`), nil
}

func (r *RefText) UnmarshalText([]byte) error {
	*r = 13
	return nil
}

// ValText has Marshaler methods with value receiver.
type ValText int

func (ValText) MarshalText() ([]byte, error) {
	return []byte(`"val"`), nil
}

func TestRefValMarshal(t *testing.T) {
	var s = struct {
		R0 Ref
		R1 *Ref
		R2 RefText
		R3 *RefText
		V0 Val
		V1 *Val
		V2 ValText
		V3 *ValText
	}{
		R0: 12,
		R1: new(Ref),
		R2: 14,
		R3: new(RefText),
		V0: 13,
		V1: new(Val),
		V2: 15,
		V3: new(ValText),
	}
	const want = `{"R0":"ref","R1":"ref","R2":"\"ref\"","R3":"\"ref\"","V0":"val","V1":"val","V2":"\"val\"","V3":"\"val\""}`
	b, err := Marshal(&s)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// C implements Marshaler and returns unescaped JSON.
type C int

func (C) MarshalJSON() ([]byte, error) {
	return []byte(`"<&>"`), nil
}

// CText implements Marshaler and returns unescaped text.
type CText int

func (CText) MarshalText() ([]byte, error) {
	return []byte(`"<&>"`), nil
}

func TestMarshalerEscaping(t *testing.T) {
	var c C
	want := `"\u003c\u0026\u003e"`
	b, err := Marshal(c)
	if err != nil {
		t.Fatalf("Marshal(c): %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("Marshal(c) = %#q, want %#q", got, want)
	}

	var ct CText
	want = `"\"\u003c\u0026\u003e\""`
	b, err = Marshal(ct)
	if err != nil {
		t.Fatalf("Marshal(ct): %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("Marshal(ct) = %#q, want %#q", got, want)
	}
}

func TestAnonymousFields(t *testing.T) {
	tests := []struct {
		label     string     // Test name
		makeInput func() any // Function to create input value
		want      string     // Expected JSON output
	}{{
		// Both S1 and S2 have a field named X. From the perspective of S,
		// it is ambiguous which one X refers to.
		// This should not serialize either field.
		label: "AmbiguousField",
		makeInput: func() any {
			type (
				S1 struct{ x, X int }
				S2 struct{ x, X int }
				S  struct {
					S1
					S2
				}
			)
			return S{S1{1, 2}, S2{3, 4}}
		},
		want: `{}`,
	}, {
		label: "DominantField",
		// Both S1 and S2 have a field named X, but since S has an X field as
		// well, it takes precedence over S1.X and S2.X.
		makeInput: func() any {
			type (
				S1 struct{ x, X int }
				S2 struct{ x, X int }
				S  struct {
					S1
					S2
					x, X int
				}
			)
			return S{S1{1, 2}, S2{3, 4}, 5, 6}
		},
		want: `{"X":6}`,
	}, {
		// Unexported embedded field of non-struct type should not be serialized.
		label: "UnexportedEmbeddedInt",
		makeInput: func() any {
			type (
				myInt int
				S     struct{ myInt }
			)
			return S{5}
		},
		want: `{}`,
	}, {
		// Exported embedded field of non-struct type should be serialized.
		label: "ExportedEmbeddedInt",
		makeInput: func() any {
			type (
				MyInt int
				S     struct{ MyInt }
			)
			return S{5}
		},
		want: `{"MyInt":5}`,
	}, {
		// Unexported embedded field of pointer to non-struct type
		// should not be serialized.
		label: "UnexportedEmbeddedIntPointer",
		makeInput: func() any {
			type (
				myInt int
				S     struct{ *myInt }
			)
			s := S{new(myInt)}
			*s.myInt = 5
			return s
		},
		want: `{}`,
	}, {
		// Exported embedded field of pointer to non-struct type
		// should be serialized.
		label: "ExportedEmbeddedIntPointer",
		makeInput: func() any {
			type (
				MyInt int
				S     struct{ *MyInt }
			)
			s := S{new(MyInt)}
			*s.MyInt = 5
			return s
		},
		want: `{"MyInt":5}`,
	}, {
		// Exported fields of embedded structs should have their
		// exported fields be serialized regardless of whether the struct types
		// themselves are exported.
		label: "EmbeddedStruct",
		makeInput: func() any {
			type (
				s1 struct{ x, X int }
				S2 struct{ y, Y int }
				S  struct {
					s1
					S2
				}
			)
			return S{s1{1, 2}, S2{3, 4}}
		},
		want: `{"X":2,"Y":4}`,
	}, {
		// Exported fields of pointers to embedded structs should have their
		// exported fields be serialized regardless of whether the struct types
		// themselves are exported.
		label: "EmbeddedStructPointer",
		makeInput: func() any {
			type (
				s1 struct{ x, X int }
				S2 struct{ y, Y int }
				S  struct {
					*s1
					*S2
				}
			)
			return S{&s1{1, 2}, &S2{3, 4}}
		},
		want: `{"X":2,"Y":4}`,
	}, {
		// Exported fields on embedded unexported structs at multiple levels
		// of nesting should still be serialized.
		label: "NestedStructAndInts",
		makeInput: func() any {
			type (
				MyInt1 int
				MyInt2 int
				myInt  int
				s2     struct {
					MyInt2
					myInt
				}
				s1 struct {
					MyInt1
					myInt
					s2
				}
				S struct {
					s1
					myInt
				}
			)
			return S{s1{1, 2, s2{3, 4}}, 6}
		},
		want: `{"MyInt1":1,"MyInt2":3}`,
	}, {
		// If an anonymous struct pointer field is nil, we should ignore
		// the embedded fields behind it. Not properly doing so may
		// result in the wrong output or reflect panics.
		label: "EmbeddedFieldBehindNilPointer",
		makeInput: func() any {
			type (
				S2 struct{ Field string }
				S  struct{ *S2 }
			)
			return S{}
		},
		want: `{}`,
	}}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			b, err := Marshal(tt.makeInput())
			if err != nil {
				t.Fatalf("Marshal() = %v, want nil error", err)
			}
			if string(b) != tt.want {
				t.Fatalf("Marshal() = %q, want %q", b, tt.want)
			}
		})
	}
}

type BugA struct {
	S string
}

type BugB struct {
	BugA
	S string
}

type BugC struct {
	S string
}

// Legal Go: We never use the repeated embedded field (S).
type BugX struct {
	A int
	BugA
	BugB
}

// golang.org/issue/16042.
// Even if a nil interface value is passed in, as long as
// it implements Marshaler, it should be marshaled.
type nilJSONMarshaler string

func (nm *nilJSONMarshaler) MarshalJSON() ([]byte, error) {
	if nm == nil {
		return Marshal("0zenil0")
	}
	return Marshal("zenil:" + string(*nm))
}

// golang.org/issue/34235.
// Even if a nil interface value is passed in, as long as
// it implements encoding.TextMarshaler, it should be marshaled.
type nilTextMarshaler string

func (nm *nilTextMarshaler) MarshalText() ([]byte, error) {
	if nm == nil {
		return []byte("0zenil0"), nil
	}
	return []byte("zenil:" + string(*nm)), nil
}

// See golang.org/issue/16042 and golang.org/issue/34235.
func TestNilMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: nil, want: `null`},
		{v: new(float64), want: `0`},
		{v: []any(nil), want: `null`},
		{v: []string(nil), want: `null`},
		{v: map[string]string(nil), want: `null`},
		{v: []byte(nil), want: `null`},
		{v: struct{ M string }{"gopher"}, want: `{"M":"gopher"}`},
		{v: struct{ M Marshaler }{}, want: `{"M":null}`},
		{v: struct{ M Marshaler }{(*nilJSONMarshaler)(nil)}, want: `{"M":"0zenil0"}`},
		{v: struct{ M any }{(*nilJSONMarshaler)(nil)}, want: `{"M":null}`},
		{v: struct{ M encoding.TextMarshaler }{}, want: `{"M":null}`},
		{v: struct{ M encoding.TextMarshaler }{(*nilTextMarshaler)(nil)}, want: `{"M":"0zenil0"}`},
		{v: struct{ M any }{(*nilTextMarshaler)(nil)}, want: `{"M":null}`},
	}

	for _, tt := range testCases {
		out, err := Marshal(tt.v)
		if err != nil || string(out) != tt.want {
			t.Errorf("Marshal(%#v) = %#q, %#v, want %#q, nil", tt.v, out, err, tt.want)
			continue
		}
	}
}

// Issue 5245.
func TestEmbeddedBug(t *testing.T) {
	v := BugB{
		BugA{"A"},
		"B",
	}
	b, err := Marshal(v)
	if err != nil {
		t.Fatal("Marshal:", err)
	}
	want := `{"S":"B"}`
	got := string(b)
	if got != want {
		t.Fatalf("Marshal: got %s want %s", got, want)
	}
	// Now check that the duplicate field, S, does not appear.
	x := BugX{
		A: 23,
	}
	b, err = Marshal(x)
	if err != nil {
		t.Fatal("Marshal:", err)
	}
	want = `{"A":23}`
	got = string(b)
	if got != want {
		t.Fatalf("Marshal: got %s want %s", got, want)
	}
}

type BugD struct { // Same as BugA after tagging.
	XXX string `json:"S"`
}

// BugD's tagged S field should dominate BugA's.
type BugY struct {
	BugA
	BugD
}

// Test that a field with a tag dominates untagged fields.
func TestTaggedFieldDominates(t *testing.T) {
	v := BugY{
		BugA{"BugA"},
		BugD{"BugD"},
	}
	b, err := Marshal(v)
	if err != nil {
		t.Fatal("Marshal:", err)
	}
	want := `{"S":"BugD"}`
	got := string(b)
	if got != want {
		t.Fatalf("Marshal: got %s want %s", got, want)
	}
}

// There are no tags here, so S should not appear.
type BugZ struct {
	BugA
	BugC
	BugY // Contains a tagged S field through BugD; should not dominate.
}

func TestDuplicatedFieldDisappears(t *testing.T) {
	v := BugZ{
		BugA{"BugA"},
		BugC{"BugC"},
		BugY{
			BugA{"nested BugA"},
			BugD{"nested BugD"},
		},
	}
	b, err := Marshal(v)
	if err != nil {
		t.Fatal("Marshal:", err)
	}
	want := `{}`
	got := string(b)
	if got != want {
		t.Fatalf("Marshal: got %s want %s", got, want)
	}
}

func TestIssue10281(t *testing.T) {
	type Foo struct {
		N Number
	}
	x := Foo{Number(`invalid`)}

	b, err := Marshal(&x)
	if err == nil {
		t.Errorf("Marshal(&x) = %#q; want error", b)
	}
}

var encodeStringTests = []struct {
	in  string
	out string
}{
	{"\x00", `"\u0000"`},
	{"\x01", `"\u0001"`},
	{"\x02", `"\u0002"`},
	{"\x03", `"\u0003"`},
	{"\x04", `"\u0004"`},
	{"\x05", `"\u0005"`},
	{"\x06", `"\u0006"`},
	{"\x07", `"\u0007"`},
	{"\x08", `"\u0008"`},
	{"\x09", `"\t"`},
	{"\x0a", `"\n"`},
	{"\x0b", `"\u000b"`},
	{"\x0c", `"\u000c"`},
	{"\x0d", `"\r"`},
	{"\x0e", `"\u000e"`},
	{"\x0f", `"\u000f"`},
	{"\x10", `"\u0010"`},
	{"\x11", `"\u0011"`},
	{"\x12", `"\u0012"`},
	{"\x13", `"\u0013"`},
	{"\x14", `"\u0014"`},
	{"\x15", `"\u0015"`},
	{"\x16", `"\u0016"`},
	{"\x17", `"\u0017"`},
	{"\x18", `"\u0018"`},
	{"\x19", `"\u0019"`},
	{"\x1a", `"\u001a"`},
	{"\x1b", `"\u001b"`},
	{"\x1c", `"\u001c"`},
	{"\x1d", `"\u001d"`},
	{"\x1e", `"\u001e"`},
	{"\x1f", `"\u001f"`},
}

func TestEncodeString(t *testing.T) {
	for _, tt := range encodeStringTests {
		b, err := Marshal(tt.in)
		if err != nil {
			t.Errorf("Marshal(%q): %v", tt.in, err)
			continue
		}
		out := string(b)
		if out != tt.out {
			t.Errorf("Marshal(%q) = %#q, want %#q", tt.in, out, tt.out)
		}
	}
}

type jsonbyte byte

func (b jsonbyte) MarshalJSON() ([]byte, error) { return tenc(`{"JB":%d}`, b) }

type textbyte byte

func (b textbyte) MarshalText() ([]byte, error) { return tenc(`TB:%d`, b) }

type jsonint int

func (i jsonint) MarshalJSON() ([]byte, error) { return tenc(`{"JI":%d}`, i) }

type textint int

func (i textint) MarshalText() ([]byte, error) { return tenc(`TI:%d`, i) }

func tenc(format string, a ...any) ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, a...)
	return buf.Bytes(), nil
}

type textfloat float64

func (f textfloat) MarshalText() ([]byte, error) { return tenc(`TF:%0.2f`, f) }

// Issue 13783
func TestEncodeBytekind(t *testing.T) {
	testdata := []struct {
		data any
		want string
	}{
		{byte(7), "7"},
		{jsonbyte(7), `{"JB":7}`},
		{textbyte(4), `"TB:4"`},
		{jsonint(5), `{"JI":5}`},
		{textint(1), `"TI:1"`},
		{[]byte{0, 1}, `"AAE="`},
		{[]jsonbyte{0, 1}, `[{"JB":0},{"JB":1}]`},
		{[][]jsonbyte{{0, 1}, {3}}, `[[{"JB":0},{"JB":1}],[{"JB":3}]]`},
		{[]textbyte{2, 3}, `["TB:2","TB:3"]`},
		{[]jsonint{5, 4}, `[{"JI":5},{"JI":4}]`},
		{[]textint{9, 3}, `["TI:9","TI:3"]`},
		{[]int{9, 3}, `[9,3]`},
		{[]textfloat{12, 3}, `["TF:12.00","TF:3.00"]`},
	}
	for _, d := range testdata {
		js, err := Marshal(d.data)
		if err != nil {
			t.Error(err)
			continue
		}
		got, want := string(js), d.want
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}
}

var re = regexp.MustCompile

// syntactic checks on form of marshaled floating point numbers.
var badFloatREs = []*regexp.Regexp{
	re(`p`),                     // no binary exponential notation
	re(`^\+`),                   // no leading + sign
	re(`^-?0[^.]`),              // no unnecessary leading zeros
	re(`^-?\.`),                 // leading zero required before decimal point
	re(`\.(e|$)`),               // no trailing decimal
	re(`\.[0-9]+0(e|$)`),        // no trailing zero in fraction
	re(`^-?(0|[0-9]{2,})\..*e`), // exponential notation must have normalized mantissa
	re(`e[0-9]`),                // positive exponent must be signed
	re(`e[+-]0`),                // exponent must not have leading zeros
	re(`e-[1-6]$`),              // not tiny enough for exponential notation
	re(`e+(.|1.|20)$`),          // not big enough for exponential notation
	re(`^-?0\.0000000`),         // too tiny, should use exponential notation
	re(`^-?[0-9]{22}`),          // too big, should use exponential notation
	re(`[1-9][0-9]{16}[1-9]`),   // too many significant digits in integer
	re(`[1-9][0-9.]{17}[1-9]`),  // too many significant digits in decimal
	// below here for float32 only
	re(`[1-9][0-9]{8}[1-9]`),  // too many significant digits in integer
	re(`[1-9][0-9.]{9}[1-9]`), // too many significant digits in decimal
}

func TestMarshalFloat(t *testing.T) {
	t.Parallel()
	nfail := 0
	test := func(f float64, bits int) {
		vf := any(f)
		if bits == 32 {
			f = float64(float32(f)) // round
			vf = float32(f)
		}
		bout, err := Marshal(vf)
		if err != nil {
			t.Errorf("Marshal(%T(%g)): %v", vf, vf, err)
			nfail++
			return
		}
		out := string(bout)

		// result must convert back to the same float
		g, err := strconv.ParseFloat(out, bits)
		if err != nil {
			t.Errorf("Marshal(%T(%g)) = %q, cannot parse back: %v", vf, vf, out, err)
			nfail++
			return
		}
		if f != g || fmt.Sprint(f) != fmt.Sprint(g) { // fmt.Sprint handles ±0
			t.Errorf("Marshal(%T(%g)) = %q (is %g, not %g)", vf, vf, out, float32(g), vf)
			nfail++
			return
		}

		bad := badFloatREs
		if bits == 64 {
			bad = bad[:len(bad)-2]
		}
		for _, re := range bad {
			if re.MatchString(out) {
				t.Errorf("Marshal(%T(%g)) = %q, must not match /%s/", vf, vf, out, re)
				nfail++
				return
			}
		}
	}

	var (
		bigger  = math.Inf(+1)
		smaller = math.Inf(-1)
	)

	var digits = "1.2345678901234567890123"
	for i := len(digits); i >= 2; i-- {
		if testing.Short() && i < len(digits)-4 {
			break
		}
		for exp := -30; exp <= 30; exp++ {
			for _, sign := range "+-" {
				for bits := 32; bits <= 64; bits += 32 {
					s := fmt.Sprintf("%c%se%d", sign, digits[:i], exp)
					f, err := strconv.ParseFloat(s, bits)
					if err != nil {
						log.Fatal(err)
					}
					next := math.Nextafter
					if bits == 32 {
						next = func(g, h float64) float64 {
							return float64(math.Nextafter32(float32(g), float32(h)))
						}
					}
					test(f, bits)
					test(next(f, bigger), bits)
					test(next(f, smaller), bits)
					if nfail > 50 {
						t.Fatalf("stopping test early")
					}
				}
			}
		}
	}
	test(0, 64)
	test(math.Copysign(0, -1), 64)
	test(0, 32)
	test(math.Copysign(0, -1), 32)
}

type marshalPanic struct{}

func (marshalPanic) MarshalJSON() ([]byte, error) { panic(0xdead) }

func TestMarshalPanic(t *testing.T) {
	defer func() {
		if got := recover(); !reflect.DeepEqual(got, 0xdead) {
			t.Errorf("panic() = (%T)(%v), want 0xdead", got, got)
		}
	}()
	Marshal(&marshalPanic{})
	t.Error("Marshal should have panicked")
}

func TestMarshalUncommonFieldNames(t *testing.T) {
	v := struct {
		A0, À, Aβ int
	}{}
	b, err := Marshal(v)
	if err != nil {
		t.Fatal("Marshal:", err)
	}
	want := `{"A0":0,"À":0,"Aβ":0}`
	got := string(b)
	if got != want {
		t.Fatalf("Marshal: got %s want %s", got, want)
	}
}

type marshalTruncate struct {
	Sr string `json:"sr"`

	Blr []byte `json:"blr,random"`

	Slr []string `json:"slr,random"`

	Mr map[string]*marshalTruncate `json:"mr"`

	Str *marshalTruncate `json:"str"`
}

var truncateExpected = `{
 "sr": "ab...26 chars",
 "blr": "YWI=...52 bytes",
 "slr": [
  "0 ...28 chars",
  "1 ...54 chars",
  "...7 elems"
 ],
 "mr": {
  "a": {
   "sr": "ab...26 chars",
   "blr": "YWI=...26 bytes",
   "slr": [
    "0 ...28 chars",
    "1 ...54 chars",
    "...7 elems"
   ],
   "mr": null,
   "str": null
  },
  "b": null,
  "...4 pairs": "4"
 },
 "str": {
  "sr": "ab...26 chars",
  "blr": "YWI=...26 bytes",
  "slr": [
   "0 ...28 chars",
   "1 ...54 chars",
   "...7 elems"
  ],
  "mr": null,
  "str": null
 }
}`

func TestTruncateMarshal(t *testing.T) {
	var o = marshalTruncate{
		Sr:  "abcdefghihklmnopqrstuvwxyz",
		Blr: []byte("abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		Slr: []string{"0 abcdefghihklmnopqrstuvwxyz",
			"1 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"2 abcdefghihklmnopqrstuvwxyzA",
			"3 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"4 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"5 abcdefghihklmnopqrstuvwxyzA",
			"6 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		Mr: map[string]*marshalTruncate{
			"a": {
				Sr:  "abcdefghihklmnopqrstuvwxyz",
				Blr: []byte("abcdefghihklmnopqrstuvwxyz"),
				Slr: []string{"0 abcdefghihklmnopqrstuvwxyz",
					"1 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
					"2 abcdefghihklmnopqrstuvwxyzA",
					"3 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
					"4 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
					"5 abcdefghihklmnopqrstuvwxyzA",
					"6 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"},
				Mr:  nil,
				Str: nil,
			},
			"b": nil,
			"c": nil,
			"d": nil,
		},
		Str: &marshalTruncate{
			Sr:  "abcdefghihklmnopqrstuvwxyz",
			Blr: []byte("abcdefghihklmnopqrstuvwxyz"),
			Slr: []string{"0 abcdefghihklmnopqrstuvwxyz",
				"1 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
				"2 abcdefghihklmnopqrstuvwxyzA",
				"3 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
				"4 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
				"5 abcdefghihklmnopqrstuvwxyzA",
				"6 abcdefghihklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"},
			Mr:  nil,
			Str: nil,
		},
	}

	got, err := MarshalIndent(&o, "", " ", WithEncOptsTruncate(2))
	if err != nil {
		t.Fatal(err)
	}
	if got := string(got); got != truncateExpected {
		t.Errorf(" got: %s\nwant: %s\n", got, optionalsExpected)
	}
}
