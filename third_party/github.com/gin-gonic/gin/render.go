// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import (
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	os_ "github.com/searKing/golang/go/os"
	filepath_ "github.com/searKing/golang/go/path/filepath"
	"github.com/sirupsen/logrus"
)

var defaultFilePerm = os.FileMode(0664)
var defaultDirPerm = os.FileMode(0755)

func Render(loggerProvider func() logrus.FieldLogger, tmplFilePath string, htmlFilePath string, data interface{}) gin.HandlerFunc {
	var tmplCache *template.Template
	var tmplCacheOnce sync.Once
	var tmplCacheMutex sync.RWMutex
	// clear files to be generated
	_ = os.RemoveAll(htmlFilePath)
	// Setup the graceful shutdown handler (traps SIGINT and SIGTERM)
	go func() {
		var stopChan = make(chan os.Signal)

		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

		<-stopChan
		if err := os.RemoveAll(htmlFilePath); err != nil {
			return
		}
	}()
	tmplName := filepath.Base(tmplFilePath)

	tmplProvider := func() *template.Template {
		tmplCacheOnce.Do(func() {
			tmplCache = template.Must(template.New(tmplName).ParseFiles(tmplFilePath))
		})
		return tmplCache
	}

	return func(ctx *gin.Context) {
		var logger logrus.FieldLogger
		if loggerProvider != nil {
			logger = loggerProvider()
		}
		if logger == nil {
			logger = logrus.New()
		}
		logger = logger.WithField("tmpl", tmplFilePath)

		if ctx != nil {
			ctx.Header("Content-Type", "text/html")
			ctx.Header("charset", "UTF-8")
		}
		isHtmlFileExist := func() bool {
			tmplCacheMutex.RLock()
			defer tmplCacheMutex.RUnlock()
			isExist, _ := os_.PathExists(htmlFilePath)
			return isExist
		}
		if isHtmlFileExist() {
			if ctx != nil {
				ctx.File(htmlFilePath)
			}
			return
		}

		// generate html file
		tmplCacheMutex.Lock()
		defer tmplCacheMutex.Unlock()
		if isExist, _ := os_.PathExists(htmlFilePath); isExist {
			if ctx != nil {
				ctx.File(htmlFilePath)
			}
			return
		}

		if err := filepath_.TouchAll(htmlFilePath, defaultDirPerm); err != nil {
			logger.WithError(err).Error("Failed executing template")
			if ctx != nil {
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			}
			_ = os.RemoveAll(htmlFilePath)
			return
		}

		htmlFile, err := os.OpenFile(htmlFilePath, os.O_CREATE|os.O_WRONLY, defaultFilePerm)
		if err != nil {
			logger.WithError(err).Error("Failed executing template")

			if ctx != nil {
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			}
			_ = os.RemoveAll(htmlFilePath)
			return
		}
		defer htmlFile.Close()

		err = tmplProvider().ExecuteTemplate(htmlFile, "index", data)
		if err != nil {
			logger.WithError(err).Error("Failed executing template")
			if ctx != nil {
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			}

			_ = os.RemoveAll(htmlFilePath)

			return
		}

		if ctx != nil {
			ctx.File(htmlFilePath)
		}
		return
	}
}
