package ddb

import "strconv"

// BasicListPath path is a type for building a path into a list of basic types
type BasicListPath string

// At append a list index to the path and returns the whole path
func (p BasicListPath) At(i int) string {
	return string(p) + "[" + strconv.Itoa(i) + "]"
}

// ListPath is a type for building paths into a list of messages
type ListPath[T interface{ Set(v string) T }] string

// At appends a list index to the path and returns T
func (p ListPath[T]) At(i int) T {
	var v T
	return v.Set(string(p) + "[" + strconv.Itoa(i) + "]")
}
