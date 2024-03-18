// Package emap
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package emap

import (
	"unsafe"

	"github.com/hyphennn/gcescape/internal"
)

type EMap[K comparable, V any] struct {
	data uintptr
	size uintptr
	cap  int
	len  int

	m map[uintptr]uintptr

	hasher internal.Hasher
}

type item[K comparable, V any] struct {
	k    K
	v    V
	pre  uintptr
	next uintptr
}

func NewEmap() {
	f := func(unsafe.Pointer, uintptr) uintptr
	e := EMap[int, int]{}
	e.hasher = f
}

func WithHasher() {

}
