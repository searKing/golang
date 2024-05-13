// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin/render"
)

// TemplateHTML contains template reference and its name with given interface object.
type TemplateHTML struct {
	Template *template.Template
	Files    []string
	Glob     string

	FuncMap template.FuncMap
	Name    string // Data's Name in tmpl
	Data    any

	once   sync.Once
	Delims *render.Delims
}

var htmlContentType = []string{"text/html; charset=utf-8"}

// Render (TemplateHTML) executes template and writes its result with custom ContentType for response.
func (r TemplateHTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	r.once.Do(func() {
		if r.Template == nil {
			r.Template = template.New("")
		}

		if r.Delims != nil {
			r.Template.Delims(r.Delims.Left, r.Delims.Right)
		}

		if r.FuncMap != nil {
			r.Template.Funcs(r.FuncMap)
		}

		if len(r.Files) > 0 {
			r.Template = template.Must(r.Template.ParseFiles(r.Files...))
		}
		if r.Glob != "" {
			r.Template = template.Must(r.Template.ParseGlob(r.Glob))
		}
	})

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

// WriteContentType (TemplateHTML) writes TemplateHTML ContentType.
func (r TemplateHTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
