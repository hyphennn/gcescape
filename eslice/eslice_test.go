// Package eslice
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package eslice

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestESlice(t *testing.T) {
	e := MakeESlice[int](1)
	fmt.Println(e.Len())
	fmt.Println(e.data)
	s := MakeESlice[int64](3)
	fmt.Printf("%p\n", unsafe.Pointer(s.data))
	s = Append(s, 100, 200, 300)
	fmt.Printf("%p\n", unsafe.Pointer(s.data))
	fmt.Println(s.Get(1))
	fmt.Println(s.IsDangling())
	s = Append(s, 400)
	fmt.Println(s.IsDangling())
	s.Set(1, 500)
	fmt.Println(s.Get(1))
	fmt.Println(s.IsDangling())
	fmt.Println(s.Len())
	s = Append(s, 400, 500)
	s.Free()
	fmt.Println(s.IsDangling())

}
