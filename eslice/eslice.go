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

func (s *ESlice[T]) Len() int {
	return s.len
}

func (s *ESlice[T]) Cap() int {
	return s.cap
}

func (s *ESlice[T]) Get(i int) T {
	if i >= s.len {
		panic("out of range for eslice")
	}
	return s.get(i)
}

func (s *ESlice[T]) get(i int) T {
	addr := s.data + uintptr(i)*s.size
	return *(*T)(unsafe.Pointer(addr))
}

func (s *ESlice[T]) Set(i int, t T) {
	if i >= s.len {
		panic("out of range for eslice")
	}
	s.set(i, t)
}

func (s *ESlice[T]) set(i int, t T) {
	addr := s.data + uintptr(i)*s.size
	p := (*T)(unsafe.Pointer(addr))
	*p = t
}

func (s *ESlice[T]) Free() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.freed {
		return
	}
	err := internal.MemFreeMunmap(s.data, s.offset)
	if err != nil {
		panic("free eslice failed: " + err.Error())
	}
	s.freed = true
}

func (s *ESlice[T]) IsDangling() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.freed
}

// 如果不做此检查，则出现悬垂指针时会直接 fatal error
// todo 是否有必要做检查？还是交给用户？
func (s *ESlice[T]) checkDangling() {
	if s.IsDangling() {
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
	e.mu.Lock()
	e.freed = true
	e.mu.Unlock()

	ndata, offset, err := internal.ScaleG[T](e.data, e.cap, aim, nil)
	if err != nil {
		panic("scale e slice failed: " + err.Error())
	}
	return &ESlice[T]{
		data:   ndata,
		offset: offset,
		size:   e.size,
		len:    e.len,
		cap:    aim,
		freed:  false,
		mu:     sync.Mutex{},
	}
}
