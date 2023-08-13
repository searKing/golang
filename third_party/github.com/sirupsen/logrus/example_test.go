// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"encoding/json"
	"log"

	"github.com/sirupsen/logrus"
)

func ExampleNewFactoryFromFile() {
	fromFile, err := NewFactoryFromFile("./testdata/log.json", json.Unmarshal)
	//NewFactoryFromFile("./testdata/log.yaml", yaml.Unmarshal)
	if err != nil {
		log.Fatalf("read file failed: %s", err)
		return
	}
	err = fromFile.Apply()
	if err != nil {
		log.Fatalf("read file failed: %s", err)
	}

	// [INFO ][20230814 00:28:27.055471] [70427] [logrus.factory.slog.go:94](Apply) add rotation wrapper for log, path=./testdata/log/example, mute_directly_output=true, rotate_size_in_byte=0, duration=1h0m0s, max_age=24h0m0s, max_count=0
	// [INFO ][20230814 00:28:27.056335] [70427] [example_test.go:26](ExampleNewFactoryFromFile) Hello World
	// [WARN ][20230814 00:28:27.056381] [70427] [example_test.go:27](ExampleNewFactoryFromFile) Hello World
	logrus.Infof("Hello World")
	logrus.Warnf("Hello World")

	// Output:
}
