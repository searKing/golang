package os

import (
	"os"
	"path/filepath"
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

func MakdirAllIfNeeded(path string, perm os.FileMode) (created bool, err error) {
	dir, _ := filepath.Split(path)
	has, err := PathExists(dir)
	if err != nil {
		return false, err
	}
	if has {
		return false, nil
	}
	err = os.MkdirAll(dir, perm)
	if err != nil {
		return false, err
	}
	return true, nil
}
