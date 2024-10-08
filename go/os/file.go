// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	filepath_ "github.com/searKing/golang/go/path/filepath"
)

const (
	DefaultPermissionFile       os.FileMode = 0666
	DefaultPermissionDirectory  os.FileMode = 0755
	DefaultFlagCreateIfNotExist             = os.O_RDWR | os.O_CREATE
	DefaultFlagCreateTruncate               = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	DefaultFlagCreate                       = DefaultFlagCreateTruncate
	DefaultFlagCreateAppend                 = os.O_RDWR | os.O_CREATE | os.O_APPEND
	DefaultFlagLock                         = os.O_RDWR | os.O_CREATE | os.O_EXCL
)

func GetAbsBinDir() (dir string, err error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// RemoveIfExist removes the named file or (empty) directory.
// If the path does not exist, RemoveIfExist returns nil (no error).
// If there is an error, it will be of type *PathError.
func RemoveIfExist(name string) error {
	err := os.Remove(name)
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	return err
}

// ChtimesNow changes the access and modification times of the named
// file with Now, similar to the Unix utime() or utimes() functions.
//
// The underlying filesystem may truncate or round the values to a
// less precise time unit.
// If there is an error, it will be of type *PathError.
func ChtimesNow(name string) error {
	now := time.Now()
	return os.Chtimes(name, now, now)
}

// MakeAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// If the dir does not exist, it is created with mode 0755 (before umask).
func MakeAll(name string) error {
	return os.MkdirAll(name, DefaultPermissionDirectory)
}

// Make creates a directory named path and returns nil,
// or else returns an error.
// If the dir does not exist, it is created with mode 0755 (before umask).
func Make(name string) error {
	return os.Mkdir(name, DefaultPermissionDirectory)
}

// RenameAll renames (moves) oldpath to newpath.
// If newpath already exists and is not a directory, Rename replaces it.
// OS-specific restrictions may apply when oldpath and newpath are in different directories.
// If there is an error, it will be of type *LinkError.
// If the dir does not exist, it is created with mode 0755 (before umask).
func RenameAll(oldpath, newpath string) error {
	return RenameFileAll(oldpath, newpath, DefaultPermissionDirectory)
}

// RenameFileAll is the generalized open call; most users will use RenameAll instead.
// It renames (moves) oldpath to newpath.
func RenameFileAll(oldpath, newpath string, dirperm os.FileMode) error {
	// file or dir not exists
	if _, err := os.Stat(newpath); err != nil {
		dir, _ := filepath.Split(newpath)

		if dir != "" {
			// mkdir -p dir
			if err := os.MkdirAll(dir, dirperm); err != nil {
				return err
			}
		}
	}
	return os.Rename(oldpath, newpath)
}

// TouchAll creates the named file or dir. If the file already exists, it is touched to now.
// If the file does not exist, it is created with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0755 (before umask).
func TouchAll(path string) (*os.File, error) {
	f, err := OpenFileAll(path, os.O_WRONLY|os.O_CREATE, DefaultPermissionDirectory, DefaultPermissionFile)
	if err != nil {
		return nil, err
	}
	if err := ChtimesNow(path); err != nil {
		defer f.Close()
		return nil, err
	}
	return f, nil
}

// CreateAll creates or truncates the named file or dir. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0755 (before umask).
func CreateAll(path string) (*os.File, error) {
	return OpenFileAll(path, DefaultFlagCreate, DefaultPermissionDirectory, DefaultPermissionFile)
}

// CreateAllIfNotExist creates the named file or dir. If the file does not exist, it is created
// with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0755 (before umask).
// If path is already a directory, CreateAllIfNotExist does nothing and returns nil.
func CreateAllIfNotExist(path string) (*os.File, error) {
	return OpenFileAll(path, DefaultFlagCreateIfNotExist, DefaultPermissionDirectory, DefaultPermissionFile)
}

// AppendAllIfNotExist appends the named file or dir. If the file does not exist, it is created
// with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0755 (before umask).
// If path is already a directory, CreateAllIfNotExist does nothing and returns nil.
func AppendAllIfNotExist(path string) (*os.File, error) {
	return OpenFileAll(path, DefaultFlagCreateAppend, DefaultPermissionDirectory, DefaultPermissionFile)
}

// OpenAll opens the named file or dir for reading. If successful, methods on
// the returned file or dir can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func OpenAll(path string) (*os.File, error) {
	return OpenFileAll(path, os.O_RDONLY, 0, 0)
}

// LockAll creates the named file or dir. If the file already exists, error returned.
// If the file does not exist, it is created with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0755 (before umask).
func LockAll(path string) (*os.File, error) {
	return OpenFileAll(path, DefaultFlagLock, DefaultPermissionDirectory, DefaultPermissionFile)
}

// OpenFileAll is the generalized open call; most users will use OpenAll
// or CreateAll instead. It opens the named file or directory with specified flag
// (O_RDONLY etc.).
// If the file does not exist, and the O_CREATE flag is passed, it is created with mode fileperm (before umask).
// If the directory does not exist, it is created with mode dirperm (before umask).
// If successful, methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.
func OpenFileAll(path string, flag int, dirperm, fileperm os.FileMode) (*os.File, error) {
	dir, file := filepath.Split(path)
	// file or dir exists
	if _, err := os.Stat(path); err == nil {
		return os.OpenFile(path, flag, 0)
	}

	if dir != "" {
		// mkdir -p dir
		if err := os.MkdirAll(dir, dirperm); err != nil {
			return nil, err
		}
	}

	// create file if needed
	if file == "" {
		return nil, nil
	}
	return os.OpenFile(path, flag, fileperm)
}

// CopyAll creates or truncates the dst file or dir, filled with content from src file.
// If the dst file already exists, it is truncated.
// If the dst file does not exist, it is created with mode 0666 (before umask).
// If the dst dir does not exist, it is created with mode 0755 (before umask).
func CopyAll(dst string, src string) error {
	return CopyFileAll(dst, src, DefaultFlagCreate, DefaultPermissionDirectory, DefaultPermissionFile)
}

// CopyAppendAll creates or appends the dst file or dir, filled with content from src file.
// If the dst file already exists, it is truncated.
// If the dst file does not exist, it is created with mode 0666 (before umask).
// If the dst dir does not exist, it is created with mode 0755 (before umask).
func CopyAppendAll(dst string, src string) error {
	return CopyFileAll(dst, src, DefaultFlagCreateAppend, DefaultPermissionDirectory, DefaultPermissionFile)
}

// CopyFileAll is the generalized open call; most users will use CopyAll
// or AppendAll instead. It opens the named file or directory with specified flag
// (O_RDONLY etc.).
// If the dst file does not exist, and the O_CREATE flag is passed, it is created with mode fileperm (before umask).
// If the dst directory does not exist, it is created with mode dirperm (before umask).
// If successful, methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.
func CopyFileAll(dst string, src string, flag int, dirperm, fileperm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := OpenFileAll(dst, flag, dirperm, fileperm)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Copy creates or truncates the dst file or dir, filled with content from src file.
// If the dst file already exists, it is truncated.
// If the dst file does not exist, it is created with mode 0666 (before umask).
// If the dst dir does not exist, it is created with mode 0755 (before umask).
// parent dirs will not be created, otherwise, use CopyAll instead.
func Copy(dst string, src string) error {
	return CopyFile(dst, src, DefaultFlagCreate, DefaultPermissionFile)
}

// Append creates or appends the dst file or dir, filled with content from src file.
// If the dst file already exists, it is truncated.
// If the dst file does not exist, it is created with mode 0666 (before umask).
// If the dst dir does not exist, it is created with mode 0755 (before umask).
// parent dirs will not be created, otherwise, use AppendAll instead.
func Append(dst string, src string) error {
	return CopyFile(dst, src, DefaultFlagCreateAppend, DefaultPermissionFile)
}

// CopyFile is the generalized open call; most users will use Copy
// or Append instead. It opens the named file or directory with specified flag
// (O_RDONLY etc.).
// CopyFile copies from src to dst.
// parent dirs will not be created, otherwise, use CopyFileAll instead.
func CopyFile(dst string, src string, flag int, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, flag, perm)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CopyTruncateAll truncates the original src file in place after creating a copy dst, instead of moving the
// src file to dir and optionally creating a new one src.  It can be used when some program  can‐
// not  be  told  to  close its logfile and thus might continue writing (appending) to the
// previous log file forever.  Note that there is a very small time slice between  copying
// the  file  and  truncating it, so some logging data might be lost.
// parent dirs will be created with dirperm if not exist.
func CopyTruncateAll(dst string, src string) error {
	return CopyTruncateFileAll(dst, src, DefaultFlagCreate, DefaultPermissionDirectory, DefaultPermissionFile, 0)
}

// AppendTruncateAll truncates the original src file in place after appending or creating a copy dst,
// instead of moving the src file to dir and optionally creating a new one src.
// parent dirs will be created with dirperm if not exist.
func AppendTruncateAll(dst string, src string) error {
	return CopyTruncateFileAll(dst, src, DefaultFlagCreateAppend, DefaultPermissionDirectory, DefaultPermissionFile, 0)
}

// CopyTruncateFileAll is the generalized open call; most users will use CopyTruncateAll or
// AppendTruncateAll instead. It opens the named file or directory with specified flag (O_RDONLY etc.).
// CopyTruncateFileAll copies from src to dst and truncates src.
// parent dirs will be created with dirperm if not exist.
func CopyTruncateFileAll(dst string, src string, flag int, dirperm, fileperm os.FileMode, size int64) error {
	if err := CopyFileAll(dst, src, flag, dirperm, fileperm); err != nil {
		return err
	}
	return os.Truncate(src, size)
}

// CopyTruncate truncates the original src file in place after creating a copy dst, instead of moving the
// src file to dir and optionally creating a new one src.  It can be used when some program  can‐
// not  be  told  to  close its logfile and thus might continue writing (appending) to the
// previous log file forever.  Note that there is a very small time slice between  copying
// the  file  and  truncating it, so some logging data might be lost.
func CopyTruncate(dst string, src string) error {
	return CopyTruncateFile(dst, src, DefaultFlagCreate, DefaultPermissionFile, 0)
}

// AppendTruncate truncates the original src file in place after appending or creating a copy dst,
// instead of moving the src file to dir and optionally creating a new one src.
func AppendTruncate(dst string, src string) error {
	return CopyTruncateFile(dst, src, DefaultFlagCreateAppend, DefaultPermissionFile, 0)
}

// CopyTruncateFile is the generalized open call; most users will use CopyTruncate or
// AppendTruncate instead. It opens the named file or directory with specified flag (O_RDONLY etc.).
// CopyTruncateFile copies from src to dst and truncates src.
// parent dirs will not be created, otherwise, use CopyTruncateFileAll instead.
// CopyTruncateFile = CopyFile(src->dst) + Truncate(src)
func CopyTruncateFile(dst string, src string, flag int, perm os.FileMode, size int64) error {
	if err := CopyFile(dst, src, flag, perm); err != nil {
		return err
	}
	return os.Truncate(src, size)
}

// CopyRenameAll makes a copy of the src file, but don't change the original src at all.
// This option can be used, for instance, to make a snapshot of the current log file,
// or when some other utility needs to truncate or parse the file.
// parent dirs will be created with dirperm if not exist.
func CopyRenameAll(dst string, src string) error {
	return CopyRenameFileAll(dst, src, DefaultFlagCreateIfNotExist, DefaultPermissionDirectory, DefaultPermissionFile)
}

// CopyRenameFileAll is the generalized open call; most users will use CopyRenameAll instead.
// It makes a copy of the src file, but don't change the original src at all.
// CopyRenameFileAll renames from src to dst and creates src if not exist.
// parent dirs will be created with dirperm if not exist.
// CopyRenameFileAll = RenameFileAll(src->dst) + OpenFile(src)
func CopyRenameFileAll(dst string, src string, flag int, dirperm, fileperm os.FileMode) error {
	if err := RenameFileAll(src, dst, dirperm); err != nil {
		return err
	}
	f, err := os.OpenFile(src, flag, fileperm)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// CopyRename makes a copy of the src file, but don't change the original src at all.  This option can  be
// used,  for  instance,  to  make  a snapshot of the current log file, or when some other
// utility needs to truncate or parse the file.
func CopyRename(dst string, src string) error {
	return CopyRenameFile(dst, src, DefaultFlagCreate, DefaultPermissionFile)
}

// CopyRenameFile is the generalized open call; most users will use CopyRename instead.
// It opens the named file or directory with specified flag (O_RDONLY etc.).
// CopyTruncateFile copies from src to dst and truncates src.
// parent dirs will not be created, otherwise, use CopyRenameFileAll instead.
func CopyRenameFile(dst string, src string, flag int, perm os.FileMode) error {
	if err := os.Rename(src, dst); err != nil {
		return err
	}
	f, err := os.OpenFile(src, flag, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// SameFile reports whether fi1 and fi2 describe the same file.
// Overload os.SameFile by file path
func SameFile(fi1, fi2 string) bool {
	stat1, err := os.Stat(fi1)
	if err != nil {
		return false
	}

	stat2, err := os.Stat(fi2)
	if err != nil {
		return false
	}
	return os.SameFile(stat1, stat2)
}

// ReLink creates or replaces newname as a hard link to the oldname file.
// If there is an error, it will be of type *LinkError.
func ReLink(oldname, newname string) error {
	tempLink, err := os.CreateTemp(filepath.Dir(newname), filepath.Base(newname))
	if err != nil {
		return err
	}
	if err = tempLink.Close(); err != nil {
		return err
	}

	if err = os.Remove(tempLink.Name()); err != nil {
		return err
	}

	defer os.Remove(tempLink.Name())
	// keep mode the same if newname already exists.
	if fi, err := os.Stat(newname); err == nil {
		if err := os.Chmod(tempLink.Name(), fi.Mode()); err != nil {
			return err
		}
	}

	if err := os.Link(oldname, tempLink.Name()); err != nil {
		return err
	}
	return os.Rename(tempLink.Name(), newname)

}

// ReSymlink creates or replace newname as a symbolic link to oldname.
// If there is an error, it will be of type *LinkError.
func ReSymlink(oldname, newname string) error {
	tempLink, err := os.CreateTemp(filepath.Dir(newname), filepath.Base(newname))
	if err != nil {
		return err
	}
	if err = tempLink.Close(); err != nil {
		return err
	}

	if err = os.Remove(tempLink.Name()); err != nil {
		return err
	}

	defer os.Remove(tempLink.Name())
	if err := os.Symlink(oldname, tempLink.Name()); err != nil {
		return err
	}
	// keep mode the same if newname already exists.
	if fi, err := os.Stat(newname); err == nil {
		if err := os.Chmod(tempLink.Name(), fi.Mode()); err != nil {
			return err
		}
	}

	return os.Rename(tempLink.Name(), newname)
}

// NextFile creates a new file, opens the file for reading and writing,
// and returns the resulting *os.File.
// The filename is generated by taking pattern and adding a seq increased
// string to the end. If pattern includes a "*", the random string
// replaces the last "*".
// Multiple programs calling NextFile simultaneously
// will not choose the same file. The caller can use f.Name()
// to find the pathname of the file. It is the caller's responsibility
// to remove the file when no longer needed.
func NextFile(pattern string, seq int) (f *os.File, seqUsed int, err error) {
	// prefixAndSuffix splits pattern by the last wildcard "*", if applicable,
	// returning prefix as the part before "*" and suffix as the part after "*".
	prefix, suffix := prefixAndSuffix(pattern)

	try := 0
	for {
		seqUsed = seq + try
		name := fmt.Sprintf("%s%d%s", prefix, seqUsed, suffix)
		f, err = LockAll(name)
		if os.IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return nil, 0, &os.PathError{Op: "nextfile", Path: prefix + "*" + suffix, Err: os.ErrExist}
		}
		return f, seqUsed, err
	}
}

// MaxSeq return max seq set by NextFile
// split pattern by the last wildcard "*"
func MaxSeq(pattern string) (prefix string, seq int, suffix string) {
	return MaxSeqFunc(pattern, func(name string) bool { return false })
}

// MaxSeqFunc return hit seq or max seq else set by NextFile
// split pattern by the last wildcard "*"
// All errors that arise visiting files and directories are filtered by fn:
// see the [filepath.WalkFunc] documentation for details.
func MaxSeqFunc(pattern string, handler func(name string) bool) (prefix string, seq int, suffix string) {
	// prefixAndSuffix splits pattern by the last wildcard "*", if applicable,
	// returning prefix as the part before "*" and suffix as the part after "*".
	prefix, suffix = prefixAndSuffix(pattern)

	var hitOrMaxSeq int
	// ignore [os.ErrBadPattern], taken as hitOrMaxSeq is 0
	_ = filepath_.WalkGlob(fmt.Sprintf("%s*%s", prefix, suffix), func(path string) error {
		hit := handler(path)
		// filepath.Clean fix ./xxx -> xxx
		seqStr := strings.TrimSuffix(strings.TrimPrefix(path, filepath.Clean(prefix)), suffix)
		if seq, err := strconv.Atoi(seqStr); err == nil {
			if seq > hitOrMaxSeq || hit {
				hitOrMaxSeq = seq
			}
		} else if hit {
			hitOrMaxSeq = 0
		}

		if hit {
			return filepath.SkipAll
		}
		return nil
	})
	return prefix, hitOrMaxSeq, suffix
}

// prefixAndSuffix splits pattern by the last wildcard "*", if applicable,
// returning prefix as the part before "*" and suffix as the part after "*".
func prefixAndSuffix(pattern string) (prefix, suffix string) {
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}
	return prefix, suffix
}

// WriteAll writes data to a file named by filename.
// If the file does not exist, WriteAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteAll creates it with 0755 (before umask)
// otherwise WriteAll truncates it before writing, without changing permissions.
func WriteAll(filename string, data []byte) error {
	return WriteFileAll(filename, data, DefaultPermissionDirectory, DefaultPermissionFile)
}

// WriteFileAll is the generalized open call; most users will use WriteAll instead.
// It writes data to a file named by filename.
// If the file does not exist, WriteFileAll creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAll creates it with permissions dirperm (before umask)
// otherwise WriteFileAll truncates it before writing, without changing permissions.
func WriteFileAll(filename string, data []byte, dirperm, fileperm os.FileMode) error {
	return WriteFileAllFrom(filename, bytes.NewReader(data), dirperm, fileperm)
}

// WriteAllFrom writes data to a file named by filename from r until EOF or error.
// If the file does not exist, WriteAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteAll creates it with 0755 (before umask)
// otherwise WriteAll truncates it before writing, without changing permissions.
func WriteAllFrom(filename string, r io.Reader) error {
	return WriteFileAllFrom(filename, r, DefaultPermissionDirectory, DefaultPermissionFile)
}

// WriteFileAllFrom is the generalized open call; most users will use WriteAllFrom instead.
// It writes data to a file named by filename from r until EOF or error.
// If the file does not exist, WriteFileAllFrom creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAllFrom creates it with permissions dirperm (before umask)
// otherwise WriteFileAllFrom truncates it before writing, without changing permissions.
func WriteFileAllFrom(filename string, r io.Reader, dirperm, fileperm os.FileMode) error {
	f, err := OpenFileAll(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, dirperm, fileperm)
	if err != nil {
		return err
	}
	_, err = f.ReadFrom(r)
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
	return AppendFileAll(filename, data, DefaultPermissionDirectory, DefaultPermissionFile)
}

// AppendFileAll is the generalized open call; most users will use AppendAll instead.
// It appends data to a file named by filename.
// If the file does not exist, AppendFileAll creates it with permissions fileperm (before umask)
// If the dir does not exist, AppendFileAll creates it with permissions dirperm (before umask)
// otherwise AppendFileAll appends it before writing, without changing permissions.
func AppendFileAll(filename string, data []byte, dirperm, fileperm os.FileMode) error {
	return AppendFileAllFrom(filename, bytes.NewReader(data), dirperm, fileperm)
}

// AppendAllFrom appends data to a file named by filename from r until EOF or error.
// If the file does not exist, AppendAllFrom creates it with mode 0666 (before umask)
// If the dir does not exist, AppendAllFrom creates it with 0755 (before umask)
// (before umask); otherwise AppendAllFrom appends it before writing, without changing permissions.
func AppendAllFrom(filename string, r io.Reader) error {
	return AppendFileAllFrom(filename, r, DefaultPermissionDirectory, DefaultPermissionFile)
}

// AppendFileAllFrom is the generalized open call; most users will use AppendFileFrom instead.
// It appends data to a file named by filename from r until EOF or error.
// If the file does not exist, AppendFileAllFrom creates it with permissions fileperm (before umask)
// If the dir does not exist, AppendFileAllFrom creates it with permissions dirperm (before umask)
// otherwise AppendFileAllFrom appends it before writing, without changing permissions.
func AppendFileAllFrom(filename string, r io.Reader, dirperm, fileperm os.FileMode) error {
	f, err := OpenFileAll(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, dirperm, fileperm)
	if err != nil {
		return err
	}
	_, err = f.ReadFrom(r)
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
	return WriteRenameFileAll(filename, data, DefaultPermissionDirectory)
}

// WriteRenameFileAll is the generalized open call; most users will use WriteRenameAll instead.
// WriteRenameFileAll is safer than WriteFileAll as before Write finished, nobody can find the unfinished file.
// It writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameFileAll creates it with permissions fileperm
// If the dir does not exist, WriteRenameFileAll creates it with permissions dirperm
// (before umask); otherwise WriteRenameFileAll truncates it before writing, without changing permissions.
func WriteRenameFileAll(filename string, data []byte, dirperm os.FileMode) error {
	return WriteRenameFileAllFrom(filename, bytes.NewReader(data), dirperm)
}

// WriteRenameAllFrom writes data to a temp file from r until EOF or error, and rename to the new file named by filename.
// WriteRenameAllFrom is safer than WriteAllFrom as before Write finished, nobody can find the unfinished file.
// If the file does not exist, WriteRenameAllFrom creates it with mode 0666 (before umask)
// If the dir does not exist, WriteRenameAllFrom creates it with 0755 (before umask)
// otherwise WriteRenameAllFrom truncates it before writing, without changing permissions.
func WriteRenameAllFrom(filename string, r io.Reader) error {
	return WriteRenameFileAllFrom(filename, r, DefaultPermissionDirectory)
}

// WriteRenameFileAllFrom is the generalized open call; most users will use WriteRenameAllFrom instead.
// WriteRenameFileAllFrom is safer than WriteRenameAllFrom as before Write finished, nobody can find the unfinished file.
// It writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameFileAllFrom creates it with permissions fileperm
// If the dir does not exist, WriteRenameFileAllFrom creates it with permissions dirperm
// (before umask); otherwise WriteRenameFileAllFrom truncates it before writing, without changing permissions.
func WriteRenameFileAllFrom(filename string, r io.Reader, dirperm os.FileMode) error {
	dir, file := filepath.Split(filename)
	if dir == "" {
		dir = "."
	}
	pattern := ".*.rename"
	if file != "" {
		pattern = fmt.Sprintf(".%s.*.rename", file)
	}
	tempFile, err := TempAll(dir, pattern)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath) // remove if rename failed
	_, err = tempFile.ReadFrom(r)
	if err != nil {
		return err
	}
	return RenameFileAll(tempFilePath, filename, dirperm)
}

// TempAll creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
// If the file does not exist, TempAll creates it with mode 0600 (before umask)
// If the dir does not exist, TempAll creates it with 0755 (before umask)
// otherwise TempAll truncates it before writing, without changing permissions.
func TempAll(dir, pattern string) (f *os.File, err error) {
	return TempFileAll(dir, pattern, DefaultPermissionDirectory)
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
	return os.CreateTemp(dir, pattern)
}

// MkdirTempAll creates a new temporary directory in the directory dir
// and returns the pathname of the new directory.
// The new directory's name is generated by adding a random string to the end of pattern.
// If pattern includes a "*", the random string replaces the last "*" instead.
// If dir is the empty string, MkdirTemp uses the default directory for temporary files, as returned by TempDir.
// Multiple programs or goroutines calling MkdirTemp simultaneously will not choose the same directory.
// It is the caller's responsibility to remove the directory when it is no longer needed.
// If the dir does not exist, TempAll creates it with 0755 (before umask)
func MkdirTempAll(dir, pattern string) (string, error) {
	return MkdirTempDirAll(dir, pattern, DefaultPermissionDirectory)
}

// MkdirTempDirAll is the generalized open call; most users will use MkdirTempAll instead.
// If the directory does not exist, it is created with mode dirperm (before umask).
func MkdirTempDirAll(dir, pattern string, dirperm os.FileMode) (string, error) {
	if dir != "" {
		// mkdir -p dir
		if err := os.MkdirAll(dir, dirperm); err != nil {
			return "", err
		}
	}
	return os.MkdirTemp(dir, pattern)
}
