// Code generated by protoc-gen-dynamodb. DO NOT EDIT.

// Package messagev1ddb holds generated schema structure
package messagev1ddb

import expression "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

// OtherKitchenPath allows for constructing type-safe expression names
type OtherKitchenPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p OtherKitchenPath) WithDynamoNameBuilder(n expression.NameBuilder) OtherKitchenPath {
	p.NameBuilder = n
	return p
}

// AnotherKitchen returns 'p' with the attribute name appended and allow subselecting nested message
func (p OtherKitchenPath) AnotherKitchen() KitchenPath {
	return KitchenPath{p.AppendName(expression.Name("16"))}
}

// OtherTimer appends the path being build
func (p OtherKitchenPath) OtherTimer() expression.NameBuilder {
	return p.AppendName(expression.Name("17"))
}
