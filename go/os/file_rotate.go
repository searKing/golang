// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	filepath_ "github.com/searKing/golang/go/path/filepath"
	time_ "github.com/searKing/golang/go/time"
)

type RotateMode int

const (
	// create new rotate file directly
	RotateModeNew RotateMode = iota

	// Make a copy of the log file, but don't change the original at all. This option can be
	// used, for instance, to make a snapshot of the current log file, or when some other
	// utility needs to truncate or parse the file. When this option is used, the create
	// option will have no effect, as the old log file stays in place.
	RotateModeCopyRename RotateMode = iota

	// Truncate the original log file in place after creating a copy, instead of moving the
	// old log file and optionally creating a new one. It can be used when some program can‐
	// not be told to close its rotatefile and thus might continue writing (appending) to the
	// previous log file forever. Note that there is a very small time slice between copying
	// the file and truncating it, so some logging data might be lost. When this option is
	// used, the create option will have no effect, as the old log file stays in place.
	RotateModeCopyTruncate RotateMode = iota
)

// logrotate reads everything about the log files it should be handling from the series of con‐
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

	mu            sync.Mutex
	usingSeq      int // file rotated by size limit meet
	usingFilePath string
	usingFile     *os.File
}

func NewRotateFile(layout string) *RotateFile {
	return &RotateFile{
		FilePathRotateLayout: layout,
		RotateFileGlob:       fileGlobFromStrftimeLayout(time_.LayoutTimeToSimilarStrftime(layout)),
		RotateInterval:       24 * time.Hour,
	}
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
	// Guard against concurrent writes
	f.mu.Lock()
	defer f.mu.Unlock()

	out, err := f.getWriterLocked(false, false)
	if err != nil {
		return 0, fmt.Errorf("acquite rotated file :%w", err)
	}

	return out.Write(b)
}

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (f *RotateFile) WriteString(s string) (n int, err error) {
	return f.Write([]byte(s))
}

// WriteAt writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
//
// If file was opened with the O_APPEND flag, WriteAt returns an error.
func (f *RotateFile) WriteAt(b []byte, off int64) (n int, err error) {
	// Guard against concurrent writes
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.WriteAt(b, off)
}

// Close satisfies the io.Closer interface. You must
// call this method if you performed any writes to
// the object.
func (f *RotateFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.usingFile == nil {
		return nil
	}

	defer func() { f.usingFile = nil }()
	return f.usingFile.Close()
}

// Rotate forcefully rotates the file. If the generated file name
// clash because file already exists, a numeric suffix of the form
// ".1", ".2", ".3" and so forth are appended to the end of the log file
//
// This method can be used in conjunction with a signal handler so to
// emulate servers that generate new log files when they receive a SIGHUP
func (f *RotateFile) Rotate(forceRotate bool) error {
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
	return nextSeqFileName(f.filePathByRotateTime(), f.usingSeq)
}

func (f *RotateFile) filePathByRotate(forceRotate bool) (name string, byTime, bySize bool) {
	name = f.filePathByRotateTime()
	fi, err := os.Stat(name)

	// startup
	if f.usingFilePath == "" {
		if f.ForceNewFileOnStartup {
			// instead of just using the regular time layout,
			// we create a new file name using names such as "foo.1", "foo.2", "foo.3", etc
			name, f.usingSeq = nextSeqFileName(name, f.usingSeq+1)
			return name, false, true
		}
		name, f.usingSeq = maxSeqFileName(name)
		return name, true, false
	}

	// rotate by time
	if name != trimSeqFromNextFileName(f.usingFilePath, f.usingSeq) {
		if os.IsNotExist(err) {
			// rotate by time, reset part seq of rotate by size
			name_, seq := maxSeqFileName(name)
			f.usingSeq = seq
			return name_, true, false
		}
		// instead of just using the regular time layout,
		// we create a new file name using names such as "foo.1", "foo.2", "foo.3", etc
		name, f.usingSeq = nextSeqFileName(name, f.usingSeq+1)
		return name, false, true
	}
	// rotate by size
	if forceRotate || (err == nil && (f.RotateSize > 0 && fi.Size() > f.RotateSize)) {
		// instead of just using the regular time layout,
		// we create a new file name using names such as "foo.1", "foo.2", "foo.3", etc
		name, f.usingSeq = nextSeqFileName(name, f.usingSeq+1)
		return name, false, true
	}
	return name, false, false
}

func (f *RotateFile) getWriterLocked(bailOnRotateFail, forceRotate bool) (io.Writer, error) {
	newName, byTime, bySize := f.filePathByRotate(forceRotate)
	if !byTime && !bySize {
		return f.usingFile, nil
	}
	if f.PreRotateHandler != nil {
		f.PreRotateHandler(f.usingFilePath)
	}
	newFile, err := f.rotateLocked(newName)
	if f.PostRotateHandler != nil {
		f.PostRotateHandler(f.usingFilePath)
	}
	if err != nil {
		if bailOnRotateFail {
			// Failure to rotate is a problem, but it's really not a great
			// idea to stop your application just because you couldn't rename
			// your log.
			//
			// We only return this error when explicitly needed (as specified by bailOnRotateFail)
			//
			// However, we *NEED* to close `fh` here
			if newFile != nil {
				_ = newFile.Close()
				newFile = nil
			}
			return nil, err
		}
	}
	if newFile == nil {
		return f.usingFile, nil
	}

	if f.usingFile != nil {
		_ = f.usingFile.Close()
		f.usingFile = nil
	}
	f.usingFile = newFile
	f.usingFilePath = newName

	return f.usingFile, nil
}

// file may not be nil if err is nil
func (f *RotateFile) rotateLocked(newName string) (*os.File, error) {
	var err error
	// if we got here, then we need to create a file
	switch f.RotateMode {
	case RotateModeCopyRename:
		// for which open the file, and write file not by RotateFile
		// CopyRenameFileAll = RenameFileAll(src->dst) + OpenFile(src)
		// usingFilePath->newName + recreate usingFilePath
		err = CopyRenameAll(newName, f.usingFilePath)
	case RotateModeCopyTruncate:
		// for which open the file, and write file not by RotateFile
		// CopyTruncateFile = CopyFile(src->dst) + Truncate(src)
		// usingFilePath->newName + truncate usingFilePath
		err = CopyTruncateAll(newName, f.usingFilePath)
	case RotateModeNew:
		// for which open the file, and write file by RotateFile
		fallthrough
	default:
	}
	if err != nil {
		return nil, err
	}
	file, err := AppendAllIfNotExist(newName)
	if err != nil {
		return nil, err
	}

	// link -> filename
	if f.FileLinkPath != "" {
		if err := ReSymlink(newName, f.FileLinkPath); err != nil {
			return nil, err
		}
	}

	now := time.Now()

	// find old files
	var filesNotExpired []string
	filesExpired, err := filepath_.GlobFunc(f.RotateFileGlob, func(name string) bool {
		fi, err := os.Stat(name)
		if err != nil {
			return false
		}

		fl, err := os.Lstat(name)
		if err != nil {
			return false
		}
		if f.MaxAge <= 0 {
			filesNotExpired = append(filesNotExpired, name)
			return false
		}

		if now.Sub(fi.ModTime()) < f.MaxAge {
			filesNotExpired = append(filesNotExpired, name)
			return false
		}

		if fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			return false
		}
		return true
	})
	if err != nil {
		return file, err
	}

	var filesExceedMaxCount []string
	if f.MaxCount > 0 && len(filesNotExpired) > 0 {
		removeCount := len(filesNotExpired) - f.MaxCount
		if removeCount < 0 {
			removeCount = 0
		}
		sort.Sort(rotateFileSlice(filesNotExpired))
		filesExceedMaxCount = filesNotExpired[:removeCount]
	}

	go func() {
		// unlink files on a separate goroutine
		for _, path := range filesExpired {
			os.Remove(path)
		}
		for _, path := range filesExceedMaxCount {
			os.Remove(path)
		}
	}()

	return file, nil
}

// foo.txt, 0 -> foo.txt
// foo.txt, 1 -> foo.txt.[1,2,...], which is not exist and seq is max
func nextSeqFileName(name string, seq int) (string, int) {
	// A new file has been requested. Instead of just using the
	// regular strftime pattern, we create a new file name using
	// generational names such as "foo.1", "foo.2", "foo.3", etc
	nf, seqUsed, err := NextFile(name+".*", seq)
	if err != nil {
		return name, seq
	}
	defer nf.Close()
	if seqUsed == 0 {
		return name, seqUsed
	}
	return nf.Name(), seqUsed
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

// foo.txt.* -> foo.txt.[1,2,...], which is exist and seq is max
func maxSeqFileName(name string) (string, int) {
	prefix, seq, suffix := MaxSeq(name + ".*")
	if seq == 0 {
		return name, seq
	}
	return fmt.Sprintf("%s%d%s", prefix, seq, suffix), seq
}

// sort filename by mode time and ascii in increase order
type rotateFileSlice []string

func (s rotateFileSlice) Len() int {
	return len(s)
}
func (s rotateFileSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s rotateFileSlice) Less(i, j int) bool {
	fi, err := os.Stat(s[i])
	if err != nil {
		return false
	}
	fj, err := os.Stat(s[j])
	if err != nil {
		return false
	}
	if fi.ModTime().Equal(fj.ModTime()) {
		if len(s[i]) == len(s[j]) {
			return s[i] < s[j]
		}
		return len(s[i]) > len(s[j]) // foo.1, foo.2, ..., foo
	}
	return fi.ModTime().Before(fj.ModTime())
}
