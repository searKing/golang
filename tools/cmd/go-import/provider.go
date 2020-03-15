// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"sync"
	"text/template"
)

func importTmplProvider(name string) func() *template.Template {
	var tmplCache *template.Template
	var tmplCacheOnce sync.Once
	tmplProvider := func() *template.Template {
		tmplCacheOnce.Do(func() {
			tmplCache = template.Must(template.New(name).Parse(string(MustAsset(name))))
		})
		return tmplCache
	}
	return tmplProvider
}
