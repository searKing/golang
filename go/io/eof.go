package io

import "io"

type eofReader struct{}

func (eofReader) Read([]byte) (int, error) {
	return 0, io.EOF
}

// EOFReader returns a Reader that return EOF anytime.
func EOFReader() io.Reader {
	return eofReader{}
}
