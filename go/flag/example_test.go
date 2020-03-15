// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag_test

import (
	"flag"
	"fmt"
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
	fmt.Printf("infos before parse: %q\n", infos)
	_ = fs.Parse([]string{"-i", "golang", "-i", "flag", "-i", "string slice"})
	fmt.Printf("infos after parse: %q\n", infos)
	// Output:
	// -i value
	//     	info arrays (default &["hello" "world"])
	// infos before parse: ["hello" "world"]
	// infos after parse: ["golang" "flag" "string slice"]

}
