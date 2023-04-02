package ddbreg

import (
	"fmt"
	"reflect"
)

// ErrTypeNotRegistered is returned when a type is not registered in the registry being used
type ErrTypeNotRegistered struct{ typ reflect.Type }

func (e ErrTypeNotRegistered) Error() string { return fmt.Sprintf("type not registered: %s", e.typ) }

func errTypeNotRegistered(typ reflect.Type) error {
	return fmt.Errorf("%w", ErrTypeNotRegistered{typ})
}

type ErrIndexNotAllowed struct {
	idx  int
	info FieldInfo
}

func (e ErrIndexNotAllowed) Error() string {
	return fmt.Sprintf("indexing '%d' not allowed on %s", e.idx, e.info)
}
func errIndexNotAllowed(idx int, info FieldInfo) error {
	return fmt.Errorf("%w", ErrIndexNotAllowed{idx, info})
}

type ErrFieldNotAllowed struct {
	field string
	info  FieldInfo
}

func (e ErrFieldNotAllowed) Error() string {
	return fmt.Sprintf("field selecting '%s' not allowed on %s", e.field, e.info)
}
func errFieldNotAllowed(field string, info FieldInfo) error {
	return fmt.Errorf("%w", ErrFieldNotAllowed{field, info})
}

type ErrUnknownField struct {
	field string
	info  FieldInfo
}

func (e ErrUnknownField) Error() string {
	return fmt.Sprintf("unknown field '%s' of %s", e.field, e.info)
}
func errUnknownField(field string, info FieldInfo) error {
	return fmt.Errorf("%w", ErrUnknownField{field, info})
}
