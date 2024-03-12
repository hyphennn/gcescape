// Package internal
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/12
package internal

import (
	"syscall"
	"fmt"
)

type SysCallErr struct {
	Errno syscall.Errno
	Msg   string
}

func (s *SysCallErr) Error() string {
	return fmt.Sprintf("sys call failed: errno: %d, info: %s", s.Errno, s.Msg)
}

func NewSysCallErr(e syscall.Errno, m string) *SysCallErr {
	return &SysCallErr{e, m}
}
