package ddb

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// P is the path to a basic field
type P struct{ v string }

// N returns the path as a namebuilder for direct use in expressions
func (p P) N() expression.NameBuilder {
	return expression.Name(p.String())
}

// String formats  the path correctly
func (p P) String() string { return strings.TrimPrefix(p.v, ".") }

// Set set the path value
func (p P) Set(v string) P { p.v = v; return p }

// BasicListP is the path to a list of basic types
type BasicListP struct{ v string }

// String formats the path and returns it
func (p BasicListP) String() string { return strings.TrimPrefix(p.v, ".") }

// N returns the path as a namebuilder for direct use in expressions
func (p BasicListP) N() expression.NameBuilder {
	return expression.Name(p.String())
}

// Set sets the path value
func (p BasicListP) Set(v string) BasicListP { p.v = v; return p }

// At returns the path to the basic type at the provided index
func (p BasicListP) At(i int) P {
	return P{p.v + "[" + strconv.Itoa(i) + "]"}
}

// ListP is the path to a list of messages
type ListP[T interface{ Set(v string) T }] struct{ v string }

// N returns the path as a namebuilder for direct use in expressions
func (p ListP[T]) N() expression.NameBuilder {
	return expression.Name(p.String())
}

// String formats the path and returns it
func (p ListP[T]) String() string { return strings.TrimPrefix(p.v, ".") }

// Set sets the path value
func (p ListP[T]) Set(v string) ListP[T] { p.v = v; return p }

// At returns the path to the basic type at the provided index
func (p ListP[T]) At(i int) T {
	var v T
	return v.Set(p.v + "[" + strconv.Itoa(i) + "]")
}
