// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package os

import (
	"fmt"
	"os"
	"syscall"
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
//
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa364937(v=vs.85).aspx
func DiskUsage(path string) (total uint64, free uint64, avail uint64, inodes uint64, inodesFree uint64, err error) {
	// Stat to know if the path exists.
	if _, err = os.Stat(path); err != nil {
		return 0, 0, 0, 0, 0, err
	}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	// Extract values safely
	// BOOL WINAPI GetDiskFreeSpaceEx(
	// _In_opt_  LPCTSTR         lpDirectoryName,
	// _Out_opt_ PULARGE_INTEGER lpFreeBytesAvailable,
	// _Out_opt_ PULARGE_INTEGER lpTotalNumberOfBytes,
	// _Out_opt_ PULARGE_INTEGER lpTotalNumberOfFreeBytes
	// );
	_, _, _ = GetDiskFreeSpaceEx.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))

	if uint64(lpTotalNumberOfFreeBytes) > uint64(lpTotalNumberOfBytes) {
		return 0, 0, 0, 0, 0, fmt.Errorf("detected free space (%d) > total disk space (%d), fs corruption at (%s). please run 'fsck'",
			uint64(lpTotalNumberOfFreeBytes), uint64(lpTotalNumberOfBytes), path)
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
	_, _, _ = GetDiskFreeSpace.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpSectorsPerCluster)),
		uintptr(unsafe.Pointer(&lpBytesPerSector)),
		uintptr(unsafe.Pointer(&lpNumberOfFreeClusters)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfClusters)))

	return uint64(lpTotalNumberOfBytes), uint64(lpTotalNumberOfFreeBytes), lpFreeBytesAvailable, uint64(lpTotalNumberOfClusters), uint64(lpNumberOfFreeClusters), err
}
