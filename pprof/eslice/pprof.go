// Package main
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()
	time.Sleep(time.Second)
	UseESlice()
	//UseNormalSlice()
	//UseNormalPtrSlice()
	time.Sleep(1000 * time.Second)
}
