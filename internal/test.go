// Package internal
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/4/2
package internal

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

func GenTestType() *TestType {
	return &TestType{
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
	}
}
