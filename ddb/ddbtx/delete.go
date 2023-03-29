package ddbtx

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
)

// DeleteOption allows configuring delete calls
type DeleteOption func(*deleteOpts)

// DeleteIfExists adds a condition check that checks if the item exists
func DeleteIfExists() DeleteOption {
	return func(po *deleteOpts) {
		po.ifExists = true
	}
}

// Delete adds a DeleteItem to the transaction
func (wtx WriteTx) Delete(it Item, os ...DeleteOption) WriteTx {
	var del types.Delete
	var exprb expression.Builder

	opts := applyDeleteOptions(os...)
	avm, err := it.MarshalDynamoItem()
	if err != nil {
		return wtx.errorf("failed to marshal: %w", err)
	}

	del.Key, err = ddbpath.SelectMapValues(avm, it.DynamoKeyNames()...)
	if err != nil {
		return wtx.errorf("failed to select key attribute values: %w", err)
	}

	if opts.ifExists {
		exprb = withShouldExist(exprb, it)
	}

	expr, err := exprb.Build()
	if err != nil {
		return wtx.errorf("failed to build expression: %w", err)
	}

	del.ExpressionAttributeNames = expr.Names()
	del.ExpressionAttributeValues = expr.Values()
	del.ConditionExpression = expr.Condition()

	wtx.its = append(wtx.its, types.TransactWriteItem{Delete: &del})
	return wtx
}

// deleteOptions holds optional configurations for delete
type deleteOpts struct {
	ifExists bool
}

// applyDeleteOptions applies delete option defaults and overwites whatever is configured
func applyDeleteOptions(pos ...DeleteOption) (po deleteOpts) {
	for _, dof := range pos {
		dof(&po)
	}
	return
}
