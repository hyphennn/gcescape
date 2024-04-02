// Package emap
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/29
package emap

import (
	"unsafe"

	"github.com/hyphennn/gcescape/internal"
)

const (
	threshold = 256
)

type item[K comparable, V any] struct {
	k K
	v V
}

type bucket struct {
	tophash [8]uint8
	vs      [8]int
	next    *bucket
}

type bpool struct {
	data   uintptr
	offset uintptr
	size   uintptr
	len    int
	cap    int

	u *stk

	next *bpool
}

func newBpool(c int) *bpool {
	data, offset, err := internal.MemAllocMmap[bucket](c)
	if err != nil {
		panic("make bpool failed: " + err.Error())
	}
	var b bucket
	return &bpool{
		data:   data,
		offset: offset,
		size:   unsafe.Sizeof(b),
		len:    0,
		cap:    c,
		u: &stk{
			s: make([]int, 10),
			t: 0,
		},
		next: nil,
	}
}

func (b *bpool) get() *bucket {
	var ret *bucket
	if b.u.empty() {
		// 没有未使用空间
		if b.len == b.cap {
			if b.next == nil {
				b.scale()
			}
			ret = b.next.get()
			goto clear
		}
		p := b.data + b.size*uintptr(b.len)
		b.len++
		ret = (*bucket)(unsafe.Pointer(p))
		goto clear
	}
	ret = (*bucket)(unsafe.Pointer(uintptr(b.u.pop())))
	goto clear
clear:
	for i := 0; i < 8; i++ {
		ret.tophash[i] = 0
		ret.vs[i] = 0
	}
	ret.next = nil
	return ret
}

func (b *bpool) scale() {
	b.next = newBpool(b.cap)
}

// p不存在于pool中的时候，不会有任何事情发生
func (b *bpool) free(p uintptr) {
	if b.data > p || b.offset < p {
		if b.next != nil {
			b.next.free(p)
		}
	}
	b.u.push(int(p))
}

func (b *bpool) freeAll() {
	err := internal.MemFreeMunmap(b.data, b.offset)
	if err != nil {
		panic("free b pool failed: " + err.Error())
	}
}

type stk struct {
	s []int
	t int
}

func (s *stk) push(v int) {
	if s.t >= len(s.s) {
		s.s = append(s.s, v)
	} else {
		s.s[s.t] = v
	}
	s.t++
}

func (s *stk) pop() int {
	s.t--
	return s.s[s.t]
}

func (s *stk) empty() bool {
	return s.t == 0
}
