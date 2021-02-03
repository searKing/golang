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
	// persist using config
	f := v.ConfigFileUsed() // ./conf/.sole.yaml
	if f == "" {
		return fmt.Errorf("persist skiped, for no config file used")
	}
	dir := filepath.Dir(f)
	base := filepath.Base(f)
	ext := filepath.Ext(f)
	name := strings.TrimPrefix(strings.TrimSuffix(base, ext), ".")

	configFileUsing := filepath.Join(dir, ".use."+name+".yaml") // /root/.use.sole.yaml

	err := v.WriteConfigAs(configFileUsing)
	if err != nil {
		return fmt.Errorf("write using config file[%s]: %w", filepath_.Pathify(configFileUsing), err)
	}
	return nil
}

func PersistGlobalConfig() error {
	return PersistConfig(viper.GetViper())
}
