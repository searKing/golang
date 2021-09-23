package ioutil

import (
	"io/ioutil"
	"os"
	"path/filepath"

	os_ "github.com/searKing/golang/go/os"
)

// WriteAll writes data to a file named by filename.
// If the file does not exist, WriteAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteAll creates it with 0755 (before umask)
// otherwise WriteAll truncates it before writing, without changing permissions.
func WriteAll(filename string, data []byte) error {
	return WriteFileAll(filename, data, 0755, 0666)
}

// WriteFileAll is the generalized open call; most users will use WriteAll instead.
// It writes data to a file named by filename.
// If the file does not exist, WriteFileAll creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAll creates it with permissions dirperm (before umask)
// otherwise WriteFileAll truncates it before writing, without changing permissions.
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
// If the file does not exist, WriteFileAll creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAll creates it with permissions dirperm (before umask)
// otherwise WriteFileAll appends it before writing, without changing permissions.
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

// WriteRenameAll writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteRenameAll creates it with 0755 (before umask)
// otherwise WriteRenameAll truncates it before writing, without changing permissions.
func WriteRenameAll(filename string, data []byte) error {
	return WriteRenameFileAll(filename, data, 0755)
}

// WriteRenameFileAll is the generalized open call; most users will use WriteRenameAll instead.
// WriteRenameFileAll is safer than WriteFileAll as before Write finished, nobody can find the unfinished file.
// It writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameFileAll creates it with permissions fileperm
// If the dir does not exist, WriteRenameFileAll creates it with permissions dirperm
// (before umask); otherwise WriteRenameFileAll truncates it before writing, without changing permissions.
func WriteRenameFileAll(filename string, data []byte, dirperm os.FileMode) error {
	tempDir := filepath.Dir(filename)
	if tempDir != "" {
		// mkdir -p dir
		if err := os.MkdirAll(tempDir, dirperm); err != nil {
			return err
		}
	}

	tempFile, err := ioutil.TempFile(tempDir, "")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	tempFilePath := filepath.Join(tempDir, tempFile.Name())
	defer os.Remove(tempFilePath)
	_, err = tempFile.Write(data)
	if err != nil {
		return err
	}
	return os_.RenameFileAll(tempFilePath, filename, dirperm)
}

// TempAll creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
// If the file does not exist, TempAll creates it with mode 0600 (before umask)
// If the dir does not exist, TempAll creates it with 0755 (before umask)
// otherwise TempAll truncates it before writing, without changing permissions.
func TempAll(dir, pattern string) (f *os.File, err error) {
	return TempFileAll(dir, pattern, 0755)
}

// TempFileAll is the generalized open call; most users will use TempAll instead.
// If the directory does not exist, it is created with mode dirperm (before umask).
func TempFileAll(dir, pattern string, dirperm os.FileMode) (f *os.File, err error) {
	if dir != "" {
		// mkdir -p dir
		if err := os.MkdirAll(dir, dirperm); err != nil {
			return nil, err
		}
	}
	return ioutil.TempFile(dir, pattern)
}
