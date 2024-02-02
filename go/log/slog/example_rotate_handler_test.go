// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/searKing/golang/go/log/slog/internal/slogtest"
)

func ExampleNewRotateTextHandler() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	path := "" // If path is empty, the default os.Stdout are used. path can be relative path or absolute path.
	var slogOpt slog.HandlerOptions
	slogOpt.ReplaceAttr = slogtest.RemoveTime

	// if set: path = /tmp/logs/slog.log
	// slog.log -> /tmp/logs/slog.20240202180000.log
	// /tmp/logs/slog.20240202170000.log
	// /tmp/logs/slog.20240202180000.log
	// ...
	rotateOpts := []RotateOption{
		WithRotateRotateInterval(time.Hour),
		WithRotateMaxCount(3),
		WithRotateMaxAge(24 * time.Hour),
		// Below is default options.
		// WithRotateFilePathRotateStrftime(".%Y%m%d%H%M%S.log"), // time layout in strftime format to format rotate file.
		// WithRotateFilePathRotateLayout(".20060102150405.log"), // time layout in golang format to format rotate file.
		// WithRotateFileLinkPath(filepath.Base(path) + ".log"), // the symbolic link name that gets linked to the current file name being used.
	}
	// ...
	{
		fmt.Printf("----rotate_text----\n")
		handler, err := NewRotateTextHandler(path,
			&slogOpt, // If opts is nil, the default options are used.
			rotateOpts...)
		if err != nil {
			panic("failed to create rotate text handler:" + err.Error())
		}
		logger := slog.New(handler)
		logger.Info("rotate text message")
	}

	// Output:
	// ----rotate_text----
	// level=INFO msg="rotate text message"
}

func ExampleNewRotateJSONHandler() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	path := "" // If path is empty, the default os.Stdout are used. path can be relative path or absolute path.
	var slogOpt slog.HandlerOptions
	slogOpt.ReplaceAttr = slogtest.RemoveTime

	// if set: path = /tmp/logs/slog.log
	// slog.log -> /tmp/logs/slog.20240202180000.log
	// /tmp/logs/slog.20240202170000.log
	// /tmp/logs/slog.20240202180000.log
	// ...
	rotateOpts := []RotateOption{
		WithRotateRotateInterval(time.Hour),
		WithRotateMaxCount(3),
		WithRotateMaxAge(24 * time.Hour),
		// Below is default options.
		// WithRotateFilePathRotateStrftime(".%Y%m%d%H%M%S.log"), // time layout in strftime format to format rotate file.
		// WithRotateFilePathRotateLayout(".20060102150405.log"), // time layout in golang format to format rotate file.
		// WithRotateFileLinkPath(filepath.Base(path) + ".log"), // the symbolic link name that gets linked to the current file name being used.
	}
	// ...
	{
		fmt.Printf("----rotate_json----\n")
		handler, err := NewRotateJSONHandler(path,
			&slogOpt, // If opts is nil, the default options are used.
			rotateOpts...)
		if err != nil {
			panic("failed to create rotate json handler:" + err.Error())
		}
		logger := slog.New(handler)
		logger.Info("rotate json message")
	}

	// Output:
	// ----rotate_json----
	// {"level":"INFO","msg":"rotate json message"}
}

func ExampleNewRotateGlogHandler() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	path := "" // If path is empty, the default os.Stdout are used. path can be relative path or absolute path.
	var slogOpt slog.HandlerOptions
	slogOpt.ReplaceAttr = slogtest.RemoveTime

	// if set: path = /tmp/logs/slog.log
	// slog.log -> /tmp/logs/slog.20240202180000.log
	// /tmp/logs/slog.20240202170000.log
	// /tmp/logs/slog.20240202180000.log
	// ...
	rotateOpts := []RotateOption{
		WithRotateRotateInterval(time.Hour),
		WithRotateMaxCount(3),
		WithRotateMaxAge(24 * time.Hour),
		// Below is default options.
		// WithRotateFilePathRotateStrftime(".%Y%m%d%H%M%S.log"), // time layout in strftime format to format rotate file.
		// WithRotateFilePathRotateLayout(".20060102150405.log"), // time layout in golang format to format rotate file.
		// WithRotateFileLinkPath(filepath.Base(path) + ".log"), // the symbolic link name that gets linked to the current file name being used.
	}
	// ...
	{
		fmt.Printf("----rotate_glog----\n")
		handler, err := NewRotateGlogHandler(path,
			&slogOpt, // If opts is nil, the default options are used.
			rotateOpts...)
		if err != nil {
			panic("failed to create rotate glog handler:" + err.Error())
		}
		logger := slog.New(handler)
		logger.Info("rotate glog message")
	}

	// Output:
	// ----rotate_glog----
	// I 0] rotate glog message
}

func ExampleNewRotateGlogHumanHandler() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	path := "" // If path is empty, the default os.Stdout are used. path can be relative path or absolute path.
	var slogOpt slog.HandlerOptions
	slogOpt.ReplaceAttr = slogtest.RemoveTime

	// if set: path = /tmp/logs/slog.log
	// slog.log -> /tmp/logs/slog.20240202180000.log
	// /tmp/logs/slog.20240202170000.log
	// /tmp/logs/slog.20240202180000.log
	// ...
	rotateOpts := []RotateOption{
		WithRotateRotateInterval(time.Hour),
		WithRotateMaxCount(3),
		WithRotateMaxAge(24 * time.Hour),
		// Below is default options.
		// WithRotateFilePathRotateStrftime(".%Y%m%d%H%M%S.log"), // time layout in strftime format to format rotate file.
		// WithRotateFilePathRotateLayout(".20060102150405.log"), // time layout in golang format to format rotate file.
		// WithRotateFileLinkPath(filepath.Base(path) + ".log"), // the symbolic link name that gets linked to the current file name being used.
	}
	// ...
	{
		fmt.Printf("----rotate_glog_human----\n")
		handler, err := NewRotateGlogHumanHandler(path,
			&slogOpt, // If opts is nil, the default options are used.
			rotateOpts...)
		if err != nil {
			panic("failed to create rotate human glog handler:" + err.Error())
		}
		logger := slog.New(handler)
		logger.Info("rotate glog_human message")
	}

	// Output:
	// ----rotate_glog_human----
	// [INFO ] [0] rotate glog_human message
}

func ExampleNewRotateHandler() {
	getPid = func() int { return 0 } // set pid to zero for test
	defer func() { getPid = os.Getpid }()

	path := "" // If path is empty, the default os.Stdout are used. path can be relative path or absolute path.
	var slogOpt slog.HandlerOptions
	slogOpt.ReplaceAttr = slogtest.RemoveTime

	// if set: path = /tmp/logs/slog.log
	// slog.log -> /tmp/logs/slog.20240202180000.log
	// /tmp/logs/slog.20240202170000.log
	// /tmp/logs/slog.20240202180000.log
	// ...
	rotateOpts := []RotateOption{
		WithRotateRotateInterval(time.Hour),
		WithRotateMaxCount(3),
		WithRotateMaxAge(24 * time.Hour),
		// Below is default options.
		// WithRotateFilePathRotateStrftime(".%Y%m%d%H%M%S.log"), // time layout in strftime format to format rotate file.
		// WithRotateFilePathRotateLayout(".20060102150405.log"), // time layout in golang format to format rotate file.
		// WithRotateFileLinkPath(filepath.Base(path) + ".log"), // the symbolic link name that gets linked to the current file name being used.
	}
	// ...
	{
		fmt.Printf("----rotate_text----\n")
		handler, err := NewRotateTextHandler(path,
			&slogOpt, // If opts is nil, the default options are used.
			rotateOpts...)
		if err != nil {
			panic("failed to create rotate text handler:" + err.Error())
		}
		logger := slog.New(handler)
		logger.Info("rotate text message")
	}
	// ...
	{
		fmt.Printf("----rotate_json----\n")
		handler, err := NewRotateJSONHandler(path,
			&slogOpt, // If opts is nil, the default options are used.
			rotateOpts...)
		if err != nil {
			panic("failed to create rotate json handler:" + err.Error())
		}
		logger := slog.New(handler)
		logger.Info("rotate json message")
	}
	// ...
	{
		fmt.Printf("----rotate_glog----\n")
		handler, err := NewRotateGlogHandler(path,
			&slogOpt, // If opts is nil, the default options are used.
			rotateOpts...)
		if err != nil {
			panic("failed to create rotate glog handler:" + err.Error())
		}
		logger := slog.New(handler)
		logger.Info("rotate glog message")
	}
	// ...
	{
		fmt.Printf("----rotate_glog_human----\n")
		handler, err := NewRotateGlogHumanHandler(path,
			&slogOpt, // If opts is nil, the default options are used.
			rotateOpts...)
		if err != nil {
			panic("failed to create rotate human glog handler:" + err.Error())
		}
		logger := slog.New(handler)
		logger.Info("rotate glog_human message")
	}
	// ...
	{
		fmt.Printf("----multi_rotate[text-json-glog-glog_human]----\n")
		logger := slog.New(MultiHandler(
			func() slog.Handler {
				handler, err := NewRotateTextHandler(path,
					&slogOpt, // If opts is nil, the default options are used.
					rotateOpts...)
				if err != nil {
					panic("failed to create rotate text handler:" + err.Error())
				}
				return handler
			}(),
			func() slog.Handler {
				handler, err := NewRotateJSONHandler(path,
					&slogOpt, // If opts is nil, the default options are used.
					rotateOpts...)
				if err != nil {
					panic("failed to create rotate json handler:" + err.Error())
				}
				return handler
			}(),
			func() slog.Handler {
				handler, err := NewRotateGlogHandler(path,
					&slogOpt, // If opts is nil, the default options are used.
					rotateOpts...)
				if err != nil {
					panic("failed to create rotate glog handler:" + err.Error())
				}
				return handler
			}(),
			func() slog.Handler {
				handler, err := NewRotateGlogHumanHandler(path,
					&slogOpt, // If opts is nil, the default options are used.
					rotateOpts...)
				if err != nil {
					panic("failed to create rotate glog_human handler:" + err.Error())
				}
				return handler
			}(),
		))
		logger.Info("rotate multi_rotate[text-json-glog-glog_human] message")
	}

	// Output:
	// ----rotate_text----
	// level=INFO msg="rotate text message"
	// ----rotate_json----
	// {"level":"INFO","msg":"rotate json message"}
	// ----rotate_glog----
	// I 0] rotate glog message
	// ----rotate_glog_human----
	// [INFO ] [0] rotate glog_human message
	// ----multi_rotate[text-json-glog-glog_human]----
	// level=INFO msg="rotate multi_rotate[text-json-glog-glog_human] message"
	// {"level":"INFO","msg":"rotate multi_rotate[text-json-glog-glog_human] message"}
	// I 0] rotate multi_rotate[text-json-glog-glog_human] message
	// [INFO ] [0] rotate multi_rotate[text-json-glog-glog_human] message
}
