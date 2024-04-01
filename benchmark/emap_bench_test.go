// Package benchmark
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/29
package benchmark

import (
	"strconv"
	"testing"

	"github.com/hyphennn/gcescape/emap"
)

func BenchmarkEmapSet(b *testing.B) {
	s := emap.MakeEMap[int, int](1000000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(i, i)
	}
}

func BenchmarkMapSet(b *testing.B) {
	s := make(map[int]int, 100000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s[i] = i
	}
}

const mapGetSize = 1000000

func BenchmarkEMapGet(b *testing.B) {
	s := emap.MakeEMap[string, *TestType](mapGetSize)
	for i := 0; i < mapGetSize; i++ {
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, ok := s.Get(strconv.Itoa(i))
		_ = v
		_ = ok
	}
}

func BenchmarkMapGet(b *testing.B) {
	s := make(map[string]*TestType, mapGetSize)
	for i := 0; i < mapGetSize; i++ {
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
