package io

type WriterFunc func(p []byte) (n int, err error)

func (f WriterFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

type WriterFuncPrintfLike func(format string, args ...interface{})

func (f WriterFuncPrintfLike) Write(p []byte) (n int, err error) {
	f("%s", string(p))
	return len(p), nil
}
