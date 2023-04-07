// Package ddbtable allows generated code to register DynamoDB table structure
package ddbtable

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"google.golang.org/protobuf/proto"
)

// defaultRegistry holds the registry in which table placement is registered by default
var defaultRegistry = NewRegistry()

// Registry allows for registering table placment information from Protobuf options
// on messages for wich Dynamo marshalling is being generated.
type Registry struct {
	tables map[string]*TablePlacement
}

// NewRegistry inits an empty registry
func NewRegistry() *Registry {
	return &Registry{tables: map[string]*TablePlacement{}}
}

// KeyDescriptor describes the key's attribute name and type
type KeyDescriptor struct {
	// Name of the attribute that is marked as the key
	AttrName string
	// AttrType is the kind of attribute value
	AttrType expression.DynamoDBAttributeType
}

// CheckCompatibility checks if the ye descriptor of 'd' is compatible with 'other'
func (d *KeyDescriptor) CheckCompatibility(other *KeyDescriptor) error {
	if (d == nil && other != nil) || (d != nil && other == nil) {
		return fmt.Errorf("both must be defined, or both must not be defined")
	}

	if d.AttrName != other.AttrName {
		return fmt.Errorf("attribute name is different: %s != %s", d.AttrName, other.AttrName)
	}
	if d.AttrType != other.AttrType {
		return fmt.Errorf("attribute type is different: %s != %s", d.AttrType, other.AttrType)
	}
	return nil
}

// GlobalSecondaryIndexPlacement describes a message being place in table's GSI
type GlobalSecondaryIndexPlacement struct {
	// Partition key of the gsi
	PartitionKey KeyDescriptor
	// SortKey for the gsi
	SortKey *KeyDescriptor
	// Other attr names that are project in the GSI
	OtherAttrNames []string
}

// CheckCompatibility will error if 'p' is not comability with the table placement constraint of 'other'
func (p *GlobalSecondaryIndexPlacement) CheckCompatibility(other *GlobalSecondaryIndexPlacement) (err error) {
	return
}

// LocalSecondaryIndexPlacement describes a message being place in table's GSI
type LocalSecondaryIndexPlacement struct {
	// Sort key of this LSI
	SortKey KeyDescriptor
	// Other attributes that are projected on the index
	OtherAttrNames []string
}

// CheckCompatibility will error if 'p' is not comability with the table placement constraint of 'other'
func (p *LocalSecondaryIndexPlacement) CheckCompatibility(other *LocalSecondaryIndexPlacement) (err error) {
	return
}

// TablePlacement describes how a Protobuf message is placed in one-or-more tables with
// fields being part of any GSIs or LSIs
type TablePlacement struct {
	// All tables this message may be placed in
	TableNames []string
	// Name and type of the Attribute that should be considered the base table partition key
	PartitionKey KeyDescriptor
	// Name and type of the attribute that should be considered the base tables' sort key
	SortKey *KeyDescriptor
	// Global Secondary Indexes the item is placed in
	GlobalSecondaryIdxs map[string]*GlobalSecondaryIndexPlacement
	// Local Secondary Indexes the item is placed in
	LocalSecondaryIdxs map[string]*LocalSecondaryIndexPlacement
}

// CheckCompatibility will error if 'p' is not comability with the table placement constraint of 'other'
func (p *TablePlacement) CheckCompatibility(other *TablePlacement) (err error) {
	if err = p.PartitionKey.CheckCompatibility(&other.PartitionKey); err != nil {
		return fmt.Errorf("partition key %s is not compatible with %s: %w", p.PartitionKey, other.PartitionKey, err)
	}

	if err = p.SortKey.CheckCompatibility(other.SortKey); err != nil {
		return fmt.Errorf("sort key %s is not compatible with %s: %w", p.SortKey, other.SortKey, err)
	}

	return nil
}

// TableCreates returns create table inputs from registered tables.
func (r *Registry) TableCreates() []*dynamodb.CreateTableInput {
	// (key)attribute definitions
	// key schema
	// GlobalSecondaryIndexes
	// LocalSecondaryIndexes
	return nil
}

// Register registers the message's placement into a DynamoDB table. It validates the information in the
// placement and returns an error if it violtes the constrains given other messages that maybe also be
// placed in the table.
func (r *Registry) Register(m proto.Message, p *TablePlacement) (err error) {
	for _, tname := range p.TableNames {
		existing, ok := r.tables[tname]
		if !ok {
			// first placement for the table, nothing to check, just add to registry
			r.tables[tname] = p
			continue
		}

		if err = p.CheckCompatibility(existing); err != nil {
			return fmt.Errorf("table '%s': %w", tname, err)
		}

		// check and addd lsi placements
		for lsiName, lsip := range p.LocalSecondaryIdxs {
			existingLsi, ok := existing.LocalSecondaryIdxs[lsiName]
			if !ok {
				// first placmenent for the lsi, nothing to check against
				existing.LocalSecondaryIdxs[lsiName] = lsip
				continue
			}

			if err = lsip.CheckCompatibility(existingLsi); err != nil {
				return fmt.Errorf("table '%s', lsi '%s': %w", tname, lsiName, err)
			}

			// @TODO update the existing lsi, with projections(?)
			existing.LocalSecondaryIdxs[lsiName] = existingLsi
		}

		// check and add gsi placement
		for gsiName, gsip := range p.GlobalSecondaryIdxs {
			existingGsi, ok := existing.GlobalSecondaryIdxs[gsiName]
			if !ok {
				// first placmenent for the gsi, nothing to check, just add to registry
				existing.GlobalSecondaryIdxs[gsiName] = gsip
				continue
			}

			if err = gsip.CheckCompatibility(existingGsi); err != nil {
				return fmt.Errorf("table '%s', gsi '%s': %w", tname, gsiName, err)
			}
			// @TODO add something? to a map? update projections?

			existing.GlobalSecondaryIdxs[gsiName] = existingGsi
		}

		// check with existing:
		//  - PK has different attribute name
		//  - PK has

		// checks if other messages have conflicting placement
		// - pk field with different name
		// - sk field with different name
		// - messages in the same gsi with wrong sk pk names
		// - mesasge als in the lsi with wrong sk name
		_ = existing
		r.tables[tname] = existing
	}

	// @TODO each messages M that is placed into a table, that have field being part of a LSI/GSI
	// if a field "F" is marked as the pk/sk of the lsi, all other messages that have the same attr
	// name must mark it in the same way.
	// Or, put another way. If one message defines a field to be the sk/pk of an index, all other
	// messages must do the same
	// @TODO same with pk/sk, if one attribute encodes as pk/sk under "sk/pk" all other messages
	// @TODO error when one message of a table just has a PK, while others have a SK as well
	// @TODO error if one message has a different AttributeValueType for pk/sk (base table and indices)

	return nil
}

// Register on the default registry
func Register(typ proto.Message, mp *TablePlacement) error {
	return defaultRegistry.Register(typ, mp)
}

// MustRegister will register the message's placement or panics
func MustRegister(typ proto.Message, mp *TablePlacement) {
	if err := Register(typ, mp); err != nil {
		panic("failed to register: " + err.Error())
	}
}
