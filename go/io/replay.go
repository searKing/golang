package io

import "io"

// ReadReplayer is the interface that groups the basic Read and Replay methods.
type ReadReplayer interface {
	io.Reader
	Replay() ReadReplayer
}

// ReplayReader returns a Reader that allows replay and read from
// the provided input reader.
// data is buffered always.
// buffered data is taken first, if Replay() is called.
func ReplayReader(r io.Reader) ReadReplayer {
	sr := &replayReader{
		ReadSniffer: SniffReader(r).Sniff(true),
	}
	return sr
}

type replayReader struct {
	ReadSniffer
}

func (r *replayReader) Replay() ReadReplayer {
	r.Sniff(false).Sniff(true)
	return r
}
