// Package eslice
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package eslice

import (
	"sync"
	"unsafe"

	"github.com/hyphennn/gcescape/internal"
)

const (
	threshold = 256
)

type ESlice[T any] struct {
	data   uintptr
	offset uintptr
	size   uintptr
	len    int
	cap    int
	freed  bool
	mu     sync.Mutex
}

// MakeESlice ESlice Expansion is very expensive.
// It is recommended to set the length during initialization to avoid expansion.
func MakeESlice[T any](args ...int) *ESlice[T] {
	if len(args) == 0 {
		panic("invalid args when calling make eslice")
	}
	l, c := args[0], args[0]
	if len(args) > 1 {
		if args[0] > args[1] {
			panic("len can not bigger than cap when calling make eslice")
		}
		c = args[1]
	}
	var t T
	ptr, offset, err := internal.MemAllocMmap[T](c)
	if err != nil {
		panic("make eslice failed: " + err.Error())
	}

	return &ESlice[T]{
		data:   ptr,
		offset: offset,
		size:   unsafe.Sizeof(t),
		len:    l,
		cap:    c,
		freed:  false,
		mu:     sync.Mutex{},
	}
}

func (e *ESlice[T]) Len() int {
	return e.len
}

func (e *ESlice[T]) Cap() int {
	return e.cap
}

func (e *ESlice[T]) Get(i int) T {
	if i >= e.len {
		panic("out of range for eslice")
	}
	return e.get(i)
}

func (e *ESlice[T]) get(i int) T {
	addr := e.data + uintptr(i)*e.size
	return *(*T)(unsafe.Pointer(addr))
}

func (e *ESlice[T]) Set(i int, t T) {
	if i >= e.len {
		panic("out of range for eslice")
	}
	e.set(i, t)
}

func (e *ESlice[T]) set(i int, t T) {
	addr := e.data + uintptr(i)*e.size
	p := (*T)(unsafe.Pointer(addr))
	*p = t
}

func (e *ESlice[T]) Free() {
	e.mu.Lock()
	defer e.mu.Unlock()
	err := internal.MemFreeMunmap(e.data, e.offset)
	if err != nil {
		panic("free eslice failed: " + err.Error())
	}
	e.freed = true
}

func (e *ESlice[T]) IsDangling() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.freed
}

// 如果不做此检查，则出现悬垂指针时会直接 fatal error
// todo 是否有必要做检查？还是交给用户？
func (e *ESlice[T]) checkDangling() {
	if e.IsDangling() {
		panic("dangling eslice")
	}
}

// Append If expansion occurs during the process, it will cause the original eslice to dangling.
// Try not to use the original eslice again.
func Append[T any](e *ESlice[T], ts ...T) *ESlice[T] {
	aim := len(ts) + e.len
	if aim > e.cap {
		e = scale(e, aim)
	}
	st := e.len
	for i := st; i < aim; i++ {
		e.len++
		e.set(i, ts[i-st])
	}
	return e
}

func scale[T any](e *ESlice[T], aim int) *ESlice[T] {
	if aim <= e.cap {
		return e
	}
	if aim > e.cap*2 {
		return realScale(e, aim)
	}
	if e.cap < threshold {
		return realScale(e, e.cap*2)
	}
	for e.cap < aim {
		var err error
		e = realScale(e, e.cap+(e.cap+3*threshold)/4)
		if err != nil {
			return e
		}
	}
	return e
}

func realScale[T any](e *ESlice[T], aim int) *ESlice[T] {
	e1 := MakeESlice[T](aim)
	for i := 0; i < e.len; i++ {
		e1.set(i, e.get(i))
	}
	e.Free()
	return e1
}
