package ddbpath

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// Kind describes the registered message field for validation
type Kind int

// String returns a human readable version
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
	// BasicKind are non-composit fields
	BasicKind Kind = 0
	// MapKind describe map fields
	MapKind Kind = 1
	// ListKind describe list fields
	ListKind Kind = 2
	// AnyKind describes a field that may hold any value
	AnyKind Kind = 3
)

// FieldInfo describes information about a message field for validation
type FieldInfo struct {
	Kind Kind
	Ref  reflect.Type
}

var registry = map[reflect.Type]map[string]FieldInfo{}

// RegisterMessage registers a message's path struct for validation
func RegisterMessage(m interface {
	AppendName(field expression.NameBuilder) expression.NameBuilder
}, fields map[string]FieldInfo) {
	registry[reflect.TypeOf(m)] = fields
}

// Validate paths agains a registered type that implements name builder. Types from the generated ddbpath
// package will be automatically registered.
func Validate(nb interface {
	AppendName(field expression.NameBuilder) expression.NameBuilder
}, paths ...string) (err error) {
	typ := reflect.TypeOf(nb)
	els := make([]PathElement, 32) // allocate space for upto 32 element deep paths, this Dynamo's max
	for _, p := range paths {
		if els, err = AppendParsePath(p, els[:0]); err != nil {
			return fmt.Errorf("failed to parse path '%s': %w", p, err)
		}

		if err = validatePath(typ, els); err != nil {
			return fmt.Errorf("failed to validate path '%s': %w", p, err)
		}
	}
	return
}

// validatePath validates a single set of parsed paths against registered path types
func validatePath(typ reflect.Type, els []PathElement) (err error) {
	var field string
	var index int

	var currField FieldInfo
	var i int
	for ; i < len(els) && typ != nil; i++ {
		field, index = els[i].Field, els[i].Index
		fields := registry[typ]
		if fields == nil {
			return fmt.Errorf("type not registered: %v", typ)
		}

		// the ValuePath holds arbitrary json so it is the exception, any
		// path into it is valid
		if typ == reflect.TypeOf(ValuePath{}) {
			return nil
		}

		if index >= 0 {

			// either did not enter a "list" field, or it is not of the list kind
			if currField.Kind != ListKind {
				return fmt.Errorf("index '[%d]' into non-list: %s", index, currField.Kind)
			}

			typ = currField.Ref
			currField = FieldInfo{}
		} else {

			// if we entered a map field, any field value is valid
			if currField.Kind == MapKind {
				typ = currField.Ref
				currField = FieldInfo{}
				continue
			}

			// if not, the typ must have the field defined
			var isValid bool
			currField, isValid = fields[field]
			switch {
			case !isValid:
				return fmt.Errorf("non-existing field '%s' on: %v", field, typ)
			case currField.Ref != nil:
				// field holds another type that we "move into"
				typ = currField.Ref
			default:
				// field holds a basic type, can't traver further
				typ = nil
			}
		}
	}

	// if we did not iterate until the last element it means the path is "too deep"
	// and this is also invalid.
	if i < len(els) && currField.Kind != AnyKind {
		return fmt.Errorf("path (or index) '%s' on basic type", els[i].Field)
	}

	return nil
}
