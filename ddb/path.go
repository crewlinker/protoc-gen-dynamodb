package ddb

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// P is the path to a basic field
type P struct{ v string }

// N returns the path as a namebuilder for direct use in expressions
func (p P) N() expression.NameBuilder { return expression.Name(p.String()) }

// K returns the path as a keybuilder for direct use in expressions
func (p P) K() expression.KeyBuilder { return expression.Key(p.String()) }

// String formats  the path correctly
func (p P) String() string { return strings.TrimPrefix(p.v, ".") }

// Val returns the raw value of the path, without formatting it
func (p P) Val() string { return p.v }

// Set set the path value
func (p P) Set(v string) P { p.v = v; return p }

// BasicListP is the path to a list of basic types
type BasicListP struct{ P }

// Set sets the path value
func (p BasicListP) Set(v string) BasicListP { p.P = p.P.Set(v); return p }

// At returns the path to the basic type at the provided index
func (p BasicListP) At(i int) P {
	return P{p.Val() + "[" + strconv.Itoa(i) + "]"}
}

// ListP is the path to a list of messages
type ListP[T interface{ Set(v string) T }] struct{ P }

// Set sets the path value
func (p ListP[T]) Set(v string) ListP[T] { p.P = p.P.Set(v); ; return p }

// At returns the path to the basic type at the provided index
func (p ListP[T]) At(i int) T {
	var v T
	return v.Set(p.Val() + "[" + strconv.Itoa(i) + "]")
}

// BasicMapP is the path to a list of basic types
type BasicMapP struct{ P }

// Set sets the path value
func (p BasicMapP) Set(v string) BasicMapP { p.P = p.P.Set(v); ; return p }

// Key returns the path to the basic type at the provided index
func (p BasicMapP) Key(k string) P {
	return P{p.Val() + "." + k}
}

// MapP is the path to a map of messages
type MapP[T interface{ Set(v string) T }] struct{ P }

// Set sets the path value
func (p MapP[T]) Set(v string) MapP[T] { p.P = p.P.Set(v); return p }

// Key returns the path to the basic type at the provided index
func (p MapP[T]) Key(k string) T {
	var v T
	return v.Set(p.Val() + "." + k)
}
