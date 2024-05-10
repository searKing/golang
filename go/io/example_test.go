// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	io_ "github.com/searKing/golang/go/io"
)

func ExampleSniffReader() {
	r := strings.NewReader("MSG:some io.Reader stream to be read\n")
	sniff := io_.SniffReader(r)

	// start sniffing
	sniff.Sniff(true)
	// sniff "MSG:"
	printSniff(sniff, len("MSG:"))
	fmt.Printf("\n")

	// stop sniffing
	sniff.Sniff(false)
	printSniff(sniff, len("MSG:"))
	fmt.Printf("\n")

	// start sniffing again
	sniff.Sniff(true)
	// sniff "io.Reader"
	printSniff(sniff, len("some"))
	fmt.Printf("\n")

	// stop sniffing
	sniff.Sniff(false)
	printAll(sniff)

	// Output:
	// MSG:
	// MSG:
	// some
	// some io.Reader stream to be read
}
func ExampleReplayReader() {
	r := strings.NewReader("MSG:some io.Reader stream to be read")
	replayR := io_.ReplayReader(r)

	// print "MSG:"
	printSniff(replayR, len("MSG:"))
	fmt.Printf("\n")

	// start replay
	replayR.Replay()
	// print "MSG:"
	printSniff(replayR, len("MSG:"))
	fmt.Printf("\n")

	// start replay
	replayR.Replay()
	// print "MSG:"
	printAll(replayR)
	fmt.Printf("\n")

	// start replay
	replayR.Replay()
	// print "MSG:"
	printAll(replayR)
	fmt.Printf("\n")

	// Output:
	// MSG:
	// MSG:
	// MSG:some io.Reader stream to be read
	// MSG:some io.Reader stream to be read
}

func ExampleEOFReader() {
	r := io_.EOFReader()

	printAll(r)

	// Output:
	//
}

func ExampleWatchReader() {
	r := strings.NewReader("some io.Reader stream to be read\n")
	watch := io_.WatchReader(r, io_.WatcherFunc(func(p []byte, n int, err error) (int, error) {
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		fmt.Printf("%s", p[:n])
		return n, err
	}))

	printAll(watch)

	// Output:
	// some io.Reader stream to be read
	// some io.Reader stream to be read
}

func ExampleLimitReadSeeker() {
	r := strings.NewReader("some io.Reader stream to be read\n")
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}

	limit := io_.LimitReadSeeker(r, int64(len("some io.Reader stream")))

	_, _ = limit.Seek(int64(len("some io.Reader ")), io.SeekStart)
	if _, err := io.Copy(os.Stdout, limit); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n")

	_, _ = limit.Seek(int64(-len("stream to be read\n")), io.SeekEnd)
	if _, err := io.Copy(os.Stdout, limit); err != nil {
		log.Fatal(err)
	}

	// Output:
	// some io.Reader stream to be read
	// stream
	// stream
}

func ExampleDynamicReadSeeker() {
	r := strings.NewReader("some io.Reader stream to be read\n")
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}

	ignoreOff := len("some ")

	// dynamic behaves like a reader for "io.Reader stream to be read\n"
	dynamic := io_.DynamicReadSeeker(func(off int64) (reader io.Reader, e error) {
		if off >= 0 {
			off += int64(ignoreOff)
		}
		_, err := r.Seek(off, io.SeekStart)
		// to omit r's io.Seeker
		return io.MultiReader(r), err
	}, r.Size()-int64(ignoreOff))

	_, _ = dynamic.Seek(int64(len("io.Reader ")), io.SeekStart)
	if _, err := io.Copy(os.Stdout, dynamic); err != nil {
		log.Fatal(err)
	}

	_, _ = dynamic.Seek(int64(-len("stream to be read\n")), io.SeekEnd)
	if _, err := io.Copy(os.Stdout, dynamic); err != nil {
		log.Fatal(err)
	}

	// Output:
	// some io.Reader stream to be read
	// stream to be read
	// stream to be read
}

func ExampleCount() {
	cnt, tailMatch, err := io_.Count(bytes.NewReader([]byte("abcdef")), "b")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cnt: %d, tailMatch: %t\n", cnt, tailMatch)
	cnt, tailMatch, err = io_.Count(bytes.NewReader([]byte("abcdef")), "f")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cnt: %d, tailMatch: %t\n", cnt, tailMatch)
	cnt, tailMatch, err = io_.Count(bytes.NewReader([]byte("abcdef")), "cd")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cnt: %d, tailMatch: %t\n", cnt, tailMatch)
	cnt, tailMatch, err = io_.Count(bytes.NewReader([]byte("abcdef")), "cb")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cnt: %d, tailMatch: %t\n", cnt, tailMatch)

	// Output:
	// cnt: 1, tailMatch: false
	// cnt: 1, tailMatch: true
	// cnt: 1, tailMatch: false
	// cnt: 0, tailMatch: false
}

func ExampleCountSize() {
	cnt, tailMatch, err := io_.CountSize(bytes.NewReader([]byte("abcdef")), "b", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cnt: %d, tailMatch: %t\n", cnt, tailMatch)
	cnt, tailMatch, err = io_.CountSize(bytes.NewReader([]byte("abcdef")), "f", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cnt: %d, tailMatch: %t\n", cnt, tailMatch)
	cnt, tailMatch, err = io_.CountSize(bytes.NewReader([]byte("abcdef")), "cd", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cnt: %d, tailMatch: %t\n", cnt, tailMatch)
	cnt, tailMatch, err = io_.CountSize(bytes.NewReader([]byte("abcdef")), "cb", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cnt: %d, tailMatch: %t\n", cnt, tailMatch)

	// Output:
	// cnt: 1, tailMatch: false
	// cnt: 1, tailMatch: true
	// cnt: 1, tailMatch: false
	// cnt: 0, tailMatch: false
}

func ExampleCountLines() {
	cnt, err := io_.CountLines(bytes.NewReader([]byte("abc\ndef")))
	if err != nil {
		log.Fatal(err)
	}
	if cnt != 2 {
		log.Fatalf("got %d, want 2", cnt)
	}
	fmt.Printf("cnt: %d\n", cnt)

	// Output:
	// cnt: 2
}

func printSniff(r io.Reader, n int) {
	b := make([]byte, n)
	n, err := r.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b[:n])
}

func printAll(r io.Reader) {
	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
}
