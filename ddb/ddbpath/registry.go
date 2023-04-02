package ddbpath

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// NameBuilder inteface is implemented by generated name building structs
type NameBuilder interface {
	AppendName(field expression.NameBuilder) expression.NameBuilder
}

// FieldKind describes the kind of field
type FieldKind int

const (
	// FieldKindUndefined means the kind was unknown
	FieldKindUndefined FieldKind = iota
	// FiledKindSingle represents a single instance of something
	FieldKindSingle
	// FieldKindList means a list of something
	FieldKindList
	// FieldKindMap means a map of something
	FieldKindMap
)

// String provides human readable form for the kind
func (fi FieldKind) String() string {
	switch fi {
	case FieldKindUndefined:
		return "_undefined"
	case FieldKindSingle:
		return "Single"
	case FieldKindList:
		return "List"
	case FieldKindMap:
		return "Map"
	default:
		panic("unsupported")
	}
}

// FieldInfo of a field on a message
type FieldInfo struct {
	Kind    FieldKind    // list, map, basic, any etc
	Message reflect.Type // field holds a non-basic type, or nil if its a basic type
}

// NoInfo is the FieldInfo zero value
var NoInfo = FieldInfo{}

// String returns a human readable form of the field info
func (fi FieldInfo) String() string {
	if fi.Message == nil {
		return fi.Kind.String()
	}
	return fmt.Sprintf("%s<%s>", fi.Kind, fi.Message)
}

// Registry holds type information so validation of string paths can happen efficiently
type Registry struct {
	infos map[reflect.Type]registryItem
}

// NewRegistry inits an empty registry
func NewRegistry() Registry {
	return Registry{infos: make(map[reflect.Type]registryItem)}
}

// registryItem represents a single registered message
type registryItem struct {
	fields map[string]FieldInfo
}

// fieldsOf a registered item
func (r Registry) fieldsOf(typ reflect.Type) (fi map[string]FieldInfo, ok bool) {
	it, ok := r.infos[typ]
	if !ok {
		return fi, false
	}
	return it.fields, true
}

// FieldsOf returns field information of a name builder implementation
func (r Registry) FieldsOf(nb NameBuilder) (fi map[string]FieldInfo, ok bool) {
	return r.fieldsOf(reflect.TypeOf(nb))
}

// Register a name builder with the registry for efficient validation. It panics if the typ is
// already registered.
func (r Registry) Register(nb NameBuilder, fields map[string]FieldInfo) {
	typ := reflect.TypeOf(nb)
	if _, ok := r.infos[typ]; ok {
		panic(fmt.Sprintf("ddbpath: type '%s' is already registered for validation", typ))
	}
	r.infos[typ] = registryItem{fields: fields}
}

// Traverse a message name builder 'nb' via path 'p' and return field info and any fields.
func (r Registry) Traverse(nb NameBuilder, p string) (fi FieldInfo, flds map[string]FieldInfo, err error) {
	typ := reflect.TypeOf(nb)
	els := make([]PathElement, 32)
	if els, err = AppendParsePath(p, els[:0]); err != nil {
		return fi, flds, fmt.Errorf("failed to parse path '%s': %w", p, err)
	}
	return r.traverse(typ, els)
}

// Validate given the types in the registry
func (r Registry) Validate(nb NameBuilder, paths ...string) (err error) {
	typ := reflect.TypeOf(nb)
	els := make([]PathElement, 32) // allocate space for upto 32 element deep paths, this Dynamo's max
	for _, p := range paths {
		if els, err = AppendParsePath(p, els[:0]); err != nil {
			return fmt.Errorf("failed to parse path '%s': %w", p, err)
		}

		if err = r.validate(typ, els); err != nil {
			return fmt.Errorf("failed to validate path '%s': %w", p, err)
		}
	}
	return
}

// defaultRegistry allows validation agains the default registry
var defaultRegistery = NewRegistry()

// Validate each path in 'paths' given the path naming struct 'nb' as the root
// message where the paths are started from agains the default registery.
func Validate(nb NameBuilder, paths ...string) error {
	return defaultRegistery.Validate(nb, paths...)
}

// Register a generated name building struct with the default registry. It panics if the
// type is already registered.
func Register(nb NameBuilder, fields map[string]FieldInfo) { defaultRegistery.Register(nb, fields) }
