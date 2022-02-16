// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	io_ "github.com/searKing/golang/go/io"
)

// The algorithm uses at most sniffLen bytes to make its decision.
const sniffLen = 512

// ContentType implements the algorithm described
// at https://mimesniff.spec.whatwg.org/ to determine the
// Content-Type of the given data. It considers at most the
// first 512 bytes of data from r. ContentType always returns
// a valid MIME type: if it cannot determine a more specific one, it
// returns "application/octet-stream".
// ContentType is based on http.DetectContentType.
func ContentType(r io.Reader, name string) (ctype string, bufferedContent io.Reader, err error) {
	ctype = mime.TypeByExtension(filepath.Ext(name))
	if ctype == "" && r != nil {
		// read a chunk to decide between utf-8 text and binary
		var buf [sniffLen]byte
		var n int
		if readSeeker, ok := r.(io.Seeker); ok {
			n, _ = io.ReadFull(r, buf[:])
			_, err = readSeeker.Seek(0, io.SeekStart) // rewind to output whole file
			if err != nil {
				err = errors.New("seeker can't seek")
				return "", r, err
			}
		} else {
			contentBuffer := bufio.NewReader(r)
			sniffed, err := contentBuffer.Peek(sniffLen)
			if err != nil {
				err = errors.New("reader can't read")
				return "", contentBuffer, err
			}
			n = copy(buf[:], sniffed)
			r = contentBuffer
		}
		ctype = http.DetectContentType(buf[:n])
	}
	return ctype, r, nil
}

// ServeContent replies to the request using the content in the
// provided Reader. The main benefit of ServeContent over io.Copy
// is that it handles Range requests properly, sets the MIME type, and
// handles If-Match, If-Unmodified-Since, If-None-Match, If-Modified-Since,
// and If-Range requests.
//
// If the response's Content-Type header is not set, ServeContent
// first tries to deduce the type from name's file extension and,
// if that fails, falls back to reading the first block of the content
// and passing it to DetectContentType.
// The name is otherwise unused; in particular it can be empty and is
// never sent in the response.
//
// If modtime is not the zero time or Unix epoch, ServeContent
// includes it in a Last-Modified header in the response. If the
// request includes an If-Modified-Since header, ServeContent uses
// modtime to decide whether the content needs to be sent at all.
//
// If the content's Seek method work: ServeContent uses
// a seek to the end of the content to determine its size, and the param size is ignored. The same as http.ServeFile
// If the content's Seek method doesn't work: ServeContent uses the param size
// to generate a onlySizeSeekable as a pseudo io.ReadSeeker. If size < 0, use chunk or connection close instead
//
// If the caller has set w's ETag header formatted per RFC 7232, section 2.3,
// ServeContent uses it to handle requests using If-Match, If-None-Match, or If-Range.
//
// Note that *os.File implements the io.ReadSeeker interface.
func ServeContent(w http.ResponseWriter, r *http.Request, name string, modtime time.Time, content io.Reader, size int64) {
	readseeker, seekable := content.(io.ReadSeeker)

	// generate a onlySizeSeekable as a pseudo io.ReadSeeker
	if !seekable {

		rangeReq := r.Header.Get("Range")
		if rangeReq != "" {
			ranges, err := parseRange(rangeReq, size)
			if err != nil {
				if err == errNoOverlap {
					w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", size))
				}
				http.Error(w, err.Error(), http.StatusRequestedRangeNotSatisfiable)
				return
			}
			for _, r := range ranges {
				if r.start != 0 {
					// only Range: bytes=0- is supported for none seekable reader
					http.Error(w, "range is not support", http.StatusRequestedRangeNotSatisfiable)
					return
				}
			}
		}

		// Content-Type must be set here, avoid sniff in http.ServeContent for onlySizeSeekable later
		// If Content-Type isn't set, use the file's extension to find it, but
		// if the Content-Type is unset explicitly, do not sniff the type.
		ctypes, haveType := w.Header()["Content-Type"]
		var ctype string
		if !haveType {
			var err error
			ctype, content, err = ContentType(content, name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if len(ctypes) > 0 {
			ctype = ctypes[0]
		}

		w.Header().Set("Content-Type", ctype)

		if size < 0 {
			// to reject unsupported Range
			w.Header().Del("Content-Length")
			w.Header().Set("Content-Encoding", "chunked")

			// Use HTTP Trunk or connection close later
			defer func() {
				if r.Method != "HEAD" {
					_, _ = io.Copy(w, content)
				}
			}()
		}

		readseeker = newOnlySizeSeekable(content, size)
	}

	if size >= 0 {
		readseeker = io_.LimitReadSeeker(readseeker, size)
	}

	if stater, ok := content.(io_.Stater); ok {
		if fi, err := stater.Stat(); err == nil {
			modtime = fi.ModTime()
		}
	}
	http.ServeContent(w, r, name, modtime, readseeker)

	// Use HTTP Trunk or connection close by defer

	return
}

// can only be used for ServeContent
type onlySizeSeekable struct {
	r      io.Reader
	size   int64
	offset int64
}

func newOnlySizeSeekable(r io.Reader, size int64) *onlySizeSeekable {
	return &onlySizeSeekable{
		r:    r,
		size: size,
	}
}

func (s *onlySizeSeekable) Seek(offset int64, whence int) (int64, error) {
	if offset != 0 {
		return 0, os.ErrInvalid
	}
	if whence == io.SeekStart {
		s.offset = 0
		return s.offset, nil
	}
	if whence == io.SeekEnd {
		s.offset = s.size
		return s.offset, nil
	}
	if whence == io.SeekCurrent {
		return s.offset, nil
	}
	return s.offset, os.ErrInvalid
}

func (s *onlySizeSeekable) Read(p []byte) (n int, err error) {
	n, err = s.r.Read(p)
	s.offset += int64(n)
	return n, err
}
