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
)

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

func UseESlice() {
	s := eslice.MakeESlice[*TestType](0, 10000000)

	for i := 0; i < 10000000; i++ {
		s = eslice.Append(s, &TestType{
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

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}

func UseNormalSlice() {
	s := make([]TestType, 0, 10000000)

	for i := 0; i < 10000000; i++ {
		s = append(s, TestType{
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

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}

func UseNormalPtrSlice() {
	s := make([]*TestType, 10000000)

	for i := 0; i < 10000000; i++ {
		s = append(s, &TestType{
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

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}

func UseObject() {
	s := [10000000]*TestType{}
	for i := 0; i < 10000000; i++ {
		s[i] = new(TestType)
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
	s := [10000000]eobject.EObject[TestType]{}
	for i := 0; i < 10000000; i++ {
		s[i] = eobject.MakeEObject[TestType]()
	}

	for i := 0; i < 10; i++ {
		st := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(st))
		time.Sleep(time.Second)
	}

	runtime.KeepAlive(s)
}
