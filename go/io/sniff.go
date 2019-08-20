package io

import (
	"bytes"
	"io"
)

// ReadCloser is the interface that groups the basic Read and Close methods.
type ReadSniffer interface {
	io.Reader
	Sniff(sniffing bool)
}

type sinffReader struct {
	source io.Reader
	buffer bytes.Buffer

	selectorF DynamicReaderFunc

	sniffing bool
}

func (sr *sinffReader) Sniff(sniffing bool) {
	if sr.sniffing == sniffing {
		return
	}
	sr.sniffing = sniffing
	if sniffing {
		// We don't need the buffer anymore.
		// Reset it to release the internal slice.
		sr.buffer = bytes.Buffer{}
		sr.selectorF = func() io.Reader {
			return io.TeeReader(sr.source, &sr.buffer)
		}
		return
	}
	sr.selectorF = func() io.Reader {
		return io.MultiReader(&sr.buffer, sr.source)
	}
}

func (sr *sinffReader) Read(p []byte) (n int, err error) {
	return sr.selectorF.Read(p)
}

// SniffReader returns a Reader that allows sniff and read from
// the provided input reader.
// data is buffered if Sniff(true) is called.
// buffered data is taken first, if Sniff(false) is called.
func SniffReader(r io.Reader) ReadSniffer {
	sr := &sinffReader{}
	sr.source = WatchReader(r, WatcherFunc(func(p []byte, n int, err error) (int, error) {
		if err == io.EOF {
			sr.buffer = bytes.Buffer{}
		}
		return n, err
	}))
	return sr
}
