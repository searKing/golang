// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"flag"
	"fmt"
	"sync"
)

// StringSliceVarWithFlagSet defines a []string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
func StringSliceVarWithFlagSet(f *flag.FlagSet, p *[]string, name string, value []string, usage string) {
	f.Var(newStringSliceValue(value, p), name, usage)
}

// StringSliceWithFlagSet defines an []string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
func StringSliceWithFlagSet(f *flag.FlagSet, name string, value []string, usage string) *[]string {
	p := new([]string)
	StringSliceVarWithFlagSet(f, p, name, value, usage)
	return p
}

// StringSliceVar defines a []string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
func StringSliceVar(p *[]string, name string, value []string, usage string) {
	flag.CommandLine.Var(newStringSliceValue(value, p), name, usage)
}

// StringSlice defines an []string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
func StringSlice(name string, value []string, usage string) *[]string {
	return StringSliceWithFlagSet(flag.CommandLine, name, value, usage)
}

// stringSliceValue is a flag.Value that accumulates strings.
// e.g. --flag=one --flag=two would produce []string{"one", "two"}.
type stringSliceValue struct {
	vars *[]string
	once sync.Once
}

func newStringSliceValue(val []string, p *[]string) *stringSliceValue {
	*p = val
	return &stringSliceValue{vars: p}
	//return (*stringSliceValue)(p)
}

func (ss *stringSliceValue) Get() interface{} { return *ss.vars }

func (ss *stringSliceValue) String() string { return fmt.Sprintf("%q", ss.vars) }

func (ss *stringSliceValue) Set(s string) error {
	ss.once.Do(func() {
		*(ss.vars) = nil
	})
	*(ss.vars) = append(*(ss.vars), s)
	return nil
}
