// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	strings_ "github.com/searKing/golang/go/strings"
	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Proto     Log
	Validator *validator.Validate

	viper     *viper.Viper
	viperKeys []string
}

type completedConfig struct {
	*Config

	// for Complete Only
	completeError error
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	var c Config
	c.SetDefaults()
	return &c
}

// NewViperConfig returns a Config struct with the global viper instance
// key representing a subtree of this instance.
// NewViperConfig is case-insensitive for a key.
func NewViperConfig(v *viper.Viper, keys ...string) *Config {
	c := NewConfig()
	c.viper = v
	c.viperKeys = keys
	return c
}

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	if c.viper != nil {
		err := viper_.UnmarshalKeysViper(c.viper, c.viperKeys, &c.Proto)
		if err != nil {
			return CompletedConfig{&completedConfig{completeError: err}}
		}
	}
	if c.Validator == nil {
		c.Validator = validator.New()
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

// Apply creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (c completedConfig) Apply() error {
	if c.completeError != nil {
		return c.completeError
	}
	err := c.Validate()
	if err != nil {
		return err
	}
	return c.install()
}

func (c *completedConfig) install() error {
	if c.Proto.GetFormat() == Log_json {
		logrus.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: logrus_.ShortCallerPrettyfier,
		})
	} else if c.Proto.GetFormat() == Log_text {
		logrus.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: logrus_.ShortCallerPrettyfier,
			DisableColors:    true,
		})
	} else if c.Proto.GetFormat() == Log_glog || c.Proto.GetFormat() == Log_glog_human {
		var formatter *logrus_.GlogFormatter
		if c.Proto.GetFormat() == Log_glog {
			formatter = logrus_.NewGlogFormatter()
		} else {
			formatter = logrus_.NewGlogEnhancedFormatter()
		}
		formatter.DisableColors = true

		var truncate = func(s string, n int) string {
			if len(s) <= n {
				return s
			}
			var buf strings.Builder
			buf.WriteString(fmt.Sprintf("size: %d, string: ", len(s)))
			buf.WriteString(strings_.Truncate(s, n))
			return buf.String()
		}

		if size := int(c.Proto.GetTruncateMessageSizeTo()); size > 0 {
			formatter.MessageStringFunc = func(value interface{}) string {
				stringVal, ok := value.(string)
				if !ok {
					stringVal = fmt.Sprint(value)
				}
				return truncate(stringVal, size)
			}
		}

		if size := int(c.Proto.GetTruncateKeySizeTo()); size > 0 {
			formatter.KeyStringFunc = func(key string) string {
				return truncate(key, size)
			}
		}

		if size := int(c.Proto.GetTruncateValueSizeTo()); size > 0 {
			formatter.ValueStringFunc = func(value interface{}) string {
				stringVal, ok := value.(string)
				if !ok {
					stringVal = fmt.Sprint(value)
				}
				return truncate(stringVal, size)
			}
		}

		logrus.SetFormatter(formatter)
	}

	level, err := logrus.ParseLevel(c.Proto.GetLevel().String())
	if err != nil {
		level = logrus.InfoLevel
		logrus.WithField("module", "log").WithField("log_level", c.Proto.GetLevel()).
			WithError(err).
			Warnf("malformed log level, use %s instead", level)
	}
	logrus.SetLevel(level)

	var RotateDuration = c.Proto.GetRotationDuration().AsDuration()
	var RotateMaxAge = c.Proto.GetRotationMaxAge().AsDuration()
	var RotateSizeInByte = c.Proto.GetRotationSizeInByte()
	var RotateMaxCount = int(c.Proto.GetRotationMaxCount())

	logrus.SetReportCaller(c.Proto.GetReportCaller())

	muteDirectlyOutputLogLevel, err := logrus.ParseLevel(c.Proto.GetMuteDirectlyOutputLevel().String())
	if err != nil {
		muteDirectlyOutputLogLevel = logrus.WarnLevel
		logrus.WithField("module", "log").
			WithField("mute_directly_output_log_level", c.Proto.GetMuteDirectlyOutputLevel()).
			WithError(err).
			Warnf("malformed log level, use %s instead", muteDirectlyOutputLogLevel)
	}

	if err := logrus_.WithRotate(logrus.StandardLogger(),
		c.Proto.GetPath(),
		logrus_.WithRotateInterval(RotateDuration),
		logrus_.WithMaxCount(RotateMaxCount),
		logrus_.WithMaxAge(RotateMaxAge),
		logrus_.WithRotateSize(RotateSizeInByte),
		logrus_.WithMuteDirectlyOutput(c.Proto.GetMuteDirectlyOutput()),
		logrus_.WithMuteDirectlyOutputLogLevel(muteDirectlyOutputLogLevel)); err != nil {
		logrus.WithField("module", "log").
			WithField("path", c.Proto.GetPath()).
			WithField("duration", RotateDuration).
			WithField("max_count", RotateMaxCount).
			WithField("max_age", RotateMaxAge).
			WithField("rotate_size_in_byte", RotateSizeInByte).
			WithField("mute_directly_output", c.Proto.GetMuteDirectlyOutput()).
			WithError(err).Error("add rotation wrapper for log")
		return err
	}
	logrus.WithField("module", "log").
		WithField("path", c.Proto.GetPath()).
		WithField("duration", RotateDuration).
		WithField("max_count", RotateMaxCount).
		WithField("max_age", RotateMaxAge).
		WithField("mute_directly_output", c.Proto.GetMuteDirectlyOutput()).
		Infof("add rotation wrapper for log")
	return nil
}
