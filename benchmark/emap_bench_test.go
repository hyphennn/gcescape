// Package benchmark
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/29
package benchmark

import (
	"strconv"
	"testing"

	"github.com/hyphennn/gcescape/emap"
	"github.com/hyphennn/gcescape/internal"
)

const mapSize = 100000

func BenchmarkEmapSet(b *testing.B) {
	s := emap.MakeEMap[int, int](mapSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(i%mapSize, i)
	}
}

func BenchmarkMapSet(b *testing.B) {
	s := make(map[int]int, mapSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s[i%mapSize] = i
	}
}

func BenchmarkEMapGet(b *testing.B) {
	s := emap.MakeEMap[string, *internal.TestType](mapSize)
	for i := 0; i < mapSize; i++ {
		s.Set(strconv.Itoa(i), internal.GenTestType())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, ok := s.Get(strconv.Itoa(i))
		_ = v
		_ = ok
	}
}

func BenchmarkMapGet(b *testing.B) {
	s := make(map[string]*internal.TestType, mapSize)
	for i := 0; i < mapSize; i++ {
		s[strconv.Itoa(i)] = internal.GenTestType()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, ok := s[strconv.Itoa(i)]
		_ = v
		_ = ok
	}
}
