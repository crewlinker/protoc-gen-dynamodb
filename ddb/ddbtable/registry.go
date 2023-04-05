// Package ddbtable allows generated code to register DynamoDB table structure
package ddbtable

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"google.golang.org/protobuf/proto"
)

// defaultRegistry holds the registry in which table placement is registered by default
var defaultRegistry = NewRegistry()

// Registry allows for registering table placment information from Protobuf options
// on messages for wich Dynamo marshalling is being generated.
type Registry struct{}

// NewRegistry inits an empty registry
func NewRegistry() *Registry {
	return &Registry{}
}

// KeyPlacement describes the key's attribute name and type
type KeyPlacement struct {
	AttrName string
	AttrType expression.DynamoDBAttributeType
}

// GlobalSecondaryIndexPlacement describes a message being place in table's GSI
type GlobalSecondaryIndexPlacement struct {
	// Partition key of the gsi
	PartitionKey KeyPlacement
	// SortKey for the gsi
	SortKey KeyPlacement
	// Other attr names that are project in the GSI
	OtherAttrNames []string
}

// LocalSecondaryIndexPlacement describes a message being place in table's GSI
type LocalSecondaryIndexPlacement struct {
	// Sort key of this LSI
	SortKey KeyPlacement
	// Other attributes that are projected on the index
	OtherAttrNames []string
}

// MessagePlacement describes how a Protobuf message is placed in one-or-more tables with
// fields being part of any GSIs or LSIs
type MessagePlacement struct {
	// All tables this message may be placed in
	TableNames []string
	// Name and type of the Attribute that should be considered the base table partition key
	PartitionKey KeyPlacement
	// Name and type of the attribute that should be considered the base tables' sort key
	SortKey KeyPlacement
	// Global Secondary Indexes the item is placed in
	GlobalSecondaryIdxs map[string]GlobalSecondaryIndexPlacement
	// Local Secondary Indexes the item is placed in
	LocalSecondaryIdxs map[string]LocalSecondaryIndexPlacement
}

// Register registers the message's placement into a DynamoDB table. It validates the information in the
// placement and returns an error if it violtes the constrains given other messages that maybe also be
// placed in the table.
func (r *Registry) Register(typ proto.Message, mp MessagePlacement) error {

	// @TODO each messages M that is placed into a table, that have field being part of a LSI/GSI
	// if a field "F" is marked as the pk/sk of the lsi, all other messages that have the same attr
	// name must mark it in the same way.
	// Or, put another way. If one message defines a field to be the sk/pk of an index, all other
	// messages must do the same
	// @TODO same with pk/sk, if one attribute encodes as pk/sk under "sk/pk" all other messages

	return nil
}

// Register on the default registry
func Register(typ proto.Message, mp MessagePlacement) error {
	return defaultRegistry.Register(typ, mp)
}

// MustRegister will register the message's placement or panics
func MustRegister(typ proto.Message, mp MessagePlacement) {
	if err := Register(typ, mp); err != nil {
		panic("failed to register: " + err.Error())
	}
}
