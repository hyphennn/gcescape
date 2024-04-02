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

	"github.com/hyphennn/gcescape/internal"
)

func TestEmapPprof(t *testing.T) {
	c := 10000
	s := MakeEMap[string, *internal.TestType](c)
	for i := 0; i < c; i++ {
		si := strconv.Itoa(i)
		s.Set(si, internal.GenTestType())
	}

	os.Remove("cpu.pprof")
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0777)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < c; i++ {
		si := strconv.Itoa(i)
		v, ok := s.Get(si)
		_, _ = v, ok
	}
}

func TestMapPprof(t *testing.T) {
	c := 1000000
	s := make(map[string]*internal.TestType, c)
	for i := 0; i < c; i++ {
		s[strconv.Itoa(i)] = internal.GenTestType()
	}

	os.Remove("cpu2.pprof")
	f, _ := os.OpenFile("cpu2.pprof", os.O_CREATE|os.O_RDWR, 0777)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < c; i++ {
		v, ok := s[strconv.Itoa(i)]
		_ = v
		_ = ok
	}
}
