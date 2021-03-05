package viper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	errors_ "github.com/searKing/golang/go/errors"
	proto_ "github.com/searKing/golang/third_party/github.com/golang/protobuf/proto"
)

// merge sequence: protos..., file, env, replace if member has been set
func MergeAll(v *viper.Viper, cfgFile string, envPrefix string, protos ...proto.Message) error {
	// read default config from protobuf
	if err := MergeConfigFromProtoMessages(v, "", protos...); err != nil {
		return fmt.Errorf("merge config from proto messages: %w", err)
	}
	// read from file
	if err := MergeConfigFromFile(v, cfgFile); err != nil {
		return fmt.Errorf("merge config from file[%s]: %w", cfgFile, err)
	}

	// read in environment variables that match
	MergeConfigFromENV(v, envPrefix)
	return nil
}

// read from file
func MergeConfigFromFile(v *viper.Viper, cfgFile string) error {
	if cfgFile == "" {
		return nil
	}
	// enable ability to specify config file via flag
	v.SetConfigFile(cfgFile)

	return v.MergeInConfig()
}

// read from env
func MergeConfigFromENV(v *viper.Viper, envPrefix string) {
	// read in environment variables that match
	v.AutomaticEnv()          // read in environment variables that match
	v.SetEnvPrefix(envPrefix) // will be uppercase automatically
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// read from protobuf
// merge protos into viper one by one, replace if member has been set
// that is, later proto appeared, higher priority proto has
func MergeConfigFromProtoMessages(v *viper.Viper, configType string, protos ...proto.Message) error {
	v.SetConfigType("yaml")
	defer v.SetConfigType(configType)
	var errs []error
	for _, p := range protos {
		protoMap, err := proto_.ToGolangMap(p)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		protoBytes, err := yaml.Marshal(protoMap)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// merge if not exist
		if err := v.MergeConfig(bytes.NewReader(protoBytes)); err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return errors_.Multi(errs...)
}
