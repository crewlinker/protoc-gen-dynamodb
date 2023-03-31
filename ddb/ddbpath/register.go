package ddbpath

import "reflect"

type Kind int

func (k Kind) String() string {
	switch k {
	case BasicKind:
		return "Basic"
	case MapKind:
		return "Map"
	case ListKind:
		return "List"
	default:
		panic("unsupported")
	}
}

const (
	BasicKind Kind = 0
	MapKind   Kind = 1
	ListKind  Kind = 2
)

type FieldInfo struct {
	Kind Kind
	Ref  reflect.Type
}

var registry = map[reflect.Type]map[string]FieldInfo{}

func RegisterMessage(name reflect.Type, fields map[string]FieldInfo) {
	registry[name] = fields
}
