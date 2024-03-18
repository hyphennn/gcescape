// Package benchmark
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package benchmark

import (
	"testing"

	"github.com/hyphennn/gcescape/eobject"
)

func BenchmarkEObjectNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = eobject.MakeEObject[int]()
	}
}

func BenchmarkObjectNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = new(int)
	}
}

func BenchmarkEObjectSet(b *testing.B) {
	a := eobject.MakeEObject[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Set(1)
		c := a.Value()
		_ = c
	}
}

func BenchmarkObjectSet(b *testing.B) {
	a := 0
	_ = a
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a = 1
		c := a
		_ = c
	}
}
