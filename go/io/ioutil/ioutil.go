package ioutil

import (
	"os"

	os_ "github.com/searKing/golang/go/os"
)

// WriteAll writes data to a file named by filename.
// If the file does not exist, WriteAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteAll creates it with 0755 (before umask)
// (before umask); otherwise WriteAll truncates it before writing, without changing permissions.
func WriteAll(filename string, data []byte) error {
	return WriteFileAll(filename, data, 0755, 0666)
}

// WriteFileAll is the generalized open call; most users will use WriteAll instead.
// It writes data to a file named by filename.
// If the file does not exist, WriteFileAll creates it with permissions fileperm
// If the dir does not exist, WriteFileAll creates it with permissions dirperm
// (before umask); otherwise WriteFileAll truncates it before writing, without changing permissions.
func WriteFileAll(filename string, data []byte, dirperm, fileperm os.FileMode) error {
	f, err := os_.OpenFileAll(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, dirperm, fileperm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// AppendAll appends data to a file named by filename.
// If the file does not exist, AppendAll creates it with mode 0666 (before umask)
// If the dir does not exist, AppendAll creates it with 0755 (before umask)
// (before umask); otherwise AppendAll appends it before writing, without changing permissions.
func AppendAll(filename string, data []byte) error {
	return AppendFileAll(filename, data, 0755, 0666)
}

// AppendFileAll appends data to a file named by filename.
// If the file does not exist, WriteFileAll creates it with permissions fileperm
// If the dir does not exist, WriteFileAll creates it with permissions dirperm
// (before umask); otherwise WriteFileAll appends it before writing, without changing permissions.
func AppendFileAll(filename string, data []byte, dirperm, fileperm os.FileMode) error {
	f, err := os_.OpenFileAll(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, dirperm, fileperm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
