package json

import (
	"encoding/json"
	"io/ioutil"
)

func ReadConfigFile(name string, v interface{}) error {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
