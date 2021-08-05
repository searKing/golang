// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

// ResolveOneAddrOptionWithPickerOption append PickOption to picker
func ResolveOneAddrOptionWithPickerOption(opts ...PickOption) ResolveOneAddrOption {
	return ResolveOneAddrOptionFunc(func(opt *resolveOneAddr) {
		opt.Picker = append(opt.Picker, opts...)
	})
}
