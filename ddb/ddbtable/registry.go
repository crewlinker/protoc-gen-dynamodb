// Package ddbtable allows generated code to register table structure
package ddbtable

import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

// GlobalIndex describes a Global Secondary Index (GSI) on a DynamoDB table
type GlobalIndex struct {
	Name         string
	PartitionKey *Attribute
	SortKey      *Attribute
}

// Attribute describes attribute by name and type (String, Bool, Number, etc)
type Attribute struct {
	Name string
	Type expression.DynamoDBAttributeType
}

// Table describes the structure of a DynamoDB table
type Table struct {
	Name          string
	PartitionKey  *Attribute
	SortKey       *Attribute
	EntityType    *Attribute
	GlobalIndexes []*GlobalIndex
}

// Registry holds table descriptions
type Registry struct {
	tables map[string]*Table
}

// NewRegistry inits an empty registry
func NewRegistry() *Registry {
	return &Registry{map[string]*Table{}}
}

// TableDef returns the table definition (if any) of a table with the given name.
func (r *Registry) TableDef(name string) (tbl *Table, ok bool) {
	tbl, ok = r.tables[name]
	return
}

// Register the table definition. If a table with the same name is
// already registered, it is overwritten.
func (r *Registry) Register(tbl *Table) {
	r.tables[tbl.Name] = tbl
}

// TableDef returs the table definition (if any) from the defaulte registry
func TableDef(name string) (tbl *Table, ok bool) {
	return defaultRegistry.TableDef(name)
}

// Register a table description with the default registry
func Register(tbl *Table) {
	defaultRegistry.Register(tbl)
}

// default registry
var defaultRegistry = NewRegistry()
