// Package benchmark
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package benchmark

import (
	"testing"

	"github.com/hyphennn/gcescape/eslice"
)

func BenchmarkEsliceAppend(b *testing.B) {
	s := eslice.MakeESlice[int64](0, 1000000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = eslice.Append(s, 1)
	}
}

func BenchmarkSliceAppend(b *testing.B) {
	s := make([]int64, 0, 1000000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = append(s, 1)
	}
}

func BenchmarkMakeESlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = eslice.MakeESlice[int64](0, 1)
	}
}

func BenchmarkMakeSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = make([]int64, 0, 1)
	}
}

func BenchmarkESliceAppendScale(b *testing.B) {
	s := eslice.MakeESlice[int64](1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = eslice.Append(s, 1)
	}
}

func BenchmarkSliceAppendScale(b *testing.B) {
	s := make([]int64, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = append(s, 1)
	}
}
