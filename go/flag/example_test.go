package flag_test

import (
	"flag"
	"os"

	flag_ "github.com/searKing/golang/go/flag"
)

func ExampleStringSliceVar() {
	var infos []string

	ResetForTesting(nil)
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	flag_.StringSliceVarWithFlagSet(fs, &infos, "i", []string{"hello", "world"}, "info arrays")
	fs.PrintDefaults()
	// Output:
	// -i value
	//     	info arrays (default &["hello" "world"])

}
