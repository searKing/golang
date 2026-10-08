// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	filepath_ "github.com/searKing/golang/go/path/filepath"
	time_ "github.com/searKing/golang/go/time"
)

// RotateMode represents a way to rotate file.
// see details: https://man7.org/linux/man-pages/man8/logrotate.8.html
type RotateMode int

const (
	// RotateModeNew create new rotate file directly
	// Immediately after rotation (before the postrotate script is run) the log file is created (with
	// the same name as the log file just rotated).
	// see option: `create` in https://man7.org/linux/man-pages/man8/logrotate.8.html
	RotateModeNew RotateMode = iota

	// RotateModeCopyRename rename the log file to the new rotate file, copy it back as the old rotate file,
	// then truncate the new rotate file. Rotate files are named by the time of their contents as RotateModeNew,
	// while the log file keeps its inode, so that some program which can not be told to close its log file
	// continues writing (appending) to the new rotate file.
	// Note that there is a very small time slice between copying the file and truncating it,
	// so some logging data written by other programs might be lost.
	RotateModeCopyRename RotateMode = iota

	// RotateModeCopyTruncate Truncate the original log file in place after creating a copy, instead of moving the
	// old log file and optionally creating a new one. It can be used when some program can‐
	// not be told to close its rotatefile and thus might continue writing (appending) to the
	// previous log file forever. Note that there is a very small time slice between copying
	// the file and truncating it, so some logging data might be lost. When this option is
	// used, the create option will have no effect, as the old log file stays in place.
	//
	// The log file is named by RotateFile.CopyTruncateFilePath if set, or keeps the name of the first rotate file
	// otherwise, and the copies are named by the time of rotation, that is, each copy contains logs of the
	// previous rotate interval, as `dateext` without `dateyesterday` in logrotate.
	// Use RotateModeCopyRename if the log file need not keep its name.
	// see option: `copytruncate` in https://man7.org/linux/man-pages/man8/logrotate.8.html
	RotateModeCopyTruncate RotateMode = iota
)

// RotateFile logrotate reads everything about the log files it should be handling from the series of con‐
// figuration files specified on the command line.  Each configuration file can set global
// options (local definitions override global ones, and later definitions override earlier ones)
// and specify rotatefiles to rotate. A simple configuration file looks like this:
type RotateFile struct {
	RotateMode           RotateMode
	FilePathPrefix       string // FilePath = FilePathPrefix + now.Format(filePathRotateLayout)
	FilePathRotateLayout string // Time layout to format rotate file

	RotateFileGlob string // file glob to clean

	// sets the symbolic link name that gets linked to the current file name being used.
	FileLinkPath string

	// CopyTruncateFilePath sets the fixed name of the log file being written for RotateModeCopyTruncate,
	// which is never removed by cleaning. If empty, the log file keeps the name of the first rotate file.
	// The log file is copied on startup only if ForceNewFileOnStartup is set.
	// It takes no effect for other RotateMode.
	CopyTruncateFilePath string

	// Rotate files are rotated until RotateInterval expired before being removed
	// take effects if only RotateInterval is bigger than 0.
	RotateInterval time.Duration

	// Rotate files are rotated if they grow bigger then size bytes.
	// take effects if only RotateSize is bigger than 0.
	RotateSize int64

	// max age of a log file before it gets purged from the file system.
	// Remove rotated logs older than duration. The age is only checked if the file is
	// to be rotated.
	// take effects if only MaxAge is bigger than 0.
	MaxAge time.Duration

	// Rotate files are rotated MaxCount times before being removed
	// take effects if only MaxCount is bigger than 0.
	MaxCount int

	// Force File Rotate when start up
	ForceNewFileOnStartup bool

	// PreRotateHandler called before file rotate
	// name means file path rotated
	PreRotateHandler func(name string)

	// PostRotateHandler called after file rotate
	// name means file path rotated
	PostRotateHandler func(name string)

	// CleanErrorHandler called with the error of cleaning rotate files if any, such as failure to remove.
	// It is called without RotateFile locked, so it can write to the RotateFile, such as logging.
	// It is called in background for cleaning on rotate, maybe after Close returns.
	CleanErrorHandler func(err error)

	cleaning               atomic.Bool
	cleanWg                sync.WaitGroup // cleaning on rotate, waited by Close
	mu                     sync.Mutex
	writingSeq             int // seq of writingFilePathRotated, file rotated by size limit meet
	writingFilePath        string
	writingFile            *os.File
	writingFilePathRotated string // rotated file path, same as writingFilePath except for copytruncate
}

func NewRotateFile(layout string) *RotateFile {
	return NewRotateFileWithStrftime(time_.LayoutTimeToSimilarStrftime(layout))
}

func NewRotateFileWithStrftime(strftimeLayout string) *RotateFile {
	return &RotateFile{
		FilePathRotateLayout: time_.LayoutStrftimeToSimilarTime(strftimeLayout),
		RotateFileGlob:       fileGlobFromStrftimeLayout(strftimeLayout),
		RotateInterval:       24 * time.Hour,
	}
}

func fileGlobFromStrftimeLayout(strftimeLayout string) string {
	var regexps = []*regexp.Regexp{
		regexp.MustCompile(`%[%+A-Za-z]`),
		regexp.MustCompile(`\*+`),
	}
	globPattern := strftimeLayout
	for _, re := range regexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}
	return globPattern + `*`
}

func (f *RotateFile) Write(b []byte) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}
	// Guard against concurrent writes
	f.mu.Lock()
	defer f.mu.Unlock()

	out, err := f.getWriterLocked(false, false)
	if err != nil {
		return 0, fmt.Errorf("acquire rotated file: %w", err)
	}

	return out.Write(b)
}

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (f *RotateFile) WriteString(s string) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}
	return f.Write([]byte(s))
}

// WriteAt writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
//
// If file was opened with the O_APPEND flag, WriteAt returns an error.
func (f *RotateFile) WriteAt(b []byte, off int64) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}
	// Guard against concurrent writes
	f.mu.Lock()
	defer f.mu.Unlock()

	out, err := f.getWriterLocked(false, false)
	if err != nil {
		return 0, fmt.Errorf("acquire rotated file: %w", err)
	}
	if w, ok := out.(io.WriterAt); ok {
		return w.WriteAt(b, off)
	}
	return 0, fmt.Errorf("WriteAt is not supported")
}

// Close satisfies the io.Closer interface. You must
// call this method if you performed any writes to
// the object.
// Close waits for cleaning of rotate files in background.
func (f *RotateFile) Close() error {
	if err := f.checkValid("close"); err != nil {
		return err
	}
	err, cleanErr := f.close()
	f.handleCleanError(cleanErr)
	return err
}

// close closes the writing file, then cleans rotate files, with f locked.
func (f *RotateFile) close() (closeErr, cleanErr error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	// files are not touched by cleaning on rotate after Close returns
	f.cleanWg.Wait()

	if f.writingFile == nil { // maybe file is closed or not open
		return nil, nil
	}
	closeErr = f.writingFile.Close()
	f.writingFile = nil
	return closeErr, f.serializedClean("")
}

// Rotate forcefully rotates the file. If the generated file name
// clash because file already exists, a numeric suffix of the form
// ".1", ".2", ".3" and so forth are appended to the end of the log file
//
// This method can be used in conjunction with a signal handler so to
// emulate servers that generate new log files when they receive a SIGHUP
func (f *RotateFile) Rotate(forceRotate bool) error {
	if err := f.checkValid("rotate"); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, err := f.getWriterLocked(true, forceRotate); err != nil {
		return err
	}
	return nil
}

func (f *RotateFile) filePathByRotateTime() string {
	// create a new file name using the regular time layout
	return f.FilePathPrefix + time_.TruncateByLocation(time.Now(), f.RotateInterval).Format(f.FilePathRotateLayout)
}

func (f *RotateFile) filePathByRotateSize() (name string, seq int) {
	// instead of just using the regular time layout,
	// we create a new file name using names such as "foo.1", "foo.2", "foo.3", etc
	return nextSeqFileName(f.filePathByRotateTime(), f.writingSeq)
}

func (f *RotateFile) filePathByRotate(forceRotate bool) (name string, seq int, byTime, bySize bool) {
	// name using the regular time layout, without seq
	name = f.filePathByRotateTime()
	// startup
	if f.writingFilePath == "" {
		if f.ForceNewFileOnStartup {
			// instead of just using the regular time layout,
			// we create a new file name using names such as "foo", "foo.1", "foo.2", "foo.3", etc
			name, seq = nextSeqFileName(name, f.writingSeq)
			return name, seq, false, true
		}
		name, seq = maxSeqFileName(name)
		return name, seq, true, false
	}

	// rotate by time
	// compare expect time with current rotated file, as writingFilePath stays the same for copytruncate
	if name != trimSeqFromNextFileName(f.writingFilePathRotated, f.writingSeq) {
		if forceRotate {
			// instead of just using the regular time layout,
			// we create a new file name using names such as "foo", "foo.1", "foo.2", "foo.3", etc
			name, seq = nextSeqFileName(name, 0)
			return name, seq, true, false
		}
		name, seq = maxSeqFileName(name)
		return name, seq, true, false
	}

	// determine if rotate by size

	// recreate file if file not exist, removed by other process for example
	usingFileInfo, err := os.Stat(f.writingFilePath)
	if os.IsNotExist(err) {
		name = f.writingFilePath
		seq = f.writingSeq
		return name, seq, false, false
	}

	// rotate by size
	// compare rotate size with current using file
	if forceRotate || (err == nil && (f.RotateSize > 0 && usingFileInfo.Size() > f.RotateSize)) {
		// instead of just using the regular time layout,
		// we create a new file name using names such as "foo", "foo.1", "foo.2", "foo.3", etc
		name, seq = nextSeqFileName(name, f.writingSeq)
		return name, seq, false, true
	}
	name = f.writingFilePath
	seq = f.writingSeq
	return name, seq, false, false
}

func (f *RotateFile) makeUsingFileReadyLocked() (err error) {
	// using file exist, close this file if not ready to use
	if f.writingFile != nil {
		diskFileInfo, err := os.Stat(f.writingFile.Name())
		if err == nil {
			usingFileInfo, statErr := f.writingFile.Stat()
			if statErr == nil {
				if os.SameFile(diskFileInfo, usingFileInfo) {
					return nil
				}
			}
		}
		// file not exist or not the same file, recreate the file and file link
		_ = f.writingFile.Close()
		f.writingFile = nil
	}

	// using file not exist, recreate the file and file link
	file, err := AppendAllIfNotExist(f.writingFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = file.Close()
		}
	}()

	if err := f.symlinkLocked(f.writingFilePath); err != nil {
		return err
	}
	f.writingFile = file
	return nil

}

func (f *RotateFile) getWriterLocked(bailOnRotateFail, forceRotate bool) (out io.Writer, err error) {
	newName, newSeq, byTime, bySize := f.filePathByRotate(forceRotate)
	needRotate := byTime || bySize
	if !needRotate { // no need to rotate by time and size
		err = f.makeUsingFileReadyLocked()
		if err != nil {
			return nil, err
		}
		return f.writingFile, nil
	}

	// rotate the file
	if f.PreRotateHandler != nil {
		f.PreRotateHandler(f.writingFilePath)
	}
	newFile, err := f.rotateLocked(newName)
	if err != nil {
		if bailOnRotateFail {
			// Failure to rotate is a problem, but it's really not a great idea
			// to stop your application just because you couldn't rename your log.
			//
			// We only return this error when explicitly needed (as specified by bailOnRotateFail)
			return nil, err
		}
		// no file can be written, it's an error explicitly
		if f.writingFile == nil {
			return nil, err
		}
		return f.writingFile, nil
	}

	// swap file
	if f.writingFile != nil {
		_ = f.writingFile.Close()
		f.writingFile = nil
	}
	f.writingFilePathRotated = newName
	f.writingFile = newFile
	f.writingFilePath = newFile.Name()
	f.writingSeq = newSeq
	if f.PostRotateHandler != nil {
		f.PostRotateHandler(f.writingFilePath)
	}

	return f.writingFile, nil
}

// file may not be nil if err is nil
func (f *RotateFile) rotateLocked(newName string) (_ *os.File, err error) {
	// if we got here, then we need to create a file
	oldName := f.writingFilePath
	var writeName string
	var needRotate bool
	switch {
	case f.RotateMode == RotateModeCopyTruncate && f.CopyTruncateFilePath != "":
		// newName may be an existing copy on startup, which must not be overwritten unless forced
		needRotate = newName != f.writingFilePathRotated && (oldName != "" || f.ForceNewFileOnStartup)
		oldName = f.CopyTruncateFilePath
		writeName = oldName
	case f.RotateMode == RotateModeCopyTruncate && oldName != "":
		writeName = oldName
		needRotate = newName != f.writingFilePathRotated
	default:
		writeName = newName
		needRotate = newName != f.writingFilePath
	}
	if needRotate {
		var err error
		var mode string
		switch f.RotateMode {
		case RotateModeCopyRename:
			mode = "copy_rename"
			_, err = os.Stat(oldName)
			if err == nil {
				err = f.copyRenameTruncateLocked(newName, oldName)
			} else if os.IsNotExist(err) {
				err = nil
			}
		case RotateModeCopyTruncate:
			mode = "copy_truncate"
			_, err = os.Stat(oldName)
			if err == nil {
				err = CopyTruncateAll(newName, oldName)
			} else if os.IsNotExist(err) {
				err = nil
				// newName may be created empty to reserve a seq, remove it as nothing is copied
				if newName != writeName {
					if fi, statErr := os.Stat(newName); statErr == nil && fi.Size() == 0 {
						_ = os.Remove(newName)
					}
				}
			}
		case RotateModeNew:
			// for which open the file, and write file by RotateFile
			mode = "new"
		}
		if err != nil {
			return nil, fmt.Errorf("failed to %s file %s to %s: %w", mode, oldName, newName, err)
		}
	}

	// make sure new file exists
	file, err := AppendAllIfNotExist(writeName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = file.Close()
		}
	}()

	if err := f.symlinkLocked(writeName); err != nil {
		return nil, err
	}
	// unlink files on a separate goroutine
	f.cleanWg.Add(1)
	go func() {
		err := f.serializedClean(writeName)
		f.cleanWg.Done()
		// not waited by Close, as it may write to f
		f.handleCleanError(err)
	}()

	return file, nil
}

// unlink files
// expect run on a separate goroutine
func (f *RotateFile) serializedClean(protectedPath string) error {
	// running already, ignore duplicate clean
	if !f.cleaning.CompareAndSwap(false, true) {
		return nil
	}
	defer f.cleaning.Store(false)

	now := time.Now()

	// find old files
	var filesNotExpired []rotateFile
	filesExpired, err := filepath_.GlobFunc(f.FilePathPrefix+f.RotateFileGlob, func(name string) bool {
		if protectedPath != "" && name == protectedPath {
			return false
		}
		if f.RotateMode == RotateModeCopyTruncate && f.CopyTruncateFilePath != "" &&
			filepath.Clean(name) == filepath.Clean(f.CopyTruncateFilePath) {
			return false
		}

		// skip symbolic links, such as FileLinkPath, neither counted nor removed
		fi, err := os.Lstat(name)
		if err != nil || fi.Mode()&os.ModeSymlink != 0 {
			return false
		}

		if f.MaxAge <= 0 || now.Sub(fi.ModTime()) < f.MaxAge {
			filesNotExpired = append(filesNotExpired, rotateFile{name: name, modTime: fi.ModTime()})
			return false
		}
		return true
	})
	if err != nil {
		return err
	}

	var filesExceedMaxCount []rotateFile
	if f.MaxCount > 0 && len(filesNotExpired) > f.MaxCount {
		sortRotateFiles(filesNotExpired)
		filesExceedMaxCount = filesNotExpired[:len(filesNotExpired)-f.MaxCount]
	}
	var errs []error
	for _, path := range filesExpired {
		err = os.Remove(path)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, file := range filesExceedMaxCount {
		err = os.Remove(file.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// copyRenameTruncateLocked is CopyRenameTruncateAll of the writing file src.
// A file opened can not be renamed on windows, so the writing file is closed before,
// and reopened if failed, to keep writing to src as other platforms.
func (f *RotateFile) copyRenameTruncateLocked(dst, src string) error {
	if runtime.GOOS != "windows" || f.writingFile == nil {
		return CopyRenameTruncateAll(dst, src)
	}
	_ = f.writingFile.Close()
	f.writingFile = nil
	err := CopyRenameTruncateAll(dst, src)
	if err != nil {
		_ = f.makeUsingFileReadyLocked()
	}
	return err
}

func (f *RotateFile) handleCleanError(err error) {
	if err != nil && f.CleanErrorHandler != nil {
		f.CleanErrorHandler(err)
	}
}

// checkValid checks whether f is valid for use.
// If not, it returns an appropriate error, perhaps incorporating the operation name op.
func (f *RotateFile) checkValid(op string) error {
	if f == nil {
		return os.ErrInvalid
	}
	return nil
}

// foo.txt, 0 -> foo.txt
// foo.txt, 1 -> foo.txt.[1,2,...], which is not exist and seq is max
func nextSeqFileName(name string, seq int) (string, int) {
	// Special case: if seq is 0, we don't want to append 0 to the file name.
	if seq == 0 {
		nf, err := LockAll(name)
		if err == nil {
			_ = nf.Close()
		}
		if !os.IsExist(err) {
			return name, seq
		}
		seq++
	}

	// A new file has been requested. Instead of just using the
	// regular strftime pattern, we create a new file name using
	// generational names such as "foo.1", "foo.2", "foo.3", etc
	nf, seqUsed, err := NextFile(name+".*", seq)
	if err != nil {
		return name, seq
	}
	_ = nf.Close()
	return nf.Name(), seqUsed
}

// symlinkLocked links FileLinkPath to name if FileLinkPath is set,
// creating the dir of FileLinkPath if not exist.
func (f *RotateFile) symlinkLocked(name string) error {
	if f.FileLinkPath == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(f.FileLinkPath), DefaultPermissionDirectory); err != nil {
		return err
	}
	return ReSymlink(symlinkTarget(name, f.FileLinkPath), f.FileLinkPath)
}

// symlinkTarget returns the target of a symbolic link at linkPath to name.
// A relative target is resolved against the dir of the link, not the working dir,
// so a relative name is converted to be relative to the dir of linkPath.
func symlinkTarget(name, linkPath string) string {
	if filepath.IsAbs(name) {
		return name
	}
	absName, err := filepath.Abs(name)
	if err != nil {
		return name
	}
	absLinkDir, err := filepath.Abs(filepath.Dir(linkPath))
	if err != nil {
		return absName
	}
	rel, err := filepath.Rel(absLinkDir, absName)
	if err != nil {
		return absName
	}
	return rel
}

// foo.txt -> foo.txt
// foo.txt.1 -> foo.txt
// foo.txt.1.1 -> foo.txt.1
func trimSeqFromNextFileName(name string, seq int) string {
	if seq == 0 {
		return name
	}
	return strings.TrimSuffix(name, fmt.Sprintf(".%d", seq))
}

// foo.txt.* -> foo.txt.[1,2,...], which exists and seq is max
func maxSeqFileName(name string) (string, int) {
	prefix, seq, suffix := MaxSeq(name + ".*")
	if seq == 0 {
		return name, seq
	}
	return fmt.Sprintf("%s%d%s", prefix, seq, suffix), seq
}

// rotateFile is a rotate file with ModTime, which is stat once for sorting.
type rotateFile struct {
	name    string
	modTime time.Time
}

// sortRotateFiles sorts files by ModTime and name in increase order.
func sortRotateFiles(files []rotateFile) {
	slices.SortFunc(files, func(a, b rotateFile) int {
		if c := a.modTime.Compare(b.modTime); c != 0 {
			return c
		}
		return compareRotateFileName(a.name, b.name)
	})
}

// compareRotateFileName compares names as foo.1, foo.2, ..., foo.9, foo.10, ..., foo,
// by ascii, except that runs of digits are compared by numeric value,
// and a name is after the longer names it prefixes.
func compareRotateFileName(a, b string) int {
	isDigit := func(c byte) bool { return '0' <= c && c <= '9' }
	for a != "" && b != "" {
		if !isDigit(a[0]) || !isDigit(b[0]) {
			if a[0] != b[0] {
				return cmp.Compare(a[0], b[0])
			}
			a, b = a[1:], b[1:]
			continue
		}
		i, j := 0, 0
		for i < len(a) && isDigit(a[i]) {
			i++
		}
		for j < len(b) && isDigit(b[j]) {
			j++
		}
		na, nb := strings.TrimLeft(a[:i], "0"), strings.TrimLeft(b[:j], "0")
		if c := cmp.Compare(len(na), len(nb)); c != 0 {
			return c
		}
		if c := strings.Compare(na, nb); c != 0 {
			return c
		}
		// the same value, such as 01 and 1, the longer first as a prefixed name
		if i != j {
			return cmp.Compare(j, i)
		}
		a, b = a[i:], b[j:]
	}
	return cmp.Compare(len(b), len(a))
}
