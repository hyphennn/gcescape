// Package emap
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/29
package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/hyphennn/gcescape/emap"
	"github.com/hyphennn/gcescape/internal"
)

func UseEmap() {
	s := emap.MakeEMap[string, *internal.TestType](1000000)

	for i := 0; i < 1000000; i++ {
		s.Set(strconv.Itoa(i), internal.GenTestType())
	}

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}

func UseNormalmap() {
	s := make(map[string]*internal.TestType, 2000000)

	for i := 0; i < 1000000; i++ {
		s[strconv.Itoa(i)] = internal.GenTestType()
	}

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}
