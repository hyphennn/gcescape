// Package eobject
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package eobject

import (
	"unsafe"

	"github.com/hyphennn/gcescape/internal"
)

type EObject[T any] uintptr

func MakeEObject[T any]() EObject[T] {
	ptr, _, err := internal.MemAllocMmap[T](1)
	if err != nil {
		panic("make eobject failed: " + err.Error())
	}
	return EObject[T](ptr)
}

func (e EObject[T]) Value() T {
	return *(*T)(unsafe.Pointer(e))
}

func (e EObject[T]) Set(t T) {
	p := (*T)(unsafe.Pointer(e))
	*p = t
}

func (e EObject[T]) Free() {
	err := internal.MemFreeMunmapG[T](uintptr(e), 1)
	if err != nil {
		panic("free eobject failed: " + err.Error())
	}
}
