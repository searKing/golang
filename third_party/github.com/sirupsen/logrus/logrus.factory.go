// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Format 日志格式
//
//go:generate stringer -type=Format
type Format int32

const (
	FormatJson      Format = 0
	FormatText      Format = 1
	FormatGlog      Format = 2
	FormatGlogHuman Format = 3
)

// Convert the Level to a string. E.g. FormatJson becomes "json".
func (f Format) String() string {
	if b, err := f.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

// ParseFormat takes a string format and returns the Logrus log format constant.
func ParseFormat(lvl string) (Format, error) {
	switch strings.ToLower(lvl) {
	case "json":
		return FormatJson, nil
	case "text":
		return FormatText, nil
	case "glog":
		return FormatGlog, nil
	case "glog_human":
		return FormatGlogHuman, nil
	}

	var l Format
	return l, fmt.Errorf("not a valid logrus Format: %q", lvl)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *Format) UnmarshalText(text []byte) error {
	l, err := ParseFormat(string(text))
	if err != nil {
		return err
	}

	*f = l

	return nil
}

func (f Format) MarshalText() ([]byte, error) {
	switch f {
	case FormatJson:
		return []byte("json"), nil
	case FormatText:
		return []byte("text"), nil
	case FormatGlog:
		return []byte("glog"), nil
	case FormatGlogHuman:
		return []byte("glog_human"), nil
	}

	return nil, fmt.Errorf("not a valid logrus format %d", f)
}

// FactoryConfig 日志工厂函数配置
type FactoryConfig struct {
	Level                   logrus.Level  `json:"level,omitempty" yaml:"level"`                                           // sets the logger level, E.g. panic, fatal, error, warn, info, debug, trace
	Format                  Format        `json:"format,omitempty" yaml:"format"`                                         // sets the logger format, E.g. json, text, glog, glog_human
	Path                    string        `json:"path,omitempty" yaml:"path"`                                             // sets the log file path prefix, E.g. "./log/" + filepath.Base(os.Args[0]).
	RotationDuration        time.Duration `json:"rotation_duration,omitempty" yaml:"rotation_duration"`                   // Rotate files are rotated until RotateInterval expired before being removed, E.g.
	RotationSizeInByte      int64         `json:"rotation_size_in_byte,omitempty" yaml:"rotation_size_in_byte"`           // 日志循环最大分片大小,单位为Byte
	RotationMaxCount        int           `json:"rotation_max_count,omitempty" yaml:"rotation_max_count"`                 // 日志循环覆盖保留分片个数
	RotationMaxAge          time.Duration `json:"rotation_max_age,omitempty" yaml:"rotation_max_age"`                     // 文件最大保存时间
	ReportCaller            bool          `json:"report_caller,omitempty" yaml:"report_caller"`                           // 调用者堆栈
	MuteDirectlyOutput      bool          `json:"mute_directly_output,omitempty" yaml:"mute_directly_output"`             // warn及更高级别日志是否打印到标准输出
	MuteDirectlyOutputLevel logrus.Level  `json:"mute_directly_output_level,omitempty" yaml:"mute_directly_output_level"` // 标准输出日志最低打印等级
	TruncateMessageSizeTo   int           `json:"truncate_message_size_to,omitempty" yaml:"truncate_message_size_to"`     // 日志 message 最大长度，超长则截断; 当前仅glog和glog_human模式生效
	TruncateKeySizeTo       int           `json:"truncate_key_size_to,omitempty" yaml:"truncate_key_size_to"`             // 日志键值对的key最大长度，超长则截断; 当前仅glog和glog_human模式生效
	TruncateValueSizeTo     int           `json:"truncate_value_size_to,omitempty" yaml:"truncate_value_size_to"`         // 日志键值对的value最大长度，超长则截断; 当前仅glog和glog_human模式生效
}

// SetDefaults sets sensible values for unset fields in config. This is
// exported for testing: Configs passed to repository functions are copied and have
// default values set automatically.
func (fc *FactoryConfig) SetDefaults() {
	fc.Level = logrus.InfoLevel
	fc.Format = FormatGlogHuman
	fc.Path = "./log/" + filepath.Base(os.Args[0])
	fc.RotationDuration = 24 * time.Hour
	fc.RotationMaxCount = 0
	fc.RotationMaxAge = 7 * 24 * time.Hour
	fc.RotationSizeInByte = 0
	fc.ReportCaller = false
	fc.MuteDirectlyOutput = true
	fc.MuteDirectlyOutputLevel = logrus.FatalLevel
}

type Factory struct {
	// it's better to keep FactoryConfig as a private attribute,
	// thanks to that we are always sure that our configuration is not changed in the not allowed way
	fc FactoryConfig
}

func NewFactory(fc FactoryConfig) Factory {
	return Factory{fc: fc}
}

func (f Factory) Config() FactoryConfig {
	return f.fc
}
