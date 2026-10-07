// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

// Code borrowed from https://github.com/minio/minio/blob/master/internal/disk/stat_windows.go

package os

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	// GetDiskFreeSpaceEx - https://msdn.microsoft.com/en-us/library/windows/desktop/aa364937(v=vs.85).aspx
	// Retrieves information about the amount of space that is available on a disk volume,
	// which is the total amount of space, the total amount of free space, and the total
	// amount of free space available to the user that is associated with the calling thread.
	GetDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")

	// GetDiskFreeSpace - https://msdn.microsoft.com/en-us/library/windows/desktop/aa364935(v=vs.85).aspx
	// Retrieves information about the specified disk, including the amount of free space on the disk.
	GetDiskFreeSpace = kernel32.NewProc("GetDiskFreeSpaceW")
)

// DiskUsage returns total and free bytes available in a directory, e.g. `C:\`.
// It returns free space available to the user (including quota limitations)
// inodes and inodesFree are total and free clusters, 0 if unknown.
//
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa364937(v=vs.85).aspx
func DiskUsage(path string) (total int64, free int64, avail int64, inodes int64, inodesFree int64, err error) {
	// Stat to know if the path exists.
	if _, err = os.Stat(path); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	if err := windows.GetDiskFreeSpaceEx(pathPtr, &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes); err != nil {
		return 0, 0, 0, 0, 0, &os.PathError{Op: "GetDiskFreeSpaceEx", Path: path, Err: err}
	}

	if totalNumberOfFreeBytes > totalNumberOfBytes {
		return 0, 0, 0, 0, 0, fmt.Errorf("detected free space (%d) > total disk space (%d), fs corruption at (%s). please run 'fsck'",
			totalNumberOfFreeBytes, totalNumberOfBytes, path)
	}

	// Return values of GetDiskFreeSpace()
	lpSectorsPerCluster := uint32(0)
	lpBytesPerSector := uint32(0)
	lpNumberOfFreeClusters := uint32(0)
	lpTotalNumberOfClusters := uint32(0)

	// Extract values safely
	// BOOL WINAPI GetDiskFreeSpace(
	//   _In_  LPCTSTR lpRootPathName,
	//   _Out_ LPDWORD lpSectorsPerCluster,
	//   _Out_ LPDWORD lpBytesPerSector,
	//   _Out_ LPDWORD lpNumberOfFreeClusters,
	//   _Out_ LPDWORD lpTotalNumberOfClusters
	// );
	// clusters are left 0 if failed, as lpRootPathName may be required to be the root of a disk.
	_, _, _ = GetDiskFreeSpace.Call(uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&lpSectorsPerCluster)),
		uintptr(unsafe.Pointer(&lpBytesPerSector)),
		uintptr(unsafe.Pointer(&lpNumberOfFreeClusters)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfClusters)))

	return int64(totalNumberOfBytes), int64(totalNumberOfFreeBytes), int64(freeBytesAvailable), int64(lpTotalNumberOfClusters), int64(lpNumberOfFreeClusters), nil
}
