// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goimport

import (
	"embed"
	"sync"
	"text/template"
)

//go:embed tmpl
var tmplFS embed.FS

func importTmplProvider(name string) func() *template.Template {
	var tmplCache *template.Template
	var tmplCacheOnce sync.Once
	tmplProvider := func() *template.Template {
		tmplCacheOnce.Do(func() {
			text, err := tmplFS.ReadFile(name)
			if err != nil {
				panic(err)
			}
			tmplCache = template.Must(template.New(name).Parse(string(text)))
		})
		return tmplCache
	}
	return tmplProvider
}
