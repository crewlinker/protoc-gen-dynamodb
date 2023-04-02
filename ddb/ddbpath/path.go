// Package ddbpath provides logic for building and parsing Dynamo document paths
package ddbpath

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// List of basic type(s)
type List struct{ expression.NameBuilder }

// Index into a list of basic types
func (p List) Index(i int) expression.NameBuilder {
	return p.AppendName(expression.Name(fmt.Sprintf(`[%d]`, i)))
}

// ItemList is a list of nested items
type ItemList[T interface {
	WithDynamoNameBuilder(expression.NameBuilder) T
}] struct{ expression.NameBuilder }

// Index into a list of items
func (p ItemList[T]) Index(i int) T {
	var v T
	return v.WithDynamoNameBuilder(p.AppendName(expression.Name(fmt.Sprintf(`[%d]`, i))))
}

// Map of basic type(s)
type Map struct{ expression.NameBuilder }

// Key into a map of basic types
func (p Map) Key(k string) expression.NameBuilder {
	return p.AppendName(expression.Name(k))
}

// ItemMap is a list of nested items
type ItemMap[T interface {
	WithDynamoNameBuilder(expression.NameBuilder) T
}] struct{ expression.NameBuilder }

// Key into a list of items
func (p ItemMap[T]) Key(k string) T {
	var v T
	return v.WithDynamoNameBuilder(p.AppendName(expression.Name(k)))
}

// register list and map
func init() {
	Register(List{}, map[string]FieldInfo{})
	Register(Map{}, map[string]FieldInfo{})
}
