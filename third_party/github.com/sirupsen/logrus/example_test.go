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

	logrus.Infof("Hello World")
	logrus.Warnf("Hello World")
	// [INFO ] [20220616 14:41:08.978628] [15900] [example_test.go:23](ExampleNewFactoryFromFile) Hello World
	// [WARN ] [20220616 14:41:08.978745] [15900] [example_test.go:24](ExampleNewFactoryFromFile) Hello World

	// Output:
}
