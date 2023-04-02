package ddbpath

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// register our well-known paths
func init() {
	Register(ValuePath{}, map[string]FieldInfo{})
	Register(AnyPath{}, map[string]FieldInfo{
		"1": {Kind: FieldKindSingle},
		"2": {Kind: FieldKindSingle, Message: reflect.TypeOf(ValuePath{})},
	})
	Register(FieldMaskPath{}, map[string]FieldInfo{
		"1": {Kind: FieldKindList},
	})
}

// AnyPath is registered to support path validation into anypb structs
type AnyPath struct{ expression.NameBuilder }

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p AnyPath) WithDynamoNameBuilder(n expression.NameBuilder) AnyPath {
	p.NameBuilder = n
	return p
}

// TypeURL appends the path of the type url
func (p AnyPath) TypeURL() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// Value appends the path of the value
func (p AnyPath) Value() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// ValuePath is registered to support path validation into structpb's value fields. It has no
// fields but is special in that it will accept any path into it.
type ValuePath struct{ expression.NameBuilder }

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p ValuePath) WithDynamoNameBuilder(n expression.NameBuilder) ValuePath {
	p.NameBuilder = n
	return p
}

// FieldMaskPath is registered to support path validation of fieldmask
type FieldMaskPath struct{ expression.NameBuilder }

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p FieldMaskPath) WithDynamoNameBuilder(n expression.NameBuilder) FieldMaskPath {
	p.NameBuilder = n
	return p
}

// Masks appends the path of the value
func (p FieldMaskPath) Masks() List {
	return List{NameBuilder: p.AppendName(expression.Name("1"))}
}
