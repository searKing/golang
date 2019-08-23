package flag_test

import (
	"flag"
	"os"
)

// Additional routines compiled into the package only during testing.

var DefaultUsage = flag.Usage

// ResetForTesting clears all flag state and sets the usage function as directed.
// After calling ResetForTesting, parse errors in flag handling will not
// exit the program.
func ResetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.Usage = flag.Usage
	flag.Usage = usage
}
