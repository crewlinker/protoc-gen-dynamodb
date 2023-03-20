package ddb

import (
	"fmt"
	"strings"
)

// Path into a Dynamo protobuf message
type Path struct{ val string }

// SetTo allows 'p' to be set to the values of 'p2' by external implementations
// for generic types since it doesn't allow access to shared fields.
func (p *Path) SetTo(p2 Path) {
	p.val = p2.val
}

// String formats the path as a string
func (p Path) String() string {
	return p.val
}

// Append another element to the path, it modifies 'p' and returns it.
func (p *Path) Append(e string) Path {
	switch {
	case strings.HasPrefix(e, "[") && strings.HasSuffix(e, "]"):
		p.val += e
	default:
		p.val += "." + e
	}

	return *p
}

// ListPath provides the ability to append an index path element
// for list members.
type ListPath[T any, TP interface {
	SetTo(Path)
	*T
}] struct {
	Path
}

func (p ListPath[T, TP]) Index(i int) TP {
	vp := TP(new(T)) // init pointer to
	vp.SetTo(p.Append(fmt.Sprintf("[%d]", i)))
	return vp
}
