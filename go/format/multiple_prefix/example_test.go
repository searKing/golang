package multiple_prefix_test

import (
	"fmt"

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
