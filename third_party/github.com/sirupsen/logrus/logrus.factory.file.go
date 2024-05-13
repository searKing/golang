// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"os"
	"time"

	time_ "github.com/searKing/golang/go/time"
	"github.com/sirupsen/logrus"
)

// file
type file struct {
	Log struct {
		Level                   logrus.Level   `json:"level,omitempty" yaml:"level"`                                           // sets the logger level, E.g. panic, fatal, error, warn, info, debug, trace
		Format                  Format         `json:"format,omitempty" yaml:"format"`                                         // sets the logger format, E.g. json, text, glog, glog_human
		Path                    string         `json:"path,omitempty" yaml:"path"`                                             // sets the log file path prefix, E.g. "./log/" + filepath.Base(os.Args[0]).
		RotationDuration        time_.Duration `json:"rotation_duration,omitempty" yaml:"rotation_duration"`                   // Rotate files are rotated until RotateInterval expired before being removed, E.g.
		RotationSizeInByte      int64          `json:"rotation_size_in_byte,omitempty" yaml:"rotation_size_in_byte"`           // 日志循环最大分片大小,单位为Byte
		RotationMaxCount        int            `json:"rotation_max_count,omitempty" yaml:"rotation_max_count"`                 // 日志循环覆盖保留分片个数
		RotationMaxAge          time_.Duration `json:"rotation_max_age,omitempty" yaml:"rotation_max_age"`                     // 文件最大保存时间
		ReportCaller            bool           `json:"report_caller,omitempty" yaml:"report_caller"`                           // 调用者堆栈
		MuteDirectlyOutput      bool           `json:"mute_directly_output,omitempty" yaml:"mute_directly_output"`             // warn及更高级别日志是否打印到标准输出
		MuteDirectlyOutputLevel logrus.Level   `json:"mute_directly_output_level,omitempty" yaml:"mute_directly_output_level"` // 标准输出日志最低打印等级
		TruncateMessageSizeTo   int            `json:"truncate_message_size_to,omitempty" yaml:"truncate_message_size_to"`     // 日志 message 最大长度，超长则截断; 当前仅glog和glog_human模式生效
		TruncateKeySizeTo       int            `json:"truncate_key_size_to,omitempty" yaml:"truncate_key_size_to"`             // 日志键值对的key最大长度，超长则截断; 当前仅glog和glog_human模式生效
		TruncateValueSizeTo     int            `json:"truncate_value_size_to,omitempty" yaml:"truncate_value_size_to"`         // 日志键值对的value最大长度，超长则截断; 当前仅glog和glog_human模式生效
	} `json:"log,omitempty" yaml:"log"`
}

// NewFactoryFromFile reads factory from file, parsed by unmarshal
func NewFactoryFromFile(name string, unmarshal func(data []byte, v any) error) (Factory, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return Factory{}, err
	}

	var def FactoryConfig
	def.SetDefaults()

	var f = file{
		Log: struct {
			Level                   logrus.Level   `json:"level,omitempty" yaml:"level"`
			Format                  Format         `json:"format,omitempty" yaml:"format"`
			Path                    string         `json:"path,omitempty" yaml:"path"`
			RotationDuration        time_.Duration `json:"rotation_duration,omitempty" yaml:"rotation_duration"`
			RotationSizeInByte      int64          `json:"rotation_size_in_byte,omitempty" yaml:"rotation_size_in_byte"`
			RotationMaxCount        int            `json:"rotation_max_count,omitempty" yaml:"rotation_max_count"`
			RotationMaxAge          time_.Duration `json:"rotation_max_age,omitempty" yaml:"rotation_max_age"`
			ReportCaller            bool           `json:"report_caller,omitempty" yaml:"report_caller"`
			MuteDirectlyOutput      bool           `json:"mute_directly_output,omitempty" yaml:"mute_directly_output"`
			MuteDirectlyOutputLevel logrus.Level   `json:"mute_directly_output_level,omitempty" yaml:"mute_directly_output_level"`
			TruncateMessageSizeTo   int            `json:"truncate_message_size_to,omitempty" yaml:"truncate_message_size_to"`
			TruncateKeySizeTo       int            `json:"truncate_key_size_to,omitempty" yaml:"truncate_key_size_to"`
			TruncateValueSizeTo     int            `json:"truncate_value_size_to,omitempty" yaml:"truncate_value_size_to"`
		}{
			Level:                   def.Level,
			Format:                  def.Format,
			Path:                    def.Path,
			RotationDuration:        time_.Duration(def.RotationDuration),
			RotationSizeInByte:      def.RotationSizeInByte,
			RotationMaxCount:        def.RotationMaxCount,
			RotationMaxAge:          time_.Duration(def.RotationMaxAge),
			ReportCaller:            def.ReportCaller,
			MuteDirectlyOutput:      def.MuteDirectlyOutput,
			MuteDirectlyOutputLevel: def.MuteDirectlyOutputLevel,
			TruncateMessageSizeTo:   def.TruncateMessageSizeTo,
			TruncateKeySizeTo:       def.TruncateKeySizeTo,
			TruncateValueSizeTo:     def.TruncateValueSizeTo,
		},
	}
	err = unmarshal(data, &f)
	if err != nil {
		return Factory{}, err
	}

	return NewFactory(FactoryConfig{
		Level:                   f.Log.Level,
		Format:                  f.Log.Format,
		Path:                    f.Log.Path,
		RotationDuration:        time.Duration(f.Log.RotationDuration),
		RotationSizeInByte:      f.Log.RotationSizeInByte,
		RotationMaxCount:        f.Log.RotationMaxCount,
		RotationMaxAge:          time.Duration(f.Log.RotationMaxAge),
		ReportCaller:            f.Log.ReportCaller,
		MuteDirectlyOutput:      f.Log.MuteDirectlyOutput,
		MuteDirectlyOutputLevel: f.Log.MuteDirectlyOutputLevel,
		TruncateMessageSizeTo:   f.Log.TruncateMessageSizeTo,
		TruncateKeySizeTo:       f.Log.TruncateKeySizeTo,
		TruncateValueSizeTo:     f.Log.TruncateValueSizeTo,
	}), err
}
