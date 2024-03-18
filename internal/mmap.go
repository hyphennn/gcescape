// Package internal
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package internal

import (
	"syscall"
	"unsafe"
)

func MemAllocMmap[T any](len int) (uintptr, uintptr, error) {
	fd := -1
	var t T
	offset := uintptr(len) * unsafe.Sizeof(t)
	data, _, errno := syscall.Syscall6(
		syscall.SYS_MMAP,
		0,
		offset,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE,
		uintptr(fd),
		0,
	)
	if errno != 0 {
		return 0, 0, errno
	}
	return data, offset, nil
}

func MemFreeMunmap(addr uintptr, offset uintptr) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_MUNMAP,
		addr,
		offset,
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}

func MemFreeMunmapG[T any](addr uintptr, len int) error {
	var t T
	offset := uintptr(len) * unsafe.Sizeof(t)
	_, _, errno := syscall.Syscall(
		syscall.SYS_MUNMAP,
		addr,
		offset,
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}
