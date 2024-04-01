// Package emap
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/29
package emap

import (
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
)

func Test111(t *testing.T) {
	s := MakeEMap[string, *TestType](10000000)
	for i := 0; i < 10000000; i++ {
		s.Set(strconv.Itoa(i), &TestType{
			Str:    "1",
			Map:    map[string]string{"1": "1"},
			Value:  0,
			Str2:   "1",
			Str3:   "1",
			Str4:   "1",
			Str5:   "1",
			Str6:   "1",
			Str7:   "1",
			Str8:   "1",
			Value2: 0,
			Value3: 0,
			Value4: 0,
			Value5: 0,
			Value6: 0,
			Value7: 0,
		})
	}

	os.Remove("cpu.pprof")
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0777)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < 10000000; i++ {
		v, ok := s.Get(strconv.Itoa(i))
		_ = v
		_ = ok
	}
	//avg, n, m := 0, 0, 0
	//
	//fmt.Println(s.len)
	//mm := map[int]int{}
	//for _, b := range s.m {
	//	l := 0
	//	for {
	//		for j := 0; j < 8; j++ {
	//			if b[j] == -1 || b[j] == 0 {
	//				continue
	//			}
	//			it := (*item[int, int])(unsafe.Pointer(s.geti(b[j])))
	//			if it.k == 0 {
	//				fmt.Println(11)
	//			}
	//			mm[it.k]++
	//			l++
	//		}
	//		// 当前桶未找到
	//		if b[8] == 0 {
	//			break
	//		} else {
	//			// 有下一个桶
	//			b = (*bucket)(unsafe.Pointer(uintptr(b[8])))
	//			continue
	//		}
	//	}
	//	if l != 0 {
	//		avg += l
	//		n++
	//	}
	//	if l > m {
	//		m = l
	//	}
	//}
	//fmt.Println(avg)
	//fmt.Println(n)
	//fmt.Println(float64(avg) / float64(n))
	//fmt.Println(m)

	//for i := 0; i < 256; i++ {
	//	if mm[i] == 0 {
	//		fmt.Println("is 0")
	//		fmt.Println(i)
	//	}
	//	if mm[i] > 1 {
	//		fmt.Println("is big")
	//		fmt.Println(i, mm[i])
	//	}
	//}
}

func Test222(t *testing.T) {
	s := make(map[string]*TestType, 10000000)
	for i := 0; i < 10000000; i++ {
		s[strconv.Itoa(i)] = &TestType{
			Str:    "1",
			Map:    map[string]string{"1": "1"},
			Value:  0,
			Str2:   "1",
			Str3:   "1",
			Str4:   "1",
			Str5:   "1",
			Str6:   "1",
			Str7:   "1",
			Str8:   "1",
			Value2: 0,
			Value3: 0,
			Value4: 0,
			Value5: 0,
			Value6: 0,
			Value7: 0,
		}
	}

	os.Remove("cpu2.pprof")
	f, _ := os.OpenFile("cpu2.pprof", os.O_CREATE|os.O_RDWR, 0777)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < 10000000; i++ {
		v, ok := s[strconv.Itoa(i)]
		_ = v
		_ = ok
	}
}

type TestType struct {
	Str    string
	Map    map[string]string
	Value  int
	Str2   string
	Str3   string
	Str4   string
	Str5   string
	Str6   string
	Str7   string
	Str8   string
	Value2 int
	Value3 int
	Value4 int
	Value5 int
	Value6 int
	Value7 int
}
