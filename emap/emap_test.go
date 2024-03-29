// Package emap
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/29
package emap_test

import (
	"fmt"
	"testing"

	"github.com/hyphennn/gcescape/emap"
)

func TestEMap(t *testing.T) {
	e := emap.MakeEMap[int, int](3)
	fmt.Println(e.Get(1))
	e.Set(1, 2)
	fmt.Println(e.Get(1))
	e.Set(2, 3)
	e.Set(3, 3)
	e.Set(4, 3)
	e.Set(5, 3)
	e.Set(5, 3)
	e.Set(6, 3)
	fmt.Println(e.Get(4))
	fmt.Println(e.Delete(4))
	fmt.Println(e.Delete(5))
	fmt.Println(e.Get(4))
	e.Set(7, 6)
	fmt.Println(e.Get(6))
}
