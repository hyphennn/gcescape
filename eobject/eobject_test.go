// Package eobject
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package eobject

import (
	"fmt"
	"testing"
)

func TestEObject(t *testing.T) {
	e := MakeEObject[int]()
	fmt.Println(e)
	fmt.Println(e.Value())
	e.Set(3)
	fmt.Println(e.Value())
	e.Set(4)
	fmt.Println(e.Value())
	e.Free()
	fmt.Println(e)
}
