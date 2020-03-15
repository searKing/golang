// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiple_prefix_test

import (
	"fmt"
	"math/big"

	"github.com/searKing/golang/go/format/multiple_prefix"
)

func ExampleDecimalFormatFloat() {
	fmt.Printf("%s\n", multiple_prefix.DecimalFormatFloat(1234.56789, 4))

	// Output:
	// 1.2346k
}

func ExampleDecimalFormatInt() {
	fmt.Printf("%s\n", multiple_prefix.DecimalFormatInt(-1234, 4))

	// Output:
	// -1.234k
}

func ExampleDecimalFormatUint() {
	fmt.Printf("%s\n", multiple_prefix.DecimalFormatUint(1234, 4))

	// Output:
	// 1.234k
}

func ExampleDecimalFormatInt64() {
	fmt.Printf("%s\n", multiple_prefix.DecimalFormatInt64(-123456789, 4))

	// Output:
	// -123.4568M
}

func ExampleDecimalFormatUint64() {
	fmt.Printf("%s\n", multiple_prefix.DecimalFormatUint64(123456789, 4))

	// Output:
	// 123.4568M
}

func ExampleBinaryFormatFloat() {
	fmt.Printf("%s\n", multiple_prefix.BinaryFormatFloat(1024.1024, 4))

	// Output:
	// 1.0001Ki
}

func ExampleBinaryFormatInt() {
	fmt.Printf("%s\n", multiple_prefix.BinaryFormatInt(-1024*1024, 4))

	// Output:
	// -1Mi
}

func ExampleBinaryFormatUint() {
	fmt.Printf("%s\n", multiple_prefix.BinaryFormatUint(1024*10000, 4))

	// Output:
	// 9.7656Mi
}

func ExampleBinaryFormatInt64() {
	fmt.Printf("%s\n", multiple_prefix.BinaryFormatInt64(-1024*1024, 4))

	// Output:
	// -1Mi
}

func ExampleBinaryFormatUint64() {
	fmt.Printf("%s\n", multiple_prefix.BinaryFormatUint64(1024*1024, 4))

	// Output:
	// 1Mi
}

func ExampleSplitDecimal() {
	s := "+1234.567890\tkBHello\tWorld"

	gotNumber, gotPrefix, gotUnparsed := multiple_prefix.SplitDecimal(s)
	fmt.Printf("%s\n", s)
	fmt.Printf("Number:%s\n", gotNumber)
	fmt.Printf("Symbol:%s\n", gotPrefix.Symbol())
	fmt.Printf("Unparsed:%s\n", gotUnparsed)

	// Output:
	// +1234.567890	kBHello	World
	// Number:+1234.567890
	// Symbol:k
	// Unparsed:BHello	World
}

func ExampleSplitBinary() {
	s := "+1234.567890 KiBHelloWorld"

	gotNumber, gotPrefix, gotUnparsed := multiple_prefix.SplitBinary(s)
	fmt.Printf("%s\n", s)
	fmt.Printf("Number:%s\n", gotNumber)
	fmt.Printf("Symbol:%s\n", gotPrefix.Symbol())
	fmt.Printf("Unparsed:%s\n", gotUnparsed)

	// Output:
	// +1234.567890 KiBHelloWorld
	// Number:+1234.567890
	// Symbol:Ki
	// Unparsed:BHelloWorld
}

func ExampleDecimalMultiplePrefix_SetFloat64() {
	n := 1234.56789
	mp := multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetFloat64(n)

	fmt.Printf("%g\n", n)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// 1234.56789
	// Base:10
	// Power:3
	// Symbol:k
	// Name:kilo
}

func ExampleDecimalMultiplePrefix_SetInt64() {
	n := int64(-2 * 1000 * 1000)
	mp := multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetInt64(n)

	fmt.Printf("%d\n", n)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// -2000000
	// Base:10
	// Power:6
	// Symbol:M
	// Name:mega
}

func ExampleDecimalMultiplePrefix_SetUint64() {
	n := uint64(2 * 1000 * 1000)
	mp := multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetUint64(n)

	fmt.Printf("%d\n", n)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// 2000000
	// Base:10
	// Power:6
	// Symbol:M
	// Name:mega
}

func ExampleDecimalMultiplePrefix_SetBigFloat() {
	n := big.NewFloat(1234.5678)
	mp := multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetBigFloat(n)

	f, _ := n.Float64()
	fmt.Printf("%g\n", f)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// 1234.5678
	// Base:10
	// Power:3
	// Symbol:k
	// Name:kilo
}

func ExampleDecimalMultiplePrefix_SetBigInt() {
	n := big.NewInt(-2 * 1000 * 1000)
	mp := multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetBigInt(n)

	fmt.Printf("%d\n", n.Int64())
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// -2000000
	// Base:10
	// Power:6
	// Symbol:M
	// Name:mega
}

func ExampleDecimalMultiplePrefix_SetBigRat() {
	n := big.NewRat(-4*1000*1000, 2)
	mp := multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetBigRat(n)

	f, _ := n.Float64()

	fmt.Printf("%g\n", f)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// -2e+06
	// Base:10
	// Power:6
	// Symbol:M
	// Name:mega
}

func ExampleBinaryMultiplePrefix_SetFloat64() {
	n := 1234.56789
	mp := multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetFloat64(n)

	fmt.Printf("%g\n", n)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// 1234.56789
	// Base:2
	// Power:10
	// Symbol:Ki
	// Name:kibi
}

func ExampleBinaryMultiplePrefix_SetInt64() {
	n := int64(-2 * 1024 * 1024)
	mp := multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetInt64(n)

	fmt.Printf("%d\n", n)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// -2097152
	// Base:2
	// Power:20
	// Symbol:Mi
	// Name:mebi
}

func ExampleBinaryMultiplePrefix_SetUint64() {
	n := uint64(2 * 1024 * 1024)
	mp := multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetUint64(n)

	fmt.Printf("%d\n", n)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// 2097152
	// Base:2
	// Power:20
	// Symbol:Mi
	// Name:mebi
}

func ExampleBinaryMultiplePrefix_SetBigFloat() {
	n := big.NewFloat(1234.5678)
	mp := multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetBigFloat(n)

	f, _ := n.Float64()
	fmt.Printf("%g\n", f)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// 1234.5678
	// Base:2
	// Power:10
	// Symbol:Ki
	// Name:kibi
}

func ExampleBinaryMultiplePrefix_SetBigInt() {
	n := big.NewInt(-2 * 1024 * 1024)
	mp := multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetBigInt(n)

	fmt.Printf("%d\n", n.Int64())
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// -2097152
	// Base:2
	// Power:20
	// Symbol:Mi
	// Name:mebi
}

func ExampleBinaryMultiplePrefix_SetBigRat() {
	n := big.NewRat(-4*1024*1024, 2)
	mp := multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetBigRat(n)

	f, _ := n.Float64()

	fmt.Printf("%g\n", f)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// -2.097152e+06
	// Base:2
	// Power:20
	// Symbol:Mi
	// Name:mebi
}

func ExampleDecimalMultiplePrefix_SetPrefix() {
	prefix := "kilo"
	mp := multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetPrefix(prefix)

	fmt.Printf("%s\n", prefix)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())

	prefix = "k"
	mp = multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetPrefix(prefix)

	fmt.Printf("%s\n", prefix)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// kilo
	// Base:10
	// Power:3
	// Symbol:k
	// Name:kilo
	// k
	// Base:10
	// Power:3
	// Symbol:k
	// Name:kilo
}

func ExampleBinaryMultiplePrefix_SetPrefix() {
	prefix := "kibi"
	mp := multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetPrefix(prefix)

	fmt.Printf("%s\n", prefix)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())

	prefix = "Ki"
	mp = multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetPrefix(prefix)

	fmt.Printf("%s\n", prefix)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// kibi
	// Base:2
	// Power:10
	// Symbol:Ki
	// Name:kibi
	// Ki
	// Base:2
	// Power:10
	// Symbol:Ki
	// Name:kibi
}

func ExampleDecimalMultiplePrefix_SetPower() {
	p := 3
	mp := multiple_prefix.DecimalMultiplePrefixTODO.Copy().SetPower(p)

	fmt.Printf("%d\n", p)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// 3
	// Base:10
	// Power:3
	// Symbol:k
	// Name:kilo

}

func ExampleBinaryMultiplePrefix_SetPower() {
	p := 10
	mp := multiple_prefix.BinaryMultiplePrefixTODO.Copy().SetPower(p)

	fmt.Printf("%d\n", p)
	fmt.Printf("Base:%d\n", mp.Base())
	fmt.Printf("Power:%d\n", mp.Power())
	fmt.Printf("Symbol:%s\n", mp.Symbol())
	fmt.Printf("Name:%s\n", mp.Name())
	// Output:
	// 10
	// Base:2
	// Power:10
	// Symbol:Ki
	// Name:kibi
}

/*
// for go vet

func ExampleMultiplePrefix_Base() {
	fmt.Printf("%d\n", multiple_prefix.DecimalMultiplePrefixTODO.Copy().Base())
	fmt.Printf("%d\n", multiple_prefix.BinaryMultiplePrefixTODO.Copy().Base())

	// Output:
	// 10
	// 2
}

func ExampleMultiplePrefix_Factor() {
	// Decimal
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixYocto.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixZepto.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixAtto.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixFemto.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixPico.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixNano.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixMicro.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixMilli.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixDeci.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixCenti.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixOne.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixHecto.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixDeka.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixKilo.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixMega.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixGiga.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixTera.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixPeta.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixExa.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixZetta.Copy().Factor())
	fmt.Printf("%g\n", multiple_prefix.DecimalMultiplePrefixYotta.Copy().Factor())

	// Binary
	fmt.Printf("%b\n", int64(multiple_prefix.BinaryMultiplePrefixOne.Copy().Factor()))
	fmt.Printf("%b\n", int64(multiple_prefix.BinaryMultiplePrefixKibi.Copy().Factor()))
	fmt.Printf("%b\n", int64(multiple_prefix.BinaryMultiplePrefixMebi.Copy().Factor()))
	fmt.Printf("%b\n", int64(multiple_prefix.BinaryMultiplePrefixGibi.Copy().Factor()))
	fmt.Printf("%b\n", int64(multiple_prefix.BinaryMultiplePrefixTebi.Copy().Factor()))
	fmt.Printf("%b\n", int64(multiple_prefix.BinaryMultiplePrefixPebi.Copy().Factor()))
	fmt.Printf("%b\n", int64(multiple_prefix.BinaryMultiplePrefixExbi.Copy().Factor()))

	// Output:
	// 1.0000000000000001e-24
	// 1e-21
	// 1e-18
	// 1e-15
	// 1e-12
	// 1e-09
	// 1e-06
	// 0.001
	// 0.01
	// 0.1
	// 1
	// 10
	// 100
	// 1000
	// 1e+06
	// 1e+09
	// 1e+12
	// 1e+15
	// 1e+18
	// 1e+19
	// 1e+21
	// 1
	// 10000000000
	// 100000000000000000000
	// 1000000000000000000000000000000
	// 10000000000000000000000000000000000000000
	// 100000000000000000000000000000000000000000000000000
	// 1000000000000000000000000000000000000000000000000000000000000
}

func ExampleMultiplePrefix_Name() {
	// Decimal
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixYocto.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixZepto.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixAtto.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixFemto.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixPico.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixNano.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixMicro.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixMilli.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixDeci.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixCenti.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixOne.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixHecto.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixDeka.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixKilo.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixMega.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixGiga.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixTera.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixPeta.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixExa.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixZetta.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixYotta.Copy().Name())

	// Binary
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixOne.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixKibi.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixMebi.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixGibi.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixTebi.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixPebi.Copy().Name())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixExbi.Copy().Name())

	// Output:
	// yocto
	// atto
	// zepto
	// femto
	// pico
	// nano
	// micro
	// milli
	// deci
	// centi
	//
	// hecto
	// deka
	// kilo
	// mega
	// giga
	// tera
	// peta
	// exa
	// zetta
	// yotta
	//
	// kibi
	// mebi
	// gibi
	// tebi
	// pebi
	// exbi
}

func ExampleMultiplePrefix_Symbol() {
	// Decimal
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixYocto.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixZepto.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixAtto.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixFemto.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixPico.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixNano.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixMicro.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixMilli.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixDeci.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixCenti.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixOne.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixHecto.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixDeka.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixKilo.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixMega.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixGiga.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixTera.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixPeta.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixExa.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixZetta.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.DecimalMultiplePrefixYotta.Copy().Symbol())

	// Binary
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixOne.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixKibi.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixMebi.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixGibi.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixTebi.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixPebi.Copy().Symbol())
	fmt.Printf("%s\n", multiple_prefix.BinaryMultiplePrefixExbi.Copy().Symbol())

	// Output:
	// y
	// z
	// a
	// f
	// p
	// n
	// μ
	// m
	// m
	// m
	//
	// h
	// da
	// k
	// M
	// G
	// T
	// P
	// E
	// Z
	// Y
	//
	// Ki
	// Mi
	// Gi
	// Ti
	// Pi
	// Ei
}
*/
