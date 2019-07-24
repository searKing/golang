package json

import (
	"encoding/json"
	"io/ioutil"
)

const permissions = 0666

func WriteConfigFile(name string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, data, permissions)
}
