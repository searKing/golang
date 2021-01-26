package os

import (
	"os"
)

type FileInfos []os.FileInfo

func (s FileInfos) Len() int {
	return len(s)
}

func (s FileInfos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FileInfos) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}

// WalkFileInfo is a wrapper for sort of filepath.WalkFunc
type WalkFileInfo struct {
	Path     string
	FileInfo os.FileInfo
}

type WalkFileInfos []WalkFileInfo

func (w WalkFileInfos) Len() int {
	return len(w)
}

func (w WalkFileInfos) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w WalkFileInfos) Less(i, j int) bool {
	return w[i].Path < w[j].Path
}
