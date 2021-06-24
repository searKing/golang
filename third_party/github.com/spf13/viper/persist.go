// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	filepath_ "github.com/searKing/golang/go/path/filepath"
)

// PersistConfig writes config using into .use.<name>.yaml
func PersistConfig(v *viper.Viper) error {
	configFileUsing := DefaultPersisConfigPath(v)
	return PersistConfigTo(v, configFileUsing)
}

// DefaultPersisConfigPath returns the given file name to persist as YAML
func DefaultPersisConfigPath(v *viper.Viper) string {
	// persist using config
	f := v.ConfigFileUsed() // ./conf/.sole.yaml
	if f == "" {
		return ".use.yaml"
	}
	dir := filepath.Dir(f)
	base := filepath.Base(f)
	ext := filepath.Ext(f)
	name := strings.TrimPrefix(strings.TrimSuffix(base, ext), ".")

	return filepath_.Pathify(filepath.Join(dir, ".use."+name+".yaml")) // /root/.use.sole.yaml
}

// PersistConfigTo writes the completed component config into the given file name as YAML
func PersistConfigTo(v *viper.Viper, filename string) error {
	if filename == "" {
		return fmt.Errorf("persist skiped, for no config file used")
	}
	err := v.WriteConfigAs(filename)
	if err != nil {
		return fmt.Errorf("write using config file[%s]: %w", filename, err)
	}
	return nil
}

func PersistGlobalConfig() error {
	return PersistConfig(viper.GetViper())
}
