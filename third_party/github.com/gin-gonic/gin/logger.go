// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const DateTime = time.DateTime

// LogFormatter is the log format function [gin.Logger] middleware uses.
func LogFormatter(layout string) func(param gin.LogFormatterParams) string {
	return LogFormatterWithExtra(layout, nil)
}

// LogFormatterWithExtra is the log format function [gin.Logger] middleware uses with extra append {path}.
func LogFormatterWithExtra(layout string, getExtra func(param gin.LogFormatterParams) string) func(param gin.LogFormatterParams) string {
	return func(param gin.LogFormatterParams) string {
		// code borrowed from https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L141
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		var extra string
		if getExtra != nil {
			extra = getExtra(param)
		}
		if layout == "" {
			return fmt.Sprintf("[GIN] %3d %s| %13v | %15s |%s %-7s %s %#v%s\n%s",
				param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				extra,
				param.ErrorMessage,
			)
		}
		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v%s\n%s",
			param.TimeStamp.Format(layout),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			extra,
			param.ErrorMessage,
		)
	}
}
