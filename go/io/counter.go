package io

import (
	"bufio"
	"bytes"
	"io"

	bytes_ "github.com/searKing/golang/go/bytes"
)

// Count counts the number of non-overlapping instances of sep in r.
// If sep is an empty slice, Count returns 1 + the number of UTF-8-encoded code points in s.
// tailMatch returns true if tail bytes match sep
func Count(r io.Reader, sep string) (cnt int, tailMatch bool, err error) {
	return CountSize(r, sep, bufio.MaxScanTokenSize)
}

// CountSize counts the number of non-overlapping instances of sep in r.
// tailMatch returns true if tail bytes match sep
func CountSize(r io.Reader, sep string, size int) (cnt int, tailMatch bool, err error) {
	// special case
	if len(sep) == 0 || len(sep) == 1 {
		return CountAnySize(r, sep, size)
	}

	if size < len(sep) {
		size = len(sep)
	}
	var count int
	buf := make([]byte, size)

	// buffered bytes
	var bufferedPos int
	for {
		n, err := r.Read(buf[bufferedPos:])
		if n > 0 {
			bufferedPos += n
			cnt, index := bytes_.CountIndex(buf[:bufferedPos], []byte(sep))
			count += cnt

			// store the next index to do match
			var tailIndex int
			if index >= 0 {
				// skip matched
				tailIndex = index + len(sep)
			} else {
				// skip first byte
				tailIndex = 1
			}

			copyCnt := bufferedPos - tailIndex
			if copyCnt >= len(sep) {
				copyCnt = len(sep) - 1
			}

			// buffer tail bytes if any
			if copyCnt > 0 {
				copy(buf[:copyCnt], buf[tailIndex:tailIndex+copyCnt])
				bufferedPos = copyCnt
			} else {
				bufferedPos = 0
			}
		}
		if err == io.EOF {
			return count, bufferedPos == 0, nil
		}
		if err != nil {
			return count, true, err
		}
	}
}

// CountAnySize counts the number of non-overlapping instances of sep in r.
// If sep is an empty slice, Count returns 1 + the number of UTF-8-encoded code points in s.
func CountAny(r io.Reader, sep string) (cnt int, tailMatch bool, err error) {
	return CountAnySize(r, sep, bufio.MaxScanTokenSize)
}

// CountAnySize counts the number of non-overlapping instances of sep in r.
// If sep is an empty slice, Count returns 1 + the number of UTF-8-encoded code points in s.
func CountAnySize(r io.Reader, sep string, size int) (cnt int, tailMatch bool, err error) {
	var count int
	if size < 1 {
		size = 1
	}

	buf := make([]byte, size)
	var lastByte byte
	for {
		n, err := r.Read(buf)
		if n > 0 {
			lastByte = buf[n-1]
			count += bytes.Count(buf[:n], []byte(sep))
		}
		if err == io.EOF {
			if bytes.ContainsAny([]byte(sep), string(lastByte)) {
				return count, true, nil
			}
			return count, false, nil
		}
		if err != nil {
			return count, true, err
		}
	}
}

// CountLines counts the number of lines by \n.
func CountLines(r io.Reader) (lines int, err error) {
	return CountLinesSize(r, bufio.MaxScanTokenSize)
}

// CountLinesSize counts the number of lines by \n.
func CountLinesSize(r io.Reader, size int) (lines int, err error) {
	cnt, tailMatch, err := CountSize(r, "\n", size)
	// take care of ending line without '\n'
	if !tailMatch {
		cnt++
	}
	return cnt, err
}
