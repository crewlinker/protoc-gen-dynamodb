package ddbtable

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CreateTableOption allows configuration of how the table will be created
type CreateTableOption func(*createOpts)

// WithProvisioning allows customization of how the tables read and write capacity are provisioned. If the
// gsiName is not empty it can be specified for global secondary indexes.
func WithProvisioning(pf func(gsiName string) (types.BillingMode, *types.ProvisionedThroughput)) CreateTableOption {
	return func(co *createOpts) {
		co.provisionFunc = pf
	}
}

// WithTableName allows explicitely defining the table name instead of taking it directly from the
// protobuf definition.
func WithTableName(name string) CreateTableOption {
	return func(co *createOpts) {
		co.tableName = name
	}
}

// TableCreate will formulate a table creation input for the v2 sdk. It return nil if no table
// with the given name is registered.
func (r *Registry) TableCreate(name string, opts ...CreateTableOption) (cti *dynamodb.CreateTableInput) {
	def, ok := r.TableDef(name)
	if !ok {
		return nil
	}

	o := applyCreateOpts(opts...)
	if o.tableName == "" {
		o.tableName = def.Name
	}

	cti = &dynamodb.CreateTableInput{TableName: &o.tableName}
	cti.BillingMode, cti.ProvisionedThroughput = o.provisionFunc("")
	cti.KeySchema, cti.AttributeDefinitions = withKeyAttr(cti.KeySchema, cti.AttributeDefinitions,
		def.PartitionKey, types.KeyTypeHash)
	if def.SortKey != nil {
		cti.KeySchema, cti.AttributeDefinitions = withKeyAttr(cti.KeySchema, cti.AttributeDefinitions,
			def.SortKey, types.KeyTypeRange)
	}

	for _, gsid := range def.GlobalIndexes {
		gsi := types.GlobalSecondaryIndex{
			IndexName: aws.String(gsid.Name)}
		gsi.Projection = &types.Projection{ProjectionType: types.ProjectionTypeAll}
		gsi.KeySchema, cti.AttributeDefinitions = withKeyAttr(gsi.KeySchema, cti.AttributeDefinitions,
			gsid.PartitionKey, types.KeyTypeHash)
		if gsid.SortKey != nil {
			gsi.KeySchema, cti.AttributeDefinitions = withKeyAttr(gsi.KeySchema, cti.AttributeDefinitions,
				gsid.SortKey, types.KeyTypeRange)
		}
		_, gsi.ProvisionedThroughput = o.provisionFunc(gsid.Name)
		cti.GlobalSecondaryIndexes = append(cti.GlobalSecondaryIndexes, gsi)
	}

	return
}

// createOpts hold options for fromulating a table create input
type createOpts struct {
	provisionFunc func(gsiName string) (types.BillingMode, *types.ProvisionedThroughput)
	tableName     string
}

// applyCreateOpts applies the options with defaults
func applyCreateOpts(os ...CreateTableOption) (opts createOpts) {
	os = append([]CreateTableOption{
		WithProvisioning(func(gsiName string) (types.BillingMode, *types.ProvisionedThroughput) {
			return types.BillingModeProvisioned, &types.ProvisionedThroughput{
				ReadCapacityUnits: aws.Int64(1), WriteCapacityUnits: aws.Int64(1),
			}
		}),
	}, os...)

	for _, of := range os {
		of(&opts)
	}
	return
}
