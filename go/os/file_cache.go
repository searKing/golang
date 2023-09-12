// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/searKing/golang/go/crypto/md5"
	filepath_ "github.com/searKing/golang/go/path/filepath"
)

// CacheFile is a package cache(Eventual consistency), backed by a file system directory tree.
//
// It is safe for multiple processes on a single machine to use the
// same cache directory in a local file system simultaneously.
// They will coordinate using operating system file locks and may
// duplicate effort but will not corrupt the cache.
//
// However, it is NOT safe for multiple processes on different machines
// to share a cache directory (for example, if the directory were stored
// in a network file system). File locking is notoriously unreliable in
// network file systems and may not suffice to protect the cache.
//
//go:generate go-option -type "CacheFile"
type CacheFile struct {
	BucketRootDir string // cache root dir
	// generate bucket key from key(file path)
	// bucket key should not contain any of the magic characters recognized by [filepath.Match]
	// otherwise, bucket key will be escaped by MD5CacheKey
	// see: https://github.com/golang/go/issues/13516
	BucketKeyFunc func(key string) string

	CacheMetaExt      string        // the file name extension used by path. ".cache" if empty
	CacheExpiredAfter time.Duration // Cache file expiration time, lazy expire cache files base on cache URL modification time
}

func NewCacheFile(opts ...CacheFileOption) *CacheFile {
	var f CacheFile
	f.ApplyOptions(opts...)
	if f.CacheMetaExt == "" {
		f.CacheMetaExt = ".cache"
	}
	if k := f.BucketKeyFunc; k != nil {
		f.BucketKeyFunc = func(key string) string {
			bk := k(key)
			if hasMeta(bk) {
				return MD5CacheKey(bk)
			}
			return bk
		}
	}
	return &f
}

func (f *CacheFile) BucketKey(name string) string {
	if f.BucketKeyFunc != nil {
		return f.BucketKeyFunc(name)
	}
	return MD5CacheKey(name)
}

// Get looks up the file in the cache and returns
// the cache name of the corresponding data file.
func (f *CacheFile) Get(name string) (cacheFilePath, cacheMetaPath string, hit bool, err error) {
	cacheFilePath = filepath.Join(f.BucketRootDir, f.BucketKey(name))
	cacheMetaPath = cacheFilePath + f.CacheMetaExt

	cacheFilePath, cacheMetaPath, err = f.getCacheMeta(name, cacheFilePath+".*"+f.CacheMetaExt)
	if err != nil {
		return "", "", false, err
	}

	// double check
	// handle special case below:
	// 1. cache hit
	// <protect> 2. cache removed by other process or goroutine -- not controllable
	// 3. refresh cache's ModTime
	{
		info, err_ := os.Stat(cacheMetaPath)
		if err_ != nil {
			hit = false
			return
		}
		// violate cache file if cache expired
		expired := time.Since(info.ModTime()) > f.CacheExpiredAfter
		if expired {
			hit = false
			return
		}
		info, err_ = os.Stat(cacheFilePath)
		if err_ != nil {
			hit = false
			return
		}
		if info.Size() == 0 {
			hit = false
			return
		}
	}

	// STEP3 cache url not conflict, refresh ModTime of cache file and cache file's metadata
	hit = true
	_ = ChtimesNow(cacheMetaPath)
	_ = ChtimesNow(cacheFilePath)
	return
}

func (f *CacheFile) Put(name string, r io.Reader) (cacheFilePath string, refreshed bool, err error) {
	cacheFilePath, cacheMetaPath, hit, err := f.Get(name)
	if err != nil {
		return "", false, err
	}
	if hit {
		return cacheFilePath, false, nil
	}

	err = WriteRenameAllFrom(cacheFilePath, r)
	if err != nil {
		return "", false, fmt.Errorf("failed to create cache file: %w", err)
	}
	err = WriteRenameAll(cacheMetaPath, []byte(name))
	if err != nil {
		_ = os.Remove(cacheFilePath)
		return "", false, fmt.Errorf("failed to create cache meta: %w", err)
	}
	return cacheFilePath, true, nil
}

func (f *CacheFile) getCacheMeta(key, cacheMetaPathPattern string) (cacheFilePath, cacheMetaPath string, err error) {
	var hitMeta bool

	// STEP1 clean expired cache file
	_ = filepath_.WalkGlob(cacheMetaPathPattern, func(path string) error {
		cacheMetaPath = path
		cacheFilePath = strings.TrimSuffix(cacheMetaPath, f.CacheMetaExt)
		info, err := os.Stat(cacheMetaPath)
		if err == nil {
			// violate cache file if cache expired
			expired := time.Since(info.ModTime()) > f.CacheExpiredAfter
			if expired {
				_ = os.Remove(cacheFilePath)
				_ = os.Remove(cacheMetaPath)
				return nil
			}
		}
		return nil
	})

	// STEP2 search for cache file in cache open list
	_ = filepath_.WalkGlob(cacheMetaPathPattern, func(path string) error {
		cacheMetaPath = path
		cacheFilePath = strings.TrimSuffix(cacheMetaPath, f.CacheMetaExt)
		// verify whether if cache key in cache file is match
		keyInCache, _ := os.ReadFile(cacheMetaPath)
		if string(keyInCache) == key {
			hitMeta = true
			return filepath.SkipAll
		}
		// cache key conflict, continue search cache file list
		return nil
	})
	if hitMeta {
		return
	}

	// STEP3 add new cache file to cache open list
	// foo.txt.* -> foo.txt.[0,1,2,...], which exists and seq is max
	nf, _, err := NextFile(cacheMetaPathPattern, 0)
	if err != nil {
		return "", "", fmt.Errorf("failed to open next cache meta: %w", err)
	}
	defer nf.Close()
	// STEP3 cache url not conflict, refresh ModTime of cache file and cache file's metadata
	cacheMetaPath = nf.Name()
	cacheFilePath = strings.TrimSuffix(cacheMetaPath, f.CacheMetaExt)
	_, err = nf.WriteString(key)
	if err != nil {
		return "", "", fmt.Errorf("failed to write next cache meta: %w", err)
	}
	_ = os.Remove(cacheFilePath)
	return
}

func MD5CacheKey(s string) string {
	// Special CASE 1: filename-as-part-of-a-query-string
	// http://foo.com?url=http://bar.com/kitty.jpg&filename=kitty.jpg
	// https://stackoverflow.com/questions/28915717/how-does-one-safely-pass-a-url-and-filename-as-part-of-a-query-string
	// Special CASE 2: filename-as-url-path，but different by various version in query
	// http://bucket.s3.amazonaws.com/my-image.jpg?versionId=L4kqtJlcpXroDTDmpUMLUo
	// https://docs.aws.amazon.com/AmazonS3/latest/userguide/RetrievingObjectVersions.html
	// https://cloud.tencent.com/document/product/436/19883
	return md5.SumHex(s)
}
