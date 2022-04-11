// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logs

import (
	"os"
	"path/filepath"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

// SetDefaults assigns default values for the Config
func (c *Config) SetDefaults() {
	c.Proto = Log{
		Level:                   Log_info,
		Format:                  Log_glog_human,
		Path:                    "./log/" + filepath.Base(os.Args[0]),
		RotationDuration:        durationpb.New(24 * time.Hour),
		RotationMaxCount:        0,
		RotationMaxAge:          durationpb.New(7 * 24 * time.Hour),
		RotationSizeInByte:      0,
		ReportCaller:            false,
		MuteDirectlyOutput:      true,
		MuteDirectlyOutputLevel: Log_warn,
	}
}
