package ddbpath

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

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
	if i < len(els) {
		return fmt.Errorf("path (or index) '%s' on basic type", els[i].Field)
	}

	return nil
}
