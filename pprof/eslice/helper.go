// Package main
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/hyphennn/gcescape/eobject"
	"github.com/hyphennn/gcescape/eslice"
	"github.com/hyphennn/gcescape/internal"
)

func UseESlice() {
	s := eslice.MakeESlice[*internal.TestType](0, 10000000)

	for i := 0; i < 10000000; i++ {
		s = eslice.Append(s, internal.GenTestType())
	}

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}

func UseNormalSlice() {
	s := make([]internal.TestType, 0, 10000000)

	for i := 0; i < 10000000; i++ {
		s = append(s, *internal.GenTestType())
	}

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}

func UseNormalPtrSlice() {
	s := make([]*internal.TestType, 10000000)

	for i := 0; i < 10000000; i++ {
		s = append(s, internal.GenTestType())
	}

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}

func UseObject() {
	s := [10000000]*internal.TestType{}
	for i := 0; i < 10000000; i++ {
		s[i] = new(internal.TestType)
	}

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}

func UseEObject() {
	s := [10000000]eobject.EObject[internal.TestType]{}
	for i := 0; i < 10000000; i++ {
		s[i] = eobject.MakeEObject[internal.TestType]()
	}

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}
