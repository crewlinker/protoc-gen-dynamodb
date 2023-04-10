// Code generated by protoc-gen-dynamodb. DO NOT EDIT.

// Package modelv2ddbpath holds generated code for working with Dynamo document paths
package modelv2ddbpath

import (
	expression "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	ddbpath "github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
	"reflect"
)

// RoomPath allows for constructing type-safe expression names
type RoomPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p RoomPath) WithDynamoNameBuilder(n expression.NameBuilder) RoomPath {
	p.NameBuilder = n
	return p
}

// Number appends the path being build
func (p RoomPath) Number() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}
func init() {
	ddbpath.Register(RoomPath{}, map[string]ddbpath.FieldInfo{"1": {Kind: ddbpath.FieldKindSingle}})
}

// StrangeKeysPath allows for constructing type-safe expression names
type StrangeKeysPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p StrangeKeysPath) WithDynamoNameBuilder(n expression.NameBuilder) StrangeKeysPath {
	p.NameBuilder = n
	return p
}

// Hash appends the path being build
func (p StrangeKeysPath) Hash() expression.NameBuilder {
	return p.AppendName(expression.Name("23"))
}

// Range appends the path being build
func (p StrangeKeysPath) Range() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// Kind appends the path being build
func (p StrangeKeysPath) Kind() expression.NameBuilder {
	return p.AppendName(expression.Name("300"))
}

// Gsi1Pk appends the path being build
func (p StrangeKeysPath) Gsi1Pk() expression.NameBuilder {
	return p.AppendName(expression.Name("34"))
}

// Room returns 'p' with the attribute name appended and allow subselecting nested message
func (p StrangeKeysPath) Room() RoomPath {
	return RoomPath{NameBuilder: p.AppendName(expression.Name("400"))}
}
func init() {
	ddbpath.Register(StrangeKeysPath{}, map[string]ddbpath.FieldInfo{
		"1":   {Kind: ddbpath.FieldKindSingle},
		"23":  {Kind: ddbpath.FieldKindSingle},
		"300": {Kind: ddbpath.FieldKindSingle},
		"34":  {Kind: ddbpath.FieldKindSingle},
		"400": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(RoomPath{}),
		},
	})
}
