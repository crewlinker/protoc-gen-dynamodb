package ddbpath

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// register our well-known paths
func init() {
	RegisterMessage(reflect.TypeOf(AnyPath{}), map[string]FieldInfo{
		"1": {Kind: BasicKind},
		"2": {Kind: BasicKind}, // @TODO should be a maptype that allows for any json path
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