// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)
import os_ "github.com/searKing/golang/go/os"

func ExampleNewRotateFile() {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := os_.NewRotateFile("log/test.2006-01-02-15-04-05.log")
	defer file.Close()
	file.FilePathPrefix = dir + string(filepath.Separator)
	file.MaxCount = 5
	file.RotateInterval = 5 * time.Second
	file.MaxAge = time.Hour
	file.FileLinkPath = filepath.Join(dir, "log", "s.log")
	for i := range 10000 {
		time.Sleep(1 * time.Millisecond)
		file.WriteString(time.Now().String())
		if err := file.Rotate(false); err != nil {
			fmt.Printf("%d, err: %v\n", i, err)
		}
	}
}

func ExampleNewRotateFileWithStrftime() {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := os_.NewRotateFileWithStrftime("log/test.%Y-%m-%d-%H-%M-%S.log")
	defer file.Close()
	file.FilePathPrefix = dir + string(filepath.Separator)
	file.MaxCount = 5
	file.RotateInterval = 5 * time.Second
	file.MaxAge = time.Hour
	file.FileLinkPath = filepath.Join(dir, "log", "s.log")
	for i := range 10000 {
		time.Sleep(1 * time.Millisecond)
		file.WriteString(time.Now().String())
		if err := file.Rotate(false); err != nil {
			fmt.Printf("%d, err: %v\n", i, err)
		}
	}
}

func ExampleDiskUsage() {
	total, free, avail, inodes, inodesFree, err := os_.DiskUsage("/tmp")
	if err != nil {
		return
	}

	fmt.Printf("total :%d B, free: %d B, avail: %d B, inodes: %d, inodesFree: %d", total, free, avail, inodes, inodesFree)
	// total :499963174912 B, free: 57534603264 B, avail: 57534603264 B, inodes: 566386444, inodesFree: 561861360
}

func ExampleReadDirN() {
	files, err := os_.ReadDirN(".", 1)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func ExampleNewCacheFile() {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := os_.NewCacheFile(os_.WithCacheFileBucketRootDir(filepath.Join(dir, "log")),
		os_.WithCacheFileCacheExpiredAfter(10*time.Millisecond),
		os_.WithCacheFileBucketKeyFunc(func(url string) string {
			return "always conflict key"
		}))

	for i := range 10000 {
		time.Sleep(1 * time.Millisecond)
		_, _, err := file.Put(fmt.Sprintf("cache%d", i), strings.NewReader(strconv.Itoa(i)))
		if err != nil {
			fmt.Printf("%d, err: %v\n", i, err)
		}
	}
}
