package ddbpath

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// register our well-known paths
func init() {
	RegisterMessage(AnyPath{}, map[string]FieldInfo{
		"1": {Kind: BasicKind},
		"2": {Kind: AnyKind},
	})
	RegisterMessage(ValuePath{}, map[string]FieldInfo{})
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
