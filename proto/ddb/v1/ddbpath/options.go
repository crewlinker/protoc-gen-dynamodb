// Code generated by protoc-gen-dynamodb. DO NOT EDIT.

// Package ddbv1ddbpath holds generated code for working with Dynamo document paths
package ddbv1ddbpath

import (
	expression "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	ddbpath "github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
)

// FieldOptionsPath allows for constructing type-safe expression names
type FieldOptionsPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p FieldOptionsPath) WithDynamoNameBuilder(n expression.NameBuilder) FieldOptionsPath {
	p.NameBuilder = n
	return p
}

// Name appends the path being build
func (p FieldOptionsPath) Name() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// Pk appends the path being build
func (p FieldOptionsPath) Pk() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// Sk appends the path being build
func (p FieldOptionsPath) Sk() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}

// Omit appends the path being build
func (p FieldOptionsPath) Omit() expression.NameBuilder {
	return p.AppendName(expression.Name("4"))
}

// Set appends the path being build
func (p FieldOptionsPath) Set() expression.NameBuilder {
	return p.AppendName(expression.Name("5"))
}

// Embed appends the path being build
func (p FieldOptionsPath) Embed() expression.NameBuilder {
	return p.AppendName(expression.Name("6"))
}
func init() {
	ddbpath.Register(FieldOptionsPath{}, map[string]ddbpath.FieldInfo{
		"1": {Kind: ddbpath.FieldKindSingle},
		"2": {Kind: ddbpath.FieldKindSingle},
		"3": {Kind: ddbpath.FieldKindSingle},
		"4": {Kind: ddbpath.FieldKindSingle},
		"5": {Kind: ddbpath.FieldKindSingle},
		"6": {Kind: ddbpath.FieldKindSingle},
	})
}
