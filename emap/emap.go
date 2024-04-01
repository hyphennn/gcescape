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
	h, s := internal.GetRuntimeHasher[K]()

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
		m: make(map[uintptr]uintptr, cap),
		//s: eslice.MakeESlice[uintptr](1 << 24),
		//s:      make([]uintptr, 1<<24),

		p:      newBpool(cap),
		hasher: h,
		seed:   s,
		freed:  false,
		mu:     sync.Mutex{},
	}
}

type EMap[K comparable, V any] struct {
	data   uintptr
	offset uintptr
	size   uintptr
	cap    int
	len    int

	u *stk

	m map[uintptr]uintptr
	//s []uintptr
	//s *eslice.ESlice[uintptr]

	p *bpool

	hasher internal.Hasher
	seed   uintptr

	freed bool
	mu    sync.Mutex
}

func (m *EMap[K, V]) Set(k K, v V) {
	// 查找是否存在，存在则直接更新 v
	h := internal.Tophash(m.hasher(unsafe.Pointer(&k), m.seed))
	b, ok := m.m[h]
	if !ok {
		// 直接新增
		b = m.p.get()
		p := m.findUnusedDataSet(k, v)
		*(*uintptr)(unsafe.Pointer(b)) = uintptr(p)
		m.m[h] = b
		return
	}
	//b := (*bucket)(unsafe.Pointer(bp))
	i := 0
	// first empty bucket index
	var feb uintptr
	var febi int
	for {
		for i = 0; i < 8; i++ {
			bi := *(*int)(unsafe.Pointer(b + uintptr(i*8)))
			if bi == 0 {
				// 当前位置为空，直接在当前位置新增
				if feb == 0 {
					goto insert
				} else {
					goto reuse
				}
			}
			if bi == -1 {
				if feb == 0 {
					feb = b
					febi = i
				}
				continue
			}
			it := (*item[K, V])(unsafe.Pointer(m.geti(bi)))
			if it.k == k {
				// 找到目标，直接更新 v
				it.v = v
				return
			}
		}
		// 当前桶未找到
		b8 := *(*uintptr)(unsafe.Pointer(b + uintptr(8*8)))
		if b8 == 0 {
			// 没有下一个桶，说明所有桶中均未找到，并且当前桶满
			if feb != 0 {
				// 前面的 bucket 存在可重用的空位
				goto reuse
			}
			goto newbucket
		} else {
			// 有下一个桶
			b = b8
			continue
		}
	}
insert:
	*(*int)(unsafe.Pointer(b + uintptr(i*8))) = m.findUnusedDataSet(k, v)
	return
reuse:
	// 未找到，且当前桶不为空
	*(*int)(unsafe.Pointer(feb + uintptr(febi*8))) = m.findUnusedDataSet(k, v)
	return
newbucket:
	nbp := m.p.get()
	*(*int)(unsafe.Pointer(nbp)) = m.findUnusedDataSet(k, v)
	*(*uintptr)(unsafe.Pointer(b + uintptr(8*8))) = nbp
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

	i := (*item[K, V])(unsafe.Pointer(m.geti(idx + 1)))
	i.k, i.v = k, v
	return idx + 1
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
	h := internal.Tophash(m.hasher(unsafe.Pointer(&k), m.seed))
	b, ok := m.m[h]
	if !ok {
		return
	}

	for {
		for i := 0; i < 8; i++ {
			bi := *(*int)(unsafe.Pointer(b + uintptr(i*8)))
			if bi == 0 {
				return
			}
			if bi == -1 {
				continue
			}
			it := (*item[K, V])(unsafe.Pointer(m.geti(bi)))
			if it.k == k {
				ok = true
				v = it.v
				return
			}
		}
		// 当前桶未找到
		b8 := *(*uintptr)(unsafe.Pointer(b + uintptr(8*8)))
		if b8 == 0 {
			return
		} else {
			// 有下一个桶
			b = b8
			continue
		}
	}
}

func (m *EMap[K, V]) geti(i int) uintptr {
	//return add(m.data, uintptr(i-1)*m.size)
	return m.data + uintptr(i-1)*m.size
}

func (m *EMap[K, V]) Delete(k K) bool {
	h := internal.Tophash(m.hasher(unsafe.Pointer(&k), m.seed))
	b, ok := m.m[h]
	if !ok {
		return false
	}

	//b := (*bucket)(unsafe.Pointer(bp))
	var pb uintptr
	i := 0
	for {
		for i = 0; i < 8; i++ {
			bi := (*int)(unsafe.Pointer(b + uintptr(i*8)))
			if *bi == 0 {
				return false
			}
			if *bi == -1 {
				continue
			}
			p := m.geti(*bi)
			it := (*item[K, V])(unsafe.Pointer(p))
			if it.k == k {
				// 真实删除逻辑
				m.u.push(*bi - 1)
				*bi = -1
				if i == 0 {
					// 此时桶中不再有任何元素
					for i := 0; i < 9; i++ {
						*bi = 0
					}
					m.p.free(b)
					if pb == 0 {
						delete(m.m, h)
					} else {
						*(*uintptr)(unsafe.Pointer(pb + uintptr(8*8))) = 0
					}
				}
				return true
			}
		}
		// 当前桶未找到
		b8 := *(*uintptr)(unsafe.Pointer(b + uintptr(8*8)))
		if b8 == 0 {
			// 没有下一个桶
			return false
		} else {
			// 有下一个桶
			pb = b
			b = b8
			continue
		}
	}
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

	m.p.freeAll()

	m.freed = true
	return
}

func (m *EMap[K, V]) IsDangling() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.freed
}
