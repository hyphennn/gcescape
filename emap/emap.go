// Package emap
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package emap

import (
	"sync"
	"unsafe"

	"github.com/hyphennn/gcescape/internal"
)

// MakeEMap 此版本不支持扩容
// 后续会支持，但扩容代价极大，扩容不会导致指针更改，和标准库保持一致
// 任何并发写操作都是极其危险的，因为存在绕过 Runtime 的内核调用，Runtime 和编译器无法兜底
func MakeEMap[K comparable, V any](cap int) *EMap[K, V] {
	data, offset, err := internal.MemAllocMmap[item[K, V]](cap)
	if err != nil {
		panic("make emap failed: " + err.Error())
	}

	b := uint8(0)
	for internal.OverLoadFactor(cap, b) {
		b++
	}

	var it item[K, V]
	return &EMap[K, V]{
		data:   data,
		offset: offset,
		size:   unsafe.Sizeof(it),
		cap:    cap,
		len:    0,
		u: &stk{
			// 常数可以调整
			s: make([]int, 10),
			t: 0,
		},
		// 常数可以调整
		m: make(map[K]int, cap),

		freed: false,
		mu:    sync.Mutex{},
	}
}

type EMap[K comparable, V any] struct {
	data   uintptr
	offset uintptr
	size   uintptr
	cap    int
	len    int

	u *stk

	m map[K]int

	freed bool
	mu    sync.Mutex
}

func (m *EMap[K, V]) Set(k K, v V) {
	// 查找是否存在，存在则直接更新 v
	i, ok := m.m[k]
	if !ok {
		// 直接新增
		p := m.findUnusedDataSet(k, v)
		m.m[k] = p
		return
	}
	p := (*V)(unsafe.Pointer(m.geti(i)))
	*p = v
	return
}

// 返回的不是绝对地址而是 index
func (m *EMap[K, V]) findUnusedDataSet(k K, v V) int {
	var idx int
	if m.u.empty() {
		if m.len == m.cap {
			m.autoScale()
		}
		idx = m.len
		m.len++
	} else {
		idx = m.u.pop()
	}

	i := (*item[K, V])(unsafe.Pointer(m.geti(idx)))
	i.k, i.v = k, v
	return idx
}

func (m *EMap[K, V]) autoScale() {
	if m.cap < threshold {
		m.realScale(m.cap * 2)
		return
	}
	m.realScale(m.cap + (m.cap+3*threshold)/4)
	return
}

func (m *EMap[K, V]) realScale(aim int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, offset, err := internal.ScaleG[item[K, V]](m.data, m.cap, aim, &m.data)
	if err != nil {
		panic("scale e map failed: " + err.Error())
	}
	m.offset = offset
	m.cap = aim
}

func (m *EMap[K, V]) Get(k K) (v V, ok bool) {
	i, ok := m.m[k]
	if !ok {
		return
	}
	v = *(*V)(unsafe.Pointer(m.geti(i)))
	ok = true
	return
}

func (m *EMap[K, V]) geti(i int) uintptr {
	return m.data + uintptr(i)*m.size
}

func (m *EMap[K, V]) freei(i int) {
	m.u.push(i)
}

func (m *EMap[K, V]) Delete(k K) bool {
	i, ok := m.m[k]
	if !ok {
		return false
	}
	m.freei(i)
	delete(m.m, k)
	return true
}

// Free 只会释放向内核申请的部分，剩余由gc管理
// 一旦 free，任何指针都可能是悬垂的
func (m *EMap[K, V]) Free() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.freed {
		return
	}

	err := internal.MemFreeMunmap(m.data, m.offset)
	if err != nil {
		panic("free e map data failed: " + err.Error())
	}

	m.freed = true
	return
}

func (m *EMap[K, V]) IsDangling() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.freed
}
