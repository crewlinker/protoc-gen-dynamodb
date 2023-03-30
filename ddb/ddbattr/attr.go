// Package ddbattr provides helper code for building type-safe Dynamo paths
package ddbattr

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// List of basic type(s)
type List struct{ expression.NameBuilder }

// Index into a list of basic types
func (a List) Index(i int) expression.NameBuilder {
	return a.AppendName(expression.Name(fmt.Sprintf(`[%d]`, i)))
}

// ItemList is a list of nested items
type ItemList[T interface {
	WithDynamoNameBuilder(expression.NameBuilder) T
}] struct{ expression.NameBuilder }

// Index into a list of items
func (a ItemList[T]) Index(i int) T {
	var v T
	return v.WithDynamoNameBuilder(a.AppendName(expression.Name(fmt.Sprintf(`[%d]`, i))))
}

// Map of basic type(s)
type Map struct{ expression.NameBuilder }

// Key into a map of basic types
func (a Map) Key(k string) expression.NameBuilder {
	return a.AppendName(expression.Name(k))
}

// ItemMap is a list of nested items
type ItemMap[T interface {
	WithDynamoNameBuilder(expression.NameBuilder) T
}] struct{ expression.NameBuilder }

// Key into a list of items
func (a ItemMap[T]) Key(k string) T {
	var v T
	return v.WithDynamoNameBuilder(a.AppendName(expression.Name(k)))
}
